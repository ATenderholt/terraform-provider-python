package python_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				Config: config,
			},
		},
		WorkingDir: "",
	})
}

const config = `
provider "python" {
  pip_command = "pip3.10"
}

data "python_aws_lambda" "test" {}
`
