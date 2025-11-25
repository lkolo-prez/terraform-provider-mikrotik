# Terraform Provider Examples for MikroTik RouterOS v7

This directory contains Terraform configuration examples demonstrating RouterOS v7 features.

## Directory Structure

```
examples/
├── README.md                      # This file
├── basic/                         # Basic configurations
│   ├── bridge-vlan-filtering/    # VLAN filtering with hardware offload
│   ├── dhcp-server/              # DHCP server setup
│   └── wireguard-vpn/            # WireGuard VPN
├── routing/                       # Routing examples
│   ├── bgp-v7/                   # BGP with new v7 API
│   ├── ospf-v7/                  # OSPF with templates
│   ├── vrf/                      # VRF (Virtual Routing)
│   └── policy-routing/           # Policy-based routing
├── firewall/                      # Firewall configurations
│   ├── basic-protection/         # Basic firewall rules
│   ├── raw-table/                # RAW table (pre-CT)
│   └── advanced-nat/             # Complex NAT scenarios
├── wireless/                      # WiFi configurations
│   ├── wifi-v7-ap/               # New WiFi system (802.11ax)
│   └── wireless-legacy/          # Legacy wireless (v6 compat)
├── advanced/                      # Advanced features
│   ├── queue-cake/               # CAKE queueing
│   ├── veth-containers/          # Virtual Ethernet for containers
│   └── dns-doh/                  # DNS over HTTPS
└── complete/                      # Complete infrastructure
    ├── small-office/             # Small office setup
    ├── isp-edge/                 # ISP edge router
    └── datacenter-tor/           # Top-of-Rack switch
```

## Prerequisites

1. **Install Terraform**:
   ```bash
   # Download from https://www.terraform.io/downloads
   terraform version  # Should be 1.0+
   ```

2. **Install MikroTik Provider**:
   ```bash
   # Will be downloaded automatically from Terraform Registry
   # Or build locally:
   go build -o terraform-provider-mikrotik.exe
   ```

3. **RouterOS Access**:
   - RouterOS 7.14.3 or newer
   - API access enabled
   - Admin credentials

## Quick Start

### 1. Choose an Example

```bash
cd examples/basic/bridge-vlan-filtering
```

### 2. Configure Provider

Create `provider.tf`:

```hcl
terraform {
  required_providers {
    mikrotik = {
      source  = "ddelnano/mikrotik"
      version = "~> 1.0"
    }
  }
}

provider "mikrotik" {
  host     = var.mikrotik_host     # e.g., "192.168.88.1:8728"
  username = var.mikrotik_username # e.g., "admin"
  password = var.mikrotik_password
  tls      = false
}
```

### 3. Set Variables

Create `terraform.tfvars`:

```hcl
mikrotik_host     = "192.168.88.1:8728"
mikrotik_username = "admin"
mikrotik_password = ""
```

**Security Note**: Never commit `terraform.tfvars` with real credentials!

### 4. Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Preview changes
terraform plan

# Apply configuration
terraform apply

# Destroy resources (when done)
terraform destroy
```

## Example Categories

### Basic Examples

Perfect for getting started:

- **Bridge VLAN Filtering**: Hardware-accelerated VLANs on CRS3xx/CRS5xx
- **DHCP Server**: Complete DHCP setup with pools and static leases
- **WireGuard VPN**: Site-to-site or road warrior VPN

### Routing Examples

Advanced routing configurations:

- **BGP v7**: New connection/template system
- **OSPF v7**: Instance/area templates
- **VRF**: Separate routing tables
- **Policy Routing**: Rule-based routing

### Firewall Examples

Security configurations:

- **Basic Protection**: Essential firewall rules
- **RAW Table**: Pre-connection tracking filtering
- **Advanced NAT**: Hairpin NAT, port forwarding

### Wireless Examples

WiFi configurations:

- **WiFi v7**: New 802.11ax system
- **Legacy Wireless**: Backward compatibility

### Advanced Examples

Cutting-edge features:

- **CAKE Queueing**: Modern QoS (RouterOS 7+)
- **veth Interfaces**: Container networking
- **DNS DoH**: DNS over HTTPS

### Complete Examples

Real-world scenarios:

- **Small Office**: Router + WiFi + DHCP + Firewall
- **ISP Edge**: BGP + Firewall + QoS
- **Datacenter ToR**: VLAN + LACP + OSPF

## Best Practices

### 1. Use Variables

```hcl
variable "lan_network" {
  description = "LAN network CIDR"
  type        = string
  default     = "192.168.1.0/24"
}
```

### 2. Use Modules

Organize reusable configurations:

```hcl
module "firewall" {
  source = "./modules/firewall"
  
  wan_interface = "ether1"
  lan_interface = "bridge1"
}
```

### 3. State Management

For production, use remote state:

```hcl
terraform {
  backend "s3" {
    bucket = "my-terraform-state"
    key    = "mikrotik/prod.tfstate"
    region = "us-east-1"
  }
}
```

### 4. Testing

Test in lab before production:

```bash
# Use separate workspaces
terraform workspace new lab
terraform workspace select lab
terraform apply

terraform workspace select production
terraform apply
```

## Common Patterns

### Data Sources

Query existing resources:

```hcl
data "mikrotik_system_resource" "router" {}

output "routeros_version" {
  value = data.mikrotik_system_resource.router.version
}
```

### Loops

Create multiple similar resources:

```hcl
locals {
  vlans = {
    10 = "management"
    20 = "servers"
    30 = "clients"
  }
}

resource "mikrotik_interface_vlan7" "vlans" {
  for_each = local.vlans
  
  interface = "bridge1"
  vlan_id   = each.key
  name      = "vlan${each.key}"
  comment   = each.value
}
```

### Conditionals

Optional resources:

```hcl
resource "mikrotik_wireless_interface" "wifi" {
  count = var.enable_wifi ? 1 : 0
  
  name = "wlan1"
  ssid = var.wifi_ssid
  # ...
}
```

## Troubleshooting

### Connection Issues

```bash
# Test API access
curl -k https://192.168.88.1/rest/system/resource
```

### Debug Mode

```bash
export TF_LOG=DEBUG
terraform apply
```

### Common Errors

1. **"Resource not found"**: Check RouterOS version compatibility
2. **"Permission denied"**: Ensure user has admin rights
3. **"Invalid parameter"**: Check RouterOS 7 documentation for parameter changes

## Migration from v6

If migrating configurations from RouterOS v6:

1. Review [MIGRATION_ROUTEROS7.md](../MIGRATION_ROUTEROS7.md)
2. Check [ROUTEROS7_COVERAGE.md](../ROUTEROS7_COVERAGE.md) for feature availability
3. Update resource names (e.g., `bgp_instance` → `bgp_connection`)
4. Test in lab environment first!

## Contributing

To add new examples:

1. Create a new directory under appropriate category
2. Include:
   - `main.tf` - Main configuration
   - `variables.tf` - Input variables
   - `outputs.tf` - Outputs
   - `README.md` - Documentation
3. Test thoroughly on RouterOS 7.14.3+
4. Submit pull request

## Support

- [Provider Documentation](../README.md)
- [RouterOS v7 Cheat Sheet](../ROUTEROS7_CHEATSHEET.md)
- [Feature Coverage Matrix](../ROUTEROS7_COVERAGE.md)
- [GitHub Issues](https://github.com/ddelnano/terraform-provider-mikrotik/issues)

---

**Last Updated**: November 25, 2025  
**Provider Version**: 1.0.0  
**RouterOS Versions**: 7.14.3, 7.16.2+
