resource "azurerm_service_plan" "main" {
  name                = "${var.prefix}-asp"
  resource_group_name = var.resource_group_name
  location            = var.location
  os_type             = "Linux"
  sku_name            = var.service_plan_sku
}

resource "azurerm_linux_web_app" "main" {
  name                = "${var.prefix}-appservice"
  resource_group_name = var.resource_group_name
  location            = var.location
  service_plan_id     = azurerm_service_plan.main.id
  https_only          = true

  site_config {
    minimum_tls_version = "1.2"
    always_on           = true
    application_stack {
      dotnet_version = "8.0"
    }
  }
}
