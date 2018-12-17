GO ?= go
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
pkgs         = $(shell $(GO) list ./... | grep -v /vendor/)
PREFIX ?= _outputs

DOCKERFILE	?= Dockerfile
DOCKER_REPO ?= aixeshunter
DOCKER_IMAGE_NAME ?= prometheus-plugin
DOCKER_IMAGE_TAG ?= v0.1

.PHONY: build
build:
	@echo ">> building binaries"
	CGO_ENABLED=0 $(GO) build -o $(DOCKER_IMAGE_NAME) cmd/main.go

.PHONY: docker
docker: 
	@echo ">> building docker image from $(DOCKERFILE)"
	@docker build -t "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)
