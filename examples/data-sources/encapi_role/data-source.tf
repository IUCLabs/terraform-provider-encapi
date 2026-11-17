data "encapi_role" "vault" {
  name = "roles::infra::storage::vault"
}

output "vault_default_params" {
  value = jsondecode(data.encapi_role.vault.default_params)
}
