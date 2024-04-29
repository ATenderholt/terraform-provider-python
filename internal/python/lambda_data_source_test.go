package python_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"testing"
)

const basicExample = `
provider "python" {
  pip_command = "pip3.10"
}

data "python_aws_lambda" "test" {
  source_dir  = "test-fixtures/example_without_deps"
  archive_path = "output/example_without_deps.zip"
}
`

func TestAccAwsLambda_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:                false,
		PreCheck:                  nil,
		ProtoV6ProviderFactories:  protoV6ProviderFactories(),
		PreventPostDestroyRefresh: false,
		CheckDestroy:              nil,
		ErrorCheck:                nil,
		Steps: []resource.TestStep{
			{
				Config: basicExample,
				Check: resource.ComposeAggregateTestCheckFunc(
					testFileExists("output/example_without_deps.zip"),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "archive_base64sha256", hexToBase64("bef7f0ecaa3caa9168df5c4845da02e9d16a033875b92c3bb64cd78a1afc3448")),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "dependencies_base64sha256", ""),
				),
			},
		},
		WorkingDir: "",
	})
}

const example = `
provider "python" {
  pip_command = "pip3.11"
}

data "python_aws_lambda" "test" {
  source_dir        = "test-fixtures/example"
  archive_path      = "output/example.zip"
  dependencies_path = "output/example_deps.zip"
}
`

func TestAccAwsLambda_WithDependencies(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:                false,
		PreCheck:                  nil,
		ProtoV6ProviderFactories:  protoV6ProviderFactories(),
		PreventPostDestroyRefresh: false,
		CheckDestroy:              nil,
		ErrorCheck:                nil,
		Steps: []resource.TestStep{
			{
				Config: example,
				Check: resource.ComposeAggregateTestCheckFunc(
					testFileExists("output/example.zip"),
					testFileExists("output/example_deps.zip"),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "archive_base64sha256", hexToBase64("842611c6d40cc437abda689b68204416172152e5b70072d7a681e510ca08f40f")),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "dependencies_base64sha256", hexToBase64("2a4ba9a6524ed60b9aa69fdb49b300030e4cee0d45474d86100b0a6551bf9571")),
				),
			},
		},
		WorkingDir: "",
	})
}

func testFileExists(path string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := os.Stat(path)
		if err != nil {
			return err
		}
		return nil
	}
}
