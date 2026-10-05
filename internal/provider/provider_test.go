package provider_test

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/dbanck/terraform-provider-concept/internal/provider"
)

// testProviderFactories runs the provider in-process for terraform-plugin-testing.
//
// The provider doesn't talk to any remote API, so all tests that use these
// factories run with resource.UnitTest and don't need TF_ACC to be set. They
// do require a Terraform binary, either on the PATH or via
// TF_ACC_TERRAFORM_PATH.
var testProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"concept": func() (tfprotov6.ProviderServer, error) {
		return provider.ConceptProvider{}, nil
	},
}
