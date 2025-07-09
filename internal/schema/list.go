package schema

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func PetListResourceSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:        "count",
					Type:        tftypes.Number,
					Description: "How many pets this list block generates. Defaults to 5",
					Optional:    true,
				},
			},
		},
	}
}
