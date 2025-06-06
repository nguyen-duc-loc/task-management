resource "random_string" "rd" {
  length  = 6
  special = false
  upper   = false
  numeric = false
}

resource "aws_secretsmanager_secret" "secrets" {
  name = "${var.app_name}-${random_string.rd.result}"
}

resource "random_string" "db_password" {
  length  = 32
  special = false
}

resource "aws_secretsmanager_secret_version" "secrets_version" {
  secret_id = aws_secretsmanager_secret.secrets.id
  secret_string = jsonencode({
    JWT_SECRET_KEY = var.jwt_secret_key
    DB_PASSWORD    = random_string.db_password.result
  })
}
