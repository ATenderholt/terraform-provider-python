package main

import (
	"context"
	"flag"
	"github.com/ATenderholt/terraform-provider-python-package/python"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "github.com/ATenderholt/python_package",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), python.New(), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
