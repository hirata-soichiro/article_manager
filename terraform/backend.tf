terraform {
  backend "s3" {
    bucket = "article-manager-terraform-state"
    key    = "terraform.tfstate"
    region = "ap-northeast-1"
    use_lockfile = true
  }
}
