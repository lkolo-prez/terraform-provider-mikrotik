# Firewall NAT Resource

Creates and manages a firewall NAT rule for network address translation.

## Features

- **Masquerade**: Internet sharing with automatic source IP
- **Destination NAT**: Port forwarding to internal servers
- **Source NAT**: Static source IP translation
- **Netmap**: 1:1 NAT mapping for DMZ
- **Redirect**: Transparent proxy redirection
- **Connection Tracking**: Match by connection state, marks
- **Advanced Matching**: Time-based, rate-limited, conditional NAT
- **Logging**: Per-rule logging with custom prefixes

## Use Cases

- Internet sharing (masquerade)
- Port forwarding (dst-nat)
- Load balancing across multiple WANs
- DMZ hosting with 1:1 NAT
- Transparent proxy setup
- VPN split tunneling
- Hairpin NAT for internal access

## Example Usage

### Basic Internet Sharing (Masquerade)

```terraform
resource "mikrotik_firewall_nat" "masquerade_wan" {
  chain         = "srcnat"
  action        = "masquerade"
  out_interface = "ether1-wan"
  comment       = "Internet sharing"
}
```

### Port Forwarding (HTTP Server)

```terraform
resource "mikrotik_firewall_nat" "http_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "80"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.100"
  to_ports     = "80"
  comment      = "Forward HTTP to web server"
}
```

### SSH with Custom Port

```terraform
resource "mikrotik_firewall_nat" "ssh_custom" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "2222"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.50"
  to_ports     = "22"
  comment      = "SSH on port 2222"
}
```

### Hairpin NAT

```terraform
resource "mikrotik_firewall_nat" "hairpin" {
  chain       = "srcnat"
  action      = "masquerade"
  src_address = "192.168.1.0/24"
  dst_address = "192.168.1.100"
  protocol    = "tcp"
  dst_port    = "80"
  comment     = "Internal access to public IP"
}
```

### Multi-WAN NAT

```terraform
resource "mikrotik_firewall_nat" "wan1" {
  chain           = "srcnat"
  action          = "masquerade"
  connection_mark = "wan1-conn"
  out_interface   = "ether1"
  comment         = "WAN1 masquerade"
}

resource "mikrotik_firewall_nat" "wan2" {
  chain           = "srcnat"
  action          = "masquerade"
  connection_mark = "wan2-conn"
  out_interface   = "ether2"
  comment         = "WAN2 masquerade"
}
```

### NAT Exception (No NAT for VPN)

```terraform
resource "mikrotik_firewall_nat" "vpn_no_nat" {
  chain            = "srcnat"
  action           = "accept"
  src_address_list = "vpn-clients"
  dst_address_list = "local-networks"
  comment          = "No NAT for VPN traffic"
}
```

## Argument Reference

### Required Arguments

- `chain` - (Required) NAT chain: `srcnat` (source NAT) or `dstnat` (destination NAT).
- `action` - (Required) NAT action: `masquerade`, `dst-nat`, `src-nat`, `netmap`, `redirect`, `accept`, `passthrough`, `jump`, `return`, `same`.

### Optional Arguments - Basic

- `disabled` - (Optional) Whether the rule is disabled. Default: `false`.
- `comment` - (Optional) Comment for the rule.

### Optional Arguments - Matching Criteria

- `src_address` - (Optional) Source IP/network (CIDR). Example: `192.168.1.0/24`.
- `dst_address` - (Optional) Destination IP/network (CIDR).
- `src_address_list` - (Optional) Source address list name.
- `dst_address_list` - (Optional) Destination address list name.
- `protocol` - (Optional) IP protocol: `tcp`, `udp`, `icmp`, or number.
- `src_port` - (Optional) Source port(s). Example: `80` or `80-90` or `80,443`.
- `dst_port` - (Optional) Destination port(s).
- `in_interface` - (Optional) Input interface name.
- `out_interface` - (Optional) Output interface name.
- `in_interface_list` - (Optional) Input interface list name.
- `out_interface_list` - (Optional) Output interface list name.

### Optional Arguments - Connection Tracking

- `connection_state` - (Optional) Connection state. Example: `new,established,related`.
- `connection_nat_state` - (Optional) Connection NAT state: `srcnat`, `dstnat`.
- `connection_mark` - (Optional) Connection mark to match.
- `packet_mark` - (Optional) Packet mark to match.
- `routing_mark` - (Optional) Routing mark to match.

### Optional Arguments - NAT Parameters

- `to_addresses` - (Optional) NAT to IP address(es). Required for `dst-nat` and `src-nat`. Example: `192.168.1.100` or `192.168.1.1-192.168.1.10`.
- `to_ports` - (Optional) NAT to port(s). Example: `8080` or `8000-9000`.

### Optional Arguments - Logging

- `log` - (Optional) Enable logging. Default: `false`.
- `log_prefix` - (Optional) Log message prefix.

### Optional Arguments - Advanced

- `icmp_options` - (Optional) ICMP options. Example: `0:0`.
- `limit` - (Optional) Rate limit. Example: `5/1m,10:packet`.
- `time` - (Optional) Time range. Example: `08:00-17:00,mon,tue,wed,thu,fri`.
- `random` - (Optional) Random probability (1-99).
- `hotspot` - (Optional) HotSpot status: `none`, `auth`, `http`, `https`.
- `content_type` - (Optional) HTTP content type.
- `layer7_protocol` - (Optional) Layer 7 protocol name.
- `tcp_flags` - (Optional) TCP flags. Example: `syn,!ack`.
- `tcp_mss` - (Optional) TCP MSS range. Example: `500-1500`.
- `dst_limit` - (Optional) Destination limit.
- `packet_size` - (Optional) Packet size range. Example: `64-128`.
- `src_address_type` - (Optional) Source address type.
- `dst_address_type` - (Optional) Destination address type.

### Computed Attributes

- `id` - The unique identifier of the NAT rule.
- `bytes` - Number of bytes matched (computed).
- `packets` - Number of packets matched (computed).
- `dynamic` - Whether rule is dynamic (computed).
- `invalid` - Whether rule is invalid (computed).

## Import

Firewall NAT rules can be imported using the ID:

```bash
terraform import mikrotik_firewall_nat.example *1
```

## NAT Actions

### masquerade
Automatic source NAT with interface IP. Best for dynamic WAN IPs.

```terraform
action = "masquerade"
# No to_addresses needed
```

### dst-nat
Destination NAT for port forwarding.

```terraform
action       = "dst-nat"
to_addresses = "192.168.1.100"
to_ports     = "80"  # Optional port change
```

### src-nat
Source NAT with static IP. For static WAN IPs.

```terraform
action       = "src-nat"
to_addresses = "203.0.113.10"
```

### netmap
1:1 NAT mapping for DMZ.

```terraform
action       = "netmap"
to_addresses = "192.168.100.0/24"
```

### redirect
Transparent proxy redirection.

```terraform
action   = "redirect"
to_ports = "3128"  # Proxy port
```

### accept
Accept packet without NAT (NAT exception).

```terraform
action = "accept"
```

## Common Patterns

### Internet Sharing
```terraform
chain         = "srcnat"
action        = "masquerade"
out_interface = "wan-interface"
```

### Port Forwarding
```terraform
chain        = "dstnat"
action       = "dst-nat"
protocol     = "tcp"
dst_port     = "80"
to_addresses = "192.168.1.100"
```

### Hairpin NAT
```terraform
chain       = "srcnat"
action      = "masquerade"
src_address = "LAN-network"
dst_address = "LAN-server"
```

### VPN No-NAT
```terraform
chain            = "srcnat"
action           = "accept"
src_address_list = "vpn-clients"
dst_address_list = "local-networks"
```

## Best Practices

1. **Order Matters**: NAT rules are processed top-to-bottom. Place specific rules before general ones.
2. **Use Accept**: Add NAT exceptions (action=accept) before general masquerade rules.
3. **Hairpin NAT**: Always implement for internal access to port-forwarded services.
4. **Logging**: Use sparingly - high traffic rules can flood logs.
5. **Connection State**: Use `connection_state` for efficiency.
6. **Comments**: Always add descriptive comments.

## Troubleshooting

### Port Forwarding Not Working

1. Check firewall filter rules (INPUT/FORWARD chains)
2. Verify `to_addresses` is correct
3. Ensure service is running on target
4. Check NAT rule order
5. Use logging to debug

### Hairpin NAT Issues

- Ensure hairpin rule is AFTER port forward rule
- Use `masquerade` in srcnat chain
- Match exact dst_address and port

### Multi-WAN NAT

- Mark connections before NAT
- Use separate NAT rules per WAN
- Match on `connection_mark`

## Performance Considerations

- Use interface lists for multiple interfaces
- Avoid unnecessary matching criteria
- Use connection state when possible
- Minimize logging on high-traffic rules

## Security Notes

- Never expose RDP (3389) directly - use VPN
- Change default SSH port (22)
- Use port knocking for sensitive services
- Monitor logs for unusual NAT activity
- Implement rate limiting on exposed ports

## RouterOS Version Compatibility

- Minimum RouterOS version: 7.0
- All NAT features available in v7.x
- Some advanced matching requires 7.6+

## References

- [MikroTik Firewall NAT Documentation](https://help.mikrotik.com/docs/display/ROS/NAT)
- [RouterOS Firewall Best Practices](https://help.mikrotik.com/docs/display/ROS/Firewall)
