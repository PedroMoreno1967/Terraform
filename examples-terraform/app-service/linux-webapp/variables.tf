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

variable "service_plan_sku" {
  type    = string
  default = "S1"
}
