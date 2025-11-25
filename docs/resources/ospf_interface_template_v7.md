# mikrotik_ospf_interface_template_v7 (Resource)

Manages OSPF interface template configuration on RouterOS v7.

Interface templates define how OSPF operates on matched interfaces or networks. Templates control network type, cost, timers, authentication, and other interface-specific OSPF parameters.

## Features

- **Flexible Matching**: Match by network prefix OR interface name
- **Network Types**: Broadcast, Point-to-Point, PTMP, NBMA, Virtual Link
- **Authentication**: MD5, SHA1, SHA256, SHA384, SHA512
- **Tunable Timers**: Hello, dead, retransmit intervals
- **Passive Mode**: Advertise without forming adjacencies
- **DR/BDR Control**: Priority-based election

## Example Usage

### Basic Interface Template (Network Match)

```terraform
resource "mikrotik_ospf_instance_v7" "main" {
  name      = "main_ospf"
  version   = "2"
  router_id = "1.1.1.1"
}

resource "mikrotik_ospf_area_v7" "backbone" {
  name     = "backbone"
  instance = mikrotik_ospf_instance_v7.main.name
  area_id  = "0.0.0.0"
}

# Enable OSPF on networks
resource "mikrotik_ospf_interface_template_v7" "lan" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = [
    "192.168.1.0/24",
    "10.0.0.0/8"
  ]
  cost     = 10
  comment  = "LAN interfaces"
}
```

### Interface Name Matching

```terraform
# Match specific interfaces by name
resource "mikrotik_ospf_interface_template_v7" "wan" {
  area       = mikrotik_ospf_area_v7.backbone.name
  interfaces = ["ether1", "ether2"]
  type       = "ptp"  # Point-to-point
  cost       = 100
  comment    = "WAN uplinks"
}
```

### Point-to-Point Links

```terraform
# Direct router-to-router links (no DR/BDR needed)
resource "mikrotik_ospf_interface_template_v7" "p2p" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["10.255.0.0/30", "10.255.0.4/30"]
  type     = "ptp"
  cost     = 10
  priority = 0  # Not used in ptp, but can set
  comment  = "Point-to-point links"
}
```

### OSPF Authentication (SHA256)

```terraform
resource "mikrotik_ospf_interface_template_v7" "secure" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["192.168.100.0/24"]
  auth     = "sha256"
  auth_key = "MySecureOSPFPassword2024"
  auth_id  = 1
  comment  = "Authenticated OSPF"
}
```

### MD5 Authentication with Key Rollover

```terraform
# Key 1 (current)
resource "mikrotik_ospf_interface_template_v7" "auth_key1" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["10.20.0.0/24"]
  auth     = "md5"
  auth_key = "current_key_2024"
  auth_id  = 1
}

# Key 2 (new, for rollover)
resource "mikrotik_ospf_interface_template_v7" "auth_key2" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["10.20.0.0/24"]
  auth     = "md5"
  auth_key = "new_key_2024"
  auth_id  = 2
}
```

### Passive Interface (Advertise Only)

```terraform
# Advertise network but don't form adjacencies
resource "mikrotik_ospf_interface_template_v7" "passive_lan" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["192.168.200.0/24"]
  passive  = true  # No Hello packets sent
  comment  = "Internal LAN - no OSPF neighbors"
}
```

### Custom Timers (Fast Convergence)

```terraform
resource "mikrotik_ospf_interface_template_v7" "fast" {
  area               = mikrotik_ospf_area_v7.backbone.name
  networks           = ["10.10.10.0/24"]
  hello_interval     = "1s"   # Faster hellos
  dead_interval      = "3s"   # Faster failure detection
  retransmit_interval = "2s"
  comment            = "Fast convergence link"
}
```

### DR/BDR Priority Control

```terraform
# High priority - prefer as DR
resource "mikrotik_ospf_interface_template_v7" "dr_preferred" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["192.168.50.0/24"]
  priority = 255  # Highest priority, likely DR
  comment  = "Prefer this router as DR"
}

# Never become DR/BDR
resource "mikrotik_ospf_interface_template_v7" "no_dr" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["192.168.60.0/24"]
  priority = 0  # Never elected
  comment  = "Low-power router, never DR"
}
```

### NBMA Network Configuration

```terraform
resource "mikrotik_ospf_interface_template_v7" "nbma" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["10.50.0.0/24"]
  type     = "nbma"
  priority = 128
  comment  = "Frame Relay or similar NBMA"
}

# Static neighbors required for NBMA (separate resource in future)
```

### Virtual Link (Connect Remote Area to Backbone)

```terraform
resource "mikrotik_ospf_area_v7" "transit" {
  name     = "transit_area"
  instance = mikrotik_ospf_instance_v7.main.name
  area_id  = "5.5.5.5"
  type     = "default"
}

# Virtual link through transit area
resource "mikrotik_ospf_interface_template_v7" "vlink" {
  area               = mikrotik_ospf_area_v7.backbone.name
  type               = "virtual-link"
  vlink_transit_area = "transit_area"
  vlink_neighbor_id  = "10.0.0.2"  # Remote ABR router-id
  auth               = "sha256"     # Authentication recommended
  auth_key           = "vlink_secret"
  auth_id            = 1
  comment            = "Virtual link to remote area"
}
```

### Complete Enterprise Example

```terraform
resource "mikrotik_ospf_instance_v7" "enterprise" {
  name      = "enterprise_ospf"
  version   = "2"
  router_id = "172.16.0.1"
}

resource "mikrotik_ospf_area_v7" "backbone" {
  name     = "backbone"
  instance = mikrotik_ospf_instance_v7.enterprise.name
  area_id  = "0.0.0.0"
}

# Internal LAN - Passive
resource "mikrotik_ospf_interface_template_v7" "internal_lan" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["10.0.0.0/8"]
  passive  = true
  cost     = 10
  comment  = "Internal subnets"
}

# Core routers - Point-to-point, SHA256 auth
resource "mikrotik_ospf_interface_template_v7" "core_links" {
  area               = mikrotik_ospf_area_v7.backbone.name
  networks           = ["172.16.0.0/16"]
  type               = "ptp"
  cost               = 10
  auth               = "sha256"
  auth_key           = var.ospf_auth_key
  auth_id            = 1
  hello_interval     = "5s"
  dead_interval      = "20s"
  comment            = "Core router links"
}

# Edge routers - Broadcast, lower priority
resource "mikrotik_ospf_interface_template_v7" "edge_routers" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["192.168.0.0/16"]
  type     = "broadcast"
  cost     = 100
  priority = 50  # Lower than core routers
  auth     = "sha256"
  auth_key = var.ospf_auth_key
  auth_id  = 1
  comment  = "Edge router segments"
}
```

## Argument Reference

### Required

- `area` (String) - Name of the OSPF area this template belongs to.

### Interface/Network Matching (At least one required)

- `networks` (List of String) - List of network prefixes in CIDR format to match. Example: `["192.168.1.0/24", "10.0.0.0/8"]`.
- `interfaces` (List of String) - List of interface names to match directly. Example: `["ether1", "bridge1"]`.

### Optional - General

- `type` (String) - Network type. Valid values:
  - `"broadcast"` - Ethernet, DR/BDR election **(default)**
  - `"ptp"` - Point-to-point, no DR/BDR
  - `"ptmp"` - Point-to-multipoint, no DR/BDR
  - `"nbma"` - Non-broadcast multi-access, requires static neighbors
  - `"ptmp-broadcast"` - PTMP with automatic neighbor discovery
  - `"virtual-link"` - Virtual link through transit area
- `disabled` (Boolean) - Disable this template. Default: `false`.
- `comment` (String) - Comment for this template.

### Optional - Cost and Priority

- `cost` (Number) - Interface cost (metric). Lower is better. Range: 1-65535. Default: `10`.
- `priority` (Number) - DR/BDR election priority (0-255). Higher wins. `0` = never become DR/BDR. Default: `128`.
- `passive` (Boolean) - Passive mode: advertise network but don't send/receive OSPF packets. Default: `false`.

### Optional - Authentication

- `auth` (String) - Authentication type. Valid values:
  - `"none"` - No authentication **(default)**
  - `"simple"` - Plaintext password (insecure, not recommended)
  - `"md5"` - MD5 cryptographic authentication
  - `"sha1"` - SHA1 authentication
  - `"sha256"` - SHA256 authentication (recommended)
  - `"sha384"` - SHA384 authentication
  - `"sha512"` - SHA512 authentication
- `auth_key` (String, Sensitive) - Authentication password/key. Required when `auth` is not `"none"`.
- `auth_id` (Number) - Key ID for MD5/SHA authentication (1-255). Used for key rollover.

### Optional - Timers

- `hello_interval` (String) - Hello packet interval. RouterOS time format (e.g., `"10s"`, `"30s"`). Default: `"10s"`.
- `dead_interval` (String) - Dead interval (neighbor down after missing hellos). Should be 4× `hello_interval`. Default: `"40s"`.
- `retransmit_interval` (String) - LSA retransmission interval. Default: `"5s"`.
- `transmit_delay` (String) - LSA transmission delay (age increment). Default: `"1s"`.
- `wait_time` (String) - Time to wait before DR/BDR election on startup. Default: same as `dead_interval`.

### Optional - Virtual Link

- `vlink_transit_area` (String) - Transit area name for virtual link. Required when `type="virtual-link"`.
- `vlink_neighbor_id` (String) - Remote router ID (IPv4 format). Required when `type="virtual-link"`.

## Attribute Reference

- `id` (String) - Unique identifier (.id).
- `dynamic` (Boolean) - Whether dynamically created (computed).
- `invalid` (Boolean) - Whether configuration is invalid (computed).

## Import

OSPF interface templates can be imported by ID (they don't have unique names):

```shell
terraform import mikrotik_ospf_interface_template_v7.example *1A
```

## Notes

### Network vs. Interface Matching

**Network Matching (`networks`):**
- Matches interfaces with IPs in specified subnets
- More flexible (survives interface renames)
- Recommended for most scenarios

**Interface Matching (`interfaces`):**
- Matches by exact interface name
- Useful for unnumbered interfaces or specific control
- Brittle if interfaces renamed

### Network Type Selection

| Type | Use Case | DR/BDR | Neighbor Discovery | Example |
|------|----------|--------|-------------------|---------|
| **broadcast** | Ethernet, WiFi | Yes | Multicast | LAN segments |
| **ptp** | Direct links | No | Multicast | Router-to-router |
| **ptmp** | Hub-spoke | No | Manual config | DMVPN |
| **ptmp-broadcast** | Hub-spoke | No | Multicast | RouterOS v6 compat |
| **nbma** | Frame Relay | Yes | Manual config | Legacy WAN |
| **virtual-link** | Backbone extension | N/A | Configured | Logical link |

### Authentication Best Practices

**Recommendations:**
1. **Use SHA256 or better** for production
2. **Avoid `simple` auth** (plaintext, easily sniffed)
3. **Use key rollover** with different `auth_id` for seamless updates
4. **Secure auth_key storage** (use Terraform variables, Vault, etc.)

**Key Rollover Process:**
1. Add new key with different `auth_id` on all routers
2. Wait for propagation (few minutes)
3. Remove old key template
4. All routers accept both keys during transition

### Timer Tuning

**Default Timers (hello=10s, dead=40s):**
- Good balance for most networks
- Convergence time: 40-50 seconds

**Fast Convergence (hello=1s, dead=3s):**
- Detect failures in 3 seconds
- **Warning**: Increases CPU usage and Hello traffic
- Use only on stable, low-latency links

**Slow Convergence (hello=30s, dead=120s):**
- Reduces CPU and bandwidth
- Use for low-speed WAN links or unstable connections

**Rules:**
- `dead_interval` should be **4× `hello_interval`**
- Too fast timers may cause flapping on high-latency links

### DR/BDR Election

**Priority Rules:**
- **Highest priority wins** (priority 128 > priority 100)
- **Tie-breaker**: Highest router ID
- **Priority 0**: Never becomes DR or BDR (DROther only)

**When to Adjust Priority:**
- **Core routers**: Higher priority (200-255) to ensure stable DR
- **Edge/low-power routers**: Lower priority (1-50) or 0
- **Equal routers**: Leave at default (128)

**DR/BDR Functions:**
- **DR**: Floods LSAs to all routers (224.0.0.5)
- **BDR**: Standby, takes over if DR fails
- **DROther**: Only exchanges with DR/BDR

### Cost Calculation

**Cisco Formula (optional):**
```
cost = 100,000,000 / bandwidth_bps
```

**RouterOS Default:** `cost = 10` for all interfaces

**Common Cost Values:**
- 1 = 100 Gbps+ (very fast)
- 10 = 10 Gbps (default, good for most)
- 100 = 1 Gbps
- 1000 = 100 Mbps
- 10000 = 10 Mbps

**Use Cost to Prefer Paths:**
- Lower cost = preferred path
- Equal cost = ECMP (load balancing)

### Passive Interfaces

**When to Use:**
- **User LANs**: Advertise subnets, no neighbors expected
- **DMZ segments**: Security isolation
- **Loopback interfaces**: Advertise addresses only

**Benefits:**
- Reduced CPU (no Hello processing)
- Security (no unexpected adjacencies)
- Faster convergence (less complexity)

### Performance Considerations

**CPU Impact:**
- Hello packets: Minimal (every 10s)
- SPF calculation: High (on topology change)
- Authentication: SHA512 > SHA256 > MD5 > none

**Memory Impact:**
- Each template: ~1-2 KB
- Per-interface state: ~10-20 KB

**Recommendations:**
- **Limit templates**: Consolidate similar interfaces
- **Use passive mode**: Where possible
- **Tune timers**: Balance convergence vs. overhead

## See Also

- [mikrotik_ospf_instance_v7](ospf_instance_v7.md) - Configure OSPF instances
- [mikrotik_ospf_area_v7](ospf_area_v7.md) - Configure OSPF areas
- [MikroTik OSPF Documentation](https://help.mikrotik.com/docs/display/ROS/OSPF)
