data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "main" {
  name                       = lower("${var.prefix}kv${substr(md5(var.resource_group_name), 0, 8)}")
  resource_group_name        = var.resource_group_name
  location                   = var.location
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 90
  purge_protection_enabled   = false
}

resource "azurerm_key_vault_access_policy" "admin" {
  key_vault_id = azurerm_key_vault.main.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = var.admin_object_id

  key_permissions         = ["Get", "List", "Create", "Delete"]
  secret_permissions      = ["Get", "List", "Set", "Delete"]
  certificate_permissions = ["Get", "List", "Create", "Import", "Delete"]
}
