# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.project_name}-cluster"

  tags = {
    Name        = "${var.project_name}-cluster"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Backend Task Definition
resource "aws_ecs_task_definition" "api" {
  family                   = "${var.project_name}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "api"
      image = "${aws_ecr_repository.backend.repository_url}:latest"

      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_app.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "api"
        }
      }

      environment = [
        {
          name  = "DB_HOST"
          value = aws_db_instance.main.address
        },
        {
          name  = "PORT"
          value = "8080"
        }
      ]

      secrets = [
        {
          name      = "DB_NAME"
          valueFrom = aws_ssm_parameter.db_name.arn
        },
        {
          name      = "DB_USER"
          valueFrom = aws_ssm_parameter.db_app_user.arn
        },
        {
          name      = "DB_PASSWORD"
          valueFrom = data.aws_ssm_parameter.db_app_password.arn
        },
        {
          name      = "GEMINI_API_KEY"
          valueFrom = data.aws_ssm_parameter.gemini_api_key.arn
        },
        {
          name      = "GOOGLE_BOOKS_API_KEY"
          valueFrom = data.aws_ssm_parameter.google_books_api_key.arn
        }
      ]

      essential = true
    }
  ])

  tags = {
    Name        = "${var.project_name}-api-task"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Frontend Task Definition
resource "aws_ecs_task_definition" "frontend" {
  family                   = "${var.project_name}-frontend"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "frontend"
      image = "${aws_ecr_repository.frontend.repository_url}:latest"

      portMappings = [
        {
          containerPort = 3000
          protocol      = "tcp"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_app.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "frontend"
        }
      }

      environment = [
        {
          name  = "NEXT_PUBLIC_API_URL"
          value = "http://localhost:8080"
        }
      ]

      essential = true
    }
  ])

  tags = {
    Name        = "${var.project_name}-frontend-task"
    Environment = var.environment
    Project     = var.project_name
  }
}
