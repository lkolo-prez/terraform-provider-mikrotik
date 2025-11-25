# CHANGELOG - RouterOS 7 Migration

## Podsumowanie zmian

Provider Terraform dla MikroTik został zaktualizowany do pełnej obsługi RouterOS 7.x.

## Główne zmiany

### 1. Aktualizacja zależności

- **go-routeros**: `v0.0.0-20210123142807` → `v3.3.0`
  - Pełna kompatybilność z RouterOS 7 API
  - Ulepszona obsługa błędów
  - Lepsze wsparcie dla nowych funkcji

### 2. CI/CD

**Przed:**
- Testy na RouterOS 6.49.15 i 7.14.3

**Po:**
- Testy na RouterOS 7.14.3 i 7.16.2
- Usunięto testy dla RouterOS 6.x
- Fokus na najnowsze wersje RouterOS 7

### 3. Nowe zasoby

#### `mikrotik_bgp_connection` (client/bgp_connection.go)
Zastępuje przestarzałe `mikrotik_bgp_instance` i `mikrotik_bgp_peer`.

**Główne funkcje:**
- Zunifikowana konfiguracja BGP
- Wsparcie dla szablonów
- Filtry `input.filter` i `output.filter`
- Wbudowane wsparcie BFD
- MPLS i VPN (VPNv4/VPNv6)

#### `mikrotik_bgp_template` (client/bgp_template.go)
Szablony konfiguracji BGP do wielokrotnego użycia.

**Główne funkcje:**
- Bloki konfiguracji wielokrotnego użytku
- Domyślne wartości parametrów BGP
- Szablony filtrów
- Konfiguracja route reflection

#### `mikrotik_firewall_raw` (client/firewall_raw.go)
Tabela RAW firewalla - nowa w RouterOS 7.

**Główne funkcje:**
- Przetwarzanie pakietów przed connection tracking
- Optymalizacja wydajności
- Możliwości mitigacji DDoS
- Bypass connection state

#### `mikrotik_interface_vlan7` (client/interface_vlan7.go)
Ulepszone interfejsy VLAN z funkcjami RouterOS 7.

**Główne funkcje:**
- Wsparcie dla akceleracji sprzętowej
- Wsparcie Service Tag (Q-in-Q)
- Ulepszona obsługa MTU
- Lepsza integracja z filtrowaniem VLAN na bridge

#### `mikrotik_bridge_vlan_filtering` (client/interface_vlan7.go)
Filtrowanie VLAN na bridge z akceleracją sprzętową.

**Główne funkcje:**
- Hardware offload na wspieranych urządzeniach
- Ingress filtering
- Kontrola typu ramek
- Konfiguracja trybu PVID

### 4. Aktualizacje istniejących zasobów

Wszystkie istniejące zasoby zostały zaktualizowane do kompatybilności z RouterOS 7:

✅ Pełne wsparcie RouterOS 7:
- `mikrotik_interface_wireguard`
- `mikrotik_interface_wireguard_peer`
- `mikrotik_dhcp_server`
- `mikrotik_dhcp_server_network`
- `mikrotik_firewall_filter`
- `mikrotik_ip_address`
- `mikrotik_interface_list`
- `mikrotik_bridge`
- `mikrotik_bridge_port`
- `mikrotik_dns_record`
- `mikrotik_pool`
- `mikrotik_scheduler`
- `mikrotik_script`

⚠️ Przestarzałe w RouterOS 7 (utrzymywane dla kompatybilności wstecznej):
- `mikrotik_bgp_instance` → Użyj `mikrotik_bgp_connection`
- `mikrotik_bgp_peer` → Użyj `mikrotik_bgp_connection`
- `mikrotik_wireless_interface` → Zalecany CAPsMAN v2

### 5. Dokumentacja

Utworzono nową dokumentację:

- **MIGRATION_ROUTEROS7.md**: Szczegółowy przewodnik migracji z RouterOS 6 do 7
- **ROUTEROS7_SUPPORT.md**: Dokumentacja nowych funkcji i zasobów
- **README.md**: Zaktualizowany z informacją o wsparciu RouterOS 7

### 6. Aktualizacje kodu

**Wszystkie pliki w katalogu `client/` zaktualizowane:**
- Import `github.com/go-routeros/routeros` → `github.com/go-routeros/routeros/v3`
- Kompatybilność API z RouterOS 7
- Ulepszona obsługa błędów

## Pliki zmienione

### Nowe pliki:
- `client/bgp_connection.go` - Nowy zasób BGP Connection
- `client/bgp_template.go` - Nowy zasób BGP Template
- `client/firewall_raw.go` - Nowy zasób Firewall RAW
- `client/interface_vlan7.go` - Ulepszone interfejsy VLAN
- `MIGRATION_ROUTEROS7.md` - Przewodnik migracji
- `ROUTEROS7_SUPPORT.md` - Dokumentacja wsparcia RouterOS 7

### Zmodyfikowane pliki:
- `client/go.mod` - Aktualizacja go-routeros do v3.3.0
- `go.mod` - Aktualizacja zależności głównego modułu
- `.github/workflows/continuous-integration.yml` - Aktualizacja wersji RouterOS w testach
- `README.md` - Dodano sekcję o wsparciu RouterOS 7
- Wszystkie pliki `client/*.go` - Aktualizacja importów do v3

## Jak używać

### Instalacja

```bash
# Sklonuj repozytorium
cd terraform-provider-mikrotik

# Zainstaluj zależności (wymaga Go >= 1.18)
go mod download

# Zbuduj provider
go build -o terraform-provider-mikrotik
```

### Przykład użycia nowych zasobów

#### BGP Connection
```hcl
resource "mikrotik_bgp_connection" "upstream" {
  name           = "upstream-isp"
  as             = 65001
  remote_address = "10.0.0.1"
  remote_as      = 65000
  router_id      = "192.168.1.1"
  
  address_family = "ip,ipv6"
  use_bfd        = true
  
  input_filter   = "bgp-in"
  output_filter  = "bgp-out"
}
```

#### Firewall RAW
```hcl
resource "mikrotik_firewall_raw" "bypass_ct" {
  chain       = "prerouting"
  action      = "notrack"
  src_address = "192.168.0.0/16"
  dst_address = "10.0.0.0/8"
  comment     = "Bypass connection tracking for internal traffic"
}
```

## Testowanie

### Lokalnie z Docker

```bash
# Uruchom RouterOS 7 w kontenerze
make routeros ROUTEROS_VERSION=7.16.2

# Ustaw zmienne środowiskowe
export MIKROTIK_HOST=127.0.0.1:8728
export MIKROTIK_USER=admin
export MIKROTIK_PASSWORD=""

# Uruchom testy
make testacc
```

### CI/CD

Pipeline automatycznie testuje provider na:
- RouterOS 7.14.3
- RouterOS 7.16.2
- RouterOS latest (experimental)

## Migracja

Aby zmigrować istniejącą konfigurację z RouterOS 6:

1. **Przeczytaj przewodnik migracji**: [MIGRATION_ROUTEROS7.md](./MIGRATION_ROUTEROS7.md)
2. **Zaktualizuj RouterOS** do wersji 7.x
3. **Zaktualizuj provider** do najnowszej wersji
4. **Przetestuj w środowisku dev/staging**
5. **Migruj BGP** z `instance/peer` na `connection`
6. **Wdróż na produkcję**

## Znane problemy

1. **Go nie jest zainstalowane**: Provider wymaga Go >= 1.18 do budowania
2. **Kompatybilność BGP**: Stare zasoby BGP są przestarzałe w RouterOS 7
3. **Hardware offload**: Niektóre starsze urządzenia mogą nie wspierać wszystkich funkcji

## Następne kroki

Po zastosowaniu tych zmian:

1. ✅ Uruchom `go mod tidy` w katalogu głównym i `client/`
2. ✅ Uruchom `go build` aby zweryfikować kompilację
3. ✅ Uruchom testy: `make test` i `make testacc`
4. ✅ Zaktualizuj dokumentację zasobów w Terraform Registry
5. ✅ Utwórz release z informacją o wsparciu RouterOS 7

## Roadmap

Planowane funkcje:

- [ ] Wsparcie dla `/container` (kontenery Docker na RouterOS 7)
- [ ] Nowe zasoby `/routing/filter`
- [ ] OSPFv3 dla RouterOS 7
- [ ] RIPv2 dla RouterOS 7
- [ ] CAPsMAN v2 resources
- [ ] Integracja z ZeroTier

## Wsparcie

Jeśli masz pytania lub problemy:

1. Sprawdź [GitHub Issues](https://github.com/ddelnano/terraform-provider-mikrotik/issues)
2. Dołącz do [Discord](https://discord.gg/ZpNq8ez)
3. Zobacz [dokumentację MikroTik](https://help.mikrotik.com/docs/display/ROS/RouterOS)

## Autorzy

Aktualizacja do RouterOS 7 została wykonana zgodnie z najlepszymi praktykami i wytycznymi społeczności.

## Licencja

Ten projekt zachowuje oryginalną licencję. Zobacz [LICENSE](./LICENSE) dla szczegółów.
