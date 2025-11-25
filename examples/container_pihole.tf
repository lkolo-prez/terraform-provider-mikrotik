# Pi-hole DNS Container on MikroTik
#
# This example demonstrates running Pi-hole as a container on RouterOS
# for network-wide ad blocking and DNS management.

# Global container configuration
resource "mikrotik_container_config" "main" {
  registry_url = "https://registry-1.docker.io"
  tmpdir       = "disk1/containers/tmp"
  memory_high  = 256000000 # 256MB global soft limit
}

# Network bridge for containers
resource "mikrotik_bridge" "containers" {
  name = "containers"
}

# IP address for container bridge
resource "mikrotik_ip_address" "container_gateway" {
  address   = "172.17.0.1/24"
  interface = mikrotik_bridge.containers.name
}

# Veth interface for Pi-hole
# Note: In real setup, use proper veth interface resource when available
# For now, this must be pre-configured via CLI:
# /interface/veth/add name=veth-pihole address=172.17.0.2/24 gateway=172.17.0.1

# NAT for outgoing container traffic
resource "mikrotik_firewall_nat_rule" "container_masquerade" {
  chain       = "srcnat"
  action      = "masquerade"
  src_address = "172.17.0.0/24"
}

# Port forwarding for Pi-hole web interface
resource "mikrotik_firewall_nat_rule" "pihole_web" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_address  = "192.168.88.1" # Router's LAN IP
  dst_port     = "8080"
  to_addresses = "172.17.0.2"
  to_ports     = "80"
}

# Environment variables for Pi-hole
# Note: Must be pre-configured via CLI:
# /container/envs/add list=ENV_PIHOLE key=TZ value="Europe/Warsaw"
# /container/envs/add list=ENV_PIHOLE key=FTLCONF_webserver_api_password value="admin123"
# /container/envs/add list=ENV_PIHOLE key=DNSMASQ_USER value="root"

# Volume mounts for Pi-hole
# Note: Must be pre-configured via CLI:
# /container/mounts/add name=MOUNT_PIHOLE_ETC src=disk1/volumes/pihole/etc dst=/etc/pihole
# /container/mounts/add name=MOUNT_PIHOLE_DNSMASQ src=disk1/volumes/pihole/dnsmasq dst=/etc/dnsmasq.d

# Pi-hole container
resource "mikrotik_container" "pihole" {
  name         = "pihole"
  remote_image = "pihole/pihole:latest"
  interface    = "veth-pihole"
  root_dir     = "disk1/containers/pihole"

  # Reference pre-configured environment variables and mounts
  envlist = "ENV_PIHOLE"
  mounts  = "MOUNT_PIHOLE_ETC,MOUNT_PIHOLE_DNSMASQ"

  # Pi-hole specific settings
  hostname = "pihole.local"
  dns      = "1.1.1.1,8.8.8.8" # Upstream DNS

  # Lifecycle
  logging       = true
  start_on_boot = true

  depends_on = [
    mikrotik_container_config.main,
    mikrotik_bridge.containers,
    mikrotik_ip_address.container_gateway,
  ]
}

# Output access information
output "pihole_web_access" {
  value = "Access Pi-hole web interface at http://192.168.88.1:8080/admin"
}

output "pihole_container_ip" {
  value = "172.17.0.2"
}

# Setup Instructions:
# 
# 1. Enable container mode (requires physical access):
#    /system/device-mode/update container=yes
#    Press reset button to confirm
# 
# 2. Attach external storage and format:
#    /disk/format-drive usb1 file-system=ext4 label=disk1
# 
# 3. Create veth interface (CLI only for now):
#    /interface/veth/add name=veth-pihole address=172.17.0.2/24 gateway=172.17.0.1
# 
# 4. Add veth to bridge (CLI only for now):
#    /interface/bridge/port add bridge=containers interface=veth-pihole
# 
# 5. Configure environment variables (CLI only for now):
#    /container/envs/add list=ENV_PIHOLE key=TZ value="Europe/Warsaw"
#    /container/envs/add list=ENV_PIHOLE key=FTLCONF_webserver_api_password value="admin123"
#    /container/envs/add list=ENV_PIHOLE key=DNSMASQ_USER value="root"
# 
# 6. Configure volume mounts (CLI only for now):
#    /container/mounts/add name=MOUNT_PIHOLE_ETC src=disk1/volumes/pihole/etc dst=/etc/pihole
#    /container/mounts/add name=MOUNT_PIHOLE_DNSMASQ src=disk1/volumes/pihole/dnsmasq dst=/etc/dnsmasq.d
# 
# 7. Apply Terraform configuration:
#    terraform apply
# 
# 8. Wait for container to download and extract (check status):
#    terraform show mikrotik_container.pihole | grep status
# 
# 9. Configure devices to use Pi-hole:
#    - Set DNS server to 192.168.88.1
#    - Or update DHCP server to distribute Pi-hole as DNS
