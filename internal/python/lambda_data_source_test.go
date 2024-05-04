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

const exampleNoExtraArgs = `
provider "python" {
  pip_command = "pip3.11"
}

data "python_aws_lambda" "test" {
  source_dir        = "test-fixtures/example"
  archive_path      = "output/example.zip"
  dependencies_path = "output/example_deps.zip"
}
`

func TestAccAwsLambda_WithDependencies_NoExtraArgs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:                false,
		PreCheck:                  nil,
		ProtoV6ProviderFactories:  protoV6ProviderFactories(),
		PreventPostDestroyRefresh: false,
		CheckDestroy:              nil,
		ErrorCheck:                nil,
		Steps: []resource.TestStep{
			{
				Config: exampleNoExtraArgs,
				Check: resource.ComposeAggregateTestCheckFunc(
					testFileExists("output/example.zip"),
					testFileExists("output/example_deps.zip"),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "archive_base64sha256", hexToBase64("bb75ac39bdf74b70cc6f8770a9077598e95bf2caf44d96dca70f250759a6b160")),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "dependencies_base64sha256", hexToBase64("33c8a87069446f87db42525345ae848a5cdd957942935e5e57f2b3251880536c")),
				),
			},
		},
		WorkingDir: "",
	})
}

const exampleExtraArgs = `
provider "python" {
  pip_command = "pip3.11"
}

data "python_aws_lambda" "test" {
  source_dir        = "test-fixtures/example"
  archive_path      = "output/example.zip"
  dependencies_path = "output/example_deps_extra_args.zip"
  extra_args        = "--platform=manylinux_2_17_i686 --only-binary=:all:"
}
`

func TestAccAwsLambda_WithDependencies_ExtraArgs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:                false,
		PreCheck:                  nil,
		ProtoV6ProviderFactories:  protoV6ProviderFactories(),
		PreventPostDestroyRefresh: false,
		CheckDestroy:              nil,
		ErrorCheck:                nil,
		Steps: []resource.TestStep{
			{
				Config: exampleExtraArgs,
				Check: resource.ComposeAggregateTestCheckFunc(
					testFileExists("output/example.zip"),
					testFileExists("output/example_deps_extra_args.zip"),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "archive_base64sha256", hexToBase64("bb75ac39bdf74b70cc6f8770a9077598e95bf2caf44d96dca70f250759a6b160")),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "dependencies_base64sha256", hexToBase64("9e55e018464f2d93dfa85b2ff0daa85346b47137b346837d1c2e8ce8e5c49cca")),
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
