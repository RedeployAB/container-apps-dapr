resource "azurerm_storage_account" "main" {
  name                = var.storage_account_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location
  tags                = var.tags

  account_tier             = var.storage_account_tier
  account_replication_type = var.storage_account_replication_type

  enable_https_traffic_only       = true
  allow_nested_items_to_be_public = false
}

resource "azurerm_role_assignment" "output" {
  scope                = azurerm_storage_account.main.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.main.principal_id
}

resource "azurerm_storage_container" "output" {
  name                  = var.storage_container_name
  storage_account_name  = azurerm_storage_account.main.name
  container_access_type = "private"
}
