resource "encapi_node" "example" {
  certname    = "ausyd1nxvm2000.main.unkin.net"
  role        = "roles::infra::storage::vault"
  environment = "testing"

  # Optional per-node overrides (win over the role's default_params).
  params = jsonencode({
    vault_role = "leader"
  })
}
