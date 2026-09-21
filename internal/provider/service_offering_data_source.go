package provider

import (
	"context"
	"encoding/json"

	"github.com/achachw/terraform-provider-servicenow/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*ServiceOfferingDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*ServiceOfferingDataSource)(nil)

type ServiceOfferingDataSource struct {
	client *client.Client
}

type ServiceOfferingDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	View              types.String `tfsdk:"view"`
	SysID             types.String `tfsdk:"sys_id"`
	ShortDescription  types.String `tfsdk:"short_description"`
	OperationalStatus types.String `tfsdk:"operational_status"`
	OwnedBy           types.String `tfsdk:"owned_by"`
	ManagedBy         types.String `tfsdk:"managed_by"`
	RawJSON           types.String `tfsdk:"raw_json"`
}

func NewServiceOfferingDataSource() datasource.DataSource {
	return &ServiceOfferingDataSource{}
}

func (d *ServiceOfferingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_offering"
}

func (d *ServiceOfferingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches one ServiceNow service offering by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Terraform datasource ID. Same value as sys_id.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Service offering name.",
			},
			"view": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ServiceNow view used for the table request.",
			},
			"sys_id": schema.StringAttribute{
				Computed:    true,
				Description: "ServiceNow sys_id.",
			},
			"short_description": schema.StringAttribute{
				Computed:    true,
				Description: "Service offering short description.",
			},
			"operational_status": schema.StringAttribute{
				Computed:    true,
				Description: "Service offering operational status.",
			},
			"owned_by": schema.StringAttribute{
				Computed:    true,
				Description: "Service offering owner display value or sys_id.",
			},
			"managed_by": schema.StringAttribute{
				Computed:    true,
				Description: "Service offering manager display value or sys_id.",
			},
			"raw_json": schema.StringAttribute{
				Computed:    true,
				Description: "Raw ServiceNow record as JSON.",
			},
		},
	}
}

func (d *ServiceOfferingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	serviceNowClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.Client provider data.",
		)
		return
	}

	d.client = serviceNowClient
}

func (d *ServiceOfferingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServiceOfferingDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	view := state.View.ValueString()
	if view == "" {
		view = "automate"
	}

	offering, err := d.client.GetServiceOfferingByName(ctx, state.Name.ValueString(), view)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Service Offering", err.Error())
		return
	}

	rawJSON, err := json.Marshal(offering.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Encode Service Offering", err.Error())
		return
	}

	state.ID = types.StringValue(offering.SysID)
	state.View = types.StringValue(view)
	state.SysID = types.StringValue(offering.SysID)
	state.Name = types.StringValue(offering.Name)
	state.ShortDescription = types.StringValue(offering.ShortDescription)
	state.OperationalStatus = types.StringValue(offering.OperationalStatus)
	state.OwnedBy = types.StringValue(offering.OwnedBy)
	state.ManagedBy = types.StringValue(offering.ManagedBy)
	state.RawJSON = types.StringValue(string(rawJSON))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
