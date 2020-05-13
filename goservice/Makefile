SHELL             := /usr/bin/env bash
DEV_REGISTRY      ?= quay.io
IMAGE_BASE_NAME   ?= goservice
PROJECT_NAME      ?= service_preview
DOCKER_REPO       ?= $(DEV_REGISTRY)/$(PROJECT_NAME)/$(IMAGE_BASE_NAME)
GIT_COMMIT        := $(shell git rev-parse --verify HEAD)
TAG               ?= $(GIT_COMMIT)
IMAGE_NAME        := $(DOCKER_REPO)-goservice:$(TAG)
IMAGE_NAME_LATEST := $(DOCKER_REPO)-goservice:latest

# command for pushing images (can be replaced by "gcloud docker push" when pushing to gcloud)
DOCKER_PUSH       ?= docker push

########################################################################################

.PHONY: image.build image.push

image.build:
	@echo ">>> Building goservice image"
	docker build \
		-t $(IMAGE_NAME) \
		-t $(IMAGE_NAME_LATEST) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-f Dockerfile .
	@echo ">>> done!"

image.push: image.build
	@echo ">>> Pushing goservice image"
	$(DOCKER_PUSH) $(IMAGE_NAME)
	$(DOCKER_PUSH) $(IMAGE_NAME_LATEST)
	@echo ">>> done!"

