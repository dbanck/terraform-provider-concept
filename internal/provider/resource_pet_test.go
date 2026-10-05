package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// petNameRegexp matches a pet name consisting of exactly n words.
func petNameRegexp(n int) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`^[a-z]+(-[a-z]+){%d}$`, n-1))
}

func TestPet_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "concept_pet" "test" {}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("concept_pet.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("concept_pet.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("length"), knownvalue.Int64Exact(2)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("id"), knownvalue.StringRegexp(petNameRegexp(2))),
					statecheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("length"), knownvalue.Int64Exact(2)),
					statecheck.ExpectIdentityValueMatchesState("concept_pet.test", tfjsonpath.New("id")),
					statecheck.ExpectIdentityValue("concept_pet.test", tfjsonpath.New("legs"), knownvalue.Int64Func(func(v int64) error {
						if v < 2 || v > 7 {
							return fmt.Errorf("expected legs between 2 and 7, got %d", v)
						}
						return nil
					})),
				},
			},
			{
				// Re-applying the same configuration must neither change the
				// pet nor its identity.
				Config: `resource "concept_pet" "test" {}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestPet_length(t *testing.T) {
	sameID := statecheck.CompareValue(compare.ValuesSame())
	differentID := statecheck.CompareValue(compare.ValuesDiffer())

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "concept_pet" "test" { length = 3 }`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("length"), knownvalue.Int64Exact(3)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("id"), knownvalue.StringRegexp(petNameRegexp(3))),
					sameID.AddStateValue("concept_pet.test", tfjsonpath.New("id")),
					differentID.AddStateValue("concept_pet.test", tfjsonpath.New("id")),
				},
			},
			{
				// Removing the optional+computed length keeps the prior value
				Config: `resource "concept_pet" "test" {}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("length"), knownvalue.Int64Exact(3)),
					sameID.AddStateValue("concept_pet.test", tfjsonpath.New("id")),
				},
			},
			{
				// Changing the length requires a new pet
				Config: `resource "concept_pet" "test" { length = 1 }`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("concept_pet.test", plancheck.ResourceActionDestroyBeforeCreate),
						plancheck.ExpectUnknownValue("concept_pet.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("length"), knownvalue.Int64Exact(1)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("id"), knownvalue.StringRegexp(petNameRegexp(1))),
					statecheck.ExpectIdentityValueMatchesState("concept_pet.test", tfjsonpath.New("id")),
					differentID.AddStateValue("concept_pet.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func TestPet_invalidLength(t *testing.T) {
	for name, length := range map[string]string{
		"zero":       "0",
		"negative":   "-1",
		"fractional": "1.5",
	} {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      `resource "concept_pet" "test" { length = ` + length + ` }`,
						ExpectError: regexp.MustCompile(`Invalid pet length`),
					},
				},
			})
		})
	}
}

func TestPet_import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "concept_pet" "test" { length = 3 }`,
			},
			{
				ResourceName:      "concept_pet.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportCommandWithID,
				ImportStateVerify: true,
			},
			{
				ResourceName:    "concept_pet.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
			},
			{
				// Also verifies that the identity (including legs) is
				// preserved during import
				ResourceName:    "concept_pet.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

// TestPet_importFromSearch simulates the Terraform Search workflow: the
// configuration below is what `terraform query -generate-config-out` produces
// for a concept_pet list result. Planning it must result in a no-op import.
func TestPet_importFromSearch(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "concept_pet" "pets_0" {
  provider = concept
  length   = 2
}

import {
  to       = concept_pet.pets_0
  provider = concept
  identity = {
    id   = "probable-salmon"
    legs = 4
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("concept_pet.pets_0", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("concept_pet.pets_0", tfjsonpath.New("id"), knownvalue.StringExact("probable-salmon")),
						plancheck.ExpectKnownValue("concept_pet.pets_0", tfjsonpath.New("length"), knownvalue.Int64Exact(2)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity("concept_pet.pets_0", map[string]knownvalue.Check{
						"id":   knownvalue.StringExact("probable-salmon"),
						"legs": knownvalue.Int64Exact(4),
					}),
				},
			},
		},
	})
}

// TestPet_importIdentityWithoutLegs verifies that the optional identity
// attribute legs can be omitted on import.
func TestPet_importIdentityWithoutLegs(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "concept_pet" "test" {}

import {
  to       = concept_pet.test
  identity = {
    id = "cat"
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// The imported length (1) differs from the default (2),
						// but as length isn't configured, the imported value is
						// kept.
						plancheck.ExpectResourceAction("concept_pet.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("concept_pet.test", tfjsonpath.New("length"), knownvalue.Int64Exact(1)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue("concept_pet.test", tfjsonpath.New("id"), knownvalue.StringExact("cat")),
					statecheck.ExpectIdentityValue("concept_pet.test", tfjsonpath.New("legs"), knownvalue.NotNull()),
				},
			},
		},
	})
}
