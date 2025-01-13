package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zanda256/commitsmart-task/app/services/user-api/v1/handlers"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	v1 "github.com/Zanda256/commitsmart-task/business/web/v1"
	"github.com/Zanda256/commitsmart-task/foundation/keystore"
	"github.com/Zanda256/commitsmart-task/foundation/logger"
	"github.com/Zanda256/commitsmart-task/foundation/web"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//mongoScheme = "mongodb+srv://"
	mongoScheme = "mongodb://"
)

var (
	build = "develop"
)

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT ******")
		},
	}

	traceIDFunc := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "USER-API", traceIDFunc, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		return
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	expvar.NewString("build").Set(build)

	cfg := struct {
		Web struct {
			ReadTimeout     time.Duration
			WriteTimeout    time.Duration
			IdleTimeout     time.Duration
			ShutdownTimeout time.Duration
			APIHost         string
			// APIPort   string
			// DebugHost string
		}

		MongoDb struct {
			Url                  string
			UsersCollectionName  string
			UsersAPIDatabaseName string
			UsersMongoUser       string
			UsersMongoPassword   string
			UsersDEKAltName      string
			CMKeyPath            string
		}
	}{}

	cfg.Web.APIHost = getConfStrVal("API_HOST", "0.0.0.0:3000")

	cfg.MongoDb.Url = getConfStrVal("MONGO_DB_URL", "localhost:27017")
	cfg.MongoDb.UsersCollectionName = getConfStrVal("MONGO_DB_USER_COLLECTION_NAME", "users")
	cfg.MongoDb.UsersAPIDatabaseName = getConfStrVal("MONGO_DB_DATABASE_NAME", "user-db")
	cfg.MongoDb.UsersMongoUser = getConfStrVal("MONGO_DB_USERNAME", "user")
	cfg.MongoDb.UsersMongoPassword = getConfStrVal("MONGO_DB_PASSWORD", "pass")
	cfg.MongoDb.UsersDEKAltName = getConfStrVal("USER_DEK_ALIAS", "user-dek-key-1")
	cfg.MongoDb.CMKeyPath = getConfStrVal("USER_CMKey_PATH", "user-cmk-key-1")

	// -------------------------------------------------------------------------
	fmt.Printf("cfg.MongoDb.CMKeyPath : %s", cfg.MongoDb.CMKeyPath)
	fmt.Printf("cfg.MongoDb.CMKeyPath : %s", cfg.MongoDb.CMKeyPath)

	// -------------------------------------------------------------------------
	// Start up db
	keystore.UserDEKeyAlias = cfg.MongoDb.UsersDEKAltName

	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilSliceAsEmpty:   true,
	}

	mongoUrl := mongoScheme + cfg.MongoDb.UsersMongoUser + ":" + cfg.MongoDb.UsersMongoPassword + "@" + cfg.MongoDb.Url

	mongoRegistry := documentStore.UUIDTypeRegistry()

	clientOpts := options.Client().
		ApplyURI(mongoUrl).
		SetBSONOptions(bsonOpts).
		SetTimeout(5 * time.Second).
		SetRegistry(mongoRegistry)
	store, err := documentStore.StartEncryptedDB(clientOpts)
	if err != nil {
		return fmt.Errorf("mongodb connection failed: %q", err.Error())
	}
	db := store.Client.Database(cfg.MongoDb.UsersAPIDatabaseName)
	err = documentStore.StatusCheck(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to ping mongodb: %q", err.Error())
	}

	// -------------------------------------------------------------------------

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := v1.APIMuxConfig{
		Build:              build,
		Shutdown:           shutdown,
		Log:                log,
		DbClients:          store,
		UserDbName:         db.Name(),
		UserCollectionName: cfg.MongoDb.UsersCollectionName,
	}

	apiMux := v1.APIMux(cfgMux, handlers.Routes{})

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	log.Info(ctx, "starting service", "version", build)

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		_ = store.EncryptMgr.Close(ctx)

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func getConfStrVal(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
