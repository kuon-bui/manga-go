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

FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM builder AS api-builder
RUN go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/dev/main.go -o swagger-docs/ --parseDependency --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/dev

FROM builder AS queue-builder
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/queue

FROM alpine:3.22 AS production-base

RUN apk add --no-cache ca-certificates \
  && addgroup -S app \
  && adduser -S -G app app

WORKDIR /app

COPY --chown=app:app resources ./resources

USER app

FROM production-base AS production-api
COPY --from=api-builder --chown=app:app /out/app ./app
EXPOSE 8080
ENTRYPOINT ["/app/app"]

FROM production-base AS production-queue
COPY --from=queue-builder --chown=app:app /out/app ./app
ENTRYPOINT ["/app/app"]

