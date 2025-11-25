# RouterOS v7 Cheat Sheet

> **Source**: Based on [3zzy's RouterOS v7 Cheatsheet](https://gist.github.com/3zzy/61e356f0bfcd2918d271836e30d80698)  
> **Official Documentation**: [MikroTik RouterOS v7 Upgrade Guide](https://help.mikrotik.com/docs/spaces/ROS/pages/115736772/Upgrading+to+v7)

This cheatsheet combines the most important commands, changes, and best practices for RouterOS v7, focusing on differences from v6. Content is organized by functional area for easy reference.

## Table of Contents

- [I. General System & Updates](#i-general-system--updates)
- [II. Interfaces](#ii-interfaces)
- [III. IP Addressing & Services](#iii-ip-addressing--services)
- [IV. Routing](#iv-routing)
- [V. Firewall](#v-firewall)
- [VI. Queues](#vi-queues)
- [VII. Tools](#vii-tools)
- [VIII. Scripting](#viii-scripting)
- [IX. Wireless](#ix-wireless)
- [X. PPP](#x-ppp)
- [XI. System](#xi-system)
- [XII. Files](#xii-files)
- [XIII. Log](#xiii-log)

---

## I. General System & Updates

### Check RouterOS Version

```routeros
/system resource print
```

### Upgrade RouterOS (v7+)

```routeros
/system package update check-for-updates
/system package update install
```

**New option**: `install ignore-missing` - Upgrades only RouterOS main package, omitting missing packages during manual upgrades

### Upgrade RouterBOOT (after RouterOS upgrade)

```routeros
/system routerboard upgrade
/system reboot
```

### Device Mode

```routeros
/system/device-mode/update container=yes
```

Controls features/packages on device including container and wifi. Confirm with reset button or power cycle to enable features.

### System Note (Login Banner)

```routeros
/system note set note="Your Text Here"
```

### System History

```routeros
/system/history
```

Shows exact CLI commands executed during "Undo" or "Redo". New verbose mode not available in v6.

### Console Settings

```routeros
/console/settings
```

New option `sanitize-names` for file names.

### Scripting Enhancements

```routeros
# Run Script
/system script run script_name

# Enable error reporting
:onerror do={}

# Access local variables
:local myVar
```

Improved error handling and variable scoping.

---

## II. Interfaces

### General Interface Commands

```routeros
/interface print  # list all interfaces
/interface enable [find name=ether1]
/interface disable [find name=ether2]
/interface monitor [find name=ether1]
```

### Ethernet Interfaces

```routeros
/interface ethernet set [find name=ether1] auto-negotiation=no speed=1Gbps full-duplex=yes
```

**Changes**:
- `slave` flag is deprecated
- Introduces `ethernet-defaults`
- Many changes to switch chip features, PoE, auto-negotiation, and speed settings

### Interface Lists

```routeros
/interface list
```

**New**: Interface list members can be added to a bridge (not possible in v6)

### New Interface Types

#### wifi (802.11ax and newer, replaces wireless for supported devices)

```routeros
/interface/wifi/radio
/interface/wifi/channel
/interface/wifi/configuration
/interface/wifi/security
/interface/wifi/access-list
```

#### wifiwave2 (802.11ac Wave 2 and 802.11ax, alternative package)

```routeros
/interface/wifiwave2
/interface/wifiwave2/access-list
/interface/wifiwave2/capsman
```

#### veth (Virtual Ethernet)

```routeros
/interface veth add name=veth1 address=192.168.10.1/24 gateway=192.168.10.254
```

#### wg (WireGuard)

```routeros
/interface wireguard
add listen-port=13231 private-key="YOUR_PRIVATE_KEY" name=wg-interface1
```

#### vtx (VLAN Tunneling eXtensions)

No direct v6 equivalent.

### Bridge (Important Changes)

#### VLAN Filtering (Highly Recommended)

```routeros
/interface bridge
add name=bridge1 vlan-filtering=yes

/interface bridge port
add bridge=bridge1 interface=ether2 pvid=10  # Access port for VLAN 10
add bridge=bridge1 interface=ether3 pvid=20  # Access port for VLAN 20
add bridge=bridge1 interface=ether1           # Trunk port (usually)

/interface bridge vlan
add bridge=bridge1 tagged=ether1 vlan-ids=10,20,30  # Allowed VLANs on trunk
```

#### Hardware Offloading

- Check `hw=yes` flag on bridge ports
- Essential for CRS3xx/CRS5xx

#### STP/RSTP/MSTP

```routeros
/interface bridge set bridge1 protocol-mode=rstp  # or mstp
```

Set individual port properties like `edge`, `point-to-point`, `external-fdb`

#### IGMP/MLD Snooping

```routeros
/interface bridge
set bridge1 igmp-snooping=yes

/interface bridge port
set [find interface=ether1] multicast-router=disabled
```

#### MVRP

```routeros
/interface bridge set bridge1 mvrp=yes
/interface bridge port set ether1 mvrp-registar-state=fixed
```

### VLAN Interface

```routeros
/interface vlan add interface=bridge1 name=vlan10 vlan-id=10
```

**New property**: `use-service-tag`

### Bonding (LAG)

```routeros
/interface bonding
add mode=802.3ad name=bond1 slaves=ether1,ether2
```

### EoIP, EoIPv6, GRE, IPIP

```routeros
/interface eoip
/interface eoipv6
/interface gre
/interface ipip
```

Several parameter changes in v7 (e.g., removal of `arp` parameter and addition of new options for EoIP)

### Loop Protect

```routeros
/interface ethernet
set ether1 loop-protect=on
```

### LTE

```routeros
/interface lte
```

**New**: `/interface lte apn` for APN profiles configuration

### PPPoE Client

```routeros
/interface pppoe-client
```

Added support for several new authentication methods. Removed PAP authentication.

---

## III. IP Addressing & Services

### IPv4 Address

```routeros
/ip address add address=192.168.1.1/24 interface=ether1
```

### IPv6 Address

```routeros
/ipv6 address add address=2001:db8::1/64 interface=ether1
```

### DHCP Client

```routeros
/ip dhcp-client add interface=ether1 disabled=no
```

Options and parameter names adjusted from v6.

### DHCP Server

```routeros
/ip pool add name=dhcp_pool ranges=192.168.1.10-192.168.1.100
/ip dhcp-server network add address=192.168.1.0/24 gateway=192.168.1.1 dns-server=192.168.1.1
/ip dhcp-server add interface=ether1 address-pool=dhcp_pool name=dhcp1
```

Options and parameter names adjusted from v6.

### DHCPv6 Server

```routeros
/ipv6 dhcp-server
add address-pool=ipv6pool1 disabled=no interface=vlan1 name=dhcpv6_1
```

### IP Pool

```routeros
/ip pool add name=my-pool ranges=192.168.10.10-192.168.10.20
```

### IP Services

```routeros
/ip service
set telnet disabled=yes
set www-ssl certificate=my-cert.crt disabled=no  # Enable HTTPS
set winbox address=192.168.88.0/24 # Restrict Winbox access
```

**New**: `address` property to specify from which address a service is accessible

### DNS Configuration

```routeros
/ip dns
set servers=1.1.1.1,8.8.8.8
set allow-remote-requests=yes # Enable DNS cache

# Static DNS Entry
/ip dns static add name=example.com address=192.168.1.10
```

Adds support for regexp and more record types.

### DoH (DNS over HTTPS)

```routeros
/ip dns set use-doh-server=https://cloudflare-dns.com/dns-query verify-doh-cert=yes
```

**New in v7**, no v6 equivalent.

### Neighbor Discovery

```routeros
/ip/neighbor/discovery-settings
```

**New property**: `discover-interface-list`

### IP Proxy

```routeros
/ip proxy
```

Changes to `cache-size` parameter (v6 provided size in KiB)

### IP Socks

```routeros
/ip socks
```

Removed `connection-idle-timeout` parameter

### IPv6 Neighbor Discovery

```routeros
/ipv6/nd
```

Adds parameters like `ra-delay`, `reachable-time`

---

## IV. Routing

### Static Route

```routeros
/ip route add dst-address=0.0.0.0/0 gateway=192.168.1.254
```

### Routing - Major Overhaul

New top-level `/routing` menu for all routing-related configuration. Replaces separate menus like `/ip route`, `/ipv6 route`, `/routing ospf`, etc.

### Routing Tables

```routeros
/routing table add name=myVrf fib
```

`FIB` parameter indicates if the table should push routes to FIB.

### Routing Rules

```routeros
/routing rule
add dst-address=8.8.8.8/32 action=lookup-only-in-table table=Table1
```

Replaces `/ip route rule`, more versatile rule-based routing.

### Routing Route

```routeros
/routing/route
```

Shows all routes (all address families) and detailed route information. Read-only, replaces functionality of old print commands.

### VRF (Virtual Routing and Forwarding)

```routeros
# Create a new VRF
/routing/table add name=myVrf fib

# Assign an interface to the VRF
/ip vrf add interface=ether2 vrf=myVrf

# Add route to a VRF
/ip route add dst-address=10.0.0.0/24 gateway=192.168.1.1@main routing-table=myVrf
```

### OSPF - Complete Redesign

```routeros
/routing ospf instance add name=default-v2 router-id=0.0.0.1
/routing ospf area add name=backbone instance=default-v2 area-id=0.0.0.0
/routing ospf interface-template add networks=192.168.1.0/24 area=backbone
```

**Changes**:
- Combines OSPFv2 and OSPFv3 into a single menu
- Uses instance and area templates
- Introduces "interface-template" concept
- Monitor output moved to dedicated menus

### BGP - Complete Redesign

```routeros
/routing/bgp
```

**Changes**:
- Uses `connection`, `template`, and `session` sub-menus for structured approach
- Includes many new options and changes to attribute handling

### RIP

```routeros
/routing/rip
```

Similar structure to v6, but with new parameters and options.

### Route Filters

```routeros
/routing filter
```

**Completely new system** using script-like syntax. Much more powerful and flexible than v6's routing filters.

---

## V. Firewall

### Basic Firewall Concepts (v7)

- **connection-state**: `established`, `related`, `new`, `invalid`, `untracked`
- **Chains**: `input`, `forward`, `output` (plus `prerouting`, `postrouting` for NAT/mangle)
- **Actions**: `accept`, `drop`, `reject`, `log`, `add-dst-to-address-list`, etc.

### Basic Firewall Example (IPv4)

```routeros
/ip firewall filter
add action=accept chain=input connection-state=established,related,untracked comment="accept established,related,untracked"
add action=drop chain=input connection-state=invalid comment="drop invalid"
add action=accept chain=input in-interface-list=LAN comment="accept lan to router"
add action=drop chain=input comment="drop all else"

add action=fasttrack-connection chain=forward connection-state=established,related comment="fasttrack"
add action=accept chain=forward connection-state=established,related comment="accept established,related, untracked"
add action=drop chain=forward connection-state=invalid comment="drop invalid"
```

**Changes**:
- Significant changes to structure and available matchers/actions
- Adds new connection-state like `untracked`

### NAT (Network Address Translation)

```routeros
/ip firewall nat
add action=masquerade chain=srcnat out-interface-list=WAN

# Port Forwarding Example
add action=dst-nat chain=dstnat protocol=tcp dst-port=80 to-addresses=192.168.1.10 to-ports=80
```

### Mangle

```routeros
/ip firewall mangle
add action=mark-connection chain=prerouting connection-state=new in-interface=ether1 new-connection-mark=my_conn_mark
add action=mark-packet chain=forward connection-mark=my_conn_mark new-packet-mark=my_pkt_mark
```

---

## VI. Queues

### Simple Queues

```routeros
/queue simple
add name=queue1 target=192.168.1.10/32 max-limit=1M/2M
```

**Note**: Not working if Fasttrack is used

### Queue Tree

More advanced, hierarchical queuing. Requires packet marks.

### Queue Types

```routeros
/queue type
add name=pcq-download kind=pcq pcq-rate=1M pcq-classifier=dst-address
```

**Available types**: `pcq`, `red`, `sfq`, `fifo`, `cake`, `fq_codel`

---

## VII. Tools

### Ping

```routeros
/ping 8.8.8.8
```

### Traceroute

```routeros
/tool traceroute 8.8.8.8
```

### Torch (Real-time traffic monitoring)

```routeros
/tool torch interface=ether1
```

### Packet Sniffer

```routeros
/tool sniffer set filter-interface=ether1 file-name=capture.pcap
/tool sniffer start
/tool sniffer stop
```

### Bandwidth Test

```routeros
/tool bandwidth-test address=192.168.1.2 user=admin password=""
```

Added `local-tx-speed`, `remote-tx-speed`. Changed `random-data` behavior.

### Traffic Generator

```routeros
/tool/traffic-generator
```

Significant changes to stream, packet-templates, and stats configurations.

### Profile

```routeros
/tool profile
```

### Netwatch

```routeros
/tool netwatch add host=192.168.88.1 interval=30s \
  up-script="/system script run upScript" \
  down-script="/system script run downScript"
```

### Fetch

```routeros
/tool fetch url="https://example.com/file.txt"
```

---

## VIII. Scripting

### Basic Syntax

```routeros
:local myVar "Hello"
:put $myVar

:if (condition) do={
    # commands
} else={
    # commands
}

:foreach i in=[/interface find] do={ :put $i }
```

### Error Handling

```routeros
:onerror do={
  :log error "An error occurred"
}
```

More robust error handling than v6.

### Variable Scope

Improved variable scoping:
- Local variables with `:local`
- Global variables remain available

---

## IX. Wireless

### Three Different Wireless Packages

#### 1. `/interface wireless` (Legacy)

Traditional wireless configuration similar to RouterOS v6. For older hardware.

```routeros
/interface wireless scan
/interface wireless registration-table
/interface wireless security-profiles
/interface wireless set [find name=wlan1] ...
/interface wireless monitor
```

#### 2. `/interface wifi` (New - 802.11ax and newer)

New configuration paradigm for 802.11ax and some newer 802.11ac devices.

```routeros
/interface/wifi/radio          # Basic radio settings
/interface/wifi/channel        # Channel configurations
/interface/wifi/configuration  # Core configuration profiles
/interface/wifi/security       # Security profiles
/interface/wifi/provisioning   # Rules for applying configurations
/interface/wifi/access-list    # MAC address-based access control
```

**Example (simple WPA2-PSK AP)**:

```routeros
/interface/wifi/channel
add name=ch-5ghz-1 frequency=5180,5200,5220,5240 width=20/40/80mhz

/interface/wifi/security
add name=sec-home authentication-types=wpa2-psk,wpa3-psk \
    group-ciphers=aes-ccm pairwise-ciphers=aes-ccm \
    passphrase="MyWiFiPassword"

/interface/wifi/configuration
add name=conf-home mode=ap ssid=MyHomeWiFi country=your_country \
    channel=ch-5ghz-1 security=sec-home

/interface/wifi/provisioning
add action=create-dynamic-enabled radio-mac=00:00:00:00:00:00 \
    master-configuration=conf-home

# Enable the radio
/interface/wifi/radio
set [find name="radio1"] band=5ghz-ax
```

#### 3. `/interface wifiwave2` (Alternative - 802.11ac Wave 2 and 802.11ax)

Another different configuration paradigm for older devices not supported by "wifi" package.

```routeros
/interface/wifiwave2
/interface/wifiwave2/access-list
/interface/wifiwave2/capsman
```

**Example (simple WPA2-PSK AP)**:

```routeros
/interface/wifiwave2
set [find default-name=wifi1] \
    configuration.country=Latvia \
    configuration.mode=ap \
    configuration.ssid=MikroTik \
    security.authentication-types=wpa2-psk,wpa3-psk \
    security.passphrase=MyPassword
```

---

## X. PPP

### PPP Client

```routeros
/interface ppp-client
```

### PPP Secret

```routeros
/ppp secret
```

### OpenVPN

- Added more encryption ciphers
- Added UDP mode

### L2TPv3

Now supported in v7, not available in v6.

---

## XI. System

### System Resource

```routeros
/system resource print
```

Output and available information has changed.

### System Package

```routeros
/system package
```

Changes in how updates are handled.

### System Scheduler

```routeros
/system scheduler
```

No significant changes from v6.

### System Script

```routeros
/system script
```

Improved error handling and variable scoping.

### System Routerboard

```routeros
/system routerboard
```

- Settings reorganized
- Added `upgrade-firmware` for RouterBOOT upgrade instead of using reboot

### Reset Configuration

```routeros
/system/reset-configuration
```

Added several new parameters.

---

## XII. Files

### Basic File Management

```routeros
/file print
/file add
/file remove
/file upload
/file download
```

Largely unchanged from v6.

---

## XIII. Log

### Basic Logging

```routeros
/log print
```

Mostly unchanged from v6.

### Log Actions

```routeros
/log action
```

Configure different actions for different log topics.

---

## Important Considerations

⚠️ **Critical Notes**:

1. **Backups**: ALWAYS make a backup before upgrading, especially from v6 to v7
2. **Downgrades**: Downgrading from v7 to v6 is NOT directly supported
3. **Testing**: Test configurations in a lab environment before deploying to production
4. **Documentation**: Refer to official RouterOS v7 documentation for specific features
5. **Default Configuration**: Review carefully as it might have changed significantly
6. **Hardware Offloading**: Pay close attention to hardware offloading in v7, especially for bridges and VLANs
7. **Queues**: Simple queues do not work if Fasttrack is used

---

## Additional Resources

- [Official RouterOS v7 Documentation](https://help.mikrotik.com/docs/spaces/ROS/overview)
- [RouterOS v7 Upgrade Guide](https://help.mikrotik.com/docs/spaces/ROS/pages/115736772/Upgrading+to+v7)
- [MikroTik Wiki](https://wiki.mikrotik.com/)
- [MikroTik Forum](https://forum.mikrotik.com/)

---

**Note**: This cheatsheet provides a starting point. Always refer to the official RouterOS v7 documentation for the most up-to-date and complete information.

**Last Updated**: November 25, 2025  
**RouterOS Version**: 7.14.3 - 7.16.2
