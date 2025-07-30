package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-boxer/pkg/generated/api"
)

func getIssuerClient(request datasource.ConfigureRequest, response *datasource.ConfigureResponse) *issuer.Client {
	if request.ProviderData == nil {
		return nil
	}
	data, ok := request.ProviderData.(*BoxerProviderData)
	if !ok {
		response.Diagnostics.AddError(
			"Invalid Provider Data",
			"The provider data must be of type *BoxerProviderData, but was %s. This is most likely the bug in the provider implementation.",
		)
		return nil
	}
	if data.issuerClient == nil {
		response.Diagnostics.AddError(
			"Invalid Issuer Client",
			"The issuer client must not be nil. This is most likely the bug in the provider implementation.",
		)
		return nil
	}
	client := data.issuerClient
	return client
}

func readPlan(ctx context.Context, basePlan tfsdk.Plan, diagnostics *diag.Diagnostics) (*identityProviderResourceModel, error) {
	var plan identityProviderResourceModel
	diags := basePlan.Get(ctx, &plan)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return nil, fmt.Errorf("error getting plan")
	}
	return &plan, nil
}
