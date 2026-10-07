# Core Plugins, Libraries, and Dependencies Specification

Dokumen ini mendefinisikan seluruh pustaka (_dependencies_), driver, dan utilitas inti yang digunakan dalam proyek Go (Golang) + Gin untuk sistem WebRTC Signaling dan REST API.

---

## 1. Runtime & Core Framework

| Package       | Modul / Path               | Versi Disarankan  | Deskripsi & Peran                                                                                                        |
| :------------ | :------------------------- | :---------------- | :----------------------------------------------------------------------------------------------------------------------- |
| **Go**        | Runtime                    | `1.22+` / `1.23+` | Runtime utama dengan dukungan enhanced routing dan manajemen memori efisien.                                             |
| **Gin Gonic** | `github.com/gin-gonic/gin` | `^v1.10.0`        | HTTP web framework berkecepatan tinggi dengan routing berbasis HttpRouter, middleware chain, dan JSON serializer bawaan. |

---

## 2. Real-Time & WebRTC Communication

| Package                          | Modul / Path                   | Peran & Alasan Pemilihan                                                                                                                                                           |
| :------------------------------- | :----------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Gorilla WebSocket**            | `github.com/gorilla/websocket` | Standar industri untuk implementasi protokol WebSocket di Go. Mendukung penanganan konkurensi goroutine per-koneksi (_readPump_ / _writePump_) dan pengontrolan buffer buffer I/O. |
| **Pion WebRTC (Opsional/Addon)** | `github.com/pion/webrtc/v4`    | Implementasi WebRTC native Go murni. Digunakan jika di masa depan backend bertindak sebagai SFU/MCU (bukan sekadar Signaling Server).                                              |
| **Pion TURN / STUN**             | `github.com/pion/turn/v3`      | Generator kredensial TURN ephemeral HMAC-SHA1 dan utilitas verifikasi konektivitas ICE.                                                                                            |

---

## 3. Database Layer (MySQL)

| Package                | Modul / Path                                                                    | Peran & Alasan Pemilihan                                                                                                                                                                                        |
| :--------------------- | :------------------------------------------------------------------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Go-MySQL-Driver**    | `github.com/go-sql-driver/mysql`                                                | Driver native MySQL standar untuk `database/sql`. Mendukung parsing tipe data `DATETIME`, connection pooling, dan performa tinggi.                                                                              |
| **sqlx** atau **GORM** | `github.com/jmoiron/sqlx` _(Rekomendasi)_ <br>`gorm.io/gorm` _(Alternatif ORM)_ | **sqlx:** Abstraksi ringan di atas `database/sql` dengan fitur struct scanning yang cepat tanpa overhead ORM berat.<br>**GORM:** Dipilih jika membutuhkan migrasi otomatis dan relasi entitas deklaratif cepat. |
| **Golang-Migrate**     | `github.com/golang-migrate/migrate/v4`                                          | CLI & Go library untuk mengeksekusi file migrasi skema SQL versi terstruktur (`up.sql` / `down.sql`).                                                                                                           |

---

## 4. Security, Auth, & Cryptography

| Package               | Modul / Path                   | Peran & Alasan Pemilihan                                                                                                       |
| :-------------------- | :----------------------------- | :----------------------------------------------------------------------------------------------------------------------------- |
| **Golang-JWT**        | `github.com/golang-jwt/jwt/v5` | Pembuatan, penandatanganan (_signing_), dan verifikasi token autentikasi JWT (HS256/RS256) untuk handshake HTTP dan WebSocket. |
| **Crypto (x/crypto)** | `golang.org/x/crypto/bcrypt`   | Hashing kata sandi pengguna secara aman dengan algoritma BCrypt yang tahan serangan brute-force.                               |
| **Google UUID**       | `github.com/google/uuid`       | Pembuatan UUID v4 untuk identitas entitas unik (_User ID_, _Room ID_, _Session ID_).                                           |

---

## 5. Configuration & Environment

| Package                | Modul / Path               | Peran & Alasan Pemilihan                                                                                   |
| :--------------------- | :------------------------- | :--------------------------------------------------------------------------------------------------------- |
| **Godotenv**           | `github.com/joho/godotenv` | Membaca file konfigurasi `.env` ke environment variable sistem saat mode pengembangan lokal.               |
| **Viper** _(Opsional)_ | `github.com/spf13/viper`   | Pengelolaan konfigurasi dinamis berbasis JSON/YAML/ENV untuk lingkungan multi-stage (staging, production). |

---

## 6. Observability & Logging

| Package                  | Modul / Path                                | Peran & Alasan Pemilihan                                                                                                                    |
| :----------------------- | :------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------ |
| **Zap** atau **Zerolog** | `go.uber.org/zap` / `github.com/rs/zerolog` | Structured JSON logger dengan alokasi memori mendekati nol (_zero-allocation_). Sangat krusial untuk mencatat jutaan pesan event WebSocket. |

---

## 7. Development Tools & Utility CLI

- **Air (`github.com/air-verse/air`):** Utilitas live-reload untuk Go apps saat kode diubah di environment WSL.
- **Golangci-lint:** Linter agregator untuk menjaga kualitas dan kepatuhan standar sintaks Go.
- **Docker & Docker Compose:** Menjalankan instance MySQL 8.x dan Coturn server secara lokal.

---

## 8. Ringkasan Perintah Instalasi Awal

```bash
# Inisialisasi dependensi inti
go get -u github.com/gin-gonic/gin
go get -u github.com/gorilla/websocket
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jmoiron/sqlx
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/google/uuid
go get -u github.com/joho/godotenv
go get -u go.uber.org/zap
```
