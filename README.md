Simple REST API Finance Server
==========================

**Stack**:

- GO 1.20.4
- Postgres 15.3
- Fiber
- SimpleCache
- Swagger
- JWT

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

1. Set database configs
2. Set TwelveData API key
3. Set DB password salt
4. Set JWT HMAC secret

## To run:

```shell
go run cmd/app/main.go
```
