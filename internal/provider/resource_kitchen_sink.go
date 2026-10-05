package provider

import (
	"fmt"
	"math/rand/v2"

	"github.com/dbanck/terraform-provider-concept/internal/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func newKitchenSinkIdentity(id tftypes.Value) (*tfprotov6.ResourceIdentityData, error) {
	identity, err := newObject(schema.KitchenSinkIdentitySchema().ValueType(), map[string]tftypes.Value{
		"id": id,
	})
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ResourceIdentityData{IdentityData: identity}, nil
}

func planKitchenSink(request *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	typ := schema.KitchenSinkResourceSchema().ValueType()

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

	if prior == nil {
		// Create
		planned, err := newObject(typ, map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		})
		if err != nil {
			return nil, err
		}
		return &tfprotov6.PlanResourceChangeResponse{PlannedState: planned}, nil
	}

	// There are no configurable attributes, so there is never anything to update
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

func applyKitchenSink(request *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	typ := schema.KitchenSinkResourceSchema().ValueType()

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
		// Update, nothing to do
		return &tfprotov6.ApplyResourceChangeResponse{
			NewState:    request.PlannedState,
			Private:     request.PlannedPrivate,
			NewIdentity: request.PlannedIdentity,
		}, nil
	}

	// Create
	id := tftypes.NewValue(tftypes.String, fmt.Sprintf("ks-%d", rand.IntN(1_000_000)))
	newState, err := newObject(typ, map[string]tftypes.Value{
		"id": id,
	})
	if err != nil {
		return nil, err
	}
	identity, err := newKitchenSinkIdentity(id)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.ApplyResourceChangeResponse{
		NewState:    newState,
		NewIdentity: identity,
	}, nil
}

func importKitchenSink(request *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	_, id, diags := importID(request, schema.KitchenSinkIdentitySchema().ValueType())
	if diags != nil {
		return &tfprotov6.ImportResourceStateResponse{Diagnostics: diags}, nil
	}

	idValue := tftypes.NewValue(tftypes.String, id)
	state, err := newObject(schema.KitchenSinkResourceSchema().ValueType(), map[string]tftypes.Value{
		"id": idValue,
	})
	if err != nil {
		return nil, err
	}
	identity, err := newKitchenSinkIdentity(idValue)
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
