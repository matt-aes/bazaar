SHELL             := /usr/bin/env bash
DEV_REGISTRY      ?= docker.io/brucehorn
PROJECT_NAME      ?= service_preview
DOCKER_REPO       ?= $(DEV_REGISTRY)/$(PROJECT_NAME)

# command for pushing images (can be replaced by "gcloud docker push" when pushing to gcloud)
DOCKER_PUSH       ?= docker push

# Export these variable to the sub-makefiles.
export

#########################################################################

# Phony targets for make: these are not files and will always run when invoked.
.PHONY: build_images push_images app images inventory
all:    build_images push_images

build_images:
	@echo ">>> Building images..."
	make -C specsservice      build_image
	make -C inventoryservice  build_image
	make -C imageservice      build_image
	make -C appservice        build_image
	@echo ">>> images built! You can 'make push_images' now."

# note: the `image.push` in the subdirs trigger a build automatically
push_images: build_images
	@echo ">>> Pushing images..."
	make -C specsservice      push_image
	make -C inventoryservice  push_image
	make -C imageservice      push_image
	make -C appservice        push_image
	@echo ">>> images pushed! You can 'make deploy' now."

# Deploy all services once images have been built.
deploy:
	@echo ">>> deploying to cluster..."
	make -C specsservice     deploy
	make -C inventoryservice deploy
	make -C imageservice     deploy
	make -C appservice       deploy
	@echo ""
	@echo ">>> applying the traffic manager and configuring RBAC..."
	kubectl apply -f k8s/traffic-agent-rbac.yaml

# Build, push, deploy individual services.
app:
	make -C appservice       build_image
	make -C appservice       push_image
	make -C appservice       deploy

images:
	make -C imageservice     build_image
	make -C imageservice     push_image
	make -C imageservice     deploy

inventory:
	make -C inventoryservice build_image
	make -C inventoryservice push_image
	make -C inventoryservice deploy

specs:
	make -C specsservice     build_image
	make -C specsservice     push_image
	make -C specsservice     deploy
