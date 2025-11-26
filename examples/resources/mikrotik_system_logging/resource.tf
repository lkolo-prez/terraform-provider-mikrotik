# Example 1: Basic Firewall Logging to Remote Syslog
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

# Example 2: Critical System Events to Email
resource "mikrotik_system_logging_action" "email_critical" {
  name   = "email-alerts"
  target = "email"
}

resource "mikrotik_system_logging" "critical_alerts" {
  topics = "critical,error"
  action = mikrotik_system_logging_action.email_critical.name
  prefix = "CRITICAL"
}

# Example 3: BGP Events to Dedicated Log Server
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

# Example 4: Wireless Events to Disk
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

# Example 5: DHCP Events to Memory
resource "mikrotik_system_logging_action" "dhcp_memory" {
  name   = "dhcp-mem"
  target = "memory"
  memory = "500"
}

resource "mikrotik_system_logging" "dhcp_logs" {
  topics = "dhcp,info"
  action = mikrotik_system_logging_action.dhcp_memory.name
}

# Example 6: Script Execution Logs
resource "mikrotik_system_logging_action" "script_logs" {
  name                = "script-disk"
  target              = "disk"
  disk_file_name      = "scripts"
  disk_file_count     = "3"
  disk_lines_per_file = "2000"
}

resource "mikrotik_system_logging" "script_logging" {
  topics = "script,error,warning,info"
  action = mikrotik_system_logging_action.script_logs.name
  prefix = "SCRIPT"
}

# Example 7: Account/Authentication Events (Security Audit)
resource "mikrotik_system_logging_action" "security_syslog" {
  name            = "security-audit"
  target          = "remote"
  remote          = "audit.company.local:514"
  bsd_syslog      = true
  syslog_facility = "authpriv"
}

resource "mikrotik_system_logging" "auth_logs" {
  topics = "account,info"
  action = mikrotik_system_logging_action.security_syslog.name
  prefix = "AUTH"
}

# Example 8: VPN (IPSec) Logs
resource "mikrotik_system_logging" "ipsec_logs" {
  topics = "ipsec,info"
  action = mikrotik_system_logging_action.remote_syslog.name
  prefix = "VPN"
}

# Example 9: Multi-Topic Debugging
resource "mikrotik_system_logging_action" "debug_all" {
  name                = "debug-disk"
  target              = "disk"
  disk_file_name      = "debug-all"
  disk_file_count     = "10"
  disk_lines_per_file = "10000"
}

resource "mikrotik_system_logging" "debug_logging" {
  topics = "firewall,nat,route,bgp,ospf,wireless"
  action = mikrotik_system_logging_action.debug_all.name
  prefix = "DEBUG"
}

# Example 10: Web Proxy Logs
resource "mikrotik_system_logging" "web_proxy_logs" {
  topics = "web-proxy,info"
  action = mikrotik_system_logging_action.remote_syslog.name
  prefix = "PROXY"
}

# Example 11: Disabled Logging Rule (Template)
resource "mikrotik_system_logging" "disabled_example" {
  topics   = "system,info"
  action   = mikrotik_system_logging_action.remote_syslog.name
  prefix   = "DISABLED"
  disabled = true
}

# Example 12: Graylog Integration (Complete Stack)
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

# Example 13: Console Echo (Development/Troubleshooting)
resource "mikrotik_system_logging_action" "console" {
  name   = "console-echo"
  target = "echo"
}

resource "mikrotik_system_logging" "console_firewall" {
  topics = "firewall,info"
  action = mikrotik_system_logging_action.console.name
}

# Example 14: Backup Logging to Secondary Server
resource "mikrotik_system_logging_action" "backup_syslog" {
  name   = "backup-logs"
  target = "remote"
  remote = "192.168.1.20:514"
}

resource "mikrotik_system_logging" "backup_critical" {
  topics = "critical,error"
  action = mikrotik_system_logging_action.backup_syslog.name
  prefix = "BACKUP"
}

# Example 15: Rate-Limited High-Volume Logs (Memory Buffer)
resource "mikrotik_system_logging_action" "rate_limit_mem" {
  name   = "rate-limit-buffer"
  target = "memory"
  memory = "1000"
}

resource "mikrotik_system_logging" "high_volume_firewall" {
  topics = "firewall,info"
  action = mikrotik_system_logging_action.rate_limit_mem.name
}
