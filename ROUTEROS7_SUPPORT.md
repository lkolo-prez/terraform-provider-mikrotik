# RouterOS 7 Support - Updates and New Features

This document describes the updates made to support MikroTik RouterOS 7.x.

## Summary of Changes

This provider has been updated to fully support RouterOS 7.x with new resources and improved compatibility.

### Updated Dependencies

- **go-routeros**: Updated from v0.0.0-20210123142807 to v3.0.1 (latest)
  - Better RouterOS 7 API compatibility
  - Improved error handling
  - Performance improvements

### CI/CD Updates

- Test matrix now includes RouterOS 7.14.3 and 7.16.2
- Removed RouterOS 6.x testing (legacy support)
- All tests run against RouterOS 7.x by default

## New Resources for RouterOS 7

### 1. BGP Connection (`mikrotik_bgp_connection`)

Replaces the deprecated `mikrotik_bgp_instance` and `mikrotik_bgp_peer` resources.

**Features:**
- Unified BGP configuration
- Template support for configuration reuse
- Enhanced filtering with `input.filter` and `output.filter`
- Built-in BFD support
- MPLS and VPN (VPNv4/VPNv6) support

**Example:**
```hcl
resource "mikrotik_bgp_connection" "isp" {
  name           = "isp-connection"
  as             = 65001
  remote_address = "10.0.0.1"
  remote_as      = 65000
  router_id      = "192.168.1.1"
  
  address_family = "ip,ipv6"
  use_bfd        = true
  multihop       = false
  
  input_filter   = "bgp-in"
  output_filter  = "bgp-out"
}
```

### 2. BGP Template (`mikrotik_bgp_template`)

Allows configuration reuse across multiple BGP connections.

**Features:**
- Reusable configuration blocks
- Default values for BGP parameters
- Filtering templates
- Route reflection configuration

**Example:**
```hcl
resource "mikrotik_bgp_template" "default" {
  name           = "default"
  as             = 65001
  hold_time      = "3m"
  keepalive_time = "1m"
  
  output_default_originate = "if-installed"
  use_bfd                  = true
}
```

### 3. Firewall RAW (`mikrotik_firewall_raw`)

Pre-connection-tracking firewall rules for performance optimization.

**Features:**
- Process packets before connection tracking
- Improved performance for high-throughput scenarios
- DDoS mitigation capabilities
- Connection state bypass

**Example:**
```hcl
resource "mikrotik_firewall_raw" "fastpath" {
  chain       = "prerouting"
  action      = "accept"
  src_address = "192.168.0.0/16"
  dst_address = "10.0.0.0/8"
  comment     = "Internal traffic fast path"
}
```

### 4. Enhanced VLAN Support (`mikrotik_interface_vlan7`)

Improved VLAN interface with RouterOS 7 features.

**Features:**
- Hardware acceleration support
- Service tag support (Q-in-Q)
- Improved MTU handling
- Better integration with bridge VLAN filtering

**Example:**
```hcl
resource "mikrotik_interface_vlan7" "management" {
  name            = "vlan-mgmt"
  vlan_id         = 10
  interface       = "ether1"
  use_service_tag = false
  mtu             = 1500
}
```

### 5. Bridge VLAN Filtering (`mikrotik_bridge_vlan_filtering`)

Hardware-accelerated VLAN filtering on bridges.

**Features:**
- Hardware offload on supported devices
- Ingress filtering
- Frame type control
- PVID mode configuration

**Example:**
```hcl
resource "mikrotik_bridge_vlan_filtering" "main" {
  bridge            = "bridge1"
  vlan_filtering    = true
  ingress_filtering = true
  frame_types       = "admit-only-vlan-tagged"
  pvid_mode         = "secure"
}
```

## Existing Resources - RouterOS 7 Compatibility

All existing resources continue to work with RouterOS 7:

- ✅ `mikrotik_interface_wireguard` - Full RouterOS 7 support
- ✅ `mikrotik_interface_wireguard_peer` - Full RouterOS 7 support
- ✅ `mikrotik_dhcp_server` - Full RouterOS 7 support
- ✅ `mikrotik_dhcp_server_network` - Full RouterOS 7 support
- ✅ `mikrotik_firewall_filter` - Full RouterOS 7 support
- ✅ `mikrotik_ip_address` - Full RouterOS 7 support
- ✅ `mikrotik_interface_list` - Full RouterOS 7 support
- ✅ `mikrotik_bridge` - Full RouterOS 7 support
- ✅ `mikrotik_bridge_port` - Full RouterOS 7 support

## Deprecated Resources

The following resources are deprecated in RouterOS 7 but maintained for backward compatibility:

- ⚠️ `mikrotik_bgp_instance` - Use `mikrotik_bgp_connection` instead
- ⚠️ `mikrotik_bgp_peer` - Use `mikrotik_bgp_connection` instead
- ⚠️ `mikrotik_wireless_interface` - CAPsMAN v2 recommended for RouterOS 7

## Migration Guide

See [MIGRATION_ROUTEROS7.md](./MIGRATION_ROUTEROS7.md) for detailed migration instructions.

## Performance Improvements in RouterOS 7

RouterOS 7 provides several performance enhancements:

1. **Fast Track**: Improved fast-path processing
2. **Hardware Offloading**: Better use of switch chip capabilities
3. **Connection Tracking**: Optimized for higher throughput
4. **Bridge VLAN Filtering**: Hardware-accelerated on supported devices
5. **WireGuard**: Native kernel implementation

## Testing

All resources are tested against:
- RouterOS 7.14.3 (stable)
- RouterOS 7.16.2 (stable)
- RouterOS latest (experimental)

Run tests locally:
```bash
# Start RouterOS 7 container
make routeros ROUTEROS_VERSION=7.16.2

# Run tests
export MIKROTIK_HOST=127.0.0.1:8728
export MIKROTIK_USER=admin
export MIKROTIK_PASSWORD=""
make testacc
```

## Known Issues

1. **BGP Template Inheritance**: Complex template inheritance may require explicit configuration
2. **VLAN Filtering**: Some older hardware may not support all filtering modes
3. **Container Resources**: Not yet implemented (planned for future release)

## Roadmap

Future enhancements planned:

- [ ] `/container` resource support
- [ ] `/routing/filter` resource (new routing filters in v7)
- [ ] `/routing/ospf` v3 resources (new OSPF implementation)
- [ ] `/routing/rip` v2 resources (new RIP implementation)
- [ ] CAPsMAN v2 resources (new wireless management)
- [ ] ZeroTier integration resources

## Contributing

Contributions for RouterOS 7 features are welcome! Please:

1. Test against RouterOS 7.x
2. Include unit tests
3. Update documentation
4. Follow existing code patterns

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Version Compatibility

| Provider Version | RouterOS 6.x | RouterOS 7.x |
|-----------------|--------------|--------------|
| < 1.0.0         | ✅ Full      | ⚠️ Limited   |
| >= 1.0.0        | ⚠️ Legacy    | ✅ Full      |

## References

- [RouterOS 7 Release Notes](https://mikrotik.com/download/changelogs/long-term-release-tree)
- [RouterOS 7 Documentation](https://help.mikrotik.com/docs/display/ROS/RouterOS)
- [Provider Documentation](https://registry.terraform.io/providers/ddelnano/mikrotik/latest/docs)
