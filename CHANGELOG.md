# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **System Backup & Disaster Recovery (v1.8.0)**:
  - `mikrotik_system_backup` - System backup creation and management for DR, compliance, and pre-change snapshots
  - Auto-generated backup names with timestamps
  - Password-protected encrypted backups
  - Backup file tracking and verification
  - Integration examples: Pre-change backups, compliance snapshots, CI/CD automation
  - Comprehensive documentation (700+ lines) with security best practices
  - 15+ real-world examples for backup strategies
- **System Logging & Monitoring (v1.7.0)**:
  - `mikrotik_system_logging_action` - Log destination configuration (remote syslog, disk, memory, email, echo)
  - `mikrotik_system_logging` - Log topic routing with 20+ topics (firewall, bgp, system, wireless, etc.)
  - `mikrotik_snmp` - SNMP service configuration (SNMPv1/v2c, traps, contact, location)
  - `mikrotik_snmp_community` - SNMP community access control (read/write permissions, IP filtering)
  - Integration examples: Graylog, ELK, Splunk, Zabbix, PRTG, LibreNMS, Nagios, Observium
  - Comprehensive documentation (850+ lines) with enterprise patterns
  - 30+ real-world examples for system logging
  - 25+ real-world examples for SNMP monitoring
- **High Availability & NAT (v1.6.0)**:
  - `mikrotik_vrrp_interface` - VRRP v3 configuration (IPv4/IPv6, authentication, failover)
  - `mikrotik_firewall_nat` - NAT rule management (srcnat, dstnat, masquerade)
- **WiFi 6 / 802.11ax Support (v1.5.0)**:
  - `mikrotik_wifi_access_list` - WiFi access control lists
  - `mikrotik_wifi_channel` - WiFi channel configuration (2.4/5/6 GHz, DFS)
  - `mikrotik_wifi_configuration` - WiFi configuration profiles
  - `mikrotik_wifi_datapath` - WiFi data path configuration
  - `mikrotik_wifi_security` - WiFi security profiles (WPA2, WPA3, Enterprise)
  - `mikrotik_wireless_interface` - Legacy wireless interface
- **Routing & Container Support (v1.4.0)**:
  - `mikrotik_routing_filter_chain` - Routing filter chain management
  - `mikrotik_routing_filter_rule` - Routing filter rules (BGP, OSPF filtering)
  - `mikrotik_container` - Docker container support
  - `mikrotik_ospf_instance_v7` - OSPF v7 instance configuration
  - `mikrotik_ospf_area_v7` - OSPF v7 area configuration
  - `mikrotik_ospf_interface_template_v7` - OSPF v7 interface templates
- Full RouterOS 7.20+ BGP support with new resources:
  - `mikrotik_bgp_instance_v7` - BGP instance configuration
  - `mikrotik_bgp_connection` - BGP peer connections with templates
  - `mikrotik_bgp_template` - Reusable BGP connection templates
  - `mikrotik_bgp_session` - Data source for monitoring active sessions
- Comprehensive BGP v7 examples (6 production scenarios)
- Migration guide from BGP v6 to v7
- Performance optimizations (90% API reduction, 100x cache speedup)
- Automated release workflow with semantic versioning
- Enhanced CI/CD with multi-version Go testing (1.21, 1.22, 1.23)

### Changed
- Provider now includes 45 resources (up from 26 in v1.3.0)
- Updated GoReleaser to v2 format with improved changelog generation
- Modernized GitHub Actions workflows (Go 1.23, latest actions)
- Simplified documentation structure

### Deprecated
- `mikrotik_bgp_instance` - Use `mikrotik_bgp_instance_v7` for RouterOS 7.20+
- `mikrotik_bgp_peer` - Use `mikrotik_bgp_connection` for RouterOS 7.20+

### Fixed
- CI/CD compilation errors in BGP v7 tests
- Documentation redundancy and verbosity

## [0.9.1] - 2024-XX-XX

### Previous Releases
See [GitHub Releases](https://github.com/lkolo-prez/terraform-provider-mikrotik/releases) for complete history.

[Unreleased]: https://github.com/lkolo-prez/terraform-provider-mikrotik/compare/v0.9.1...HEAD
[0.9.1]: https://github.com/lkolo-prez/terraform-provider-mikrotik/releases/tag/v0.9.1
