# Check to see if we can use ash, in Alpine images, or default to BASH.
#docker compose rm -v -f ./zarf/docker/docker-compose.yml
#curl -il -X POST http://localhost:3000/v1/users
	# curl -il -X POST -H 'Content-Type: application/json' -d '{"name":"bill","email":"b@gmail.com","roles":["ADMIN"],"department":"IT","password":"123","passwordConfirm":"123"}' http://localhost:3000/v1/users
#rebuild-compose-project: rebuild-service-container service-compose

SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Define dependencies
BASE_IMAGE_NAME := commitsmart
SERVICE_NAME    := user-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
GO_VERSION := 1.23.1

run:
	go run -tags cse app/services/user-api/main.go

#build:
#    cd app/services/user-api/ && go build -tags cse -ldflags "-X main.build=test-run"

rebuild-compose-project:

tidy:
	go mod tidy

curl-create:
	curl -il -X POST -H 'Content-Type: application/json' -d '{"name":"bill","email":"b@gmail.com","department":"IT","credit_card":"72635 37734 90273"}' http://localhost:3000/v1/users

service-compose:
	docker compose -p commitsmart-project -f ./zarf/docker/docker-compose.yml up -d

push-containers:
	docker tag commitsmart/user-api:0.0.1  zanda256/user-api:0.0.1 && docker push zanda256/user-api:0.0.1

service-compose-test:
	docker compose --dry-run -f ./zarf/docker/docker-compose.yml up

stop-compose-project:
	docker compose -f ./zarf/docker/docker-compose.yml stop user-api mongodb

rebuild-service-container:
	docker build \
		--platform linux/amd64 \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
