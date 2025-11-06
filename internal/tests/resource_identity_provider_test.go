package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"io"
	"net/http"
	"net/url"
	"terraform-provider-boxer/internal"
	"testing"
)

func TestAccExampleWidget_basic(t *testing.T) {

	//var idpBefore, idpAfter identityProviderResource
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	protoV6ProviderFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"boxer": func() (tfprotov6.ProviderServer, error) {
			return providerFactory(), nil
		},
	}
	token := getToken()
	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		//CheckDestroy:             testAccCheckExampleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResource(token, rName),
				//ConfigStateChecks: []statecheck.StateCheck{
				//	stateCheckExampleResourceExists("example_widget.foo", &widgetBefore),
				//},
			},
			//{
			//	Config: testAccExampleResource_removedPolicy(rName),
			//	ConfigStateChecks: []statecheck.StateCheck{
			//		stateCheckExampleResourceExists("example_widget.foo", &widgetAfter),
			//	},
			//},
		},
	})
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getToken() string {
	form := url.Values{}
	form.Add("client_id", "test_client")
	form.Add("client_secret", "test_client_secret")
	form.Add("username", "test_root")
	form.Add("password", "test-root-password")
	form.Add("grant_type", "password")

	resp, err := http.PostForm("http://localhost:5555/auth/realms/master/protocol/openid-connect/token", form)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		panic(err)
	}
	return tr.AccessToken
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

func providerFactory() tfprotov6.ProviderServer {
	p, _ := providerserver.NewProtocol6WithError(internal.BoxerProvider{})()
	return p
}
