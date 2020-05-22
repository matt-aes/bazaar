# SHELL, DOCKER_REPO defined in ../Makefile
SERVICE_NAME      := inventoryservice

GIT_COMMIT        := $(shell git rev-parse --verify HEAD)
TAG               := $(GIT_COMMIT)
IMAGE_NAME        := $(DOCKER_REPO)-$(SERVICE_NAME):$(TAG)
IMAGE_NAME_LATEST := $(DOCKER_REPO)-$(SERVICE_NAME):latest

########################################################################################

# Phony targets for make: these are not files and will always run when invoked.
.PHONY: build_image push_image

build_image:
	@echo ">>> Building " $(IMAGE_NAME)
	@docker build \
		-t $(IMAGE_NAME) \
		-t $(IMAGE_NAME_LATEST) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-f Dockerfile .
	@echo ">>> done!"

push_image: build_image
	$(DOCKER_PUSH) $(IMAGE_NAME)
	$(DOCKER_PUSH) $(IMAGE_NAME_LATEST)
	@echo ">>> done!"

deploy:
	@echo ">>> deploying inventory_service to cluster..."
	kubectl apply -f k8s/inventory_service.yaml
