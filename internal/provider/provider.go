// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	v1 "terraform-provider-servicepipe/internal/pkg/sdkv1"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ServicepipeProvider satisfies various provider interfaces.
var _ provider.Provider = &servicepipeProvider{}

// servicepipeProvider defines the provider implementation.
type servicepipeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &servicepipeProvider{
			version: version,
		}
	}
}

// servicepipeProviderModel describes the provider data model.
type servicepipeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *servicepipeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "servicepipe"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *servicepipeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Base url to work with auth API. https://api.servicepipe.ru/api/v1 used by default provider attribute",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Service api token",
			},
		},
	}
}

// Configure prepares a Servicepipe API client for data sources and resources.
func (p *servicepipeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring servicepipe client")
	// Retrieve provider data from configuration
	var config servicepipeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown servicepipe API Token",
			"The provider cannot create the servicepipe API client as there is an unknown configuration value for the servicepipe API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the servicepipe_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing servicepipe API Endpoint",
			"The provider cannot create the servicepipe API client as there is a missing or empty value for the servicepipe API endpoint. "+
				"Set the endpoint value in the configuration or use the servicepipe_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		// endpoint = config.Endpoint.ValueString()
	}

	if config.Token.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing servicepipe API Token",
			"The provider cannot create the servicepipe API client as there is a missing or empty value for the servicepipe API token. "+
				"Set the token value in the configuration or use the servicepipe_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "servicepipe_endpoint", config.Endpoint.ValueString())
	ctx = tflog.SetField(ctx, "servicepipe_token", config.Token.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "servicepipe_token")

	tflog.Debug(ctx, "Creating servicepipe client")

	// Create a new servicepipe client using the configuration values
	client := v1.NewClientV1(config.Token.ValueString(), config.Endpoint.ValueString())
	ok, _, err := client.Echo(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create servicepipe API Client",
			"An unexpected error occurred when creating the servicepipe API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"servicepipe Client Error: "+err.Error(),
		)
		return
	}

	if !ok {
		resp.Diagnostics.AddError(
			"Unable to connect to servicepipe API Client",
			"An unexpected error occurred when creating the servicepipe API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"servicepipe Client Error: "+err.Error(),
		)
		return
	}

	// Make the servicepipe client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured servicepipe client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *servicepipeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewL7resourceResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *servicepipeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
