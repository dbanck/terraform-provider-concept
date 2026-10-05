package provider

import (
	"context"
	"fmt"

	"github.com/dbanck/terraform-provider-concept/internal/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// resourceStateType returns the state type of the given managed resource type.
func resourceStateType(typeName string) (tftypes.Type, bool) {
	switch typeName {
	case "concept_pet":
		return schema.PetResourceSchema().ValueType(), true
	case "concept_kitchen_sink":
		return schema.KitchenSinkResourceSchema().ValueType(), true
	}
	return nil, false
}

// resourceIdentityType returns the identity type of the given managed resource type.
func resourceIdentityType(typeName string) (tftypes.Type, bool) {
	switch typeName {
	case "concept_pet":
		return schema.PetIdentitySchema().ValueType(), true
	case "concept_kitchen_sink":
		return schema.KitchenSinkIdentitySchema().ValueType(), true
	}
	return nil, false
}

// ignoreUndefinedAttributes allows reading state that was written with a
// schema that contained attributes that have been removed since.
var ignoreUndefinedAttributes = tfprotov6.UnmarshalOpts{
	ValueFromJSONOpts: tftypes.ValueFromJSONOpts{
		IgnoreUndefinedAttributes: true,
	},
}

func (p ConceptProvider) UpgradeResourceIdentity(ctx context.Context, request *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	typ, ok := resourceIdentityType(request.TypeName)
	if !ok {
		return &tfprotov6.UpgradeResourceIdentityResponse{Diagnostics: unknownResourceType(request.TypeName)}, nil
	}

	value := tftypes.NewValue(typ, nil)
	if request.RawIdentity != nil {
		var err error
		value, err = request.RawIdentity.UnmarshalWithOpts(typ, ignoreUndefinedAttributes)
		if err != nil {
			return &tfprotov6.UpgradeResourceIdentityResponse{
				Diagnostics: errorDiag("Unable to upgrade resource identity", err.Error()),
			}, nil
		}
	}

	identity, err := tfprotov6.NewDynamicValue(typ, value)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.UpgradeResourceIdentityResponse{
		UpgradedIdentity: &tfprotov6.ResourceIdentityData{IdentityData: &identity},
	}, nil
}

func (p ConceptProvider) UpgradeResourceState(ctx context.Context, request *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	typ, ok := resourceStateType(request.TypeName)
	if !ok {
		return &tfprotov6.UpgradeResourceStateResponse{Diagnostics: unknownResourceType(request.TypeName)}, nil
	}

	value := tftypes.NewValue(typ, nil)
	if request.RawState != nil {
		var err error
		value, err = request.RawState.UnmarshalWithOpts(typ, ignoreUndefinedAttributes)
		if err != nil {
			return &tfprotov6.UpgradeResourceStateResponse{
				Diagnostics: errorDiag("Unable to upgrade resource state", err.Error()),
			}, nil
		}
	}

	state, err := tfprotov6.NewDynamicValue(typ, value)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.UpgradeResourceStateResponse{UpgradedState: &state}, nil
}

// ReadResource returns the current state unchanged, as the resources of this
// provider don't exist anywhere outside of the Terraform state. If the state
// doesn't contain an identity yet (e.g. it was created by an older version of
// this provider), it is derived from the state.
func (p ConceptProvider) ReadResource(ctx context.Context, request *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	typ, ok := resourceStateType(request.TypeName)
	if !ok {
		return &tfprotov6.ReadResourceResponse{Diagnostics: unknownResourceType(request.TypeName)}, nil
	}

	identity := request.CurrentIdentity
	if identity == nil || identity.IdentityData == nil {
		state, err := unmarshalObject(request.CurrentState, typ)
		if err != nil {
			return nil, err
		}

		if state != nil {
			switch request.TypeName {
			case "concept_pet":
				identity, err = newPetIdentity(state["id"], randomLegs())
			case "concept_kitchen_sink":
				identity, err = newKitchenSinkIdentity(state["id"])
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return &tfprotov6.ReadResourceResponse{
		NewState:    request.CurrentState,
		Private:     request.Private,
		NewIdentity: identity,
	}, nil
}

func (p ConceptProvider) PlanResourceChange(ctx context.Context, request *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	switch request.TypeName {
	case "concept_pet":
		return planPet(request)
	case "concept_kitchen_sink":
		return planKitchenSink(request)
	}
	return &tfprotov6.PlanResourceChangeResponse{Diagnostics: unknownResourceType(request.TypeName)}, nil
}

func (p ConceptProvider) ApplyResourceChange(ctx context.Context, request *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	switch request.TypeName {
	case "concept_pet":
		return applyPet(request)
	case "concept_kitchen_sink":
		return applyKitchenSink(request)
	}
	return &tfprotov6.ApplyResourceChangeResponse{Diagnostics: unknownResourceType(request.TypeName)}, nil
}

func (p ConceptProvider) MoveResourceState(ctx context.Context, request *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	return &tfprotov6.MoveResourceStateResponse{
		Diagnostics: errorDiag(
			"Unsupported resource move",
			fmt.Sprintf("Moving resource state from %q to %q is not supported by this provider.", request.SourceTypeName, request.TargetTypeName),
		),
	}, nil
}

func (p ConceptProvider) ImportResourceState(ctx context.Context, request *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	switch request.TypeName {
	case "concept_pet":
		return importPet(request)
	case "concept_kitchen_sink":
		return importKitchenSink(request)
	}
	return &tfprotov6.ImportResourceStateResponse{Diagnostics: unknownResourceType(request.TypeName)}, nil
}

// importID returns the ID of the resource to import, either from the given
// identity or from the legacy import ID.
func importID(request *tfprotov6.ImportResourceStateRequest, identityType tftypes.Type) (map[string]tftypes.Value, string, []*tfprotov6.Diagnostic) {
	if request.Identity != nil && request.Identity.IdentityData != nil {
		identity, err := unmarshalObject(request.Identity.IdentityData, identityType)
		if err != nil {
			return nil, "", errorDiag("Invalid import identity", err.Error())
		}
		if identity == nil || !identity["id"].IsFullyKnown() || identity["id"].IsNull() {
			return nil, "", errorDiag("Invalid import identity", "The identity attribute \"id\" must be set to a known value.")
		}

		var id string
		if err := identity["id"].As(&id); err != nil {
			return nil, "", errorDiag("Invalid import identity", err.Error())
		}
		if id == "" {
			return nil, "", errorDiag("Invalid import identity", "The identity attribute \"id\" must not be empty.")
		}
		return identity, id, nil
	}

	if request.ID == "" {
		return nil, "", errorDiag("Invalid import ID", "The import ID must not be empty.")
	}
	return nil, request.ID, nil
}

func unknownResourceType(typeName string) []*tfprotov6.Diagnostic {
	return errorDiag("Unknown resource type", fmt.Sprintf("The resource type %q is not supported by this provider.", typeName))
}

func errorDiag(summary, detail string) []*tfprotov6.Diagnostic {
	return []*tfprotov6.Diagnostic{
		{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  summary,
			Detail:   detail,
		},
	}
}

// unmarshalObject decodes an object DynamicValue into its attributes. It
// returns a nil map if the value is null or missing.
func unmarshalObject(dv *tfprotov6.DynamicValue, typ tftypes.Type) (map[string]tftypes.Value, error) {
	if dv == nil {
		return nil, nil
	}

	value, err := dv.Unmarshal(typ)
	if err != nil {
		return nil, err
	}
	if value.IsNull() {
		return nil, nil
	}

	attrs := map[string]tftypes.Value{}
	if err := value.As(&attrs); err != nil {
		return nil, err
	}
	return attrs, nil
}

// newObject encodes the given attributes as DynamicValue. A nil map results in
// a null value.
func newObject(typ tftypes.Type, attrs map[string]tftypes.Value) (*tfprotov6.DynamicValue, error) {
	var value tftypes.Value
	if attrs == nil {
		value = tftypes.NewValue(typ, nil)
	} else {
		value = tftypes.NewValue(typ, attrs)
	}

	dv, err := tfprotov6.NewDynamicValue(typ, value)
	if err != nil {
		return nil, err
	}
	return &dv, nil
}
