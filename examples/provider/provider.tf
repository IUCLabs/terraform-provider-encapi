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
