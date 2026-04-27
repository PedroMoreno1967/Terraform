resource "azurerm_mssql_server" "main" {
  name                          = "${var.prefix}-sqlsvr"
  resource_group_name           = var.resource_group_name
  location                      = var.location
  version                       = "12.0"
  administrator_login           = var.administrator_login
  administrator_login_password  = var.administrator_login_password
  public_network_access_enabled = true
  minimum_tls_version           = "1.2"
}

resource "azurerm_mssql_database" "main" {
  name      = "${var.prefix}-db"
  server_id = azurerm_mssql_server.main.id
  sku_name  = "Basic"
  collation = "SQL_Latin1_General_CP1_CI_AS"
}

resource "azurerm_mssql_firewall_rule" "allow_azure" {
  name             = "allow-azure-services"
  server_id        = azurerm_mssql_server.main.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}
