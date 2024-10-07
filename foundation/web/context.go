package web

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"time"
)

type ctxKey int

const key ctxKey = 1

type pathParamsKey struct{}

// ParamsKey is the request context key under which URL params are stored.
var PathParamsKey = pathParamsKey{}

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// SetValues sets the specified Values in the context.
func SetValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}

func AttachPathParams(ctx context.Context, val *httprouter.Params) context.Context {
	ctx = context.WithValue(ctx, PathParamsKey, val)
	fmt.Printf("\nAttachPathParams : ctx : %#v\n", ctx)
	return ctx
}

func GetPathParams(ctx context.Context) interface{} {
	return ctx.Value(PathParamsKey)
}

// GetValues returns the values from the context.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v.TraceID
}

// GetTime returns the time from the context.
func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return time.Now()
	}
	return v.Now
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}
