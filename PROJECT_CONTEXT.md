# Project Context & Scope

## 1. Project Purpose
`scc2go` adalah Go library client yang dirancang untuk menjembatani ekosistem Java Spring Cloud dengan ekosistem Go. Proyek ini mempermudah aplikasi Go mengambil konfigurasi terpusat dari **Spring Cloud Config Server** dan secara otomatis memetakan seluruh properties hierarkis ke dalam **Viper** registry, sehingga developer Go dapat mengonsumsi konfigurasi menggunakan API standar `viper.Get*()`.

---

## 2. System Boundary & Scope
- **Yang Ditangani oleh Sistem Ini**:
  - Melakukan panggilan HTTP REST ke endpoint Spring Cloud Config Server (`/{application}/{profile}`).
  - Menangani autentikasi HTTP (Basic Auth, Bearer Token via Header `Authorization`).
  - Mendukung bypass TLS verifikasi sertifikat opsional (`InsecureSkipVerify`) untuk kebutuhan development/internal PKI.
  - Mem-parsing respon multi-level `PropertySources` Spring Cloud Config dan memetakannya ke konfigurasi key-value Viper dengan memperhatikan urutan prioritas (*precedence order*).
  - Menyediakan fallback lokal berbasis environment variables OS saat URL kosong atau `"local"`.
- **Yang Berada di Luar Lingkup Sistem (Eksternal)**:
  - Manajemen server dan git backend Spring Cloud Config itu sendiri.
  - Runtime aplikasi konsumen (database pooling, HTTP servers, business logic).

---

## 3. Main Actors & Personas
- **Go Application Developer**: Mengintegrasikan `scc2go` ke dalam aplikasi microservice Go untuk mengambil konfigurasi secara terpusat saat startup/bootstrap.
- **CI/CD Pipeline (GitHub Actions)**: Otomasi pengujian multi-versi Go (1.25.7, 1.26), linter security scanner (`gosec`, `staticcheck`, `govulncheck`), dan tagging rilis otomatis.

---

## 4. Important Domain Concepts & Glossary
- **Spring Cloud Config Server**: Layanan terpusat dalam ekosistem Spring yang menyediakan konfigurasi berbasis HTTP/REST untuk lingkungan terdistribusi.
- **PropertySource**: Unit koleksi key-value konfigurasi di Spring Cloud Config. Sebuah aplikasi dapat memiliki beberapa PropertySource (misalnya spesifik profile `prod`, lalu fallback profile `default`).
- **Precedence Order**: Urutan prioritas property source. Di Spring Cloud Config, item urutan awal memiliki prioritas lebih tinggi daripada item urutan akhir.
- **Viper**: Configuration management library paling populer di Go yang mendukung nested keys, flags, env vars, dan format data seperti JSON/YAML/TOML.
- **Cipher Text**: Format nilai terenkripsi di Spring Cloud Config yang biasanya diawali prefix `{cipher}` untuk menyembunyikan informasi rahasia seperti password.

---

## 5. External Systems & Integrations
- **Spring Cloud Config Server**: Upstream service penyedia konfigurasi.
- **Viper Registry**: Target internal in-memory repository penyimpanan konfigurasi aplikasi.
- **Codecov & GitHub Releases**: Ekosistem CI/CD untuk monitoring test coverage dan deployment tagging.

---

## 6. Runtime Environment & Constraints
- **Language / Runtime**: Go 1.25.0+ (diuji pada Go 1.25.7 dan Go 1.26).
- **OS Support**: Multiplatform (Linux, macOS, Windows).
- **Network Constraints**: Memerlukan konektivitas HTTP/HTTPS ke Spring Cloud Config Server dengan timeout default 5 detik dan 3 kali percobaan ulang (retry).
- **Security Scanners**: Wajib mematuhi audit `gosec`, `govulncheck`, dan `staticcheck` yang dieksekusi pada pre-push dan CI.

---

## 7. Coding Conventions & Standards
- **Linter & Formatting**:
  - `goimports` dengan flag `-local github.com/KAnggara75/scc2go`.
  - `zerolog-use-stringer` untuk konsistensi zerologging.
  - `staticcheck` dan `govulncheck` pada level pre-push git hook.
- **Commit Style**: Mengikuti panduan [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, dsb.).
- **Branching & Release**: Merge ke `main` secara otomatis memicu pembuatan git tag berbasis timestamp waktu WIB: `v0.YY.M-DHHMM`.

---

## 8. Known Limitations & Future Roadmap
- **Error Return**: Saat ini masih berupa silent failure di fungsi public `GetEnv` dan `GetEnvWithDebug`; sedang dijadwalkan untuk mengembalikan eksplisit `error` ([ADR-001](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-001-error-return-propagation-on-scc-fetch-failure)).
- **Custom Viper Instance**: Belum mendukung injeksi langsung ke instance kustom `*viper.Viper` ([ADR-002](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-002-support-for-custom-viper-instance)).
- **Cipher Decryption**: Belum mendukung pemrosesan data `{cipher}` untuk konfigurasi terenkripsi ([ADR-004](file:///Users/i/work/KAnggara75/scc2go/DECISIONS.md#adr-004-roadmap-for-cipher-text-encryptiondecryption)).
- **Context Propagation**: Belum ada dukungan parameter `context.Context` untuk pembatalan request dinamis.
