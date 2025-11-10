package issuer_tests

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"terraform-provider-boxer/internal/tests"
	"testing"
)

func TestResourceActionDiscoveryDocument_creation(t *testing.T) {
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	services := tests.NewLocalServices()
	token, err := tests.GetExternalToken(services)
	testContext := tests.NewTestContext(randomName, services, token)
	if err != nil {
		t.Fatalf("failed to get external token: %s", err)
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tests.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tests.RenderTemplate(testContext, "resource_action_discovery_document.tmpl.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"boxer_action_discovery_document.example",
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					statecheck.ExpectKnownValue(
						"boxer_action_discovery_document.example",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("www.example.com"),
					),
					statecheck.ExpectKnownValue(
						"boxer_action_discovery_document.example",
						tfjsonpath.New("routes"),
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"method": knownvalue.StringExact("GET"),
									"path":   knownvalue.StringExact("api/v1/resource"),
									"action": knownvalue.StringExact("PhotoApp::Action::\"viewPhoto\""),
								}),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"method": knownvalue.StringExact("GET"),
									"path":   knownvalue.StringExact("api/v2/resource"),
									"action": knownvalue.StringExact("PhotoApp::Action::\"viewPhoto\""),
								}),
							},
						),
					),
				},
			},
		},
	})
}
