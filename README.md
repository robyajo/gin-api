# Dokumentasi Arsitektur & Panduan Proyek (gin-api)

Backend Service berbasis **Go (Golang)** menggunakan framework **Gin** dan **Gorilla WebSocket** untuk menangani RESTful API, autentikasi berbasis peran (RBAC), ruang pertemuan, dan **WebRTC Signaling Server** berlatensi rendah untuk layanan komunikasi tatap muka online warga.

---

## 1. Daftar Isi
- [1. Daftar Isi](#1-daftar-isi)
- [2. Pohon Direktori Lengkap](#2-pohon-direktori-lengkap)
- [3. Penjelasan Detail Fungsi & Kegunaan Setiap Folder](#3-penjelasan-detail-fungsi--kegunaan-setiap-folder)
  - [3.1. Folder `cmd/` (Application Entry Point)](#31-folder-cmd-application-entry-point)
  - [3.2. Folder `config/` (Konfigurasi Terpusat)](#32-folder-config-konfigurasi-terpusat)
  - [3.3. Folder `internal/` (Core Application & Business Logic)](#33-folder-internal-core-application--business-logic)
    - [3.3.1. `internal/delivery/` (Lapisan Transport / I/O)](#331-internaldelivery-lapisan-transport--io)
    - [3.3.2. `internal/service/` (Lapisan Logika Bisnis Murni)](#332-internalservice-lapisan-logika-bisnis-murni)
    - [3.3.3. `internal/repository/` (Lapisan Akses Database)](#333-internalrepository-lapisan-akses-database)
    - [3.3.4. `internal/model/` (Entitas Basis Data & DTO)](#334-internalmodel-entitas-basis-data--dto)
  - [3.4. Folder `pkg/` (Pustaka Pendukung / Shared Packages)](#34-folder-pkg-pustaka-pendukung--shared-packages)
  - [3.5. Folder `scripts/` (Migrasi & Automasi)](#35-folder-scripts-migrasi--automasi)
  - [3.6. Folder `releases/` (Catatan Rilis Publik Harian)](#36-folder-releases-catatan-rilis-publik-harian)
  - [3.7. Folder `.agents/` (Katalog Skills AI Agent)](#37-folder-agents-katalog-skills-ai-agent)
  - [3.8. Folder `DOCS/` (Dokumentasi Teknis & Arsitektur)](#38-folder-docs-dokumentasi-teknis--arsitektur)
  - [3.9. Berkas Konfigurasi di Root](#39-berkas-konfigurasi-di-root)
- [4. Struktur 4 Role Sistem](#4-struktur-4-role-sistem)
- [5. Panduan Menjalankan Aplikasi](#5-panduan-menjalankan-aplikasi)
- [6. Dokumentasi Endpoint API](#6-dokumentasi-endpoint-api)

---

## 2. Pohon Direktori Lengkap

```text
gin-api/
├── .agents/                                # Konfigurasi kustomisasi agen AI Antigravity
│   └── skills/                             # Modul skill operasional agen
│       ├── api-testing/                    # Panduan pengujian REST & WebSocket
│       ├── database-migration/             # Panduan DDL MySQL, UUID, dan transaksi
│       ├── release-notes/                  # Protokol pencatatan rilis harian dwibahasa
│       └── webrtc-signaling/               # Panduan Hub, SDP/ICE, dan TURN ephemeral
├── cmd/                                    # Titik masuk utama eksekusi program
│   └── api/
│       └── main.go                         # Titik masuk utama server backend
├── config/                                 # Pengelolaan variabel konfigurasi aplikasi
│   └── config.go                           # Pemetaan file .env ke struct Config
├── DOCS/                                   # Dokumentasi spesifikasi sistem
│   ├── erd.md                              # Entity Relationship Diagram & skema MySQL
│   ├── plugins.md                          # Dokumentasi plugin & ekstensi
│   ├── prd.md                              # Product Requirement Document
│   └── structure.md                        # Blueprint arsitektur modular
├── internal/                               # Kode aplikasi privat (tidak diimpor dari luar)
│   ├── delivery/                           # Lapisan penerima request (HTTP / WebSocket)
│   │   ├── http/                           # Penanganan protokol HTTP
│   │   │   ├── handler/                    # Controller penerima request & response DTO
│   │   │   │   ├── response.go             # Helper standar amplop respons JSON
│   │   │   │   ├── role_handler.go         # Endpoint CRUD Role
│   │   │   │   └── user_handler.go         # Endpoint CRUD User
│   │   │   ├── middleware/                 # Middleware Gin (Autentikasi JWT, CORS, dll.)
│   │   │   │   └── auth_middleware.go      # Validasi Bearer Token JWT
│   │   │   └── router.go                   # Pendaftaran rute API Gin (/api/v1/...)
│   │   └── ws/                             # Penanganan protokol WebSocket
│   │       ├── client.go                   # Siklus hidup client WebSocket (read/write pump)
│   │       ├── hub.go                      # Pengelola ruangan, pendaftaran, dan broadcast
│   │       └── message.go                  # Skema amplop pesan signaling WebRTC (SDP/ICE)
│   ├── model/                              # Definisi struct entitas data & DTO request
│   │   ├── role.go                         # Struct Role, CreateRoleRequest, UpdateRoleRequest
│   │   ├── room.go                         # Struct Room pertemuan konsultasi
│   │   └── user.go                         # Struct User, UserResponse, Create/UpdateUser
│   ├── repository/                         # Lapisan akses data (SQL Query ke MySQL)
│   │   ├── role_repository.go              # Query database tabel roles
│   │   ├── room_repository.go              # Query database tabel rooms
│   │   └── user_repository.go              # Query database tabel users (join roles)
│   └── service/                            # Lapisan logika bisnis (Business Logic murni)
│       ├── role_service.go                 # Aturan bisnis role (validasi nama, proteksi sistem)
│       ├── turn_service.go                 # Pembuatan kredensial ephemeral TURN (HMAC-SHA1)
│       └── user_service.go                 # Aturan bisnis user (hash bcrypt, verifikasi role)
├── pkg/                                    # Pustaka utilitas yang dapat digunakan ulang
│   ├── database/                           # Konektor basis data
│   │   ├── migration.go                    # Auto-migration skema & default seeder
│   │   └── mysql.go                        # Inisialisasi koneksi pool sqlx MySQL
│   └── utils/                              # Fungsi bantuan teknis
│       └── password.go                     # Utilitas enkripsi & verifikasi hash bcrypt
├── releases/                               # Riwayat rilis publik aplikasi (1 JSON/hari)
│   └── 2026-10-08-fondasi-layanan-komunikasi.json
├── scripts/                                # Skrip pemeliharaan sistem
│   └── migrations/                         # Skrip DDL MySQL up/down terurut
│       ├── 000001_init_schema.down.sql     # Skrip rollback pembatalan tabel
│       └── 000001_init_schema.up.sql       # Skrip pembentukan tabel & indeks
├── .air.toml                               # Konfigurasi live-reloading server Air
├── .env                                    # Konfigurasi lingkungan aktif lokal
├── .env.example                            # Templat konfigurasi lingkungan
├── .gitignore                              # Berkas yang diabaikan Git
├── AGENTS.md                               # Panduan aturan baku untuk AI Agent & Pengembang
├── go.mod                                  # Definisi modul dan dependensi Go
├── go.sum                                  # Checksum integritas paket dependensi Go
├── READMI.md                               # Dokumentasi lengkap proyek ini
└── setup.sh                                # Skrip otomasi pembentukan struktur awal
```

---

## 3. Penjelasan Detail Fungsi & Kegunaan Setiap Folder

Aplikasi ini mengadopsi **Clean Layered Architecture (Arsitektur Berlapis Bersih)** dengan arah ketergantungan satu arah:
```text
[ Delivery (HTTP / WS) ] ──> [ Service (Logika Bisnis) ] ──> [ Repository (Akses DB) ]
```

### 3.1. Folder `cmd/` (Application Entry Point)
Folder ini memuat titik awal eksekusi (*entry point*) program.
* **`cmd/api/main.go`**:
  - Menginisialisasi pembacaan file konfigurasi lingkungan (`config.LoadConfig()`).
  - Membuka koneksi pool ke MySQL (`database.InitMySQL()`).
  - Menjalankan migrasi otomatis dan seeding akun awal (`database.RunMigrations()`).
  - Menghubungkan dependensi secara terstruktur (*Dependency Injection*): instansiasi Repository $\rightarrow$ diteruskan ke Service $\rightarrow$ diteruskan ke Handler $\rightarrow$ didaftarkan ke Router.
  - Menjalankan server HTTP Gin di port yang dikonfigurasikan.

### 3.2. Folder `config/` (Konfigurasi Terpusat)
Folder ini bertanggung jawab mengisolasi seluruh pembacaan konfigurasi sistem.
* **`config/config.go`**:
  - Membaca file `.env` menggunakan library `godotenv`.
  - Memetakan variabel lingkungan (`APP_PORT`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`, `TURN_SECRET`, `TURN_HOST`) ke dalam struct `Config` dengan nilai fallback default yang aman.

### 3.3. Folder `internal/` (Core Application & Business Logic)
Folder `internal/` adalah konvensi resmi Go di mana kode di dalamnya **hanya dapat diakses oleh modul internal ini** dan tidak dapat diimpor oleh modul luar. Ini menjaga integritas arsitektur.

#### 3.3.1. `internal/delivery/` (Lapisan Transport / I/O)
Lapisan ini berhadapan langsung dengan client (aplikasi web/mobile/frontend). Tugasnya hanya menerima data dari request, memvalidasi formatnya, memanggil Service, dan mengembalikan respons:
* **`internal/delivery/http/`**:
  - **`router.go`**: Mengonfigurasi engine Gin, middleware CORS global, endpoint `/health`, dan mengelompokkan grup rute `/api/v1/roles` serta `/api/v1/users`.
  - **`handler/`**:
    - **`response.go`**: Standarisasi format amplop JSON (`code`, `status`, `message`, `data`, `errors`) serta format pagination (`current_page`, `per_page`, `total`, `total_pages`).
    - **`role_handler.go`**: Menangani permintaan HTTP REST untuk operasi CRUD Role (`GetRoles`, `GetRoleByID`, `CreateRole`, `UpdateRole`, `DeleteRole`).
    - **`user_handler.go`**: Menangani permintaan HTTP REST untuk operasi CRUD User (`GetUsers` dengan query pagination & pencarian, `GetUserByID`, `CreateUser`, `UpdateUser`, `DeleteUser`).
  - **`middleware/`**:
    - **`auth_middleware.go`**: Tempat meletakkan middleware verifikasi Bearer Token JWT sebelum request diizinkan mengakses endpoint tertentu.
* **`internal/delivery/ws/`**:
  - **`hub.go`**: Mengelola daftar ruangan (`Room`), client yang sedang terhubung, serta saluran channel pendaftaran (`register`), pelepasan (`unregister`), dan penyiaran pesan (`broadcast`). Dilindungi `sync.RWMutex` agar aman dari *race condition*.
  - **`client.go`**: Mengelola koneksi satu pengguna pada WebSocket dengan 2 goroutine terpisah: `readPump()` untuk membaca pesan masuk dan `writePump()` untuk mengirimkan pesan keluar dengan mekanisme detak jantung (*heartbeat ping/pong*).
  - **`message.go`**: Definisi format amplop pesan standar WebRTC Signaling (`type`, `room_id`, `sender_id`, `target_id`, `payload`) untuk pertukaran penawaran video (*SDP Offer*), jawaban (*SDP Answer*), dan rute jaringan (*ICE Candidates*).

#### 3.3.2. `internal/service/` (Lapisan Logika Bisnis Murni)
Lapisan ini berisi seluruh aturan bisnis sistem. **Sesuai standar AGENTS.md, lapisan ini DILARANG mengimpor package `gin` atau `*gin.Context`** agar logika bisnis tetap independen dan mudah diuji secara unit:
* **`role_service.go`**:
  - Memastikan nama role tidak duplikat.
  - Memproteksi 4 role bawaan sistem (`superadmin`, `admin`, `operator`, `user`) agar tidak bisa dihapus.
  - Memastikan role yang masih memiliki pengguna terdaftar tidak dapat dihapus (*referential integrity*).
* **`user_service.go`**:
  - Memverifikasi keunikan alamat email pengguna baru/update.
  - Memvalidasi keberadaan `role_id` yang dipilih di database.
  - Mengenkripsi kata sandi menggunakan hash `bcrypt` sebelum disimpan.
  - Menghasilkan ID unik UUID v4 untuk setiap user baru.
  - Mengonversi data entitas internal ke `UserResponse` DTO yang aman (tanpa membocorkan hash password).
* **`turn_service.go`**:
  - Menghasilkan kredensial sementara (*ephemeral credentials*) server TURN menggunakan algoritma `HMAC-SHA1` dengan batasan waktu kedaluwarsa (*timestamp expiration*) untuk menembus firewall dan *Symmetric NAT*.

#### 3.3.3. `internal/repository/` (Lapisan Akses Database)
Lapisan ini bertanggung jawab langsung mengeksekusi instruksi SQL ke MySQL menggunakan library `sqlx`:
* **`role_repository.go`**: Eksekusi SQL `SELECT`, `INSERT`, `UPDATE`, `DELETE` pada tabel `roles`, serta fungsi pengecekan relasi `CountUsersByRoleID`.
* **`user_repository.go`**: Eksekusi SQL pada tabel `users` dengan query `JOIN roles` untuk mendapatkan nama role, serta mendukung filter pencarian teks (*search*) dan pembatasan halaman (*LIMIT & OFFSET*).
* **`room_repository.go`**: Tempat penyimpanan kueri manajemen ruangan pertemuan tatap muka.

#### 3.3.4. `internal/model/` (Entitas Basis Data & DTO)
Mendefinisikan tipe data struct untuk pemetaan tabel database dan validasi payload HTTP:
* **`role.go`**: Struct entitas `Role`, payload pembuatan `CreateRoleRequest`, dan `UpdateRoleRequest`.
* **`user.go`**: Struct entitas `User`, struct DTO aman `UserResponse`, payload pendaftaran `CreateUserRequest`, dan payload pembaruan `UpdateUserRequest`.
* **`room.go`**: Struct entitas `Room` dan partisipan ruang konsultasi.

### 3.4. Folder `pkg/` (Pustaka Pendukung / Shared Packages)
Folder untuk modul pembantu yang terisolasi dan dapat digunakan di berbagai bagian aplikasi:
* **`database/`**:
  - **`mysql.go`**: Membuka koneksi pool MySQL, mengatur `MaxOpenConns (50)`, `MaxIdleConns (10)`, dan batas masa aktif koneksi (`ConnMaxLifetime`).
  - **`migration.go`**: Memastikan skema tabel (`roles`, `users`, `rooms`, `room_participants`, `call_logs`) terbuat secara otomatis saat server dinyalakan, sekaligus menanamkan (*seeding*) 4 role sistem dan akun awal superadmin.
* **`utils/`**:
  - **`password.go`**: Fungsi `HashPassword()` dan `CheckPasswordHash()` menggunakan pustaka resmi `golang.org/x/crypto/bcrypt`.

### 3.5. Folder `scripts/` (Migrasi & Automasi)
Menampung berkas skrip operasional:
* **`migrations/`**:
  - **`000001_init_schema.up.sql`**: Berkas SQL DDL untuk membuat seluruh tabel dengan standar InnoDB, charset `utf8mb4_unicode_ci`, dan foreign keys.
  - **`000001_init_schema.down.sql`**: Berkas SQL untuk membatalkan (*drop table*) jika ingin melakukan rollback.
* **`setup.sh`**: Skrip bash pembentuk struktur direktori awal proyek.

### 3.6. Folder `releases/` (Catatan Rilis Publik Harian)
Folder riwayat rilis publik aplikasi. Mematuhi aturan ketat **1 berkas JSON per hari kerja** dalam format dwibahasa (`id` & `en`). Narasi di dalam berkas ini ditulis dengan bahasa pelayanan warga yang ramah publik dan **bebas dari istilah teknis internal/nama library** untuk melindungi keamanan sistem (*Security-Safe & Citizen-Facing*).

### 3.7. Folder `.agents/` (Katalog Skills AI Agent)
Memuat instruksi dan panduan operasional (*runbook*) untuk AI Agent (Antigravity) agar dapat bekerja secara konsisten pada proyek ini:
* **`skills/release-notes/`**: Panduan konsolidasi rilis harian dwibahasa dan kamus bahasa aman publik.
* **`skills/webrtc-signaling/`**: Panduan implementasi WebSocket Hub, amplop pesan SDP/ICE, dan TURN ephemeral.
* **`skills/database-migration/`**: Standar migrasi skema MySQL, transaksi, dan penanganan UUID.
* **`skills/api-testing/`**: Prosedur pengujian endpoint REST dan WebSocket menggunakan CLI.

### 3.8. Folder `DOCS/` (Dokumentasi Teknis & Arsitektur)
Dokumentasi teknis pendukung:
* **`prd.md`**: Persyaratan produk, target latensi (<50ms), dan persona pengguna layanan.
* **`erd.md`**: Diagram relasi entitas basis data (Mermaid) dan optimasi indeks.
* **`structure.md`**: Detail arsitektur modular sistem.
* **`plugins.md`**: Dokumentasi integrasi plugin dan ekstensi pihak ketiga.

### 3.9. Berkas Konfigurasi di Root
* **`AGENTS.md`**: Pedoman utama aturan pengkodean, batasan layer arsitektur, standar WebRTC, dan checklist verifikasi kode.
* **`.air.toml`**: Konfigurasi live-reload server Go sehingga saat ada file kode `.go` yang diubah, server otomatis mengompilasi ulang dan me-restart secara instan.
* **`.env` / `.env.example`**: Konfigurasi port aplikasi, kredensial database MySQL, dan rahasia TURN server.
* **`go.mod` / `go.sum`**: Pengunci versi package eksternal Go.

---

## 4. Struktur 4 Role Sistem

Aplikasi menerapkan pembagian hak akses (*Role-Based Access Control / RBAC*) dengan 4 tingkatan:

1. **`superadmin`** (ID: `11111111-1111-1111-1111-111111111111`)
   - Memiliki kendali penuh terhadap seluruh sistem.
   - Dapat menambah, mengubah, dan menghapus role kustom serta mengelola seluruh pengguna.
2. **`admin`** (ID: `22222222-2222-2222-2222-222222222222`)
   - Mengelola operasional harian, pendaftaran pengguna operator/warga, serta memantau ruangan panggilan.
3. **`operator`** (ID: `33333333-3333-3333-3333-333333333333`)
   - Petugas pelayanan dinas yang bertugas menerima dan memandu sesi konsultasi tatap muka daring dengan warga.
4. **`user`** (ID: `44444444-4444-4444-4444-444444444444`)
   - Warga masyarakat yang menggunakan aplikasi untuk mengikuti sesi panggilan video konsultasi publik.

---

## 5. Panduan Menjalankan Aplikasi

### Prasyarat
1. **Go** versi 1.22 atau lebih baru.
2. **MySQL Server** versi 8.x yang sedang aktif.
3. **Air** (opsional untuk live reload): `go install github.com/air-verse/air@latest`.

### Langkah Menjalankan:
1. Pastikan file `.env` sudah sesuai dengan kredensial database lokal Anda:
   ```env
   APP_PORT=8080
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=123
   DB_NAME=webrtc_db
   ```
2. Jalankan server dengan fitur Live Reload (Air):
   ```bash
   air
   ```
   Atau jalankan langsung dengan perintah bawaan Go:
   ```bash
   go run cmd/api/main.go
   ```
3. Saat pertama kali dijalankan, sistem akan **otomatis membuat tabel basis data, mengisi 4 role default, dan membuat akun superadmin**:
   - **Email Login**: `superadmin@pekanbaru.go.id`
   - **Password**: `Password123!`

---

## 6. Dokumentasi Endpoint API

Semua endpoint mengembalikan respons dalam format amplop JSON standar:

### A. Health Check
* **`GET /health`**
  - Mengembalikan status kesehatan server backend.

### B. Manajemen Role (`/api/v1/roles`)
| Method | Endpoint | Deskripsi | Status Sukses |
| :--- | :--- | :--- | :---: |
| `GET` | `/api/v1/roles` | Mengambil seluruh daftar role | `200 OK` |
| `GET` | `/api/v1/roles/:id` | Mengambil detail satu role | `200 OK` |
| `POST` | `/api/v1/roles` | Menambahkan role kustom baru | `201 Created` |
| `PUT` | `/api/v1/roles/:id` | Memperbarui nama/deskripsi role | `200 OK` |
| `DELETE` | `/api/v1/roles/:id` | Menghapus role (role bawaan diproteksi) | `200 OK` |

### C. Manajemen Pengguna (`/api/v1/users`)
| Method | Endpoint | Query Param | Deskripsi | Status Sukses |
| :--- | :--- | :--- | :--- | :---: |
| `GET` | `/api/v1/users` | `page`, `limit`, `search`, `role_id` | Daftar user berhalaman | `200 OK` |
| `GET` | `/api/v1/users/:id` | - | Detail profil user | `200 OK` |
| `POST` | `/api/v1/users` | - | Membuat akun pengguna baru | `201 Created` |
| `PUT` | `/api/v1/users/:id` | - | Memperbarui akun pengguna | `200 OK` |
| `DELETE` | `/api/v1/users/:id` | - | Menghapus akun pengguna | `200 OK` |

---

## 7. Pengujian & Standar Kualitas Kode

Jalankan perintah berikut sebelum menyerahkan perubahan kode:
```bash
# 1. Pastikan seluruh paket dapat dikompilasi
go build ./...

# 2. Pengujian statis analisa kode
go vet ./...

# 3. Eksekusi unit test dengan pendeteksi race condition
go test -race ./...
```
