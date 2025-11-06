package tests

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		//CheckDestroy:             testAccCheckExampleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResource(rName),
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

func testAccExampleResource(name string) string {
	return `

provider "boxer" {

  external_auth = {
    security_token = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImoyaHM1VnJpQ09RX0dJelBsVVBRV2trbjZCa2FDMng5RkFoWWZLblJNeFkifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzYyMzU2Mzk1LCJpYXQiOjE3NjIzNTI3OTUsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiYjQzOWMyZjYtNTBhYS00ZTVhLWJiZGYtY2ViN2JkNTU2MGUyIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImludGVncmF0aW9uLXRlc3RzLWJveGVyLWlzc3VlciIsInVpZCI6IjI4OGJhMjc4LTQ3ZTktNDBkOC1hZWZhLTA1ZjhmNzNiODRlOSJ9fSwibmJmIjoxNzYyMzUyNzk1LCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDppbnRlZ3JhdGlvbi10ZXN0cy1ib3hlci1pc3N1ZXIifQ.iatWjutuU1atsJLpfdz9Cg5QY9vduIsbjH7nsBwm_ENFWa8sZ_caFGEpac30QWQpE8_mK3TQow7A2ez6hgbrJhmYI14Y75KKuw4rUteYZ9i1-PkK5wSNYRQxbO2-fYjkvvYYzvzUom0GirrN8Myi_gC-9mABl6yzJrG8msBgn0cRpz0siLdPyNUK6wdJe8PJl2ubLlVArKzpmvVuQZLQun4qdnH5BGzwOlnCCvKrJK2G6TfEUEATzkbQ7XjVQLnqxsF1LmJm0840WP_D8nXMnODTsGPge1S2UdzW2ys-lqSFDYEJSxrHTeexqBS3RdLhIdl_kZbaV0N31W31amwQ3g"
    identity_provider_id = "root"
    internal_token_provider_endpoint = "http://localhost:5555/issuer"
  }

  issuer_host    = "http://localhost:5555/issuer"
  validator_host = "http://localhost:5555/validator"
}

resource "boxer_identity_provider" "example" {
  id = "provider"
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

}

func providerFactory() tfprotov6.ProviderServer {
	p, _ := providerserver.NewProtocol6WithError(internal.BoxerProvider{})()
	return p
}
