package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"io"
	"net/http"
	"net/url"
	"terraform-provider-boxer/internal"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getExternalToken(services *Services) string {
	form := url.Values{}
	form.Add("client_id", services.ExternalIdp.Credentials.ClientID)
	form.Add("client_secret", services.ExternalIdp.Credentials.ClientSecret)
	form.Add("username", services.ExternalIdp.Credentials.Username)
	form.Add("password", services.ExternalIdp.Credentials.Password)
	form.Add("grant_type", services.ExternalIdp.Credentials.GrantType)

	endpoint := fmt.Sprintf("%s/auth/realms/master/protocol/openid-connect/token", services.ExternalIdp.Endpoint)
	resp, err := http.PostForm(endpoint, form)
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

func providerFactory() tfprotov6.ProviderServer {
	p, _ := providerserver.NewProtocol6WithError(internal.BoxerProvider{})()
	return p
}

var protoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"boxer": func() (tfprotov6.ProviderServer, error) {
		return providerFactory(), nil
	},
}
