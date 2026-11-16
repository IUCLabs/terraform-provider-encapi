# Terraform Provider encapi

The encapi provider manages a Puppet External Node Classifier (ENC) served by
encapi: Puppet environments (statuses), roles with inheritable default
parameters, and per-node role assignments. It works with Terraform and OpenTofu.

## Quick start

```hcl
terraform {
  required_providers {
    encapi = {
      source = "IUCLabs/encapi"
    }
  }
}

provider "encapi" {
  endpoint = "https://encapi.example.com"
  # token defaults to the ENCAPI_WRITE_TOKEN environment variable
}

resource "encapi_status" "production" {
  name = "production"
}

resource "encapi_role" "web" {
  name           = "roles::web"
  default_params = jsonencode({ port = 8080 })
}

resource "encapi_node" "web01" {
  certname    = "web01.example.com"
  role        = encapi_role.web.name
  environment = encapi_status.production.name
}
```

Reads are unauthenticated; writes need the bearer token.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0 or [OpenTofu](https://opentofu.org) >= 1.6
- [Go](https://go.dev/doc/install) >= 1.25 (to build the provider)

## Building

```sh
make build   # or: go build
make test
make docs    # regenerate docs/ from the schema and examples/
```

## Documentation

See [`docs/`](docs/index.md) for the provider, resource and data source
reference, and [`examples/`](examples/) for usage.

## License

[MIT](LICENSE)
