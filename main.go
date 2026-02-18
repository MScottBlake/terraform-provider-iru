package main

import (
	"context"
	"flag"
	"log"

	"github.com/MScottBlake/terraform-provider-iru/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run the docs generation tool, check its documentation for more information on how it works:
// http://github.com/hashicorp/terraform-plugin-docs
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/MScottBlake/iru",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New("0.1.0"), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
