resource "azurerm_storage_account" "sa" {
  name                     = "sacautomatego2021"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "asp" {
  name                = "asp-cloud-automate-go"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  kind                = "functionapp"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "fn" {
  name                       = "fnapp-cloud-automate-go"
  location                   = azurerm_resource_group.main.location
  resource_group_name        = azurerm_resource_group.main.name
  app_service_plan_id        = azurerm_app_service_plan.asp.id
  storage_account_name       = azurerm_storage_account.sa.name
  storage_account_access_key = azurerm_storage_account.sa.primary_access_key

  os_type = "linux"
  app_settings = {
    "FUNCTIONS_WORKER_RUNTIME" = "custom"
    "SUBSCRIPTION_ID"          = data.azurerm_client_config.current.subscription_id
  }
  version = "~4"
  site_config {
    cors {
      allowed_origins = ["https://portal.azure.com"]
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
