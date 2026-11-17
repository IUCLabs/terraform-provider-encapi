resource "encapi_node" "example" {
  certname    = "node01.example.com"
  role        = "roles::infra::storage::vault"
  environment = "testing"

  # Optional per-node overrides (win over the role's default_params).
  params = jsonencode({
    vault_role = "leader"
  })
}
