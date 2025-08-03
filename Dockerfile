FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o shortener ./cmd/url-shortener/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

COPY --from=builder /app/shortener .
COPY --from=builder /app/config ./config

ENV CONFIG_PATH=./config/local.yaml
EXPOSE 6969

ENTRYPOINT ["./shortener"]
