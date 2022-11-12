DOCKER = docker
DOCKER-COMPOSE = $(DOCKER) compose

GO = go
GOLANGCI-LINT = golangci-lint

test:
	@bin/$@

fmt:
	@$(GO) $@ ./...

download vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

dl: download
ven: vendor
ver: verify
format: fmt

.PHONY: fmt test download vendor verify lint dl ven ver format
