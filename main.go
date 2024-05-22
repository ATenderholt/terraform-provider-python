package main

import (
	"context"
	"flag"
	"github.com/ATenderholt/terraform-provider-python/internal/python"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "github.com/ATenderholt/python",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), python.New(), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
