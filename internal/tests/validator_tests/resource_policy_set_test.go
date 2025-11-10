package issuer_tests

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"terraform-provider-boxer/internal/tests"
	"terraform-provider-boxer/internal/tests/assertions"
	"testing"
)

func TestResourceCedarPolicySet_creation(t *testing.T) {

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
				Config: tests.RenderTemplate(testContext, "resource_policy_set.tmpl.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"boxer_policy_set.example",
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					assertions.ValidatePolicySetCedarIsParseable("boxer_policy_set.example"),
					assertions.ValidatePolicySetJsonIsParseable("boxer_policy_set.example"),
				},
			},
		},
	})
}
