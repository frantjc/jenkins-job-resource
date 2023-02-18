DOCKER = docker
DOCKER-COMPOSE = $(DOCKER) compose
GIT = git
GO = go
GOLANGCI-LINT = golangci-lint

SEMVER ?= 0.1.1

test:
	@bin/$@

fmt:
	@$(GO) $@ ./...

download vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

release:
	@$(GIT) tag -a v$(SEMVER) -m v$(SEMVER)
	@$(GIT) push --follow-tags

dl: download
ven: vendor
ver: verify
format: fmt

.PHONY: fmt test download vendor verify lint release dl ven ver format
