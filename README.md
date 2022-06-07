# App Invitation Service

## How to run this project

You need to install Go version 1.18 and an IDE/editor such as Goland or VSCode to run the project.

```bash
# create `.env` file
cp .env.template .env
# run docker-compose to create MySQL database and Go server
make start
# down all services
make stop
# run unit test locally
make test
```

## Description

### Project structure

This project has 5 Domain layers:
* Model Layer
* Storage Layer
* Repository Layer
* Business Layer
* Transport Layer

![project architecture diagram](https://i.postimg.cc/8zfZW6sW/clean-arch-diagram.png)

### API Endpoints

The Go server will run default on port `8000`.

- POST `/api/v1/register`: create a new user with email and password
- POST `/api/v1/login`: login with email and password
- POST `/api/v1/login/invitation`: login with an invitation token
- GET `/api/v1/users/invitation`: generate an invitation token
- GET `/api/v1/token/validation?invite_token=`: validate an invitation token
- DELETE `/api/v1/token/{invite_token}`: delete an invitation token

## Documentation

Swagger
