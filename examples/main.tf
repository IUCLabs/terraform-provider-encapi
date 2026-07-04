terraform {
  required_providers {
    encapi = {
      source  = "git.unkin.net/unkin/encapi"
      version = "0.0.1"
    }
  }
}

provider "encapi" {
  endpoint = "https://encapi.k8s.syd1.au.unkin.net"
  # token defaults to the ENCAPI_WRITE_TOKEN environment variable
}

# ---------------------------------------------------------------------------
# Statuses (Puppet environments / Cobbler "status")
# ---------------------------------------------------------------------------

resource "encapi_status" "testing" {
  name        = "testing"
  description = "Implicit environment; dropped from ENC output so agents use their default."
}

resource "encapi_status" "production" {
  name = "production"
}

resource "encapi_status" "development" {
  name = "development"
}

# ---------------------------------------------------------------------------
# Roles, with inheritable default params (use jsonencode to keep value types)
# ---------------------------------------------------------------------------

resource "encapi_role" "vault" {
  name        = "roles::infra::storage::vault"
  description = "HashiCorp Vault cluster member"
}

resource "encapi_role" "minio" {
  name        = "roles::infra::storage::minio"
  description = "MinIO object storage node"

  default_params = jsonencode({
    minio_pool = "pool1"
    replicas   = 4
  })
}

# ---------------------------------------------------------------------------
# Node assignments
# ---------------------------------------------------------------------------

resource "encapi_node" "vault0" {
  certname    = "ausyd1nxvm2000.main.unkin.net"
  role        = encapi_role.vault.name
  environment = encapi_status.testing.name
}

resource "encapi_node" "minio0" {
  certname    = "ausyd1nxvm2100.main.unkin.net"
  role        = encapi_role.minio.name
  environment = encapi_status.production.name

  # per-node override of the role default
  params = jsonencode({
    minio_pool = "pool2"
  })
}

# ---------------------------------------------------------------------------
# Data sources
# ---------------------------------------------------------------------------

data "encapi_node" "vault0" {
  certname = encapi_node.vault0.certname
}

output "vault0_role" {
  value = data.encapi_node.vault0.role
}
