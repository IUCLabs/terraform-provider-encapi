package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &nodeDataSource{}

type nodeDataSource struct {
	client *apiClient
}

func NewNodeDataSource() datasource.DataSource { return &nodeDataSource{} }

func (d *nodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (d *nodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up the role/environment/params assigned to a Puppet node.",
		Attributes: map[string]schema.Attribute{
			"certname":    schema.StringAttribute{Required: true, Description: "Puppet certname (fqdn)."},
			"role":        schema.StringAttribute{Computed: true, Description: "Assigned role."},
			"environment": schema.StringAttribute{Computed: true, Description: "Assigned environment/status."},
			"params":      schema.StringAttribute{Computed: true, Description: "Per-node params as JSON."},
		},
	}
}

func (d *nodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *nodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg nodeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out nodeAPI
	if err := d.client.get(ctx, "/api/v1/nodes/"+pathEscape(cfg.Certname.ValueString()), &out); err != nil {
		resp.Diagnostics.AddError("read node failed", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, nodeResourceModel{
		Certname:    types.StringValue(out.Certname),
		Role:        types.StringValue(out.Role),
		Environment: types.StringValue(out.Environment),
		Params:      paramsToJSON(out.Params),
	})...)
}
