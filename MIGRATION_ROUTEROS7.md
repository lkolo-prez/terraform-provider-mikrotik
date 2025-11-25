# Migration Guide: RouterOS 6 to RouterOS 7

This guide helps you migrate your Terraform MikroTik provider configuration from RouterOS 6 to RouterOS 7.

## Overview of Changes

RouterOS 7 introduced significant architectural changes, particularly in the BGP implementation. The legacy BGP commands (`/routing/bgp/instance` and `/routing/bgp/peer`) have been replaced with a new, more flexible system.

## Breaking Changes

### 1. BGP Configuration (RouterOS 7+)

#### Old Approach (RouterOS 6, deprecated in v7):

```hcl
resource "mikrotik_bgp_instance" "main" {
  name      = "main"
  as        = 65530
  router_id = "172.16.0.1"
  
  redistribute_connected = true
  redistribute_static    = true
}

resource "mikrotik_bgp_peer" "upstream" {
  name           = "upstream-peer"
  instance       = mikrotik_bgp_instance.main.name
  remote_address = "10.0.0.1"
  remote_as      = 65531
  ttl            = "255"
}
```

#### New Approach (RouterOS 7):

```hcl
# Optional: Create a template for reusable configuration
resource "mikrotik_bgp_template" "default" {
  name              = "default-template"
  as                = 65530
  router_id         = "172.16.0.1"
  hold_time         = "3m"
  keepalive_time    = "1m"
  address_families  = "ip,ipv6"
  
  output_default_originate = "always"
}

# Create BGP connection
resource "mikrotik_bgp_connection" "upstream" {
  name           = "upstream-connection"
  as             = 65530
  remote_address = "10.0.0.1"
  remote_as      = 65531
  router_id      = "172.16.0.1"
  
  templates      = mikrotik_bgp_template.default.name
  
  ttl            = "255"
  address_family = "ip,ipv6"
  
  output_default_originate = "always"
  output_network           = "connected,static"
  
  use_bfd        = true
  multihop       = true
}
```

**Key Differences:**
- **Unified Resource**: BGP Instance and Peer are merged into `bgp_connection`
- **Templates**: New template system for configuration reuse
- **Input/Output Filters**: More granular control with `input.filter` and `output.filter`
- **Network Advertisement**: Use `output.network` instead of `redistribute_*` options
- **Enhanced Features**: Built-in BFD support, MPLS, VPNv4/v6

### 2. New Firewall Features

#### RAW Table (RouterOS 7+)

RouterOS 7 introduces the RAW firewall table for pre-connection-tracking processing:

```hcl
resource "mikrotik_firewall_raw" "prerouting_accept" {
  chain            = "prerouting"
  action           = "accept"
  src_address      = "192.168.1.0/24"
  dst_address      = "10.0.0.0/8"
  protocol         = "tcp"
  dst_port         = "80,443"
  comment          = "Fast path for trusted internal traffic"
}
```

**Use Cases:**
- Bypass connection tracking for performance
- Early packet filtering before connection state
- DDoS mitigation
- Fast-path optimization

### 3. Enhanced VLAN Support

RouterOS 7 improves VLAN handling with hardware acceleration:

```hcl
# Modern VLAN interface
resource "mikrotik_interface_vlan7" "vlan100" {
  name            = "vlan100"
  vlan_id         = 100
  interface       = "ether1"
  use_service_tag = false
  mtu             = 1500
  comment         = "Management VLAN"
}

# Bridge VLAN filtering (hardware accelerated)
resource "mikrotik_bridge_vlan_filtering" "main_bridge" {
  bridge             = "bridge1"
  vlan_filtering     = true
  ingress_filtering  = true
  frame_types        = "admit-only-vlan-tagged"
  pvid_mode          = "secure"
}
```

## Migration Steps

### Step 1: Upgrade RouterOS

1. Backup your current configuration
2. Upgrade to RouterOS 7.x (latest stable recommended)
3. Test all functionality after upgrade

### Step 2: Update Terraform Provider

```bash
# Update your terraform configuration
terraform init -upgrade
```

Update your `terraform` block:

```hcl
terraform {
  required_providers {
    mikrotik = {
      source  = "ddelnano/mikrotik"
      version = "~> 1.0"  # Use latest version supporting RouterOS 7
    }
  }
}
```

### Step 3: Migrate BGP Configuration

1. **Export existing BGP configuration** from RouterOS 6
2. **Create migration script** to convert to new format
3. **Test in staging** before production

Example migration script:

```bash
#!/bin/bash
# Convert old BGP instance to new connection format

OLD_INSTANCE="main"
OLD_PEER="upstream-peer"
NEW_CONNECTION="upstream-connection"

# Export old config
/routing bgp instance export file=old-bgp-instance
/routing bgp peer export file=old-bgp-peer

# After manual conversion, import new config
/routing bgp connection import file=new-bgp-connection
```

### Step 4: Leverage New Features

Take advantage of RouterOS 7 improvements:

1. **Fast Track**: Enable for improved throughput
2. **BFD**: Use for faster BGP convergence
3. **Hardware VLAN Filtering**: Enable on supported hardware
4. **IPv6**: Enhanced IPv6 support throughout

## Compatibility Matrix

| Feature | RouterOS 6 | RouterOS 7 | Notes |
|---------|-----------|-----------|-------|
| Legacy BGP (instance/peer) | ✅ Supported | ⚠️ Deprecated | Will be removed in future |
| New BGP (connection/template) | ❌ Not Available | ✅ Supported | Recommended |
| Firewall RAW Table | ❌ Not Available | ✅ Supported | Performance feature |
| Hardware VLAN Filtering | ⚠️ Limited | ✅ Full Support | Better performance |
| WireGuard | ❌ Not Available | ✅ Supported | Native VPN |
| Container | ❌ Not Available | ✅ Supported | Run containers on router |

## Testing Your Migration

### 1. State Verification

After migration, verify your Terraform state:

```bash
terraform plan
# Should show no changes if migration is correct
```

### 2. Functionality Testing

Test critical paths:

```bash
# Test BGP connectivity
/routing bgp connection print
/routing bgp session print

# Test firewall rules
/ip firewall raw print
/ip firewall filter print

# Test VLAN functionality
/interface vlan print
/interface bridge vlan print
```

### 3. Rollback Plan

Always have a rollback plan:

1. Keep RouterOS 6 configuration backup
2. Document all changes
3. Test rollback procedure in staging

## Common Issues and Solutions

### Issue: BGP Sessions Won't Establish

**Symptoms**: BGP connections show "idle" or "active" state

**Solutions**:
1. Check `router-id` is set correctly
2. Verify `remote.address` and `local.address`
3. Enable `multihop` if needed
4. Check firewall rules (port 179)

### Issue: VLAN Traffic Not Working

**Symptoms**: No traffic passing through VLAN interfaces

**Solutions**:
1. Ensure `vlan-filtering=yes` on bridge
2. Check `bridge vlan` table is configured
3. Verify `pvid` settings on bridge ports
4. Check `frame-types` settings

### Issue: Firewall Performance Degradation

**Symptoms**: Throughput lower after migration

**Solutions**:
1. Enable FastTrack in `/ip firewall filter`
2. Use RAW table for early packet filtering
3. Review connection tracking settings
4. Check hardware offloading status

## Additional Resources

- [MikroTik RouterOS 7 Documentation](https://help.mikrotik.com/docs/display/ROS/RouterOS)
- [BGP in RouterOS 7](https://help.mikrotik.com/docs/display/ROS/BGP)
- [Firewall Configuration](https://help.mikrotik.com/docs/display/ROS/Firewall)
- [Bridge VLAN Filtering](https://help.mikrotik.com/docs/display/ROS/Bridging+and+Switching)

## Getting Help

If you encounter issues:

1. Check the [GitHub Issues](https://github.com/ddelnano/terraform-provider-mikrotik/issues)
2. Join the [Discord Community](https://discord.gg/ZpNq8ez)
3. Review MikroTik's official documentation
4. Post detailed logs and configuration when asking for help

## Version Support Policy

- **RouterOS 6.x**: Legacy support (maintenance mode)
- **RouterOS 7.x**: Full support (actively developed)
- **Migration Period**: Both versions supported during transition

Plan your migration within 6-12 months for best experience and support.
