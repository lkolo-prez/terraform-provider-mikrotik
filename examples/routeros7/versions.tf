terraform {
  required_version = ">= 1.5"

  required_providers {
    mikrotik = {
      source  = "lkolo-prez/mikrotik"
      version = ">= 1.8.0"
    }
  }
}

provider "mikrotik" {
  host     = var.mikrotik_host
  username = var.mikrotik_username
  password = var.mikrotik_password
  tls      = true
  insecure = var.mikrotik_insecure
}
