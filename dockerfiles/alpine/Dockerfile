ARG base_image=alpine:edge
ARG build_image=golang

FROM ${base_image} as base
FROM ${build_image} as build

FROM build AS builder
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

FROM base AS resource
RUN apk update
RUN apk add --no-cache ca-certificates
COPY --from=builder /assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
COPY --from=builder /tests/ /tests/

ARG SKIP_TESTS
ARG JENKINS_URL
ARG JENKINS_JOB
ARG JENKINS_JOB_TOKEN
ARG JENKINS_USERNAME
ARG JENKINS_API_TOKEN

RUN set -e; for test in /tests/*.test; do \
		$test -ginkgo.v; \
	done

FROM resource
