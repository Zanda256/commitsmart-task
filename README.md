# commitsmart-task
A CRUD application with a client facing REST API that utilizes mongodb database for data persistence

## Project dependencies
- go 1.22.*
- mongodb 7.*

## How to run the applicaton
Create a folder on your computer Clone repository to your computer into created folder

SSH
```git clone git@github.com:Zanda256/commitsmart-task.git```

Navigate into created folder
```cd commitsmart-task```

Download dependent modules
`go get .` followed by `go mod tidy`

To run the RESTful API on your local machine , use the following command. (WARNING: This setup is not compatible with apple silicon processor machines)
`make run`

To run the API as docker compose project, use the following command
`make service-compose`

Test the dummy endpoint using `curl`
`curl -il -X POST -H 'Content-Type: application/json' -d '{"name":"bill","email":"b@gmail.com","department":"IT","credit_card":"72635 37734 90273"}' http://localhost:3000/v1/users`

