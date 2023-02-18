ARG base_image=alpine:3.16
ARG build_image=golang:1.20-alpine3.16

FROM ${build_image} AS build
ENV CGO_ENABLED 0
WORKDIR $GOPATH/src/github.com/frantjc/jenkins-job-resource
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /assets/check ./cmd/check
RUN go build -o /assets ./cmd/in
RUN go build -o /assets ./cmd/out
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM ${base_image} AS resource
RUN apk update && \
	apk add --no-cache ca-certificates
COPY --from=build /assets/ /opt/resource/

FROM resource AS test
COPY --from=build /tests/ /tests/
ARG JENKINS_URL
ARG JENKINS_JOB
ARG JENKINS_JOB_TOKEN
ARG JENKINS_USERNAME
ARG JENKINS_API_TOKEN
ARG JENKINS_JOB_ARTIFACT
ARG JENKINS_JOB_RESULT
RUN set -e; for test in /tests/*.test; do \
		$test -ginkgo.v; \
	done

FROM resource
