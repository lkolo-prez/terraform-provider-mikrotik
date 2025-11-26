# GitHub Issues - Provider Development Roadmap

> **Auto-generated Issue Templates for Systematic Development**  
> Use this document to create GitHub Issues systematically

---

## 🏷️ Labels to Create

```bash
# Priority Labels
gh label create "P0-critical" --color "d73a4a" --description "Critical priority - blocking production use"
gh label create "P1-high" --color "ff9800" --description "High priority - important for enterprise"
gh label create "P2-medium" --color "fbca04" --description "Medium priority - nice to have"
gh label create "P3-low" --color "0e8a16" --description "Low priority - future consideration"

# Type Labels
gh label create "routeros-v7" --color "1d76db" --description "RouterOS v7 specific feature"
gh label create "enhancement" --color "a2eeef" --description "New feature or request"
gh label create "bug" --color "d73a4a" --description "Something isn't working"
gh label create "documentation" --color "0075ca" --description "Improvements or additions to documentation"

# Area Labels
gh label create "area:routing" --color "5319e7" --description "Routing (BGP, OSPF, Filters, VRF)"
gh label create "area:wifi" --color "5319e7" --description "WiFi / Wireless"
gh label create "area:firewall" --color "5319e7" --description "Firewall rules"
gh label create "area:interfaces" --color "5319e7" --description "Network interfaces"
gh label create "area:vpn" --color "5319e7" --description "VPN (WireGuard, IPsec, L2TP, etc.)"
gh label create "area:system" --color "5319e7" --description "System configuration"

# Status Labels
gh label create "status:in-progress" --color "c5def5" --description "Work in progress"
gh label create "status:blocked" --color "e99695" --description "Blocked by dependencies"
gh label create "status:ready" --color "0e8a16" --description "Ready for implementation"
```

---

## 📋 Milestones to Create

```bash
gh milestone create "Q1 2025 - Routing Foundation" --due-date "2025-03-31" --description "VRF, Routing Filter, OSPF v3"
gh milestone create "Q2 2025 - WiFi & Infrastructure" --due-date "2025-06-30" --description "WiFi 6, Container Support"
gh milestone create "Q3 2025 - Enhancements" --due-date "2025-09-30" --description "Queue Types, ZeroTier, VXLAN"
gh milestone create "Q4 2025 - Completeness" --due-date "2025-12-31" --description "VPN VRF, Advanced Features"
```

---

## 📝 Issues to Create

### PHASE 1: Routing Foundation (Q1 2025)

#### Issue #1: VRF / Routing Table Support

```bash
gh issue create \
  --title "[P0] Implement VRF / Routing Table support (mikrotik_routing_table)" \
  --label "P0-critical,routeros-v7,enhancement,area:routing" \
  --milestone "Q1 2025 - Routing Foundation" \
  --body "## Feature Description

**RouterOS Path:** \`/routing/table/\`  
**Priority:** P0 - CRITICAL  
**Estimated Effort:** 1-2 weeks  
**Attributes:** ~15

## Why Critical?

- Already referenced in BGP v7 examples (\`vrf = \"main\"\`) but NOT implemented
- Required for enterprise VPN/MPLS deployments
- Blocks advanced BGP and OSPF features
- Simplest P0 feature - QUICK WIN

## Use Case

Multi-tenancy, VRF isolation, MPLS L3VPN, BGP route separation per customer.

## Proposed Resources

- [ ] \`mikrotik_routing_table\` - Routing table/VRF management
- [ ] \`mikrotik_routing_table_fib\` - FIB information (data source)

## Example Configuration

\`\`\`hcl
resource \"mikrotik_routing_table\" \"customer_a\" {
  name     = \"customer_a\"
  fib      = true
  disabled = false
  comment  = \"Customer A VRF\"
}

resource \"mikrotik_routing_table\" \"management\" {
  name = \"management\"
  fib  = true
}

# Use in BGP
resource \"mikrotik_bgp_instance_v7\" \"customer_a_bgp\" {
  name          = \"customer_a\"
  as            = 65001
  router_id     = \"10.0.0.1\"
  routing_table = mikrotik_routing_table.customer_a.name
  vrf           = mikrotik_routing_table.customer_a.name
}
\`\`\`

## RouterOS CLI Example

\`\`\`
/routing/table
add name=customer_a fib=true comment=\"Customer A VRF\"
add name=management fib=true

/routing/bgp/instance
add name=customer_a as=65001 routing-table=customer_a vrf=customer_a
\`\`\`

## Key Attributes

**mikrotik_routing_table:**
- \`name\` (required) - Table name
- \`fib\` - Use FIB for forwarding
- \`disabled\` - Enable/disable
- \`comment\` - Description
- Statistics (read-only): route count, active routes

## Implementation Plan

1. Create \`client/routing_table.go\` with CRUD operations
2. Create \`mikrotik/resource_routing_table.go\` 
3. Create \`mikrotik/data_source_routing_table.go\`
4. Add tests: \`resource_routing_table_test.go\`
5. Update documentation
6. Add examples

## Dependencies

None (can be implemented immediately)

## Testing Requirements

- Medium complexity
- Validate VRF isolation
- Test with BGP integration
- Verify FIB selection

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 32)
- Coverage Matrix: ROUTEROS7_COVERAGE.md (line 156)
- BGP VRF Example: examples/bgp_v7_vrf.tf"
```

---

#### Issue #2: Routing Filter (New v7 System)

```bash
gh issue create \
  --title "[P0] Implement Routing Filter v7 (mikrotik_routing_filter_*)" \
  --label "P0-critical,routeros-v7,enhancement,area:routing" \
  --milestone "Q1 2025 - Routing Foundation" \
  --body "## Feature Description

**RouterOS Path:** \`/routing/filter/\`  
**Priority:** P0 - CRITICAL  
**Estimated Effort:** 2-3 weeks  
**Attributes:** ~50 total

## Why Critical?

- Completely redesigned in RouterOS v7
- Essential for BGP route filtering and manipulation
- Required for production BGP deployments
- Replaces old chain-based system with new rule-based syntax

## Use Case

BGP route filtering, prefix filtering, AS-path filtering, community manipulation, route-map functionality.

## Proposed Resources

- [ ] \`mikrotik_routing_filter_rule\` - Filter rules (new syntax)
- [ ] \`mikrotik_routing_filter_select_chain\` - Selection chains
- [ ] \`mikrotik_routing_filter_chain\` - Filter chains

## Example Configuration

\`\`\`hcl
# Deny default route
resource \"mikrotik_routing_filter_rule\" \"deny_default\" {
  chain    = \"bgp_in\"
  rule     = \"if (dst == 0.0.0.0/0) { reject }\"
  disabled = false
}

# Accept customer routes with community
resource \"mikrotik_routing_filter_rule\" \"accept_customer\" {
  chain    = \"bgp_in\"
  rule     = \"if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }\"
  disabled = false
}

# Set local-pref for preferred routes
resource \"mikrotik_routing_filter_rule\" \"set_localpref\" {
  chain = \"bgp_in\"
  rule  = \"if (bgp-communities includes 65001:200) { set bgp-local-pref 200; accept }\"
}

resource \"mikrotik_routing_filter_chain\" \"bgp_filtering\" {
  name           = \"bgp_in\"
  dynamic_chain  = true
}

# Use in BGP connection
resource \"mikrotik_bgp_connection\" \"peer\" {
  name          = \"peer1\"
  remote_address = \"10.0.0.2\"
  remote_as     = 65002
  input_filter  = mikrotik_routing_filter_chain.bgp_filtering.name
}
\`\`\`

## RouterOS CLI Example

\`\`\`
/routing/filter/rule
add chain=bgp_in rule=\"if (dst == 0.0.0.0/0) { reject }\"
add chain=bgp_in rule=\"if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }\"
add chain=bgp_in rule=\"if (bgp-communities includes 65001:200) { set bgp-local-pref 200; accept }\"

/routing/filter/chain
add name=bgp_in dynamic-chain=true

/routing/bgp/connection
set [find name=peer1] input-filter=bgp_in
\`\`\`

## Key Attributes

**mikrotik_routing_filter_rule (~30 attrs):**
- \`chain\` (required) - Chain name
- \`rule\` (required) - Filter rule in new v7 syntax
- \`disabled\` - Enable/disable

**Rule Syntax Supports:**
- Prefix matching: \`dst\`, \`dst-len\`, \`prefix-length\`
- BGP: \`bgp-as-path\`, \`bgp-communities\`, \`bgp-ext-communities\`, \`bgp-local-pref\`, \`bgp-med\`
- OSPF: \`ospf-type\`, \`ospf-tag\`
- Actions: \`accept\`, \`reject\`, \`jump\`, \`return\`
- Set operations: \`set bgp-local-pref\`, \`set bgp-med\`, \`set bgp-weight\`

## Implementation Plan

1. Study new v7 filter syntax (major change from v6)
2. Create \`client/routing_filter.go\` with rule parsing
3. Create resources for rule, chain, select-chain
4. Implement rule validation
5. Add comprehensive tests (complex logic)
6. Document new syntax with examples
7. Migration guide from v6 filters

## Dependencies

- None (standalone feature)

## Testing Requirements

- HIGH complexity - rule syntax validation
- Integration tests with BGP
- Test prefix matching, AS-path, communities
- Verify action execution (accept/reject/jump)

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 18)
- MikroTik Docs: https://help.mikrotik.com/docs/display/ROS/Routing+Filters"
```

---

#### Issue #3: OSPF v3 Redesign

```bash
gh issue create \
  --title "[P0] Implement OSPF v3 redesign (mikrotik_ospf_*_v7)" \
  --label "P0-critical,routeros-v7,enhancement,area:routing" \
  --milestone "Q1 2025 - Routing Foundation" \
  --body "## Feature Description

**RouterOS Path:** \`/routing/ospf/\`  
**Priority:** P0 - HIGH  
**Estimated Effort:** 2-3 weeks  
**Attributes:** ~60 total

## Why Critical?

- OSPF is the second most popular routing protocol (after BGP)
- Completely redesigned in RouterOS v7 (unified v2/v3)
- New interface template concept
- Required for enterprise networks

## Use Case

Internal routing, data center fabric, campus networks, multi-area OSPF.

## Proposed Resources

- [ ] \`mikrotik_ospf_instance_v7\` - OSPF instance (replaces v6)
- [ ] \`mikrotik_ospf_area_v7\` - OSPF areas
- [ ] \`mikrotik_ospf_interface_template_v7\` - Interface templates (NEW in v7)

## Example Configuration

\`\`\`hcl
resource \"mikrotik_ospf_instance_v7\" \"main\" {
  name                  = \"default\"
  version               = 2
  router_id             = \"10.0.0.1\"
  vrf                   = \"main\"
  
  redistribute_connected = true
  redistribute_static    = true
  redistribute_bgp       = false
  
  disabled = false
  comment  = \"Main OSPF instance\"
}

resource \"mikrotik_ospf_area_v7\" \"backbone\" {
  name     = \"backbone\"
  area_id  = \"0.0.0.0\"
  instance = mikrotik_ospf_instance_v7.main.name
  type     = \"default\"
  disabled = false
}

resource \"mikrotik_ospf_area_v7\" \"branch\" {
  name        = \"branch_offices\"
  area_id     = \"0.0.0.10\"
  instance    = mikrotik_ospf_instance_v7.main.name
  type        = \"stub\"
  stub_cost   = 10
  no_summaries = true
}

resource \"mikrotik_ospf_interface_template_v7\" \"lan_interfaces\" {
  area      = mikrotik_ospf_area_v7.backbone.name
  networks  = [\"10.0.1.0/24\", \"10.0.2.0/24\"]
  cost      = 10
  priority  = 1
  type      = \"broadcast\"
  
  auth      = \"sha256\"
  auth_key  = \"MySecurePassword123\"
  
  hello_interval      = 10
  dead_interval       = 40
  retransmit_interval = 5
  
  passive   = false
  disabled  = false
}

resource \"mikrotik_ospf_interface_template_v7\" \"wan_p2p\" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = [\"192.168.100.0/30\"]
  type     = \"ptp\"
  cost     = 100
}
\`\`\`

## RouterOS CLI Example

\`\`\`
/routing/ospf/instance
add name=default version=2 router-id=10.0.0.1 vrf=main \\
    redistribute-connected=yes redistribute-static=yes

/routing/ospf/area
add name=backbone area-id=0.0.0.0 instance=default type=default
add name=branch_offices area-id=0.0.0.10 instance=default type=stub stub-cost=10

/routing/ospf/interface-template
add area=backbone networks=10.0.1.0/24,10.0.2.0/24 cost=10 type=broadcast \\
    auth=sha256 hello-interval=10 dead-interval=40
add area=backbone networks=192.168.100.0/30 type=ptp cost=100
\`\`\`

## Key Attributes

**mikrotik_ospf_instance_v7 (~20 attrs):**
- \`name\`, \`version\` (2, 3)
- \`router_id\`, \`domain_id\`
- \`redistribute_*\` (connected, static, bgp, rip)
- \`vrf\`, \`routing_table\`

**mikrotik_ospf_area_v7 (~15 attrs):**
- \`name\`, \`area_id\`, \`instance\`
- \`type\` (default, stub, nssa)
- \`default_cost\`, \`no_summaries\`
- \`stub_cost\`, \`nssa_translator\`

**mikrotik_ospf_interface_template_v7 (~25 attrs):**
- \`area\`, \`networks\`
- \`type\` (broadcast, ptp, ptmp, nbma, virtual-link)
- \`cost\`, \`priority\`, \`passive\`
- \`auth\` (none, simple, md5, sha1, sha256, sha384, sha512)
- \`hello_interval\`, \`dead_interval\`, \`retransmit_interval\`

## Implementation Plan

1. Deprecate old v6 OSPF resources
2. Create \`client/ospf_v7.go\` for all three resources
3. Implement instance, area, interface-template resources
4. Add authentication support (multiple hash types)
5. Implement area type validation (default/stub/nssa)
6. Add comprehensive tests (multi-area setup)
7. Document migration from v6 to v7

## Dependencies

- VRF/Routing Table (#1) - for VRF support
- Routing Filter (#2) - for redistribution filtering

## Testing Requirements

- HIGH complexity
- Multi-area OSPF testing
- Authentication testing (MD5, SHA256)
- Area types (stub, NSSA)
- Interface templates with different networks

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 68)
- Coverage Matrix: ROUTEROS7_COVERAGE.md (line 172)
- MikroTik Docs: https://help.mikrotik.com/docs/display/ROS/OSPF"
```

---

### PHASE 2: WiFi & Infrastructure (Q2 2025)

#### Issue #4: WiFi 6 (802.11ax) System ✅ IMPLEMENTED v1.5.0

```bash
gh issue create \
  --title "[P0] Implement WiFi 6 / 802.11ax system (mikrotik_wifi_*)" \
  --label "P0-critical,routeros-v7,enhancement,area:wifi" \
  --milestone "Q2 2025 - WiFi & Infrastructure" \
  --body "## Feature Description

**RouterOS Path:** \`/interface/wifi/\`  
**Priority:** P0 - HIGH  
**Estimated Effort:** 3-4 weeks  
**Attributes:** ~120 total across 6 resources  
**Status:** ✅ COMPLETE - Implemented in v1.5.0

## Why Critical?

- WiFi 6 (802.11ax) is the modern wireless standard
- Completely NEW stack in RouterOS v7 (replaces \`/interface/wireless\`)
- WPA3 support
- 6 GHz band support
- Required for modern deployments

## Use Case

Modern wireless access points, WiFi 6 networks, WPA3 security, multi-SSID, guest networks.

## Implemented Resources

- [x] \`mikrotik_interface_wifi\` - WiFi interface (15 attributes)
- [x] \`mikrotik_wifi_configuration\` - WiFi profiles/configs (40+ attributes)
- [x] \`mikrotik_wifi_datapath\` - Data path settings (bridge, VLAN) (20 attributes)
- [x] \`mikrotik_wifi_security\` - Security profiles (WPA3) (20+ attributes)
- [x] \`mikrotik_wifi_channel\` - Channel configuration (10 attributes)
- [x] \`mikrotik_wifi_access_list\` - Access control (MAC filtering) (12 attributes)

## Example Configuration

\`\`\`hcl
# Security profile with WPA3
resource \"mikrotik_wifi_security\" \"wpa3_enterprise\" {
  name                    = \"wpa3-enterprise\"
  authentication_types    = [\"wpa3-eap\"]
  encryption              = \"ccmp\"
  pmf                     = \"required\"
  eap_methods             = [\"eap-tls\"]
  eap_radius_accounting   = true
  eap_radius_server       = \"192.168.1.10\"
  eap_radius_secret       = var.radius_secret
}

# Channel configuration for 5GHz WiFi 6
resource \"mikrotik_wifi_channel\" \"ch36_80mhz\" {
  name              = \"5ghz-ch36-80\"
  band              = \"5ghz-ax\"
  frequency         = 5180
  width             = 80
  secondary_frequency = 5210
}

# WiFi configuration profile
resource \"mikrotik_wifi_configuration\" \"office_5ghz\" {
  name              = \"office-5ghz\"
  mode              = \"ap\"
  ssid              = \"OfficeWiFi\"
  country           = \"poland\"
  security          = mikrotik_wifi_security.wpa3_enterprise.name
  
  # WiFi 6 specific
  he_guard_interval = \"long\"
  he_frame_format   = \"he-su\"
  
  # 802.11 parameters
  supported_rates   = [\"6Mbps\", \"12Mbps\", \"24Mbps\"]
  basic_rates       = [\"6Mbps\"]
  
  # Transmit power
  tx_power          = 20
  tx_power_mode     = \"default\"
  
  hide_ssid         = false
  disabled          = false
}

# Datapath for VLAN isolation
resource \"mikrotik_wifi_datapath\" \"vlan10\" {
  name              = \"vlan10-datapath\"
  bridge            = \"bridge1\"
  vlan_id           = 10
  vlan_mode         = \"use-tag\"
  client_isolation  = false
  arp               = \"enabled\"
}

# WiFi interface
resource \"mikrotik_wifi_interface\" \"wifi1\" {
  name          = \"wifi1\"
  configuration = mikrotik_wifi_configuration.office_5ghz.name
  datapath      = mikrotik_wifi_datapath.vlan10.name
  channel       = mikrotik_wifi_channel.ch36_80mhz.name
  mac_address   = \"AA:BB:CC:DD:EE:01\"
  disabled      = false
}

# Guest WiFi with different security
resource \"mikrotik_wifi_configuration\" \"guest\" {
  name     = \"guest-wifi\"
  mode     = \"ap\"
  ssid     = \"Guest\"
  security = mikrotik_wifi_security.wpa2_psk.name
}

resource \"mikrotik_interface_wifi\" \"wifi1_guest\" {
  name               = \"wifi1-guest\"
  master_interface   = mikrotik_interface_wifi.wifi1.name
  configuration      = mikrotik_wifi_configuration.guest.name
  datapath           = mikrotik_wifi_datapath.guest_isolated.name
}

# Access list for MAC filtering
resource \"mikrotik_wifi_access_list\" \"whitelist\" {
  mac_address   = \"11:22:33:44:55:66\"
  action        = \"accept\"
  signal_range  = \"-80..0\"
}
\`\`\`

## Key Attributes Summary

**mikrotik_interface_wifi (~30 attrs)**
**mikrotik_wifi_configuration (~40 attrs)** - SSID, security, WiFi 6 features
**mikrotik_wifi_security (~20 attrs)** - WPA3, EAP, RADIUS
**mikrotik_wifi_datapath (~20 attrs)** - Bridge, VLAN integration
**mikrotik_wifi_channel (~10 attrs)** - Band, frequency, width

## Implementation Plan

1. Study new WiFi stack architecture (major redesign)
2. Implement in order: security → channel → configuration → datapath → interface
3. Create comprehensive client API for all resources
4. Implement WPA3 support
5. Add WiFi 6 (802.11ax) specific attributes
6. Testing requires physical hardware
7. Document migration from old wireless system

## Dependencies

- Bridge/VLAN resources (already implemented)
- Interface management

## Testing Requirements

- VERY HIGH complexity
- Requires physical WiFi 6 hardware
- Multiple client device testing
- WPA3 authentication testing
- Channel switching testing
- Virtual AP testing

## Hardware Requirements

- RouterOS 7.13+ device with WiFi 6 support
- Examples: wAP ax, cAP ax, Audience

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 119)
- Coverage Matrix: ROUTEROS7_COVERAGE.md (line 60)
- MikroTik Docs: https://help.mikrotik.com/docs/display/ROS/WiFi"
```

---

#### Issue #5: Container Support

```bash
gh issue create \
  --title "[P1] Implement Container support (mikrotik_container*)" \
  --label "P1-high,routeros-v7,enhancement,area:system" \
  --milestone "Q2 2025 - WiFi & Infrastructure" \
  --body "## Feature Description

**RouterOS Path:** \`/container/\`  
**Priority:** P1 - HIGH  
**Estimated Effort:** 1-2 weeks  
**Attributes:** ~35 total

## Why Important?

- RouterOS 7.4+ supports OCI containers
- Modern application deployment on routers
- Run Pi-hole, monitoring agents, custom apps
- Enables edge computing use cases

## Use Case

Run containers on MikroTik: DNS services (Pi-hole), monitoring (Prometheus), logging, custom applications.

## Proposed Resources

- [ ] \`mikrotik_container_config\` - Global container settings
- [ ] \`mikrotik_container\` - Container instance
- [ ] \`mikrotik_container_env\` - Environment variables
- [ ] \`mikrotik_container_mount\` - Volume mounts

## Example Configuration

\`\`\`hcl
# Global container configuration
resource \"mikrotik_container_config\" \"main\" {
  registry_url = \"https://registry-1.docker.io\"
  tmpdir       = \"disk1/containers/tmp\"
  ram_high     = 256
  enabled      = true
}

# Pi-hole DNS container
resource \"mikrotik_container\" \"pihole\" {
  name         = \"pihole\"
  remote_image = \"pihole/pihole:latest\"
  interface    = \"veth-pihole\"
  root_dir     = \"disk1/containers/pihole\"
  
  mounts = [
    \"disk1/pihole/etc:/etc/pihole\",
    \"disk1/pihole/dnsmasq:/etc/dnsmasq.d\"
  ]
  
  envlist = [
    \"TZ=Europe/Warsaw\",
    \"WEBPASSWORD=admin123\",
    \"DNS1=1.1.1.1\",
    \"DNS2=8.8.8.8\"
  ]
  
  cmd        = []
  logging    = true
  disabled   = false
  
  depends_on = [mikrotik_container_config.main]
}

# Prometheus monitoring
resource \"mikrotik_container\" \"prometheus\" {
  name         = \"prometheus\"
  remote_image = \"prom/prometheus:latest\"
  interface    = \"veth-monitoring\"
  root_dir     = \"disk1/containers/prometheus\"
  
  mounts = [
    \"disk1/prometheus/config:/etc/prometheus\",
    \"disk1/prometheus/data:/prometheus\"
  ]
  
  cmd = [
    \"--config.file=/etc/prometheus/prometheus.yml\",
    \"--storage.tsdb.path=/prometheus\"
  ]
}
\`\`\`

## RouterOS CLI Example

\`\`\`
/container/config
set registry-url=https://registry-1.docker.io tmpdir=disk1/containers/tmp ram-high=256

/container
add remote-image=pihole/pihole:latest interface=veth-pihole root-dir=disk1/containers/pihole \\
    mounts=disk1/pihole/etc:/etc/pihole,disk1/pihole/dnsmasq:/etc/dnsmasq.d \\
    envlist=TZ=Europe/Warsaw,WEBPASSWORD=admin123 logging=yes
\`\`\`

## Key Attributes

**mikrotik_container_config (~8 attrs):**
- \`registry_url\`, \`tmpdir\`
- \`ram_high\`, \`ram_low\`
- \`enabled\`

**mikrotik_container (~15 attrs):**
- \`name\`, \`remote_image\`, \`tag\`, \`digest\`
- \`interface\`, \`root_dir\`
- \`cmd\`, \`entrypoint\`, \`workdir\`
- \`mounts\`, \`envlist\`
- \`hostname\`, \`domainname\`
- \`logging\`, \`disabled\`

## Implementation Plan

1. Create \`client/container.go\` with CRUD operations
2. Implement config resource (global settings)
3. Implement container resource with mount/env support
4. Add image pull validation
5. Add lifecycle management (start/stop/restart)
6. Test with common containers (Pi-hole, Prometheus)
7. Document storage requirements and limitations

## Dependencies

- veth interface support (may need implementation)

## Testing Requirements

- Medium complexity
- Requires USB storage or disk
- Docker registry connectivity
- Container lifecycle testing

## Hardware Requirements

- RouterOS 7.4+
- ARM64 or x86 device
- USB storage or internal disk
- Example devices: RB5009, CCR2004, x86 routers

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 201)
- MikroTik Docs: https://help.mikrotik.com/docs/display/ROS/Container"
```

---

### PHASE 3: Enhancements (Q3 2025)

#### Issue #6: Queue Types (CAKE, fq_codel)

```bash
gh issue create \
  --title "[P2] Add CAKE and fq_codel queue types to mikrotik_queue_type" \
  --label "P2-medium,routeros-v7,enhancement,area:system" \
  --milestone "Q3 2025 - Enhancements" \
  --body "## Feature Description

**RouterOS Path:** \`/queue/type/\`  
**Priority:** P2 - MEDIUM  
**Estimated Effort:** 1 week  
**Attributes:** ~20 new

## Why Important?

- Modern queue algorithms for better QoS
- CAKE (Common Applications Kept Enhanced) - Smart queue management
- fq_codel (Fair Queue CoDel) - Reduces bufferbloat
- Improves latency and fairness

## Proposed Changes

Extend existing \`mikrotik_queue_type\` resource with:
- [ ] CAKE queue type support
- [ ] fq_codel queue type support

## Example Configuration

\`\`\`hcl
# CAKE queue for 100 Mbps connection
resource \"mikrotik_queue_type\" \"cake_100mbit\" {
  name      = \"cake-100m\"
  kind      = \"cake\"
  
  # CAKE specific
  bandwidth = \"100M\"
  overhead  = 18        # ATM overhead
  rtt       = \"100ms\"
  nat       = true
  wash      = true
  ack_filter = false
  atm       = false
}

# fq_codel for low-latency
resource \"mikrotik_queue_type\" \"fq_codel_low_latency\" {
  name     = \"fq-codel-ll\"
  kind     = \"fq-codel\"
  
  # fq_codel specific
  target   = \"5ms\"
  interval = \"100ms\"
  quantum  = 1514
  flows    = 1024
  ecn      = true
}

# Use in simple queue
resource \"mikrotik_simple_queue\" \"customer_a\" {
  name        = \"customer_a\"
  target      = \"10.0.1.0/24\"
  max_limit   = \"100M/100M\"
  queue_type  = mikrotik_queue_type.cake_100mbit.name
}
\`\`\`

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 332)
- Coverage Matrix: ROUTEROS7_COVERAGE.md (line 206)"
```

---

#### Issue #7: ZeroTier Support

```bash
gh issue create \
  --title \"[P1] Implement ZeroTier integration (mikrotik_zerotier)\" \
  --label \"P1-high,routeros-v7,enhancement,area:vpn\" \
  --milestone \"Q3 2025 - Enhancements\" \
  --body \"## Feature Description

**RouterOS Path:** \`/zerotier/\`  
**Priority:** P1 - MEDIUM  
**Estimated Effort:** 3-5 days  
**Attributes:** ~12

## Why Important?

- Native ZeroTier client in RouterOS 7.1+
- SD-WAN and overlay networking
- Easy site-to-site connectivity
- Alternative to traditional VPN

## Example Configuration

\`\`\`hcl
resource \"mikrotik_zerotier\" \"overlay_network\" {
  name            = \"zt-overlay\"
  instance        = \"d5e5fb7e7e7e7e7e\"
  disabled        = false
  allow_default   = false
  allow_global    = false
  allow_managed   = true
  port            = 9993
  comment         = \"ZeroTier overlay for branch offices\"
}

# Use ZeroTier interface in firewall
resource \"mikrotik_firewall_filter\" \"allow_zerotier\" {
  chain           = \"input\"
  in_interface    = mikrotik_zerotier.overlay_network.name
  action          = \"accept\"
}
\`\`\`

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 260)"
```

---

#### Issue #8: VXLAN Interface

```bash
gh issue create \
  --title \"[P3] Implement VXLAN interface (mikrotik_interface_vxlan)\" \
  --label \"P3-low,routeros-v7,enhancement,area:interfaces\" \
  --milestone \"Q3 2025 - Enhancements\" \
  --body \"## Feature Description

**RouterOS Path:** \`/interface/vxlan/\`  
**Priority:** P3 - MEDIUM  
**Estimated Effort:** 1 week  
**Attributes:** ~12

## Example Configuration

\`\`\`hcl
resource \"mikrotik_interface_vxlan\" \"vxlan100\" {
  name           = \"vxlan100\"
  vni            = 100
  local_address  = \"10.0.0.1\"
  remote_address = \"10.0.0.2\"
  port           = 4789
  mtu            = 1450
  disabled       = false
}
\`\`\`

## References

- Gap Analysis: ROUTEROS7_GAP_ANALYSIS.md (line 477)"
```

---

## 🚀 Quick Start Commands

\`\`\`bash
# Clone and navigate
cd terraform-provider-mikrotik

# Create all labels
gh label create \"P0-critical\" --color \"d73a4a\"
gh label create \"P1-high\" --color \"ff9800\"
gh label create \"routeros-v7\" --color \"1d76db\"
gh label create \"enhancement\" --color \"a2eeef\"
gh label create \"area:routing\" --color \"5319e7\"
gh label create \"area:wifi\" --color \"5319e7\"

# Create milestones
gh milestone create \"Q1 2025 - Routing Foundation\" --due-date \"2025-03-31\"
gh milestone create \"Q2 2025 - WiFi & Infrastructure\" --due-date \"2025-06-30\"

# Create issues (copy from templates above)
# Each issue command is ready to run
\`\`\`

---

## 📈 Progress Tracking

Use GitHub Projects board:

\`\`\`bash
gh project create --title \"RouterOS v7 Provider Roadmap\" --owner lkolo-prez
\`\`\`

Add columns:
- 📋 Backlog
- 🔄 In Progress  
- 🧪 Testing
- ✅ Done

---

**Total Issues to Create:** 8 (P0: 4, P1: 2, P2: 1, P3: 1)  
**Estimated Timeline:** 6 months (Q1-Q3 2025)  
**Target Coverage:** 80%+ by Q4 2025
