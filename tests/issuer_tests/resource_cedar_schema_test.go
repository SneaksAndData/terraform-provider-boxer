package issuer_tests

import (
	helpers "terraform-provider-boxer/tests"
	"terraform-provider-boxer/tests/assertions"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestResourceCedarSchema_creation(t *testing.T) {
	const resourceAddress = "boxer_issuer_cedar_schema.example"
	const templateName = "resource_cedar_schema/resource_cedar_schema.tmpl.tf"

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	services := helpers.NewLocalServices()
	token, err := helpers.GetExternalToken(services)
	if err != nil {
		t.Fatalf("failed to get external token: %s", err)
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: helpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: helpers.RenderTemplate(helpers.NewTestContext(randomName, services, token), templateName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					assertions.ValidateSchemaIsParseable(resourceAddress),
				},
			},
		},
	})
}

func TestResourceCedarSchema_creation_validated(t *testing.T) {
	const resourceAddress = "boxer_issuer_cedar_schema.example"
	const templateName = "resource_cedar_schema_validated/resource_cedar_schema_validated.tmpl.tf"

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	services := helpers.NewLocalServices()
	token, err := helpers.GetExternalToken(services)
	if err != nil {
		t.Fatalf("failed to get external token: %s", err)
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: helpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: helpers.RenderTemplate(helpers.NewTestContext(randomName, services, token), templateName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					assertions.ValidateSchemaIsParseable(resourceAddress),
				},
			},
		},
	})
}
