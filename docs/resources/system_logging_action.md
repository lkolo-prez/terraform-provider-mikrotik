# mikrotik_system_logging_action

Creates and manages system logging actions that define where and how logs are stored or transmitted.

## Features

- **Multiple Targets:** Remote syslog, disk, memory, email, console echo
- **Remote Syslog:** BSD syslog format, custom facilities, source IP
- **Disk Logging:** File rotation, configurable lines per file
- **Memory Logging:** Circular buffer with configurable size
- **Email Alerts:** Critical event notifications
- **Syslog Integration:** Graylog, ELK, Splunk, syslog-ng compatible

## Example Usage

### Remote Syslog Server

```hcl
resource "mikrotik_system_logging_action" "remote_syslog" {
  name   = "remote-syslog"
  target = "remote"
  remote = "192.168.1.10:514"
}
```

### Remote Syslog with BSD Format

```hcl
resource "mikrotik_system_logging_action" "remote_bsd" {
  name            = "remote-bsd-syslog"
  target          = "remote"
  remote          = "syslog.example.com:514"
  bsd_syslog      = true
  syslog_facility = "local0"
  src_address     = "192.168.1.1"
}
```

### Disk Logging with Rotation

```hcl
resource "mikrotik_system_logging_action" "disk_logs" {
  name                = "disk-logger"
  target              = "disk"
  disk_file_name      = "system-logs"
  disk_file_count     = "5"
  disk_lines_per_file = "5000"
}
```

### Memory Logging (Circular Buffer)

```hcl
resource "mikrotik_system_logging_action" "memory_logs" {
  name   = "memory-logger"
  target = "memory"
  memory = "500"
}
```

### Email Alerts for Critical Events

```hcl
resource "mikrotik_system_logging_action" "email_alerts" {
  name   = "critical-email"
  target = "email"
}
```

### Graylog/ELK Integration

```hcl
resource "mikrotik_system_logging_action" "graylog" {
  name            = "graylog-input"
  target          = "remote"
  remote          = "graylog.company.local:12201"
  bsd_syslog      = true
  syslog_facility = "local1"
}
```

## Argument Reference

### Required Arguments

- `name` - (Required, String, ForceNew) Name of the logging action. Used by logging rules to reference this action.
- `target` - (Required, String) Log target type. Valid values: `disk`, `echo`, `email`, `memory`, `remote`.

### Optional Arguments

#### Remote Syslog Configuration

- `remote` - (Optional, String) Remote syslog server address and port. Format: `ip:port` or `hostname:port`. Example: `192.168.1.10:514`. Required when `target="remote"`.
- `remote_port` - (Optional, String) Remote syslog port (deprecated, use `remote` with `ip:port` format instead). Default: `514`.
- `bsd_syslog` - (Optional, Bool) Use BSD syslog format for remote logging. Default: `false`.
- `syslog_facility` - (Optional, String) Syslog facility code. Default: `daemon`. Valid values:
  - System: `kern`, `user`, `mail`, `daemon`, `auth`, `syslog`, `lpr`, `news`
  - Services: `uucp`, `cron`, `authpriv`, `ftp`, `ntp`, `security`, `console`
  - Local: `local0`, `local1`, `local2`, `local3`, `local4`, `local5`, `local6`, `local7`
- `src_address` - (Optional, String) Source IP address for remote syslog connections. Uses router's egress interface IP by default.

#### Disk Logging Configuration

- `disk_file_name` - (Optional, String) Log file name for disk target. Example: `logs`. Only applicable when `target="disk"`.
- `disk_file_count` - (Optional, String) Number of log files for rotation. Default: `2`. Only applicable when `target="disk"`.
- `disk_lines_per_file` - (Optional, String) Maximum lines per log file before rotation. Default: `1000`. Only applicable when `target="disk"`.

#### Memory Logging Configuration

- `memory` - (Optional, String) Number of log lines to keep in memory. Default: `100`. Only applicable when `target="memory"`.
- `remember` - (Optional, Bool) Remember logs on disk even after reboot (for memory target). Default: `false`.

## Attribute Reference

- `id` - The unique identifier of the logging action (RouterOS internal ID).

## Import

Logging actions can be imported using their name:

```bash
terraform import mikrotik_system_logging_action.example remote-syslog
```

## Log Target Types

### remote
Send logs to remote syslog server via UDP/514. Supports BSD syslog format, custom facilities, and source IP binding.

**Use cases:**
- Centralized logging (Graylog, ELK, Splunk)
- Security audit trails
- Multi-site log aggregation
- Compliance requirements

**Configuration:**
- Set `remote` to `ip:port`
- Enable `bsd_syslog` for standard syslog format
- Choose `syslog_facility` based on log category
- Optionally set `src_address` for multi-homed routers

### disk
Store logs to local disk files with rotation. Useful for local troubleshooting and historical analysis.

**Use cases:**
- Local log retention
- Offline analysis
- Network-isolated devices
- Development/debugging

**Configuration:**
- Set `disk_file_name` (file name without extension)
- Configure `disk_file_count` for rotation (2-10 typical)
- Set `disk_lines_per_file` based on storage capacity

### memory
Store logs in RAM circular buffer. Fast access, limited retention, cleared on reboot.

**Use cases:**
- Real-time troubleshooting
- High-volume temporary logging
- Resource-constrained devices
- Short-term debugging

**Configuration:**
- Set `memory` to number of lines (100-5000 typical)
- Enable `remember` to persist across reboots

### email
Send log messages via email. Requires email server configuration in RouterOS.

**Use cases:**
- Critical event alerts
- Notifications to administrators
- Emergency notifications
- Low-frequency important events

**Configuration:**
- Ensure `/tool e-mail` is configured in RouterOS
- Typically used only for critical/error topics

### echo
Display logs on console (serial, Telnet, SSH). Useful for interactive troubleshooting.

**Use cases:**
- Interactive debugging
- Development testing
- Real-time event monitoring
- Console-based troubleshooting

## Syslog Facilities

Syslog facilities categorize log messages for filtering and routing:

| Facility | Code | Use Case |
|----------|------|----------|
| `kern` | 0 | Kernel messages |
| `user` | 1 | User-level messages |
| `mail` | 2 | Mail system |
| `daemon` | 3 | System daemons (default) |
| `auth` | 4 | Security/authentication |
| `syslog` | 5 | Syslog internal |
| `lpr` | 6 | Line printer subsystem |
| `news` | 7 | Network news |
| `uucp` | 8 | UUCP subsystem |
| `cron` | 9 | Clock daemon |
| `authpriv` | 10 | Security/authentication (private) |
| `ftp` | 11 | FTP daemon |
| `ntp` | 12 | NTP subsystem |
| `security` | 13 | Security/audit |
| `console` | 14 | Console |
| `local0-7` | 16-23 | Custom use (recommended for MikroTik) |

**Best Practice:** Use `local0-7` facilities for different log types:
- `local0` - Firewall logs
- `local1` - Routing (BGP, OSPF)
- `local2` - Wireless/WiFi
- `local3` - VPN (IPSec, WireGuard)
- `local4` - System events
- `local5` - Security/auth
- `local6` - Scripts
- `local7` - Custom applications

## Common Patterns

### Enterprise Syslog Server

```hcl
resource "mikrotik_system_logging_action" "corporate_syslog" {
  name            = "corp-logs"
  target          = "remote"
  remote          = "syslog.corp.local:514"
  bsd_syslog      = true
  syslog_facility = "local0"
  src_address     = var.router_mgmt_ip
}
```

### Redundant Logging (Primary + Backup)

```hcl
resource "mikrotik_system_logging_action" "primary_syslog" {
  name   = "primary-logs"
  target = "remote"
  remote = "syslog-primary.local:514"
}

resource "mikrotik_system_logging_action" "backup_syslog" {
  name   = "backup-logs"
  target = "remote"
  remote = "syslog-backup.local:514"
}

# Send critical logs to both servers
resource "mikrotik_system_logging" "critical_primary" {
  topics = "critical,error"
  action = mikrotik_system_logging_action.primary_syslog.name
}

resource "mikrotik_system_logging" "critical_backup" {
  topics = "critical,error"
  action = mikrotik_system_logging_action.backup_syslog.name
}
```

### Disk + Remote (Hybrid)

```hcl
# Local disk for all logs
resource "mikrotik_system_logging_action" "local_disk" {
  name                = "local-logs"
  target              = "disk"
  disk_file_name      = "all-logs"
  disk_file_count     = "7"
  disk_lines_per_file = "10000"
}

# Remote syslog for important logs only
resource "mikrotik_system_logging_action" "remote_important" {
  name   = "remote-important"
  target = "remote"
  remote = "syslog.local:514"
}

# Log everything locally
resource "mikrotik_system_logging" "all_local" {
  topics = "firewall,system,bgp,wireless,account"
  action = mikrotik_system_logging_action.local_disk.name
}

# Send critical to remote
resource "mikrotik_system_logging" "critical_remote" {
  topics = "critical,error"
  action = mikrotik_system_logging_action.remote_important.name
}
```

## Best Practices

1. **Use BSD Syslog Format:** Enable `bsd_syslog = true` for standard syslog servers (Graylog, ELK, Splunk).

2. **Choose Appropriate Facilities:** Use `local0-7` facilities to categorize MikroTik logs for easier filtering.

3. **Set Source IP:** Configure `src_address` for multi-homed routers to ensure consistent log source.

4. **Monitor Disk Space:** When using `target="disk"`, ensure sufficient storage and configure rotation.

5. **Memory Logs for Debugging:** Use `target="memory"` for temporary high-volume debugging, then disable.

6. **Email for Critical Only:** Limit `target="email"` to critical/error topics to avoid email flooding.

7. **Redundant Logging:** For production, send critical logs to both local disk and remote syslog.

8. **Naming Convention:** Use descriptive action names: `graylog-input`, `firewall-logs`, `security-audit`.

## Troubleshooting

### Logs Not Appearing on Remote Server

1. **Check Connectivity:**
   ```bash
   /ping 192.168.1.10
   ```

2. **Verify Firewall:**
   ```bash
   /ip firewall filter print
   # Ensure UDP/514 is allowed outbound
   ```

3. **Test with Console Echo:**
   ```hcl
   resource "mikrotik_system_logging_action" "test_console" {
     name   = "test"
     target = "echo"
   }
   
   resource "mikrotik_system_logging" "test_logs" {
     topics = "firewall,info"
     action = mikrotik_system_logging_action.test_console.name
   }
   ```

4. **Check RouterOS Logs:**
   ```bash
   /log print where topics~"system"
   ```

### Disk Full Errors

1. **Check Disk Usage:**
   ```bash
   /file print
   ```

2. **Reduce Log Retention:**
   ```hcl
   disk_file_count     = "2"   # Fewer files
   disk_lines_per_file = "1000" # Fewer lines
   ```

3. **Switch to Memory/Remote:**
   - Use `target="memory"` for temporary logs
   - Use `target="remote"` for centralized storage

### Missing Log Messages

1. **Verify Logging Rules:** Ensure `mikrotik_system_logging` resources reference correct action names.

2. **Check Topics:** Verify correct topic names (case-sensitive):
   ```hcl
   topics = "firewall,info"  # Correct
   topics = "Firewall,Info"  # Wrong
   ```

3. **Confirm Action Exists:**
   ```bash
   /system logging action print
   ```

## Security Notes

- **Syslog is UDP:** Not encrypted, can be intercepted. Use VPN tunnels for sensitive logs.
- **Source IP Spoofing:** Without VPN, source IP can be spoofed. Consider IPSec or WireGuard.
- **Sensitive Data:** Avoid logging passwords, keys, or sensitive configuration in clear text.
- **Access Control:** Restrict syslog server access to authorized IPs only.
- **TLS Syslog:** RouterOS doesn't support native TLS syslog. Use VPN for encryption.

## Performance Considerations

- **High-Volume Logs:** Use `target="memory"` or local disk for firewall logs with high traffic.
- **Remote Latency:** Network latency affects remote syslog performance. Use local logging as primary.
- **CPU Usage:** Excessive logging (debug topics) can impact router CPU. Use debug logging sparingly.
- **Bandwidth:** High-volume remote syslog consumes bandwidth. Monitor and adjust as needed.

## Related Resources

- `mikrotik_system_logging` - Routes log topics to actions
- `/tool e-mail` - Email server configuration (configured outside Terraform)

## See Also

- [RouterOS Logging Documentation](https://help.mikrotik.com/docs/display/ROS/Log)
- [Syslog RFC 5424](https://datatracker.ietf.org/doc/html/rfc5424)
- [Graylog GELF Input](https://docs.graylog.org/docs/gelf)
- [ELK Syslog Input](https://www.elastic.co/guide/en/logstash/current/plugins-inputs-syslog.html)
