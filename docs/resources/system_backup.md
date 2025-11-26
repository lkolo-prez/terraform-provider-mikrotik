---
page_title: "mikrotik_system_backup Resource - terraform-provider-mikrotik"
subcategory: "System"
description: |-
  Creates and manages RouterOS system backup files for disaster recovery and configuration snapshots.
---

# mikrotik_system_backup (Resource)

Creates a system backup file on RouterOS using the `/system/backup/save` command. The backup is stored as a `.backup` file on the router's file system and can be used for disaster recovery, configuration rollback, or compliance purposes.

**Important Notes:**
- Backups are **immutable** - any changes to `name` or `password` will destroy and recreate the backup
- This resource **creates** backup files but does **not download** them - use FTP/SFTP to retrieve files
- Password is **required for restore** - store securely (e.g., HashiCorp Vault, AWS Secrets Manager)
- Backup files include full system configuration (except user files)
- Restore operation is manual via RouterOS `/system/backup/load` command

## Features

- **Automated Backup Creation**: Create backups via Terraform apply
- **Auto-generated Names**: Optional timestamp-based naming
- **Encryption**: Password-protected backups (default)
- **State Tracking**: Monitors backup file existence and size
- **Compliance**: Automated configuration snapshots
- **Disaster Recovery**: Pre-change backups with verification

## Example Usage

### Basic Backup (Auto-generated Name)

```hcl
resource "mikrotik_system_backup" "basic" {
  # Name will be auto-generated: terraform-backup-20251126-143022
}
```

### Named Backup

```hcl
resource "mikrotik_system_backup" "nightly" {
  name = "nightly-backup"
}
```

### Encrypted Backup with Password

```hcl
resource "mikrotik_system_backup" "encrypted" {
  name     = "secure-backup"
  password = var.backup_password
}
```

### Backup Without Encryption

```hcl
resource "mikrotik_system_backup" "plain" {
  name         = "plain-backup"
  dont_encrypt = true
}
```

### Pre-Change Backup

```hcl
resource "mikrotik_system_backup" "pre_firewall_update" {
  name     = "before-firewall-${formatdate("YYYYMMDD-HHmmss", timestamp())}"
  password = var.backup_password
}

resource "mikrotik_ip_firewall_filter" "critical_rule" {
  depends_on = [mikrotik_system_backup.pre_firewall_update]
  
  chain    = "forward"
  action   = "drop"
  protocol = "tcp"
  dst_port = "445"
}
```

### Compliance Backup (Daily)

```hcl
resource "mikrotik_system_backup" "compliance" {
  name     = "compliance-${formatdate("YYYYMMDD", timestamp())}"
  password = var.backup_password
}

output "backup_verification" {
  value = {
    name         = mikrotik_system_backup.compliance.name
    size         = mikrotik_system_backup.compliance.size
    created_at   = mikrotik_system_backup.compliance.creation_time
  }
}
```

### Multi-Router Backups

```hcl
resource "mikrotik_system_backup" "router1" {
  provider = mikrotik.router1
  name     = "router1-backup"
  password = var.backup_password
}

resource "mikrotik_system_backup" "router2" {
  provider = mikrotik.router2
  name     = "router2-backup"
  password = var.backup_password
}
```

### Weekly Backup Rotation

```hcl
resource "mikrotik_system_backup" "weekly" {
  count    = 4
  name     = "weekly-${count.index + 1}"
  password = var.backup_password
}
```

## Schema

### Required

This resource has no required arguments. If `name` is not provided, it will be auto-generated.

### Optional

- `name` (String) - Backup filename (without `.backup` extension). If not provided, auto-generated as `terraform-backup-YYYYMMDD-HHmmss`. **NOTE:** Changing this will destroy and recreate the backup.
- `password` (String, Sensitive) - Encryption password for the backup. Required for restore operations. Store securely (e.g., Vault). **NOTE:** Changing this will destroy and recreate the backup.
- `dont_encrypt` (Boolean) - Disable backup encryption. Default: `false`. Use with caution - backups will be readable as plain text.

### Read-Only (Computed)

- `size` (String) - Backup file size (e.g., "128.5KiB", "2.1MiB")
- `creation_time` (String) - Timestamp when backup was created (RouterOS format: "jan/26/2025 14:30:22")
- `id` (String) - Terraform resource identifier (same as backup filename)

## Common Patterns

### Disaster Recovery

```hcl
# Create backup before major infrastructure changes
resource "mikrotik_system_backup" "dr_pre_upgrade" {
  name     = "dr-before-upgrade-v7.18"
  password = var.dr_password
}

# Download backup using null_resource (requires SSH access)
resource "null_resource" "download_backup" {
  depends_on = [mikrotik_system_backup.dr_pre_upgrade]
  
  provisioner "local-exec" {
    command = "scp admin@${var.router_ip}:/${mikrotik_system_backup.dr_pre_upgrade.name}.backup ./backups/"
  }
}
```

### Scheduled Backups (Terraform Cloud)

```hcl
# Use Terraform Cloud scheduled runs to create daily backups
resource "mikrotik_system_backup" "daily" {
  name     = "daily-${formatdate("YYYYMMDD", timestamp())}"
  password = var.backup_password
  
  lifecycle {
    ignore_changes = [name]  # Prevent recreation on every plan
  }
}
```

### Compliance Snapshots

```hcl
# Quarterly compliance backups
resource "mikrotik_system_backup" "quarterly" {
  name     = "compliance-Q${formatdate("Q-YYYY", timestamp())}"
  password = var.compliance_password
}

# Upload to S3 for long-term retention
resource "null_resource" "upload_to_s3" {
  depends_on = [mikrotik_system_backup.quarterly]
  
  provisioner "local-exec" {
    command = <<-EOT
      scp admin@${var.router_ip}:/${mikrotik_system_backup.quarterly.name}.backup /tmp/
      aws s3 cp /tmp/${mikrotik_system_backup.quarterly.name}.backup s3://compliance-backups/mikrotik/
    EOT
  }
}
```

### Pre-Maintenance Backup

```hcl
variable "maintenance_type" {
  description = "Type of maintenance (firewall-update, os-upgrade, etc.)"
  type        = string
}

resource "mikrotik_system_backup" "pre_maintenance" {
  name     = "pre-${var.maintenance_type}-${formatdate("YYYYMMDD-HHmmss", timestamp())}"
  password = var.backup_password
}

# All maintenance changes depend on backup
resource "mikrotik_ip_firewall_filter" "maintenance_rule" {
  depends_on = [mikrotik_system_backup.pre_maintenance]
  # ... rule configuration
}
```

## Best Practices

### Encryption & Password Management

**Always encrypt production backups:**
```hcl
resource "mikrotik_system_backup" "production" {
  name     = "prod-backup"
  password = var.backup_password  # Never hardcode passwords
}
```

**Use external secret management:**
```hcl
# HashiCorp Vault
data "vault_generic_secret" "backup_password" {
  path = "secret/mikrotik/backup"
}

resource "mikrotik_system_backup" "secure" {
  name     = "secure-backup"
  password = data.vault_generic_secret.backup_password.data["password"]
}
```

**AWS Secrets Manager:**
```hcl
data "aws_secretsmanager_secret_version" "backup_password" {
  secret_id = "mikrotik/backup/password"
}

resource "mikrotik_system_backup" "secure" {
  name     = "secure-backup"
  password = jsondecode(data.aws_secretsmanager_secret_version.backup_password.secret_string)["password"]
}
```

### Backup Retention & Rotation

**Daily rotation (7 days):**
```hcl
resource "mikrotik_system_backup" "daily" {
  count    = 7
  name     = "daily-${formatdate("E", timeadd(timestamp(), "${count.index * 24}h"))}"  # Mon, Tue, Wed...
  password = var.backup_password
  
  lifecycle {
    ignore_changes = [name]
  }
}
```

**Manual cleanup (use RouterOS scheduler or external script):**
```bash
# SSH to router and remove old backups
ssh admin@192.168.88.1 '/file remove [find name~"daily-" where creation-time<7d]'
```

### Pre-Change Verification

**Always create backup before critical changes:**
```hcl
resource "mikrotik_system_backup" "pre_change" {
  name     = "before-change-${var.change_id}"
  password = var.backup_password
}

# Wait for backup to complete
resource "time_sleep" "wait_for_backup" {
  depends_on = [mikrotik_system_backup.pre_change]
  create_duration = "5s"
}

# Critical changes depend on backup + wait
resource "mikrotik_ip_route" "critical_route" {
  depends_on = [time_sleep.wait_for_backup]
  # ... route configuration
}
```

### Downloading Backups

**This resource does NOT download files.** Use one of these methods:

**SCP (Secure Copy):**
```bash
scp admin@192.168.88.1:/backup-name.backup ./local-path/
```

**FTP (if enabled on RouterOS):**
```bash
ftp 192.168.88.1
# login, cd /, get backup-name.backup
```

**Terraform provisioner:**
```hcl
resource "null_resource" "download" {
  depends_on = [mikrotik_system_backup.my_backup]
  
  provisioner "local-exec" {
    command = "scp admin@${var.router_ip}:/${mikrotik_system_backup.my_backup.name}.backup ./backups/"
  }
}
```

### CI/CD Integration

**GitHub Actions example:**
```yaml
name: RouterOS Backup
on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM UTC

jobs:
  backup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: hashicorp/setup-terraform@v3
      
      - name: Create Backup
        env:
          TF_VAR_router_host: ${{ secrets.ROUTER_HOST }}
          TF_VAR_router_username: ${{ secrets.ROUTER_USERNAME }}
          TF_VAR_router_password: ${{ secrets.ROUTER_PASSWORD }}
          TF_VAR_backup_password: ${{ secrets.BACKUP_PASSWORD }}
        run: |
          terraform init
          terraform apply -auto-approve
      
      - name: Download Backup
        run: |
          scp ${{ secrets.ROUTER_USERNAME }}@${{ secrets.ROUTER_HOST }}:/daily-backup.backup ./
      
      - name: Upload to S3
        run: |
          aws s3 cp daily-backup.backup s3://backups/mikrotik/$(date +%Y%m%d)/
```

## Restoring from Backup

**This resource only CREATES backups.** To restore:

### Manual Restore (WinBox/SSH)

1. Upload backup file to router (FTP/SCP)
2. Via WinBox: Files → Upload
3. Via SSH: System → Backup → Load
   ```bash
   /system/backup/load name=backup-name password=your-password
   ```
4. Router will reboot automatically

### Automated Restore (SSH)

```bash
# Upload backup
scp backup-name.backup admin@192.168.88.1:/

# Load backup (router will reboot)
ssh admin@192.168.88.1 '/system/backup/load name=backup-name password=your-password'
```

### Restore Considerations

- **Router will reboot** after loading backup
- **All current configuration is replaced** (irreversible)
- **Password required** if backup was encrypted
- **RouterOS version must match** or be compatible
- **User files are NOT included** in backup

## Troubleshooting

### Backup Not Created

**Problem:** Resource shows as created but file doesn't exist

```hcl
# Check if backup file exists
data "mikrotik_file" "check_backup" {
  name = "${mikrotik_system_backup.my_backup.name}.backup"
}

output "backup_exists" {
  value = data.mikrotik_file.check_backup.id != null ? "YES" : "NO"
}
```

**Solutions:**
- Increase wait time: Add `time_sleep` resource after backup
- Check disk space: `/system/resource/print` (free-hdd-space)
- Check file permissions: Ensure user has write access to /file

### Password Lost

**Problem:** Cannot restore backup because password is forgotten

**Solutions:**
- **Prevention:** Store passwords in secret management system (Vault, AWS Secrets Manager)
- **Recovery:** If password lost, backup is unrecoverable (by design)
- **Best Practice:** Document password storage location in runbook

### Backup Size Too Large

**Problem:** Backup file is unexpectedly large

**Causes:**
- RouterOS logs included in backup (check `/log/print`)
- Large user files (use `/file/print` to list)

**Solutions:**
- Clear old logs: `/log/print follow=no where message~"pattern"` then remove
- Exclude user files: Only system config is in backup (user files stored separately)
- Check actual backup size:
  ```hcl
  output "backup_size_bytes" {
    value = mikrotik_system_backup.my_backup.size
  }
  ```

### Import Existing Backup

**Problem:** Backup file already exists on router, want to import to Terraform

```bash
# Import backup by name (without .backup extension)
terraform import mikrotik_system_backup.existing_backup existing-backup-name
```

**Note:** Only imports state, does not recreate backup file.

### Backup Creation Timeout

**Problem:** Resource creation times out waiting for backup file

**Causes:**
- Slow router (large config, slow storage)
- High CPU usage during backup creation

**Solutions:**
- Use `time_sleep` resource for longer wait:
  ```hcl
  resource "time_sleep" "wait_longer" {
    depends_on = [mikrotik_system_backup.slow_backup]
    create_duration = "10s"
  }
  ```
- Check router resources: `/system/resource/print`

### Debugging

Enable debug logging:
```hcl
provider "mikrotik" {
  host     = var.router_host
  username = var.router_username
  password = var.router_password
  
  # Enable debug logging
  insecure = false
}
```

Check RouterOS logs:
```bash
ssh admin@192.168.88.1 '/log/print follow=no where topics~"system"'
```

## Security Considerations

1. **Always use encryption** (password parameter) for production backups
2. **Store passwords securely** (Vault, AWS Secrets Manager, Azure Key Vault)
3. **Never commit passwords** to version control
4. **Rotate passwords regularly** (creates new backup with new password)
5. **Download and store backups externally** (S3, Azure Blob, on-premises storage)
6. **Test restore process** regularly to ensure backups are valid
7. **Use TLS for provider connection** to protect credentials in transit
8. **Limit backup file access** on RouterOS (use user groups)

## Performance

- **Creation time:** 2-10 seconds (depends on config size)
- **File size:** 50 KiB - 500 KiB (typical)
- **CPU impact:** Minimal (backup runs in background)
- **Storage:** Stored in RouterOS /file system (check free space with `/system/resource/print`)

## Limitations

- Does **not** download backup files (use SCP/FTP)
- Does **not** support restore operation (use RouterOS `/system/backup/load`)
- Does **not** include user files (only system configuration)
- **Immutable** - changes to name/password recreate backup
- **No partial backups** - always full system config
- Password **cannot be retrieved** after creation (RouterOS security)

## Import

Import existing backup file by name (without `.backup` extension):

```bash
terraform import mikrotik_system_backup.my_backup backup-name
```

**Note:** Only imports Terraform state, does not modify backup file.

## Related Resources

- `mikrotik_file` - Manage router files (future resource)
- `mikrotik_system_script` - Automate backup creation via scheduler
- `time_sleep` - Add delays for backup creation verification

## References

- [RouterOS System Backup Documentation](https://help.mikrotik.com/docs/display/ROS/Backup)
- [Disaster Recovery Best Practices](https://wiki.mikrotik.com/wiki/Manual:System/Backup)
- [RouterOS File Management](https://help.mikrotik.com/docs/display/ROS/File)
