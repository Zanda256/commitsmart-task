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

To run the RESTful API use the following command
`make run`

Test the dummy endpoint using `curl`
`curl -il -X POST http://localhost:3000/v1/users`

To run the API as docker compose project, use the following command
`make service-compose`
