# `middleware/request_id.go`

Request tracing middleware: reuses a client-supplied `X-Request-ID` header or generates a UUID, stores it as `"requestID"` in the Gin context, and echoes it back on the response header.

Applied globally in `routes/routes.go` (`router.Use(middleware.RequestID())`).
