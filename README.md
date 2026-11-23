# Mikrotik provider for Terraform — extended for RouterOS 7

> **About this repository.** This is an extended version of [ddelnano/terraform-provider-mikrotik](https://github.com/ddelnano/terraform-provider-mikrotik) (MIT), maintained by [Łukasz Kołodziej](https://github.com/lkolo-prez) ([AIRCLOUD](https://aircloud.com.pl)). All credit for the original provider goes to its authors and contributors.
>
> **Added here (11.2025, up to v1.8.0):**
> - OSPF v3 (RouterOS 7)
> - Container support
> - Routing filters (RouterOS 7)
> - VRRP interfaces for high availability
> - Wi-Fi resource management (RouterOS 7)
> - System backup and disaster recovery
> - A unified CI/CD pipeline (test → build → release on tag) and dependency fixes
>
> Releases: [GitHub Releases](https://github.com/lkolo-prez/terraform-provider-mikrotik/releases). The original documentation on the Terraform Registry below describes the upstream provider.


## Intro

This is a terraform provider for managing resources on your RouterOS device. To see what resources and data sources are supported, please see the [documentation](https://registry.terraform.io/providers/ddelnano/mikrotik/latest/docs) on the terraform registry.

## RouterOS 7 Support

✅ **Full RouterOS 7.x support is now available!**

This provider has been updated to work with MikroTik RouterOS 7.x, including:
- New BGP implementation (`/routing/bgp/connection` and `/routing/bgp/template`)
- Firewall RAW table for pre-connection tracking
- Enhanced VLAN support with hardware acceleration
- Improved WireGuard integration
- All modern RouterOS 7 features

**Migration Guide**: See [MIGRATION_ROUTEROS7.md](./MIGRATION_ROUTEROS7.md) for upgrading from RouterOS 6.x

**What's New**: See [ROUTEROS7_SUPPORT.md](./ROUTEROS7_SUPPORT.md) for detailed information about RouterOS 7 features

**Documentation Index**: See [INDEX.md](./INDEX.md) for complete documentation overview

**RouterOS 7 Cheat Sheet**: See [ROUTEROS7_CHEATSHEET.md](./ROUTEROS7_CHEATSHEET.md) for command reference

**Feature Coverage**: See [ROUTEROS7_COVERAGE.md](./ROUTEROS7_COVERAGE.md) for implementation status

**Release Process**: See [RELEASE_PROCESS.md](./RELEASE_PROCESS.md) for automated release workflow

## Support

You can discuss any issues you have or feature requests in [Discord](https://discord.gg/ZpNq8ez).

## Donations

If you get value out this project and want to show your support you can find me on [patreon](https://www.patreon.com/ddelnano).

## Building provider locally

Requirements:
* [Go](https://go.dev/doc/install) >= 1.21 (required for `log/slog` support)
* [Terraform]() >= 0.14

To build the provider with `make`:
```shell
$ make build
```
which creates a `terraform-provider-mikrotik` binary in repository's root folder.

or build with `go` compiler:
```shell
$ go build -o terraform-provider-mikrotik
```

To use locally built provider, Terraform should be aware of its binary.

It could be done with custom CLI config file:
```hcl
# custom.tfrc

provider_installation {
    dev_overrides {
        "ddelnano/mikrotik" = "/path/to/clones/repository/terraform-provider-mikrotik"
    }

    direct {}
}
```
The [dev_overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) section is available since Terraform `0.14`.

Finally, tell Terraform CLI to use custom confiuration by exporting environment variable:
```shell
$ export TF_CLI_CONFIG_FILE=path/to/custom.tfrc
```

**NOTE**: with `dev_overrides` it is not possible to run `terraform init` (see [official docs](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)) so you should immediately use `terraform plan` and `terraform apply` without initializing.

## Contributing

### Dependencies
- RouterOS. See which versions are supported by what is tested in [CI](.github/workflows/continuous-integration.yml)
- Terraform 0.12+


For code generation of boilerplate code, see [codegen Readme](./cmd/mikrotik-codegen/internal/codegen/README.md)

### Testing

The provider is tested with Terraform's acceptance testing framework. As long as you have a RouterOS device you should be able to run them. Please be aware it will create resources on your device! Code that is accepted by the project will not be destructive for anything existing on your router but be careful when changing test code!

In order to run the tests you will need to set the following environment variables:
```bash
export MIKROTIK_HOST=router-hostname:8728
export MIKROTIK_USER=username
# Please be aware this will put your password in your bash history and is not safe
export MIKROTIK_PASSWORD=password
```

After those environment variables are set you can run the tests with the following command:
```bash
make testacc
```

### Testing without MikroTik hardware

If you do not have MikroTik hardware or virtual machine with pre-installed RouterOS, you still have a way to run tests locally.

To make this happen, install [Docker](https://www.docker.com) on your developer machine, and run from the root of the repository:
```sh
$ make routeros
```
It will start RouterOS container locally and make its API server available at `127.0.0.1:8728`

Just export connection settings
```sh
export MIKROTIK_HOST=127.0.0.1:8728
export MIKROTIK_USER=admin
export MIKROTIK_PASSWORD=""
```

and you are ready to run tests with
```sh
$ make test
```

You can use specific RouterOS version by passing `ROUTEROS_VERSION` argument
```sh
$ make routeros ROUTEROS_VERSION="6.49beta54"
```

or even
```sh
$ make routeros ROUTEROS_VERSION=latest
```

To cleanup everything, just run:
```sh
$ make routeros-clean
```
