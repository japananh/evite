# App Invitation Service

## How to run project locally?

You need to install Go version 1.18 and an IDE/editor such as Goland or VSCode to run the project.

```bash
# create `.env` file
cp .env.template .env
# run docker-compose to create mySQL database and server
make start
# down all services
make stop
# run unit test locally
make test
```

## Project structure

### Folder structure

![project architecture diagram](https://imge.cloud/images/2022/06/07/rZ1cup.png)

### API Endpoints

- POST `/api/v1/register`: create a new user with email and password
- POST `/api/v1/login`: login with email and password
- POST `/api/v1/login/invitation`: login with an invitation token
- GET `/api/v1/users/invitation`: generate an invitation token
- GET `/api/v1/token/validation?invite_token=`: validate an invitation token
- DELETE `/api/v1/token/{invite_token}`: delete an invitation token

## Documentation

Swagger
