# syntax=docker/dockerfile:1.7

# ----------------------------------------------------------------------------
# Estágio 1: compilar o binário do PocketBase
# ----------------------------------------------------------------------------
FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o /out/pocketbase .
RUN mkdir -p /out/app/pb_public /out/app/pb_hooks /out/app/pb_migrations /out/app/pb_data \
    && chown -R 1001:1001 /out/app

# ----------------------------------------------------------------------------
# Estágio 2: imagem final com PocketBase
# ----------------------------------------------------------------------------
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/pocketbase /app/pocketbase
COPY --from=builder --chown=1001:1001 /out/app /app

ENV TZ=America/Sao_Paulo \
    PB_PORT=8090

EXPOSE 8090

VOLUME ["/app/pb_data"]

USER 1001:1001

CMD ["/app/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/app/pb_data"]
