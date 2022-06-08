# App Invitation Service

## Requirements

- [x] The App Admin generates an invitation token (6 - 12 digit alphanumberic string)
- [x] The invitation token is then used to log in
- [x] The APIs should be REST-ful
- [x] The admin endpoints should be authenticated with JWT
- [x] Invite tokens to expire after 7 days
- [x] Invite tokens can be recalled (disabled)
- [x] A public endpoint for validating the invite token
- [x] Develop the API in Golang
- [x] Frameworks/Libraries: Gin, Gorm,
- [x] Use in-memory storage for the tokens (Redis)
- [x] Use an actual DB
- [x] Deployment instructions are written in `README.md` 
- [x] Write tests (unit tests and integration tests)
- [ ] An admin can get an overview of active and inactive tokens
- [ ] Document the APIs in Swagger
- [ ] The invite token validation logic needs to be throttled (limit the requests coming from a
  specific client)

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
* Storage Layer: interact with databases such as RLDB/NoSQL or File System or Remote API.
* Repository Layer (optional): provides data for the Business Layer.
* Business Layer: business logic happened here.
* Transport Layer: receive HTTP requests from Client, parse data (if needed).

![project architecture diagram](https://i.postimg.cc/8zfZW6sW/clean-arch-diagram.png)

### API Endpoints

The Go server will run default on port `8000`.

- GET `/api/v1/users/invitation`: Admin generates an invitation token
- POST `/api/v1/login/invitation`: login with an invitation token
- GET `/api/v1/token/validation?invitation_token=`: validate an invitation token
- GET `/api/v1/token/invitation?status=`: Admin gets invitation token by status
- PATCH `/api/v1/token/invitation/:invitation_token`: Admin disable/enable an invitation token
- POST `/api/v1/register`: create a new user with email and password
- POST `/api/v1/login`: login with email and password

### Documentation

Swagger
