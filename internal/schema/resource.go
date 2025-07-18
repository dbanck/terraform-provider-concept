package schema

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func PetResourceSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Description: "The resource `concept_pet` generates random pet names. It is the necessary companion " +
				"resource to the `concept_pet` list resource.\n",
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:            "length",
					Type:            tftypes.Number,
					Description:     "The length (in words) of the pet name. Defaults to 2",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
					Computed:        true,
				},
				{
					Name:            "id",
					Type:            tftypes.String,
					Description:     "The random pet name.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Computed:        true,
				},
			},
		},
	}
}
