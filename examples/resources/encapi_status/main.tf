# Statuses map to Puppet environments. Define as many as your estate uses;
# a node can only be pinned to one that exists.
resource "encapi_status" "testing" {
  name = "testing"
}

resource "encapi_status" "production" {
  name        = "production"
  description = "Production environment"
}

resource "encapi_status" "development" {
  name = "development"
}
