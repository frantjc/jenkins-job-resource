FROM golang:alpine AS builder

WORKDIR /src

COPY go.mod .
# COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /assets/check ./cmd/check \
    && go build -o /assets/in ./cmd/in \
    && go build -o /assets/out ./cmd/out

FROM alpine:edge AS resource

RUN apk update
COPY --from=builder assets/ opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
# run tests

FROM resource
