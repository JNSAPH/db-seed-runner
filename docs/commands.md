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