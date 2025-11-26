# Basic backup (auto-generated name)
resource "mikrotik_system_backup" "basic" {
}

# Named backup
resource "mikrotik_system_backup" "nightly" {
  name = "nightly-backup"
}

# Encrypted backup with password
resource "mikrotik_system_backup" "encrypted" {
  name     = "secure-backup"
  password = "strongpassword123"
}

# Backup without encryption
resource "mikrotik_system_backup" "plain" {
  name         = "plain-backup"
  dont_encrypt = true
}

# Pre-change backup (manual)
resource "mikrotik_system_backup" "pre_change" {
  name     = "before-firewall-update"
  password = var.backup_password
}

# Compliance backup with timestamp
resource "mikrotik_system_backup" "compliance" {
  name     = "compliance-${formatdate("YYYYMMDD", timestamp())}"
  password = var.backup_password
}

# Weekly backup (scheduled via external scheduler)
resource "mikrotik_system_backup" "weekly" {
  name     = "weekly-${formatdate("YYYY-'W'ww", timestamp())}"
  password = var.backup_password
}

# Backup before major changes
resource "mikrotik_system_backup" "pre_upgrade" {
  name     = "pre-upgrade-${var.current_version}"
  password = var.backup_password
}

# Daily backup with rotation (use count for multiple backups)
resource "mikrotik_system_backup" "daily" {
  count    = 7
  name     = "daily-${formatdate("E", timeadd(timestamp(), "${count.index * 24}h"))}"
  password = var.backup_password
}

# Integration with Terraform workspace
resource "mikrotik_system_backup" "workspace" {
  name     = "backup-${terraform.workspace}"
  password = var.backup_password
}

# Multi-router backup
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

# Backup with computed attributes
resource "mikrotik_system_backup" "monitored" {
  name     = "monitored-backup"
  password = var.backup_password
}

output "backup_size" {
  value = mikrotik_system_backup.monitored.size
}

output "backup_created_at" {
  value = mikrotik_system_backup.monitored.creation_time
}

# Disaster recovery backup
resource "mikrotik_system_backup" "dr" {
  name     = "dr-${formatdate("YYYYMMDD-HHmmss", timestamp())}"
  password = var.dr_backup_password
}

# Note: To download backup files, use RouterOS FTP/SFTP after creation:
# scp admin@192.168.88.1:/nightly-backup.backup ./backups/
