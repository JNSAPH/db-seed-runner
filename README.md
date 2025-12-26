# DB-Seed-Runner

## Overview

DB-Seed-Runner is a Helm chart for deploying a Go-based PostgreSQL seeding runner. It is designed to help automate the process of seeding PostgreSQL databases in Kubernetes environments.


> **Note:** This chart is primarily intended to be used with Terraform for automated deployments as part of your infrastructure-as-code workflows. For full documentation and advanced usage, see the [README.md in the main branch](https://github.com/jnsaph/db-seed-runner/blob/main/README.md).

## Prerequisites

- Kubernetes 1.19+
- Helm 3.x
- A running PostgreSQL instance accessible from your cluster


## Installation

### Using Terraform


You can use the [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs) in Terraform to deploy this chart as part of your infrastructure code. Example configuration:

```hcl
# ConfigMap containing your SQL seed scripts
resource "kubernetes_config_map" "db_seed_sql" {
	metadata {
		name = "db-seed-sql"
	}
	data = {
		"0001_payload_cms.sql" = file("${path.module}/0001_payload_cms.sql")
		# Add your SQL files here. The key is the filename, the value is the file content.
		# The order of execution is determined by the key's alphabetical order.
		"0001_init.sql" = file("${path.module}/0001_init.sql")
		# "0002_more_data.sql" = file("${path.module}/0002_more_data.sql")
	}
}

# Helm release to deploy db-seed-runner
resource "helm_release" "db-seed-runner" {
	name       = "db-seed-runner"
	repository = "https://jnsaph.github.io/db-seed-runner"
	chart      = "db-seed-runner"
	version    = "0.0.7"

	# Ensure the ConfigMap is created before the Helm release
	depends_on = [kubernetes_config_map.db_seed_sql]

	# These options ensure the job is re-run if the config map changes
	force_update   = true
	recreate_pods  = true

	values = [
		yamlencode({
			seed = {
				sqlConfigMapName = kubernetes_config_map.db_seed_sql.metadata[0].name
				checksum         = kubernetes_config_map.db_seed_sql.metadata[0].generation # Forces job rerun on config change
			}

			database = {
				engine   = "postgres" # Supported: postgres
				host     = var.db_host
				port     = var.db_port
				dbname   = "postgres" # Change as needed, ideally use the root/default database and user
				user     = var.db_username
				password = var.db_password
			}
		})
	]
}
```

### Manual Helm Installation

Add the Helm repository:

```sh
helm repo add db-seed-runner https://jnsaph.github.io/db-seed-runner/
helm repo update
```

Install the chart:

```sh
helm install my-seed-runner db-seed-runner/db-seed-runner
```

## Upgrading

To upgrade your release:

```sh
helm upgrade my-seed-runner db-seed-runner/db-seed-runner
```

## Uninstalling

To uninstall the release:

```sh
helm uninstall my-seed-runner
```

## Configuration

The following table lists the configurable parameters of the chart and their default values. You can override these values using `--set key=value` or by providing a custom `values.yaml` file.

| Parameter         | Description                        | Default |
|-------------------|------------------------------------|---------|
| image.repository  | Image repository                   |         |
| image.tag         | Image tag                          |         |
| postgres.host     | PostgreSQL host                    |         |
| postgres.port     | PostgreSQL port                    |   5432  |
| postgres.user     | PostgreSQL user                    |         |
| postgres.password | PostgreSQL password                |         |
| postgres.database | PostgreSQL database                |         |

Example:

```sh
helm install my-seed-runner db-seed-runner/db-seed-runner \
	--set postgres.host=mydb.example.com \
	--set postgres.user=myuser \
	--set postgres.password=mypassword \
	--set postgres.database=mydb
```

## Chart Source

This chart is hosted at: https://jnsaph.github.io/db-seed-runner/

## License

Distributed under the MIT License. See `LICENSE` for more information.