# ---------- Base: toolchain + module metadata ----------
FROM golang:1.27-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./

# ---------- Deps: vendored modules + compiled binaries ----------
# Vendored (no network needed at build time) — refresh with `go mod vendor`.
FROM base AS deps

COPY vendor/ ./vendor/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags='-s -w' -o /server . && \
    CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags='-s -w' -o /migrate ./cmd/migrate

# ---------- Development: live reload with Air ----------
FROM base AS development

ENV CGO_ENABLED=0
ENV GIN_MODE=debug

RUN go install github.com/air-verse/air@v1.67.4

# Pre-fetch modules into the image cache: the bind mount overlays /app
# at runtime, and external network may be unavailable there.
COPY vendor/ ./vendor/
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

# ---------- Production: minimal non-root runtime ----------
FROM alpine:3.21 AS production

ENV GIN_MODE=release

# ca-certificates is required for SMTP TLS (Gmail) and any HTTPS calls.
RUN apk add --no-cache ca-certificates && \
    adduser -D -h /app appuser

WORKDIR /app

COPY --from=deps /server /app/server
COPY --from=deps /migrate /app/migrate
# golang-migrate file source (file://migrations) needs the SQL on disk.
# Email templates are go:embed'ed into the binary, no copy needed.
COPY --from=deps /app/migrations /app/migrations

# Pre-create the uploads dir owned by appuser so the named volume
# inherits working ownership (avatars are written at runtime).
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["/app/server"]
