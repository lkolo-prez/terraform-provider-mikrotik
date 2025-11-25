---
page_title: "mikrotik_routing_table Resource - terraform-provider-mikrotik"
subcategory: "Routing"
description: |-
  Creates a MikroTik Routing Table / VRF for RouterOS v7.
---

# mikrotik_routing_table (Resource)

Creates a MikroTik Routing Table (VRF) for RouterOS v7.

Routing tables in MikroTik RouterOS v7 act as Virtual Routing and Forwarding (VRF) instances, allowing you to maintain multiple independent routing tables on a single router. This is essential for:

- **Customer isolation** - ISPs can keep customer routes separate
- **Multi-tenant environments** - Different routing contexts per tenant
- **BGP route separation** - Associate BGP instances with specific VRFs
- **Management plane isolation** - Separate management traffic routing

## Example Usage

### Basic Routing Table / VRF

```terraform
resource "mikrotik_routing_table" "example" {
  name    = "example_vrf"
  fib     = "main"
  comment = "Example VRF for testing"
}
```

### Multi-VRF BGP Setup

```terraform
# Create VRF for Customer A
resource "mikrotik_routing_table" "customer_a" {
  name    = "customer_a"
  fib     = "main"
  comment = "Customer A VRF - Isolated routing instance"
}

# BGP instance using Customer A VRF
resource "mikrotik_bgp_instance_v7" "customer_a_bgp" {
  name          = "customer_a_bgp"
  as            = 65001
  router_id     = "10.0.1.1"
  routing_table = mikrotik_routing_table.customer_a.name
  vrf           = mikrotik_routing_table.customer_a.name
  
  redistribute_connected = true
  comment = "BGP instance for Customer A"
}

# BGP connection in Customer A VRF
resource "mikrotik_bgp_connection" "customer_a_peer" {
  name           = "customer_a_peer"
  remote_address = "10.0.1.2"
  remote_as      = 65003
  instance       = mikrotik_bgp_instance_v7.customer_a_bgp.name
  routing_table  = mikrotik_routing_table.customer_a.name
}
```

### Disabled Routing Table (Future Use)

```terraform
resource "mikrotik_routing_table" "future_vrf" {
  name     = "future_vrf"
  fib      = "main"
  disabled = true
  comment  = "Disabled VRF - reserved for future customer"
}
```

## Argument Reference

### Required

- `name` (String) - The name of the routing table. This will be used as the VRF identifier in other resources (e.g., BGP instances).

### Optional

- `fib` (String) - Forwarding Information Base (FIB) table to use. Default: `""` (empty string, uses default FIB).
- `disabled` (Boolean) - Whether the routing table is disabled. When disabled, the table won't be used for forwarding. Default: `false`.
- `comment` (String) - Comment for the routing table. Default: `""`.

### Read-Only

- `id` (String) - The unique identifier of the routing table.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the routing table resource.

## Import

Routing tables can be imported using their name:

```shell
terraform import mikrotik_routing_table.example customer_a
```

## RouterOS Documentation

For more information about Routing Tables in RouterOS v7, see the [official MikroTik documentation](https://help.mikrotik.com/docs/display/ROS/Routing+Table).

## Notes

- **RouterOS v7 only** - This resource requires RouterOS version 7 or higher. For RouterOS v6, routing tables are not user-configurable.
- **BGP Integration** - Use the `vrf` and `routing_table` attributes in `mikrotik_bgp_instance_v7` and `mikrotik_bgp_connection` resources to associate BGP with a specific VRF.
- **Default Table** - RouterOS creates a default "main" routing table automatically. You cannot delete or modify the main table.
- **FIB Selection** - The `fib` attribute determines which forwarding table is used. Typically set to "main" or left empty for default behavior.
