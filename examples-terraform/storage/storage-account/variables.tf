variable "resource_group_name" {
  type = string
}

variable "location" {
  type = string
}

variable "storage_account_name" {
  description = "Globally unique lowercase storage account name"
  type        = string
}

variable "ip_rules" {
  description = "CIDRs allowed to access storage"
  type        = list(string)
  default     = ["23.45.1.0/30"]
}
