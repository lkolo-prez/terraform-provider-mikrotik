# BGP v7 Examples

Production-ready Terraform configurations for RouterOS 7.20+ BGP.

## Examples

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

```bash
cd examples/bgp
terraform init
terraform plan
terraform apply
```

## Resources

- `mikrotik_bgp_instance_v7` - BGP instance (AS, router-id, VRF)
- `mikrotik_bgp_connection` - BGP neighbor/peer
- `mikrotik_bgp_template` - Reusable configuration
- `mikrotik_bgp_session` - Session monitoring (data source)

## Migration from v6

Old (RouterOS 6):
```hcl
resource "mikrotik_bgp_instance" "old" {
  name = "default"
  as   = 65000
}
resource "mikrotik_bgp_peer" "old" {
  instance       = "default"
  remote_address = "10.0.0.2"
  remote_as      = 65001
}
```

New (RouterOS 7.20+):
```hcl
resource "mikrotik_bgp_instance_v7" "new" {
  name = "default"
  as   = 65000
}
resource "mikrotik_bgp_connection" "new" {
  name           = "peer-1"
  instance       = "default"
  remote_address = "10.0.0.2"
  remote_as      = 65001
  connect        = true
}
```

See [MIGRATION.md](./MIGRATION.md) for details.

---

**Documentation**: [MikroTik BGP](https://help.mikrotik.com/docs/display/ROS/BGP) | [Provider GitHub](https://github.com/lkolo-prez/terraform-provider-mikrotik)

