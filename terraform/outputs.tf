output "hostname" {
  description = "RDS instance hostname"
  value       = aws_db_instance.postgres_database.address
}

output "port" {
  description = "RDS instance port"
  value       = aws_db_instance.postgres_database.port
}

output "secret_name" {
  description = "RDS secret name (use AWS Secret Manager to fetch the value)"
  value       = aws_secretsmanager_secret.database_secret.name
}

output "secret_arn" {
  description = "RDS secret ARN"
  value       = aws_secretsmanager_secret_version.database_secret.arn
}
