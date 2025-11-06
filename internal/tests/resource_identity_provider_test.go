package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"testing"
)

func TestAccExampleWidget_basic(t *testing.T) {
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	token := getExternalToken()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResource(token, randomName),
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
						knownvalue.StringExact("http://localhost:8080/realms/master/"),
					),
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("user_id_claim"),
						knownvalue.StringExact("preferred_username"),
					),
					statecheck.ExpectKnownValue(
						"boxer_identity_provider.example",
						tfjsonpath.New("issuers"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("http://localhost:8080/realms/master")}),
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

func testAccExampleResource(token string, name string) string {
	configurationTemplate := `
provider "boxer" {

  external_auth = {
    security_token = "%s"
    identity_provider_id = "keycloak"
    internal_token_provider_endpoint = "http://localhost:5555/issuer"
  }

  issuer_host    = "http://localhost:5555/issuer"
  validator_host = "http://localhost:5555/validator"
}

resource "boxer_identity_provider" "example" {
  id = "%s"
  user_id_claim = "preferred_username"
  discovery_url = "http://localhost:8080/realms/master/"
  issuers = [
    "http://localhost:8080/realms/master",
  ]
  audiences = [
    "account"
  ]
}
`

	return fmt.Sprintf(configurationTemplate, token, name)

}
