resource "azurerm_container_app_environment" "main" {
  name                = var.container_app_environment_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location

  log_analytics_workspace_id = local.log_analytics_workspace_id

  tags = var.tags
}

resource "azurerm_container_app_environment_dapr_component" "pubsub" {
  name                         = "reports"
  container_app_environment_id = azurerm_container_app_environment.main.id
  component_type               = "pubsub.azure.servicebus"
  version                      = "v1"
  scopes                       = var.pubsub_dapr_scopes

  metadata {
    name  = "namespaceName"
    value = "${azurerm_servicebus_namespace.main.name}.servicebus.windows.net"
  }

  metadata {
    name  = "azureClientId"
    value = azurerm_user_assigned_identity.main.client_id
  }
}

resource "azurerm_container_app_environment_dapr_component" "output" {
  name                         = "reports-output"
  container_app_environment_id = azurerm_container_app_environment.main.id
  component_type               = "bindings.azure.blobstorage"
  version                      = "v1"
  scopes                       = var.output_dapr_scopes

  metadata {
    name  = "accountName"
    value = azurerm_storage_account.main.name
  }

  metadata {
    name  = "containerName"
    value = var.storage_container_name
  }

  metadata {
    name  = "azureClientId"
    value = azurerm_user_assigned_identity.main.client_id
  }
}
