# Project TODO & Technical Debt

## 1. Immediate Tasks
_Tugas atau perbaikan mendesak yang mempengaruhi reliability, correctness, atau API contract._
- [ ] **[BUG] Error Propagation on SCC Failure**: Perbaiki silent failure pada [`GetEnv`](file:///Users/i/work/KAnggara75/scc2go/scc2go.go#L43) dan [`GetEnvWithDebug`](file:///Users/i/work/KAnggara75/scc2go/scc2go.go#L47). Sediakan mekanisme atau varian fungsi yang mengembalikan `error` eksplisit saat fetch HTTP atau unmarshal JSON gagal ([ADR-001](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-001-error-return-propagation-on-scc-fetch-failure)).
- [ ] **Custom Viper Instance Support**: Refactor penulisan konfigurasi agar tidak terkunci pada singleton global `viper`, melainkan mendukung penyuntikan ke instance `*viper.Viper` kustom ([ADR-002](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-002-support-for-custom-viper-instance)).
- [ ] **Global Zerolog Mutability**: Fungsi `GetEnvWithDebug` memodifikasi global setting zerolog (`zerolog.TimeFieldFormat = zerolog.TimeFormatUnix`) di [`scc2go.go:48`](file:///Users/i/work/KAnggara75/scc2go/scc2go.go#L48), yang berpotensi memiliki efek samping pada format logging di aplikasi utama pengguna.

## 2. Existing Code Annotations (TODO / FIXME)
_Daftar TODO/FIXME yang tercantum langsung di dalam source code._
- _Tidak ditemukan marker TODO, FIXME, HACK, XXX, atau BUG di codebase._
- Anotasi security scanner yang teridentifikasi:
  - [`scc2go.go:129`](file:///Users/i/work/KAnggara75/scc2go/scc2go.go#L129): `// #nosec G402 -- caller explicitly opted in` (Pengecualian gosec untuk opsi `InsecureSkipVerify`).

## 3. Technical Debt & Structural Improvements
_Pekerjaan arsitektural/refactoring jangka panjang untuk maintainability sistem._
- [ ] **Cipher Text Decryption Enhancement**: Implementasikan dukungan dekripsi untuk nilai properti bertanda `{cipher}` dari Spring Cloud Config ([ADR-004](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-004-roadmap-for-cipher-text-encryptiondecryption)).
- [ ] **Missing Context Propagation**: Integrasikan `context.Context` ke dalam pemanggilan Resty HTTP client di `getSCC` agar mendukung timeout dinamis dan graceful cancellation dari caller.
- [x] **Adopsi Semantic Versioning Otomatis**: Pipeline CI telah dioptimalkan mengadopsi standar SemVer 2.0.0 via `anothrNick/github-tag-action` berbasis Conventional Commits ([ADR-005](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-005-automated-semantic-versioning-semver-with-conventional-commits)).
