# syntax=docker/dockerfile:1.7

# ----------------------------------------------------------------------------
# Estágio 1: compilar o binário do PocketBase
# ----------------------------------------------------------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache build-base gcc musl-dev

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/pocketbase .

# ----------------------------------------------------------------------------
# Estágio 2: imagem final com PocketBase
# ----------------------------------------------------------------------------
FROM alpine:3.22.1

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1001 pocketbase \
    && adduser -u 1001 -G pocketbase -s /bin/sh -D pocketbase

WORKDIR /app

COPY --from=builder /out/pocketbase /app/pocketbase
RUN chmod +x /app/pocketbase

RUN mkdir -p /app/pb_public /app/pb_hooks /app/pb_migrations /app/pb_data \
    && chown -R pocketbase:pocketbase /app

ARG URL_LIBSQL_TURSO
ENV TZ=America/Sao_Paulo \
    PB_PORT=8090 \
    URL_LIBSQL_TURSO=${URL_LIBSQL_TURSO}

EXPOSE 8090

VOLUME ["/app/pb_data"]

USER pocketbase

CMD ["/app/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/app/pb_data"]
