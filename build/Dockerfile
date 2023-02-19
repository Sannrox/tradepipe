# syntax=docker/dockerfile:1

## Build
FROM golang:1.17-buster AS build

ARG TARGET
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN cd $TARGET && go build -o /cmd

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /cmd /cmd

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/cmd"]


