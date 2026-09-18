# `routes/contact_routes.go`

Contact route table on the protected group.

## `RegisterContactRoutes(router *gin.RouterGroup)`

Creates the `/contacts` subgroup with `POST /`, `GET /`, `GET /:id`, `PUT /:id`, `DELETE /:id` → the `handlers` contact functions. The signature takes the already-authenticated `/api` group (see `routes/routes.go`), so contacts inherit JWT protection.
