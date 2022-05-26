variable "username" {
  description = "Database user name"
  type        = string
  default     = "postgres"
  sensitive   = true
}

variable "database_name" {
  description = "The name of the database for compliance data"
  type        = string
  default     = "compliance"
}

variable "storage" {
  description = "Allocated storage in gigabytes"
  type        = number
  default     = 20
}
