# mikrotik_container

Manages OCI containers on RouterOS v7.4+. Run containerized applications like Pi-hole, Prometheus, monitoring agents, or custom services on your MikroTik router.

## Features

- **OCI Compatible**: Works with Docker Hub, GCR, Quay, and other OCI registries
- **Offline Support**: Import containers from tar files
- **Volume Mounts**: Persistent storage for container data
- **Environment Variables**: Configure containers with environment variables
- **Networking**: Full veth interface support with bridge/NAT or host networking
- **Resource Limits**: CPU and memory constraints
- **Lifecycle Management**: Auto-start on boot, auto-restart on failure

## Requirements

- RouterOS v7.4+ with container package installed
- Container mode enabled (requires physical access first time)
- External storage (USB/SATA recommended) for container data
- Pre-configured veth interface for networking

## Example Usage

### Basic Container (Pi-hole DNS)

```hcl
# Global container configuration
resource "mikrotik_container_config" "main" {
  registry_url = "https://registry-1.docker.io"
  tmpdir       = "disk1/containers/tmp"
  memory_high  = 256000000  # 256MB
}

# Veth interface for container networking
resource "mikrotik_interface_veth" "pihole" {
  name    = "veth-pihole"
  address = "172.17.0.2/24"
  gateway = "172.17.0.1"
}

# Pi-hole container
resource "mikrotik_container" "pihole" {
  name         = "pihole"
  remote_image = "pihole/pihole:latest"
  interface    = mikrotik_interface_veth.pihole.name
  root_dir     = "disk1/containers/pihole"
  
  envlist = "ENV_PIHOLE"  # Pre-configured in /container/envs
  mounts  = "MOUNT_PIHOLE_ETC,MOUNT_PIHOLE_DNSMASQ"
  
  logging       = true
  start_on_boot = true
  
  depends_on = [mikrotik_container_config.main]
}
```

### Prometheus Monitoring

```hcl
resource "mikrotik_container" "prometheus" {
  name         = "prometheus"
  remote_image = "prom/prometheus:latest"
  interface    = "veth-monitoring"
  root_dir     = "disk1/containers/prometheus"
  
  mounts = "MOUNT_PROM_CONFIG,MOUNT_PROM_DATA"
  
  cmd = "--config.file=/etc/prometheus/prometheus.yml,--storage.tsdb.path=/prometheus"
  
  memory_high = 512000000  # 512MB limit
  cpu_list    = "0-1"       # Use CPU cores 0 and 1
  
  logging       = true
  start_on_boot = true
}
```

### Offline Container Import

```hcl
resource "mikrotik_container" "custom_app" {
  name      = "myapp"
  file      = "disk1/images/myapp.tar"  # Pre-uploaded tarball
  interface = "veth-app"
  root_dir  = "disk1/containers/myapp"
  
  hostname = "myapp.local"
  dns      = "1.1.1.1,8.8.8.8"
  
  start_on_boot         = true
  auto_restart_interval = "30s"
}
```

## Argument Reference

### Required

- `name` - (Required, Forces new resource) Container name, must be unique
- `interface` - (Required, Forces new resource) Veth interface name for networking
- `root_dir` - (Required, Forces new resource) Root filesystem directory (must be on external storage)

### Optional

**Image Source** (choose one):
- `remote_image` - Container image from registry (e.g., `pihole/pihole:latest`)
- `file` - Path to local container tarball (e.g., `disk1/pihole.tar`)

**Runtime Configuration**:
- `cmd` - Command arguments for entrypoint
- `entrypoint` - Override container entrypoint (e.g., `/bin/sh`)
- `workdir` - Working directory for commands
- `user` - User and group (e.g., `1000:1000` or `www-data`). Set `0:0` for root.

**Data & Configuration**:
- `mounts` - Comma-separated mount names (must be pre-configured)
- `envlist` - Environment variable list name (must be pre-configured)
- `dns` - Custom DNS servers (comma-separated)
- `hostname` - Container hostname
- `domain_name` - Domain name

**Resource Limits**:
- `memory_high` - RAM limit in bytes (soft limit with throttling)
- `cpu_list` - CPU cores allowed (e.g., `0-1` for cores 0 and 1)
- `devices` - Pass-through devices (e.g., `/dev/ttyUSB0`)

**Lifecycle**:
- `logging` - Enable logging to RouterOS logs (default: `false`)
- `start_on_boot` - Auto-start on device boot (default: `false`)
- `auto_restart_interval` - Auto-restart interval on failure (e.g., `10s`)
- `stop_signal` - Linux signal for stop (default: `15` = SIGTERM)

**Other**:
- `comment` - Description

### Computed

- `id` - Container unique ID
- `tag` - Image tag (read-only)
- `digest` - Image digest (read-only)
- `status` - Container status (`stopped`, `running`, `extracting`, etc.)

## Import

Containers can be imported by name:

```
terraform import mikrotik_container.pihole pihole
```

## Notes

### Storage Requirements

- **Always use external storage** (USB/SATA) for `root_dir`
- Built-in flash storage is too small and slow for containers
- Recommended: 100MB/s sequential read/write, 10K random IOPS

### Security Considerations

⚠️ **WARNING**: Enabling container mode reduces device security:
- Container mode requires physical access to enable (reset button)
- Containers run with elevated privileges
- Malicious containers can compromise RouterOS
- Only use trusted container images

### Status Values

- `stopped` - Container not running
- `running` - Container is active
- `extracting` - Image being extracted (first-time setup)
- `error` - Container failed

### Networking Modes

**Bridge with NAT** (recommended):
- Containers share veth interface
- Use firewall NAT for port forwarding
- Containers isolated by default

**Layer2/Host**:
- Container directly on LAN bridge
- All ports exposed
- Less secure but slightly faster

## See Also

- `mikrotik_container_config` - Global container settings
- Environment variables: `/container/envs` (CLI)
- Volume mounts: `/container/mounts` (CLI)
- [MikroTik Container Documentation](https://help.mikrotik.com/docs/display/ROS/Container)
