locals {
  servicebus_namespace_authorization_rule_name = (
    var.messaging_system == "queue"
    ? "queue-scaling"
    : "topic-scaling"
  )
}

resource "azurerm_servicebus_namespace" "main" {
  name                = var.servicebus_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location
  tags                = var.tags

  sku = "Standard"
}

resource "azurerm_role_assignment" "servicebus" {
  scope                = azurerm_servicebus_namespace.main.id
  role_definition_name = "Azure Service Bus Data Owner"
  principal_id         = azurerm_user_assigned_identity.main.principal_id
}

resource "azurerm_servicebus_queue" "create" {
  count = var.messaging_system == "queue" ? 1 : 0

  name         = "create"
  namespace_id = azurerm_servicebus_namespace.main.id
}

resource "azurerm_servicebus_topic" "create" {
  count = var.messaging_system == "pubsub" ? 1 : 0

  name         = "create"
  namespace_id = azurerm_servicebus_namespace.main.id
}

resource "azurerm_servicebus_namespace_authorization_rule" "scaling" {
  name         = local.servicebus_namespace_authorization_rule_name
  namespace_id = azurerm_servicebus_namespace.main.id

  manage = true
  listen = true
  send   = true
}
