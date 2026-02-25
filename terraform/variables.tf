variable "project_name" {
  type        = string
  description = "Project name for resource naming and tagging"
  default     = "article-manager"
}

variable "environment" {
  type        = string
  description = "Environment name"
  default     = "production"
}

variable "aws_region" {
  type        = string
  description = "AWS region to deploy resources"
  default     = "ap-northeast-1"
}

variable "vpc_cidr" {
  type        = string
  description = "CIDR block for VPC"
  default     = "10.0.0.0/16"
}

# Phase 3で使用開始: RDS・Secrets Manager構築時に有効化
# variable "db_master_password" {
#   type        = string
#   description = "Master password for RDS MySQL database"
#   sensitive   = true
# }

# variable "gemini_api_key" {
#   type        = string
#   description = "Google Gemini API key for AI features"
#   sensitive   = true
# }

# variable "google_books_api_key" {
#   type        = string
#   description = "Google Books API key for book recommendations"
#   sensitive   = true
# }

# Phase 6で使用開始: Route 53設定時に有効化
# variable "domain_name" {
#   type        = string
#   description = "Domain name for the application"
#   default     = ""
# }
