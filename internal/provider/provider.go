package provider

import (
	"context"
	"os"

	"github.com/achachw/terraform-provider-servicenow/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = (*ServiceNowProvider)(nil)

type ServiceNowProvider struct {
	version string
}

type ServiceNowProviderModel struct {
	BaseURL  types.String `tfsdk:"base_url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ServiceNowProvider{version: version}
	}
}

func (p *ServiceNowProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "servicenow"
	resp.Version = p.version
}

func (p *ServiceNowProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for querying ServiceNow records.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "ServiceNow instance base URL. Can also be set with SERVICENOW_BASE_URL.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "ServiceNow username. Can also be set with SERVICENOW_USERNAME.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "ServiceNow password. Can also be set with SERVICENOW_PASSWORD.",
			},
		},
	}
}

func (p *ServiceNowProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ServiceNowProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseURL := os.Getenv("SERVICENOW_BASE_URL")
	username := os.Getenv("SERVICENOW_USERNAME")
	password := os.Getenv("SERVICENOW_PASSWORD")

	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if baseURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing ServiceNow Base URL",
			"Set base_url in the provider configuration or SERVICENOW_BASE_URL in the environment.",
		)
	}
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing ServiceNow Username",
			"Set username in the provider configuration or SERVICENOW_USERNAME in the environment.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing ServiceNow Password",
			"Set password in the provider configuration or SERVICENOW_PASSWORD in the environment.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	serviceNowClient, err := client.New(baseURL, username, password)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create ServiceNow Client", err.Error())
		return
	}

	resp.DataSourceData = serviceNowClient
}

func (p *ServiceNowProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *ServiceNowProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServiceOfferingDataSource,
		NewUserGroupDataSource,
	}
}
