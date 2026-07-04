resource "encapi_role" "minio" {
  name        = "roles::infra::storage::minio"
  description = "MinIO object storage node"

  # Inheritable defaults, merged into every node with this role.
  # Use jsonencode so numbers/bools keep their type in the ENC output.
  default_params = jsonencode({
    minio_pool = "pool1"
    replicas   = 4
    tls        = true
  })
}
