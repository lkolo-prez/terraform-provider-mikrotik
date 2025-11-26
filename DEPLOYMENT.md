# MikroTik Terraform Provider - Deployment Guide

Complete guide for installing, configuring, and deploying the MikroTik Terraform Provider.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
  - [Terraform Registry (Recommended)](#terraform-registry-recommended)
  - [Local Development Build](#local-development-build)
  - [Manual Binary Installation](#manual-binary-installation)
- [Provider Configuration](#provider-configuration)
- [Authentication](#authentication)
- [TLS/SSL Configuration](#tlsssl-configuration)
- [Quick Start](#quick-start)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)
- [Upgrade Guide](#upgrade-guide)

---

## Prerequisites

### Required Software

- **Terraform**: v1.0.0 or later ([Download](https://www.terraform.io/downloads))
- **Go**: v1.21 or later (for building from source) ([Download](https://go.dev/dl/))
- **Git**: For cloning repository ([Download](https://git-scm.com/downloads))

### MikroTik RouterOS Requirements

- **RouterOS Version**: v7.x (7.1 or later recommended)
- **API Access**: RouterOS API service enabled
- **User Permissions**: Admin or API group with full access
- **Network Access**: TCP/8728 (API) or TCP/8729 (API-SSL)

### Enable RouterOS API

```bash
# Connect to RouterOS via SSH or terminal
/ip service enable api
/ip service set api address=0.0.0.0/0  # Or restrict to management subnet

# For TLS/SSL (recommended for production):
/ip service enable api-ssl
/ip service set api-ssl certificate=your-certificate
```

---

## Installation Methods

### Terraform Registry (Recommended)

**When available** (provider published to Terraform Registry):

```hcl
terraform {
  required_providers {
    mikrotik = {
      source  = "ddelnano/mikrotik"
      version = "~> 1.7.0"
    }
  }
}

provider "mikrotik" {
  host     = "192.168.1.1:8728"
  username = "admin"
  password = "your-password"
  tls      = false
}
```

Run:
```bash
terraform init
```

### Local Development Build

**For development or testing latest features:**

#### 1. Clone Repository

```bash
git clone https://github.com/ddelnano/terraform-provider-mikrotik.git
cd terraform-provider-mikrotik
```

#### 2. Build Provider

```bash
# Build for your platform
go build -o terraform-provider-mikrotik

# Or build with version info
go build -ldflags="-X main.version=1.7.0" -o terraform-provider-mikrotik
```

#### 3. Configure Development Override

Create `~/.terraformrc` (Linux/Mac) or `%APPDATA%\terraform.rc` (Windows):

```hcl
provider_installation {
  dev_overrides {
    "ddelnano/mikrotik" = "C:/path/to/terraform-provider-mikrotik"
  }

  # For all other providers, use Terraform Registry
  direct {}
}
```

**Note:** With dev overrides, `terraform init` is NOT required.

#### 4. Test Provider

```bash
cd examples/complete
terraform plan
```

### Manual Binary Installation

**For CI/CD or air-gapped environments:**

#### 1. Download/Build Binary

```bash
# Download release (when available)
wget https://github.com/ddelnano/terraform-provider-mikrotik/releases/download/v1.7.0/terraform-provider-mikrotik_1.7.0_linux_amd64.tar.gz
tar -xzf terraform-provider-mikrotik_1.7.0_linux_amd64.tar.gz

# Or build from source
go build -o terraform-provider-mikrotik
```

#### 2. Install to Terraform Plugin Directory

**Linux/Mac:**
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ddelnano/mikrotik/1.7.0/linux_amd64/
cp terraform-provider-mikrotik ~/.terraform.d/plugins/registry.terraform.io/ddelnano/mikrotik/1.7.0/linux_amd64/
```

**Windows (PowerShell):**
```powershell
New-Item -Path "$env:APPDATA\terraform.d\plugins\registry.terraform.io\ddelnano\mikrotik\1.7.0\windows_amd64" -ItemType Directory -Force
Copy-Item terraform-provider-mikrotik.exe "$env:APPDATA\terraform.d\plugins\registry.terraform.io\ddelnano\mikrotik\1.7.0\windows_amd64\"
```

#### 3. Use in Terraform

```hcl
terraform {
  required_providers {
    mikrotik = {
      source  = "ddelnano/mikrotik"
      version = "1.7.0"
    }
  }
}
```

```bash
terraform init
```

---

## Provider Configuration

### Basic Configuration

```hcl
provider "mikrotik" {
  host     = "192.168.1.1:8728"  # RouterOS IP:port
  username = "admin"
  password = "your-password"
  tls      = false
}
```

### Configuration with TLS/SSL

```hcl
provider "mikrotik" {
  host     = "192.168.1.1:8729"  # API-SSL port
  username = "admin"
  password = "your-password"
  tls      = true
  insecure = true  # For self-signed certificates
}
```

### Configuration with Environment Variables

```bash
export MIKROTIK_HOST="192.168.1.1:8728"
export MIKROTIK_USER="admin"
export MIKROTIK_PASSWORD="your-password"
```

```hcl
provider "mikrotik" {
  # Uses environment variables
}
```

### Multiple RouterOS Devices

```hcl
provider "mikrotik" {
  alias    = "router1"
  host     = "192.168.1.1:8728"
  username = "admin"
  password = "password1"
}

provider "mikrotik" {
  alias    = "router2"
  host     = "192.168.2.1:8728"
  username = "admin"
  password = "password2"
}

# Use specific provider
resource "mikrotik_ip_address" "router1_mgmt" {
  provider = mikrotik.router1
  address  = "10.0.0.1/24"
  interface = "ether1"
}

resource "mikrotik_ip_address" "router2_mgmt" {
  provider = mikrotik.router2
  address  = "10.0.1.1/24"
  interface = "ether1"
}
```

---

## Authentication

### User Permissions

Create dedicated API user on RouterOS:

```bash
# Create API user with full access
/user group add name=terraform policy=api,read,write,policy,test

# Create user
/user add name=terraform password=secure-password group=terraform

# Verify permissions
/user print detail where name=terraform
```

### Password Security

**Option 1: Terraform Variables**

```hcl
variable "mikrotik_password" {
  type      = string
  sensitive = true
}

provider "mikrotik" {
  host     = "192.168.1.1:8728"
  username = "terraform"
  password = var.mikrotik_password
}
```

```bash
# Set via environment
export TF_VAR_mikrotik_password="secure-password"
terraform apply
```

**Option 2: HashiCorp Vault**

```hcl
data "vault_generic_secret" "mikrotik" {
  path = "secret/mikrotik/router1"
}

provider "mikrotik" {
  host     = "192.168.1.1:8728"
  username = data.vault_generic_secret.mikrotik.data["username"]
  password = data.vault_generic_secret.mikrotik.data["password"]
}
```

**Option 3: AWS Secrets Manager**

```hcl
data "aws_secretsmanager_secret_version" "mikrotik" {
  secret_id = "mikrotik/router1"
}

locals {
  mikrotik_creds = jsondecode(data.aws_secretsmanager_secret_version.mikrotik.secret_string)
}

provider "mikrotik" {
  host     = "192.168.1.1:8728"
  username = local.mikrotik_creds.username
  password = local.mikrotik_creds.password
}
```

---

## TLS/SSL Configuration

### Generate Self-Signed Certificate (RouterOS)

```bash
# Generate certificate
/certificate add name=api-cert common-name=router.local key-size=2048
/certificate sign api-cert

# Wait for certificate to be issued
/certificate print

# Enable API-SSL
/ip service set api-ssl certificate=api-cert
/ip service enable api-ssl

# Optionally disable non-TLS API
/ip service disable api
```

### Terraform Configuration with TLS

**Self-Signed Certificate:**
```hcl
provider "mikrotik" {
  host     = "192.168.1.1:8729"
  username = "admin"
  password = "password"
  tls      = true
  insecure = true  # Skip certificate verification
}
```

**Trusted CA Certificate:**
```hcl
provider "mikrotik" {
  host       = "router.example.com:8729"
  username   = "admin"
  password   = "password"
  tls        = true
  insecure   = false
  ca_certificate = file("${path.module}/ca-cert.pem")
}
```

### Certificate Best Practices

1. **Use TLS in Production**: Always encrypt API traffic
2. **Valid Certificates**: Use Let's Encrypt or internal CA
3. **Rotate Regularly**: Update certificates before expiration
4. **Restrict Access**: Firewall API-SSL to management subnet
5. **Monitor Expiration**: Alert before certificate expires

---

## Quick Start

### 1. Create Terraform Configuration

**`main.tf`:**
```hcl
terraform {
  required_version = ">= 1.0"
  required_providers {
    mikrotik = {
      source  = "ddelnano/mikrotik"
      version = "~> 1.7.0"
    }
  }
}

provider "mikrotik" {
  host     = var.mikrotik_host
  username = var.mikrotik_username
  password = var.mikrotik_password
  tls      = false
}

# Example: Configure IP address
resource "mikrotik_ip_address" "management" {
  address   = "10.0.0.1/24"
  interface = "ether1"
  comment   = "Management IP"
}

# Example: System logging
resource "mikrotik_system_logging_action" "remote_syslog" {
  name        = "remote-syslog"
  target      = "remote"
  remote      = "192.168.1.100"
  bsd_syslog  = true
}

resource "mikrotik_system_logging" "firewall" {
  topics = "firewall,info"
  action = mikrotik_system_logging_action.remote_syslog.name
}

# Example: SNMP monitoring
resource "mikrotik_snmp" "monitoring" {
  enabled        = true
  contact        = "noc@company.com"
  location       = "Datacenter A"
  trap_version   = "2"
  trap_target    = "192.168.1.50"
}

resource "mikrotik_snmp_community" "zabbix" {
  name        = "zabbix"
  read_access = true
  address     = "192.168.1.50/32"
}
```

**`variables.tf`:**
```hcl
variable "mikrotik_host" {
  type        = string
  description = "MikroTik router IP and port"
  default     = "192.168.1.1:8728"
}

variable "mikrotik_username" {
  type    = string
  default = "admin"
}

variable "mikrotik_password" {
  type      = string
  sensitive = true
}
```

**`terraform.tfvars`:**
```hcl
mikrotik_host     = "192.168.1.1:8728"
mikrotik_username = "terraform"
mikrotik_password = "your-secure-password"
```

### 2. Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Validate configuration
terraform validate

# Preview changes
terraform plan

# Apply configuration
terraform apply

# View state
terraform show
```

---

## Production Deployment

### Directory Structure

```
mikrotik-infrastructure/
├── environments/
│   ├── production/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── terraform.tfvars
│   │   └── backend.tf
│   ├── staging/
│   │   └── ...
│   └── development/
│       └── ...
├── modules/
│   ├── networking/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── firewall/
│   │   └── ...
│   ├── monitoring/
│   │   └── ...
│   └── wifi/
│       └── ...
└── README.md
```

### Remote State Backend

**S3 Backend:**
```hcl
terraform {
  backend "s3" {
    bucket         = "company-terraform-state"
    key            = "mikrotik/production/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}
```

**Terraform Cloud:**
```hcl
terraform {
  cloud {
    organization = "your-organization"
    workspaces {
      name = "mikrotik-production"
    }
  }
}
```

### CI/CD Integration

**GitHub Actions Example:**

```yaml
name: Terraform MikroTik Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: 1.7.0
      
      - name: Terraform Init
        run: terraform init
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
      
      - name: Terraform Plan
        run: terraform plan
        env:
          TF_VAR_mikrotik_host: ${{ secrets.MIKROTIK_HOST }}
          TF_VAR_mikrotik_username: ${{ secrets.MIKROTIK_USER }}
          TF_VAR_mikrotik_password: ${{ secrets.MIKROTIK_PASSWORD }}
      
      - name: Terraform Apply
        if: github.ref == 'refs/heads/main'
        run: terraform apply -auto-approve
        env:
          TF_VAR_mikrotik_host: ${{ secrets.MIKROTIK_HOST }}
          TF_VAR_mikrotik_username: ${{ secrets.MIKROTIK_USER }}
          TF_VAR_mikrotik_password: ${{ secrets.MIKROTIK_PASSWORD }}
```

### Security Best Practices

1. **Never commit credentials**: Use variables and secrets
2. **Enable TLS**: Always use API-SSL in production
3. **Restrict API access**: Firewall to management subnet only
4. **Use dedicated user**: Create `terraform` user with minimal permissions
5. **State encryption**: Encrypt Terraform state at rest
6. **Audit logging**: Enable RouterOS logging for API access
7. **Review plans**: Always review `terraform plan` before apply
8. **Change tracking**: Use Git for version control

---

## Troubleshooting

### Connection Issues

**Problem:** `Error connecting to RouterOS`

**Solutions:**
```bash
# Verify API service is running
/ip service print

# Check firewall rules
/ip firewall filter print where chain=input

# Test connectivity
ping 192.168.1.1
telnet 192.168.1.1 8728

# Verify credentials
ssh admin@192.168.1.1
```

### Authentication Errors

**Problem:** `authentication failed`

**Solutions:**
```bash
# Verify user exists
/user print

# Check user permissions
/user group print

# Reset password
/user set terraform password=new-password

# Check IP restrictions
/user print detail where name=terraform
```

### TLS Certificate Errors

**Problem:** `certificate verify failed`

**Solutions:**
```hcl
# Option 1: Use insecure (self-signed)
provider "mikrotik" {
  tls      = true
  insecure = true
}

# Option 2: Provide CA certificate
provider "mikrotik" {
  tls            = true
  insecure       = false
  ca_certificate = file("ca-cert.pem")
}
```

### State Lock Errors

**Problem:** `state locked`

**Solutions:**
```bash
# Check lock status
terraform force-unlock <lock-id>

# Or wait for lock timeout (usually 10 minutes)
```

### Resource Not Found

**Problem:** `resource not found after creation`

**Solutions:**
- Wait a few seconds (RouterOS propagation delay)
- Verify resource was created in RouterOS CLI
- Check if resource ID changed
- Use `terraform refresh` to sync state

### Enable Debug Logging

```bash
# Enable Terraform debug logging
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform-debug.log

terraform apply

# Review logs
cat terraform-debug.log
```

---

## Upgrade Guide

### Upgrading Provider Version

#### 1. Review Changelog

Read [CHANGELOG.md](CHANGELOG.md) for breaking changes.

#### 2. Update Version Constraint

```hcl
terraform {
  required_providers {
    mikrotik = {
      source  = "ddelnano/mikrotik"
      version = "~> 1.7.0"  # Update version
    }
  }
}
```

#### 3. Upgrade Provider

```bash
terraform init -upgrade
```

#### 4. Test Changes

```bash
terraform plan  # Review changes
terraform apply # Apply if acceptable
```

### Breaking Changes

#### v1.7.0 (System Logging)
- No breaking changes
- New resources: `mikrotik_system_logging`, `mikrotik_system_logging_action`, `mikrotik_snmp`, `mikrotik_snmp_community`

#### v1.6.0 (VRRP & NAT)
- No breaking changes
- New resources: `mikrotik_vrrp_interface`, `mikrotik_firewall_nat`

#### v1.5.0 (WiFi 6)
- WiFi resources use new `/interface/wifi/` path (RouterOS v7)
- Legacy wireless uses `/interface/wireless/`

#### v1.4.0 (Routing Filters)
- OSPF v7 replaces OSPF v6
- Routing filters use new `/routing/filter/` path

### Rollback Procedure

```bash
# Downgrade provider version
terraform {
  required_providers {
    mikrotik = {
      version = "1.6.0"  # Previous version
    }
  }
}

# Re-initialize
terraform init -upgrade

# Restore previous state
terraform apply
```

---

## Additional Resources

- **GitHub Repository**: [terraform-provider-mikrotik](https://github.com/ddelnano/terraform-provider-mikrotik)
- **Issue Tracker**: [GitHub Issues](https://github.com/ddelnano/terraform-provider-mikrotik/issues)
- **Examples**: [examples/](https://github.com/ddelnano/terraform-provider-mikrotik/tree/main/examples)
- **Documentation**: [docs/resources/](https://github.com/ddelnano/terraform-provider-mikrotik/tree/main/docs/resources)
- **MikroTik Documentation**: [MikroTik Wiki](https://help.mikrotik.com/docs/)
- **Terraform Documentation**: [Terraform Registry](https://registry.terraform.io/providers/ddelnano/mikrotik/latest/docs)

---

## Support

- **Community**: GitHub Discussions
- **Issues**: GitHub Issues (bug reports, feature requests)
- **Security**: Report vulnerabilities privately to maintainers

---

## License

This provider is released under the MIT License. See [LICENSE](LICENSE) for details.
