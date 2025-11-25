---
page_title: "mikrotik_routing_filter_chain Resource - terraform-provider-mikrotik"
subcategory: "Routing"
description: |-
  Manages a RouterOS v7 routing filter chain that groups multiple filter rules.
---

# mikrotik_routing_filter_chain (Resource)

Manages a RouterOS v7 routing filter chain. Chains group multiple filter rules together and can be referenced by BGP connections, OSPF instances, and other routing protocols.

## Key Features

- **Rule Grouping**: Organize multiple filter rules under a single chain name
- **Protocol Integration**: Referenced by BGP, OSPF, and other routing protocols
- **Dynamic Chains**: Support for system-managed dynamic chains
- **Reusability**: One chain can be used by multiple BGP connections or routing instances

## Example Usage

### Basic Filter Chain

```terraform
resource "mikrotik_routing_filter_chain" "bgp_in" {
  name     = "bgp-input-filter"
  dynamic  = false
  disabled = false
  comment  = "BGP input filtering chain"
}
```

### Dynamic Chain

```terraform
resource "mikrotik_routing_filter_chain" "dynamic_filter" {
  name     = "bgp-dynamic"
  dynamic  = true
  comment  = "Dynamic chain managed by system"
}
```

### Chain with Rules

```terraform
resource "mikrotik_routing_filter_chain" "bgp_filtering" {
  name     = "bgp-in"
  dynamic  = true
  comment  = "Complete BGP input filtering"
}

resource "mikrotik_routing_filter_rule" "deny_default" {
  chain   = mikrotik_routing_filter_chain.bgp_filtering.name
  rule    = "if (dst == 0.0.0.0/0) { reject }"
  comment = "Deny default route"
}

resource "mikrotik_routing_filter_rule" "accept_customer" {
  chain   = mikrotik_routing_filter_chain.bgp_filtering.name
  rule    = "if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }"
  comment = "Accept customer routes"
}

resource "mikrotik_routing_filter_rule" "accept_all" {
  chain   = mikrotik_routing_filter_chain.bgp_filtering.name
  rule    = "accept"
  comment = "Accept all other routes"
}
```

### Use Chain in BGP Connection

```terraform
resource "mikrotik_routing_filter_chain" "bgp_in" {
  name    = "bgp-in-filter"
  dynamic = true
}

resource "mikrotik_routing_filter_chain" "bgp_out" {
  name    = "bgp-out-filter"
  dynamic = true
}

resource "mikrotik_bgp_connection" "isp" {
  name           = "isp-primary"
  template       = "default"
  remote_address = "192.0.2.1"
  remote_as      = 64512
  local_as       = 65001
  
  # Use filter chains
  input_filter  = mikrotik_routing_filter_chain.bgp_in.name
  output_filter = mikrotik_routing_filter_chain.bgp_out.name
}
```

### Multi-Stage Filtering with Jump

```terraform
# Main chain
resource "mikrotik_routing_filter_chain" "bgp_main" {
  name    = "bgp-main"
  comment = "Main BGP filtering chain"
}

# Customer-specific chain
resource "mikrotik_routing_filter_chain" "customer_a" {
  name    = "customer-a-filter"
  comment = "Customer A specific filtering"
}

# Main chain rule that jumps to customer chain
resource "mikrotik_routing_filter_rule" "check_customer" {
  chain   = mikrotik_routing_filter_chain.bgp_main.name
  rule    = "if (bgp-communities includes 65001:100) { jump customer-a-filter }"
  comment = "Jump to customer A filtering"
}

# Customer chain rules
resource "mikrotik_routing_filter_rule" "customer_accept" {
  chain   = mikrotik_routing_filter_chain.customer_a.name
  rule    = "if (dst in 10.100.0.0/16) { set bgp-local-pref 200; accept }"
  comment = "Accept customer A prefixes with high priority"
}

resource "mikrotik_routing_filter_rule" "customer_return" {
  chain   = mikrotik_routing_filter_chain.customer_a.name
  rule    = "return"
  comment = "Return to main chain"
}
```

### Separate Input/Output Chains

```terraform
# Input filtering chain
resource "mikrotik_routing_filter_chain" "bgp_in" {
  name    = "bgp-in"
  dynamic = true
  comment = "BGP input filtering - routes received from peers"
}

resource "mikrotik_routing_filter_rule" "in_deny_default" {
  chain = mikrotik_routing_filter_chain.bgp_in.name
  rule  = "if (dst == 0.0.0.0/0) { reject }"
}

resource "mikrotik_routing_filter_rule" "in_deny_bogons" {
  chain = mikrotik_routing_filter_chain.bgp_in.name
  rule  = "if (dst in 192.168.0.0/16 || dst in 10.0.0.0/8) { reject }"
}

resource "mikrotik_routing_filter_rule" "in_accept" {
  chain = mikrotik_routing_filter_chain.bgp_in.name
  rule  = "accept"
}

# Output filtering chain
resource "mikrotik_routing_filter_chain" "bgp_out" {
  name    = "bgp-out"
  dynamic = true
  comment = "BGP output filtering - routes advertised to peers"
}

resource "mikrotik_routing_filter_rule" "out_only_our_as" {
  chain = mikrotik_routing_filter_chain.bgp_out.name
  rule  = "if (bgp-as-path ~ \"^$\") { accept }"
  comment = "Only advertise routes originating from our AS"
}

resource "mikrotik_routing_filter_rule" "out_reject_all" {
  chain = mikrotik_routing_filter_chain.bgp_out.name
  rule  = "reject"
  comment = "Reject everything else"
}
```

### Complete Enterprise BGP Setup

```terraform
# ============================================================================
# Filter Chains
# ============================================================================

# Primary ISP input filter
resource "mikrotik_routing_filter_chain" "isp_primary_in" {
  name    = "isp-primary-in"
  dynamic = true
  comment = "Primary ISP input filtering"
}

# Backup ISP input filter (lower preference)
resource "mikrotik_routing_filter_chain" "isp_backup_in" {
  name    = "isp-backup-in"
  dynamic = true
  comment = "Backup ISP input filtering"
}

# Output filter (what we advertise)
resource "mikrotik_routing_filter_chain" "our_prefixes" {
  name    = "our-prefixes-out"
  dynamic = true
  comment = "Advertise only our prefixes"
}

# ============================================================================
# Primary ISP Rules
# ============================================================================

resource "mikrotik_routing_filter_rule" "primary_deny_default" {
  chain = mikrotik_routing_filter_chain.isp_primary_in.name
  rule  = "if (dst == 0.0.0.0/0) { reject }"
}

resource "mikrotik_routing_filter_rule" "primary_deny_bogons" {
  chain = mikrotik_routing_filter_chain.isp_primary_in.name
  rule  = "if (dst in 192.168.0.0/16 || dst in 10.0.0.0/8 || dst in 172.16.0.0/12) { reject }"
}

resource "mikrotik_routing_filter_rule" "primary_set_pref" {
  chain = mikrotik_routing_filter_chain.isp_primary_in.name
  rule  = "set bgp-local-pref 200; accept"
  comment = "Prefer primary ISP"
}

# ============================================================================
# Backup ISP Rules
# ============================================================================

resource "mikrotik_routing_filter_rule" "backup_deny_default" {
  chain = mikrotik_routing_filter_chain.isp_backup_in.name
  rule  = "if (dst == 0.0.0.0/0) { reject }"
}

resource "mikrotik_routing_filter_rule" "backup_deny_bogons" {
  chain = mikrotik_routing_filter_chain.isp_backup_in.name
  rule  = "if (dst in 192.168.0.0/16 || dst in 10.0.0.0/8 || dst in 172.16.0.0/12) { reject }"
}

resource "mikrotik_routing_filter_rule" "backup_set_pref" {
  chain = mikrotik_routing_filter_chain.isp_backup_in.name
  rule  = "set bgp-local-pref 100; accept"
  comment = "Lower preference for backup ISP"
}

# ============================================================================
# Output Rules (advertise our networks)
# ============================================================================

resource "mikrotik_routing_filter_rule" "out_our_networks" {
  chain = mikrotik_routing_filter_chain.our_prefixes.name
  rule  = "if (dst in 203.0.113.0/24) { accept }"
  comment = "Advertise our public network"
}

resource "mikrotik_routing_filter_rule" "out_reject_all" {
  chain = mikrotik_routing_filter_chain.our_prefixes.name
  rule  = "reject"
  comment = "Don't advertise anything else"
}

# ============================================================================
# BGP Connections
# ============================================================================

resource "mikrotik_bgp_connection" "isp_primary" {
  name           = "isp-primary"
  template       = "default"
  remote_address = "198.51.100.1"
  remote_as      = 64512
  local_as       = 65001
  
  input_filter  = mikrotik_routing_filter_chain.isp_primary_in.name
  output_filter = mikrotik_routing_filter_chain.our_prefixes.name
}

resource "mikrotik_bgp_connection" "isp_backup" {
  name           = "isp-backup"
  template       = "default"
  remote_address = "198.51.100.2"
  remote_as      = 64513
  local_as       = 65001
  
  input_filter  = mikrotik_routing_filter_chain.isp_backup_in.name
  output_filter = mikrotik_routing_filter_chain.our_prefixes.name
}
```

## Schema

### Required

- `name` (String) Name of the filter chain. Must be unique.

### Optional

- `comment` (String) Comment for the routing filter chain.
- `disabled` (Boolean) Whether the chain is disabled. Disabled chains are not evaluated. Default: `false`
- `dynamic` (Boolean) Whether the chain is dynamic. Dynamic chains can be modified by the system. Default: `false`

### Read-Only

- `id` (String) Unique ID of this resource.

## Import

Import using the chain name:

```shell
terraform import mikrotik_routing_filter_chain.example bgp-in-filter
```

## Notes

- **Chain Names**: Choose descriptive names like `bgp-in`, `bgp-out`, `customer-a-filter`
- **Rule Order**: Rules within a chain are evaluated in the order they appear in RouterOS
- **Dynamic Chains**: Set `dynamic = true` for chains that should be modifiable by the system
- **Empty Chains**: Chains can exist without rules and be populated later
- **Deletion**: Cannot delete chains that are referenced by active BGP connections or rules
- **Multiple References**: One chain can be used by multiple BGP connections

## Use Cases

### 1. **Basic BGP Filtering**
Single chain for input/output filtering with simple accept/reject rules.

### 2. **Multi-ISP Setup**
Separate chains for each ISP with different preferences and filtering policies.

### 3. **Customer-Specific Filtering**
Jump to customer-specific chains based on BGP communities or AS-path.

### 4. **Traffic Engineering**
Complex chains with multiple conditions to control route preference and path selection.

### 5. **Security Filtering**
Deny bogon networks, default routes, and implement prefix length limits.

## Related Resources

- [`mikrotik_routing_filter_rule`](routing_filter_rule.md) - Filter rules within chains
- [`mikrotik_bgp_connection`](bgp_connection.md) - BGP connections that reference chains
- [`mikrotik_bgp_template`](bgp_template.md) - BGP templates with filter configuration

## References

- [MikroTik RouterOS v7 Routing Filters Documentation](https://help.mikrotik.com/docs/display/ROS/Routing+Filters)
- [BGP Filtering Best Practices](https://help.mikrotik.com/docs/display/ROS/BGP+Filtering)
