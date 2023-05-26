Simple REST API Finance Server
==========================

**Stack**:

- GO 1.20.4
- Postgres 15.3
- Fiber
- SimpleCache
- Swagger
- JWT
- RabbitMQ 3.11.16
- gRPC

**Reference data**: TwelveData

## API Methods:

```
/auth - autharization endpoints
```

- POST sign-up `/signup`
- POST sign-in `/signin`
- Get refresh token `/refresh`

```
/api/v1/symbols - api endpoints
```

- GET all available symbols
- POST new symbol
- PUT update symbol
- GET symbol by name `/:symbol`
- DELETE symbol by name `/:symbol`

## Before run:

1. Check and set up your configs in [config file](config/config.yaml)
2. Set up your secrets: DB Password, Twelve Data API key, salt for users passwords, jwt secret and RabbitMQ connection
   string. You can check [example.env](example.env)

## To run:

```shell
go run cmd/app/main.go
```
