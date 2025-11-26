# mikrotik_system_logging

Creates and manages system logging rules that route log topics to specific actions (destinations).

## Features

- **Topic Routing:** Route specific log topics to different destinations
- **Multiple Topics:** Combine multiple topics in a single rule
- **Prefix Support:** Add prefixes for log categorization
- **Flexible Actions:** Reference any logging action by name
- **Disabled Rules:** Temporarily disable without deleting

## Example Usage

### Basic Firewall Logging to Remote Syslog

```hcl
resource "mikrotik_system_logging_action" "remote_syslog" {
  name   = "remote-syslog"
  target = "remote"
  remote = "192.168.1.10:514"
}

resource "mikrotik_system_logging" "firewall_logs" {
  topics = "firewall,info"
  action = mikrotik_system_logging_action.remote_syslog.name
  prefix = "FW"
}
```

### Critical System Events to Email

```hcl
resource "mikrotik_system_logging_action" "email_critical" {
  name   = "email-alerts"
  target = "email"
}

resource "mikrotik_system_logging" "critical_alerts" {
  topics = "critical,error"
  action = mikrotik_system_logging_action.email_critical.name
  prefix = "CRITICAL"
}
```

### BGP Events to Dedicated Log Server

```hcl
resource "mikrotik_system_logging_action" "bgp_server" {
  name            = "bgp-logs"
  target          = "remote"
  remote          = "10.0.0.50:514"
  bsd_syslog      = true
  syslog_facility = "daemon"
}

resource "mikrotik_system_logging" "bgp_logging" {
  topics = "bgp"
  action = mikrotik_system_logging_action.bgp_server.name
  prefix = "BGP"
}
```

### Wireless Events to Disk

```hcl
resource "mikrotik_system_logging_action" "wireless_disk" {
  name                = "wifi-logs"
  target              = "disk"
  disk_file_name      = "wireless"
  disk_file_count     = "5"
  disk_lines_per_file = "5000"
}

resource "mikrotik_system_logging" "wireless_logs" {
  topics = "wireless,info,debug"
  action = mikrotik_system_logging_action.wireless_disk.name
  prefix = "WIFI"
}
```

### Graylog Integration (Complete Stack)

```hcl
resource "mikrotik_system_logging_action" "graylog" {
  name            = "graylog-input"
  target          = "remote"
  remote          = "graylog.company.local:12201"
  bsd_syslog      = true
  syslog_facility = "local1"
}

resource "mikrotik_system_logging" "graylog_firewall" {
  topics = "firewall,info"
  action = mikrotik_system_logging_action.graylog.name
  prefix = "FW"
}

resource "mikrotik_system_logging" "graylog_system" {
  topics = "system,error,critical"
  action = mikrotik_system_logging_action.graylog.name
  prefix = "SYS"
}

resource "mikrotik_system_logging" "graylog_routing" {
  topics = "bgp,ospf,route"
  action = mikrotik_system_logging_action.graylog.name
  prefix = "ROUTING"
}
```

## Argument Reference

### Required Arguments

- `topics` - (Required, String) Comma-separated list of log topics to match. Examples: `firewall,info`, `system,error,critical`, `bgp`. See [Log Topics](#log-topics) for available topics.
- `action` - (Required, String) Reference to a logging action (destination) by name. Must match an existing `mikrotik_system_logging_action` resource.

### Optional Arguments

- `prefix` - (Optional, String) Optional prefix to add to log messages. Useful for filtering or categorizing logs. Example: `FW`, `BGP`, `AUTH`.
- `disabled` - (Optional, Bool) Whether the logging rule is disabled. Default: `false`.

## Attribute Reference

- `id` - The unique identifier of the logging rule (RouterOS internal ID).

## Import

Logging rules can be imported using their RouterOS ID:

```bash
terraform import mikrotik_system_logging.example *3
```

## Log Topics

RouterOS supports numerous log topics. Combine multiple topics with commas.

### Common Topics

| Topic | Description |
|-------|-------------|
| `account` | User login/logout, authentication |
| `bgp` | BGP routing protocol events |
| `critical` | Critical system events |
| `debug` | Debug-level messages (very verbose) |
| `dhcp` | DHCP server/client events |
| `error` | Error messages |
| `firewall` | Firewall rule matches |
| `info` | Informational messages |
| `ipsec` | IPSec VPN events |
| `l2tp` | L2TP VPN events |
| `ospf` | OSPF routing protocol events |
| `ppp` | PPP connection events |
| `pptp` | PPTP VPN events |
| `route` | Routing table changes |
| `script` | Script execution events |
| `sstp` | SSTP VPN events |
| `system` | System events |
| `warning` | Warning messages |
| `web-proxy` | Web proxy events |
| `wireless` | Wireless/WiFi events |

### Severity Levels

RouterOS doesn't use numeric syslog severity levels but has topic-based filtering:

| Topic | Similar to Syslog | Use Case |
|-------|-------------------|----------|
| `critical` | Emergency/Alert | Immediate attention required |
| `error` | Error | Error conditions |
| `warning` | Warning | Warning conditions |
| `info` | Informational | Normal informational messages |
| `debug` | Debug | Debug-level messages |

### Topic Combinations

Combine topics to filter by both category and severity:

```hcl
# Firewall informational messages
topics = "firewall,info"

# System errors and critical events
topics = "system,error,critical"

# BGP debug logging (very verbose)
topics = "bgp,debug"

# All routing protocols
topics = "bgp,ospf,route"

# All VPN events
topics = "ipsec,l2tp,pptp,sstp"

# Wireless events (info and debug)
topics = "wireless,info,debug"

# Script errors
topics = "script,error"

# Authentication events
topics = "account,info"
```

## Common Patterns

### Security Audit Logging

```hcl
resource "mikrotik_system_logging_action" "security_audit" {
  name            = "security-logs"
  target          = "remote"
  remote          = "audit.company.local:514"
  bsd_syslog      = true
  syslog_facility = "authpriv"
}

# Authentication events
resource "mikrotik_system_logging" "auth_logs" {
  topics = "account,info"
  action = mikrotik_system_logging_action.security_audit.name
  prefix = "AUTH"
}

# Firewall blocks (add action=drop log=yes in firewall rules)
resource "mikrotik_system_logging" "firewall_blocks" {
  topics = "firewall,info"
  action = mikrotik_system_logging_action.security_audit.name
  prefix = "FW-BLOCK"
}

# VPN connections
resource "mikrotik_system_logging" "vpn_logs" {
  topics = "ipsec,l2tp,pptp,sstp"
  action = mikrotik_system_logging_action.security_audit.name
  prefix = "VPN"
}
```

### Network Operations Center (NOC) Monitoring

```hcl
resource "mikrotik_system_logging_action" "noc_server" {
  name   = "noc-logs"
  target = "remote"
  remote = "noc.company.local:514"
}

# Critical and error events
resource "mikrotik_system_logging" "noc_critical" {
  topics = "critical,error"
  action = mikrotik_system_logging_action.noc_server.name
  prefix = "CRITICAL"
}

# Routing changes
resource "mikrotik_system_logging" "noc_routing" {
  topics = "bgp,ospf,route"
  action = mikrotik_system_logging_action.noc_server.name
  prefix = "ROUTING"
}

# Interface state changes
resource "mikrotik_system_logging" "noc_interfaces" {
  topics = "system,info"
  action = mikrotik_system_logging_action.noc_server.name
  prefix = "INTERFACE"
}
```

### Development/Troubleshooting

```hcl
resource "mikrotik_system_logging_action" "debug_disk" {
  name                = "debug-logs"
  target              = "disk"
  disk_file_name      = "debug"
  disk_file_count     = "5"
  disk_lines_per_file = "10000"
}

# Debug all routing protocols (VERY VERBOSE)
resource "mikrotik_system_logging" "routing_debug" {
  topics   = "bgp,ospf,route,debug"
  action   = mikrotik_system_logging_action.debug_disk.name
  prefix   = "ROUTING-DEBUG"
  disabled = true  # Enable only when troubleshooting
}

# Debug wireless issues
resource "mikrotik_system_logging" "wireless_debug" {
  topics   = "wireless,debug"
  action   = mikrotik_system_logging_action.debug_disk.name
  prefix   = "WIFI-DEBUG"
  disabled = true  # Enable only when troubleshooting
}
```

### Compliance/Audit Trail

```hcl
# Local disk for audit retention
resource "mikrotik_system_logging_action" "audit_disk" {
  name                = "audit-trail"
  target              = "disk"
  disk_file_name      = "audit"
  disk_file_count     = "30"  # 30 days retention
  disk_lines_per_file = "10000"
}

# Remote backup
resource "mikrotik_system_logging_action" "audit_remote" {
  name   = "audit-backup"
  target = "remote"
  remote = "audit.company.local:514"
}

# Log to both local and remote
resource "mikrotik_system_logging" "audit_local" {
  topics = "account,firewall,ipsec,system"
  action = mikrotik_system_logging_action.audit_disk.name
  prefix = "AUDIT"
}

resource "mikrotik_system_logging" "audit_remote" {
  topics = "account,firewall,ipsec,system"
  action = mikrotik_system_logging_action.audit_remote.name
  prefix = "AUDIT"
}
```

## Best Practices

1. **Use Prefixes:** Add meaningful prefixes to categorize logs: `FW`, `BGP`, `AUTH`, `VPN`.

2. **Combine Topics:** Use comma-separated topics to filter by category and severity: `firewall,info`, `system,error,critical`.

3. **Disabled for Debug:** Create debug logging rules with `disabled = true`. Enable temporarily when troubleshooting.

4. **Redundant Critical Logs:** Send critical/error topics to multiple actions (local + remote) for redundancy.

5. **Minimize Debug Logging:** Debug topics are very verbose. Use sparingly and only when needed.

6. **Action References:** Always use `mikrotik_system_logging_action.<name>.name` to ensure action exists before logging rule.

7. **Topic Validation:** Verify topic names match RouterOS documentation (case-sensitive).

8. **Separate by Purpose:** Create separate logging rules for different purposes (security, troubleshooting, monitoring).

## Troubleshooting

### Logs Not Appearing

1. **Check Action Exists:**
   ```bash
   /system logging action print
   ```

2. **Verify Logging Rule:**
   ```bash
   /system logging print
   ```

3. **Confirm Topics Match:**
   ```bash
   # Generate test log
   /log warning "Test message"
   
   # Check if logged
   /log print where message~"Test"
   ```

4. **Check Disabled Status:**
   ```hcl
   disabled = false  # Ensure not disabled
   ```

### Too Many Logs

1. **Refine Topics:** Remove verbose topics like `debug`:
   ```hcl
   # Instead of
   topics = "firewall,debug"  # Very verbose
   
   # Use
   topics = "firewall,info"   # Reasonable volume
   ```

2. **Use Firewall Logging Selectively:** Add `log=yes` only to important firewall rules, not all rules.

3. **Disable When Not Needed:**
   ```hcl
   disabled = true  # Disable debug logging after troubleshooting
   ```

### Missing Specific Events

1. **Check Topic Coverage:** Ensure relevant topics are included:
   ```hcl
   # Missing BGP events? Add bgp topic
   topics = "bgp,route"
   ```

2. **Verify Severity:** Include appropriate severity levels:
   ```hcl
   # Missing errors? Add error topic
   topics = "system,error,critical"
   ```

3. **Check RouterOS Version:** Some topics may not exist in older RouterOS versions.

## Performance Considerations

- **Debug Topics Impact CPU:** Debug logging can significantly impact router CPU. Use only when troubleshooting.
- **Firewall Logging Volume:** Firewall logging on high-traffic routers generates massive log volume. Use selectively.
- **Topic Specificity:** Use specific topics instead of broad logging to reduce CPU and bandwidth usage.
- **Local vs Remote:** Local disk/memory logging is faster than remote syslog. Use local for high-volume logs.

## Related Resources

- `mikrotik_system_logging_action` - Defines where logs are sent
- `mikrotik_firewall_filter` - Firewall rules with `log=yes` parameter
- `mikrotik_script` - Scripts that generate log messages

## See Also

- [RouterOS Logging Documentation](https://help.mikrotik.com/docs/display/ROS/Log)
- [RouterOS Log Topics](https://help.mikrotik.com/docs/display/ROS/Log#Log-Topics)
- [Syslog Best Practices](https://www.rsyslog.com/doc/master/configuration/index.html)
