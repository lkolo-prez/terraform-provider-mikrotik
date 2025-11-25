# Migration Guide: BGP v6 → BGP v7
**From RouterOS 6.x to RouterOS 7.20+**

---

## Table of Contents

1. [Overview](#overview)
2. [Why Migrate?](#why-migrate)
3. [Breaking Changes](#breaking-changes)
4. [Step-by-Step Migration](#step-by-step-migration)
5. [Resource Mapping](#resource-mapping)
6. [Common Patterns](#common-patterns)
7. [Troubleshooting](#troubleshooting)

---

## Overview

RouterOS 7 introduced a completely redesigned BGP implementation with significant architectural changes. The old `mikrotik_bgp_instance` and `mikrotik_bgp_peer` resources are **deprecated** and will be removed in a future version.

### Timeline
- **RouterOS 6.x**: Use `mikrotik_bgp_instance` + `mikrotik_bgp_peer`
- **RouterOS 7.0-7.19**: Migration period (both versions work)
- **RouterOS 7.20+**: Use `mikrotik_bgp_instance_v7` + `mikrotik_bgp_connection` (RECOMMENDED)
- **Future**: Old resources will be removed

---

## Why Migrate?

### New Features in BGP v7

1. **Templates** - Reusable configuration across multiple connections
2. **VRF Support** - Native VRF and routing table integration
3. **MPLS/VPN** - Built-in Layer 3 VPN support with route distinguishers
4. **Connection Modes** - Explicit `listen` and `connect` modes
5. **BFD Integration** - Fast failure detection
6. **Enhanced Filtering** - Better input/output filter control
7. **Improved Session Monitoring** - Rich session state information
8. **Address Family Support** - Multi-protocol BGP (IPv4, IPv6, VPNv4, VPNv6)

### Old BGP v6 Limitations

- ❌ No template support
- ❌ Limited VRF integration
- ❌ No native MPLS/VPN
- ❌ No BFD support
- ❌ Limited session monitoring
- ❌ Single address family focus

---

## Breaking Changes

### 1. Resource Names

| BGP v6 (Old) | BGP v7 (New) |
|--------------|--------------|
| `mikrotik_bgp_instance` | `mikrotik_bgp_instance_v7` |
| `mikrotik_bgp_peer` | `mikrotik_bgp_connection` |
| N/A | `mikrotik_bgp_template` (NEW) |
| N/A | `mikrotik_bgp_session` (data source, NEW) |

### 2. Connection Attributes

#### BGP v6 (mikrotik_bgp_peer)
```hcl
resource "mikrotik_bgp_peer" "peer" {
  instance        = "default"
  remote_address  = "10.0.0.2"
  remote_as       = 65001
  # Automatic bidirectional connection
}
```

#### BGP v7 (mikrotik_bgp_connection)
```hcl
resource "mikrotik_bgp_connection" "peer" {
  name           = "peer-1"          # NEW: Explicit name required
  instance       = "default"
  remote_address = "10.0.0.2"
  remote_as      = 65001
  connect        = true              # NEW: Must specify direction
  # OR
  # listen = true                    # For incoming connections
}
```

### 3. Instance Configuration

#### BGP v6 (mikrotik_bgp_instance)
```hcl
resource "mikrotik_bgp_instance" "main" {
  name       = "default"
  as         = 65000
  router_id  = "10.0.0.1"
  # Limited routing table control
}
```

#### BGP v7 (mikrotik_bgp_instance_v7)
```hcl
resource "mikrotik_bgp_instance_v7" "main" {
  name           = "default"
  as             = 65000
  router_id      = "10.0.0.1"
  vrf            = "customer-vrf"    # NEW: VRF support
  routing_table  = "customer-table"  # NEW: Routing table control
  cluster_id     = "10.0.0.1"        # NEW: Route reflector support
}
```

### 4. Removed/Changed Attributes

**Removed in v7:**
- `multihop` moved from instance to connection
- `route_reflect` moved from instance to connection
- Peer-level timers unified in templates

**Added in v7:**
- `templates` - Apply reusable templates
- `use_bfd` - BFD integration
- `use_mpls` - MPLS support
- `route_distinguisher` - VPN support
- `address_families` - Multi-protocol support
- Connection modes (`listen`, `connect`)

---

## Step-by-Step Migration

### Step 1: Backup Current Configuration

```bash
# Export current BGP configuration
/routing/bgp/instance/export file=bgp-backup-v6
/routing/bgp/peer/export file=bgp-peer-backup-v6

# In Terraform
terraform state pull > terraform-state-backup.json
```

### Step 2: Upgrade RouterOS

```bash
# Check current version
/system/package/print

# Upgrade to RouterOS 7.20+
# Download from https://mikrotik.com/download
/system/package/update check-for-updates
/system/package/update install

# Reboot required
/system/reboot
```

### Step 3: Update Terraform Provider

```hcl
# In your terraform configuration
terraform {
  required_providers {
    mikrotik = {
      source  = "terraform-provider-mikrotik/mikrotik"
      version = ">= 1.0.0"  # Use latest version
    }
  }
}
```

```bash
terraform init -upgrade
```

### Step 4: Create New Resources (Parallel)

**DO NOT** delete old resources yet. Create new ones alongside:

```hcl
# OLD (keep for now)
resource "mikrotik_bgp_instance" "old" {
  name      = "default"
  as        = 65000
  router_id = "10.0.0.1"
}

resource "mikrotik_bgp_peer" "old_peer" {
  instance       = "default"
  remote_address = "10.0.0.2"
  remote_as      = 65001
}

# NEW (create alongside)
resource "mikrotik_bgp_instance_v7" "new" {
  name      = "main-v7"  # Different name to avoid conflict
  as        = 65000
  router_id = "10.0.0.1"
}

resource "mikrotik_bgp_connection" "new_peer" {
  name           = "peer-1"
  instance       = mikrotik_bgp_instance_v7.new.name
  remote_address = "10.0.0.2"
  remote_as      = 65001
  connect        = true
}
```

### Step 5: Test New Configuration

```bash
# Apply new resources
terraform apply

# Verify BGP sessions
/routing/bgp/session/print
/routing/bgp/connection/print

# Check received routes
/ip/route/print where bgp

# Monitor both sessions
/routing/bgp/session/monitor [find] once
```

### Step 6: Switch Traffic (Careful!)

```bash
# On RouterOS, prefer new session
/routing/bgp/connection/set [find name=peer-1] disabled=no

# Disable old peer (graceful shutdown)
/routing/bgp/peer/set [find remote-address=10.0.0.2] disabled=yes

# Wait for convergence (check route tables)
/ip/route/print where bgp
```

### Step 7: Remove Old Resources

After confirming new BGP is working:

```hcl
# Remove old resources from Terraform
# Comment out or delete:
# resource "mikrotik_bgp_instance" "old" { ... }
# resource "mikrotik_bgp_peer" "old_peer" { ... }

# Update dependencies
# resource "mikrotik_bgp_connection" "new_peer" {
#   # Remove lifecycle prevent_destroy if set
# }
```

```bash
# Remove from state
terraform state rm mikrotik_bgp_peer.old_peer
terraform state rm mikrotik_bgp_instance.old

# Clean up on RouterOS (optional)
/routing/bgp/peer/remove [find instance=default]
/routing/bgp/instance/remove [find name=default]
```

---

## Resource Mapping

### Simple eBGP Peering

#### BGP v6 (Old)
```hcl
resource "mikrotik_bgp_instance" "main" {
  name      = "default"
  as        = 65000
  router_id = "10.0.0.1"
}

resource "mikrotik_bgp_peer" "peer1" {
  instance       = "default"
  remote_address = "10.0.0.2"
  remote_as      = 65001
  
  depends_on = [mikrotik_bgp_instance.main]
}
```

#### BGP v7 (New)
```hcl
resource "mikrotik_bgp_instance_v7" "main" {
  name      = "default"
  as        = 65000
  router_id = "10.0.0.1"
}

resource "mikrotik_bgp_connection" "peer1" {
  name           = "to-peer1"
  instance       = mikrotik_bgp_instance_v7.main.name
  remote_address = "10.0.0.2"
  remote_as      = 65001
  connect        = true
  
  address_families = "ip"
}
```

---

### iBGP with Route Reflection

#### BGP v6 (Old)
```hcl
resource "mikrotik_bgp_instance" "main" {
  name               = "default"
  as                 = 65000
  router_id          = "10.0.0.1"
  route_reflect      = true
  client_to_client_reflection = true
}

resource "mikrotik_bgp_peer" "client1" {
  instance       = "default"
  remote_address = "10.0.0.2"
  remote_as      = 65000
  route_reflect  = true
}

resource "mikrotik_bgp_peer" "client2" {
  instance       = "default"
  remote_address = "10.0.0.3"
  remote_as      = 65000
  route_reflect  = true
}
```

#### BGP v7 (New)
```hcl
resource "mikrotik_bgp_instance_v7" "main" {
  name      = "default"
  as        = 65000
  router_id = "10.0.0.1"
  
  client_to_client_reflection = true
  cluster_id                  = "10.0.0.1"
}

# Create reusable template for RR clients
resource "mikrotik_bgp_template" "rr_clients" {
  name = "route-reflector-clients"
  as   = 65000
  
  route_reflect = true
  
  address_families = "ip,ipv6"
  capabilities     = "mp,refresh,as4"
}

resource "mikrotik_bgp_connection" "client1" {
  name           = "rr-client-1"
  instance       = mikrotik_bgp_instance_v7.main.name
  remote_address = "10.0.0.2"
  remote_as      = 65000
  
  templates  = [mikrotik_bgp_template.rr_clients.name]
  local_role = "route-reflector"
  multihop   = true
}

resource "mikrotik_bgp_connection" "client2" {
  name           = "rr-client-2"
  instance       = mikrotik_bgp_instance_v7.main.name
  remote_address = "10.0.0.3"
  remote_as      = 65000
  
  templates  = [mikrotik_bgp_template.rr_clients.name]
  local_role = "route-reflector"
  multihop   = true
}
```

**Benefits of v7 approach:**
- ✅ Template reduces duplication
- ✅ Explicit route reflector role
- ✅ Easier to add new clients
- ✅ Consistent configuration

---

### Multiple Peers (Scaling)

#### BGP v6 (Old) - Repetitive
```hcl
resource "mikrotik_bgp_peer" "peer1" {
  instance       = "default"
  remote_address = "10.0.0.2"
  remote_as      = 65001
  # ... 20 lines of config ...
}

resource "mikrotik_bgp_peer" "peer2" {
  instance       = "default"
  remote_address = "10.0.0.3"
  remote_as      = 65002
  # ... same 20 lines repeated ...
}

# Repeat for 10+ peers = 200+ lines
```

#### BGP v7 (New) - Template-based
```hcl
# Define template once
resource "mikrotik_bgp_template" "customers" {
  name = "customer-peers"
  as   = 65000
  
  # All common configuration here (20 lines)
  address_families = "ip"
  capabilities     = "mp,refresh,as4"
  hold_time        = "3m"
  keepalive_time   = "1m"
  # ... etc ...
}

# Apply to multiple peers (5 lines each)
resource "mikrotik_bgp_connection" "peer1" {
  name           = "customer-1"
  instance       = mikrotik_bgp_instance_v7.main.name
  remote_address = "10.0.0.2"
  remote_as      = 65001
  templates      = [mikrotik_bgp_template.customers.name]
}

resource "mikrotik_bgp_connection" "peer2" {
  name           = "customer-2"
  instance       = mikrotik_bgp_instance_v7.main.name
  remote_address = "10.0.0.3"
  remote_as      = 65002
  templates      = [mikrotik_bgp_template.customers.name]
}

# 10 peers = 70 lines vs 200+ lines
```

**Reduction:** ~65% less code for scaled deployments

---

## Common Patterns

### Pattern 1: Adding VRF Support

```hcl
# Create VRF-aware BGP instance
resource "mikrotik_bgp_instance_v7" "customer_a" {
  name          = "customer-a"
  as            = 65000
  router_id     = "10.0.0.1"
  
  # NEW in v7
  vrf           = "customer-a-vrf"
  routing_table = "customer-a"
}

resource "mikrotik_bgp_connection" "customer_a_site1" {
  name               = "cust-a-site1"
  instance           = mikrotik_bgp_instance_v7.customer_a.name
  remote_address     = "172.16.1.1"
  remote_as          = 65100
  
  # VRF and MPLS
  vrf                = "customer-a-vrf"
  route_distinguisher = "65000:1001"
  use_mpls           = true
  address_families   = "vpnv4"
}
```

### Pattern 2: Enabling BFD

```hcl
resource "mikrotik_bgp_connection" "critical_peer" {
  name           = "critical-link"
  instance       = mikrotik_bgp_instance_v7.main.name
  remote_address = "10.0.0.2"
  remote_as      = 65001
  
  # Enable BFD for sub-second failure detection
  use_bfd = true
  
  # Reduced timers for faster convergence
  hold_time      = "90s"
  keepalive_time = "30s"
}
```

### Pattern 3: Session Monitoring

```hcl
# Data source for monitoring (NEW in v7)
data "mikrotik_bgp_session" "all_sessions" {
  depends_on = [
    mikrotik_bgp_connection.peer1,
    mikrotik_bgp_connection.peer2
  ]
}

output "bgp_status" {
  value = {
    established = data.mikrotik_bgp_session.all_sessions.established
    state       = data.mikrotik_bgp_session.all_sessions.state
    uptime      = data.mikrotik_bgp_session.all_sessions.uptime
    prefix_count = data.mikrotik_bgp_session.all_sessions.prefix_count
  }
}
```

---

## Troubleshooting

### Issue 1: "Resource not found" after migration

**Problem:** Old BGP resources return errors after RouterOS upgrade.

**Solution:**
```bash
# Check RouterOS version
/system/package/print

# If 7.20+, old /routing/bgp paths changed
# Use new resources:
terraform state rm mikrotik_bgp_instance.old
terraform state rm mikrotik_bgp_peer.old_peer

# Import new resources
terraform import mikrotik_bgp_instance_v7.new /routing/bgp/instance/main-v7
terraform import mikrotik_bgp_connection.new_peer /routing/bgp/connection/peer-1
```

### Issue 2: Connection stays in "Connect" state

**Problem:** BGP connection not establishing in v7.

**Causes & Solutions:**

1. **Missing connection direction:**
```hcl
# WRONG
resource "mikrotik_bgp_connection" "peer" {
  # No listen or connect specified
}

# CORRECT
resource "mikrotik_bgp_connection" "peer" {
  connect = true  # For outgoing
  # OR
  # listen = true  # For incoming
}
```

2. **Multihop not set:**
```hcl
resource "mikrotik_bgp_connection" "peer" {
  remote_address = "10.0.0.2"  # Not directly connected
  multihop       = true         # Required!
  ttl            = 255
}
```

### Issue 3: Templates not applying

**Problem:** BGP template changes don't affect connections.

**Solution:**
```hcl
# Ensure template is created first
resource "mikrotik_bgp_template" "example" {
  name = "my-template"
  # ... config ...
}

resource "mikrotik_bgp_connection" "peer" {
  templates = [mikrotik_bgp_template.example.name]  # Reference by name
  
  # Explicit dependency
  depends_on = [mikrotik_bgp_template.example]
}

# Recreate connection to apply template changes
terraform taint mikrotik_bgp_connection.peer
terraform apply
```

### Issue 4: Deprecation warnings

**Problem:** Terraform shows deprecation warnings.

```
│ Warning: Deprecated Resource
│ 
│   with mikrotik_bgp_instance.old,
│   on main.tf line 10, in resource "mikrotik_bgp_instance" "old":
│   10: resource "mikrotik_bgp_instance" "old" {
│ 
│ Use mikrotik_bgp_instance_v7 for RouterOS 7.20+. This resource will be
│ removed in a future version.
```

**Solution:** Migrate to v7 resources as shown in this guide.

### Issue 5: "Name" attribute required

**Problem:** BGP v6 didn't require connection names, v7 does.

```hcl
# BGP v6 (automatic naming)
resource "mikrotik_bgp_peer" "peer" {
  instance       = "default"
  remote_address = "10.0.0.2"  # Used as identifier
}

# BGP v7 (explicit naming)
resource "mikrotik_bgp_connection" "peer" {
  name           = "to-peer-10.0.0.2"  # Required!
  instance       = "default"
  remote_address = "10.0.0.2"
}
```

**Best practice:** Use descriptive names like `to-isp-as65001`, `customer-site-a`, `ibgp-core2`.

---

## Migration Checklist

- [ ] Backup current Terraform state
- [ ] Backup RouterOS BGP configuration
- [ ] Upgrade RouterOS to 7.20+
- [ ] Update Terraform provider to latest version
- [ ] Create new BGP v7 resources (parallel to old)
- [ ] Test new BGP sessions establish correctly
- [ ] Verify routes are received/advertised
- [ ] Switch traffic to new BGP sessions
- [ ] Monitor for 24-48 hours
- [ ] Remove old BGP v6 resources from Terraform
- [ ] Clean up old BGP configuration on RouterOS
- [ ] Update documentation and runbooks
- [ ] Celebrate! 🎉

---

## Additional Resources

- [BGP v7 Examples](./README.md)
- [MikroTik BGP Documentation](https://help.mikrotik.com/docs/display/ROS/BGP)
- [Provider GitHub Issues](https://github.com/lkolo-prez/terraform-provider-mikrotik/issues)
- [RouterOS 7 Coverage Analysis](../../ROUTEROS7_GAP_ANALYSIS.md)

---

**Need Help?** Open an issue on GitHub: https://github.com/lkolo-prez/terraform-provider-mikrotik/issues

**Last Updated:** December 2024  
**Applies to:** RouterOS 7.20+, Provider v1.0+
