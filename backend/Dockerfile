FROM golang:1.24-alpine AS backend_builder
WORKDIR /backend
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:3.21
WORKDIR /backend
COPY --from=backend_builder /backend/main .
COPY --from=backend_builder /backend/start.sh .
COPY --from=backend_builder /backend/wait-for.sh .
COPY cmd/migration ./cmd/migration

EXPOSE 8080

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/bin/migrate