resource "azurerm_log_analytics_workspace" "main" {
  count = (
    var.log_analytics_workspace_name != null && var.log_analytics_workspace_id == null
    ? 1
    : 0
  )

  name                = var.log_analytics_workspace_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location
  tags                = var.tags

  sku               = var.log_analytics_workspace_sku
  retention_in_days = var.log_analytics_workspace_retention_in_days
}
