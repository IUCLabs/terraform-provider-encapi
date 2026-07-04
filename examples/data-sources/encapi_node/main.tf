data "encapi_node" "web0" {
  certname = "ausyd1nxvm2000.main.unkin.net"
}

output "web0_role" {
  value = data.encapi_node.web0.role
}

output "web0_environment" {
  value = data.encapi_node.web0.environment
}
