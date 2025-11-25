# Setup Instructions for Automated Releases

## Konfiguracja sekretów GitHub (wymagane przed pierwszym release)

### 1. Wygeneruj klucz GPG (jeśli jeszcze nie masz)

```bash
# Generowanie klucza GPG
gpg --full-generate-key

# Wybierz:
# - Typ: RSA and RSA
# - Rozmiar: 4096 bits
# - Ważność: 0 (nie wygasa) lub według potrzeb
# - Podaj imię, email, komentarz
# - Ustaw hasło (PASSPHRASE)
```

### 2. Eksportuj klucze GPG

```bash
# Znajdź ID swojego klucza
gpg --list-secret-keys --keyid-format=long

# Przykładowy output:
# sec   rsa4096/ABCD1234EFGH5678 2024-01-01 [SC]
#       Fingerprint: XXXX XXXX XXXX XXXX XXXX  XXXX XXXX XXXX XXXX XXXX
# uid                 Your Name <your.email@example.com>

# Eksportuj PRYWATNY klucz (dla GitHub Secrets)
gpg --armor --export-secret-keys ABCD1234EFGH5678 > private-key.asc

# Eksportuj PUBLICZNY klucz (dla Terraform Registry)
gpg --armor --export ABCD1234EFGH5678 > public-key.asc

# Skopiuj fingerprint
gpg --fingerprint ABCD1234EFGH5678
```

### 3. Skonfiguruj GitHub Secrets

Przejdź do: `https://github.com/lkolo-prez/terraform-provider-mikrotik/settings/secrets/actions`

Dodaj 2 sekrety:

**a) GPG_PRIVATE_KEY**
```bash
# Skopiuj całą zawartość pliku (włącznie z nagłówkami)
cat private-key.asc

# Wklej do GitHub Secret GPG_PRIVATE_KEY
# Powinno zaczynać się od: -----BEGIN PGP PRIVATE KEY BLOCK-----
# I kończyć: -----END PGP PRIVATE KEY BLOCK-----
```

**b) PASSPHRASE**
```
Wpisz hasło, które ustawiłeś podczas generowania klucza GPG
```

**c) GITHUB_TOKEN** 
```
Ten sekret jest automatycznie dostarczany przez GitHub Actions - nie musisz go tworzyć
```

### 4. Opublikuj providera w Terraform Registry

**Krok 1: Zaloguj się do Terraform Registry**
- Przejdź do: https://registry.terraform.io/
- Kliknij "Sign in" → Zaloguj się przez GitHub

**Krok 2: Opublikuj providera**
- Kliknij "Publish" → "Provider"
- Wybierz repozytorium: `lkolo-prez/terraform-provider-mikrotik`
- Kliknij "Publish provider"

**Krok 3: Dodaj klucz GPG do Terraform Registry**
- W ustawieniach providera kliknij "Add Signing Key"
- Wklej zawartość `public-key.asc`
- Kliknij "Add key"

**Krok 4: Weryfikacja**
- Sprawdź czy webhook jest skonfigurowany (Settings → Webhooks)
- Powinien być webhook: `https://registry.terraform.io/...`

### 5. Test automatycznego release

**Opcja A: Commit z automatycznym tagowaniem**
```bash
# Zrób zmianę
echo "# Test" >> README.md

# Commit z konwencjonalnym formatem
git add README.md
git commit -m "feat(test): test automatic release workflow"

# Push do master
git push origin master

# Auto-release workflow automatycznie:
# 1. Wykryje "feat:" → wersja minor (v0.10.0)
# 2. Utworzy tag v0.10.0
# 3. Release workflow zbuduje artefakty
# 4. Opublikuje na GitHub Releases
# 5. Terraform Registry automatycznie pobierze nową wersję
```

**Opcja B: Ręczne utworzenie tagu (bypass auto-release)**
```bash
# Utwórz tag ręcznie
git tag -a v0.10.0 -m "Release v0.10.0: BGP v7 full implementation"

# Push tagu
git push origin v0.10.0

# Release workflow automatycznie:
# 1. Zbuduje artefakty dla v0.10.0
# 2. Podpisze checksumę GPG
# 3. Utworzy GitHub Release
# 4. Terraform Registry pobierze artefakty
```

### 6. Monitorowanie procesu

**GitHub Actions**
```
Przejdź do: https://github.com/lkolo-prez/terraform-provider-mikrotik/actions

Sprawdź:
1. "Auto Release" workflow - czy utworzył tag
2. "release" workflow - czy zbudował artefakty
3. Sprawdź logi w przypadku błędów
```

**GitHub Releases**
```
Przejdź do: https://github.com/lkolo-prez/terraform-provider-mikrotik/releases

Powinieneś zobaczyć:
- Nowy release (np. v0.10.0)
- Artifacts (.zip dla każdej platformy)
- Checksum (SHA256SUMS)
- Signature (SHA256SUMS.sig)
- Changelog z commitów
```

**Terraform Registry**
```
Przejdź do: https://registry.terraform.io/providers/ddelnano/mikrotik/latest

Po 5-10 minutach powinieneś zobaczyć nową wersję
```

## Troubleshooting

### Problem: "gpg: signing failed: No secret key"

**Rozwiązanie:**
1. Sprawdź czy GPG_PRIVATE_KEY jest poprawnie skopiowany (włącznie z nagłówkami)
2. Sprawdź czy PASSPHRASE jest poprawny
3. Sprawdź czy klucz nie wygasł: `gpg --list-keys`

### Problem: Tag już istnieje

**Rozwiązanie:**
```bash
# Usuń lokalny tag
git tag -d v0.10.0

# Usuń zdalny tag
git push origin :refs/tags/v0.10.0

# Utwórz ponownie
git tag -a v0.10.0 -m "Release v0.10.0"
git push origin v0.10.0
```

### Problem: Provider nie pojawia się w Terraform Registry

**Rozwiązanie:**
1. Sprawdź nazwę repo: MUSI być `terraform-provider-{name}`
2. Sprawdź czy `terraform-registry-manifest.json` istnieje w root
3. Sprawdź czy klucz GPG jest dodany do Terraform Registry
4. Poczekaj 10-15 minut na synchronizację
5. Sprawdź webhook logs (GitHub Settings → Webhooks)

### Problem: Weryfikacja podpisu nie działa

**Rozwiązanie:**
1. Sprawdź czy publiczny klucz w Terraform Registry pasuje do prywatnego w GitHub Secrets
2. Sprawdź fingerprint: `gpg --fingerprint YOUR_KEY_ID`
3. Prześlij ponownie publiczny klucz do Terraform Registry

## Następne kroki

Po skonfigurowaniu wszystkiego:

1. **Usuń pliki z kluczami** (WAŻNE dla bezpieczeństwa!)
   ```bash
   rm -f private-key.asc public-key.asc
   ```

2. **Utwórz pierwszy oficjalny release** z BGP v7:
   ```bash
   git commit -m "feat(bgp): complete RouterOS 7 BGP implementation
   
   - mikrotik_bgp_instance_v7: BGP instance configuration
   - mikrotik_bgp_connection: Peer connections with templates
   - mikrotik_bgp_template: Reusable connection templates
   - mikrotik_bgp_session: Active session monitoring
   - 6 production examples with migration guide
   - Full test coverage with performance optimizations"
   
   git push origin master
   # Automatycznie utworzy v0.10.0 (feat = minor bump)
   ```

3. **Ogłoś release**:
   - Napisz post na Discord
   - Update dokumentacji na Terraform Registry
   - Poinformuj użytkowników o nowej wersji

## Podsumowanie

✅ **Co zostało skonfigurowane:**
- Auto-release workflow z semantic versioning
- Release workflow z GPG signing
- GoReleaser v2 z enhanced changelog
- Terraform Registry manifest
- Dokumentacja procesu release

❌ **Co musisz jeszcze zrobić:**
1. Wygeneruj klucz GPG
2. Skonfiguruj GitHub Secrets (GPG_PRIVATE_KEY, PASSPHRASE)
3. Opublikuj providera w Terraform Registry
4. Dodaj publiczny klucz GPG do Terraform Registry
5. Przetestuj pierwszy release

⏱️ **Czas setup:** ~15-20 minut

📚 **Dokumentacja:**
- `RELEASE_PROCESS.md` - Kompletny przewodnik po procesie release
- `CHANGELOG.md` - Automatycznie generowany changelog
- `.github/workflows/auto-release.yml` - Automatyczne tagowanie
- `.github/workflows/release.yml` - Budowanie i publikacja
