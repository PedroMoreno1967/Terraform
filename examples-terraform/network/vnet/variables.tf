variable "resource_group_name" {
  description = "Resource group name where resources are deployed"
  type        = string
}

variable "location" {
  description = "Azure region"
  type        = string
}

variable "prefix" {
  description = "Prefix for naming"
  type        = string
  default     = "pm"
}

variable "vnet_address_space" {
  description = "Address space for VNet"
  type        = list(string)
  default     = ["10.10.0.0/16"]
}

variable "subnet_address_prefix" {
  description = "Subnet address prefix"
  type        = string
  default     = "10.10.1.0/24"
}
