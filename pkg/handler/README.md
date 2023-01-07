# pkg/handler

The `pkg/handler` package defines a top level "app" handler that can be embedded by micro-services.

It's best to think of micro-services as a management tool, as opposed to something that provides a technical advantage. Therefore, start with everything in `main.go`. Then spin off micro-services later as needed. Each additional service must have a corresponding `cmd/${SERVICE_NAME}/main.go` entry-point, and `internal/${SERVICE_NAME}` folder with logic specific to the service.

Shared route handler functions can be defined on the top level handler struct,for an example see `pkg/handler/docs.go`.

Response types for API endpoints must be defined in `pkg/share`, see for example `share.Response` and `share.ErrResponse`

Do not import `pkg/middleware` in this package. Services must embed the top level handler, setup middleware, and define the service route handlers, see examples in `internal/app/handler.go`
