FROM golang:1.21-rc-alpine AS builder
WORKDIR /usr/src/app
RUN apk update
RUN apk upgrade
RUN apk add build-base

# setup go env
ENV CGO_ENABLED=1
ENV GOOS=linux

ENV POSTGRES_USERNAME=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_PASSWORD_FILE=
ENV POSTGRES_HOST=localhost
ENV POSTGRES_PORT=5432 
ENV DB_NAME=db
ENV APP_NAME=app

# install go packages
COPY go.mod go.sum ./
RUN go mod download


COPY thirdparty ./thirdparty
COPY pkgs ./pkgs
COPY internal ./internal
COPY global ./global
COPY db ./db
COPY config ./config
COPY .env main.go ./
RUN go build -o ./${APP_NAME} ./main.go

FROM alpine:3.18
WORKDIR /usr/src/app

ENV POSTGRES_USERNAME=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_PASSWORD_FILE=
ENV POSTGRES_HOST=localhost
ENV POSTGRES_PORT=5432
ENV DB_NAME=db
ENV APP_NAME=app
ENV APP_PORT=8000

COPY --from=builder /usr/src/app/${APP_NAME} .
COPY views ./views
COPY config  ./config
COPY _test ./_test
RUN chmod +x ./app
ENTRYPOINT [ "./app" ]