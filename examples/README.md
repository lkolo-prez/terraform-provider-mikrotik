# Terraform Provider Examples for MikroTik RouterOS v7

This directory contains Terraform configurations for the resources implemented in this
provider. Every file and directory listed below exists in the repository.

## Directory structure

```
examples/
├── README.md                     # This file
├── provider/provider.tf          # Provider configuration
├── routeros7/                    # RouterOS 7 extensions: OSPFv3, VRRP, Wi-Fi, backup (validated)
├── resources/<resource>/         # One minimal example per resource (resource.tf, import.sh)
├── bgp/                          # BGP scenarios (01-06) + MIGRATION.md
│   ├── 01-ebgp-peering.tf
│   ├── 02-ibgp-full-mesh.tf
│   ├── 03-route-reflector.tf
│   ├── 04-vpn-mpls.tf
│   ├── 05-communities-filtering.tf
│   └── 06-graceful-restart.tf
├── routing/bgp-v7/               # BGP with the RouterOS 7 connection/template API
├── ospf_basic.tf                 # Single-area OSPF (RouterOS 7 instance/area/template)
├── ospf_enterprise.tf            # Multi-area OSPF
├── routing_filter_basic.tf       # Routing filter chains and rules
├── routing_filter_advanced.tf
├── routing_table_basic.tf        # Routing tables
├── routing_table_vrf.tf          # VRF routing tables
└── container_pihole.tf           # RouterOS container (Pi-hole)
```

### RouterOS 7 extensions

| Feature | Resources | Example |
| --- | --- | --- |
| OSPF v2/v3 | `mikrotik_ospf_instance_v7`, `mikrotik_ospf_area_v7`, `mikrotik_ospf_interface_template_v7` | [`routeros7/`](routeros7/), [`ospf_basic.tf`](ospf_basic.tf), [`ospf_enterprise.tf`](ospf_enterprise.tf) |
| VRRP | `mikrotik_interface_vrrp` | [`routeros7/`](routeros7/), [`resources/mikrotik_interface_vrrp`](resources/mikrotik_interface_vrrp/) |
| Wi-Fi (`wifi` package) | `mikrotik_wifi_security`, `mikrotik_wifi_channel`, `mikrotik_wifi_datapath`, `mikrotik_wifi_configuration`, `mikrotik_wifi_access_list`, `mikrotik_interface_wifi` | [`routeros7/`](routeros7/) |
| Backup | `mikrotik_system_backup` | [`routeros7/`](routeros7/), [`resources/mikrotik_system_backup`](resources/mikrotik_system_backup/) |
| Containers | `mikrotik_container`, `mikrotik_container_config` | [`container_pihole.tf`](container_pihole.tf) |
| Routing filters | `mikrotik_routing_filter_chain`, `mikrotik_routing_filter_rule` | [`routing_filter_basic.tf`](routing_filter_basic.tf), [`routing_filter_advanced.tf`](routing_filter_advanced.tf) |
| Routing tables / VRF | `mikrotik_routing_table` | [`routing_table_basic.tf`](routing_table_basic.tf), [`routing_table_vrf.tf`](routing_table_vrf.tf) |
| BGP (v7 API) | `mikrotik_bgp_connection`, `mikrotik_bgp_template`, `mikrotik_bgp_instance_v7`, data source `mikrotik_bgp_session` | [`bgp/`](bgp/), [`routing/bgp-v7/`](routing/bgp-v7/) |

## Prerequisites

1. Terraform 1.x.
2. The provider built from this repository (the RouterOS 7 resources are not part of the
   upstream `ddelnano/mikrotik` releases):

   ```bash
   make build      # or: go build -o terraform-provider-mikrotik
   ```

   and a `dev_overrides` CLI configuration pointing at the binary, as shown in
   [`routeros7/README.md`](routeros7/README.md).
3. RouterOS 7 with API access enabled and a user with write permissions.

Some older single-file examples declare `source = "ddelnano/mikrotik"`. Adjust the source
address to the one your `dev_overrides` block uses.

## Quick start

```bash
cd examples/routeros7
cp terraform.tfvars.example terraform.tfvars
export TF_VAR_mikrotik_password='...'
terraform init      # not needed with dev_overrides, harmless otherwise
terraform plan
terraform apply
terraform destroy   # when done
```

**Security note:** never commit `terraform.tfvars` with real credentials. Pass secrets as
`TF_VAR_*` environment variables or from a secret manager.

## Best practices

- Test in a lab router before production; use separate workspaces or state files per device.
- Use remote state with locking for shared devices.
- Keep one router (or one VRRP member) per state to avoid accidental cross-device changes.

## Troubleshooting

```bash
export TF_LOG=DEBUG
terraform apply
```

1. **"Resource not found"**: check that the RouterOS version and installed packages support the feature (e.g. `wifi` package for `mikrotik_wifi_*`).
2. **"Permission denied"**: the RouterOS user needs write rights.
3. **"Invalid parameter"**: check the RouterOS 7 documentation for parameter changes.

## Migration from v6

1. Review [MIGRATION_ROUTEROS7.md](../MIGRATION_ROUTEROS7.md).
2. Check [ROUTEROS7_COVERAGE.md](../ROUTEROS7_COVERAGE.md) for feature availability.
3. For BGP, see [bgp/MIGRATION.md](bgp/MIGRATION.md) (`mikrotik_bgp_instance`/`mikrotik_bgp_peer` → `mikrotik_bgp_connection`/`mikrotik_bgp_template`).

## Support

- [Provider documentation](../README.md)
- [RouterOS v7 cheat sheet](../ROUTEROS7_CHEATSHEET.md)
- [Feature coverage matrix](../ROUTEROS7_COVERAGE.md)
- [Issues for this fork](https://github.com/lkolo-prez/terraform-provider-mikrotik/issues)
