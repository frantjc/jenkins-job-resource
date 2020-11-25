FROM golang:latest AS builder

WORKDIR /src

COPY go.mod .
# COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /assets/check ./cmd/check \
    && go build -o /assets/in ./cmd/in \
    && go build -o /assets/out ./cmd/out

FROM ubuntu:bionic AS resource

ENV DEBAIN_FRONTEND=noninteractive
RUN apt-get update
COPY --from=builder assets/ opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
# run tests

FROM resource
