# SNMP Community Examples

# Example 1: Basic read-only community
resource "mikrotik_snmp_community" "public_readonly" {
  name        = "public"
  read_access = true
}

# Example 2: Private read-write community with IP restriction
resource "mikrotik_snmp_community" "private_readwrite" {
  name         = "private"
  security     = "private"
  read_access  = true
  write_access = true
  address      = "192.168.1.100"
}

# Example 3: Monitoring subnet access
resource "mikrotik_snmp_community" "monitoring" {
  name        = "monitoring"
  read_access = true
  address     = "10.0.0.0/8"
}

# Example 4: Zabbix specific community
resource "mikrotik_snmp_community" "zabbix" {
  name        = "zabbix"
  security    = "authorized"
  read_access = true
  address     = "192.168.1.50/32"
}

# Example 5: Multiple IP addresses (use separate resources)
resource "mikrotik_snmp_community" "server1" {
  name        = "servers"
  read_access = true
  address     = "10.10.10.1"
}

# Example 6: Disabled community
resource "mikrotik_snmp_community" "maintenance" {
  name        = "maintenance"
  read_access = true
  disabled    = true
}

# Example 7: PRTG Network Monitor
resource "mikrotik_snmp_community" "prtg" {
  name        = "prtg"
  read_access = true
  address     = "192.168.100.10"
}

# Example 8: LibreNMS monitoring
resource "mikrotik_snmp_community" "librenms" {
  name        = "librenms"
  read_access = true
  address     = "10.20.30.0/24"
}

# Example 9: Nagios monitoring
resource "mikrotik_snmp_community" "nagios" {
  name        = "nagios"
  read_access = true
  address     = "10.100.1.10"
}

# Example 10: Observium monitoring
resource "mikrotik_snmp_community" "observium" {
  name        = "observium"
  read_access = true
  address     = "172.16.0.0/24"
}

# Example 11: Security operations center
resource "mikrotik_snmp_community" "soc" {
  name        = "security"
  security    = "private"
  read_access = true
  address     = "10.200.0.0/16"
}

# Example 12: Jump server access
resource "mikrotik_snmp_community" "jump" {
  name        = "jump"
  security    = "authorized"
  read_access = true
  address     = "10.10.10.5"
}

# Example 13: Automation server with write access
resource "mikrotik_snmp_community" "automation" {
  name         = "automation"
  security     = "private"
  read_access  = true
  write_access = true
  address      = "10.50.50.50"
}

# Example 14: Central monitoring via internet
resource "mikrotik_snmp_community" "central" {
  name        = "branch-site"
  read_access = true
  address     = "203.0.113.0/24"
}

# Example 15: No IP restriction (any source)
resource "mikrotik_snmp_community" "any_source" {
  name        = "public-any"
  security    = "none"
  read_access = true
}
