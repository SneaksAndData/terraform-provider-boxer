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

func TestResourceIdentityProvider_creation(t *testing.T) {
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	services := tests2.NewLocalServices()
	token, err := tests2.GetExternalToken(services)
	if err != nil {
		t.Fatalf("failed to get external token: %s", err)
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tests2.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tests2.RenderTemplate(tests2.NewTestContext(randomName, services, token), "resource_identity_provider.tmpl"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("name"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("discovery_url"),
						knownvalue.StringExact(services.ExternalIdp.ClusterEndpoint),
					),
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("user_id_claim"),
						knownvalue.StringExact("preferred_username"),
					),
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("issuers"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact(services.ExternalIdp.Endpoint)}),
					),
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("audiences"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("account")}),
					),
				},
			},
		},
	})
}
