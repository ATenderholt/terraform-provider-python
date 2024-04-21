terraform {
  required_providers {
    python = {
      source = "registry.terraform.io/hashicorp/python"
    }
  }
}

provider "python" {
  pip_command = "pip3.10"
}

data "python_aws_lambda" "test" {
  source_dir  = "../../internal/python/test-fixtures/example"
  archive_path = "output/simple.zip"
}
