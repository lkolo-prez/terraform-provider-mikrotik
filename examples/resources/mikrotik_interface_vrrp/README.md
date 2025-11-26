# VRRP High Availability Example

This example demonstrates how to configure VRRP (Virtual Router Redundancy Protocol) for high availability using two MikroTik routers.

## Scenario

- **Master Router**: 192.168.1.10 (priority 254)
- **Backup Router**: 192.168.1.11 (priority 100)
- **Virtual IP**: 192.168.1.1 (gateway for clients)
- **VRID**: 10

## Features Demonstrated

- ✅ VRRP interface configuration
- ✅ Master/Backup router setup
- ✅ Authentication (simple password)
- ✅ Virtual IP assignment
- ✅ State change scripts (optional)
- ✅ Multiple VRRP groups for load balancing

## Prerequisites

- 2x MikroTik RouterOS 7.x devices
- Network connectivity between routers
- Terraform 1.0+ installed

## Usage

1. **Configure variables**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your passwords
   ```

2. **Initialize Terraform**:
   ```bash
   terraform init
   ```

3. **Review plan**:
   ```bash
   terraform plan
   ```

4. **Apply configuration**:
   ```bash
   terraform apply
   ```

## Verification

### Check VRRP Status on Master
```bash
/interface vrrp print detail
/interface vrrp monitor [find name=vrrp-gateway]
```

### Expected Output (Master)
```
name="vrrp-gateway" running=yes master-router=<local>
```

### Expected Output (Backup)
```
name="vrrp-gateway" running=yes master-router=192.168.1.10
```

### Test Failover
1. Disable master VRRP interface:
   ```bash
   /interface vrrp set vrrp-gateway disabled=yes
   ```

2. Check backup becomes master:
   ```bash
   /interface vrrp monitor vrrp-gateway
   ```

3. Re-enable master (with preemption):
   ```bash
   /interface vrrp set vrrp-gateway disabled=no
   ```

## Load Balancing with Multiple VRRP Groups

This example also shows how to use multiple VRRP groups for load balancing:

- **VRRP 10**: Router A is master, Router B is backup
- **VRRP 20**: Router B is master, Router A is backup

This distributes traffic across both routers while maintaining redundancy.

## Important Notes

- **Same VRID**: Both routers must use the same VRID (10 in this example)
- **Authentication**: Use strong passwords in production
- **Preemption**: Enabled by default - higher priority router will become master
- **Firewall**: Don't block VRRP (IP protocol 112)
- **Timing**: Default interval is 1 second

## Troubleshooting

### Both routers claim to be master
- Check network connectivity
- Verify VRRP packets are not blocked by firewall
- Check for network loops

### VRRP not becoming master
- Verify priority settings (higher = master)
- Check authentication passwords match
- Ensure physical interface is up
- Verify VRID matches on both routers

### Logs
```bash
/log print where topics~"vrrp"
```

## Security Recommendations

1. ✅ Always use authentication in production
2. ✅ Use strong passwords (min 16 characters)
3. ✅ Consider VRRPv3 with IPsec for encryption
4. ✅ Monitor VRRP state changes
5. ✅ Use scripts for alerting on state transitions

## References

- [MikroTik VRRP Documentation](https://help.mikrotik.com/docs/display/ROS/VRRP)
- [RFC 3768 - VRRPv2](https://tools.ietf.org/html/rfc3768)
- [RFC 5798 - VRRPv3](https://tools.ietf.org/html/rfc5798)
