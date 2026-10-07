# Project Structure: Go-Gin WebRTC Signaling & REST API

Struktur ini mengimplementasikan **Clean / Layered Architecture** yang mematuhi konvensi _Standard Go Project Layout_, memisahkan layer HTTP Delivery, Business Logic (Service), WebSocket Hub, Data Access (Repository), dan Migrasi Skema MySQL.

---

## 1. Diagram Pohon Direktori (Directory Tree)

```text
webrtc-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point server: DI wiring, database init, router boot
├── config/
│   └── config.go                   # Parsing environment variables (.env / os.Getenv)
├── internal/                       # Domain internal private (tidak dapat di-import modul eksternal)
│   ├── delivery/                   # Interface layer (Inbound adapters)
│   │   ├── http/
│   │   │   ├── handler/            # Controller REST API (Gin Context)
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── room_handler.go
│   │   │   │   └── turn_handler.go
│   │   │   ├── middleware/         # Middleware HTTP
│   │   │   │   ├── auth_middleware.go
│   │   │   │   └── cors_middleware.go
│   │   │   └── router.go           # Setup grouping & routing endpoint Gin
│   │   └── ws/                     # Real-time WebSocket Delivery
│   │       ├── client.go           # Goroutine WebSocket peer (readPump / writePump)
│   │       ├── hub.go              # WebSocket Room Manager / In-memory Message Router
│   │       ├── payload.go          # Struct sinyal envelope WebRTC (SDP / ICE candidates)
│   │       └── ws_handler.go       # Upgrade HTTP handshake ke WebSocket & registrasi
│   ├── model/                      # Struct Entity database & DTO (Data Transfer Object)
│   │   ├── call_session.go         # Model sessions & participants
│   │   ├── dto/                    # Request & Response payload definition
│   │   │   ├── auth_dto.go
│   │   │   └── room_dto.go
│   │   ├── room.go                 # Model rooms
│   │   └── user.go                 # Model users
│   ├── repository/                 # Data persistence layer (MySQL queries via sqlx)
│   │   ├── call_session_repository.go
│   │   ├── room_repository.go
│   │   └── user_repository.go
│   └── service/                    # Pure business logic layer (tidak boleh ada Gin context)
│       ├── auth_service.go
│       ├── room_service.go
│       └── turn_service.go         # Ephemeral HMAC-SHA1 credential generator
├── migrations/                     # Skema SQL DDL untuk golang-migrate
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_rooms_table.up.sql
│   ├── 000002_create_rooms_table.down.sql
│   ├── 000003_create_call_sessions_tables.up.sql
│   └── 000003_create_call_sessions_tables.down.sql
├── pkg/                            # Paket utilitas reusable lintas modul
│   ├── database/
│   │   └── mysql.go                # Koneksi connection pool MySQL (sqlx / sql.DB)
│   ├── hash/
│   │   └── password.go             # Utilitas hashing BCrypt
│   ├── jwt/
│   │   └── token.go                # Generator & parser JWT HS256/RS256
│   └── logger/
│       └── logger.go               # Inisialisasi Zap / Zerolog instance
├── .env.example                    # Template variabel lingkungan
├── .gitignore                      # Git ignore rule
├── Dockerfile                      # Multi-stage Docker build untuk Go binary
├── docker-compose.yml              # Stack lokal (MySQL 8.x + Coturn server)
├── go.mod                          # Definisi modul Go
├── go.sum                          # Checksum hash dependensi
├── Makefile                        # Shorthand CLI (make run, make migrate, make test)
└── README.md                       # Dokumentasi setup awal
```

---

## 2. Tanggung Jawab Modul Inti

| Modul                                  | Tanggung Jawab Utama                                                                                                               |
| :------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------- |
| **`cmd/api/main.go`**                  | Menghubungkan semua dependensi (DB -> Repo -> Service -> Handler -> Hub -> Gin Engine).                                            |
| **`internal/delivery/ws/hub.go`**      | Mengelola ruang (_room memory pool_), thread-safe routing pertukaran SDP Offer/Answer dan ICE candidates antar peer.               |
| **`internal/delivery/ws/client.go`**   | Mengisolasi operasi _socket I/O_ ke dalam dua goroutine terpisah: `readPump` (baca dari socket) dan `writePump` (tulis ke socket). |
| **`internal/service/turn_service.go`** | Menerbitkan kredensial STUN/TURN berbatas waktu (_ephemeral REST API standard RFC 7635_) menggunakan algoritma HMAC-SHA1.          |
| **`internal/repository/`**             | Mengeksekusi query SQL mentah / prepared statement terindeks ke database MySQL.                                                    |
| **`pkg/`**                             | Library umum yang tidak terikat pada _business logic_ spesifik domain aplikasi.                                                    |

---

## 3. Skrip Otomatisasi Pembuatan Folder (WSL Terminal)

Jalankan perintah baris tunggal ini di terminal WSL Ubuntu untuk membuat seluruh struktur folder dan file secara instan:

```bash
mkdir -p cmd/api config internal/delivery/http/handler internal/delivery/http/middleware internal/delivery/ws internal/model/dto internal/repository internal/service migrations pkg/database pkg/hash pkg/jwt pkg/logger

# Buat file-file utama
touch cmd/api/main.go config/config.go
touch internal/delivery/http/router.go
touch internal/delivery/http/handler/{auth_handler.go,room_handler.go,turn_handler.go}
touch internal/delivery/http/middleware/{auth_middleware.go,cors_middleware.go}
touch internal/delivery/ws/{client.go,hub.go,payload.go,ws_handler.go}
touch internal/model/{user.go,room.go,call_session.go}
touch internal/model/dto/{auth_dto.go,room_dto.go}
touch internal/repository/{user_repository.go,room_repository.go,call_session_repository.go}
touch internal/service/{auth_service.go,room_service.go,turn_service.go}
touch migrations/{000001_create_users_table.up.sql,000001_create_users_table.down.sql,000002_create_rooms_table.up.sql,000002_create_rooms_table.down.sql,000003_create_call_sessions_tables.up.sql,000003_create_call_sessions_tables.down.sql}
touch pkg/database/mysql.go pkg/hash/password.go pkg/jwt/token.go pkg/logger/logger.go
touch .env.example .gitignore Dockerfile docker-compose.yml Makefile README.md
```
