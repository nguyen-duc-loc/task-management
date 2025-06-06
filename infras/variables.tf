variable "region" {
  type    = string
  default = "ap-southeast-1"
}

variable "app_name" {
  type = string
}

variable "db_username" {
  type    = string
  default = "dbusername"
}

variable "jwt_secret_key" {
  type = string
}

variable "db_port" {
  type    = number
  default = 5432
}

variable "backend_port" {
  type    = number
  default = 8080
}

variable "frontend_port" {
  type    = number
  default = 80
}

variable "backend_image" {
  type = string
}

variable "frontend_image" {
  type = string
}