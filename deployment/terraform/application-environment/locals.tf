locals {
  resource_group = (
    var.deploy_resource_group
    ? azurerm_resource_group.main[0]
    : data.azurerm_resource_group.main[0]
  )
}

locals {
  log_analytics_workspace_id = (
    var.log_analytics_workspace_name != null && var.log_analytics_workspace_id == null
    ? azurerm_log_analytics_workspace.main[0].id
    : var.log_analytics_workspace_id
  )
}
