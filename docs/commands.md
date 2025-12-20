# Commands
## Add Repo to Helm
```bash
helm repo add db-seed-runner https://jnsaph.github.io/db-seed-runner
helm repo update
```

## Show Available Versions
```bash
helm search repo db-seed-runner/db-seed-runner --devel
```

## Login to GHCR
Set an environment variable `GHCR_PAT` with a GitHub Personal Access Token that has `read:packages` and `write:packages` scopes.

```bash
echo $GHCR_PAT | docker login ghcr.io -u jnsaph --password-stdin
```