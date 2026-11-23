# RouterOS 7 extensions in one place: OSPFv3, VRRP, Wi-Fi (wifi package) and
# system backup. Every resource below is implemented in this fork of the
# provider (see mikrotik/provider_framework.go).

resource "mikrotik_bridge" "lan" {
  name    = var.lan_bridge
  comment = "LAN bridge (managed by Terraform)"
}

# ---------------------------------------------------------------------------
# OSPFv3 (IPv6) with the RouterOS 7 instance / area / interface-template model
# ---------------------------------------------------------------------------
resource "mikrotik_ospf_instance_v7" "v3" {
  name                   = "ospf-v3"
  version                = "3"
  router_id              = var.router_id
  redistribute_connected = true
  comment                = "OSPFv3 for IPv6"
}

resource "mikrotik_ospf_area_v7" "backbone_v3" {
  name     = "backbone-v3"
  instance = mikrotik_ospf_instance_v7.v3.name
  area_id  = "0.0.0.0"
}

resource "mikrotik_ospf_interface_template_v7" "lan_v3" {
  area       = mikrotik_ospf_area_v7.backbone_v3.name
  interfaces = [mikrotik_bridge.lan.name]
  type       = "broadcast"
  cost       = 10
  passive    = true
  comment    = "Advertise the LAN, no adjacencies on it"
}

# ---------------------------------------------------------------------------
# VRRP gateway redundancy (apply the same code on both routers with a
# different vrrp_priority)
# ---------------------------------------------------------------------------
resource "mikrotik_interface_vrrp" "gateway" {
  name            = "vrrp-lan"
  interface       = mikrotik_bridge.lan.name
  vrid            = 10
  priority        = var.vrrp_priority
  version         = 3
  v3_protocol     = "ipv4"
  interval        = "1s"
  preemption_mode = true
  comment         = "Virtual gateway for the LAN"
}

resource "mikrotik_ip_address" "gateway_vip" {
  address   = var.vrrp_virtual_ip
  interface = mikrotik_interface_vrrp.gateway.name
  comment   = "VRRP virtual IP"
}

# ---------------------------------------------------------------------------
# Wi-Fi (RouterOS 7 "wifi" package, 802.11ax)
# ---------------------------------------------------------------------------
resource "mikrotik_wifi_security" "wpa" {
  name                 = "sec-wpa2-wpa3"
  authentication_types = "wpa2-psk,wpa3-psk"
  encryption           = "ccmp"
  passphrase           = var.wifi_passphrase
}

resource "mikrotik_wifi_channel" "ch5" {
  name              = "ch-5ghz-80"
  band              = "5ghz-ax"
  width             = "20/40/80mhz"
  skip_dfs_channels = true
}

resource "mikrotik_wifi_datapath" "lan" {
  name             = "dp-lan"
  bridge           = mikrotik_bridge.lan.name
  client_isolation = false
}

resource "mikrotik_wifi_configuration" "ap" {
  name     = "cfg-ap"
  ssid     = var.wifi_ssid
  mode     = "ap"
  country  = var.wifi_country
  security = mikrotik_wifi_security.wpa.name
}

resource "mikrotik_interface_wifi" "wifi1" {
  name          = "wifi1"
  configuration = mikrotik_wifi_configuration.ap.name
  channel       = mikrotik_wifi_channel.ch5.name
  datapath      = mikrotik_wifi_datapath.lan.name
  comment       = "5 GHz access point"
}

# ---------------------------------------------------------------------------
# Encrypted system backup taken after the changes above
# ---------------------------------------------------------------------------
resource "mikrotik_system_backup" "after_change" {
  name     = "terraform-routeros7"
  password = var.backup_password

  depends_on = [
    mikrotik_ospf_interface_template_v7.lan_v3,
    mikrotik_ip_address.gateway_vip,
    mikrotik_interface_wifi.wifi1,
  ]
}
