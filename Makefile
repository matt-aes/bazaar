SHELL             := /usr/bin/env bash
DEV_REGISTRY      ?= brucehorn
PROJECT_NAME      ?= service_preview

#########################################################################

all: image.build image.push

image.build:
	echo ">>> Building images..."
	make -C goservice  image.build
	make -C pyservice  image.build
	echo ">>> images built! You can 'make image.push' now."

# note: the `image.push` in the subdirs trigger a build automatically
image.push: image.build
	echo ">>> Pushing images..."
	make -C goservice  image.push
	make -C pyservice  image.push
	echo ">>> images pushed! You can 'make deploy' (or 'redeploy') now."
