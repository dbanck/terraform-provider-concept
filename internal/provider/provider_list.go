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
		if request.TypeName == "concept_kitchen_sink" {
			configValue, err := request.Config.Unmarshal(schema.KitchenSinkListResourceSchema().ValueType())
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
			var countInt int64 = 3 // defaultValue
			if count != nil {
				countInt, _ = count.Int64()
			}

			for i := range countInt {
				id := fmt.Sprintf("ks-%d", i)
				result := tfprotov6.ListResourceResult{
					DisplayName: fmt.Sprintf("Kitchen sink item %s", id),
					Identity:    makeKitchenSinkIdentity(id),
				}

				if request.IncludeResource {
					resourceObject, err := makeKitchenSinkResource(id)
					if err != nil {
						panic(err)
					}
					result.Resource = &resourceObject
				}

				if !push(result) {
					return
				}
			}
		} else if request.TypeName == "concept_pet" {
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

func makeKitchenSinkResource(id string) (tfprotov6.DynamicValue, error) {
	typ := schema.KitchenSinkResourceSchema().ValueType()
	value := map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, id),
	}
	return tfprotov6.NewDynamicValue(typ, tftypes.NewValue(typ, value))
}

func makeKitchenSinkIdentity(id string) *tfprotov6.ResourceIdentityData {
	typ := schema.KitchenSinkIdentitySchema().ValueType()
	value := map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, id),
	}

	identityVal, err := tfprotov6.NewDynamicValue(typ, tftypes.NewValue(typ, value))
	if err != nil {
		panic(err)
	}

	return &tfprotov6.ResourceIdentityData{
		IdentityData: &identityVal,
	}
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
