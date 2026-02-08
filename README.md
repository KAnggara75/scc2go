# SCC to Go Viper

[![Go Reference](https://pkg.go.dev/badge/github.com/KAnggara75/scc2go.svg)](https://pkg.go.dev/github.com/KAnggara75/scc2go)
[![Go CI/CD](https://github.com/KAnggara75/scc2go/actions/workflows/go.yml/badge.svg)](https://github.com/KAnggara75/scc2go/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/KAnggara75/scc2go/branch/main/graph/badge.svg)](https://codecov.io/gh/KAnggara75/scc2go)
[![Go Report Card](https://goreportcard.com/badge/github.com/KAnggara75/scc2go)](https://goreportcard.com/report/github.com/KAnggara75/scc2go)
[![Latest Release](https://img.shields.io/github/v/release/KAnggara75/scc2go)](https://github.com/KAnggara75/scc2go/releases)

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FKAnggara75%2Fscc2go.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FKAnggara75%2Fscc2go?ref=badge_large)

Package Go untuk mengambil konfigurasi aplikasi dari Spring Cloud Config Server dan langsung menyimpannya ke [Viper](https://github.com/spf13/viper) untuk kemudahan akses di aplikasi Go.


## ✨ Fitur
- Fetch konfigurasi langsung dari Spring Cloud Config Server (/application/profile)
- Support authentication (username/password)
- Integrasi otomatis dengan Viper
- Mendukung format konfigurasi: YAML, JSON, Properties (auto mapping ke Viper)
- Cocok untuk aplikasi Go yang ingin pakai central config seperti Spring di ekosistem Java


## 🚀 Instalasi

```bash
go get -u github.com/KAnggara75/scc2go
```

## 🔧 Penggunaan

### 1. Import package

```go
import (
	"github.com/KAnggara75/scc2go"
)
```

### 2. Inisialisasi dan ambil konfigurasi

```go
func init() {
	scc2go.GetEnv(os.Getenv("SCC_URL"), os.Getenv("AUTH"))
}
```

### 3. Akses konfigurasi di Viper

```go
package config

import (
	"github.com/spf13/viper"
)

func GetDBConn() string {
	return viper.GetString("db.myapp.host")
}
```

### 4. Contoh penggunaan

```go
package main
import (
    "fmt"
    "os"

    "github.com/KAnggara75/scc2go"
    "github.com/spf13/viper"
)

func init() {
	scc2go.GetEnv(os.Getenv("SCC_URL"), os.Getenv("AUTH"))
}

func main() {
	deviceId := viper.GetString("db.myapp.host")
	fmt.Println("Device ID:", deviceId)
}
```

## 📝 Konfigurasi yang Didukung
- Format: .yaml, .yml, .json, .properties
- Key mapping otomatis dari response Spring Cloud Config ke viper

## 🔒 Autentikasi
- Tambahkan field Auth jika server butuh basic auth.

```go
scc2go.GetEnv("SCC_URL", "AUTH")
```

## 📦 Versioning

This project follows [Semantic Versioning](https://semver.org/) with automated releases based on [Conventional Commits](https://www.conventionalcommits.org/).

For detailed information about versioning and commit message format, see [VERSIONING.md](VERSIONING.md).

### Quick Reference

- `feat:` - New feature (minor version bump)
- `fix:` - Bug fix (patch version bump)
- `feat!:` or `BREAKING CHANGE:` - Breaking change (major version bump)

## 🧪 Testing

Run tests with coverage:

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Current test coverage: **97.1%**

## 🤝 Kontribusi

Pull request & masukan sangat diterima!

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes using [Conventional Commits](https://www.conventionalcommits.org/)
   - `feat: add new feature`
   - `fix: resolve bug`
   - `docs: update documentation`
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

See [VERSIONING.md](VERSIONING.md) for commit message guidelines.
