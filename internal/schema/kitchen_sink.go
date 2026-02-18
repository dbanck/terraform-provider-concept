package schema

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func KitchenSinkResourceSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Description: "The resource `concept_kitchen_sink` is a minimal companion " +
				"to the kitchen sink list resource.\n",
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:            "id",
					Type:            tftypes.String,
					Description:     "Unique identifier.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Computed:        true,
				},
			},
		},
	}
}

func KitchenSinkListResourceSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:        "count",
					Type:        tftypes.Number,
					Description: "How many items to generate. Defaults to 3",
					Optional:    true,
				},
				{
					Name:            "string_attr",
					Type:            tftypes.String,
					Description:     "Plain string.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name:            "number_attr",
					Type:            tftypes.Number,
					Description:     "Numeric value.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name:            "bool_attr",
					Type:            tftypes.Bool,
					Description:     "Boolean flag.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name:            "list_attr",
					Type:            tftypes.List{ElementType: tftypes.String},
					Description:     "List of strings.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name:            "set_attr",
					Type:            tftypes.Set{ElementType: tftypes.Number},
					Description:     "Set of numbers.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name:            "map_attr",
					Type:            tftypes.Map{ElementType: tftypes.String},
					Description:     "Map of string values.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name: "object_attr",
					Type: tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"name":  tftypes.String,
							"value": tftypes.Number,
						},
					},
					Description:     "Nested object with name and value.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name:            "dynamic_attr",
					Type:            tftypes.DynamicPseudoType,
					Description:     "Accepts any type.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
				{
					Name:            "sensitive_attr",
					Type:            tftypes.String,
					Description:     "A sensitive string value.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
					Sensitive:       true,
				},
			},
		},
	}
}

func KitchenSinkIdentitySchema() *tfprotov6.ResourceIdentitySchema {
	return &tfprotov6.ResourceIdentitySchema{
		IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
			{
				Name:              "id",
				Description:       "Unique identifier.",
				Type:              tftypes.String,
				RequiredForImport: true,
			},
		},
	}
}
