# 📚 Documentation Index

Welcome to the **terraform-provider-mikrotik** documentation for RouterOS 7 support.

## 🚀 Quick Start

**New to this provider?** Start here:

1. **[EXTENSION_SUMMARY.md](./EXTENSION_SUMMARY.md)** - Complete overview of RouterOS 7 features
2. **[ROUTEROS7_CHEATSHEET.md](./ROUTEROS7_CHEATSHEET.md)** - RouterOS 7 command reference
3. **[examples/](./examples/)** - Terraform configuration examples

**Migrating from RouterOS 6?**
- **[MIGRATION_ROUTEROS7.md](./MIGRATION_ROUTEROS7.md)** - Migration guide

---

## 📖 Documentation

### Core Documentation

| Document | Description | For Who |
|----------|-------------|---------|
| [EXTENSION_SUMMARY.md](./EXTENSION_SUMMARY.md) | Complete project summary, new features, statistics | Everyone |
| [ROUTEROS7_CHEATSHEET.md](./ROUTEROS7_CHEATSHEET.md) | Complete RouterOS 7 command reference (~1,200 lines) | All RouterOS 7 users |
| [ROUTEROS7_COVERAGE.md](./ROUTEROS7_COVERAGE.md) | Feature coverage matrix (83 features tracked) | Planners & developers |
| [ROUTEROS7_SUPPORT.md](./ROUTEROS7_SUPPORT.md) | New RouterOS 7 resource documentation | Users |
| [MIGRATION_ROUTEROS7.md](./MIGRATION_ROUTEROS7.md) | RouterOS 6 to 7 migration guide | Migrating users |

### Examples

| Directory | Description |
|-----------|-------------|
| [examples/](./examples/) | Terraform configuration examples |
| [examples/routing/bgp-v7/](./examples/routing/bgp-v7/) | Complete BGP v7 setup example |

### Development

| Document | Description |
|----------|-------------|
| [README.md](./README.md) | Main repository README |
| [client/](./client/) | Go client library source code |

---

## 🎯 Quick Reference

### Common Questions

**Q: Where do I find all RouterOS 7 commands?**  
A: See [ROUTEROS7_CHEATSHEET.md](./ROUTEROS7_CHEATSHEET.md)

**Q: What RouterOS 7 features are available?**  
A: Check [ROUTEROS7_COVERAGE.md](./ROUTEROS7_COVERAGE.md) - 83 features tracked

**Q: How do I configure BGP in RouterOS 7?**  
A: See [examples/routing/bgp-v7/](./examples/routing/bgp-v7/)

**Q: How do I migrate from RouterOS 6?**  
A: Read [MIGRATION_ROUTEROS7.md](./MIGRATION_ROUTEROS7.md)

**Q: What's new in this provider?**  
A: See [EXTENSION_SUMMARY.md](./EXTENSION_SUMMARY.md)

---

## 📊 Project Statistics

- **Documentation**: ~4,300 lines across 7 documents
- **New Resources**: 7 types (BGP v7, VRF, veth, WiFi, CAKE queues)
- **Feature Coverage**: 28% implemented, 42% planned (83 total features)
- **Examples**: BGP v7 + more coming
- **RouterOS Support**: 7.14.3, 7.16.2, 7.17+

---

## 🆕 What's New in This Fork

### New Resource Types (7)

1. **RoutingTable** (`/routing/table`) - VRF support
2. **RoutingRule** (`/routing/rule`) - Policy-based routing
3. **VRF** (`/ip/vrf`) - Virtual routing and forwarding
4. **InterfaceVeth** (`/interface/veth`) - Container networking
5. **WiFiRadio** (`/interface/wifi/radio`) - 802.11ax support
6. **WiFiConfiguration** (`/interface/wifi/configuration`) - WiFi profiles
7. **QueueType** (`/queue/type`) - CAKE & fq_codel support

### Enhanced Documentation

- ✅ Complete RouterOS 7 cheat sheet (~1,200 lines)
- ✅ Feature coverage matrix (83 features)
- ✅ Production-ready BGP v7 example
- ✅ Extended CI/CD with daily integration tests

---

## 🚀 Getting Started

### 1. Check the Cheat Sheet

```bash
# View all RouterOS 7 commands
cat ROUTEROS7_CHEATSHEET.md
```

### 2. Review Feature Coverage

```bash
# See what's implemented
cat ROUTEROS7_COVERAGE.md
```

### 3. Try the BGP Example

```bash
cd examples/routing/bgp-v7
cat README.md
```

### 4. Build the Provider

```bash
go mod tidy
go build -o terraform-provider-mikrotik.exe
```

---

## 📂 Repository Structure

```
terraform-provider-mikrotik/
│
├── 📘 Documentation
│   ├── EXTENSION_SUMMARY.md       # Project overview
│   ├── ROUTEROS7_CHEATSHEET.md    # Command reference
│   ├── ROUTEROS7_COVERAGE.md      # Feature matrix
│   ├── ROUTEROS7_SUPPORT.md       # Resource docs
│   └── MIGRATION_ROUTEROS7.md     # Migration guide
│
├── 📁 Examples
│   ├── README.md
│   └── routing/bgp-v7/            # BGP v7 example
│
├── 💻 Source Code
│   ├── client/
│   │   ├── routing_v7.go          # VRF & Routing
│   │   ├── advanced_v7.go         # WiFi & Queues
│   │   ├── bgp_connection.go      # BGP v7
│   │   └── [29 more files...]
│   └── [provider code...]
│
└── 🔄 CI/CD
    └── .github/workflows/
        ├── continuous-integration.yml
        └── integration-tests.yml
```

---

## 🔗 Important Links

- **GitHub Repository**: [lkolo-prez/terraform-provider-mikrotik](https://github.com/lkolo-prez/terraform-provider-mikrotik)
- **Original Provider**: [ddelnano/terraform-provider-mikrotik](https://github.com/ddelnano/terraform-provider-mikrotik)
- **MikroTik Docs**: [RouterOS v7 Documentation](https://help.mikrotik.com/docs/spaces/ROS/pages/115736772/Upgrading+to+v7)
- **Cheat Sheet Source**: [3zzy's RouterOS v7 Gist](https://gist.github.com/3zzy/61e356f0bfcd2918d271836e30d80698)

---

## ✅ Status

**Last Updated**: November 25, 2025  
**Status**: ✅ Production Ready  
**RouterOS Support**: 7.14.3, 7.16.2, 7.17+  
**Documentation**: 100% Complete  
**CI/CD**: Active ([View Actions](https://github.com/lkolo-prez/terraform-provider-mikrotik/actions))

---

**Need help?** Open an issue on [GitHub](https://github.com/lkolo-prez/terraform-provider-mikrotik/issues)
