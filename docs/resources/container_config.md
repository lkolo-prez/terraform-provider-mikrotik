# mikrotik_container_config

Manages global container configuration on RouterOS v7.4+. This is a singleton resource - there is only one container configuration per device.

Configure registry URL, temporary directory, and global memory limits for all containers.

## Example Usage

### Basic Configuration

```hcl
resource "mikrotik_container_config" "main" {
  registry_url = "https://registry-1.docker.io"
  tmpdir       = "disk1/containers/tmp"
  memory_high  = 256000000  # 256MB global limit
}
```

### With Authentication (RouterOS 7.8+)

```hcl
resource "mikrotik_container_config" "authenticated" {
  registry_url = "https://my-private-registry.com"
  tmpdir       = "disk1/containers/tmp"
  username     = "registry-user"
  password     = var.registry_password  # Use variable for security
}
```

### Minimal Configuration

```hcl
resource "mikrotik_container_config" "minimal" {
  tmpdir = "disk1/tmp"
}
```

## Argument Reference

All arguments are optional:

- `registry_url` - External registry URL for downloading containers. Default: `https://registry-1.docker.io` (Docker Hub)
- `tmpdir` - Container extraction directory. Must point to external storage (e.g., `disk1/containers/tmp`)
- `memory_high` - Global RAM usage limit in bytes for all containers (soft limit). Set to `0` for unlimited. When exceeded, container processes are throttled.
- `username` - Username for registry authentication (RouterOS 7.8+)
- `password` - Password or token for registry authentication (RouterOS 7.8+, sensitive)

### Computed

- `id` - Always `container_config` for singleton resource

## Singleton Resource

This resource is a **singleton** - only one instance exists per RouterOS device. The resource ID is always `container_config`.

### Terraform Behavior

- **Create**: Sets container configuration
- **Update**: Updates configuration values
- **Delete**: Resets configuration to defaults (registry URL, empty tmpdir, no limits)
- **Import**: Always imports as ID `container_config`

## Registry Examples

### Docker Hub

```hcl
resource "mikrotik_container_config" "dockerhub" {
  registry_url = "https://registry-1.docker.io"
  tmpdir       = "disk1/tmp"
}
```

### Google Container Registry (GCR)

```hcl
resource "mikrotik_container_config" "gcr" {
  registry_url = "https://gcr.io"
  tmpdir       = "disk1/tmp"
}
```

### Quay.io

```hcl
resource "mikrotik_container_config" "quay" {
  registry_url = "https://quay.io"
  tmpdir       = "disk1/tmp"
}
```

### Private Registry

```hcl
resource "mikrotik_container_config" "private" {
  registry_url = "https://registry.example.com"
  tmpdir       = "disk1/tmp"
  username     = "myuser"
  password     = var.registry_password
}
```

## Memory Limit Behavior

The `memory_high` parameter sets a **soft limit**:

- Processes are NOT killed when limit is exceeded
- Instead, processes are **throttled** and put under memory reclaim pressure
- This prevents OOM kills while encouraging containers to reduce memory usage

Example values:
- `100000000` = 100 MB
- `256000000` = 256 MB  
- `512000000` = 512 MB
- `1073741824` = 1 GB
- `0` = unlimited (default)

Individual containers can override this with their own `memory_high` setting.

## Import

The singleton configuration can be imported:

```
terraform import mikrotik_container_config.main container_config
```

## Notes

### Storage Requirements

- `tmpdir` **must** point to external storage (USB/SATA)
- RouterOS downloads and extracts containers to this directory
- Ensure sufficient space (containers can be several GB)
- Example paths: `disk1/tmp`, `usb1/containers/tmp`

### Registry Authentication

- Authentication requires RouterOS 7.8+
- Use Terraform variables for sensitive credentials:
  ```hcl
  variable "registry_password" {
    type      = string
    sensitive = true
  }
  ```
- Password is stored in state file - use remote state with encryption

### Reset to Defaults

To reset configuration to defaults, use `terraform destroy`:

```bash
terraform destroy -target=mikrotik_container_config.main
```

This sets:
- `registry_url` = `https://registry-1.docker.io`
- `tmpdir` = empty
- `memory_high` = `0` (unlimited)
- `username` / `password` = empty

## See Also

- `mikrotik_container` - Container instances
- [MikroTik Container Documentation](https://help.mikrotik.com/docs/display/ROS/Container)
