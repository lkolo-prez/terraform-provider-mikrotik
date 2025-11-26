# Example 1: Remote Syslog Server (Basic)
resource "mikrotik_system_logging_action" "remote_syslog" {
  name   = "remote-syslog"
  target = "remote"
  remote = "192.168.1.10:514"
}

# Example 2: Remote Syslog with BSD Format
resource "mikrotik_system_logging_action" "remote_bsd" {
  name          = "remote-bsd-syslog"
  target        = "remote"
  remote        = "syslog.example.com:514"
  bsd_syslog    = true
  syslog_facility = "local0"
  src_address   = "192.168.1.1"
}

# Example 3: Remote Syslog for Security Events
resource "mikrotik_system_logging_action" "security_syslog" {
  name            = "security-logs"
  target          = "remote"
  remote          = "10.0.0.50:514"
  bsd_syslog      = true
  syslog_facility = "authpriv"
  src_address     = "10.0.0.1"
}

# Example 4: Disk Logging with Rotation
resource "mikrotik_system_logging_action" "disk_logs" {
  name                 = "disk-logger"
  target               = "disk"
  disk_file_name       = "system-logs"
  disk_file_count      = "5"
  disk_lines_per_file  = "5000"
}

# Example 5: Memory Logging (Circular Buffer)
resource "mikrotik_system_logging_action" "memory_logs" {
  name   = "memory-logger"
  target = "memory"
  memory = "500"
}

# Example 6: Echo to Console
resource "mikrotik_system_logging_action" "console" {
  name   = "console-echo"
  target = "echo"
}

# Example 7: Email Alerts for Critical Events
resource "mikrotik_system_logging_action" "email_alerts" {
  name   = "critical-email"
  target = "email"
}

# Example 8: Graylog/ELK Integration
resource "mikrotik_system_logging_action" "graylog" {
  name            = "graylog-input"
  target          = "remote"
  remote          = "graylog.company.local:12201"
  bsd_syslog      = true
  syslog_facility = "local1"
}

# Example 9: Splunk Integration
resource "mikrotik_system_logging_action" "splunk" {
  name            = "splunk-hec"
  target          = "remote"
  remote          = "splunk.company.local:514"
  bsd_syslog      = true
  syslog_facility = "local2"
  src_address     = "172.16.0.1"
}

# Example 10: Multi-Site Syslog (Site A)
resource "mikrotik_system_logging_action" "site_a_syslog" {
  name            = "site-a-logs"
  target          = "remote"
  remote          = "10.10.1.100:514"
  bsd_syslog      = true
  syslog_facility = "local3"
}

# Example 11: Multi-Site Syslog (Site B)
resource "mikrotik_system_logging_action" "site_b_syslog" {
  name            = "site-b-logs"
  target          = "remote"
  remote          = "10.20.1.100:514"
  bsd_syslog      = true
  syslog_facility = "local4"
}

# Example 12: Development Debug Logging (Disk)
resource "mikrotik_system_logging_action" "debug_disk" {
  name                 = "debug-logs"
  target               = "disk"
  disk_file_name       = "debug"
  disk_file_count      = "3"
  disk_lines_per_file  = "10000"
}

# Example 13: High-Volume Memory Logger
resource "mikrotik_system_logging_action" "high_volume_memory" {
  name     = "hv-memory"
  target   = "memory"
  memory   = "2000"
  remember = true
}

# Example 14: Firewall Logs to Dedicated Server
resource "mikrotik_system_logging_action" "firewall_server" {
  name            = "fw-logs-server"
  target          = "remote"
  remote          = "192.168.100.10:5140"
  bsd_syslog      = true
  syslog_facility = "security"
}

# Example 15: BGP Events to NOC
resource "mikrotik_system_logging_action" "bgp_noc" {
  name            = "bgp-to-noc"
  target          = "remote"
  remote          = "noc.company.local:514"
  bsd_syslog      = true
  syslog_facility = "daemon"
}
