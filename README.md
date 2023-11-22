# container-apps-dapr

> Example project for working with Container Apps with DAPR and Terraform

This repository contains a sample project for working with Container Apps with DAPR and Terraform.

* [Getting started](#getting-started)
  * [Prerequisites](#prerequisites)
* [Provision environment](#provision-environment)
* [Build](#build)
* [Deploy](#deploy)
  * [Push images to Azure Container Registry](#push-images-to-azure-container-registry)
  * [Deploy Container Apps](#deploy-container-apps)
* [Usage](#usage)
  * [Example with curl](#example-with-curl)


## Getting started

### Prerequisites

* Azure account (with at least Contributor role on a resource group)
* Go >=1.20
* Terraform >=1.4.0
* Azure CLI >=2.45.0

Download the project:

```sh
git clone https://github.com/RedeployAB/container-apps-dapr.git
```

## Provision environment

Create a variables file for the environment like so:

```sh
touch deployments/terraform/application-environment/terraform.tfvars
```

Add the following variables to the file:

```hcl
// If creating a new resource group, set to true.
deploy_resource_group = false

resource_group_name = "<name-of-resource-group>"
location            = "<location>"
identity_name       = "<name-of-identity>"

// Set this to create a new workspace.
log_analytics_workspace_name  = "<name-of-log-analytics>"
 // Set this to deploy to an existing workspace.
log_analytics_workspace_id    = "<log-analytics-workspace-id>"

storage_account_name    = "<name-of-storage-account>"
servicebus_name         = "<name-of-servicebus-namespace>"
container_registry_name = "<name-of-azure-container-registry>"

container_app_environment_name = "<name-of-container-app-environment>"
provision_container_apps       = false

// Messaging system to use. Options are "queue" or "pubsub".
messaging_system      = "queue"
// Name of the applications in the project.
messaging_dapr_scopes = ["endpoint", "worker"]
// Name of the "worker" application of the project.
output_dapr_scopes    = ["worker"]
```

Run Terraform:

```sh
cd deployments/terraform/application-environment
terraform plan -out=tfplan

# Verify the output and apply.
terraform apply tfplan
```

## Build

```sh
# Build the endpoint binary.
./scripts/bash/build.sh --module endpoint --version 1.0.0 --image

# Build the worker binary.
./scripts/bash/build.sh --module worker --version 1.0.0 --image
```

## Deploy

### Push images to Azure Container Registry

```sh
# Login to the container registry.
az acr login --name <name-of-azure-container-registry>

# Tag the endpoint and push it to the registry.
docker tag endpoint:1.0.0 <name-of-azure-container-registry>.azurecr.io/endpoint:1.0.0
docker push <name-of-azure-container-registry>.azurecr.io/endpoint:1.0.0

# Tag the worker and push it to the registry.
docker tag worker:1.0.0 <name-of-azure-container-registry>.azurecr.io/worker:1.0.0
docker push <name-of-azure-container-registry>.azurecr.io/worker:1.0.0
```

### Deploy Container Apps

#### Terraform

Edit `deployments/terraform/application-environment/terraform.tfvars` and add:

```hcl
provision_container_apps = true
endpoint_security_keys   = ["<key>"]
```

Run Terraform:

```sh
cd deployments/terraform/application-environment
terraform plan -out=tfplan

# Verify the output and apply.
terraform apply tfplan
```

#### Script

**Note**: Use the Terraform way of provisioning instead, this script was created when
there where certain limitiations with the Terraform modules and container app scaling.

```sh
uuid=$(uuidgen)

./deployment/scripts/deploy-container-apps.sh \
  --resource-group <name-of-resource-group> \
  --environment <name-of-container-app-environment> \
  --identity <name-of-identity> \
  --registry <name-of-azure-container-registry>.azurecr.io \
  --servicebus-namespace <name-of-servicebus-namespace> \
  --messaging-system queue \  # Supports 'pubsub' deployments. Defaults to 'queue' if omitted.
  --endpoint-version 1.0.0 \
  --worker-version 1.0.0 \
  --endpoint-api-keys $uuid
```

## Usage

```http
POST /reports

{
  "id": "12345",
  "data": "testdata" // base64 encoded
}
```

### Example with curl

```sh
# Get the URL for the `endpoint` container:
url=https://$(az containerapp show \
  --resource-group <resource-group-name> \
  --name endpoint \
  --query 'properties.configuration.ingress.fqdn' \
  --output tsv
)

# The $uuid should contain either a key set in the variable endpoint_security_keys,
# or the same $uuid as was used with the script deployment.
data=$(echo 'testdata' | base64)
curl -H "X-API-Key: $uuid" $url/reports --data "{\"id\":\"12345\",\"data\":\"$data\"}"
```
