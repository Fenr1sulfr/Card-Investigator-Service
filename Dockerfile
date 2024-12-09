FROM golang:1.21.1-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc build-base curl


COPY ./ ./

ARG MIGRATE_VERSION=4.15.2
ARG OS=linux
ARG ARCH=amd64

# Install the Migrate CLI tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.${OS}-${ARCH}.tar.gz | tar xvz && \
    mv migrate /bin/migrate


RUN go mod download

RUN make build/api

# RUN make db/migrations/up
COPY ./migrations /migrations

ENTRYPOINT [ "./bin/api" ]