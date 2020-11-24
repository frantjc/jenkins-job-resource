FROM golang:ubuntu AS builder

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
# build to /assets/[check|in|out]

FROM ubuntu:bionic AS resource

ENV DEBAIN_FRONTEND=noninteractive
RUN apt-get update
COPY --from=builder assets/ opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
# run tests

FROM resource
