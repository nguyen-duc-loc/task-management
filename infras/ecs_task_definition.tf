resource "aws_ecs_task_definition" "backend_task_definition" {
  family                   = "${var.app_name}-backend-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024"
  memory                   = "3072"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([
    {
      name         = "${var.app_name}-backend"
      image        = var.backend_image
      portMappings = [{ containerPort = var.backend_port, hostPort = var.backend_port }]
      environment = [
        {
          name  = "SERVER_PORT"
          value = tostring(var.backend_port)
        },
        {
          name  = "SERVER_ENV"
          value = "prod"
        },
        {
          name  = "DB_HOST"
          value = aws_db_instance.db_instance.address
        },
        {
          name  = "DB_DATABASE"
          value = aws_db_instance.db_instance.db_name
        },
        {
          name  = "DB_USERNAME"
          value = aws_db_instance.db_instance.username
        },
        {
          name  = "DB_SCHEMA"
          value = "public"
        },
        {
          name  = "DB_PORT"
          value = tostring(var.db_port)
        },
        {
          name  = "SECRETS_MANAGER_NAME"
          value = aws_secretsmanager_secret.secrets.name
        },
        {
          name  = "JWT_ACCESS_TOKEN_DURATION"
          value = "720h"
        }
      ]
    }
  ])
}

resource "aws_ecs_task_definition" "frontend_task_definition" {
  family                   = "${var.app_name}-frontend-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024"
  memory                   = "3072"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([
    {
      name         = "${var.app_name}-frontend"
      image        = var.frontend_image
      portMappings = [{ containerPort = var.frontend_port, hostPort = var.frontend_port }]
      environment = [
        {
          name  = "API_BASE_URL"
          value = "http://${aws_lb.backend_lb.dns_name}:${var.backend_port}"
        },
        {
          name  = "PORT"
          value = tostring(var.frontend_port)
        }
      ]
    }
  ])
}
