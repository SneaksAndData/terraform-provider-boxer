package issuer_tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	tests2 "terraform-provider-boxer/tests"
	"testing"
)

func TestDataSourceBoxerExternalIdentity_reading(t *testing.T) {
	const resourceAddress = "data.boxer_external_identity.example"
	const templateName = "data_source_boxer_external_identity/data_source_boxer_external_identity.tmpl.tf"

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	services := tests2.NewLocalServices()
	token, err := tests2.GetExternalToken(services)
	if err != nil {
		t.Fatalf("failed to get external token: %s", err)
	}
	testContext := tests2.NewTestContext(randomName, services, token)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tests2.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tests2.RenderTemplate(testContext, templateName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("identity_provider"),
						knownvalue.StringExact(randomName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("principal").AtMapKey("principal_id"),
						knownvalue.StringExact("PhotoApp::User::\"alice\""),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("principal").AtMapKey("schema_id"),
						knownvalue.StringExact(fmt.Sprintf("%s-issuer", randomName)),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("validator_schema_id"),
						knownvalue.StringExact(fmt.Sprintf("%s-validator", randomName)),
					),
				},
			},
		},
	})
}
