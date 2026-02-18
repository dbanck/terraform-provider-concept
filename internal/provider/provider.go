package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var _ tfprotov6.ProviderServerWithListResource = ConceptProvider{}

type ConceptProvider struct {
}

func (p ConceptProvider) GetMetadata(ctx context.Context, request *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	return &tfprotov6.GetMetadataResponse{
		Resources: []tfprotov6.ResourceMetadata{
			{
				TypeName: "concept_pet",
			},
			{
				TypeName: "concept_kitchen_sink",
			},
		},
		ListResources: []tfprotov6.ListResourceMetadata{
			{
				TypeName: "concept_pet",
			},
			{
				TypeName: "concept_kitchen_sink",
			},
		},
	}, nil
}

func (p ConceptProvider) ConfigureProvider(ctx context.Context, request *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	return &tfprotov6.ConfigureProviderResponse{}, nil
}

func (p ConceptProvider) ReadDataSource(ctx context.Context, request *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) StopProvider(ctx context.Context, request *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	return &tfprotov6.StopProviderResponse{}, nil
}
