package provider

import (
	"context"

	"github.com/dbanck/terraform-provider-concept/internal/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func (p ConceptProvider) GetResourceIdentitySchemas(ctx context.Context, request *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	return &tfprotov6.GetResourceIdentitySchemasResponse{
		IdentitySchemas: map[string]*tfprotov6.ResourceIdentitySchema{
			"concept_pet":          schema.PetIdentitySchema(),
			"concept_kitchen_sink": schema.KitchenSinkIdentitySchema(),
		},
	}, nil
}

func (p ConceptProvider) GetProviderSchema(ctx context.Context, request *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	return &tfprotov6.GetProviderSchemaResponse{
		Provider: &tfprotov6.Schema{
			Block: &tfprotov6.SchemaBlock{
				Description:     "Concept Provider",
				DescriptionKind: tfprotov6.StringKindPlain,
			},
		},
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"concept_pet":          schema.PetResourceSchema(),
			"concept_kitchen_sink": schema.KitchenSinkResourceSchema(),
		},
		ListResourceSchemas: map[string]*tfprotov6.Schema{
			"concept_pet":          schema.PetListResourceSchema(),
			"concept_kitchen_sink": schema.KitchenSinkListResourceSchema(),
		},
	}, nil
}
