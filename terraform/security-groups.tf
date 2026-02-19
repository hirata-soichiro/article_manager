# Security Group for ECS Tasks
# Allows inbound HTTP traffic and direct access to frontend/backend for testing
resource "aws_security_group" "ecs_tasks" {
  name        = "${var.project_name}-ecs-tasks-sg"
  description = "Security group for ECS Fargate tasks (Frontend + Backend)"
  vpc_id      = aws_vpc.main.id

  tags = {
    Name        = "${var.project_name}-ecs-tasks-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Inbound: HTTP (80) from anywhere
resource "aws_security_group_rule" "ecs_http_inbound" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ecs_tasks.id
  description       = "Allow HTTP traffic from anywhere"
}

# Inbound: Frontend (3000) from anywhere (optional, for direct access)
resource "aws_security_group_rule" "ecs_frontend_inbound" {
  type              = "ingress"
  from_port         = 3000
  to_port           = 3000
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ecs_tasks.id
  description       = "Allow direct access to Next.js frontend"
}

# Inbound: Backend (8080) from anywhere (optional, for direct access)
resource "aws_security_group_rule" "ecs_backend_inbound" {
  type              = "ingress"
  from_port         = 8080
  to_port           = 8080
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ecs_tasks.id
  description       = "Allow direct access to Go API backend"
}

# Outbound: All traffic
resource "aws_security_group_rule" "ecs_outbound" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ecs_tasks.id
  description       = "Allow all outbound traffic"
}

# Security Group for RDS
# Allows MySQL access only from ECS tasks
resource "aws_security_group" "rds" {
  name        = "${var.project_name}-rds-sg"
  description = "Security group for RDS MySQL database"
  vpc_id      = aws_vpc.main.id

  tags = {
    Name        = "${var.project_name}-rds-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Inbound: MySQL (3306) from ECS tasks only
resource "aws_security_group_rule" "rds_mysql_inbound" {
  type                     = "ingress"
  from_port                = 3306
  to_port                  = 3306
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.ecs_tasks.id
  security_group_id        = aws_security_group.rds.id
  description              = "Allow MySQL access from ECS tasks"
}
