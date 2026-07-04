package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &nodeResource{}
	_ resource.ResourceWithImportState = &nodeResource{}
)

type nodeResource struct {
	client *apiClient
}

type nodeResourceModel struct {
	Certname    types.String `tfsdk:"certname"`
	Role        types.String `tfsdk:"role"`
	Environment types.String `tfsdk:"environment"`
	Params      types.String `tfsdk:"params"`
}

func NewNodeResource() resource.Resource { return &nodeResource{} }

func (r *nodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (r *nodeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Assigns a Puppet node (by certname) to a role and environment, with optional per-node parameter overrides.",
		Attributes: map[string]schema.Attribute{
			"certname": schema.StringAttribute{
				Description: "Puppet certname (fqdn) of the host.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Description: "Role/class to assign. Must reference an existing encapi_role.",
				Required:    true,
			},
			"environment": schema.StringAttribute{
				Description: "Environment/status. Must reference an existing encapi_status.",
				Required:    true,
			},
			"params": schema.StringAttribute{
				Description: "Per-node parameter overrides as a JSON object. Use jsonencode({...}) to preserve value types. These override the role's default_params.",
				Optional:    true,
			},
		},
	}
}

func (r *nodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiClient)
	if !ok {
		resp.Diagnostics.AddError("unexpected provider data type", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *nodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan nodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.upsert(ctx, plan, &resp.Diagnostics, &resp.State)
}

func (r *nodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan nodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.upsert(ctx, plan, &resp.Diagnostics, &resp.State)
}

func (r *nodeResource) upsert(ctx context.Context, plan nodeResourceModel, diags *diag.Diagnostics, state *tfsdk.State) {
	params, err := jsonToParams(plan.Params)
	if err != nil {
		diags.AddError("invalid params", "params must be a JSON object: "+err.Error())
		return
	}
	body := nodeAPI{
		Certname:    plan.Certname.ValueString(),
		Role:        plan.Role.ValueString(),
		Environment: plan.Environment.ValueString(),
		Params:      params,
	}
	var out nodeAPI
	if err := r.client.put(ctx, "/api/v1/nodes/"+pathEscape(body.Certname), body, &out); err != nil {
		diags.AddError("upsert node failed", err.Error())
		return
	}
	diags.Append(state.Set(ctx, nodeAPIToModel(out))...)
}

func (r *nodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state nodeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out nodeAPI
	if err := r.client.get(ctx, "/api/v1/nodes/"+pathEscape(state.Certname.ValueString()), &out); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("read node failed", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, nodeAPIToModel(out))...)
}

func (r *nodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state nodeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.del(ctx, "/api/v1/nodes/"+pathEscape(state.Certname.ValueString())); err != nil {
		resp.Diagnostics.AddError("delete node failed", err.Error())
		return
	}
}

func (r *nodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("certname"), req, resp)
}

func nodeAPIToModel(api nodeAPI) nodeResourceModel {
	return nodeResourceModel{
		Certname:    types.StringValue(api.Certname),
		Role:        types.StringValue(api.Role),
		Environment: types.StringValue(api.Environment),
		Params:      paramsToJSON(api.Params),
	}
}
