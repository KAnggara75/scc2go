# Architecture Decision Records (ADR)

## ADR-001: Error Return Propagation on SCC Fetch Failure

- **Status**: Accepted
- **Date**: 2026-09-07
- **Source**: Developer interview
- **Context**:
  Fungsi `GetEnv` dan `GetEnvWithDebug` saat ini hanya mencatat error ke zerolog (`logger.Error()`) dan melakukan silent return jika HTTP request ke Spring Cloud Config gagal atau payload JSON corrupt. Developer mengonfirmasi bahwa perilaku ini merupakan bug/defisiensi. Caller aplikasi membutuhkan indikator error yang jelas untuk menentukan apakah bootstrap aplikasi boleh berlanjut atau harus fail-fast.
- **Decision**:
  Akan disediakan mekanisme/fungsi yang mengembalikan nilai `error` eksplisit saat pengambilan atau pemrosesan konfigurasi SCC gagal, sehingga caller dapat menangani kegagalan konfigurasi secara deterministik.
- **Consequences**:
  - Positif: Penanganan error menjadi transparan dan sesuai idiomatic Go (`if err != nil`).
  - Negatif/Trade-off: Membutuhkan penyesuaian signature atau penyediaan fungsi baru untuk menjaga backward compatibility dengan caller yang menggunakan signature lama tanpa return value.

---

## ADR-002: Support for Custom Viper Instance

- **Status**: Accepted
- **Date**: 2026-09-07
- **Source**: Developer interview
- **Context**:
  Saat ini `scc2go` langsung melakukan injeksi konfigurasi ke global singleton `github.com/spf13/viper`. Ini membatasi modularitas dan membuat pengujian isolasi serta arsitektur multi-instance/multi-tenant menjadi sulit.
- **Decision**:
  Mendukung injeksi konfigurasi ke instance kustom `*viper.Viper` (misalnya melalui parameter instance kustom atau builder options), sementara tetap mempertahankan opsi default untuk global viper jika diperlukan.
- **Consequences**:
  - Positif: Meningkatkan modularitas, testability, dan fleksibilitas integrasi di aplikasi Go yang menggunakan instance viper terisolasi.
  - Negatif/Trade-off: API surface library bertambah, membutuhkan refactoring pada fungsi mapping internal (`setIfNotExists`).

---

## ADR-003: Timestamp-Based Versioning Scheme

- **Status**: Superseded by ADR-005
- **Date**: 2026-09-07
- **Source**: Developer interview & CI workflow evidence (`.github/workflows/CI.yaml`)
- **Context**:
  Meskipun README menyebutkan Conventional Commits / Semver standar, pipeline CI mengotomatisasi pembuatan tag rilis dengan format tanggal dan jam berbasis zona waktu WIB: `v0.${yy}.${m}-${d}${HH}${MM}`.
- **Decision**:
  Skema versioning utama proyek untuk saat ini secara resmi menggunakan format rilis berbasis tanggal/waktu (WIB): `v0.YY.M-DHHMM` yang di-generate otomatis oleh GitHub Actions saat merge ke branch `main`.
- **Consequences**:
  - Positif: Rilis sepenuhnya terotomatisasi tanpa memerlukan manual version bump.
  - Negatif/Trade-off: Berbeda dari SemVer murni (MAJOR.MINOR.PATCH), sehingga konsumsi library bergantung pada tag timestamp tersebut.

---

## ADR-004: Roadmap for Cipher Text Encryption/Decryption

- **Status**: Accepted
- **Date**: 2026-09-07
- **Source**: Developer interview
- **Context**:
  Spring Cloud Config mendukung properti terenkripsi (format `{cipher}...`). Saat ini `scc2go` belum mendukung dekripsi lokal untuk nilai terenkripsi yang dikembalikan oleh SCC Server (atau integrasi decrypt endpoint).
- **Decision**:
  Akan ditambahkan enhancement untuk menangani enkripsi/dekripsi nilai rahasia berformat cipher text, baik dengan memanfaatkan endpoint decrypt Spring Cloud Config Server ataupun dekripsi lokal jika key disediakan.
- **Consequences**:
  - Positif: Memungkinkan aplikasi Go menangani konfigurasi rahasia (credentials, database passwords) yang terenkripsi di Spring Cloud Config.
  - Negatif/Trade-off: Menambah dependensi kriptografi atau logic request tambahan ke endpoint dekripsi SCC server.

---

## ADR-005: Automated Semantic Versioning (SemVer) with Conventional Commits

- **Status**: Accepted
- **Date**: 2026-09-16
- **Source**: Developer request (`optimalkan CI.yaml buat versioning semver`) & CI workflow (`.github/workflows/CI.yaml`)
- **Context**:
  Sebelumnya proyek menggunakan skema penomoran berbasis timestamp WIB (`v0.YY.M-DHHMM`) via shell script di GitHub Actions (ADR-003). Skema ini menyulitkan package manager Go dan tool dependensi eksternal untuk mengidentifikasi level perubahan semantik (major/minor/patch/breaking change) secara deterministik.
- **Decision**:
  Menggantikan skema timestamp dengan Semantic Versioning standar (`vMAJOR.MINOR.PATCH`) yang sepenuhnya terotomatisasi menggunakan `anothrNick/github-tag-action@v1`. Logika kenaikan versi ditentukan langsung dari analisis Conventional Commits (`feat:` -> minor, `fix:` -> patch, `BREAKING CHANGE:` / `!` -> major).
- **Consequences**:
  - Positif: Sepenuhnya mematuhi spesifikasi SemVer 2.0.0 dan ekosistem Go modules (`go get`, `go proxy`), otomatisasi rilis tetap berjalan tanpa intervensi manual.
  - Negatif/Trade-off: Memerlukan kepatuhan konsisten dari kontributor dalam menulis prefix Conventional Commits yang tepat pada setiap commit / pull request.

