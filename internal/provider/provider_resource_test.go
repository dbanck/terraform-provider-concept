package provider

import (
	"context"
	"testing"

	"github.com/dbanck/terraform-provider-concept/internal/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	petType                 = schema.PetResourceSchema().ValueType()
	petIdentityType         = schema.PetIdentitySchema().ValueType()
	kitchenSinkType         = schema.KitchenSinkResourceSchema().ValueType()
	kitchenSinkIdentityType = schema.KitchenSinkIdentitySchema().ValueType()
)

func mustObject(t *testing.T, typ tftypes.Type, attrs map[string]tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	dv, err := newObject(typ, attrs)
	if err != nil {
		t.Fatal(err)
	}
	return dv
}

func mustIdentity(t *testing.T, typ tftypes.Type, attrs map[string]tftypes.Value) *tfprotov6.ResourceIdentityData {
	t.Helper()
	return &tfprotov6.ResourceIdentityData{IdentityData: mustObject(t, typ, attrs)}
}

func mustAttrs(t *testing.T, dv *tfprotov6.DynamicValue, typ tftypes.Type) map[string]tftypes.Value {
	t.Helper()
	attrs, err := unmarshalObject(dv, typ)
	if err != nil {
		t.Fatal(err)
	}
	return attrs
}

func expectNoDiags(t *testing.T, diags []*tfprotov6.Diagnostic) {
	t.Helper()
	for _, d := range diags {
		t.Errorf("unexpected diagnostic: %s: %s", d.Summary, d.Detail)
	}
}

func expectErrorDiag(t *testing.T, diags []*tfprotov6.Diagnostic, summary string) {
	t.Helper()
	if len(diags) != 1 || diags[0].Severity != tfprotov6.DiagnosticSeverityError || diags[0].Summary != summary {
		t.Fatalf("expected a single %q error diagnostic, got %+v", summary, diags)
	}
}

func str(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }
func num(n int) tftypes.Value    { return tftypes.NewValue(tftypes.Number, n) }

func TestUpgradeResourceState_ignoresRemovedAttributes(t *testing.T) {
	resp, err := ConceptProvider{}.UpgradeResourceState(context.Background(), &tfprotov6.UpgradeResourceStateRequest{
		TypeName: "concept_pet",
		RawState: &tfprotov6.RawState{JSON: []byte(`{"id":"happy-cat","length":2,"removed":"attr"}`)},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectNoDiags(t, resp.Diagnostics)

	state := mustAttrs(t, resp.UpgradedState, petType)
	if !state["id"].Equal(str("happy-cat")) || !state["length"].Equal(num(2)) {
		t.Errorf("unexpected upgraded state: %v", state)
	}
}

func TestUpgradeResourceIdentity(t *testing.T) {
	resp, err := ConceptProvider{}.UpgradeResourceIdentity(context.Background(), &tfprotov6.UpgradeResourceIdentityRequest{
		TypeName:    "concept_kitchen_sink",
		RawIdentity: &tfprotov6.RawState{JSON: []byte(`{"id":"ks-1"}`)},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectNoDiags(t, resp.Diagnostics)

	identity := mustAttrs(t, resp.UpgradedIdentity.IdentityData, kitchenSinkIdentityType)
	if !identity["id"].Equal(str("ks-1")) {
		t.Errorf("unexpected upgraded identity: %v", identity)
	}
}

func TestReadResource_keepsIdentity(t *testing.T) {
	identity := mustIdentity(t, petIdentityType, map[string]tftypes.Value{"id": str("happy-cat"), "legs": num(4)})

	resp, err := ConceptProvider{}.ReadResource(context.Background(), &tfprotov6.ReadResourceRequest{
		TypeName:        "concept_pet",
		CurrentState:    mustObject(t, petType, map[string]tftypes.Value{"id": str("happy-cat"), "length": num(2)}),
		CurrentIdentity: identity,
	})
	if err != nil {
		t.Fatal(err)
	}
	expectNoDiags(t, resp.Diagnostics)

	if resp.NewIdentity != identity {
		t.Errorf("expected the identity to be returned unchanged")
	}
}

func TestReadResource_backfillsMissingIdentity(t *testing.T) {
	t.Run("concept_pet", func(t *testing.T) {
		resp, err := ConceptProvider{}.ReadResource(context.Background(), &tfprotov6.ReadResourceRequest{
			TypeName:     "concept_pet",
			CurrentState: mustObject(t, petType, map[string]tftypes.Value{"id": str("happy-cat"), "length": num(2)}),
		})
		if err != nil {
			t.Fatal(err)
		}
		expectNoDiags(t, resp.Diagnostics)

		identity := mustAttrs(t, resp.NewIdentity.IdentityData, petIdentityType)
		if !identity["id"].Equal(str("happy-cat")) || identity["legs"].IsNull() {
			t.Errorf("unexpected identity: %v", identity)
		}
	})

	t.Run("concept_kitchen_sink", func(t *testing.T) {
		resp, err := ConceptProvider{}.ReadResource(context.Background(), &tfprotov6.ReadResourceRequest{
			TypeName:     "concept_kitchen_sink",
			CurrentState: mustObject(t, kitchenSinkType, map[string]tftypes.Value{"id": str("ks-7")}),
		})
		if err != nil {
			t.Fatal(err)
		}
		expectNoDiags(t, resp.Diagnostics)

		identity := mustAttrs(t, resp.NewIdentity.IdentityData, kitchenSinkIdentityType)
		if !identity["id"].Equal(str("ks-7")) {
			t.Errorf("unexpected identity: %v", identity)
		}
	})
}

func TestPlanResourceChange_destroy(t *testing.T) {
	for typeName, typ := range map[string]tftypes.Type{
		"concept_pet":          petType,
		"concept_kitchen_sink": kitchenSinkType,
	} {
		t.Run(typeName, func(t *testing.T) {
			var prior map[string]tftypes.Value
			if typeName == "concept_pet" {
				prior = map[string]tftypes.Value{"id": str("happy-cat"), "length": num(2)}
			} else {
				prior = map[string]tftypes.Value{"id": str("ks-1")}
			}

			resp, err := ConceptProvider{}.PlanResourceChange(context.Background(), &tfprotov6.PlanResourceChangeRequest{
				TypeName:         typeName,
				PriorState:       mustObject(t, typ, prior),
				ProposedNewState: mustObject(t, typ, nil),
				Config:           mustObject(t, typ, nil),
			})
			if err != nil {
				t.Fatal(err)
			}
			expectNoDiags(t, resp.Diagnostics)

			if planned := mustAttrs(t, resp.PlannedState, typ); planned != nil {
				t.Errorf("expected null planned state, got %v", planned)
			}
		})
	}
}

func TestPlanResourceChange_petUnknownLength(t *testing.T) {
	unknownLength := tftypes.NewValue(tftypes.Number, tftypes.UnknownValue)
	config := mustObject(t, petType, map[string]tftypes.Value{
		"id":     tftypes.NewValue(tftypes.String, nil),
		"length": unknownLength,
	})

	// Unknown values must pass validation
	validateResp, err := ConceptProvider{}.ValidateResourceConfig(context.Background(), &tfprotov6.ValidateResourceConfigRequest{
		TypeName: "concept_pet",
		Config:   config,
	})
	if err != nil {
		t.Fatal(err)
	}
	expectNoDiags(t, validateResp.Diagnostics)

	// An unknown length might differ from the prior one, so the pet has to
	// be replaced
	resp, err := ConceptProvider{}.PlanResourceChange(context.Background(), &tfprotov6.PlanResourceChangeRequest{
		TypeName:         "concept_pet",
		PriorState:       mustObject(t, petType, map[string]tftypes.Value{"id": str("happy-cat"), "length": num(2)}),
		ProposedNewState: mustObject(t, petType, map[string]tftypes.Value{"id": str("happy-cat"), "length": unknownLength}),
		Config:           config,
	})
	if err != nil {
		t.Fatal(err)
	}
	expectNoDiags(t, resp.Diagnostics)

	planned := mustAttrs(t, resp.PlannedState, petType)
	if planned["id"].IsKnown() || planned["length"].IsKnown() {
		t.Errorf("expected id and length to be unknown, got %v", planned)
	}
	if len(resp.RequiresReplace) != 1 || !resp.RequiresReplace[0].Equal(petLengthPath) {
		t.Errorf("expected length to require replacement, got %v", resp.RequiresReplace)
	}
}

func TestImportResourceState_invalid(t *testing.T) {
	t.Run("empty ID", func(t *testing.T) {
		resp, err := ConceptProvider{}.ImportResourceState(context.Background(), &tfprotov6.ImportResourceStateRequest{
			TypeName: "concept_pet",
			ID:       "",
		})
		if err != nil {
			t.Fatal(err)
		}
		expectErrorDiag(t, resp.Diagnostics, "Invalid import ID")
	})

	t.Run("null identity id", func(t *testing.T) {
		resp, err := ConceptProvider{}.ImportResourceState(context.Background(), &tfprotov6.ImportResourceStateRequest{
			TypeName: "concept_kitchen_sink",
			Identity: mustIdentity(t, kitchenSinkIdentityType, map[string]tftypes.Value{
				"id": tftypes.NewValue(tftypes.String, nil),
			}),
		})
		if err != nil {
			t.Fatal(err)
		}
		expectErrorDiag(t, resp.Diagnostics, "Invalid import identity")
	})
}

func TestMoveResourceState_unsupported(t *testing.T) {
	resp, err := ConceptProvider{}.MoveResourceState(context.Background(), &tfprotov6.MoveResourceStateRequest{
		SourceTypeName: "random_pet",
		TargetTypeName: "concept_pet",
	})
	if err != nil {
		t.Fatal(err)
	}
	expectErrorDiag(t, resp.Diagnostics, "Unsupported resource move")
}

func TestUnknownResourceType(t *testing.T) {
	ctx := context.Background()
	p := ConceptProvider{}
	const typeName = "concept_unknown"

	upgradeState, _ := p.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{TypeName: typeName})
	expectErrorDiag(t, upgradeState.Diagnostics, "Unknown resource type")

	upgradeIdentity, _ := p.UpgradeResourceIdentity(ctx, &tfprotov6.UpgradeResourceIdentityRequest{TypeName: typeName})
	expectErrorDiag(t, upgradeIdentity.Diagnostics, "Unknown resource type")

	read, _ := p.ReadResource(ctx, &tfprotov6.ReadResourceRequest{TypeName: typeName})
	expectErrorDiag(t, read.Diagnostics, "Unknown resource type")

	plan, _ := p.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{TypeName: typeName})
	expectErrorDiag(t, plan.Diagnostics, "Unknown resource type")

	apply, _ := p.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{TypeName: typeName})
	expectErrorDiag(t, apply.Diagnostics, "Unknown resource type")

	imp, _ := p.ImportResourceState(ctx, &tfprotov6.ImportResourceStateRequest{TypeName: typeName, ID: "x"})
	expectErrorDiag(t, imp.Diagnostics, "Unknown resource type")
}
