resource "azurerm_servicebus_namespace" "main" {
  name                = var.servicebus_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location
  tags                = var.tags

  sku = "Standard"
}

resource "azurerm_servicebus_topic" "create" {
  name         = "create"
  namespace_id = azurerm_servicebus_namespace.main.id
}

resource "azurerm_role_assignment" "pubsub" {
  scope                = azurerm_servicebus_namespace.main.id
  role_definition_name = "Azure Service Bus Data Owner"
  principal_id         = azurerm_user_assigned_identity.main.principal_id
}

resource "azurerm_servicebus_namespace_authorization_rule" "pubsub_scaling" {
  name         = var.servicebus_namespace_authorization_rule_name
  namespace_id = azurerm_servicebus_namespace.main.id

  manage = true
  listen = true
  send   = true
}
