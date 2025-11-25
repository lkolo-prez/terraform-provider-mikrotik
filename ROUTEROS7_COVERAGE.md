# RouterOS v7 Feature Coverage Matrix

> **Provider Status**: RouterOS 7.14.3 - 7.16.2 Support  
> **Last Updated**: November 25, 2025

This document tracks which RouterOS v7 features are implemented in the Terraform Provider and which are planned for future releases.

## Legend

- ✅ **Fully Implemented** - Resource is production-ready
- 🟡 **Partially Implemented** - Basic functionality exists, advanced features missing
- 🔄 **In Progress** - Currently being developed
- 📋 **Planned** - Scheduled for future release
- ❌ **Not Planned** - Not suitable for Terraform or low priority
- ⚠️ **Deprecated** - Legacy v6 feature, use alternative

---

## I. General System & Updates

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| System Resource | ✅ | `mikrotik_system_resource` (data source) | Read-only system info |
| System Package | ❌ | - | Upgrades should be done manually |
| Device Mode | ❌ | - | One-time hardware configuration |
| System Note | 📋 | Planned: `mikrotik_system_note` | Login banner configuration |
| System History | ❌ | - | CLI-specific feature |
| Console Settings | ❌ | - | Not applicable to Terraform |
| System Scheduler | ✅ | `mikrotik_scheduler` | Fully implemented |
| System Script | ✅ | `mikrotik_script` | Fully implemented |

---

## II. Interfaces

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| **Ethernet** | | | |
| Basic Configuration | ✅ | Via bridge/VLAN resources | Indirect configuration |
| Ethernet Defaults | 📋 | Planned | v7-specific feature |
| Loop Protect | 📋 | Planned | v7-specific feature |
| **Bridge** | | | |
| Basic Bridge | ✅ | `mikrotik_bridge` | Fully implemented |
| Bridge Port | ✅ | `mikrotik_bridge_port` | Fully implemented |
| VLAN Filtering | ✅ | `mikrotik_interface_vlan7` | **New in v7** |
| Bridge VLAN | ✅ | `mikrotik_bridge_vlan` | VLAN tagging support |
| Hardware Offloading | 🟡 | Partial support | CRS3xx/CRS5xx specific |
| STP/RSTP/MSTP | 🟡 | Basic support in bridge | Advanced features missing |
| IGMP/MLD Snooping | 📋 | Planned | v7-specific feature |
| MVRP | 📋 | Planned | v7-specific feature |
| **VLAN** | | | |
| VLAN Interface (legacy) | ✅ | `mikrotik_vlan_interface` | v6 compatible |
| VLAN Interface (v7) | ✅ | `mikrotik_interface_vlan7` | **New in v7** |
| **WireGuard** | | | |
| WireGuard Interface | ✅ | `mikrotik_interface_wireguard` | Fully implemented |
| WireGuard Peer | ✅ | `mikrotik_interface_wireguard_peer` | Fully implemented |
| **WiFi (New System)** | | | |
| WiFi Radio | 📋 | Planned: `mikrotik_wifi_radio` | **New v7 802.11ax** |
| WiFi Channel | 📋 | Planned: `mikrotik_wifi_channel` | **New v7 feature** |
| WiFi Configuration | 📋 | Planned: `mikrotik_wifi_configuration` | **New v7 feature** |
| WiFi Security | 📋 | Planned: `mikrotik_wifi_security` | **New v7 feature** |
| WiFi Provisioning | 📋 | Planned: `mikrotik_wifi_provisioning` | **New v7 feature** |
| WiFi Access List | 📋 | Planned: `mikrotik_wifi_access_list` | **New v7 feature** |
| **WiFiWave2** | | | |
| WiFiWave2 Interface | 📋 | Planned | Alternative v7 WiFi |
| **Wireless (Legacy)** | | | |
| Wireless Interface | ✅ | `mikrotik_wireless_interface` | v6 legacy support |
| Wireless Security Profile | ✅ | `mikrotik_wireless_security_profile` | v6 legacy support |
| **Virtual Interfaces** | | | |
| veth (Virtual Ethernet) | 📋 | Planned: `mikrotik_interface_veth` | **New in v7** |
| vtx (VLAN Tunneling) | 📋 | Planned | **New in v7** |
| Bonding/LAG | 🟡 | Partial support | Basic bonding exists |
| EoIP/EoIPv6 | 🟡 | Partial support | v7 parameters missing |
| GRE | 🟡 | Partial support | v7 parameters missing |
| IPIP | 🟡 | Partial support | v7 parameters missing |
| **Interface Lists** | | | |
| Interface List | ✅ | `mikrotik_interface_list` | Fully implemented |
| Interface List Member | ✅ | `mikrotik_interface_list_member` | Fully implemented |
| **LTE** | | | |
| LTE Interface | 📋 | Planned | v7 APN profiles |
| **PPPoE** | | | |
| PPPoE Client | 📋 | Planned | v7 auth methods |

---

## III. IP Addressing & Services

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| **Addressing** | | | |
| IPv4 Address | ✅ | `mikrotik_ip_address` | Fully implemented |
| IPv6 Address | ✅ | `mikrotik_ipv6_address` | Fully implemented |
| **DHCP** | | | |
| DHCP Server | ✅ | `mikrotik_dhcp_server` | Fully implemented |
| DHCP Server Network | ✅ | `mikrotik_dhcp_server_network` | Fully implemented |
| DHCP Lease | ✅ | `mikrotik_dhcp_lease` | Static leases |
| DHCP Client | 📋 | Planned | v7 parameters |
| DHCPv6 Server | 📋 | Planned | v7 feature |
| **IP Pool** | | | |
| IP Pool | ✅ | `mikrotik_ip_pool` | Fully implemented |
| **DNS** | | | |
| DNS Settings | ✅ | `mikrotik_dns` | Fully implemented |
| DNS Static Entry | ✅ | `mikrotik_dns_record` | Fully implemented |
| DoH (DNS over HTTPS) | 📋 | Planned: DNS DoH support | **New in v7** |
| DNS Regexp Support | 📋 | Planned | v7 feature |
| **Services** | | | |
| IP Services | 📋 | Planned | SSH/HTTP/Winbox restrictions |
| **Neighbor Discovery** | | | |
| Neighbor Discovery | 📋 | Planned | v7 parameters |
| IPv6 ND | 📋 | Planned | v7 parameters |
| **Proxy/Socks** | | | |
| IP Proxy | ❌ | - | Low priority |
| IP Socks | ❌ | - | Low priority |

---

## IV. Routing

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| **Static Routes** | | | |
| IPv4 Static Route | ✅ | Via `mikrotik_bgp_*` | Basic support |
| IPv6 Static Route | 🟡 | Partial support | Needs v7 update |
| **Routing Tables** | | | |
| Routing Table | 📋 | Planned: `mikrotik_routing_table` | **New v7 VRF support** |
| **Routing Rules** | | | |
| Routing Rule | 📋 | Planned: `mikrotik_routing_rule` | **New v7 feature** |
| **VRF** | | | |
| VRF Configuration | 📋 | Planned: `mikrotik_vrf` | **New in v7** |
| **BGP** | | | |
| BGP Instance (v6) | ⚠️ | Deprecated | Use BGP Connection |
| BGP Peer (v6) | ⚠️ | Deprecated | Use BGP Connection |
| BGP Connection (v7) | ✅ | `mikrotik_bgp_connection` | **New in v7** ✅ |
| BGP Template (v7) | ✅ | `mikrotik_bgp_template` | **New in v7** ✅ |
| BGP Session Monitoring | 📋 | Planned (data source) | v7 feature |
| **OSPF** | | | |
| OSPF Instance (v7) | 📋 | Planned: `mikrotik_ospf_instance` | **Redesigned in v7** |
| OSPF Area (v7) | 📋 | Planned: `mikrotik_ospf_area` | **Redesigned in v7** |
| OSPF Interface Template | 📋 | Planned: `mikrotik_ospf_interface_template` | **New v7 concept** |
| **RIP** | | | |
| RIP Configuration | 📋 | Planned | v7 parameters |
| **Route Filters** | | | |
| Route Filter (v7) | 📋 | Planned: `mikrotik_routing_filter` | **Completely redesigned in v7** |

---

## V. Firewall

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| **Filter Rules** | | | |
| Firewall Filter (v6) | ✅ | `mikrotik_firewall_filter` | Legacy v6 support |
| Firewall Filter (v7) | 🟡 | Partial v7 support | Missing `untracked` state |
| Firewall RAW | ✅ | `mikrotik_firewall_raw` | **New in v7** ✅ |
| **NAT** | | | |
| NAT Rules | 🟡 | Partial support | v7 parameters missing |
| Port Forwarding | 🟡 | Via NAT rules | v7 updates needed |
| **Mangle** | | | |
| Mangle Rules | 🟡 | Partial support | v7 features missing |
| Connection Tracking | 📋 | Planned | v7 `untracked` state |
| **Address Lists** | | | |
| Address List | 🟡 | Partial support | v7 features missing |

---

## VI. Queues

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| **Simple Queues** | | | |
| Simple Queue | 🟡 | Partial support | Basic functionality |
| **Queue Tree** | | | |
| Queue Tree | 📋 | Planned | Hierarchical queuing |
| **Queue Types** | | | |
| PCQ | 🟡 | Via queue configuration | Indirect support |
| CAKE | 📋 | Planned: Queue type support | **New in v7** |
| fq_codel | 📋 | Planned: Queue type support | **New in v7** |
| RED/SFQ/FIFO | 🟡 | Partial support | Basic types |

---

## VII. Tools

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| Ping | ❌ | - | CLI tool, not for IaC |
| Traceroute | ❌ | - | CLI tool, not for IaC |
| Torch | ❌ | - | Monitoring tool |
| Packet Sniffer | ❌ | - | Diagnostic tool |
| Bandwidth Test | ❌ | - | Testing tool |
| Traffic Generator | ❌ | - | Testing tool |
| Profile | ❌ | - | Performance monitoring |
| Netwatch | 📋 | Planned | Useful for automation |
| Fetch | ❌ | - | Use external provisioners |

---

## VIII. Scripting

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| System Script | ✅ | `mikrotik_script` | Fully implemented |
| Script Execution | ✅ | Via `mikrotik_scheduler` | Indirect support |
| Error Handling (v7) | ✅ | Supported in scripts | v7 `:onerror` |
| Variable Scoping (v7) | ✅ | Supported in scripts | v7 improvements |

---

## IX. Wireless

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| **Legacy Wireless** | | | |
| Wireless Interface (v6) | ✅ | `mikrotik_wireless_interface` | Legacy support |
| Wireless Security (v6) | ✅ | `mikrotik_wireless_security_profile` | Legacy support |
| **WiFi (v7 - 802.11ax)** | | | |
| WiFi Radio | 📋 | Planned | **New v7 system** |
| WiFi Configuration | 📋 | Planned | **New v7 system** |
| WiFi Security | 📋 | Planned | WPA3 support |
| WiFi Provisioning | 📋 | Planned | Dynamic config |
| **WiFiWave2** | | | |
| WiFiWave2 Interface | 📋 | Planned | Alternative v7 WiFi |

---

## X. PPP

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| PPP Client | 📋 | Planned | v7 auth methods |
| PPP Secret | 📋 | Planned | User management |
| OpenVPN | 🟡 | Partial support | v7 ciphers missing |
| L2TPv3 | 📋 | Planned | **New in v7** |

---

## XI. System

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| System Resource | ✅ | `mikrotik_system_resource` (data) | Read-only |
| System Package | ❌ | - | Manual upgrades only |
| System Scheduler | ✅ | `mikrotik_scheduler` | Fully implemented |
| System Script | ✅ | `mikrotik_script` | Fully implemented |
| System Routerboard | ❌ | - | Hardware-specific |
| Reset Configuration | ❌ | - | Dangerous operation |

---

## XII. Files

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| File Management | ❌ | - | Use external provisioners |
| File Upload/Download | ❌ | - | Not suitable for Terraform |

---

## XIII. Log

| Feature | Status | Provider Resource | Notes |
|---------|--------|-------------------|-------|
| Log Configuration | 📋 | Planned | Logging setup |
| Log Actions | 📋 | Planned | Log routing |

---

## Summary Statistics

### Overall Coverage

| Status | Count | Percentage |
|--------|-------|------------|
| ✅ Fully Implemented | 23 | 28% |
| 🟡 Partially Implemented | 15 | 18% |
| 📋 Planned | 35 | 42% |
| ❌ Not Planned | 10 | 12% |
| **Total Features** | **83** | **100%** |

### Priority Features for Next Release

1. **High Priority** (Critical v7 features):
   - ✅ BGP Connection/Template (DONE)
   - ✅ Firewall RAW (DONE)
   - ✅ Interface VLAN7 (DONE)
   - 📋 Routing Filter (new system)
   - 📋 Routing Table/VRF
   - 📋 WiFi (new 802.11ax system)

2. **Medium Priority** (Enhanced v7 features):
   - 📋 OSPF Instance/Area/Templates
   - 📋 Queue Types (CAKE, fq_codel)
   - 📋 DNS DoH
   - 📋 veth Interface
   - 📋 Connection Tracking (untracked state)

3. **Low Priority** (Nice to have):
   - 📋 WiFiWave2
   - 📋 L2TPv3
   - 📋 Netwatch
   - 📋 Log Configuration

---

## How to Contribute

If you need a specific RouterOS v7 feature:

1. **Check this matrix** to see if it's planned
2. **Open an issue** on GitHub with:
   - Feature name from cheat sheet
   - Use case description
   - Example RouterOS commands
   - Priority justification
3. **Contribute code**:
   - Follow existing patterns (see `client/bgp_connection.go`)
   - Add tests
   - Update this matrix
   - Submit pull request

---

## References

- [RouterOS 7 Cheat Sheet](./ROUTEROS7_CHEATSHEET.md) - Full command reference
- [Migration Guide](./MIGRATION_ROUTEROS7.md) - Upgrade from v6
- [RouterOS 7 Support Doc](./ROUTEROS7_SUPPORT.md) - Feature documentation
- [Official MikroTik Docs](https://help.mikrotik.com/docs/spaces/ROS/pages/115736772/Upgrading+to+v7)

---

**Maintained by**: Community  
**Last Review**: November 25, 2025  
**RouterOS Versions**: 7.14.3, 7.16.2
