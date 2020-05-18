SHELL             := /usr/bin/env bash
DEV_REGISTRY      ?= docker.io
DOCKER_USER       ?= brucehorn
PROJECT_NAME      ?= service_preview
DOCKER_REPO       ?= $(DEV_REGISTRY)/$(DOCKER_USER)/$(PROJECT_NAME)

# command for pushing images (can be replaced by "gcloud docker push" when pushing to gcloud)
DOCKER_PUSH       ?= docker push

# Export these variable to the sub-makefiles.
export

#########################################################################

# Phony targets for make: these are not files and will always run when invoked.
.PHONY: build_images push_images
all:    build_images push_images

build_images:
	@echo ">>> Building images..."
	make -C goservice  build_image
	make -C pyservice  build_image
	@echo ">>> images built! You can 'make push_images' now."

# note: the `image.push` in the subdirs trigger a build automatically
push_images: build_images
	@echo ">>> Pushing images..."
	make -C goservice  push_image
	make -C pyservice  push_image
	@echo ">>> images pushed! You can 'make deploy' (or 'redeploy') now."

deploy:
	@echo ">>> deploying to cluster..."
	make -C goservice deploy
	make -C pyservice deploy