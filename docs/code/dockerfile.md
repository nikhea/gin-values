# `Dockerfile`

Multi-target build mirroring a Node-style `base → deps → development / production` layout. Select with `docker build --target <name>` (compose files do this for you).

## Stages

- **base** (`golang:1.27-alpine`, `WORKDIR /app`): toolchain + `go.mod`/`go.sum` only, so later stages share the module metadata layer.
- **deps**: copies `vendor/` + source, compiles both static binaries (`CGO_ENABLED=0`, `-mod=vendor`). Vendored on purpose: build containers may have no module-proxy access, so the build needs zero network. Refresh with `go mod vendor` after changing dependencies.
- **development**: `GIN_MODE=debug`, installs pinned Air (`air-verse/air@v1.67.4`), pre-fetches modules into the image cache (the bind mount overlays `/app` at runtime and external network may be unavailable there), copies source, `CMD ["air", "-c", ".air.toml"]`.
- **production** (`alpine:3.21`, `GIN_MODE=release`): `ca-certificates` (Gmail SMTP TLS), non-root `appuser`, copies both binaries + `migrations/` from `deps` (email templates are `go:embed`ed, no copy needed), pre-chowned `/app/uploads` so the named volume inherits writable ownership. `CMD ["/app/server"]`.

See `docker-compose.md` for how each target is run.
