data "encapi_node" "web0" {
  certname = "node01.example.com"
}

output "web0_role" {
  value = data.encapi_node.web0.role
}

output "web0_environment" {
  value = data.encapi_node.web0.environment
}
