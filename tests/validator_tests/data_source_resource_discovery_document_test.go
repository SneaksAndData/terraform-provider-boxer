package issuer_tests

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	tests2 "terraform-provider-boxer/tests"
	"testing"
)

func TestDataSourceResourceDiscoveryDocument_reading(t *testing.T) {
	t.Skip("Skipping because of the bug #56")

	const resourceAddress = "data.boxer_resource_discovery_document.example"
	const templateName = "data_source_resource_discovery_document/data_source_resource_discovery_document.tmpl.tf"

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	services := tests2.NewLocalServices()
	token, err := tests2.GetExternalToken(services)
	testContext := tests2.NewTestContext(randomName, services, token)
	if err != nil {
		t.Fatalf("failed to get external token: %s", err)
	}
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
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("www.example.com"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("routes"),
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"path":     knownvalue.StringExact("api/v1/resource"),
									"resource": knownvalue.StringExact("PhotoApp::Photo::\"vacationPhoto.jpg\""),
								}),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"path":     knownvalue.StringExact("api/v2/resource"),
									"resource": knownvalue.StringExact("PhotoApp::Photo::\"vacationPhoto.jpg\""),
								}),
							},
						),
					),
				},
			},
		},
	})
}
