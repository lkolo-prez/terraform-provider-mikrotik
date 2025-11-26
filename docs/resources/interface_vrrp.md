# VRRP Interface Resource

Creates and manages a VRRP (Virtual Router Redundancy Protocol) interface for high availability setups.

## Features

- **High Availability**: Automatic gateway failover
- **VRRPv2 and VRRPv3**: Support for both protocol versions
- **Authentication**: Simple password or AH (Authentication Header)
- **State Hooks**: Execute scripts on master/backup transitions
- **IPv4/IPv6**: Full support for both protocols (VRRPv3)
- **Preemption**: Configurable priority-based master election

## Use Cases

- Active-backup router pairs
- Gateway redundancy
- Zero-downtime network maintenance
- Multi-WAN failover
- Load balancing across multiple gateways

## Example Usage

### Basic VRRP Setup (Master Router)

```terraform
resource "mikrotik_interface_vrrp" "gateway_master" {
  name            = "vrrp-gateway"
  interface       = "ether1"
  vrid            = 10
  priority        = 254          # Higher priority = Master
  
  version         = 3
  v3_protocol     = "ipv4"
  
  authentication  = "simple"
  password        = var.vrrp_password
  
  interval        = "1s"
  preemption_mode = true
}

# Virtual IP address
resource "mikrotik_ip_address" "vrrp_vip" {
  address   = "192.168.1.1/24"
  interface = mikrotik_interface_vrrp.gateway_master.name
}
```

### Backup Router Configuration

```terraform
resource "mikrotik_interface_vrrp" "gateway_backup" {
  name            = "vrrp-gateway"
  interface       = "ether1"
  vrid            = 10           # Same VRID as master
  priority        = 100          # Lower priority = Backup
  
  version         = 3
  v3_protocol     = "ipv4"
  
  authentication  = "simple"
  password        = var.vrrp_password
  
  interval        = "1s"
  preemption_mode = true
}
```

### VRRP with State Change Scripts

```terraform
resource "mikrotik_interface_vrrp" "gateway_with_hooks" {
  name            = "vrrp-gateway"
  interface       = "ether1"
  vrid            = 10
  priority        = 254
  
  version         = 3
  v3_protocol     = "ipv4"
  
  on_master       = "/system/script/run notify-master"
  on_backup       = "/system/script/run notify-backup"
}
```

### Multiple VRRP Groups (Load Balancing)

```terraform
# Router 1: Master for VRRP 10, Backup for VRRP 20
resource "mikrotik_interface_vrrp" "vrrp10_master" {
  name      = "vrrp10"
  interface = "ether1"
  vrid      = 10
  priority  = 254
}

resource "mikrotik_interface_vrrp" "vrrp20_backup" {
  name      = "vrrp20"
  interface = "ether1"
  vrid      = 20
  priority  = 100
}

# Router 2: Backup for VRRP 10, Master for VRRP 20
# (configure with reversed priorities)
```

## Argument Reference

### Required Arguments

- `name` - (Required) Name of the VRRP interface.
- `interface` - (Required) Physical interface name that VRRP will run on.
- `vrid` - (Required) Virtual Router ID (1-255). Must be the same on both master and backup routers.

### Optional Arguments

- `priority` - (Optional) Priority of the router (1-254). Higher priority means higher chance to become master. Default: `100`.
- `version` - (Optional) VRRP version (2 or 3). Default: `3`. Version 3 supports both IPv4 and IPv6.
- `authentication` - (Optional) Authentication type: `none`, `simple`, or `ah`. Default: `none`.
- `password` - (Optional) Password for authentication (if authentication is `simple` or `ah`). Marked as sensitive.
- `interval` - (Optional) Advertisement interval. Default: `1s`. Format: integer followed by time unit (s, ms).
- `preemption_mode` - (Optional) Enable preemption mode. If true, higher priority router will become master. Default: `true`.
- `v3_protocol` - (Optional) Protocol for VRRP version 3: `ipv4` or `ipv6`. Default: `ipv4`.
- `on_backup` - (Optional) Script to execute when router transitions to BACKUP state.
- `on_master` - (Optional) Script to execute when router transitions to MASTER state.
- `disabled` - (Optional) Whether the VRRP interface is disabled. Default: `false`.
- `comment` - (Optional) Comment for the VRRP interface.

### Computed Attributes

- `id` - The unique identifier of the VRRP interface.
- `running` - Whether the VRRP interface is currently running.

## Import

VRRP interfaces can be imported using the ID or name:

```bash
terraform import mikrotik_interface_vrrp.gateway_master *1
# or
terraform import mikrotik_interface_vrrp.gateway_master vrrp-gateway
```

## Notes

### High Availability Best Practices

1. **Same VRID**: Master and backup must use the same VRID (Virtual Router ID)
2. **Authentication**: Always use authentication in production
3. **Priority Planning**: Master should have priority 254, backup 100
4. **Preemption**: Enable for automatic master takeover
5. **Monitoring**: Use `running` attribute to monitor VRRP state

### Authentication Security

- **none**: No authentication (not recommended for production)
- **simple**: Plain text password (basic security)
- **ah**: Authentication Header (stronger security, VRRPv2 only)

For VRRPv3, consider using IPsec for secure communication.

### VRRPv2 vs VRRPv3

**VRRPv2:**
- IPv4 only
- Supports AH authentication
- RFC 3768

**VRRPv3:**
- IPv4 and IPv6 support
- No AH authentication (use IPsec instead)
- RFC 5798
- Recommended for new deployments

### State Transition Scripts

Scripts specified in `on_master` and `on_backup` are executed when VRRP state changes:

```terraform
on_master = "/system/script/run master-notify"
on_backup = "/system/script/run backup-notify"
```

Useful for:
- Email notifications
- Logging state changes
- Updating monitoring systems
- Custom failover actions

### Multiple VRRP Groups

You can run multiple VRRP groups on the same interface for load balancing:

- Group 10: Router A is master, Router B is backup
- Group 20: Router B is master, Router A is backup

This distributes traffic across both routers while maintaining redundancy.

## RouterOS Version Compatibility

- Minimum RouterOS version: 7.0
- VRRP is available on all RouterOS 7.x versions
- Both VRRPv2 and VRRPv3 are supported

## Troubleshooting

### VRRP Not Becoming Master

1. Check priority values (higher = master)
2. Verify VRID matches on both routers
3. Ensure authentication passwords match
4. Check physical interface is up
5. Verify no firewall blocking VRRP (IP protocol 112)

### Both Routers Claim Master

- Split-brain condition
- Check network connectivity between routers
- Verify VRRP packets are not blocked
- Check for network loops

### Monitoring VRRP Status

```bash
# RouterOS CLI
/interface vrrp print detail
/interface vrrp monitor [find name=vrrp-gateway]
```

In Terraform, use the `running` computed attribute to monitor state.
