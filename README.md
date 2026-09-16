# SCC to Go Viper (scc2go)

[![Go Reference](https://pkg.go.dev/badge/github.com/KAnggara75/scc2go.svg)](https://pkg.go.dev/github.com/KAnggara75/scc2go)
[![Go CI/CD](https://github.com/KAnggara75/scc2go/actions/workflows/CI.yaml/badge.svg)](https://github.com/KAnggara75/scc2go/actions/workflows/CI.yaml)
[![codecov](https://codecov.io/gh/KAnggara75/scc2go/branch/main/graph/badge.svg?token=6KSAE8FQH9)](https://codecov.io/gh/KAnggara75/scc2go)
[![Go Report Card](https://goreportcard.com/badge/github.com/KAnggara75/scc2go)](https://goreportcard.com/report/github.com/KAnggara75/scc2go)
[![Latest Release](https://img.shields.io/github/v/release/KAnggara75/scc2go)](https://github.com/KAnggara75/scc2go/releases)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Client library Go minimalis dan cepat untuk mengambil konfigurasi terpusat dari **Spring Cloud Config Server** dan secara otomatis memetakannya ke registry **[Viper](https://github.com/spf13/viper)** (`viper.Get*`).

---

## ✨ Fitur Utama

- 🌐 **Integrasi Spring Cloud Config**: Mengambil konfigurasi JSON multi-profile dari endpoint `/{name}/{profile}`.
- ⚡ **Otomasi Viper Registry**: Memetakan hierarki `PropertySources` Spring Cloud Config langsung ke global Viper registry.
- 🔄 **Preservasi Urutan Prioritas (*Precedence Order*)**: Memproses property sources dengan urutan prioritas yang tepat tanpa menimpa nilai yang sudah ada sebelumnya.
- 💻 **Mode Lokal / Fallback OS Env**: Otomatis beralih ke variabel lingkungan OS saat URL kosong atau bernilai `"local"` (dengan mapping format `EXAMPLE_VAR` ke `example.var`).
- 🛡️ **Toleransi Jaringan & Retries**: Didukung oleh Resty v3 dengan konfigurasi timeout 5 detik dan *automatic retry* hingga 3 kali.
- 🔒 **Otentikasi & Opsi TLS**: Mendukung header otentikasi kustom (Basic / Bearer) dan opsi bypass verifikasi TLS (`InsecureSkipVerify`) untuk kebutuhan dev/staging.
- 🔍 **Structured Observability**: Terintegrasi dengan Zerolog dengan opsi logging debug trace.

---

## 🚀 Instalasi

```bash
go get -u github.com/KAnggara75/scc2go
```

> **Requirements**: Go 1.25+ atau Go 1.26+.

---

## 🔧 Penggunaan

### 1. Basic Bootstrap

Panggil `scc2go.GetEnv` pada fase bootstrap aplikasi (misalnya di `init()` atau sebelum server dijalankan):

```go
package main

import (
	"fmt"
	"os"

	"github.com/KAnggara75/scc2go"
	"github.com/spf13/viper"
)

func init() {
	// Format: scc2go.GetEnv(sccUrl, authHeader, disableTlsOpt...)
	scc2go.GetEnv(os.Getenv("SCC_URL"), os.Getenv("SCC_AUTH"))
}

func main() {
	// Akses konfigurasi langsung lewat Viper
	appName := viper.GetString("app.name")
	port := viper.GetInt("server.port")

	fmt.Printf("Starting %s on port %d...\n", appName, port)
}
```

### 2. Mode Debug / Trace Logging

Gunakan `GetEnvWithDebug` untuk melihat proses ekstraksi setiap key konfigurasi secara transparan via zerolog:

```go
scc2go.GetEnvWithDebug(sccUrl, authHeader, true)
```

### 3. Bypass Verifikasi TLS (Dev / Self-Signed Certs)

Untuk lingkungan pengembangan lokal atau sertifikat staging *self-signed*:

```go
scc2go.GetEnv(sccUrl, authHeader, true) // disableTls = true
```

> [!WARNING]
> Jangan mengaktifkan bypass TLS (`disableTls = true`) di lingkungan produksi.

### 4. Mode Lokal / Fallback OS Env

Jika `SCC_URL` kosong (`""`) atau bernilai `"local"`, `scc2go` otomatis memuat seluruh variabel lingkungan OS ke Viper:

```bash
export APP_PORT=8080
export DATABASE_HOST="localhost"
```

```go
scc2go.GetEnv("local", "")

// Otomatis dipetakan:
// APP_PORT       -> app.port
// DATABASE_HOST  -> database.host
viper.GetInt("app.port")      // 8080
viper.GetString("database.host") // "localhost"
```

---

## 🔒 Otentikasi

Header `auth` dioperasikan langsung sebagai nilai header `Authorization` HTTP:

- **Bearer Token**:
  ```go
  scc2go.GetEnv("https://config-server/app/prod", "Bearer my-secret-token")
  ```
- **HTTP Basic Auth**:
  ```go
  scc2go.GetEnv("https://config-server/app/prod", "Basic dXNlcjpwYXNzd29yZA==")
  ```

---

## 📦 Versioning & Releases

Proyek ini menggunakan **automated release tagging** berbasis waktu WIB (`v0.YY.M-DHHMM`) saat pull request dimerge ke branch `main`, dipadukan dengan standar **[Conventional Commits](https://www.conventionalcommits.org/)**:

- `feat:` — Menambah fungsionalitas baru
- `fix:` — Perbaikan bug
- `docs:` — Pembaruan dokumentasi
- `build:` / `ci:` — Pembaruan dependensi atau pipeline workflow
- `refactor:` — Restrukturisasi kode tanpa mengubah fungsionalitas

Lihat rilis terbaru di [GitHub Releases](https://github.com/KAnggara75/scc2go/releases).

---

## 🧪 Testing

Jalankan test suite dengan race detection dan coverage:

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🛡️ Security & Kebijakan Keamanan

Untuk informasi kebijakan keamanan, versi yang didukung, dan prosedur pelaporan kerentanan secara bertanggung jawab (*responsible disclosure*), silakan baca [SECURITY.md](SECURITY.md).

---

## 🤝 Kontribusi

Pull Request dan masukan sangat diterima!
1. Fork repositori ini.
2. Buat branch fitur (`git checkout -b feat/my-feature`).
3. Pastikan kode lulus formatting dan linter (`go test ./...`, `goimports`).
4. Commit perubahan mengikuti konvensi Conventional Commits (`git commit -m "feat: add feature"`).
5. Push ke remote branch dan ajukan Pull Request.

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah [Mozilla Public License 2.0](LICENSE).
