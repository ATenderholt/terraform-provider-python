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
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "archive_base64sha256", "sX8w3367kZUgdBbFFW0i0LQ/2zJ9QuVvm7TokKa2vto="),
				),
			},
		},
		WorkingDir: "",
	})
}

const example = `
provider "python" {
  pip_command = "pip3.10"
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
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "archive_base64sha256", hexToBase64("a29ac37520504756fed4e3d98f5a8ecbde3b56f81a7cfa0ddcb3ddecdffb1deb")),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "dependencies_base64sha256", hexToBase64("aa")),
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
