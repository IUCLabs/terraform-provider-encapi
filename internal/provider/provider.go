package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &encapiProvider{}

type encapiProvider struct {
	version string
}

type encapiProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

// New returns the provider constructor.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &encapiProvider{version: version}
	}
}

func (p *encapiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "encapi"
	resp.Version = p.version
}

func (p *encapiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage the Puppet External Node Classifier (encapi): roles, statuses (environments), and node role assignments.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The encapi server base URL (e.g. https://encapi.k8s.syd1.au.unkin.net).",
				Required:    true,
			},
			"token": schema.StringAttribute{
				Description: "Write token (bearer) for mutating operations. Defaults to the ENCAPI_WRITE_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *encapiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config encapiProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := config.Token.ValueString()
	if token == "" {
		token = os.Getenv("ENCAPI_WRITE_TOKEN")
	}

	client := newAPIClient(config.Endpoint.ValueString(), token)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *encapiProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRoleResource,
		NewStatusResource,
		NewNodeResource,
	}
}

func (p *encapiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNodeDataSource,
		NewRoleDataSource,
	}
}
