# Product Requirement Document (PRD)

## 1. Overview & Ringkasan Eksekutif

Sistem ini merupakan backend service berbasis **Go (Golang)** menggunakan framework **Gin** yang dirancang khusus untuk menangani:

- Manajemen akun pengguna dan otentikasi RESTful API.
- Manajemen ruang panggilan/pertemuan (_Rooms/Sessions_).
- Layanan **Signaling Server via WebSocket** untuk pertukaran SDP Offer/Answer dan ICE Candidates WebRTC.
- Penyediaan kredensial ephemeral STUN/TURN server untuk memfasilitasi koneksi peer-to-peer (P2P).
- Penyimpanan persisten riwayat sesi, audit, dan data relasional menggunakan **MySQL**.

---

## 2. Sasaran & Tujuan Proyek (Goals & Objectives)

- **Latensi Rendah:** Memproses pesan signaling real-time dengan latensi sub-50ms menggunakan Go Goroutines dan WebSocket channel.
- **Konektivitas Tinggi:** Memastikan keberhasilan negosiasi WebRTC pada perangkat di balik NAT/Firewall ketat melalui integrasi TURN REST API.
- **Skalabilitas & Modularitas:** Menerapkan Clean/Layered Architecture agar modul WebSocket, REST handler, dan database repository terisolasi dengan rapi.
- **Integritas Data:** Menjaga data pengguna, riwayat sesi panggilan (_call logs_), dan status ruang secara konsisten di MySQL.

---

## 3. Aktor & Peran Sistem (User Personas)

1. **User / Peer:**
   - Melakukan pendaftaran/login via REST API.
   - Membuat ruang (_create room_) atau bergabung ke ruang yang tersedia (_join room_).
   - Membuka koneksi WebSocket untuk pertukaran pesan WebRTC (SDP/ICE).
   - Melakukan audio/video call secara P2P.
2. **System / Backend Service:**
   - Memvalidasi token JWT saat handshake HTTP dan upgrade WebSocket.
   - Mengelola koneksi WebSocket dalam memori (_Hub/Room Manager_).
   - Menerbitkan kredensial STUN/TURN sementara.
   - Mencatat log durasi dan status panggilan ke MySQL.

---

## 4. Arsitektur & Alur Kerja Teknis

### 4.1. Alur WebRTC Signaling (WebSocket)

```text
[ Client A ]               [ Go Backend / Hub ]              [ Client B ]
     |                              |                              |
     |--- 1. WS Connect (Token) --->|                              |
     |                              |<--- 1. WS Connect (Token) ---|
     |--- 2. Join Room (RoomID) --->|                              |
     |                              |<--- 2. Join Room (RoomID) ---|
     |                              |                              |
     |--- 3. Send SDP Offer ------->|                              |
     |                              |--- 4. Forward SDP Offer ---->|
     |                              |                              |
     |                              |<-- 5. Send SDP Answer -------|
     |<-- 6. Forward SDP Answer ----|                              |
     |                              |                              |
     |<========== 7. ICE Candidate Exchange (via WS) =============>|
     |                                                             |
     |<================ 8. Direct WebRTC P2P Media ===============>|
```

### 4.2. Format Pesan WebSocket (JSON Payload Contract)

Setiap pesan yang dikirim dan diterima via WebSocket wajib mengikuti format amplop (_envelope_) standar:

```json
{
  "type": "offer | answer | ice-candidate | join-room | leave-room | user-joined | user-left | error",
  "room_id": "uuid-v4",
  "sender_id": "uuid-v4",
  "target_id": "uuid-v4 (opsional untuk unicast signaling)",
  "payload": {}
}
```

---

## 5. Fitur Fungsional (Functional Requirements)

### 5.1. Authentication & User Management (REST)

- `POST /api/v1/auth/register`: Mendaftarkan user baru (email, password hash Argon2/Bcrypt, full name).
- `POST /api/v1/auth/login`: Verifikasi kredensial dan menerbitkan JWT access token.
- `GET /api/v1/users/me`: Mengambil profil pengguna terautentikasi.

### 5.2. WebRTC & TURN Service (REST)

- `GET /api/v1/webrtc/ice-servers`: Menghasilkan ephemeral TURN credentials (HMAC-SHA1 time-limited) dan daftar STUN server.

### 5.3. Room Management (REST)

- `POST /api/v1/rooms`: Membuat room baru (opsional: password/kunci akses, kapasitas maksimal).
- `GET /api/v1/rooms/:id`: Memeriksa ketersediaan dan status room.
- `DELETE /api/v1/rooms/:id`: Menutup/menonaktifkan room (hanya oleh pemilik).

### 5.4. Signaling Server (WebSocket)

- `GET /ws`: Endpoint upgrade protokol HTTP ke WebSocket.
  - Memeriksa validitas JWT query token atau header.
  - Mendaftarkan client ke memori `Hub`.
- Event handling:
  - `join-room`: Menghubungkan client ke room channel tertentu di memori.
  - `offer` / `answer`: Mem-forward deskripsi sesi SDP ke peer tujuan di dalam room yang sama.
  - `ice-candidate`: Mem-forward calon kandidat jaringan ICE.
  - `disconnect`: Menghapus client dari hub, mencatat waktu keluar, dan mengirim broadcast `user-left` ke peer yang tersisa.

---

## 6. Persyaratan Non-Fungsional (Non-Functional Requirements)

1. **Performa:** Mampu menangani minimal 2.000 koneksi WebSocket konkuren pada satu instance server Go standar tanpa lonjakan CPU berarti.
2. **Keamanan:**
   - WSS (WebSocket Secure) wajib diaktifkan pada layer produksi (via Reverse Proxy Caddy/Nginx).
   - Kata sandi pengguna disimpan menggunakan hash BCrypt/Argon2id.
   - Cegah koneksi liar WebSocket dengan validasi JWT sebelum handshake upgrade disetujui.
3. **Penyimpanan:** Koneksi MySQL dioptimasi menggunakan pool (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`).
