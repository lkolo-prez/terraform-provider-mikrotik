# Release Process Guide

## Automated Release Workflow

This provider uses automated semantic versioning and GitHub Actions for releases.

## How It Works

### 1. Commit Format (Conventional Commits)

Use conventional commit messages to automatically determine version bumps:

```bash
# Patch release (v0.0.X) - Bug fixes
git commit -m "fix(bgp): correct connection template validation"

# Minor release (v0.X.0) - New features
git commit -m "feat(bgp): add BGP community filtering support"

# Major release (vX.0.0) - Breaking changes
git commit -m "feat(bgp)!: migrate to RouterOS 7 BGP API"
# or
git commit -m "feat(bgp): migrate to RouterOS 7 BGP API

BREAKING CHANGE: Legacy BGP resources are now deprecated"
```

### 2. Automatic Release Flow

```
Push to master → Auto-Release Action → Create Tag → Release Action → Publish
```

**Step-by-step:**

1. **Developer pushes to master**
   ```bash
   git add .
   git commit -m "feat(bgp): add graceful restart support"
   git push origin master
   ```

2. **Auto-Release workflow triggers**
   - Analyzes commit message
   - Determines version bump (major/minor/patch)
   - Creates and pushes new tag (e.g., `v0.10.0`)

3. **Release workflow triggers on tag**
   - Builds multi-platform binaries
   - Signs checksums with GPG
   - Generates changelog from commits
   - Creates GitHub Release with artifacts

4. **Terraform Registry auto-updates** (if published)
   - Detects new GitHub Release
   - Downloads signed artifacts
   - Publishes new provider version

### 3. Manual Release (if needed)

You can manually trigger a release with specific version bump:

1. Go to GitHub Actions → Auto Release
2. Click "Run workflow"
3. Select version bump type (major/minor/patch)
4. Click "Run workflow"

### 4. Skip CI/CD

To push without triggering release:

```bash
git commit -m "docs: update README [skip ci]"
```

## Publishing to Terraform Registry

### Prerequisites

1. **GPG Key Setup**
   ```bash
   # Generate key (if you don't have one)
   gpg --full-generate-key
   
   # Export public key for Terraform Registry
   gpg --armor --export your-email@example.com > public-key.asc
   
   # Export private key for GitHub Secrets
   gpg --armor --export-secret-keys your-key-id > private-key.asc
   ```

2. **GitHub Secrets Configuration**
   
   Add these secrets to your repository (Settings → Secrets → Actions):
   
   - `GPG_PRIVATE_KEY`: Contents of `private-key.asc`
   - `PASSPHRASE`: Your GPG key passphrase
   - `GITHUB_TOKEN`: Automatically provided by GitHub

3. **Terraform Registry Setup**
   
   a. Go to [Terraform Registry](https://registry.terraform.io/)
   
   b. Sign in with GitHub
   
   c. Click "Publish" → "Provider"
   
   d. Select your repository: `lkolo-prez/terraform-provider-mikrotik`
   
   e. Upload your GPG public key (`public-key.asc`)
   
   f. Verify webhook is configured (automatic)

### First-Time Publishing

1. **Verify all requirements**
   - ✅ Repository name: `terraform-provider-{NAME}`
   - ✅ `terraform-registry-manifest.json` exists
   - ✅ GPG key configured in Terraform Registry
   - ✅ GitHub Secrets configured
   - ✅ Valid license file (LICENSE)
   - ✅ README.md with usage examples

2. **Create initial release**
   ```bash
   # Push your code
   git add .
   git commit -m "feat: initial terraform registry publication"
   git push origin master
   
   # Wait for auto-release to create tag
   # Or manually create tag:
   git tag -a v0.10.0 -m "Release v0.10.0"
   git push origin v0.10.0
   ```

3. **Verify release**
   - Check GitHub Actions for successful build
   - Check GitHub Releases for artifacts
   - Check Terraform Registry for new version (may take 5-10 min)

## Version Strategy

Current version: `v0.9.1`

### Planned Versions

- **v0.10.0** - BGP v7 full implementation (next release)
- **v0.11.0** - Routing Filter support
- **v0.12.0** - VRF support
- **v1.0.0** - Full RouterOS 7 feature parity (stable API)

### Version Bump Rules

- **Major (vX.0.0)**: Breaking changes, API redesign, removed resources
- **Minor (v0.X.0)**: New features, new resources, backward compatible
- **Patch (v0.0.X)**: Bug fixes, documentation, performance improvements

## Commit Message Examples

```bash
# Features
git commit -m "feat(bgp): add BGP confederation support"
git commit -m "feat(firewall): add raw table support for RouterOS 7"
git commit -m "feat(routing): add OSPF v3 support"

# Bug Fixes
git commit -m "fix(bgp): correct peer state validation"
git commit -m "fix(dhcp): handle empty lease time correctly"

# Breaking Changes
git commit -m "feat(api)!: migrate to Terraform Plugin Framework v2"
git commit -m "feat(resources): remove deprecated bgp_peer resource

BREAKING CHANGE: Use mikrotik_bgp_connection instead"

# Documentation
git commit -m "docs(bgp): add migration guide for v6 to v7"
git commit -m "docs: update README with RouterOS 7 examples"

# Performance
git commit -m "perf(client): implement connection pooling"
git commit -m "perf(bgp): optimize batch operations"

# Tests
git commit -m "test(bgp): add integration tests for templates"
git commit -m "test: increase coverage to 80%"

# CI/CD
git commit -m "ci: add automated security scanning"
git commit -m "ci: upgrade to Go 1.23"
```

## Troubleshooting

### Release workflow fails

**Problem**: GPG signing fails
```
Error: gpg: signing failed: No secret key
```

**Solution**: Verify GitHub Secrets
1. Check `GPG_PRIVATE_KEY` is set correctly
2. Verify `PASSPHRASE` matches your key
3. Ensure key is not expired: `gpg --list-keys`

---

**Problem**: GoReleaser fails to build
```
Error: build failed
```

**Solution**: 
1. Test locally: `goreleaser release --snapshot --clean`
2. Check Go version compatibility
3. Verify all dependencies: `go mod tidy`

---

**Problem**: Tag already exists
```
Error: tag 'v0.10.0' already exists
```

**Solution**: Auto-release skips existing tags automatically. If manual tag creation is needed:
```bash
# Delete local tag
git tag -d v0.10.0

# Delete remote tag
git push origin :refs/tags/v0.10.0

# Recreate tag
git tag -a v0.10.0 -m "Release v0.10.0"
git push origin v0.10.0
```

### Terraform Registry issues

**Problem**: Provider not appearing in registry

**Solution**:
1. Verify repository name format: `terraform-provider-{NAME}`
2. Check `terraform-registry-manifest.json` exists in root
3. Verify GPG key is uploaded to registry
4. Wait 10-15 minutes for registry sync
5. Check GitHub webhook logs (Settings → Webhooks)

---

**Problem**: Signature verification fails

**Solution**:
1. Re-upload public key to Terraform Registry
2. Ensure GPG key matches between registry and GitHub Secrets
3. Verify key has not expired

## Testing Before Release

### Local Build Test

```bash
# Build provider locally
make build

# Test with local provider
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ddelnano/mikrotik/99.0.0/linux_amd64/
cp terraform-provider-mikrotik ~/.terraform.d/plugins/registry.terraform.io/ddelnano/mikrotik/99.0.0/linux_amd64/

# Use in Terraform
cd examples/bgp
terraform init
terraform plan
```

### GoReleaser Dry Run

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Test release process (no publish)
goreleaser release --snapshot --clean

# Check generated artifacts
ls -lah dist/
```

### CI/CD Test

```bash
# Test workflows locally with act
# Install: https://github.com/nektos/act
act -j build

# Or push to feature branch to test without releasing
git checkout -b test-release
git push origin test-release
```

## Quick Reference

```bash
# Standard release (automatic versioning)
git add .
git commit -m "feat(bgp): add new feature"
git push origin master

# Bug fix release
git commit -m "fix(bgp): correct validation"
git push origin master

# Major breaking change
git commit -m "feat!: migrate to new API"
git push origin master

# Skip release
git commit -m "docs: update [skip ci]"
git push origin master

# Manual tag creation (bypass auto-release)
git tag -a v0.10.0 -m "Release v0.10.0"
git push origin v0.10.0
```

## Resources

- [Terraform Registry - Publishing Providers](https://www.terraform.io/docs/registry/providers/publishing.html)
- [GoReleaser Documentation](https://goreleaser.com)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Actions](https://docs.github.com/en/actions)
