package provider

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// paramsToJSON renders an API params map to a canonical JSON string for
// Terraform state. Go's json.Marshal emits object keys in sorted order, which
// matches Terraform's jsonencode(), so a config written with jsonencode()
// round-trips without spurious diffs. An empty/nil map becomes a null string.
func paramsToJSON(m map[string]any) types.String {
	if len(m) == 0 {
		return types.StringNull()
	}
	b, err := json.Marshal(m)
	if err != nil {
		return types.StringNull()
	}
	return types.StringValue(string(b))
}

// jsonToParams parses a JSON-object string attribute into an API params map. A
// null/empty value yields nil. Invalid JSON yields an error the caller surfaces.
func jsonToParams(s types.String) (map[string]any, error) {
	if s.IsNull() || s.IsUnknown() || s.ValueString() == "" {
		return nil, nil
	}
	m := map[string]any{}
	if err := json.Unmarshal([]byte(s.ValueString()), &m); err != nil {
		return nil, err
	}
	return m, nil
}

// normalizeParamsJSON canonicalizes a user-supplied JSON string so plan and
// state compare equal regardless of key order or whitespace. Returns the input
// unchanged if it does not parse (validation happens elsewhere).
func normalizeParamsJSON(s types.String) types.String {
	m, err := jsonToParams(s)
	if err != nil {
		return s
	}
	return paramsToJSON(m)
}
