DOCKER_IMAGE_NAME ?= nhl_recap
DOCKER_IMAGE_TAG ?= latest
DOCKER_REGISTRY ?= ghcr.io

.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: docker_build_image
docker_build_image: ## docker build
	docker build -t ${DOCKER_REGISTRY}/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

.PHONY: docker_push_image
docker_push_image: ## docker push
	docker push ${DOCKER_REGISTRY}/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: docker_build_and_push_image
docker_build_and_push_image: docker_build_image docker_push_image ## docker build and push

.PHONY: go_mod_verify
go_mod_verify: ## go mod verify
	go mod verify

.PHONY: go_build
go_build: ## go build
	go build -o $(DOCKER_IMAGE_NAME) main.go

.PHONY: lint
lint: ## golint ./...
	golint ./...

.PHONY: test
test: ## go test -race -vet=off ./...
	go test -race -vet=off ./...