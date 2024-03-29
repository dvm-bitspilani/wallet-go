# Wallet-Go

This is a Golang port of DVM Wallet.
<br>
Currently remains largely untested, with no realtime functionality and no proper way of packaging and deployment.
<br>

This template codebase was generated by [Autostrada](https://autostrada.dev/).

## Getting started

Before running the application you will need a working PostgreSQL installation and a valid DSN (data source name) for connecting to the database.

Please open the `cmd/api/main.go` file and edit the `db-dsn` command-line flag to include your valid DSN as the default value.

```
flag.StringVar(&cfg.db.dsn, "db-dsn", "YOUR DSN GOES HERE", "postgreSQL DSN")
```

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the application using `go run ./cmd/api`:

```
$ go mod tidy
$ go generate ./ent
$ go run ./cmd/api
```

If you make a request to the `GET /status` endpoint using `curl` you should get a response like this:

```
$ curl -i localhost:4444/status
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 09 May 2022 20:46:37 GMT
Content-Length: 23

{
    "Status": "OK",
}
```

## Project structure

|                           |                                                                                                                                                                                        |
|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **`cmd/api`**             | Application-specific code (handlers, routing, middleware) for dealing with HTTP requests and responses.                                                                                |
| `↳ cmd/api/context/`      | Contains helpers for working with request context.                                                                                                                                     |
| `↳ cmd/api/errors/`       | Contains helpers for managing and responding to error conditions.                                                                                                                      |
| `↳ cmd/api/handlers/`     | Contains application HTTP handlers.                                                                                                                                                    |
| `↳ cmd/api/main.go`       | The entry point for the application. Responsible for parsing configuration settings initializing dependencies and running the server. Start here when you're looking through the code. |
| `↳ cmd/api/middleware.go` | Contains application middleware.                                                                                                                                                       |
| `↳ cmd/api/routes.go`     | Contains application route mappings.                                                                                                                                                   |
| `↳ cmd/api/server.go`     | Contains a helper functions for starting and gracefully shutting down the server.                                                                                                      |

|                         |                                                                       |
|-------------------------|-----------------------------------------------------------------------|
| **`internal`**          | Contains various helper packages used by the application.             |
| `↳ internal/database/`  | Contains database-related code (setup, connection and extra helpers). |
| `↳ internal/helpers/`   | Contains application related helpers functions.                       |
| `↳ internal/password/`  | Contains helper functions for hashing and verifying passwords.        |
| `↳ internal/request/`   | Contains helper functions for decoding JSON requests.                 |
| `↳ internal/response/`  | Contains helper functions for sending JSON responses.                 |
| `↳ internal/validator/` | Contains validation helpers.                                          |
| `↳ internal/version/`   | Contains the application version number definition.                   |



|                 |                                                                       |
|-----------------|-----------------------------------------------------------------------|
| **`ent`**       | Contains entgo related code.                                          |
| `↳ ent/schema/` | Contains all database-related models and schemas for the application. |



|               |                                             |
|---------------|---------------------------------------------|
| **`service`** | Contains model related operation functions. |


## Application version

The application version number is defined in a `Get()` function in the `internal/version/version.go` file. Feel free to change this as necessary.

```
package version

func Get() string {
    return "0.0.1"
}
```

## Changing the module path

The module path is currently set to `dvm.wallet/harsh`. If you want to change this please find and replace all instances of `dvm.wallet/harsh` in the codebase with your own module path.