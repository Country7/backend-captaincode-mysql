# Build stage
# FROM golang:1.22.2-alpine3.19 AS builder
FROM golang:1.22.2 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# Alpine Linux использует пакетный менеджер (apk), а не apt-get.
# RUN apk update && apk add curl
RUN apt-get update
RUN apt-get install -y curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
# FROM alpine:3.19
FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY db/migration ./migration

# Для Ubuntu Установка MySQL клиента
RUN apt-get update
RUN apt-get install -y mysql-client

# Для alpine:3.19 Добавьте репозиторий MySQL и установите MySQL клиент
# RUN echo 'https://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories
# RUN apk update
# RUN apk add mysql-client

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
