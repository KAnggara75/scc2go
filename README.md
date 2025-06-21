# SCC to Go Viper

[![Go Reference](https://pkg.go.dev/badge/github.com/KAnggara75/scc2go.svg)](https://pkg.go.dev/github.com/KAnggara75/scc2go)

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

## 🤝 Kontribusi
Pull request & masukan sangat diterima!
Silakan fork & buat PR jika ingin menambah fitur.

