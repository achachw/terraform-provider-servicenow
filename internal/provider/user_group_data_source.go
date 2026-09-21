package provider

import (
	"context"
	"encoding/json"

	"github.com/achachw/terraform-provider-servicenow/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*UserGroupDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*UserGroupDataSource)(nil)

type UserGroupDataSource struct {
	client *client.Client
}

type UserGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	View        types.String `tfsdk:"view"`
	SysID       types.String `tfsdk:"sys_id"`
	Description types.String `tfsdk:"description"`
	Email       types.String `tfsdk:"email"`
	Manager     types.String `tfsdk:"manager"`
	Parent      types.String `tfsdk:"parent"`
	Active      types.String `tfsdk:"active"`
	RawJSON     types.String `tfsdk:"raw_json"`
}

func NewUserGroupDataSource() datasource.DataSource {
	return &UserGroupDataSource{}
}

func (d *UserGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *UserGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches one ServiceNow user group by name from sys_user_group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Terraform datasource ID. Same value as sys_id.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "ServiceNow user group name.",
			},
			"view": schema.StringAttribute{
				Optional:    true,
				Description: "Optional ServiceNow view used for the table request.",
			},
			"sys_id": schema.StringAttribute{
				Computed:    true,
				Description: "ServiceNow sys_id.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "ServiceNow user group description.",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "ServiceNow user group email.",
			},
			"manager": schema.StringAttribute{
				Computed:    true,
				Description: "ServiceNow user group manager display value or sys_id.",
			},
			"parent": schema.StringAttribute{
				Computed:    true,
				Description: "ServiceNow parent group display value or sys_id.",
			},
			"active": schema.StringAttribute{
				Computed:    true,
				Description: "Whether the ServiceNow user group is active.",
			},
			"raw_json": schema.StringAttribute{
				Computed:    true,
				Description: "Raw ServiceNow record as JSON.",
			},
		},
	}
}

func (d *UserGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UserGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetUserGroupByName(ctx, state.Name.ValueString(), state.View.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read User Group", err.Error())
		return
	}

	rawJSON, err := json.Marshal(group.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Encode User Group", err.Error())
		return
	}

	state.ID = types.StringValue(group.SysID)
	state.SysID = types.StringValue(group.SysID)
	state.Name = types.StringValue(group.Name)
	state.Description = types.StringValue(group.Description)
	state.Email = types.StringValue(group.Email)
	state.Manager = types.StringValue(group.Manager)
	state.Parent = types.StringValue(group.Parent)
	state.Active = types.StringValue(group.Active)
	state.RawJSON = types.StringValue(string(rawJSON))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
