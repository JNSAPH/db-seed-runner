# Go Build
EXECUTABLE = db-seed-runner
BUILD_DIR = bin

# Helm Packaging
HELM_CHART_DIR = charts/db-seed-runner
HELM_RELEASE_DIR = helm-releases
PAGES_DIR = .gh-pages
REPO_URL         := https://jnsaph.github.io/db-seed-runner


build:
	go build -o bin/$(EXECUTABLE) 

run: build
	./bin/$(EXECUTABLE)


check-worktree:
	@test -d $(PAGES_DIR) || (echo "Error: $(PAGES_DIR) does not exist. Please run 'git worktree add .gh-pages gh-pages'" && exit 1)

package-helm:
	helm package $(HELM_CHART_DIR) --destination $(HELM_RELEASE_DIR)

publish-helm: check-worktree package-helm
	cp $(HELM_RELEASE_DIR)/*.tgz $(PAGES_DIR)/
	helm repo index $(PAGES_DIR) --url $(REPO_URL) --merge $(PAGES_DIR)/index.yaml