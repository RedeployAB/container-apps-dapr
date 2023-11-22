locals {
  worker_scale_queue_settings = {
    name = "queue-scale-rule"
    metadata = {
      namespace    = azurerm_servicebus_namespace.main.name
      queueName    = "create"
      messageCount = "1"
    }
  }

  worker_scale_pubsub_settings = {
    name = "topic-scale-rule"
    metadata = {
      namespace        = azurerm_servicebus_namespace.main.name
      subscriptionName = var.worker_container_app.name
      topicName        = "create"
    }
  }

  worker_scale_settings = (
    var.messaging_system == "queue"
    ? local.worker_scale_queue_settings
    : local.worker_scale_pubsub_settings
  )
}

resource "azurerm_container_app_environment" "main" {
  name                = var.container_app_environment_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location
  tags                = var.tags

  log_analytics_workspace_id = local.log_analytics_workspace_id
}

resource "azurerm_container_app_environment_dapr_component" "queue" {
  count = var.messaging_system == "queue" ? 1 : 0

  name                         = "reports"
  container_app_environment_id = azurerm_container_app_environment.main.id
  component_type               = "bindings.azure.servicebusqueues"
  version                      = "v1"
  scopes                       = var.messaging_dapr_scopes

  metadata {
    name  = "namespaceName"
    value = "${azurerm_servicebus_namespace.main.name}.servicebus.windows.net"
  }

  metadata {
    name  = "queueName"
    value = "create"
  }

  metadata {
    name  = "azureClientId"
    value = azurerm_user_assigned_identity.main.client_id
  }
}

resource "azurerm_container_app_environment_dapr_component" "pubsub" {
  count = var.messaging_system == "pubsub" ? 1 : 0

  name                         = "reports"
  container_app_environment_id = azurerm_container_app_environment.main.id
  component_type               = "pubsub.azure.servicebus"
  version                      = "v1"
  scopes                       = var.messaging_dapr_scopes

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

resource "azurerm_container_app" "endpoint" {
  count = var.provision_container_apps ? 1 : 0

  name                         = var.endpoint_container_app.name
  resource_group_name          = local.resource_group.name
  container_app_environment_id = azurerm_container_app_environment.main.id
  revision_mode                = "Single"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.main.id
    ]
  }

  registry {
    server   = azurerm_container_registry.main.login_server
    identity = azurerm_user_assigned_identity.main.id
  }

  dapr {
    app_id       = var.endpoint_container_app.name
    app_port     = var.endpoint_container_app.port
    app_protocol = "grpc"
  }

  ingress {
    external_enabled = true
    target_port      = var.endpoint_container_app.port
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  secret {
    name  = "endpoint-security-keys"
    value = join(",", var.endpoint_security_keys)
  }

  template {
    container {
      name   = var.endpoint_container_app.name
      image  = "${azurerm_container_registry.main.login_server}/${var.endpoint_container_app.image}"
      cpu    = var.endpoint_container_app.cpu
      memory = var.endpoint_container_app.memory

      env {
        name        = "ENDPOINT_SECURITY_KEYS"
        secret_name = "endpoint-security-keys"
      }

      env {
        name  = "ENDPOINT_REPORTER_TYPE"
        value = var.messaging_system
      }
    }

    min_replicas = var.endpoint_container_app.min_replicas
    max_replicas = var.endpoint_container_app.max_replicas

    http_scale_rule {
      name                = "http-scale-rule"
      concurrent_requests = 50
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].container[0].image
    ]
  }
}

resource "azurerm_container_app" "worker" {
  count = var.provision_container_apps ? 1 : 0

  name                         = var.worker_container_app.name
  resource_group_name          = local.resource_group.name
  container_app_environment_id = azurerm_container_app_environment.main.id
  revision_mode                = "Single"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.main.id
    ]
  }

  registry {
    server   = azurerm_container_registry.main.login_server
    identity = azurerm_user_assigned_identity.main.id
  }

  dapr {
    app_id       = var.worker_container_app.name
    app_port     = var.worker_container_app.port
    app_protocol = "grpc"
  }

  secret {
    name  = "servicebus-connection-string"
    value = azurerm_servicebus_namespace_authorization_rule.scaling.primary_connection_string
  }

  template {
    container {
      name   = var.worker_container_app.name
      image  = "${azurerm_container_registry.main.login_server}/${var.worker_container_app.image}"
      cpu    = var.worker_container_app.cpu
      memory = var.worker_container_app.memory

      env {
        name  = "WORKER_TYPE"
        value = var.messaging_system
      }
    }

    min_replicas = var.worker_container_app.min_replicas
    max_replicas = var.worker_container_app.max_replicas

    custom_scale_rule {
      name             = local.worker_scale_settings.name
      custom_rule_type = "azure-servicebus"
      metadata         = local.worker_scale_settings.metadata
      authentication {
        secret_name       = "servicebus-connection-string"
        trigger_parameter = "connection"
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].container[0].image
    ]
  }
}
