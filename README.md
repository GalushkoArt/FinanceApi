Simple REST API Finance Server
==========================

**Stack**:

- GO 1.20.4
- Postgres 15.3
- Fiber
- SimpleCache
- Swagger

**Reference data**: TwelveData

## API Methods:

```
/api/v1/symbols - endpoints
```

- GET all available symbols
- POST new symbol
- PUT update symbol
- GET symbol by name `\:symbol`
- DELETE symbol by name `\:symbol`

## Before run:

1. Set database configs
2. Set TwelveData API key

## To run:

```shell
go run cmd/app/main.go
```
