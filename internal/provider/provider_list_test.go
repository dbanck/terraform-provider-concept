package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestPetList_query(t *testing.T) {
	// Note: querycheck.ExpectLength only considers the summary of the last
	// completed list block, so every step only contains a single list block.
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query: true,
				Config: `
provider "concept" {}

list "concept_pet" "test" {
  provider = concept
}
`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("concept_pet.test", 5),
				},
			},
			{
				Query: true,
				Config: `
provider "concept" {}

list "concept_pet" "test" {
  provider = concept

  config {
    count = 2
  }
}
`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("concept_pet.test", 2),
				},
			},
			{
				// Pet names are random, so there must only be a single
				// result for the filters below to match exactly one result
				Query: true,
				Config: `
provider "concept" {}

list "concept_pet" "test" {
  provider         = concept
  include_resource = true

  config {
    count = 1
  }
}
`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("concept_pet.test", 1),
					querycheck.ExpectResourceDisplayName(
						"concept_pet.test",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"id":   knownvalue.StringRegexp(regexp.MustCompile(`^[a-z]+-[a-z]+$`)),
							"legs": knownvalue.NotNull(),
						}),
						knownvalue.StringRegexp(regexp.MustCompile(`^This is a [a-z]+-[a-z]+$`)),
					),
					querycheck.ExpectResourceKnownValues(
						"concept_pet.test",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"id":   knownvalue.NotNull(),
							"legs": knownvalue.NotNull(),
						}),
						[]querycheck.KnownValueCheck{
							{Path: tfjsonpath.New("length"), KnownValue: knownvalue.Int64Exact(2)},
						},
					),
				},
			},
		},
	})
}

func TestKitchenSinkList_query(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query: true,
				Config: `
provider "concept" {}

list "concept_kitchen_sink" "test" {
  provider         = concept
  include_resource = true
}
`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("concept_kitchen_sink.test", 3),
					querycheck.ExpectIdentity("concept_kitchen_sink.test", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("ks-0"),
					}),
					querycheck.ExpectIdentity("concept_kitchen_sink.test", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("ks-1"),
					}),
					querycheck.ExpectIdentity("concept_kitchen_sink.test", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("ks-2"),
					}),
					querycheck.ExpectResourceDisplayName(
						"concept_kitchen_sink.test",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"id": knownvalue.StringExact("ks-1"),
						}),
						knownvalue.StringExact("Kitchen sink item ks-1"),
					),
					querycheck.ExpectResourceKnownValues(
						"concept_kitchen_sink.test",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"id": knownvalue.StringExact("ks-2"),
						}),
						[]querycheck.KnownValueCheck{
							{Path: tfjsonpath.New("id"), KnownValue: knownvalue.StringExact("ks-2")},
						},
					),
				},
			},
		},
	})
}
