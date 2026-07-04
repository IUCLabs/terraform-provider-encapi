# terraform-provider-encapi

Terraform / OpenTofu provider for [encapi](https://git.unkin.net/unkin/encapi),
the Puppet External Node Classifier that replaces Cobbler.

It manages three resources and two data sources over the encapi HTTP API.

## Provider configuration

```hcl
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
  # token   = "..."   # defaults to $ENCAPI_WRITE_TOKEN
}
```

Reads are unauthenticated; writes require the token. Set it out-of-band via the
`ENCAPI_WRITE_TOKEN` environment variable (managed in Vault) rather than in HCL.

## Resources

| Resource         | Purpose                                                        |
|------------------|----------------------------------------------------------------|
| `encapi_status`  | A Puppet environment (Cobbler "status"): testing, production…  |
| `encapi_role`    | A class assignment target with inheritable `default_params`.   |
| `encapi_node`    | Assigns a host (certname) to a role + environment, with `params`. |

## Data sources

| Data source      | Purpose                                  |
|------------------|------------------------------------------|
| `encapi_node`    | Look up a node's role/environment/params. |
| `encapi_role`    | Look up a role and its default params.   |

## Parameters

`default_params` (role) and `params` (node) are JSON-object strings. Use
`jsonencode({...})` so numbers, bools, lists, and nested objects keep their
types in the rendered ENC document:

```hcl
resource "encapi_role" "minio" {
  name           = "roles::infra::storage::minio"
  default_params = jsonencode({ minio_pool = "pool1", replicas = 4, tls = true })
}
```

Node `params` override the role's `default_params` on key collisions.

## Releases

Tagging `vX.Y.Z` builds `terraform-provider-encapi_X.Y.Z_linux_amd64.zip` and
uploads it to the ArtifactAPI Terraform registry (`terraform-unkin` remote,
namespace `unkin/encapi`), which serves it as a GPG-signed provider registry.

See `examples/` for full usage.
