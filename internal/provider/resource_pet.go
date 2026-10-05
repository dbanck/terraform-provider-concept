package provider

import (
	"fmt"
	"math/big"
	"math/rand/v2"
	"strings"

	"github.com/dbanck/terraform-provider-concept/internal/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	petname "github.com/dustinkirkland/golang-petname"
)

const defaultPetLength = 2

var petLengthPath = tftypes.NewAttributePath().WithAttributeName("length")

func randomLegs() tftypes.Value {
	return tftypes.NewValue(tftypes.Number, rand.IntN(6)+2)
}

func newPetIdentity(id, legs tftypes.Value) (*tfprotov6.ResourceIdentityData, error) {
	identity, err := newObject(schema.PetIdentitySchema().ValueType(), map[string]tftypes.Value{
		"id":   id,
		"legs": legs,
	})
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ResourceIdentityData{IdentityData: identity}, nil
}

// petLength validates and returns the given length value as integer.
func petLength(value tftypes.Value) (int, error) {
	var length big.Float
	if err := value.As(&length); err != nil {
		return 0, err
	}
	if !length.IsInt() || length.Cmp(big.NewFloat(1)) < 0 {
		return 0, fmt.Errorf("length must be a whole number greater than or equal to 1, got %s", length.Text('f', -1))
	}
	l, _ := length.Int64()
	return int(l), nil
}

func validatePet(request *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	config, err := unmarshalObject(request.Config, schema.PetResourceSchema().ValueType())
	if err != nil {
		return nil, err
	}

	if config != nil && config["length"].IsKnown() && !config["length"].IsNull() {
		if _, err := petLength(config["length"]); err != nil {
			return &tfprotov6.ValidateResourceConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity:  tfprotov6.DiagnosticSeverityError,
						Summary:   "Invalid pet length",
						Detail:    err.Error(),
						Attribute: petLengthPath,
					},
				},
			}, nil
		}
	}

	return &tfprotov6.ValidateResourceConfigResponse{}, nil
}

func planPet(request *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	typ := schema.PetResourceSchema().ValueType()

	proposed, err := unmarshalObject(request.ProposedNewState, typ)
	if err != nil {
		return nil, err
	}
	if proposed == nil {
		// Destroy
		return &tfprotov6.PlanResourceChangeResponse{PlannedState: request.ProposedNewState}, nil
	}

	prior, err := unmarshalObject(request.PriorState, typ)
	if err != nil {
		return nil, err
	}
	config, err := unmarshalObject(request.Config, typ)
	if err != nil {
		return nil, err
	}

	// length is optional and computed: if it's not configured, we keep the
	// prior value or fall back to the default for new resources.
	length := config["length"]
	if length.IsNull() {
		if prior != nil {
			length = prior["length"]
		} else {
			length = tftypes.NewValue(tftypes.Number, defaultPetLength)
		}
	}

	if prior == nil {
		// Create
		planned, err := newObject(typ, map[string]tftypes.Value{
			"id":     tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"length": length,
		})
		if err != nil {
			return nil, err
		}
		return &tfprotov6.PlanResourceChangeResponse{PlannedState: planned}, nil
	}

	if !length.Equal(prior["length"]) {
		// A different length requires a new pet name
		planned, err := newObject(typ, map[string]tftypes.Value{
			"id":     tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"length": length,
		})
		if err != nil {
			return nil, err
		}
		return &tfprotov6.PlanResourceChangeResponse{
			PlannedState:    planned,
			RequiresReplace: []*tftypes.AttributePath{petLengthPath},
		}, nil
	}

	// No changes
	planned, err := newObject(typ, prior)
	if err != nil {
		return nil, err
	}
	return &tfprotov6.PlanResourceChangeResponse{
		PlannedState:    planned,
		PlannedPrivate:  request.PriorPrivate,
		PlannedIdentity: request.PriorIdentity,
	}, nil
}

func applyPet(request *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	typ := schema.PetResourceSchema().ValueType()

	planned, err := unmarshalObject(request.PlannedState, typ)
	if err != nil {
		return nil, err
	}
	if planned == nil {
		// Destroy
		return &tfprotov6.ApplyResourceChangeResponse{NewState: request.PlannedState}, nil
	}

	prior, err := unmarshalObject(request.PriorState, typ)
	if err != nil {
		return nil, err
	}

	if prior != nil {
		// Update, nothing to do as all changes require replacement
		return &tfprotov6.ApplyResourceChangeResponse{
			NewState:    request.PlannedState,
			Private:     request.PlannedPrivate,
			NewIdentity: request.PlannedIdentity,
		}, nil
	}

	// Create
	length, err := petLength(planned["length"])
	if err != nil {
		return &tfprotov6.ApplyResourceChangeResponse{Diagnostics: errorDiag("Invalid pet length", err.Error())}, nil
	}

	petname.NonDeterministicMode()
	id := tftypes.NewValue(tftypes.String, strings.ToLower(petname.Generate(length, "-")))

	newState, err := newObject(typ, map[string]tftypes.Value{
		"id":     id,
		"length": planned["length"],
	})
	if err != nil {
		return nil, err
	}
	identity, err := newPetIdentity(id, randomLegs())
	if err != nil {
		return nil, err
	}

	return &tfprotov6.ApplyResourceChangeResponse{
		NewState:    newState,
		NewIdentity: identity,
	}, nil
}

func importPet(request *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	identityAttrs, id, diags := importID(request, schema.PetIdentitySchema().ValueType())
	if diags != nil {
		return &tfprotov6.ImportResourceStateResponse{Diagnostics: diags}, nil
	}

	legs := randomLegs()
	if identityAttrs != nil && identityAttrs["legs"].IsKnown() && !identityAttrs["legs"].IsNull() {
		legs = identityAttrs["legs"]
	}

	idValue := tftypes.NewValue(tftypes.String, id)
	state, err := newObject(schema.PetResourceSchema().ValueType(), map[string]tftypes.Value{
		"id": idValue,
		// Pet names consist of `length` words joined by dashes
		"length": tftypes.NewValue(tftypes.Number, len(strings.Split(id, "-"))),
	})
	if err != nil {
		return nil, err
	}
	identity, err := newPetIdentity(idValue, legs)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.ImportResourceStateResponse{
		ImportedResources: []*tfprotov6.ImportedResource{
			{
				TypeName: request.TypeName,
				State:    state,
				Identity: identity,
			},
		},
	}, nil
}
