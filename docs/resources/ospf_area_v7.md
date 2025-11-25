# mikrotik_ospf_area_v7 (Resource)

Manages OSPF area configuration on RouterOS v7.

OSPF areas organize routing domains hierarchically, reducing LSA flooding and SPF calculation overhead. All areas must connect to the backbone area (0.0.0.0).

## Features

- **Area Types**: Standard, Stub, Totally Stubby, NSSA (Not-So-Stubby Area)
- **LSA Control**: Filter Type-5 external LSAs in stub/NSSA areas
- **Summarization**: Reduce routing information between areas
- **NSSA Translation**: Convert Type-7 to Type-5 LSAs at ABR

## Example Usage

### Backbone Area (Required)

```terraform
resource "mikrotik_ospf_instance_v7" "main" {
  name      = "main_ospf"
  version   = "2"
  router_id = "1.1.1.1"
}

# Backbone area - MUST exist and have area_id 0.0.0.0
resource "mikrotik_ospf_area_v7" "backbone" {
  name     = "backbone"
  instance = mikrotik_ospf_instance_v7.main.name
  area_id  = "0.0.0.0"  # Backbone always uses 0.0.0.0
  comment  = "OSPF backbone area"
}
```

### Standard Area (Multi-Area Setup)

```terraform
# Area 1 - Standard area for branch office
resource "mikrotik_ospf_area_v7" "area1" {
  name     = "branch_office"
  instance = mikrotik_ospf_instance_v7.main.name
  area_id  = "1.1.1.1"
  type     = "default"  # Allows all LSA types
  comment  = "Branch office area"
}

# Area 2 - Another standard area
resource "mikrotik_ospf_area_v7" "area2" {
  name     = "datacenter"
  instance = mikrotik_ospf_instance_v7.main.name
  area_id  = "2.2.2.2"
  type     = "default"
  comment  = "Datacenter area"
}
```

### Stub Area

```terraform
# Stub area - Blocks Type-5 external LSAs
resource "mikrotik_ospf_area_v7" "stub_area" {
  name         = "remote_site"
  instance     = mikrotik_ospf_instance_v7.main.name
  area_id      = "10.0.0.1"
  type         = "stub"
  default_cost = 100  # Cost of injected default route
  comment      = "Remote site with single exit point"
}

# Internal router in stub area
resource "mikrotik_ospf_interface_template_v7" "stub_lan" {
  area     = mikrotik_ospf_area_v7.stub_area.name
  networks = ["10.10.0.0/16"]
}
```

### Totally Stubby Area

```terraform
# Totally stubby - Blocks Type-3 summary AND Type-5 external LSAs
resource "mikrotik_ospf_area_v7" "totally_stubby" {
  name         = "small_branch"
  instance     = mikrotik_ospf_instance_v7.main.name
  area_id      = "20.0.0.1"
  type         = "stub"
  no_summaries = true   # Makes it totally stubby
  default_cost = 50
  comment      = "Minimal routing info - default route only"
}
```

### NSSA (Not-So-Stubby Area)

```terraform
# NSSA - Stub area that allows redistribution via Type-7 LSAs
resource "mikrotik_ospf_area_v7" "nssa" {
  name             = "nssa_site"
  instance         = mikrotik_ospf_instance_v7.main.name
  area_id          = "30.0.0.1"
  type             = "nssa"
  nssa_translator  = "candidate"  # Participate in translator election
  nssa_propagation = true         # Propagate Type-7 LSAs
  comment          = "Site with external connection"
}

# ASBR in NSSA can redistribute external routes
resource "mikrotik_ospf_instance_v7" "with_nssa" {
  name               = "nssa_instance"
  version            = "2"
  router_id          = "30.0.0.2"
  redistribute_static = true  # Redistributed as Type-7 in NSSA
}
```

### Complete Multi-Area Enterprise Setup

```terraform
# Core router configuration
resource "mikrotik_ospf_instance_v7" "enterprise" {
  name                   = "enterprise_ospf"
  version                = "2"
  router_id              = "172.16.0.1"
  redistribute_connected = false
  redistribute_static    = false
}

# Backbone area
resource "mikrotik_ospf_area_v7" "backbone" {
  name     = "backbone"
  instance = mikrotik_ospf_instance_v7.enterprise.name
  area_id  = "0.0.0.0"
}

# Area 1 - Data center (standard)
resource "mikrotik_ospf_area_v7" "dc" {
  name     = "datacenter"
  instance = mikrotik_ospf_instance_v7.enterprise.name
  area_id  = "0.0.0.1"
  type     = "default"
}

# Area 2 - Branch offices (stub, single exit)
resource "mikrotik_ospf_area_v7" "branches" {
  name         = "branches"
  instance     = mikrotik_ospf_instance_v7.enterprise.name
  area_id      = "0.0.0.2"
  type         = "stub"
  default_cost = 100
}

# Area 3 - Remote sites with Internet (NSSA)
resource "mikrotik_ospf_area_v7" "remote_nssa" {
  name             = "remote_sites"
  instance         = mikrotik_ospf_instance_v7.enterprise.name
  area_id          = "0.0.0.3"
  type             = "nssa"
  nssa_translator  = "candidate"
  nssa_propagation = true
}

# Interface templates for each area
resource "mikrotik_ospf_interface_template_v7" "backbone_iface" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["172.16.0.0/24"]
  cost     = 10
}

resource "mikrotik_ospf_interface_template_v7" "dc_iface" {
  area     = mikrotik_ospf_area_v7.dc.name
  networks = ["10.0.0.0/16"]
  cost     = 10
}

resource "mikrotik_ospf_interface_template_v7" "branch_iface" {
  area     = mikrotik_ospf_area_v7.branches.name
  networks = ["192.168.0.0/16"]
  cost     = 100
}
```

## Argument Reference

### Required

- `name` (String) - Name of the OSPF area. Must be unique per instance. Forces new resource on change.
- `area_id` (String) - Area ID in IPv4 address format. Backbone must use `0.0.0.0`. Example: `1.1.1.1`, `0.0.0.10`.
- `instance` (String) - Name of the OSPF instance this area belongs to. Forces new resource on change.

### Optional

- `type` (String) - Area type. Valid values:
  - `"default"` - Standard area (allows all LSA types) **(default)**
  - `"stub"` - Stub area (blocks Type-5 external LSAs)
  - `"nssa"` - Not-So-Stubby Area (allows Type-7 LSAs for redistribution)
- `disabled` (Boolean) - Disable this area. Default: `false`.
- `comment` (String) - Comment for the area.

### Stub Area Options

- `default_cost` (Number) - Cost of the default route (0.0.0.0/0) injected into stub/totally stubby areas. Only used when `type` is `"stub"`. Default: 1.
- `no_summaries` (Boolean) - Create totally stubby area (blocks Type-3 summary LSAs in addition to Type-5). Only valid for stub areas. Default: `false`.

### NSSA Options

- `nssa_translator` (String) - ABR translator role for NSSA. Only used when `type` is `"nssa"`. Valid values:
  - `"candidate"` - Participate in translator election **(default)**
  - `"yes"` - Always translate Type-7 to Type-5
  - `"no"` - Never translate
- `nssa_propagation` (Boolean) - Propagate NSSA Type-7 LSAs. Default: `true`.

## Attribute Reference

- `id` (String) - Unique identifier of the area (.id).
- `dynamic` (Boolean) - Whether this area was dynamically created (computed).
- `invalid` (Boolean) - Whether configuration is invalid (computed).

## Import

OSPF areas can be imported by name:

```shell
terraform import mikrotik_ospf_area_v7.example backbone
```

## Notes

### Area Design Guidelines

**Area Size Limits:**
- Recommended: **≤50 routers per area** for low-end hardware
- Maximum: ~100-200 routers depending on topology and CPU

**When to Use Multiple Areas:**
- Network has **>50 routers**
- Hierarchical network design (HQ → regional → branch)
- Need to **hide topology** between segments
- Want to **reduce LSA flooding** and SPF calculations

### Backbone Area Requirements

- **Backbone area (0.0.0.0) MUST exist** before creating other areas
- **All areas must connect to backbone** (directly or via virtual link)
- ABR (Area Border Router) must have at least one interface in backbone

### Area Type Selection

| Area Type | Use Case | External Routes | Summary Routes | Example |
|-----------|----------|----------------|----------------|---------|
| **Standard** | Default, full routing | ✅ Type-5 LSAs | ✅ Type-3 LSAs | HQ, datacenter |
| **Stub** | Single exit point | ❌ Blocked | ✅ Type-3 LSAs | Branch office |
| **Totally Stubby** | Minimal routing | ❌ Blocked | ❌ Blocked | Small remote site |
| **NSSA** | Stub with redistribution | ✅ Type-7 → Type-5 | ✅ Type-3 LSAs | Branch with Internet |

### LSA Types by Area

**Standard Area:**
- Type-1 (Router LSA) ✅
- Type-2 (Network LSA) ✅
- Type-3 (Summary LSA) ✅
- Type-4 (ASBR Summary) ✅
- Type-5 (External LSA) ✅

**Stub Area:**
- Type-1, Type-2, Type-3 ✅
- Type-4, Type-5 ❌
- Default route injected by ABR

**Totally Stubby:**
- Type-1, Type-2 ✅
- Type-3, Type-4, Type-5 ❌
- Only default route from ABR

**NSSA:**
- Type-1, Type-2, Type-3 ✅
- Type-7 (NSSA External) ✅
- Type-5 ❌ (converted at ABR)

### Virtual Links

If area cannot physically connect to backbone, use virtual link:

```terraform
resource "mikrotik_ospf_area_v7" "transit" {
  name     = "transit_area"
  instance = mikrotik_ospf_instance_v7.main.name
  area_id  = "5.5.5.5"
  type     = "default"  # Transit area must be standard
}

resource "mikrotik_ospf_interface_template_v7" "vlink" {
  area                = mikrotik_ospf_area_v7.backbone.name
  type                = "virtual-link"
  vlink_transit_area  = mikrotik_ospf_area_v7.transit.name
  vlink_neighbor_id   = "10.0.0.2"  # Remote ABR router-id
}
```

### Performance Impact

**Stub/NSSA Benefits:**
- **Reduced memory** (smaller link-state database)
- **Faster convergence** (less SPF calculation)
- **Less bandwidth** (fewer LSA floods)

**Tradeoff:**
- **Suboptimal routing** to external destinations (always via ABR)

## See Also

- [mikrotik_ospf_instance_v7](ospf_instance_v7.md) - Configure OSPF instances
- [mikrotik_ospf_interface_template_v7](ospf_interface_template_v7.md) - Configure OSPF interfaces
- [MikroTik OSPF Documentation](https://help.mikrotik.com/docs/display/ROS/OSPF)
