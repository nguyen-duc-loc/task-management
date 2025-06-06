terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "task-management-tttn"
    key = "task-management-state-file"
    region = "ap-southeast-1"
  }
}

provider "aws" {
  region = var.region
}