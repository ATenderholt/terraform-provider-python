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

data "python_aws_lambda" "example" {
  source_dir        = "../../internal/python/test-fixtures/example"
  archive_path      = "output/handler.zip"
  dependencies_path = "output/dependencies.zip"
  extra_args        = "--platform=manylinux_2_17_i686 --only-binary=:all:"
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

resource "aws_lambda_function" "example" {
  function_name = "provider-datasource-example"
  role          = aws_iam_role.data_source.arn
  filename      = data.python_aws_lambda.example.archive_path
  handler       = "main.handler"
  layers        = [aws_lambda_layer_version.example.arn]
  runtime       = "python3.10"

  source_code_hash = data.python_aws_lambda.example.archive_base64sha256
}

resource "aws_lambda_layer_version" "example" {
  layer_name = "provider-datasource-example"
  filename   = data.python_aws_lambda.example.dependencies_path

  compatible_runtimes = ["python3.10"]
  source_code_hash    = data.python_aws_lambda.example.dependencies_base64sha256
}
