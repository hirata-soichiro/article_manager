# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-db-subnet-group"
  subnet_ids = [aws_subnet.private_1a.id, aws_subnet.private_1c.id]

  tags = {
    Name        = "${var.project_name}-db-subnet-group"
    Project     = var.project_name
    Environment = var.environment
  }
}

# RDS MySQL Instance
resource "aws_db_instance" "main" {
  identifier = "${var.project_name}-db"

  # Engine
  engine         = "mysql"
  engine_version = "8.0"

  # Instance
  instance_class = "db.t4g.micro"

  # Storage
  storage_type      = "gp3"
  allocated_storage = 20

  # Database
  db_name  = "article_manager"
  username = "admin"
  password = data.aws_ssm_parameter.db_admin_password.value

  # Network
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible    = false

  # Backup
  backup_retention_period = 0

  # Maintenance
  maintenance_window = "mon:04:00-mon:05:00"

  # High Availability
  multi_az = false

  # Encryption
  storage_encrypted = true

  # Deletion Protection
  skip_final_snapshot      = true
  deletion_protection      = false
  delete_automated_backups = true

  tags = {
    Name        = "${var.project_name}-db"
    Project     = var.project_name
    Environment = var.environment
  }
}
