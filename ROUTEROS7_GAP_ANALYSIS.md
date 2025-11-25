# RouterOS 7 Coverage Gap Analysis
**Document Version:** 1.0  
**Date:** December 2024  
**Provider:** terraform-provider-mikrotik  
**Repository:** https://github.com/lkolo-prez/terraform-provider-mikrotik

---

## Executive Summary

This document provides a comprehensive analysis of missing RouterOS v7 features in the Terraform provider. Based on the current implementation status documented in `ROUTEROS7_COVERAGE.md`, we have:

- ✅ **27 Fully Implemented** features (32%)
- 🟡 **15 Partially Implemented** features (18%)
- 📋 **31 Planned** features (37%)
- ❌ **10 Not Planned** features (12%)
- ⚠️ **1 Deprecated** feature (1%)

**Total:** 84 tracked features

**Recent Milestone:** BGP v7 complete implementation (4 resources, 133 attributes total)

---

## Priority 1: Critical Missing Features (P0)

### 1.1 Routing Infrastructure

#### 🚨 Routing Filter (Redesigned in v7)
**Status:** 📋 Planned  
**Priority:** P0 - CRITICAL  
**Effort:** Large (2-3 weeks)  
**RouterOS Path:** `/routing/filter/`

**Description:**  
RouterOS 7 completely redesigned the routing filter system, replacing old chain-based filters with a new rule-based system supporting multiple protocols.

**Required Resources:**
1. `mikrotik_routing_filter_rule` - Filter rules
2. `mikrotik_routing_filter_select_chain` - Selection chains
3. `mikrotik_routing_filter_chain` - Filter chains

**Key Attributes (estimate ~50 total):**
- Rule matching: `prefix`, `prefix-length`, `address-family`, `protocol`
- Actions: `accept`, `reject`, `jump`, `return`
- BGP-specific: `bgp-as-path`, `bgp-communities`, `bgp-ext-communities`
- OSPF-specific: `ospf-type`, `ospf-tag`
- Manipulation: `set-bgp-local-pref`, `set-bgp-med`, `set-bgp-weight`

**Dependencies:** None  
**Blockers:** Requires deep understanding of new filter syntax  
**Testing Requirements:** High complexity - integration tests with BGP, OSPF, RIP

---

#### 🚨 Routing Table / VRF
**Status:** 📋 Planned  
**Priority:** P0 - CRITICAL  
**Effort:** Medium (1-2 weeks)  
**RouterOS Path:** `/routing/table/`

**Description:**  
RouterOS 7 introduces proper VRF support with multiple routing tables, essential for enterprise deployments.

**Required Resources:**
1. `mikrotik_routing_table` - Routing table management
2. `mikrotik_routing_table_fib` - FIB (data source)

**Key Attributes (estimate ~15 total):**
- `name` - Table name
- `fib` - FIB selection
- `disabled` - Enable/disable
- `comment` - Description
- Statistics: route count, active routes, FIB entries

**Dependencies:** None  
**Blockers:** None  
**Testing Requirements:** Medium - validate VRF isolation

**Impact:** HIGH - Required for MPLS/VPN deployments (already referenced in BGP VRF example)

---

#### 🚨 OSPF v3 Redesign
**Status:** 📋 Planned  
**Priority:** P0 - HIGH  
**Effort:** Large (2-3 weeks)  
**RouterOS Path:** `/routing/ospf/`

**Description:**  
RouterOS 7 unified OSPF v2 and v3 into a single implementation with modern architecture.

**Required Resources:**
1. `mikrotik_ospf_instance_v7` - OSPF instance (replaces v6 instance)
2. `mikrotik_ospf_area_v7` - OSPF areas (replaces v6 area)
3. `mikrotik_ospf_interface_template_v7` - Interface templates (new in v7)

**Key Attributes (estimate ~60 total):**

**Instance (20 attrs):**
- `name`, `version` (2, 3)
- `router-id`, `domain-id`
- `redistribute-connected`, `redistribute-static`, `redistribute-bgp`
- `vrf`, `routing-table`

**Area (15 attrs):**
- `name`, `area-id`, `instance`
- `type` (default, stub, nssa)
- `default-cost`, `no-summaries`
- `stub-cost`, `nssa-translator`

**Interface Template (25 attrs):**
- `type` (broadcast, ptp, ptmp, nbma, virtual-link)
- `networks` - Network list
- `cost`, `priority`, `passive`
- `authentication` (none, simple, md5, sha1, sha256, sha384, sha512)
- `hello-interval`, `dead-interval`, `retransmit-interval`

**Dependencies:** Routing Filter (for redistribution)  
**Blockers:** Complex interaction with routing filter  
**Testing Requirements:** High - multi-area, authentication, redistribution

**Impact:** HIGH - Many enterprises use OSPF

---

### 1.2 WiFi (802.11ax - WiFi 6)

#### 🚨 WiFi System Redesign
**Status:** 📋 Planned  
**Priority:** P0 - HIGH  
**Effort:** Extra Large (3-4 weeks)  
**RouterOS Path:** `/interface/wifi/`, `/interface/wifiwave2/`

**Description:**  
RouterOS 7 introduced a completely new WiFi stack supporting 802.11ax (WiFi 6), replacing the old wireless and CAPsMAN systems.

**Required Resources:**
1. `mikrotik_interface_wifi` - WiFi interface configuration
2. `mikrotik_wifi_configuration` - WiFi profiles/configs
3. `mikrotik_wifi_datapath` - Data path settings
4. `mikrotik_wifi_security` - Security profiles
5. `mikrotik_wifi_channel` - Channel configuration
6. `mikrotik_wifi_access_list` - Access control

**Key Attributes (estimate ~120 total across all resources):**

**Interface (30 attrs):**
- `name`, `disabled`, `configuration`, `datapath`, `channel`
- `master-interface` (for virtual APs)
- `mac-address`, `arp`, `mtu`

**Configuration (40 attrs):**
- `name`, `mode` (ap, station, sniffer)
- `ssid`, `hide-ssid`
- `security` - Security profile
- `country`, `installation` (indoor, outdoor, any)
- `tx-power`, `tx-power-mode`
- `supported-rates`, `basic-rates`
- 802.11ax features: `he-guard-interval`, `he-frame-format`
- WPA3 support: `wpa3-sae`, `wpa3-sae-pk`

**Datapath (20 attrs):**
- `name`, `bridge`, `vlan-id`, `vlan-mode`
- `interface-list`, `arp`, `client-isolation`

**Security (20 attrs):**
- `name`, `authentication-types`
- `encryption` (wpa-psk, wpa2-psk, wpa3-sae)
- `passphrase`, `ft`, `ft-over-ds`
- `pmf` (disabled, optional, required)
- `eap-methods`, `eap-radius-accounting`

**Channel (10 attrs):**
- `name`, `band` (2ghz-ax, 5ghz-ax, 5ghz-n/ac, 6ghz-ax)
- `frequency`, `width` (20, 40, 80, 160)
- `secondary-frequency`

**Dependencies:** Interface management, bridge, VLAN  
**Blockers:** Requires RouterOS 7.13+ hardware with WiFi 6 support  
**Testing Requirements:** VERY HIGH - requires physical WiFi hardware, multiple clients

**Impact:** CRITICAL - WiFi 6 is essential for modern deployments

---

## Priority 2: Important Missing Features (P1)

### 2.1 Container

#### Container Management
**Status:** 📋 Planned  
**Priority:** P1 - HIGH  
**Effort:** Medium (1-2 weeks)  
**RouterOS Path:** `/container/`

**Description:**  
RouterOS 7.4+ supports running OCI containers directly on the router, enabling modern application deployments.

**Required Resources:**
1. `mikrotik_container_config` - Global container settings
2. `mikrotik_container` - Container instance
3. `mikrotik_container_env` - Environment variables
4. `mikrotik_container_mount` - Volume mounts

**Key Attributes (estimate ~35 total):**

**Config (8 attrs):**
- `registry-url`, `tmpdir`
- `ram-high`, `ram-low`
- `enabled`

**Container (15 attrs):**
- `name`, `interface`, `root-dir`
- `remote-image`, `tag`, `digest`
- `cmd`, `entrypoint`
- `workdir`, `hostname`, `domainname`
- `envlist`, `mounts`, `logging`

**Dependencies:** Interface (veth), DNS  
**Blockers:** Requires RouterOS 7.4+  
**Testing Requirements:** Medium - Docker registry integration

**Impact:** HIGH - Enables modern application hosting on routers

---

### 2.2 ZeroTier

#### ZeroTier Integration
**Status:** 📋 Planned  
**Priority:** P1 - MEDIUM  
**Effort:** Small (3-5 days)  
**RouterOS Path:** `/zerotier/`

**Description:**  
RouterOS 7.1+ includes native ZeroTier client for SD-WAN and overlay networking.

**Required Resources:**
1. `mikrotik_zerotier` - ZeroTier instance
2. `mikrotik_zerotier_controller` - Optional controller

**Key Attributes (estimate ~12 total):**
- `name`, `instance`, `disabled`
- `zt-address` (read-only)
- `allow-default`, `allow-global`, `allow-managed`
- `identity`, `port`, `comment`

**Dependencies:** None  
**Blockers:** None  
**Testing Requirements:** Low - basic connectivity tests

**Impact:** MEDIUM - Useful for hybrid cloud deployments

---

### 2.3 IP Services Enhancements

#### IP Services v7 Updates
**Status:** 🟡 Partially Implemented  
**Priority:** P1 - MEDIUM  
**Effort:** Small (1 week)  
**RouterOS Path:** `/ip/service/`

**Description:**  
RouterOS 7 added new service types and VRF support for IP services.

**Required Updates to `mikrotik_ip_service`:**
- Add `vrf` attribute (new in v7)
- Add new service types: `api-ssl`, `ovpn`
- Add `certificate` attribute for SSL services

**Key New Attributes (~5 total):**
- `vrf` - VRF instance name
- `certificate` - Certificate name for SSL
- Service type validation for v7

**Dependencies:** VRF implementation  
**Blockers:** None  
**Testing Requirements:** Low

**Impact:** MEDIUM - Improves service isolation in VRF environments

---

### 2.4 Connection Tracking

#### Connection Tracking v7 Features
**Status:** 🟡 Partially Implemented  
**Priority:** P1 - MEDIUM  
**Effort:** Small (3-5 days)  
**RouterOS Path:** `/ip/firewall/connection/tracking/`

**Description:**  
RouterOS 7 added `untracked` state and improved connection tracking.

**Required Updates to `mikrotik_ip_firewall_connection_tracking`:**
- Add `untracked` state support
- Add tracking modes

**Key New Attributes (~3 total):**
- `tracking-mode` (full, none)
- `untracked-connections` - Counter (read-only)

**Dependencies:** None  
**Blockers:** None  
**Testing Requirements:** Low

**Impact:** LOW-MEDIUM - Performance optimization feature

---

## Priority 3: Advanced Features (P2)

### 3.1 Queue Types

#### CAKE and fq_codel Queue Types
**Status:** 📋 Planned  
**Priority:** P2 - MEDIUM  
**Effort:** Medium (1 week)  
**RouterOS Path:** `/queue/type/`

**Description:**  
RouterOS 7 added modern queue algorithms: CAKE (Common Applications Kept Enhanced) and fq_codel (Fair Queue Controlled Delay).

**Required Updates to `mikrotik_queue_type`:**
- Add `cake` queue type with sub-options
- Add `fq-codel` queue type with sub-options

**Key New Attributes (~20 total):**

**CAKE attributes:**
- `bandwidth` - Target bandwidth
- `overhead` - Packet overhead
- `mpu` - Minimum packet unit
- `rtt` - Round-trip time
- `atm` - ATM mode
- `nat`, `wash`, `ack-filter`

**fq_codel attributes:**
- `target` - Target queue delay
- `interval` - Update interval
- `quantum` - Quantum size
- `flows` - Number of flows
- `ecn` - ECN support

**Dependencies:** None  
**Blockers:** None  
**Testing Requirements:** High - bandwidth shaping validation

**Impact:** MEDIUM - Improves QoS capabilities

---

### 3.2 DHCP Client Options

#### DHCP Client Option List
**Status:** 📋 Planned  
**Priority:** P2 - LOW  
**Effort:** Small (2-3 days)  
**RouterOS Path:** `/ip/dhcp-client/option/`

**Description:**  
RouterOS 7 allows custom DHCP client options.

**Required Resources:**
1. `mikrotik_dhcp_client_option` - Custom DHCP options

**Key Attributes (~8 total):**
- `name`, `code`, `value`

**Dependencies:** None  
**Blockers:** None  
**Testing Requirements:** Low

**Impact:** LOW - Niche feature

---

### 3.3 User Manager

#### User Manager v5
**Status:** ❌ Not Planned  
**Priority:** P2 - LOW  
**Effort:** Extra Large (4+ weeks)  
**RouterOS Path:** `/user-manager/`

**Description:**  
User Manager for RADIUS, billing, and user authentication.

**Reason for Not Planned:** Too complex, low demand, external RADIUS servers preferred

**Impact:** LOW - Most users prefer external solutions

---

## Priority 4: Interface Enhancements (P3)

### 4.1 L2TP Client & Server v7

**Status:** 📋 Planned  
**Priority:** P3 - LOW  
**Effort:** Small (3-5 days)  
**RouterOS Path:** `/interface/l2tp-client/`, `/interface/l2tp-server/`

**Description:**  
Minor updates to L2TP with VRF support.

**Key New Attributes:**
- `vrf` - VRF instance

**Impact:** LOW

---

### 4.2 PPPoE Client & Server v7

**Status:** 📋 Planned  
**Priority:** P3 - LOW  
**Effort:** Small (3-5 days)  
**RouterOS Path:** `/interface/pppoe-client/`, `/interface/pppoe-server/`

**Description:**  
VRF support for PPPoE.

**Key New Attributes:**
- `vrf` - VRF instance

**Impact:** LOW

---

### 4.3 SSTP Client & Server v7

**Status:** 📋 Planned  
**Priority:** P3 - LOW  
**Effort:** Small (3-5 days)  
**RouterOS Path:** `/interface/sstp-client/`, `/interface/sstp-server/`

**Description:**  
VRF support for SSTP.

**Key New Attributes:**
- `vrf` - VRF instance

**Impact:** LOW

---

### 4.4 VXLAN Interface

**Status:** 📋 Planned  
**Priority:** P3 - MEDIUM  
**Effort:** Small (1 week)  
**RouterOS Path:** `/interface/vxlan/`

**Description:**  
VXLAN tunneling for overlay networks.

**Required Resources:**
1. `mikrotik_interface_vxlan` - VXLAN interface

**Key Attributes (~12 total):**
- `name`, `vni`, `port`
- `local-address`, `remote-address`
- `group`, `vteps-ip-version`
- `mtu`, `arp`, `disabled`

**Dependencies:** None  
**Blockers:** None  
**Testing Requirements:** Medium

**Impact:** MEDIUM - Overlay networking

---

## Priority 5: System Enhancements (P3)

### 5.1 Hardware Offloading

**Status:** 📋 Planned  
**Priority:** P3 - LOW  
**Effort:** Medium (1 week)  
**RouterOS Path:** `/interface/ethernet/switch/`

**Description:**  
Hardware offloading settings for supported devices.

**Impact:** LOW - Hardware-specific

---

### 5.2 LCD Display

**Status:** ❌ Not Planned  
**Priority:** P3 - VERY LOW  
**Effort:** Small  
**RouterOS Path:** `/lcd/`

**Description:**  
LCD display configuration.

**Reason for Not Planned:** Very few devices, low automation value

**Impact:** VERY LOW

---

### 5.3 Power Consumption

**Status:** ❌ Not Planned  
**Priority:** P3 - VERY LOW  
**Effort:** Small  
**RouterOS Path:** `/system/power-consumption/`

**Description:**  
Power consumption monitoring.

**Reason for Not Planned:** Read-only data, low automation value

**Impact:** VERY LOW

---

## Summary Roadmap

### Phase 1: Critical Foundation (Q1 2025)
**Estimated Time:** 6-8 weeks

1. **Routing Table/VRF** (P0) - 2 weeks
   - Foundation for BGP VRF, OSPF multi-instance
   - Blocking other features

2. **Routing Filter** (P0) - 3 weeks
   - Essential for production BGP/OSPF deployments
   - Complex implementation

3. **OSPF v3 Redesign** (P0) - 3 weeks
   - High demand, depends on routing filter

**Deliverables:**
- 3 major resource groups (8+ resources)
- ~125 total attributes
- Full test coverage
- Production examples

---

### Phase 2: Modern Infrastructure (Q2 2025)
**Estimated Time:** 4-6 weeks

1. **WiFi System** (P0) - 4 weeks
   - 6 resources, ~120 attributes
   - Requires physical hardware for testing

2. **Container Support** (P1) - 2 weeks
   - 4 resources, ~35 attributes
   - Modern application hosting

**Deliverables:**
- 10 resources
- ~155 attributes
- Hardware testing lab for WiFi
- Container registry integration

---

### Phase 3: Enhancement & Polish (Q3 2025)
**Estimated Time:** 3-4 weeks

1. **Queue Types** (P2) - 1 week
2. **ZeroTier** (P1) - 1 week
3. **VXLAN** (P3) - 1 week
4. **Service/Connection Tracking Updates** (P1) - 1 week

**Deliverables:**
- 5+ resource updates
- ~50 new attributes
- Complete QoS capabilities

---

### Phase 4: Interface Completeness (Q4 2025)
**Estimated Time:** 2-3 weeks

1. **VPN Interface Updates** (P3) - 2 weeks
   - L2TP, PPPoE, SSTP VRF support
2. **Remaining P3 Features** - 1 week

**Deliverables:**
- Full VRF support across all interfaces
- 100% RouterOS 7 core feature coverage

---

## Estimated Total Effort

| Priority | Features | Estimated Time |
|----------|----------|----------------|
| P0 (Critical) | 4 | 8-10 weeks |
| P1 (High) | 5 | 4-5 weeks |
| P2 (Medium) | 3 | 2-3 weeks |
| P3 (Low) | 10+ | 4-5 weeks |
| **Total** | **22+** | **18-23 weeks** |

**Full-time equivalent:** 4-6 months  
**Part-time (20h/week):** 8-12 months

---

## Testing Requirements

### Hardware Requirements
- ✅ Basic RouterOS devices (CHR, hAP, RB5009) - **HAVE**
- ⚠️ WiFi 6 capable devices (wAP ax, cAP ax) - **NEED**
- ⚠️ Switch with hardware offloading - **OPTIONAL**

### Software Requirements
- ✅ RouterOS 7.14, 7.16, 7.17, 7.20+ - **HAVE**
- ✅ Terraform 1.5+ - **HAVE**
- ✅ Go 1.21, 1.22, 1.23 - **HAVE**
- ⚠️ Docker registry for container tests - **NEED**
- ⚠️ ZeroTier network for ZeroTier tests - **NEED**

### Testing Lab Topology
```
                    ┌──────────────┐
                    │   Internet   │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │  Core Router │
                    │  (RB5009)    │
                    └──┬────────┬──┘
            ┌──────────┘        └──────────┐
            │                               │
     ┌──────▼───────┐              ┌───────▼──────┐
     │  WiFi Router │              │  VPN Router  │
     │  (wAP ax)    │              │  (CHR)       │
     └──────────────┘              └──────────────┘
```

---

## Contribution Guidelines

### For New Features

1. **Research Phase:**
   - Read MikroTik documentation thoroughly
   - Test on RouterOS device manually
   - Document all attributes and behaviors

2. **Implementation Phase:**
   - Create client struct in `client/`
   - Create Terraform resource in `mikrotik/`
   - Follow existing patterns (see BGP v7 implementation)

3. **Testing Phase:**
   - Unit tests in `client/*_test.go`
   - Acceptance tests in `mikrotik/*_test.go`
   - Minimum 3 test scenarios per resource

4. **Documentation Phase:**
   - Update ROUTEROS7_COVERAGE.md
   - Create example in `examples/`
   - Add inline code comments

5. **Optimization Phase:**
   - Add bulk fetch functions (List*)
   - Consider caching for read-heavy operations
   - Profile performance

---

## Conclusion

The Terraform provider has made significant progress with **32% of RouterOS 7 features fully implemented**. The recent BGP v7 implementation demonstrates the maturity of the provider architecture with:

- ✅ 4 BGP resources (133 attributes total)
- ✅ Comprehensive test coverage (20+ test cases)
- ✅ Performance optimizations (90% API call reduction)
- ✅ Production-ready examples

**Immediate priorities:**
1. Routing Filter (blocking production BGP deployments)
2. Routing Table/VRF (blocking enterprise features)
3. OSPF v7 (high demand)
4. WiFi System (modern infrastructure)

**Long-term goal:**  
Achieve 80%+ coverage of RouterOS 7 core features by Q4 2025.

---

**Document Maintained By:** Provider Development Team  
**Last Updated:** December 2024  
**Next Review:** March 2025
