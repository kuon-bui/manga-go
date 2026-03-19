FROM golang:1.25-alpine AS base

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download
RUN apk add make
RUN go install github.com/rubenv/sql-migrate/...@latest

FROM base AS queue
CMD [ "air", "-c", "queue.air.toml"]

FROM base AS api
CMD ["air", "-c", "api.air.toml"]

