# Gryphon file API
- API to send files, upload to aws s3 and return public url. For now, only sending images in base64 is supported.

## Getting started with docker-compose
- Install docker
- Install docker-compose
- Clone this project
- Run `cp .env.example .env`
- Open file .env and set environments. USER_EMAIL, USER_PASSWORD, USER_AWS_ACCESS_KEY_ID, USER_AWS_SECRET_ACCESS_KEY, USER_AWS_REGION and USER_AWS_BUCKET is a first user. Change to your data.
- Run `docker-compose build`
- Run `docker-compose up`

## Getting started without docker-compose
- Install go
- Clone this project
- Run `cp .env.example .env`
- Open file .env and change environments, set your database credentials. USER_EMAIL, USER_PASSWORD, USER_AWS_ACCESS_KEY_ID, USER_AWS_SECRET_ACCESS_KEY, USER_AWS_REGION and USER_AWS_BUCKET is a first user. Change to your data.
- Run `go run main.go`


## How to use
After service is working, access swagger in [http://localhost:8083/api/swagger/index.html#/users/post_login](http://localhost:8083/api/swagger/index.html#/users/post_login), click try it out, set your credentials defined in .env file and send request. Get the token and use it to send images.

## Other Commands
- Generate API docs `swag init --parseDependency --parseInternal`
- Run tests `go test`
- Add swag in path linus `export PATH=$(go env GOPATH)/bin:$PATH`