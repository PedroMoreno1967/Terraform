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

variable "admin_object_id" {
  type = string
}
