package provider

import (
	"context"
	"fmt"
	"math/big"
	"math/rand/v2"
	"strings"

	"github.com/dbanck/terraform-provider-concept/internal/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	petname "github.com/dustinkirkland/golang-petname"
)

func (p ConceptProvider) ListResource(ctx context.Context, request *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	petname.NonDeterministicMode()

	results := func(push func(tfprotov6.ListResourceResult) bool) {
		if request.TypeName == "concept_pet" {
			configValue, err := request.Config.Unmarshal(schema.PetListResourceSchema().ValueType())
			if err != nil {
				panic("Failed to unmarshal config")
			}
			config := map[string]tftypes.Value{}
			err = configValue.As(&config)
			if err != nil {
				panic(err)
			}
			count := big.NewFloat(0)
			err = config["count"].As(&count)
			if err != nil {
				// Ignore for now
			}
			var countInt int64 = 5 // defaultValue
			if count != nil {
				// If count is set, use it
				countInt, _ = count.Int64()
			}

			for _ = range countInt {
				pet := strings.ToLower(petname.Generate(2, "-"))
				result := tfprotov6.ListResourceResult{
					DisplayName: fmt.Sprintf("This is a %s", pet),
					Identity:    makePetIdentity(pet),
				}

				if request.IncludeResource {
					resourceObject, err := makePetResource(pet, 2)
					if err != nil {
						panic(err)
					}
					result.Resource = &resourceObject
				}

				if !push(result) {
					return
				}

			}
		}
	}

	return &tfprotov6.ListResourceServerStream{
		Results: results,
	}, nil
}

func makePetResource(name string, length int) (tfprotov6.DynamicValue, error) {
	typ := schema.PetResourceSchema().ValueType()
	value := map[string]tftypes.Value{
		"id":     tftypes.NewValue(tftypes.String, name),
		"length": tftypes.NewValue(tftypes.Number, length),
	}

	return tfprotov6.NewDynamicValue(typ, tftypes.NewValue(typ, value))
}

func makePetIdentity(name string) *tfprotov6.ResourceIdentityData {
	legs := rand.IntN(6) + 2
	typ := schema.PetIdentitySchema().ValueType()
	value := map[string]tftypes.Value{
		"id":   tftypes.NewValue(tftypes.String, name),
		"legs": tftypes.NewValue(tftypes.Number, legs),
	}

	identiyVal, err := tfprotov6.NewDynamicValue(typ, tftypes.NewValue(typ, value))
	if err != nil {
		panic(err)
	}

	return &tfprotov6.ResourceIdentityData{
		IdentityData: &identiyVal,
	}
}
