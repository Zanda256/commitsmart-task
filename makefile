# Check to see if we can use ash, in Alpine images, or default to BASH.

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

tidy:
	go mod tidy

curl-create:
	curl -i -X POST -H 'Content-Type: application/json' -d '{"name":"sekiranda","email":"seky@gmail.com","department":"IT","credit_card":"67725 37734 90273"}' http://localhost:3000/v1/users

curl-get-user:
	curl -i -X GET  http://localhost:3000/v1/users/ce937189-159d-4e5d-b307-00307eebe7d6

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

logs:
	docker logs -f commitsmart-user-api
