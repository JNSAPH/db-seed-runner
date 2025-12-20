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

# Phony Targets
.PHONY: build run docker-build docker-push docker-build-local lint-helm check-worktree check-version publish-helm push-version full-release

#
# Go Build Targets
#
build:
	go build -o bin/$(EXECUTABLE) 

run: build
	./bin/$(EXECUTABLE)

#
# Docker Targets
#

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

#
# Helm Packaging Targets
#
lint-helm:
	@echo "Linting Helm chart in $(HELM_CHART_DIR)..."
	helm lint $(HELM_CHART_DIR)
	@echo "Helm chart linting passed."

check-worktree:
	@test -d $(PAGES_DIR) || (echo "Error: $(PAGES_DIR) does not exist. Please run 'git worktree add .gh-pages gh-pages'" && exit 1)

check-version: check-worktree
	@echo "Refreshing .gh-pages worktree..."
	cd $(PAGES_DIR) && git pull origin gh-pages
	@echo "Checking if $(CHART_NAME) version $(CHART_VERSION) already exists in $(PAGES_DIR)/index.yaml..."
	@if [ -f $(PAGES_DIR)/index.yaml ]; then \
		if name="$(CHART_NAME)" version="$(CHART_VERSION)" yq -e \
			'.entries[env(name)][] | select(.version == env(version))' \
			$(PAGES_DIR)/index.yaml > /dev/null 2>&1; then \
			echo "Error: $(CHART_NAME) version $(CHART_VERSION) already exists. Change the version in $(HELM_CHART_DIR)/Chart.yaml!"; \
			exit 1; \
		fi \
	fi
	@echo "Version $(CHART_VERSION) is new. Proceeding..."

publish-helm: check-version
	@echo "Packaging Helm chart..."
	helm package $(HELM_CHART_DIR) --destination $(HELM_RELEASE_DIR)
	cp $(HELM_RELEASE_DIR)/$(CHART_NAME)-$(CHART_VERSION).tgz $(PAGES_DIR)/
	helm repo index $(PAGES_DIR) --url $(REPO_URL) --merge $(PAGES_DIR)/index.yaml
	@echo "Helm chart $(CHART_NAME) version $(CHART_VERSION) published to $(PAGES_DIR)."

# 
# Git Targets
#
push-version:
	@echo "Tagging git repository with version $(APP_VERSION)..."
	git tag $(APP_VERSION)
	git push origin $(APP_VERSION)

full-release: check-version docker-push publish-helm
	@echo "\n\n\n----------------------------------------"
	@echo "Release $(APP_VERSION) (Chart $(CHART_VERSION)) completed successfully!\n"
	@echo "Next steps: \n- Push the .gh-pages branch to GitHub \n- Run 'make push-version' to tag the release in git"
	@echo "----------------------------------------"