package issuer_tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"terraform-provider-boxer/internal/tests"
	"testing"
)

func TestDataSourceBoxerExternalIdentity_reading(t *testing.T) {
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	services := tests.NewLocalServices()
	token, err := tests.GetExternalToken(services)
	if err != nil {
		t.Fatalf("failed to get external token: %s", err)
	}
	testContext := tests.NewTestContext(randomName, services, token)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tests.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tests.RenderTemplate(testContext, "data_source_boxer_external_identity.tmpl.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.boxer_external_identity.example",
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					statecheck.ExpectKnownValue(
						"data.boxer_external_identity.example",
						tfjsonpath.New("identity_provider"),
						knownvalue.StringExact(randomName),
					),
					statecheck.ExpectKnownValue(
						"data.boxer_external_identity.example",
						tfjsonpath.New("principal").AtMapKey("principal_id"),
						knownvalue.StringExact("PhotoApp::User::\"alice\""),
					),
					statecheck.ExpectKnownValue(
						"data.boxer_external_identity.example",
						tfjsonpath.New("principal").AtMapKey("schema_id"),
						knownvalue.StringExact(fmt.Sprintf("%s-issuer", randomName)),
					),
					statecheck.ExpectKnownValue(
						"data.boxer_external_identity.example",
						tfjsonpath.New("validator_schema_id"),
						knownvalue.StringExact(fmt.Sprintf("%s-validator", randomName)),
					),
				},
			},
		},
	})
}
