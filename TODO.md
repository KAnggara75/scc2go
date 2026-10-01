# Project TODO & Technical Debt

## 1. Immediate Tasks
_Tugas atau perbaikan mendesak yang mempengaruhi reliability, correctness, atau API contract._
- [x] **[BUG] Error Propagation on SCC Failure**: Telah disediakan fungsi `Load(...) error` yang mengembalikan `error` eksplisit saat fetch HTTP atau unmarshal JSON gagal ([ADR-001](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-001-error-return-propagation-on-scc-fetch-failure)).
- [x] **Custom Viper Instance Support**: Telah disediakan opsi `WithViper(*viper.Viper)` sehingga konfigurasi dapat disuntikkan ke instance Viper terisolasi ([ADR-002](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-002-support-for-custom-viper-instance)).
- [x] **Global Zerolog Mutability**: Telah dihapus mutasi global `zerolog.TimeFieldFormat = zerolog.TimeFormatUnix` dari loader untuk menghindari efek samping pada logging aplikasi pengguna.

## 2. Existing Code Annotations (TODO / FIXME)
_Daftar TODO/FIXME yang tercantum langsung di dalam source code._
- _Tidak ditemukan marker TODO, FIXME, HACK, XXX, atau BUG di codebase._
- Anotasi security scanner yang teridentifikasi:
  - [`scc2go.go:210`](file:///Users/i/work/KAnggara75/scc2go/scc2go.go#L210): `// #nosec G402 -- caller explicitly opted in` (Pengecualian gosec untuk opsi `InsecureSkipVerify`).

## 3. Technical Debt & Structural Improvements
_Pekerjaan arsitektural/refactoring jangka panjang untuk maintainability sistem._
- [ ] **Cipher Text Decryption Enhancement**: Implementasikan dukungan dekripsi untuk nilai properti bertanda `{cipher}` dari Spring Cloud Config ([ADR-004](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-004-roadmap-for-cipher-text-encryptiondecryption)).
- [x] **Context Propagation**: Integrasikan `context.Context` ke dalam pemanggilan Resty HTTP client via `WithContext(ctx)` agar mendukung timeout dinamis dan graceful cancellation dari caller.
- [x] **Adopsi Semantic Versioning Otomatis**: Pipeline CI telah dioptimalkan mengadopsi standar SemVer 2.0.0 via `anothrNick/github-tag-action` berbasis Conventional Commits ([ADR-005](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-005-automated-semantic-versioning-semver-with-conventional-commits)).
