# Go Build
EXECUTABLE := db-seed-runner
BUILD_DIR := bin

# Helm Packaging
HELM_CHART_DIR := charts/db-seed-runner
HELM_RELEASE_DIR := helm-releases
PAGES_DIR := .gh-pages
REPO_URL := https://jnsaph.github.io/db-seed-runner
CHART_VERSION := $(shell yq '.version' $(HELM_CHART_DIR)/Chart.yaml)
APP_VERSION := $(shell yq '.appVersion' $(HELM_CHART_DIR)/Chart.yaml)
CHART_NAME := $(shell yq '.name' $(HELM_CHART_DIR)/Chart.yaml)

# Postgres
export DB_USER := postgres
export DB_PASSWORD := password

# Docker
IMAGE_REGISTRY ?= ghcr.io
IMAGE_REPO ?= jnsaph
IMAGE_NAME ?= $(CHART_NAME)
IMAGE_TAG ?= $(APP_VERSION)
FULL_IMAGE_NAME := $(IMAGE_REGISTRY)/$(IMAGE_REPO)/$(IMAGE_NAME)

build:
	go build -o bin/$(EXECUTABLE) 

run: build
	./bin/$(EXECUTABLE)

# Builds the multi-arch image and stores it in the build cache. (won't appear in local docker images)
docker-build:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(FULL_IMAGE_NAME):$(IMAGE_TAG) -t $(FULL_IMAGE_NAME):latest .

# Push multi-arch image to registry
docker-push:
	@echo "Logging in to GHCR..."
	@echo $${GHCR_PAT} | docker login $(IMAGE_REGISTRY) -u $(IMAGE_REPO) --password-stdin
	docker buildx build --platform linux/amd64,linux/arm64 -t $(FULL_IMAGE_NAME):$(IMAGE_TAG) -t $(FULL_IMAGE_NAME):latest --push .

# Build for local architecture and load into local Docker daemon
docker-build-local:
	docker buildx build --load -t $(FULL_IMAGE_NAME):$(IMAGE_TAG) -t $(FULL_IMAGE_NAME):latest .


lint-helm:
	helm lint $(HELM_CHART_DIR)

check-worktree:
	@test -d $(PAGES_DIR) || (echo "Error: $(PAGES_DIR) does not exist. Please run 'git worktree add .gh-pages gh-pages'" && exit 1)

package-helm: lint-helm
	helm package $(HELM_CHART_DIR) --destination $(HELM_RELEASE_DIR)

publish-helm: check-worktree package-helm
	@if [ -f $(PAGES_DIR)/index.yaml ]; then \
		yq -e \
		  --arg name "$(CHART_NAME)" \
		  --arg version "$(CHART_VERSION)" \
		  '.entries[$$name][] | select(.version == $$version)' \
		  $(PAGES_DIR)/index.yaml > /dev/null && \
		echo "Error: $(CHART_NAME) version $(CHART_VERSION) already exists. Change the version in $(HELM_CHART_DIR)/Chart.yaml!" && \
		exit 1 || true; \
	fi
	cp $(HELM_RELEASE_DIR)/$(CHART_NAME)-$(CHART_VERSION).tgz $(PAGES_DIR)/
	helm repo index $(PAGES_DIR) --url $(REPO_URL) --merge $(PAGES_DIR)/index.yaml

publish-helm-mac: check-worktree package-helm
	@if [ -f $(PAGES_DIR)/index.yaml ]; then \
		yq -e '.entries.$(CHART_NAME)[] | select(.version == "$(CHART_VERSION)")' \
			$(PAGES_DIR)/index.yaml > /dev/null && \
		echo "Error: $(CHART_NAME) version $(CHART_VERSION) already exists. Change the version in $(HELM_CHART_DIR)/Chart.yaml!" && \
		exit 1 || true; \
	fi
	cp $(HELM_RELEASE_DIR)/$(CHART_NAME)-$(CHART_VERSION).tgz $(PAGES_DIR)/
	helm repo index $(PAGES_DIR) --url $(REPO_URL) --merge $(PAGES_DIR)/index.yaml
