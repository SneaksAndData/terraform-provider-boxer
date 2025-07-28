package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"terraform-provider-boxer/pkg/generated/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &BoxerProvider{}
)

// BoxerProvider struct implements the Boxer Terraform provider
type BoxerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (b BoxerProvider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "boxer"
	response.Version = b.version
}

func (b BoxerProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{}
}

type boxerProviderModel struct {
	IssuerHost types.String `tfsdk:"issuer_host"`
}

func (b BoxerProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config boxerProviderModel
	diags := request.Config.Get(ctx, &config)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if config.IssuerHost.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("issuer_host"),
			"Unknown Boxer Issuer host",
			"The issuer_host parameter is mandatory"+
				"Either target apply the source of the value first, set the value statically in the configuration,"+
				"or use the BOXER_ISSUER_HOST environment variable.",
		)
	}

	if response.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("BOXER_ISSUER_HOST")

	if !config.IssuerHost.IsNull() {
		host = config.IssuerHost.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("issuer_host"),
			"Unknown Boxer Issuer host",
			"The issuer_host parameter is mandatory"+
				"Either target apply the source of the value first, set the value statically in the configuration,"+
				"or use the BOXER_ISSUER_HOST environment variable.",
		)
	}

	if response.Diagnostics.HasError() {
		return
	}

	issuerClient, err := api.NewClient(host)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to initialize Boxer Issuer Client",
			"An unexpected error occurred when creating the Boxer Issuer client. "+
				"Boxer Issuer Client Error: "+err.Error(),
		)
		return
	}

	response.DataSourceData = issuerClient
	response.ResourceData = issuerClient
}

func (b BoxerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (b BoxerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BoxerProvider{
			version: version,
		}
	}
}
