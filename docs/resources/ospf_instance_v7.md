# mikrotik_ospf_instance_v7 (Resource)

Manages OSPF v2/v3 instance configuration on RouterOS v7.

OSPF (Open Shortest Path First) is a link-state routing protocol for IP networks. RouterOS v7 provides unified OSPF implementation supporting both OSPFv2 (IPv4) and OSPFv3 (IPv6) through the same configuration interface.

## Features

- **Unified Configuration**: Single interface for OSPFv2 and OSPFv3
- **VRF Support**: Isolate OSPF routing in separate VRF instances
- **Route Redistribution**: Inject routes from other protocols (BGP, RIP, static, connected)
- **Default Route Control**: Fine-grained control over 0.0.0.0/0 origination
- **Routing Filters**: Apply filter chains to incoming/outgoing routes
- **Multi-Instance**: Run multiple OSPF instances with different domain IDs (RFC 6549)

## Example Usage

### Basic OSPFv2 Instance

```terraform
# Simple OSPFv2 backbone area
resource "mikrotik_ospf_instance_v7" "default" {
  name      = "default_v2"
  version   = "2"
  router_id = "1.1.1.1"
  comment   = "Main OSPF instance for IPv4"
}

resource "mikrotik_ospf_area_v7" "backbone" {
  name     = "backbone"
  instance = mikrotik_ospf_instance_v7.default.name
  area_id  = "0.0.0.0"
}

resource "mikrotik_ospf_interface_template_v7" "lan" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["192.168.1.0/24", "10.0.0.0/8"]
  cost     = 10
}
```

### OSPFv3 for IPv6

```terraform
resource "mikrotik_ospf_instance_v7" "ipv6" {
  name      = "ospf_v3"
  version   = "3"
  router_id = "1.1.1.1"  # Still uses IPv4 format even for OSPFv3
  comment   = "OSPF for IPv6 networks"
}

resource "mikrotik_ospf_area_v7" "backbone_v6" {
  name     = "backbone_v6"
  instance = mikrotik_ospf_instance_v7.ipv6.name
  area_id  = "0.0.0.0"
}

resource "mikrotik_ospf_interface_template_v7" "ipv6_lan" {
  area     = mikrotik_ospf_area_v7.backbone_v6.name
  networks = ["2001:db8::/64", "fd00::/8"]
  cost     = 10
}
```

### Route Redistribution

```terraform
resource "mikrotik_ospf_instance_v7" "redistributing" {
  name                   = "main_ospf"
  version                = "2"
  router_id              = "10.1.1.1"
  # Redistribute routes from other sources
  redistribute_connected = true   # Connected interfaces
  redistribute_static    = true   # Static routes
  redistribute_bgp       = false  # BGP routes (disabled)
  # Always originate default route
  originate_default      = "always"
  comment                = "OSPF with redistribution"
}
```

### VRF Isolation

```terraform
# VRF instance must exist first
resource "mikrotik_routing_table" "customer_a" {
  name = "customer_a"
  fib  = true
}

resource "mikrotik_ospf_instance_v7" "customer_ospf" {
  name          = "customer_a_ospf"
  version       = "2"
  router_id     = "172.16.1.1"
  vrf           = mikrotik_routing_table.customer_a.name
  routing_table = mikrotik_routing_table.customer_a.name
  comment       = "Isolated OSPF for Customer A"
}
```

### Multi-Instance with Domain ID

```terraform
# Primary OSPF instance
resource "mikrotik_ospf_instance_v7" "primary" {
  name      = "primary_ospf"
  version   = "2"
  router_id = "1.1.1.1"
  domain_id = "0.0.0.1"
  comment   = "Primary OSPF domain"
}

# Secondary OSPF instance (e.g., for different customer)
resource "mikrotik_ospf_instance_v7" "secondary" {
  name      = "secondary_ospf"
  version   = "2"
  router_id = "1.1.1.2"
  domain_id = "0.0.0.2"
  comment   = "Secondary OSPF domain"
}
```

### With Routing Filters

```terraform
resource "mikrotik_routing_filter_chain" "ospf_in" {
  name    = "ospf_in_filter"
  comment = "Filter incoming OSPF routes"
}

resource "mikrotik_routing_filter_rule" "reject_default" {
  chain = mikrotik_routing_filter_chain.ospf_in.name
  rule  = "if (dst == 0.0.0.0/0) { reject }"
}

resource "mikrotik_routing_filter_chain" "ospf_out" {
  name    = "ospf_out_filter"
  comment = "Filter outgoing OSPF routes"
}

resource "mikrotik_routing_filter_rule" "deny_private" {
  chain = mikrotik_routing_filter_chain.ospf_out.name
  rule  = "if (dst in 192.168.0.0/16) { reject }"
}

resource "mikrotik_ospf_instance_v7" "filtered" {
  name            = "filtered_ospf"
  version         = "2"
  router_id       = "10.0.0.1"
  in_filter_chain = mikrotik_routing_filter_chain.ospf_in.name
  out_filter_chain = mikrotik_routing_filter_chain.ospf_out.name
}
```

## Argument Reference

### Required

- `name` (String) - Name of the OSPF instance. Must be unique. Forces new resource on change.

### Optional

- `version` (String) - OSPF version: `"2"` for OSPFv2 (IPv4) or `"3"` for OSPFv3 (IPv6). Default: `"2"`.
- `router_id` (String) - Router ID in IPv4 address format (e.g., `1.1.1.1`). If not set, RouterOS auto-selects one from configured IP addresses.
- `domain_id` (String) - OSPF domain ID for multi-instance support (RFC 6549). Allows multiple OSPF instances to coexist. Format: IPv4 address (e.g., `0.0.0.1`).
- `disabled` (Boolean) - Disable this OSPF instance. Default: `false`.
- `comment` (String) - Comment for the OSPF instance.

### VRF and Routing Table

- `vrf` (String) - VRF name to associate this OSPF instance with. Requires VRF to be configured via `mikrotik_routing_table`.
- `routing_table` (String) - Routing table name to use. Leave empty for main table.

### Route Redistribution

- `redistribute_connected` (Boolean) - Redistribute directly connected routes. Default: `false`.
- `redistribute_static` (Boolean) - Redistribute static routes. Default: `false`.
- `redistribute_bgp` (Boolean) - Redistribute BGP routes. Default: `false`.
- `redistribute_rip` (Boolean) - Redistribute RIP routes. Default: `false`.
- `redistribute_ospf` (Boolean) - Redistribute routes from other OSPF instances. Default: `false`.

### Default Route Control

- `originate_default` (String) - Control default route (0.0.0.0/0) origination:
  - `"never"` - Never originate default route (default)
  - `"always"` - Always originate default route, even if not present in routing table
  - `"if-installed"` - Originate only if default route exists in routing table

### Routing Filters

- `in_filter_chain` (String) - Name of routing filter chain for incoming routes.
- `out_filter_chain` (String) - Name of routing filter chain for outgoing/redistributed routes.

## Attribute Reference

- `id` (String) - Unique identifier of the OSPF instance (.id).
- `routing_marks` (String) - Routing marks associated with this instance (computed).
- `dynamic` (Boolean) - Whether this instance was dynamically created (computed).
- `invalid` (Boolean) - Whether configuration is invalid (computed).

## Import

OSPF instances can be imported by name:

```shell
terraform import mikrotik_ospf_instance_v7.example default_v2
```

## Notes

### Router ID Selection

If `router_id` is not specified, RouterOS automatically selects one:
1. First preference: Lowest IP address from bridge interfaces
2. Second preference: Lowest IP address from all interfaces

For production deployments, **explicitly setting router_id is recommended** to ensure consistent behavior.

### Version Compatibility

- **RouterOS v7.x**: Unified OSPF configuration (`/routing/ospf`)
- **RouterOS v6.x**: Separate OSPF and OSPFv3 configurations (not supported by this resource)

### Redistribution Behavior

- Routes are redistributed as **Type-5 External LSAs** (Type-7 in NSSA areas)
- Default metric type is **Type-2** (external cost only)
- Use `out_filter_chain` to control which routes are redistributed and set attributes (metric, tag, etc.)

### Performance Considerations

- **CPU Usage**: OSPF SPF calculation is CPU-intensive. Limit area size to ~50 routers for low-end hardware.
- **Memory**: Each OSPF instance maintains full link-state database. Monitor memory usage in large networks.
- **Convergence**: Default timers (hello=10s, dead=40s) balance convergence speed vs. stability. Tune carefully.

## See Also

- [mikrotik_ospf_area_v7](ospf_area_v7.md) - Configure OSPF areas
- [mikrotik_ospf_interface_template_v7](ospf_interface_template_v7.md) - Configure OSPF interfaces
- [mikrotik_routing_filter_chain](routing_filter_chain.md) - Create routing filter chains
- [mikrotik_routing_table](routing_table.md) - Configure VRF routing tables
- [MikroTik OSPF Documentation](https://help.mikrotik.com/docs/display/ROS/OSPF)
