variable "resource_group_name" {
  type = string
}

variable "location" {
  type = string
}

variable "prefix" {
  type    = string
  default = "pm"
}

variable "administrator_login" {
  type    = string
  default = "sqladminuser"
}

variable "administrator_login_password" {
  type      = string
  sensitive = true
}
