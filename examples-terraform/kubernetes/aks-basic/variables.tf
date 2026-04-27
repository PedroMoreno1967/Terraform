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

variable "node_count" {
  type    = number
  default = 1
}

variable "vm_size" {
  type    = string
  default = "Standard_DS2_v2"
}
