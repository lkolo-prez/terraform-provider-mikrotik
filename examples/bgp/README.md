# BGP Examples for RouterOS 7.20+

This directory contains comprehensive Terraform examples demonstrating BGP v7 features in RouterOS 7.20+.

## Prerequisites

- RouterOS 7.20 or later
- Terraform 1.5+
- terraform-provider-mikrotik

## Examples Overview

### 01-ebgp-peering.tf
**External BGP (eBGP) peering between two autonomous systems**

- Two routers in different AS numbers (65001, 65002)
- TCP MD5 authentication
- Route redistribution (connected, static)
- Session monitoring with data source

**Use case:** ISP peering, customer connections, internet exchange

---

### 02-ibgp-full-mesh.tf
**Internal BGP (iBGP) full mesh topology**

- Three routers in same AS (65100)
- BGP template for shared configuration
- BFD integration for fast failure detection
- Full mesh connectivity (N*(N-1)/2 connections)

**Use case:** Small to medium enterprise networks (2-10 routers)

---

### 03-route-reflector.tf
**Route reflector (hub-and-spoke) topology**

- One route reflector + 3 clients
- Reduces connections from 6 (full mesh) to 3 (hub-spoke)
- Cluster ID configuration
- RR client template

**Use case:** Large enterprise networks (10+ routers), ISP backbones

**Topology:**
```
         ┌─────────────┐
         │   Route     │
         │  Reflector  │
         │  (10.0.2.1) │
         └──┬───┬───┬──┘
            │   │   │
    ┌───────┘   │   └───────┐
    │           │           │
┌───▼──┐    ┌───▼──┐    ┌───▼──┐
│Client│    │Client│    │Client│
│  1   │    │  2   │    │  3   │
└──────┘    └──────┘    └──────┘
```

---

### 04-vpn-mpls.tf
**BGP MPLS VPN (Layer 3 VPN)**

- PE (Provider Edge) routers with MPLS
- Multiple customer VRFs (A, B)
- Route distinguishers (RD)
- VPNv4/VPNv6 address families

**Use case:** Service provider L3VPN, multi-tenant networks

**Topology:**
```
Customer A Site 1                   Customer A Site 2
    (CE1)                               (CE3)
     │                                   │
     │         MPLS Backbone             │
     └──[PE1]──────────────[PE2]────────┘
             │              │
             │              │
            (CE2)          (CE4)
     Customer B Site 1   Customer B Site 2
```

---

### 05-communities-filtering.tf
**BGP Communities and Advanced Filtering**

- ISP-customer connection with community tagging
- Route classification (standard, backup, premium)
- Input/output community filtering
- Template-based customer configuration
- Rate limiting for customer prefixes

**Use case:** ISP traffic engineering, customer route management, policy-based routing

**Community values:**
- `65400:100` - Standard customer routes
- `65400:200` - Backup routes (lower priority)
- `65400:300` - Premium routes (higher priority)

---

### 06-graceful-restart.tf
**High Availability with Graceful Restart and BFD**

- Two core routers with HA configuration
- BGP Graceful Restart capability
- BFD for sub-second failure detection
- Preserves forwarding state during restarts
- Reduced timers for faster convergence

**Use case:** High-availability core networks, zero-downtime maintenance

**Features:**
- **Graceful Restart:** 120s restart time, 300s stale routes
- **BFD:** Sub-second link failure detection
- **Fast Convergence:** 90s hold time, 30s keepalive

---

## Quick Start

1. **Clone the repository:**
```bash
git clone https://github.com/lkolo-prez/terraform-provider-mikrotik.git
cd terraform-provider-mikrotik/examples/bgp
```

2. **Configure provider credentials:**

Edit the provider blocks in the example files:
```hcl
provider "mikrotik" {
  host     = "192.168.88.1"  # Your router IP
  username = "admin"          # Your username
  password = "admin"          # Your password
  tls      = true
  insecure = true
}
```

3. **Initialize Terraform:**
```bash
terraform init
```

4. **Plan the deployment:**
```bash
terraform plan
```

5. **Apply the configuration:**
```bash
terraform apply
```

---

## BGP Resources Reference

### mikrotik_bgp_instance_v7
**Main BGP instance configuration (RouterOS 7.20+)**

Key attributes:
- `name` - Instance name (required)
- `as` - AS number (required)
- `router_id` - BGP router ID
- `vrf` - VRF instance name
- `routing_table` - Routing table name
- `redistribute_*` - Route redistribution flags
- `client_to_client_reflection` - Enable route reflection

[Full documentation](../../docs/resources/bgp_instance_v7.md)

---

### mikrotik_bgp_connection
**BGP neighbor/connection configuration**

Key attributes:
- `name` - Connection name (required)
- `instance` - BGP instance (required)
- `remote_address` - Peer address (required)
- `remote_as` - Peer AS number (required)
- `listen` - Listen mode (incoming)
- `connect` - Connect mode (outgoing)
- `templates` - Apply BGP templates
- `address_families` - ip, ipv6, vpnv4, vpnv6
- `use_mpls` - Enable MPLS
- `use_bfd` - Enable BFD
- `vrf` - VRF name
- `route_distinguisher` - RD for VPN

[Full documentation](../../docs/resources/bgp_connection.md)

---

### mikrotik_bgp_template
**Reusable BGP configuration template**

Key attributes:
- `name` - Template name (required)
- `as` - AS number
- `capabilities` - BGP capabilities
- `route_reflect` - Enable route reflection
- `input_accept_nlri` - Accepted NLRI types
- `output_default_originate` - Default route origination
- `hold_time` - Hold timer
- `keepalive_time` - Keepalive timer

[Full documentation](../../docs/resources/bgp_template.md)

---

### mikrotik_bgp_session (data source)
**Monitor active BGP sessions**

Key attributes:
- `established` - Session established (bool)
- `state` - Session state (idle, connect, active, opensent, openconfirm, established)
- `uptime` - Session uptime
- `remote_address` - Peer address
- `remote_as` - Peer AS
- `prefix_count` - Advertised prefix count

[Full documentation](../../docs/data-sources/bgp_session.md)

---

## Advanced Features

### 1. BFD Integration
Enable fast failure detection (sub-second):
```hcl
resource "mikrotik_bgp_connection" "peer" {
  # ... basic config ...
  use_bfd = true
}
```

### 2. Multihop eBGP
For non-directly connected peers:
```hcl
resource "mikrotik_bgp_connection" "peer" {
  # ... basic config ...
  multihop = true
  ttl      = 255
}
```

### 3. Authentication
TCP MD5 authentication:
```hcl
resource "mikrotik_bgp_connection" "peer" {
  # ... basic config ...
  tcp_md5_key = "secret-key-123"
}
```

### 4. Route Filtering
Apply input/output filters:
```hcl
resource "mikrotik_bgp_connection" "peer" {
  # ... basic config ...
  input_filter  = "bgp-input-chain"
  output_filter = "bgp-output-chain"
}
```

### 5. VRF/MPLS
Layer 3 VPN configuration:
```hcl
resource "mikrotik_bgp_instance_v7" "customer_vrf" {
  name          = "customer-a"
  as            = 65000
  vrf           = "customer-a-vrf"
  routing_table = "customer-a"
}

resource "mikrotik_bgp_connection" "vpn_peer" {
  instance            = mikrotik_bgp_instance_v7.customer_vrf.name
  vrf                 = "customer-a-vrf"
  route_distinguisher = "65000:1001"
  use_mpls            = true
  address_families    = "vpnv4"
}
```

---

## Testing

### Verify BGP Session
```bash
# On RouterOS
/routing/bgp/session/print

# With Terraform
terraform console
> data.mikrotik_bgp_session.example
```

### Check Received Routes
```bash
/ip/route/print where bgp
```

### Monitor Session State
```bash
/routing/bgp/session/monitor [find] once
```

---

## Troubleshooting

### Session stuck in "Connect" state
- Check IP connectivity: `ping <peer-ip>`
- Verify firewall rules: `/ip/firewall/filter/print`
- Check TCP port 179: `telnet <peer-ip> 179`

### Session stuck in "Active" state
- Check `listen` vs `connect` mode
- Verify both sides have correct configuration
- Check `multihop` setting for non-adjacent peers

### No routes received
- Verify address families match on both sides
- Check input filters: `input-filter` and `input-accept-nlri`
- Verify peer is advertising routes

### High CPU usage
- Enable route caching (automatic in provider)
- Use BGP templates for multiple peers
- Consider route reflection for large topologies

---

## Performance Optimizations

The provider includes automatic performance optimizations:

### 1. Batch Operations
Bulk operations use `client/bgp_batch.go`:
```go
// Automatically used by provider
BatchAddConnections([]*BgpConnection)
BatchUpdateConnections([]*BgpConnection)
```

### 2. Caching
In-memory cache with `sync.RWMutex`:
- ~100x faster reads from cache
- ~90% reduction in API calls
- Thread-safe concurrent access

### 3. Bulk Fetch
Single API call for all resources:
```go
ListBgpInstancesV7()
ListBgpConnections()
ListBgpTemplates()
ListBgpSessions()
```

---

## Migration from BGP v6

### Old (RouterOS 6.x)
```hcl
resource "mikrotik_bgp_instance" "old" {
  name       = "default"
  as         = 65000
  router_id  = "10.0.0.1"
}

resource "mikrotik_bgp_peer" "old" {
  instance        = "default"
  remote_address  = "10.0.0.2"
  remote_as       = 65001
}
```

### New (RouterOS 7.20+)
```hcl
resource "mikrotik_bgp_instance_v7" "new" {
  name      = "default"
  as        = 65000
  router_id = "10.0.0.1"
}

resource "mikrotik_bgp_connection" "new" {
  name           = "peer-1"
  instance       = mikrotik_bgp_instance_v7.new.name
  remote_address = "10.0.0.2"
  remote_as      = 65001
  connect        = true
}
```

**Key differences:**
- `mikrotik_bgp_peer` → `mikrotik_bgp_connection`
- Connection requires explicit `name` attribute
- Must specify `connect = true` or `listen = true`
- More granular control over address families
- Native VRF support

**📖 Complete Migration Guide:** See [MIGRATION.md](./MIGRATION.md) for step-by-step migration instructions, troubleshooting, and common patterns.

---

## Contributing

Found an issue or have an improvement? Please open an issue or pull request:

https://github.com/lkolo-prez/terraform-provider-mikrotik

---

## Resources

- [MikroTik BGP Documentation](https://help.mikrotik.com/docs/display/ROS/BGP)
- [RouterOS 7 Release Notes](https://mikrotik.com/download/changelogs)
- [Provider GitHub Repository](https://github.com/lkolo-prez/terraform-provider-mikrotik)
- [RouterOS 7 Coverage Analysis](../../ROUTEROS7_GAP_ANALYSIS.md)

---

**Last Updated:** December 2024  
**RouterOS Version:** 7.20+  
**Provider Version:** Latest
