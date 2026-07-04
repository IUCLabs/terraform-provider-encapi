package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParamsToJSONEmptyIsNull(t *testing.T) {
	if v := paramsToJSON(nil); !v.IsNull() {
		t.Errorf("nil map -> %v, want null", v)
	}
	if v := paramsToJSON(map[string]any{}); !v.IsNull() {
		t.Errorf("empty map -> %v, want null", v)
	}
}

func TestParamsToJSONSortedKeys(t *testing.T) {
	// Go's json.Marshal sorts keys, matching Terraform's jsonencode().
	v := paramsToJSON(map[string]any{"b": 2, "a": 1})
	if v.ValueString() != `{"a":1,"b":2}` {
		t.Errorf("got %q, want sorted compact JSON", v.ValueString())
	}
}

func TestJSONToParamsRoundTrip(t *testing.T) {
	in := types.StringValue(`{"epel":"9","replicas":3,"enabled":true}`)
	m, err := jsonToParams(in)
	if err != nil {
		t.Fatal(err)
	}
	if m["epel"] != "9" || m["replicas"] != float64(3) || m["enabled"] != true {
		t.Errorf("params = %#v", m)
	}
}

func TestJSONToParamsNull(t *testing.T) {
	m, err := jsonToParams(types.StringNull())
	if err != nil || m != nil {
		t.Errorf("got %v, %v; want nil, nil", m, err)
	}
}

func TestJSONToParamsInvalid(t *testing.T) {
	if _, err := jsonToParams(types.StringValue("not json")); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestNormalizeParamsJSON(t *testing.T) {
	// reordered + whitespace should canonicalize to sorted compact form
	got := normalizeParamsJSON(types.StringValue(`{ "b": 2,  "a": 1 }`))
	if got.ValueString() != `{"a":1,"b":2}` {
		t.Errorf("normalize = %q", got.ValueString())
	}
	// invalid JSON returns input unchanged
	bad := types.StringValue("{invalid")
	if normalizeParamsJSON(bad).ValueString() != "{invalid" {
		t.Error("invalid JSON should pass through unchanged")
	}
}
