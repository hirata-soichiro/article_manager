variable "project_name" {
  description = "Project name for resource naming and tagging"
  type        = string
  default     = "article-manager"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "ap-northeast-1"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "db_master_password" {
  description = "Master password for RDS MySQL database"
  type        = string
  sensitive   = true
}

variable "gemini_api_key" {
  description = "Google Gemini API key for AI features"
  type        = string
  sensitive   = true
}

variable "google_books_api_key" {
  description = "Google Books API key for book recommendations"
  type        = string
  sensitive   = true
}

variable "domain_name" {
  description = "Domain name for the application"
  type        = string
  default     = ""
}
