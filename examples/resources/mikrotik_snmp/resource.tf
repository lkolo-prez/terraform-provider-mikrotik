# SNMP Configuration Examples

# Example 1: Basic SNMP v2c
resource "mikrotik_snmp" "basic" {
  enabled  = true
  contact  = "noc@company.com"
  location = "Datacenter A, Rack 12"
}

# Example 2: SNMP with contact and location
resource "mikrotik_snmp" "monitoring" {
  enabled  = true
  contact  = "Network Operations Center"
  location = "Building 5, Floor 3"
}

# Example 3: SNMP with trap target (Zabbix monitoring)
resource "mikrotik_snmp" "zabbix" {
  enabled        = true
  contact        = "zabbix@company.com"
  location       = "Edge Router - Site 1"
  trap_version   = "2"
  trap_community = "traps"
  trap_target    = "192.168.1.50"  # Zabbix server
  trap_generators = "interfaces,temp-exception"
}

resource "mikrotik_snmp_community" "zabbix_readonly" {
  name        = "zabbix"
  read_access = true
  address     = "192.168.1.50/32"  # Only Zabbix server
}

# Example 4: IP restricted access
resource "mikrotik_snmp" "restricted" {
  enabled  = true
  contact  = "security@company.com"
  location = "Security Operations Center"
}

# Example 5: PRTG monitoring with traps
resource "mikrotik_snmp" "prtg" {
  enabled        = true
  contact        = "noc@company.com"
  location       = "Branch Office - London"
  trap_version   = "2"
  trap_community = "prtg-traps"
  trap_target    = "192.168.100.10"
}

# Example 6: LibreNMS monitoring
resource "mikrotik_snmp" "librenms" {
  enabled         = true
  contact         = "netops@company.com"
  location        = "Core Router - DC1"
  trap_version    = "2"
  trap_community  = "librenms"
  trap_target     = "10.20.30.40"
  trap_generators = "interfaces,system"
}

# Example 7: Branch site with NAT
resource "mikrotik_snmp" "branch_site" {
  enabled        = true
  contact        = "support@company.com"
  location       = "Branch Office - New York"
  trap_version   = "2"
  trap_community = "traps"
  trap_target    = "203.0.113.50"
}

# Example 8: SNMP disabled state
resource "mikrotik_snmp" "maintenance" {
  enabled  = false
  contact  = "maintenance@company.com"
  location = "Test Lab"
}

# Example 9: Nagios monitoring
resource "mikrotik_snmp" "nagios" {
  enabled        = true
  contact        = "nagios@company.com"
  location       = "Production Network"
  trap_version   = "1"
  trap_community = "nagios-traps"
  trap_target    = "10.100.1.10"
}

# Example 10: Filtered trap generators
resource "mikrotik_snmp" "multi_trap" {
  enabled         = true
  contact         = "noc@company.com"
  location        = "Critical Infrastructure"
  trap_version    = "2"
  trap_community  = "critical"
  trap_target     = "192.168.1.50"
  trap_generators = "temp-exception,voltage-exception,psu-fail"
}

# Example 11: Security monitoring
resource "mikrotik_snmp" "security" {
  enabled  = true
  contact  = "security@company.com"
  location = "Security Perimeter"
}

# Example 12: Default (disabled)
resource "mikrotik_snmp" "disabled" {
  enabled = false
}

# Example 13: Observium
resource "mikrotik_snmp" "observium" {
  enabled        = true
  contact        = "monitoring@company.com"
  location       = "Edge Router - Paris"
  trap_version   = "2"
  trap_community = "observium"
  trap_target    = "172.16.0.100"
}

# Example 14: Automation lab
resource "mikrotik_snmp" "config_mgmt" {
  enabled  = true
  contact  = "automation@company.com"
  location = "Automation Test Lab"
}

# Example 15: Complete Zabbix stack
resource "mikrotik_snmp" "complete_stack" {
  enabled         = true
  contact         = "noc@enterprise.com"
  location        = "Data Center - Primary Site"
  trap_version    = "2"
  trap_community  = "zabbix-traps"
  trap_target     = "10.0.100.10"
  trap_generators = "interfaces,system,temp-exception,voltage-exception"
}
