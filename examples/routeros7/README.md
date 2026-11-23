# RouterOS 7 extensions: OSPFv3, VRRP, Wi-Fi and backup

One configuration that exercises the resources added in this fork for RouterOS 7:

| Area | Resources used |
| --- | --- |
| OSPFv3 (IPv6) | `mikrotik_ospf_instance_v7` (`version = "3"`), `mikrotik_ospf_area_v7`, `mikrotik_ospf_interface_template_v7` |
| VRRP | `mikrotik_interface_vrrp` (VRRPv3) + `mikrotik_ip_address` for the virtual IP |
| Wi-Fi (`wifi` package, 802.11ax) | `mikrotik_wifi_security`, `mikrotik_wifi_channel`, `mikrotik_wifi_datapath`, `mikrotik_wifi_configuration`, `mikrotik_interface_wifi` |
| Backup | `mikrotik_system_backup` (encrypted, taken after the other changes) |
| Plumbing | `mikrotik_bridge` |

The configuration passes `terraform validate` against a provider built from this repository.
It has not been applied to a router as part of this example; run it in a lab first.

## Requirements

- RouterOS 7 with the API-SSL service enabled (port 8729) and a user with write access.
- The `wifi` package and a Wi-Fi 6 radio for the Wi-Fi part (otherwise remove the Wi-Fi block from `main.tf`).
- Terraform >= 1.5.

## Build the provider and point Terraform at it

```bash
# from the repository root
make build                       # or: go build -o terraform-provider-mikrotik

cat > ~/.terraformrc-mikrotik <<RC
provider_installation {
  dev_overrides {
    "lkolo-prez/mikrotik" = "$(pwd)"
  }
  direct {}
}
RC
export TF_CLI_CONFIG_FILE=~/.terraformrc-mikrotik
```

## Run

```bash
cd examples/routeros7
cp terraform.tfvars.example terraform.tfvars   # edit host / user / SSID
export TF_VAR_mikrotik_password='...'
export TF_VAR_wifi_passphrase='...'
export TF_VAR_backup_password='...'

terraform validate
terraform plan
terraform apply
```

For a VRRP pair apply the same code to both routers (separate state or workspaces) with
`vrrp_priority = 200` on the primary and `vrrp_priority = 100` on the secondary.

## Files

| File | Content |
| --- | --- |
| `versions.tf` | Terraform and provider requirements, provider block |
| `variables.tf` | Connection, OSPF, VRRP, Wi-Fi and backup inputs (secrets are `sensitive`) |
| `main.tf` | The resources listed above |
| `outputs.tf` | OSPF instance, VRRP interface, Wi-Fi interface state, backup metadata |
| `terraform.tfvars.example` | Non-secret example values |
