FROM golang:alpine AS builder
WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV CGO_ENABLED 0
RUN go build -o /assets/check ./cmd/check \
    && go build -o /assets/in ./cmd/in \
    && go build -o /assets/out ./cmd/out

RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk update
RUN apk add --no-cache ca-certificates
COPY --from=builder assets/ opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
COPY --from=builder /tests /tests

ARG JENKINS_URL
ARG JENKINS_JOB
ARG JENKINS_JOB_TOKEN
ARG JENKINS_USERNAME
ARG JENKINS_API_TOKEN

RUN set -e; for test in /tests/*.test; do \
		$test -ginkgo.v; \
	done

FROM resource
