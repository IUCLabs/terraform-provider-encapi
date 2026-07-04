package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

type roleResource struct {
	client *apiClient
}

type roleResourceModel struct {
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	DefaultParams types.String `tfsdk:"default_params"`
}

func NewRoleResource() resource.Resource { return &roleResource{} }

func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Puppet class assignment target (e.g. roles::infra::storage::vault) with inheritable default parameters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Fully-qualified role/class name, e.g. roles::infra::storage::vault.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Human-readable description.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"default_params": schema.StringAttribute{
				Description: "Inheritable default parameters as a JSON object. Use jsonencode({...}) to preserve value types. Merged into every node carrying this role; node params win on collisions.",
				Optional:    true,
			},
		},
	}
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.upsert(ctx, plan, &resp.Diagnostics, &resp.State)
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.upsert(ctx, plan, &resp.Diagnostics, &resp.State)
}

// upsert PUTs the role and writes the returned state. Shared by Create/Update.
func (r *roleResource) upsert(ctx context.Context, plan roleResourceModel, diags *diag.Diagnostics, state *tfsdk.State) {
	params, err := jsonToParams(plan.DefaultParams)
	if err != nil {
		diags.AddError("invalid default_params", "default_params must be a JSON object: "+err.Error())
		return
	}
	body := roleAPI{Name: plan.Name.ValueString(), Description: plan.Description.ValueString(), DefaultParams: params}
	var out roleAPI
	if err := r.client.put(ctx, "/api/v1/roles/"+pathEscape(body.Name), body, &out); err != nil {
		diags.AddError("upsert role failed", err.Error())
		return
	}
	diags.Append(state.Set(ctx, roleAPIToModel(out))...)
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out roleAPI
	if err := r.client.get(ctx, "/api/v1/roles/"+pathEscape(state.Name.ValueString()), &out); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("read role failed", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, roleAPIToModel(out))...)
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.del(ctx, "/api/v1/roles/"+pathEscape(state.Name.ValueString())); err != nil {
		resp.Diagnostics.AddError("delete role failed", err.Error())
		return
	}
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func roleAPIToModel(api roleAPI) roleResourceModel {
	return roleResourceModel{
		Name:          types.StringValue(api.Name),
		Description:   types.StringValue(api.Description),
		DefaultParams: paramsToJSON(api.DefaultParams),
	}
}
