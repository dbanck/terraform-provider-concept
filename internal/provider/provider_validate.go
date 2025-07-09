package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func (p ConceptProvider) ValidateListResourceConfig(ctx context.Context, request *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	return &tfprotov6.ValidateListResourceConfigResponse{}, nil
}

func (p ConceptProvider) ValidateResourceConfig(ctx context.Context, request *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	return &tfprotov6.ValidateResourceConfigResponse{}, nil
}

func (p ConceptProvider) ValidateProviderConfig(ctx context.Context, request *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	return &tfprotov6.ValidateProviderConfigResponse{}, nil
}

func (p ConceptProvider) ValidateDataResourceConfig(ctx context.Context, request *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) ValidateEphemeralResourceConfig(ctx context.Context, request *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	panic("not implemented")
}
