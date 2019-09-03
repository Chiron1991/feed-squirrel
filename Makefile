HERE = $(shell pwd)

FRONTEND_DIR = $(HERE)/frontend
FRONTEND_BUILD_DIR = $(FRONTEND_DIR)/build

CONTAINER_TAG = feed-squirrel

build_all: build_frontend build_container

build_container:
	docker build --no-cache -t $(CONTAINER_TAG) $(HERE)

.ONESHELL:
build_frontend:
	set -e
	cd $(FRONTEND_DIR)
	rm -rf $(FRONTEND_BUILD_DIR)
	npm ci
	npm run-script build
