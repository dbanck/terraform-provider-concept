package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"

	"github.com/dbanck/terraform-provider-concept/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {
	providerFn := func() tfprotov6.ProviderServer {
		return provider.ConceptProvider{}
	}

	err := tf6server.Serve("registry.terraform.io/dbanck/concept", providerFn)

	if err != nil {
		log.Fatal(err.Error())
	}
}
