variable "deploy_resource_group" {
  type        = bool
  default     = true
  description = "Deploy a resource group or to an existing one."
}

variable "resource_group_name" {
  type        = string
  description = "Name of the resource group to deploy the application environment to."
}

variable "location" {
  type        = string
  default     = "westeurope"
  description = "Loction for the application environment."
}

variable "identity_name" {
  type        = string
  description = "Name of the user assigned identity."
}

variable "log_analytics_workspace_name" {
  type        = string
  default     = null
  description = "Name of the log analytics workspace. Use to deploy a new workspace."
}

variable "log_analytics_workspace_id" {
  type        = string
  default     = null
  description = "ID of an existing log analytics workspace. Use to deploy to an existing workspace."
}

variable "log_analytics_workspace_sku" {
  type        = string
  default     = "PerGB2018"
  description = "SKU of the log analytics workspace."
}

variable "log_analytics_workspace_retention_in_days" {
  type        = string
  default     = 30
  description = "Retention in days of the log analytics workspace."
}

variable "storage_account_name" {
  type        = string
  description = "Name of storage account for output binding."
}

variable "storage_account_tier" {
  type        = string
  default     = "Standard"
  description = "Access tier of the storage account."
}

variable "storage_account_replication_type" {
  type        = string
  default     = "LRS"
  description = "Replication type of the storage account."
}

variable "storage_container_name" {
  type        = string
  default     = "reports"
  description = "Name of the storage container for output binding."
}

variable "servicebus_name" {
  type        = string
  description = "Name of the service bus for pubsub binding."
}

variable "servicebus_namespace_authorization_rule_name" {
  type        = string
  default     = "pubsub-scaling"
  description = "Name of the service bus namespace authorization rule for pubsub binding scaling rules."
}

variable "container_registry_name" {
  type        = string
  description = "Name of the container registry for the container app environment."
}

variable "container_app_environment_name" {
  type        = string
  description = "Name of the container app environment."
}

variable "pubsub_dapr_scopes" {
  type        = list(string)
  description = "Comma separated list of container apps (names) for the pubsub binding."
}

variable "output_dapr_scopes" {
  type        = list(string)
  description = "Comma separated list of container apps (names) for the output binding."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Tags to apply to all resources."
}
