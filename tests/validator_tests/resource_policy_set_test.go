package issuer_tests

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	tests2 "terraform-provider-boxer/tests"
	assertions2 "terraform-provider-boxer/tests/assertions"
	"testing"
)

func TestResourceCedarPolicySet_creation(t *testing.T) {

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
				Config: tests2.RenderTemplate(testContext, "resource_policy_set.tmpl.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"boxer_policy_set.example",
						tfjsonpath.New("id"),
						knownvalue.StringExact(randomName),
					),
					assertions2.ValidatePolicySetCedarIsParseable("boxer_policy_set.example"),
					assertions2.ValidatePolicySetJsonIsParseable("boxer_policy_set.example"),
				},
			},
		},
	})
}
