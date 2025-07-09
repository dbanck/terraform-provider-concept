package schema

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func PetIdentitySchema() *tfprotov6.ResourceIdentitySchema {
	return &tfprotov6.ResourceIdentitySchema{
		IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
			{
				Name:              "id",
				Description:       "The random pet name.",
				Type:              tftypes.String,
				RequiredForImport: true,
			},
			{
				Name:              "legs",
				Description:       "Number of legs the pet has.",
				Type:              tftypes.Number,
				OptionalForImport: true,
			},
		},
	}
}
