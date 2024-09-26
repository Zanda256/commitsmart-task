# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)


run:
	go run app/services/user-api/main.go

tidy:
	go mod tidy

curl-create:
	curl -il -X POST http://localhost:3000/v1/users
	# curl -il -X POST -H 'Content-Type: application/json' -d '{"name":"bill","email":"b@gmail.com","roles":["ADMIN"],"department":"IT","password":"123","passwordConfirm":"123"}' http://localhost:3000/v1/users