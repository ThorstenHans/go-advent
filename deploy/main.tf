provider "azurerm" {
  features {}
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "2.85.0"
    }
  }
}

resource "azurerm_resource_group" "main" {
  name     = "rg-cloud-automation"
  location = "westeurope"
  tags = {
    "com.thorsten-hans.keeper" : ""
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_role_assignment" "contrib" {
  scope                = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  role_definition_name = "Contributor"
  principal_id         = azurerm_function_app.fn.identity[0].principal_id
}
