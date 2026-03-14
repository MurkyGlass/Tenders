# syntax=docker/dockerfile:1.7

FROM golang:1.23 AS builder

WORKDIR /app

COPY server/go.mod server/go.sum ./server/

WORKDIR /app/server
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

WORKDIR /app
COPY server ./server
COPY client ./client

WORKDIR /app/server
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/application ./src/app/main.go

FROM alpine:3.20 AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/application ./application

COPY --from=builder /app/client ./client

EXPOSE 80

CMD ["./application"]
