terraform {
  required_providers {
    python = {
      source = "registry.terraform.io/hashicorp/python"
    }
    aws = {
      source = "registry.terraform.io/hashicorp/aws"
      version = "~> 5.46.0"
    }
  }
}

provider "python" {
  pip_command = "pip3.10"
}

provider "aws" {
  region  = "us-east-1"
}

data "python_aws_lambda" "simple" {
  source_dir   = "../../internal/python/test-fixtures/example_without_deps"
  archive_path = "output/simple.zip"
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "data_source" {
  name               = "provider-datasource"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_lambda_function" "data_source_simple" {
  function_name = "provider-datasource-simple"
  role          = aws_iam_role.data_source.arn
  filename      = data.python_aws_lambda.simple.archive_path
  handler       = "main.handler"
  runtime       = "python3.10"

  source_code_hash = data.python_aws_lambda.simple.archive_base64sha256
}

output "hash" {
  value = aws_lambda_function.data_source_simple.source_code_hash
}
