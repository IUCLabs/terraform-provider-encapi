package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &roleDataSource{}

type roleDataSource struct {
	client *apiClient
}

func NewRoleDataSource() datasource.DataSource { return &roleDataSource{} }

func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up a role and its inheritable default params.",
		Attributes: map[string]schema.Attribute{
			"name":           schema.StringAttribute{Required: true, Description: "Role/class name."},
			"description":    schema.StringAttribute{Computed: true, Description: "Description."},
			"default_params": schema.StringAttribute{Computed: true, Description: "Default params as JSON."},
		},
	}
}

func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiClient)
	if !ok {
		resp.Diagnostics.AddError("unexpected provider data type", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg roleResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out roleAPI
	if err := d.client.get(ctx, "/api/v1/roles/"+pathEscape(cfg.Name.ValueString()), &out); err != nil {
		resp.Diagnostics.AddError("read role failed", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, roleAPIToModel(out))...)
}
