# Example: Simple Routing Table / VRF

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

provider "mikrotik" {
  host     = "192.168.88.1"
  username = "admin"
  password = "admin"
  tls      = false
}

# Create a simple routing table
resource "mikrotik_routing_table" "example" {
  name    = "example_vrf"
  fib     = "main"
  comment = "Example VRF for testing"
}

# Disabled routing table (won't be used for forwarding)
resource "mikrotik_routing_table" "disabled_example" {
  name     = "disabled_vrf"
  fib      = "main"
  disabled = true
  comment  = "Disabled VRF - for future use"
}

output "example_vrf_name" {
  value = mikrotik_routing_table.example.name
}
