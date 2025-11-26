---
page_title: "mikrotik_snmp_community Resource - terraform-provider-mikrotik"
subcategory: "Monitoring"
description: |-
  Manages SNMP community access control for network monitoring systems.
---

# mikrotik_snmp_community (Resource)

Manages SNMP communities on MikroTik RouterOS devices. Communities provide access control for SNMP monitoring systems using community-based authentication (SNMPv1/v2c).

Each community defines:
- **Community name** (password equivalent)
- **Security level** (none, authorized, private)
- **Access permissions** (read/write)
- **IP address restrictions**

## Features

- **Multiple Communities**: Create separate communities for different monitoring systems
- **IP Filtering**: Restrict access by source IP address or network
- **Access Control**: Separate read and write permissions
- **Security Levels**: None (any), authorized (IP filter), private (reserved for SNMPv3)
- **Flexible Management**: Enable/disable without deletion

## Example Usage

```terraform
# Basic read-only community
resource "mikrotik_snmp_community" "public" {
  name        = "public"
  read_access = true
}

# Restricted to monitoring subnet
resource "mikrotik_snmp_community" "monitoring" {
  name        = "monitoring"
  read_access = true
  address     = "10.0.0.0/8"
}

# Zabbix specific community
resource "mikrotik_snmp_community" "zabbix" {
  name        = "zabbix"
  security    = "authorized"
  read_access = true
  address     = "192.168.1.50/32"
}

# Private read-write community
resource "mikrotik_snmp_community" "private" {
  name         = "private"
  security     = "private"
  read_access  = true
  write_access = true
  address      = "192.168.1.100"
}

# PRTG Network Monitor
resource "mikrotik_snmp_community" "prtg" {
  name        = "prtg"
  read_access = true
  address     = "192.168.100.10"
}

# LibreNMS monitoring
resource "mikrotik_snmp_community" "librenms" {
  name        = "librenms"
  read_access = true
  address     = "10.20.30.0/24"
}

# Disabled community (maintenance)
resource "mikrotik_snmp_community" "maintenance" {
  name        = "test"
  read_access = true
  disabled    = true
}

# Automation server with write access
resource "mikrotik_snmp_community" "automation" {
  name         = "automation"
  security     = "private"
  read_access  = true
  write_access = true
  address      = "10.50.50.50"
}

# Security operations center
resource "mikrotik_snmp_community" "soc" {
  name        = "security"
  security    = "private"
  read_access = true
  address     = "10.200.0.0/16"
}

# No IP restriction (not recommended)
resource "mikrotik_snmp_community" "any" {
  name        = "public-any"
  security    = "none"
  read_access = true
}
```

## Argument Reference

- `name` - (Required) SNMP community name. Acts as password for SNMPv1/v2c.
- `security` - (Optional) Security level: `"none"`, `"authorized"`, `"private"`. Default: `"none"`.
  - `none` - No restrictions (not recommended for production)
  - `authorized` - Requires IP address match
  - `private` - Reserved for SNMPv3 (encrypted)
- `read_access` - (Optional) Allow read access (GET operations). Default: `true`.
- `write_access` - (Optional) Allow write access (SET operations). Default: `false`.
- `address` - (Optional) Allowed source IP address or network (CIDR). Empty = all addresses.
- `disabled` - (Optional) Disable community without deletion. Default: `false`.

## Attribute Reference

- `id` - The RouterOS internal ID of the community.

## Import

```shell
# Import by RouterOS ID
terraform import mikrotik_snmp_community.example "*1"
```

## Security Levels

| Level | Description | IP Filter Required | Use Case |
|-------|-------------|-------------------|----------|
| `none` | No restrictions | No | Testing only |
| `authorized` | IP-based restriction | Yes | Production monitoring |
| `private` | SNMPv3 encryption | Yes | Encrypted access |

## Common Patterns

### Zabbix Monitoring

```terraform
resource "mikrotik_snmp" "zabbix" {
  enabled        = true
  trap_target    = "192.168.1.50"
  trap_community = "zabbix"
}

resource "mikrotik_snmp_community" "zabbix" {
  name        = "zabbix"
  security    = "authorized"
  read_access = true
  address     = "192.168.1.50/32"
}
```

### PRTG Network Monitor

```terraform
resource "mikrotik_snmp_community" "prtg" {
  name        = "prtg"
  read_access = true
  address     = "192.168.100.10"
}
```

### LibreNMS Monitoring

```terraform
resource "mikrotik_snmp_community" "librenms" {
  name        = "librenms"
  read_access = true
  address     = "10.20.30.0/24"
}
```

### Nagios/Icinga Monitoring

```terraform
resource "mikrotik_snmp_community" "nagios" {
  name        = "nagios"
  read_access = true
  address     = "10.100.1.10"
}
```

### Multi-Server Monitoring

```terraform
# Primary monitoring server
resource "mikrotik_snmp_community" "primary" {
  name        = "monitoring"
  read_access = true
  address     = "10.0.1.10"
}

# Backup monitoring server
resource "mikrotik_snmp_community" "backup" {
  name        = "monitoring"
  read_access = true
  address     = "10.0.1.11"
}
```

### Restricted Admin Access

```terraform
resource "mikrotik_snmp_community" "admin" {
  name         = "admin-write"
  security     = "private"
  read_access  = true
  write_access = true
  address      = "192.168.1.100"  # Admin workstation only
}
```

### Central Monitoring via Internet

```terraform
resource "mikrotik_snmp_community" "central" {
  name        = "branch-site"
  read_access = true
  address     = "203.0.113.0/24"  # Central monitoring public subnet
}
```

## Best Practices

1. **Always Use IP Restrictions**: Set `address` for production environments
2. **Avoid Write Access**: Rarely needed, prefer read-only monitoring
3. **Unique Community Names**: Don't use default "public" in production
4. **Strong Names**: Use random strings, not predictable names
5. **Separate Communities**: One per monitoring system for tracking
6. **Document Purpose**: Use descriptive names (zabbix, prtg, librenms)
7. **Regular Rotation**: Change community strings periodically
8. **Monitor Access**: Log SNMP access attempts

## Security Considerations

- **Plaintext Transmission**: Community strings sent unencrypted
- **IP Spoofing**: `address` provides limited protection
- **Write Access Risk**: SET operations can change configuration
- **SNMPv3 Preferred**: Use encrypted SNMPv3 for sensitive environments
- **Management VLAN**: Isolate SNMP traffic to dedicated network
- **Firewall Rules**: Restrict SNMP access at network layer (UDP/161)
- **Disable Unused**: Set `disabled = true` instead of deleting
- **Audit Regularly**: Review community list and access logs

## Troubleshooting

**Access denied:**
```bash
# Verify community name matches
snmpwalk -v2c -c monitoring 192.168.1.1

# Check IP restriction
# Must query from allowed source IP
```

**Timeout errors:**
- Verify `mikrotik_snmp` resource has `enabled = true`
- Check firewall rules allow UDP/161
- Verify source IP matches `address` restriction
- Test with `ping` first to confirm connectivity

**Wrong security level:**
- `security = "none"` - No IP filtering (dangerous)
- `security = "authorized"` - Requires `address` match
- `security = "private"` - For SNMPv3 only

**Write operations fail:**
- Verify `write_access = true`
- Check RouterOS user permissions
- SNMPv2c write support may be limited
- Consider API for configuration changes

## Performance Notes

- Communities have minimal performance impact
- Multiple communities don't affect SNMP speed
- IP filtering checked before SNMP processing
- Disabled communities have no overhead

## Monitoring Tool Integration

### Zabbix Configuration

```hcl
resource "mikrotik_snmp_community" "zabbix" {
  name        = "zabbix"
  read_access = true
  address     = "192.168.1.50/32"
}
```

**Zabbix Template:** Use "MikroTik by SNMP" template

### PRTG Configuration

```hcl
resource "mikrotik_snmp_community" "prtg" {
  name        = "prtg"
  read_access = true
  address     = "192.168.100.10"
}
```

**PRTG Sensor:** Add "SNMP Traffic Sensor" with community "prtg"

### LibreNMS Configuration

```hcl
resource "mikrotik_snmp_community" "librenms" {
  name        = "librenms"
  read_access = true
  address     = "10.20.30.0/24"
}
```

**LibreNMS:** Add device with community "librenms", force MikroTik OS detection

### Nagios/Icinga Configuration

```hcl
resource "mikrotik_snmp_community" "nagios" {
  name        = "nagios"
  read_access = true
  address     = "10.100.1.10"
}
```

**Nagios:** Use `check_snmp` plugin with `-C nagios`

## Notes

- Community names are case-sensitive
- Multiple communities can share the same name with different `address` restrictions
- Changes take effect immediately
- `disabled = true` prevents access without deletion
- Import supported using RouterOS ID
- Works with SNMPv1 and SNMPv2c (not v3)

## See Also

- [`mikrotik_snmp`](snmp.md) - Configure SNMP service
- [MikroTik SNMP Documentation](https://help.mikrotik.com/docs/display/ROS/SNMP)
- [SNMPv2c RFC 3416](https://www.rfc-editor.org/rfc/rfc3416.html)
