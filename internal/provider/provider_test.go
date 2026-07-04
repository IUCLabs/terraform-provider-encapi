package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProviderSchema(t *testing.T) {
	p := New("test")()
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	if _, ok := resp.Schema.Attributes["endpoint"]; !ok {
		t.Error("missing endpoint attribute")
	}
	if _, ok := resp.Schema.Attributes["token"]; !ok {
		t.Error("missing token attribute")
	}
}

func TestProviderMetadata(t *testing.T) {
	p := New("1.2.3")()
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), provider.MetadataRequest{}, resp)
	if resp.TypeName != "encapi" || resp.Version != "1.2.3" {
		t.Errorf("metadata = %+v", resp)
	}
}

func TestProviderRegistersResourcesAndDataSources(t *testing.T) {
	p := New("test")()
	if len(p.Resources(context.Background())) != 3 {
		t.Error("expected 3 resources (role, status, node)")
	}
	if len(p.DataSources(context.Background())) != 2 {
		t.Error("expected 2 data sources (node, role)")
	}
}

// TestResourceSchemas validates each resource's schema compiles cleanly.
func TestResourceSchemas(t *testing.T) {
	for _, ctor := range []func() resource.Resource{NewRoleResource, NewStatusResource, NewNodeResource} {
		resp := &resource.SchemaResponse{}
		ctor().Schema(context.Background(), resource.SchemaRequest{}, resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("resource schema error: %v", resp.Diagnostics)
		}
	}
}
