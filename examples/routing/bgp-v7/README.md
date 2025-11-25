# BGP Configuration with RouterOS v7 API

This example demonstrates BGP setup using the new RouterOS v7 connection/template system.

## What's New in v7?

RouterOS v7 **completely redesigned** BGP:
- ✅ `/routing/bgp/connection` replaces `/routing/bgp/instance` and `/routing/bgp/peer`
- ✅ `/routing/bgp/template` for reusable configurations
- ✅ Better multi-protocol support (VPNv4, VPNv6, MPLS)
- ✅ Improved session management

## Architecture

```
[Your Router] ---- eBGP Session ---- [ISP Router]
     AS 65001                            AS 65000
     10.0.0.1                            10.0.0.2
```

## Resources Created

1. BGP Template (reusable configuration)
2. BGP Connection to ISP
3. Firewall rules to allow BGP (TCP/179)
4. Address list for BGP neighbors

## Usage

### 1. Configure Variables

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
```

### 2. Apply Configuration

```bash
terraform init
terraform plan
terraform apply
```

### 3. Verify BGP

```routeros
/routing/bgp/session/print
/routing/route/print where bgp=yes
```

## Prerequisites

- RouterOS 7.14.3+
- WAN interface configured
- IP connectivity to BGP peer

## Security Notes

- BGP uses TCP port 179
- Consider MD5 authentication
- Use prefix filters (see advanced example)

## Advanced Options

See `main.tf` for:
- BFD support
- Route filtering
- Multi-protocol support
- Connection templates

## References

- [RouterOS BGP Documentation](https://help.mikrotik.com/docs/display/ROS/BGP)
- [BGP Migration Guide](../../../MIGRATION_ROUTEROS7.md#bgp-changes)
