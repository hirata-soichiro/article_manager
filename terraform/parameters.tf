# 機密情報（SecureString）- Parameter Storeから参照
data "aws_ssm_parameter" "db_admin_password" {
  name = "/article-manager/db/admin-password"
}

data "aws_ssm_parameter" "db_app_password" {
  name = "/article-manager/db/app-password"
}

data "aws_ssm_parameter" "gemini_api_key" {
  name = "/article-manager/api/gemini-api-key"
}

data "aws_ssm_parameter" "google_books_api_key" {
  name = "/article-manager/api/google-books-api-key"
}

# 非機密情報 - Terraformで作成
resource "aws_ssm_parameter" "db_name" {
  name  = "/article-manager/db/name"
  type  = "String"
  value = "article_manager"

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_ssm_parameter" "db_admin_user" {
  name  = "/article-manager/db/admin-user"
  type  = "String"
  value = "admin"

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_ssm_parameter" "db_app_user" {
  name  = "/article-manager/db/app-user"
  type  = "String"
  value = "article_user"

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}
