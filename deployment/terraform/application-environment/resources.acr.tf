resource "azurerm_container_registry" "main" {
  name                = var.container_registry_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location
  tags                = var.tags

  sku = "Basic"
}

resource "azurerm_role_assignment" "acr" {
  scope                = azurerm_container_registry.main.id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_user_assigned_identity.main.principal_id
}
