# Development Roadmap - Terraform Provider MikroTik RouterOS 7

**Repository**: https://github.com/lkolo-prez/terraform-provider-mikrotik  
**Current Status**: Phase 1 - Foundation Complete  
**Last Updated**: November 25, 2025

---

## 🎯 Project Vision

Create a **comprehensive Terraform provider** for MikroTik RouterOS 7.x with:
- ✅ Full RouterOS 7 API support
- ✅ 100% feature coverage of RouterOS 7 capabilities
- ✅ Production-ready CI/CD pipeline
- ✅ Extensive documentation and examples
- ✅ Active community support

---

## 📊 Current Status (Phase 1 Complete)

### ✅ Foundation (DONE)
- [x] Updated to go-routeros v3.0.1 (latest)
- [x] Migrated to Go 1.21+ (required for log/slog)
- [x] Fixed all dependency issues
- [x] Created English-only documentation
- [x] Basic CI/CD pipeline (build & lint)
- [x] Project structure organized

### 📈 Coverage Statistics
- **Fully Implemented**: 23 resources (~28%)
- **Partially Implemented**: 15 resources (~18%)
- **Planned**: 35 resources (~42%)
- **Not Planned**: 10 resources (~12%)
- **Total Features**: 83 tracked

---

## 🚀 Development Phases

### **Phase 1: Foundation & Stabilization** ✅ COMPLETE

**Goal**: Get project building without errors, fix dependencies, organize documentation

**Tasks Completed**:
- [x] Fix go-routeros version (v3.3.0 → v3.0.1)
- [x] Update Go requirement (1.18 → 1.21+)
- [x] Remove Polish documentation
- [x] Create English documentation index
- [x] Simplify CI/CD to basic build + lint
- [x] Update README.md with correct requirements
- [x] Verify all imports and syntax

**Deliverables**:
- ✅ Project builds successfully
- ✅ No compilation errors
- ✅ English-only documentation
- ✅ Basic CI/CD pipeline passing

---

### **Phase 2: Core Routing & BGP** 🔄 IN PROGRESS

**Goal**: Implement and test core RouterOS 7 routing features

**Priority**: HIGH  
**Estimated Time**: 2-3 weeks

#### Tasks:

**2.1 BGP Implementation** (Priority: CRITICAL)
- [ ] Test existing BGP resources:
  - `mikrotik_bgp_connection` (new RouterOS 7)
  - `mikrotik_bgp_template` (new RouterOS 7)
  - `mikrotik_bgp_instance` (legacy, deprecated)
  - `mikrotik_bgp_peer` (legacy, deprecated)
- [ ] Create comprehensive BGP test suite
- [ ] Document BGP migration from v6 to v7
- [ ] Create BGP examples:
  - [ ] Simple BGP peering
  - [ ] BGP with route filtering
  - [ ] Multi-ISP BGP setup
  - [ ] BGP with BFD
- [ ] Validate against real RouterOS 7 devices

**2.2 Routing Tables & VRFs**
- [ ] Implement `mikrotik_routing_table` resource
- [ ] Implement `mikrotik_routing_rule` resource
- [ ] Implement `mikrotik_vrf` resource
- [ ] Create VRF examples
- [ ] Test routing table isolation
- [ ] Document VRF best practices

**2.3 Static & Dynamic Routing**
- [ ] Implement `mikrotik_ip_route` resource
- [ ] Implement `mikrotik_ipv6_route` resource
- [ ] Add OSPF support (if not exists)
- [ ] Create routing examples
- [ ] Test route redistribution

**Deliverables**:
- [ ] All routing resources tested
- [ ] BGP v7 fully documented
- [ ] 5+ routing examples
- [ ] Integration tests passing

---

### **Phase 3: Firewall & Security** 📋 PLANNED

**Goal**: Complete firewall implementation for RouterOS 7

**Priority**: HIGH  
**Estimated Time**: 2-3 weeks

#### Tasks:

**3.1 Firewall Rules**
- [ ] Test existing resources:
  - `mikrotik_firewall_filter`
  - `mikrotik_firewall_nat`
  - `mikrotik_firewall_mangle`
  - `mikrotik_firewall_raw` (new in v7)
- [ ] Implement missing firewall features:
  - [ ] Connection tracking helpers
  - [ ] Fast track rules
  - [ ] Address lists management
  - [ ] Port knocking
- [ ] Create firewall examples:
  - [ ] Basic edge firewall
  - [ ] Site-to-site VPN firewall
  - [ ] DDoS protection rules
  - [ ] FastTrack setup

**3.2 IPsec & VPN**
- [ ] Implement `mikrotik_ipsec_peer` resource
- [ ] Implement `mikrotik_ipsec_proposal` resource
- [ ] Implement `mikrotik_ipsec_policy` resource
- [ ] Test WireGuard resources (existing)
- [ ] Create VPN examples:
  - [ ] Site-to-site IPsec
  - [ ] WireGuard tunnel setup
  - [ ] Road warrior VPN

**3.3 Security Features**
- [ ] Implement SSH key management
- [ ] Implement certificate management
- [ ] Implement user/group management
- [ ] Create security hardening examples

**Deliverables**:
- [ ] All firewall resources tested
- [ ] 10+ firewall examples
- [ ] VPN examples complete
- [ ] Security best practices documented

---

### **Phase 4: Interfaces & Networking** 📋 PLANNED

**Goal**: Complete interface management for RouterOS 7

**Priority**: MEDIUM  
**Estimated Time**: 2-3 weeks

#### Tasks:

**4.1 Physical Interfaces**
- [ ] Test existing interface resources
- [ ] Implement interface bonding
- [ ] Implement interface bridging (enhanced)
- [ ] Create interface examples

**4.2 VLAN & Virtual Interfaces**
- [ ] Test `mikrotik_interface_vlan7` (new)
- [ ] Implement veth interfaces
- [ ] Implement VXLAN interfaces
- [ ] Create VLAN examples:
  - [ ] Tagged VLANs
  - [ ] VLAN bridging
  - [ ] Inter-VLAN routing

**4.3 WiFi 6/6E Support**
- [ ] Implement `mikrotik_wifi_radio` resource
- [ ] Implement `mikrotik_wifi_configuration` resource
- [ ] Implement `mikrotik_wifi_security` resource
- [ ] Create WiFi examples:
  - [ ] Basic AP setup
  - [ ] Multi-SSID configuration
  - [ ] Guest network isolation
  - [ ] 802.11ax features

**4.4 Advanced Interfaces**
- [ ] Implement L2TP interfaces
- [ ] Implement PPPoE interfaces
- [ ] Implement GRE tunnels
- [ ] Test container interfaces (veth)

**Deliverables**:
- [ ] All interface types supported
- [ ] WiFi 6 fully implemented
- [ ] 15+ interface examples
- [ ] Performance testing complete

---

### **Phase 5: QoS & Traffic Management** 📋 PLANNED

**Goal**: Implement modern QoS features

**Priority**: MEDIUM  
**Estimated Time**: 1-2 weeks

#### Tasks:

**5.1 Queue Systems**
- [ ] Test existing queue resources
- [ ] Implement CAKE queues (new in v7)
- [ ] Implement fq_codel queues
- [ ] Create QoS examples:
  - [ ] SQM (Smart Queue Management)
  - [ ] Per-client bandwidth limits
  - [ ] Priority queuing

**5.2 Traffic Shaping**
- [ ] Implement queue trees
- [ ] Implement simple queues
- [ ] Create traffic shaping examples
- [ ] Document QoS best practices

**Deliverables**:
- [ ] Modern QoS features complete
- [ ] 5+ QoS examples
- [ ] Performance benchmarks

---

### **Phase 6: System Management** 📋 PLANNED

**Goal**: Complete system administration features

**Priority**: LOW-MEDIUM  
**Estimated Time**: 1-2 weeks

#### Tasks:

**6.1 System Configuration**
- [ ] Test existing resources:
  - `mikrotik_script`
  - `mikrotik_scheduler`
  - DNS, DHCP, etc.
- [ ] Implement system identity
- [ ] Implement NTP client
- [ ] Implement logging
- [ ] Implement backup/restore

**6.2 Monitoring & Diagnostics**
- [ ] Implement SNMP configuration
- [ ] Implement Netflow/IPFIX
- [ ] Implement health monitoring
- [ ] Create monitoring examples

**6.3 Container Support**
- [ ] Test veth interfaces (from Phase 4)
- [ ] Implement container management
- [ ] Create container examples

**Deliverables**:
- [ ] System management complete
- [ ] Monitoring examples
- [ ] Container examples

---

### **Phase 7: Testing & Quality Assurance** 📋 PLANNED

**Goal**: Comprehensive testing and validation

**Priority**: CRITICAL  
**Estimated Time**: 2-3 weeks

#### Tasks:

**7.1 Unit Tests**
- [ ] Write unit tests for all resources
- [ ] Achieve 80%+ code coverage
- [ ] Setup coverage reporting

**7.2 Integration Tests**
- [ ] Enable full integration test suite
- [ ] Test against multiple RouterOS versions:
  - [ ] RouterOS 7.14.3 (stable)
  - [ ] RouterOS 7.16.2 (stable)
  - [ ] RouterOS 7.17+ (latest)
- [ ] Setup automated nightly tests
- [ ] Create compatibility matrix

**7.3 End-to-End Tests**
- [ ] Create real-world scenario tests:
  - [ ] Full edge router setup
  - [ ] Site-to-site VPN deployment
  - [ ] Multi-site BGP network
  - [ ] WiFi deployment
- [ ] Performance testing
- [ ] Stress testing

**7.4 Documentation Testing**
- [ ] Validate all examples work
- [ ] Test migration guides
- [ ] Verify troubleshooting steps

**Deliverables**:
- [ ] 80%+ code coverage
- [ ] All integration tests passing
- [ ] Performance benchmarks published
- [ ] Test reports automated

---

### **Phase 8: Documentation & Examples** 📋 PLANNED

**Goal**: Production-ready documentation

**Priority**: HIGH  
**Estimated Time**: 2 weeks

#### Tasks:

**8.1 Resource Documentation**
- [ ] Complete API documentation for all resources
- [ ] Add code examples to each resource
- [ ] Document all attributes and arguments
- [ ] Create troubleshooting guides

**8.2 Tutorial Content**
- [ ] Create getting started guide
- [ ] Write tutorial series:
  - [ ] Basic router setup
  - [ ] Advanced routing with BGP
  - [ ] Firewall configuration
  - [ ] WiFi deployment
  - [ ] Site-to-site VPN
- [ ] Create video tutorials (optional)

**8.3 Reference Documentation**
- [ ] Update ROUTEROS7_CHEATSHEET.md
- [ ] Update ROUTEROS7_COVERAGE.md
- [ ] Create API reference
- [ ] Document best practices

**8.4 Example Collection**
- [ ] Create example repository structure:
  ```
  examples/
  ├── basic/           # Simple configurations
  ├── routing/         # BGP, OSPF, static routes
  ├── firewall/        # Security configurations
  ├── vpn/             # IPsec, WireGuard
  ├── wifi/            # Wireless configurations
  ├── qos/             # Traffic management
  └── real-world/      # Complete deployments
  ```
- [ ] Validate all examples
- [ ] Add README to each example

**Deliverables**:
- [ ] Complete API documentation
- [ ] 50+ working examples
- [ ] Tutorial series complete
- [ ] Documentation searchable

---

### **Phase 9: Community & Release** 📋 PLANNED

**Goal**: Public release and community building

**Priority**: HIGH  
**Estimated Time**: 1-2 weeks

#### Tasks:

**9.1 Release Preparation**
- [ ] Version tagging (v1.0.0)
- [ ] Changelog creation
- [ ] Release notes
- [ ] Migration guides finalized

**9.2 Community Setup**
- [ ] Create contributing guidelines
- [ ] Setup issue templates
- [ ] Create PR templates
- [ ] Setup discussions
- [ ] Create code of conduct

**9.3 Publishing**
- [ ] Publish to Terraform Registry
- [ ] Announce on MikroTik forums
- [ ] Blog post/announcement
- [ ] Social media promotion

**9.4 Support Infrastructure**
- [ ] Setup issue triage process
- [ ] Create FAQ
- [ ] Setup Discord/Slack channel
- [ ] Create support documentation

**Deliverables**:
- [ ] v1.0.0 released
- [ ] Published to Terraform Registry
- [ ] Community active
- [ ] Support channels operational

---

### **Phase 10: Advanced Features** 📋 FUTURE

**Goal**: Cutting-edge RouterOS 7 features

**Priority**: LOW  
**Estimated Time**: Ongoing

#### Potential Features:
- [ ] MPLS support
- [ ] SD-WAN capabilities
- [ ] Advanced monitoring (Prometheus exporter)
- [ ] REST API integration
- [ ] GraphQL support
- [ ] Automated backups
- [ ] Configuration drift detection
- [ ] Multi-device orchestration
- [ ] Ansible integration
- [ ] Kubernetes operator

---

## 🎯 Success Metrics

### Coverage Goals
- **Phase 2 End**: 40% coverage (routing complete)
- **Phase 4 End**: 60% coverage (networking complete)
- **Phase 6 End**: 80% coverage (system complete)
- **Phase 9 End**: 95% coverage (production release)

### Quality Metrics
- Code coverage: >80%
- Integration test pass rate: >95%
- Documentation completeness: 100%
- Example validation: 100%

### Community Metrics
- GitHub stars: 100+
- Active contributors: 5+
- Issues resolved: <1 week average
- Community size: 50+ active users

---

## 📅 Timeline Estimate

| Phase | Duration | Target Completion |
|-------|----------|-------------------|
| Phase 1: Foundation | 1 week | ✅ Complete |
| Phase 2: Routing & BGP | 2-3 weeks | Week 4 |
| Phase 3: Firewall | 2-3 weeks | Week 7 |
| Phase 4: Interfaces | 2-3 weeks | Week 10 |
| Phase 5: QoS | 1-2 weeks | Week 12 |
| Phase 6: System | 1-2 weeks | Week 14 |
| Phase 7: Testing | 2-3 weeks | Week 17 |
| Phase 8: Documentation | 2 weeks | Week 19 |
| Phase 9: Release | 1-2 weeks | Week 21 |

**Total Estimated Time**: ~5 months to v1.0.0 release

---

## 🔄 Iterative Development Process

For each phase:

1. **Plan** (1-2 days)
   - Review ROUTEROS7_COVERAGE.md
   - Prioritize features
   - Create detailed task list

2. **Implement** (1-2 weeks)
   - Write code
   - Add unit tests
   - Update documentation

3. **Test** (2-3 days)
   - Run integration tests
   - Validate examples
   - Manual testing

4. **Review** (1 day)
   - Code review
   - Documentation review
   - Update coverage matrix

5. **Merge & Deploy** (1 day)
   - Merge to main
   - Tag release
   - Update changelog

---

## 🚀 Quick Start - Phase 2

To begin Phase 2 (Routing & BGP):

```bash
# 1. Create feature branch
git checkout -b phase2-routing-bgp

# 2. Review existing BGP code
cat client/bgp_connection.go
cat client/bgp_template.go

# 3. Create test environment
make routeros ROUTEROS_VERSION=7.16.2

# 4. Run existing tests
export MIKROTIK_HOST=127.0.0.1:8728
export MIKROTIK_USER=admin
export MIKROTIK_PASSWORD=""
cd client
go test -v -run TestBgp

# 5. Start implementing missing features
```

---

## 📚 References

- [RouterOS 7 Documentation](https://help.mikrotik.com/docs/spaces/ROS/overview)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin)
- [ROUTEROS7_COVERAGE.md](./ROUTEROS7_COVERAGE.md) - Feature tracking
- [ROUTEROS7_CHEATSHEET.md](./ROUTEROS7_CHEATSHEET.md) - Command reference
- [INDEX.md](./INDEX.md) - Documentation index

---

## 🤝 Contributing

This is an open roadmap. Contributions are welcome!

1. Pick a task from current phase
2. Create feature branch
3. Implement + test + document
4. Submit pull request
5. Update ROUTEROS7_COVERAGE.md

---

**Last Updated**: November 25, 2025  
**Current Phase**: Phase 1 Complete ✅ → Phase 2 Starting 🚀  
**Next Milestone**: BGP v7 fully tested and documented
