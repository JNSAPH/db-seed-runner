# DB-Seed-Runner
<center>
<img src="./assets/banner.png" alt="db-seed-runner" width="600px" />
</center>

# Quick Links
- [What is this?](#what-is-this)
- [Helm Chart Values](./charts/db-seed-runner/values.yaml)

# What is this?

In hyperscaler environments, database instances are often not publicly accessible. When provisioning infrastructure with tools like Terraform, this creates a challenge: how to seed a database with initial data when it can only be reached from within a private network.

`db-seed-runner` solves this by running database seed scripts from inside the same network as the database. It is a small, Go-based seeding runner packaged as a Helm chart that deploys a short-lived Kubernetes Job. The job connects to the database and executes the user-provided SQL seed scripts, then exits.

The tool is intended for bootstrap and reference data only and is designed to be safe to run in private VPC-based environments.


# Features
- Support for various Database Engines
    - PostgresSQL (and compatible variants like Amazon RDS, Aurora, Google Cloud SQL, etc.)
- Support for multiple SQL seed scripts
- Helm chart for easy deployment in Kubernetes environments

# Usage
## Terraform Example
Here is an example of how to use `db-seed-runner` in a Terraform configuration
```hcl
module "db_seed_runner" {
  source  = "path/to/db-seed-runner/helm-chart"
  version = "x.y.z"
}
```

# FAQ & Common Issues

## What will happen if the execution of a seed script fails?
It is possible to have multiple seed scripts. As shown in the example above, if one of the scripts fails, the Job will move on to the next script. This means that a failure in one script does not prevent subsequent scripts from executing. Make sure to work within the safeguards of the database engine you are using and to thoroughly test your seed scripts.

## I need my seed scripts to run in a specific order. How can I ensure that?
The `db-seed-runner` executes seed scripts in alphabetical order based on their filenames. To control the execution order, you can prefix your script filenames with numbers or letters (e.g., `0001_init.sql`, `0002_add_data.sql`, etc.).

## Error: `could not download chart: Chart.yaml file is missing`

This error occurs when Terraform detects a local directory with the same name as the chart (e.g., `db-seed-runner`) in the same directory where `terraform apply` is running. The Helm provider prioritizes the local directory over the remote repository, and if that directory doesn't contain a `Chart.yaml`, the error is raised.

**Solution:**
Rename the local Terraform module or directory to something different (e.g., `db-seed`) to avoid the conflict.

