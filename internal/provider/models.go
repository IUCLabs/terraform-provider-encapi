package provider

// API wire types, mirroring the encapi API models. They are
// duplicated here to keep the provider's dependency surface small.

type roleAPI struct {
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	DefaultParams map[string]any `json:"default_params,omitempty"`
}

type statusAPI struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type nodeAPI struct {
	Certname    string         `json:"certname"`
	Role        string         `json:"role"`
	Environment string         `json:"environment"`
	Params      map[string]any `json:"params,omitempty"`
}
