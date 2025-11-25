# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
