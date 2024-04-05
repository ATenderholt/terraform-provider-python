package python_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"testing"
)

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
				Config: config("example_without_deps"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFileExists("output/example_without_deps.zip"),
					resource.TestCheckResourceAttr("data.python_aws_lambda.test", "output_base64sha256", "383ecafd3769efb3ef66301806026e25"),
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

const configTemplate = `
provider "python" {
  pip_command = "pip3.10"
}

data "python_aws_lambda" "test" {
  source_dir  = "test-fixtures/%s"
  output_path = "output/%s.zip"
}
`

func config(name string) string {
	return fmt.Sprintf(
		configTemplate,
		name,
		name,
	)
}
