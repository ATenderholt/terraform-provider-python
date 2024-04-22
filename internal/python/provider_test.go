package python_test

import (
	"github.com/ATenderholt/terraform-provider-python/internal/python"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"python": providerserver.NewProtocol6WithError(python.New()()),
	}
}
