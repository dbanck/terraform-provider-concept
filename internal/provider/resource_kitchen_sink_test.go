package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestKitchenSink_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "concept_kitchen_sink" "test" {}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("concept_kitchen_sink.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("concept_kitchen_sink.test", tfjsonpath.New("id")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("concept_kitchen_sink.test", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`^ks-\d+$`))),
					statecheck.ExpectIdentityValueMatchesState("concept_kitchen_sink.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: `resource "concept_kitchen_sink" "test" {}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestKitchenSink_import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "concept_kitchen_sink" "test" {}`,
			},
			{
				ResourceName:      "concept_kitchen_sink.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportCommandWithID,
				ImportStateVerify: true,
			},
			{
				ResourceName:    "concept_kitchen_sink.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
			},
			{
				ResourceName:    "concept_kitchen_sink.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

// TestKitchenSink_importFromSearch simulates the Terraform Search workflow: the
// configuration below is what `terraform query -generate-config-out` produces
// for a concept_kitchen_sink list result. Planning it must result in a no-op
// import.
func TestKitchenSink_importFromSearch(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "concept_kitchen_sink" "ks_0" {
  provider = concept
}

import {
  to       = concept_kitchen_sink.ks_0
  provider = concept
  identity = {
    id = "ks-0"
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("concept_kitchen_sink.ks_0", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("concept_kitchen_sink.ks_0", tfjsonpath.New("id"), knownvalue.StringExact("ks-0")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity("concept_kitchen_sink.ks_0", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("ks-0"),
					}),
				},
			},
		},
	})
}
