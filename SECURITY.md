# Security Policy

## Supported Versions

Kami berkomitmen menjaga keamanan library `scc2go` dengan memberikan patch keamanan dan pembaruan dependensi pada versi yang aktif didukung.

| Version / Branch | Go Version | Supported | Notes |
| :--- | :--- | :---: | :--- |
| `main` (Latest Tag) | Go 1.25.x, 1.26.x | :white_check_mark: | Versi rilis tag aktif (skema rilis terotomatisasi) |
| `< v0.26.x` | Go < 1.25 | :x: | Versi terdahulu tidak lagi menerima patch keamanan |

> [!NOTE]
> Proyek ini menggunakan automated release tagging berbasis tanggal/waktu (`v0.YY.M-DHHMM`) saat merge ke branch `main`. Pastikan aplikasi Anda selalu mengonsumsi rilis tag terbaru.

---

## Security Best Practices for Users

Saat menggunakan `scc2go` di lingkungan produksi, perhatikan praktik keamanan berikut:

1. **TLS / SSL Verification**:
   - Secara default, verifikasi sertifikat TLS aktif (`InsecureSkipVerify: false`).
   - Parameter `disableTlsOpt` (atau bypass TLS) hanya ditujukan untuk pengujian lokal/internal PKI staging. **Jangan pernah mengaktifkan bypass TLS di lingkungan produksi**.
   - Codebase mematuhi audit `gosec` dengan anotasi eksplisit `#nosec G402` hanya jika dipanggil dengan parameter tersebut.

2. **Credential & Secret Protection**:
   - Hindari menyimpan token otorisasi atau kredensial rahasia secara hardcoded di kode aplikasi.
   - Gunakan environment variables atau secret vault manager untuk menginjeksikan header `Authorization` ke `scc2go.GetEnv`.

3. **Logging & Sanitization**:
   - Mode `debug=true` (`TraceLevel`) akan mencatat detail setiap konfigurasi key yang dibaca.
   - Jangan aktifkan `debug=true` di lingkungan produksi jika konfigurasi Anda memuat data sensitif.

---

## Reporting a Vulnerability

Kami sangat menghargai kerja sama komunitas dalam melaporkan kerentanan keamanan secara bertanggung jawab (*responsible disclosure*).

### Cara Melaporkan Kerentanan:
1. **GitHub Security Advisory (Disukai)**:
   - Laporkan temuan kerentanan secara privat melalui form [GitHub Security Advisory](https://github.com/KAnggara75/scc2go/security/advisories/new).
   - Tindakan ini menjaga agar detail kerentanan tidak terekspos ke publik sebelum perbaikan tersedia.
2. **Mohon Jangan Buat Public Issue**:
   - Demi keamanan seluruh pengguna, jangan mempublikasikan issue publik atau membuka pull request publik untuk kerentanan keamanan yang belum tertangani.

### Informasi yang Dibutuhkan:
- Deskripsi singkat mengenai kerentanan dan potensi dampaknya (misal: kebocoran token, TLS MITM, denial of service).
- Langkah-langkah reproduksi (PoC atau skrip pengujian).
- Versi `scc2go` dan versi Go runtime yang terdampak.
- Saran atau rekomendasi mitigasi jika tersedia.

### Waktu Respon:
- Kami berusaha mengonfirmasi laporan kerentanan dalam waktu 48 jam.
- Jika kerentanan diverifikasi, tim kami akan mempersiapkan patch perbaikan dan merilisnya bersamaan dengan advisory resmi.
