package tests

import (
	"encoding/json"
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

func getExternalToken() string {
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

func providerFactory() tfprotov6.ProviderServer {
	p, _ := providerserver.NewProtocol6WithError(internal.BoxerProvider{})()
	return p
}

var protoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"boxer": func() (tfprotov6.ProviderServer, error) {
		return providerFactory(), nil
	},
}
