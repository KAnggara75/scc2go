# Codebase Navigation Map

## `github.com/KAnggara75/scc2go` (Root Module)
- **Responsibility**: Mengambil konfigurasi aplikasi dari Spring Cloud Config Server via REST API (atau OS environment variables sebagai fallback/local mode) dan menyimpannya langsung ke Viper global registry (`github.com/spf13/viper`).
- **Entry / Key Files**:
  - [`scc2go.go`](file:///Users/i/work/KAnggara75/scc2go/scc2go.go) — Implementasi utama: struct data DTO Spring Cloud Config (`springCloudConfig`, `propertySource`), fungsi public `GetEnv` dan `GetEnvWithDebug`, HTTP client via Resty, parser property sources, dan fallback `loadFromEnv`.
  - [`scc2go_test.go`](file:///Users/i/work/KAnggara75/scc2go/scc2go_test.go) — Comprehensive test suite mencakup unit tests, mock HTTP server (`httptest`), edge case auth header, TLS bypass, dan precedence pengujian key mapping.
  - [`go.mod`](file:///Users/i/work/KAnggara75/scc2go/go.mod) — Definisi modul Go (Go 1.25.0) dan dependencies utama.
  - [`.github/workflows/CI.yaml`](file:///Users/i/work/KAnggara75/scc2go/.github/workflows/CI.yaml) — CI pipeline untuk multi-version test (Go 1.25.7 & Go 1.26), race detection, coverage upload (Codecov), linter suites (gosec, staticcheck, govulncheck, pre-commit), dan automated release tagging format `v0.yy.m-dHHMM`.
  - [`.pre-commit-config.yaml`](file:///Users/i/work/KAnggara75/scc2go/.pre-commit-config.yaml) — Standar git hook lokal dan CI: whitespace, merge-conflict check, `goimports`, `go vet`, `go mod tidy`, `zerolog-use-stringer`, serta pre-push checks (`govulncheck`, `staticcheck`, `gosec`).
- **Dependencies**:
  - `resty.dev/v3 v3.0.0-beta.6` — HTTP client dengan automatic retry, timeout, dan custom TLS config.
  - `github.com/spf13/viper v1.21.0` — Global configuration registry untuk Go applications.
  - `github.com/rs/zerolog v1.35.0` — Structured JSON/console logging dengan level switching (`InfoLevel` vs `TraceLevel`).
- **Consumers**:
  - Library ini dikonsumsi oleh service / aplikasi Go microservice yang berjalan berdampingan dengan Spring Cloud ecosystem dan memerlukan central configuration management.
- **External Integrations**:
  - **Spring Cloud Config Server**: Endpoint REST API (`/application/profile`) yang mereturn JSON representasi konfigurasi terpusat.
- **Key Notes**:
  - Property precedence: Array `PropertySources` diproses secara terbalik (`for i := len(scc.PropertySources) - 1; i >= 0; i--`), dipadukan dengan `setIfNotExists` (`if viper.IsSet(k) { return }`), sehingga property source dengan prioritas tertinggi dipertahankan dan nilai yang sudah ada sebelumnya tidak ditimpa.
  - Silent error policy: Kesalahan HTTP atau JSON unmarshal saat ini di-log menggunakan `logger.Error()` tanpa me-return error atau panic.
