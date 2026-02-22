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

# DB Parameter Group
resource "aws_db_parameter_group" "main" {
  name   = "${var.project_name}-db-parameter-group"
  family = "mysql8.0"

  # Character Set Settings
  parameter {
    name  = "character_set_server"
    value = "utf8mb4"
  }

  parameter {
    name  = "collation_server"
    value = "utf8mb4_unicode_ci"
  }

  # Full-Text Search Settings
  parameter {
    name  = "innodb_ft_min_token_size"
    value = "2"
  }

  parameter {
    name  = "ft_min_word_len"
    value = "2"
  }

  tags = {
    Name        = "${var.project_name}-db-parameter-group"
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

  # Parameter Group
  parameter_group_name = aws_db_parameter_group.main.name

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
