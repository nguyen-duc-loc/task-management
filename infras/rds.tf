resource "aws_db_subnet_group" "db_subnet_group" {
  name       = "${var.app_name}-db-subnet-group"
  subnet_ids = data.aws_subnets.default_subnets.ids
}

resource "aws_db_instance" "db_instance" {
  identifier             = "${var.app_name}-db"
  engine                 = "postgres"
  instance_class         = "db.t3.micro"
  allocated_storage      = 20
  db_name                = join("", split("-", var.app_name))
  port                   = var.db_port
  username               = var.db_username
  password               = jsondecode(aws_secretsmanager_secret_version.secrets_version.secret_string)["DB_PASSWORD"]
  skip_final_snapshot    = true
  publicly_accessible    = true
  vpc_security_group_ids = [aws_security_group.db_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.db_subnet_group.name
}
