---
page_title: "mikrotik_routing_filter_rule Resource - terraform-provider-mikrotik"
subcategory: "Routing"
description: |-
  Manages a RouterOS v7 routing filter rule for filtering and manipulating routes.
---

# mikrotik_routing_filter_rule (Resource)

Manages a RouterOS v7 routing filter rule. Routing filters in RouterOS v7 use a completely redesigned rule-based syntax for filtering and manipulating routes from BGP, OSPF, and other routing protocols.

## Key Features

- **New v7 Syntax**: Uses if-then structure instead of v6's chain-based system
- **Powerful Matching**: Supports prefix matching, BGP attributes, OSPF attributes, communities
- **Route Manipulation**: Set local-pref, MED, weight, communities, AS-path
- **Multiple Actions**: Accept, reject, jump to other chains, return
- **Protocol Support**: BGP, OSPF, RIP, connected, static routes

## Example Usage

### Basic - Deny Default Route

```terraform
resource "mikrotik_routing_filter_rule" "deny_default" {
  chain    = "bgp-in"
  rule     = "if (dst == 0.0.0.0/0) { reject }"
  disabled = false
  comment  = "Deny default route from BGP"
}
```

### BGP Community Filtering

```terraform
resource "mikrotik_routing_filter_rule" "accept_customer" {
  chain    = "bgp-in"
  rule     = "if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }"
  disabled = false
  comment  = "Accept customer routes with community tag"
}
```

### Set BGP Local Preference

```terraform
resource "mikrotik_routing_filter_rule" "set_localpref" {
  chain    = "bgp-in"
  rule     = "if (bgp-communities includes 65001:200) { set bgp-local-pref 200; accept }"
  disabled = false
  comment  = "Set high local-pref for preferred routes"
}
```

### Prefix Length Filtering

```terraform
resource "mikrotik_routing_filter_rule" "block_small_prefixes" {
  chain = "bgp-in"
  rule  = "if (dst-len > 24) { reject }"
  comment = "Block prefixes smaller than /24"
}
```

### AS-Path Filtering

```terraform
resource "mikrotik_routing_filter_rule" "block_as_path" {
  chain = "bgp-in"
  rule  = "if (bgp-as-path ~ \"^65002_\") { reject }"
  comment = "Block routes originating from AS 65002"
}
```

### Multiple Conditions with Actions

```terraform
resource "mikrotik_routing_filter_rule" "customer_routes" {
  chain = "bgp-in"
  rule  = <<-EOT
    if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100 && dst-len <= 24) {
      set bgp-local-pref 150;
      set bgp-weight 100;
      accept
    }
  EOT
  comment = "Accept and prefer customer routes"
}
```

### Jump to Another Chain

```terraform
resource "mikrotik_routing_filter_chain" "customer_chain" {
  name = "customer-filtering"
}

resource "mikrotik_routing_filter_rule" "check_customer" {
  chain = "bgp-in"
  rule  = "if (bgp-communities includes 65001:100) { jump customer-filtering }"
  comment = "Jump to customer-specific filtering"
}

resource "mikrotik_routing_filter_rule" "customer_accept" {
  chain = "customer-filtering"
  rule  = "if (dst in 10.0.0.0/8) { accept }"
  comment = "Accept customer prefixes"
}
```

### Complete BGP Filtering Example

```terraform
# Create filter chain
resource "mikrotik_routing_filter_chain" "bgp_in" {
  name     = "bgp-in-filtering"
  dynamic  = true
  comment  = "BGP input filtering chain"
}

# Deny default route
resource "mikrotik_routing_filter_rule" "deny_default" {
  chain   = mikrotik_routing_filter_chain.bgp_in.name
  rule    = "if (dst == 0.0.0.0/0) { reject }"
  comment = "Block default route"
}

# Deny bogon networks
resource "mikrotik_routing_filter_rule" "deny_bogons" {
  chain = mikrotik_routing_filter_chain.bgp_in.name
  rule  = <<-EOT
    if (dst in 192.168.0.0/16 || dst in 172.16.0.0/12 || dst in 10.0.0.0/8) {
      reject
    }
  EOT
  comment = "Block private networks"
}

# Block small prefixes (anti-DDoS)
resource "mikrotik_routing_filter_rule" "block_small" {
  chain   = mikrotik_routing_filter_chain.bgp_in.name
  rule    = "if (dst-len > 24) { reject }"
  comment = "Block prefixes longer than /24"
}

# Set local-pref based on communities
resource "mikrotik_routing_filter_rule" "pref_high" {
  chain = mikrotik_routing_filter_chain.bgp_in.name
  rule  = "if (bgp-communities includes 65001:100) { set bgp-local-pref 200; accept }"
  comment = "High priority routes"
}

resource "mikrotik_routing_filter_rule" "pref_medium" {
  chain = mikrotik_routing_filter_chain.bgp_in.name
  rule  = "if (bgp-communities includes 65001:150) { set bgp-local-pref 150; accept }"
  comment = "Medium priority routes"
}

# Accept everything else
resource "mikrotik_routing_filter_rule" "accept_all" {
  chain   = mikrotik_routing_filter_chain.bgp_in.name
  rule    = "accept"
  comment = "Accept all other routes"
}

# Use in BGP connection
resource "mikrotik_bgp_connection" "isp" {
  name           = "isp-primary"
  template       = "default"
  remote_address = "192.0.2.1"
  remote_as      = 64512
  local_as       = 65001
  
  input_filter = mikrotik_routing_filter_chain.bgp_in.name
}
```

## Rule Syntax Reference

### Matching Conditions

| Matcher | Description | Example |
|---------|-------------|---------|
| `dst` | Destination prefix | `dst == 10.0.0.0/8` |
| `dst-len` | Prefix length | `dst-len <= 24` |
| `dst in` | Prefix in range | `dst in 10.0.0.0/8` |
| `bgp-as-path` | AS path regex | `bgp-as-path ~ "^65001_"` |
| `bgp-communities` | BGP community | `bgp-communities includes 65001:100` |
| `bgp-ext-communities` | Extended community | `bgp-ext-communities includes rt:65001:100` |
| `bgp-local-pref` | Local preference | `bgp-local-pref > 100` |
| `bgp-med` | Multi-exit discriminator | `bgp-med < 50` |
| `bgp-origin` | Origin type | `bgp-origin == igp` |
| `ospf-type` | OSPF route type | `ospf-type == external` |
| `ospf-tag` | OSPF tag | `ospf-tag == 100` |
| `protocol` | Routing protocol | `protocol == bgp` |

### Actions

| Action | Description | Example |
|--------|-------------|---------|
| `accept` | Accept route | `accept` |
| `reject` | Reject route | `reject` |
| `return` | Return from chain | `return` |
| `jump` | Jump to another chain | `jump "other-chain"` |
| `set bgp-local-pref` | Set local preference | `set bgp-local-pref 200` |
| `set bgp-med` | Set MED | `set bgp-med 50` |
| `set bgp-weight` | Set weight | `set bgp-weight 100` |
| `set bgp-communities` | Set communities | `set bgp-communities 65001:100` |
| `set bgp-ext-communities` | Set extended communities | `set bgp-ext-communities rt:65001:100` |

### Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `==` | Equal | `dst == 0.0.0.0/0` |
| `!=` | Not equal | `bgp-origin != incomplete` |
| `<`, `<=`, `>`, `>=` | Comparison | `dst-len <= 24` |
| `in` | In range | `dst in 10.0.0.0/8` |
| `~` | Regex match | `bgp-as-path ~ "^65001_"` |
| `includes` | Contains (communities) | `bgp-communities includes 65001:100` |
| `&&` | Logical AND | `dst in 10.0.0.0/8 && dst-len <= 24` |
| `\|\|` | Logical OR | `dst == 0.0.0.0/0 \|\| dst == ::/0` |

## Schema

### Required

- `chain` (String) Name of the filter chain this rule belongs to. Multiple rules can share the same chain and are evaluated in order.
- `rule` (String) Filter rule expression in RouterOS v7 syntax. Must follow the if-then structure with conditions and actions.

### Optional

- `comment` (String) Comment for the routing filter rule.
- `disabled` (Boolean) Whether the rule is disabled. Disabled rules are not evaluated. Default: `false`

### Read-Only

- `id` (String) Unique ID of this resource.
- `dynamic` (Boolean) Whether the rule is dynamically created by the system.
- `invalid` (Boolean) Whether the rule syntax is invalid. If true, check the rule expression for syntax errors.

## Import

Import using the rule ID:

```shell
terraform import mikrotik_routing_filter_rule.example *1
```

## Notes

- **Rule Order Matters**: Rules in the same chain are evaluated sequentially. Place more specific rules before generic ones.
- **Syntax Validation**: RouterOS validates rule syntax when creating/updating. Invalid rules will show `invalid=true`.
- **Chain Creation**: Chains are automatically created when referenced in rules, but explicit chain resources are recommended for better control.
- **BGP Integration**: Reference filter chains in BGP connections using `input_filter` and `output_filter` attributes.
- **Performance**: Keep filter chains simple. Complex expressions can impact routing performance.
- **Migration from v6**: v7 filter syntax is completely different from v6. Manual conversion is required.

## Related Resources

- [`mikrotik_routing_filter_chain`](routing_filter_chain.md) - Manage filter chains
- [`mikrotik_bgp_connection`](bgp_connection.md) - BGP connections that use filters
- [`mikrotik_bgp_template`](bgp_template.md) - BGP templates with filter configuration

## References

- [MikroTik RouterOS v7 Routing Filters Documentation](https://help.mikrotik.com/docs/display/ROS/Routing+Filters)
- [BGP Route Filtering Examples](https://help.mikrotik.com/docs/display/ROS/BGP+Filtering)
