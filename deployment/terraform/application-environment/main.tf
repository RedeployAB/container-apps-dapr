terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.81.0"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "main" {
  count = var.deploy_resource_group ? 0 : 1

  name = var.resource_group_name
}

resource "azurerm_resource_group" "main" {
  count = var.deploy_resource_group ? 1 : 0

  name     = var.resource_group_name
  location = var.location
  tags     = var.tags
}

resource "azurerm_user_assigned_identity" "main" {
  name                = var.identity_name
  resource_group_name = local.resource_group.name
  location            = local.resource_group.location
  tags                = var.tags
}
