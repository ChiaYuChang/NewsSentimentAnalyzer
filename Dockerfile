FROM 1.21-rc-alpine AS builder
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.19.0

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mode dowload

ENV MIGRATION_PATH=./db/migration
ENV SQL_SCHEME_PATH=./db
ENV BIN_PATH=.

COPY thirdparty ./thirdparty
COPY pkgs ./pkgs
COPY internal ./internal
COPY global ./global
COPY db ./db
COPY config ./config
COPY main Makefile ./

RUN make sqlc-clean
RUN make sqlc-generate
RUN make build

FROM alpine:3.18
WORKDIR /usr/src/app
COPY --from=builder /nsa .
COPY views ./
COPY config ./
COPY _test ./