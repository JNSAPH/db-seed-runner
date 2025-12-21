# Go Build
EXECUTABLE := db-seed-runner
BUILD_DIR := bin

# Helm Packaging
HELM_CHART_DIR := charts/db-seed-runner
HELM_RELEASE_DIR := helm-releases
PAGES_DIR := .gh-pages
REPO_URL := https://jnsaph.github.io/db-seed-runner
CHART_VERSION := $(shell yq '.version' $(HELM_CHART_DIR)/Chart.yaml)
CHART_NAME := $(shell yq '.name' $(HELM_CHART_DIR)/Chart.yaml)

build:
	go build -o bin/$(EXECUTABLE) 

run: build
	./bin/$(EXECUTABLE)

lint-helm:
	helm lint $(HELM_CHART_DIR)

check-worktree:
	@test -d $(PAGES_DIR) || (echo "Error: $(PAGES_DIR) does not exist. Please run 'git worktree add .gh-pages gh-pages'" && exit 1)

package-helm: lint-helm
	helm package $(HELM_CHART_DIR) --destination $(HELM_RELEASE_DIR)

publish-helm: check-worktree package-helm
	@if [ -f $(PAGES_DIR)/index.yaml ]; then \
		yq -e '.entries.$(CHART_NAME)[] | select(.version == "$(CHART_VERSION)")' \
			$(PAGES_DIR)/index.yaml > /dev/null && \
		echo "Error: $(CHART_NAME) version $(CHART_VERSION) already exists. Change the version in $(HELM_CHART_DIR)/Chart.yaml!" && \
		exit 1 || true; \
	fi
	cp $(HELM_RELEASE_DIR)/$(CHART_NAME)-$(CHART_VERSION).tgz $(PAGES_DIR)/
	helm repo index $(PAGES_DIR) --url $(REPO_URL) --merge $(PAGES_DIR)/index.yaml
