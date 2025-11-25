# CI/CD Workflows

Uporządkowany system CI/CD dla Terraform Provider MikroTik.

## 📋 Przegląd Workflow

### 1. **Continuous Integration** (`continuous-integration.yml`)
**Trigger**: Push/PR do `master`/`main`

Podstawowy workflow budowania i testowania:
- ✅ Budowanie dla Go 1.22, 1.23, 1.24
- ✅ Weryfikacja zależności (`go mod verify`)
- ✅ Linting (`go vet`, `golangci-lint`)
- ✅ Testy jednostkowe (`go test -short`)
- ✅ Testy race condition

**Status**: Aktywny - działa przy każdym commit

---

### 2. **Documentation Validation** (`tfplugindocs.yml`)
**Trigger**: PR z zmianami w kodzie/dokumentacji

Sprawdza czy dokumentacja Terraform jest aktualna:
- ✅ Generowanie dokumentacji (`go generate`)
- ✅ Weryfikacja formatowania przykładów
- ✅ Sprawdzanie niezacommitowanych zmian

**Status**: Aktywny - działa przy PR

---

### 3. **Integration Tests** (`integration-tests.yml`)
**Trigger**: Manualny (`workflow_dispatch`)

Testy integracyjne z RouterOS:
- 🔧 Wymaga kontenera RouterOS
- 🔧 Wybór wersji RouterOS (7.14.3 - 7.17.1)
- 🔧 Pełne testy API
- 📊 Raport pokrycia funkcji

**Status**: Manualny - włączany na żądanie

**Jak uruchomić**:
```bash
# GitHub UI: Actions → Integration Tests → Run workflow
# Wybierz wersję RouterOS i Go
```

---

### 4. **Auto Release** (`auto-release.yml`)
**Trigger**: Push do `master`

Automatyczne tagowanie wersji:
- 🏷️ `feat:` → wersja minor (v1.X.0)
- 🏷️ `fix:` → wersja patch (v1.3.X)
- 🏷️ `feat!:` lub `BREAKING CHANGE:` → wersja major (vX.0.0)

**Format commitów**:
```bash
feat: add OSPF v3 support          # v1.4.0
fix: resolve dependency conflict   # v1.3.8
feat!: change API structure        # v2.0.0
```

**Status**: Aktywny - automatyczne tagowanie

---

### 5. **Release** (`release.yml`)
**Trigger**: Push tagu `v*`

Budowanie i publikacja release:
- 📦 Multi-platform binaries (Windows, Linux, macOS, ARM)
- 🔐 Podpisywanie GPG
- 📄 Generowanie changelog
- 🚀 Publikacja do GitHub Releases
- 📚 Przygotowanie dla Terraform Registry

**Wymagane secrety**:
- `GPG_PRIVATE_KEY` - klucz GPG do podpisywania
- `PASSPHRASE` - hasło do klucza GPG
- `GH_PAT` - Personal Access Token (opcjonalny)

**Status**: Aktywny - działa przy tagach

---

## 🔄 Przepływ Pracy

```mermaid
graph LR
    A[Commit] --> B[CI Build & Test]
    B --> C{Master?}
    C -->|Yes| D[Auto Tag]
    D --> E[Release v*]
    E --> F[Build Artifacts]
    F --> G[GitHub Release]
    
    H[PR] --> I[CI + Docs Check]
    
    J[Manual] --> K[Integration Tests]
```

## 📊 Aktualny Stan

| Workflow | Status | Częstotliwość |
|----------|--------|---------------|
| CI | ✅ Aktywny | Każdy commit |
| Docs | ✅ Aktywny | Każdy PR |
| Integration | 🔧 Manual | Na żądanie |
| Auto Release | ✅ Aktywny | Push master |
| Release | ✅ Aktywny | Tag v* |

## 🚀 Jak Deployować

### Normalny Feature/Fix:
```bash
git add .
git commit -m "feat: add new resource"  # lub "fix: ..."
git push origin master
# → Auto tag → Release
```

### Manual Release:
```bash
git tag -a v1.4.0 -m "Release v1.4.0"
git push origin v1.4.0
# → Release workflow
```

### Testowanie Integracyjne:
1. Idź do **Actions** → **Integration Tests**
2. Kliknij **Run workflow**
3. Wybierz wersję RouterOS
4. Zobacz wyniki i logi

## ⚙️ Konfiguracja

### Wymagane GitHub Secrets:
```yaml
GPG_PRIVATE_KEY: "-----BEGIN PGP PRIVATE KEY BLOCK-----..."
PASSPHRASE: "your-gpg-passphrase"
GH_PAT: "ghp_..." # opcjonalny
```

### Wersje Go:
- **Minimum**: 1.22
- **Recommended**: 1.24
- **CI Matrix**: 1.22, 1.23, 1.24

## 📝 Uwagi

1. **Testy integracyjne** wymagają kontenera RouterOS - uruchamiane manualnie
2. **Auto-release** używa conventional commits do wersjonowania
3. **Release workflow** wymaga prawidłowej konfiguracji GPG
4. **Dokumentacja** musi być zawsze aktualna przed merge PR

## 🔗 Linki

- [Conventional Commits](https://www.conventionalcommits.org/)
- [GoReleaser Docs](https://goreleaser.com/)
- [Terraform Registry Publishing](https://www.terraform.io/docs/registry/providers/publishing.html)
