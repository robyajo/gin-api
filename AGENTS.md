# Developer & Agent Guidelines (AGENTS.md)

Instruksi ini ditujukan untuk setiap AI Agent, LLM, atau pengembang yang bertugas membaca, mengedit, memperluas, dan memelihara kode basis (_codebase_) sistem ini.

---

## 1. Filosofi & Pola Desain (Clean Layered Architecture)

Arsitektur aplikasi mematuhi hierarki pemisahan dependensi searah:

```text
[ Delivery / HTTP & WS ] ---> [ Service / Business Logic ] ---> [ Repository / Data Access ]
```

### Aturan Batasan Layer:

1. **Layer Delivery (`internal/delivery/`)**:
   - Hanya menangani parsing payload HTTP (`*gin.Context`) atau pesan frame WebSocket.
   - Tidak boleh memuat query database (SQL) atau logika manipulasi data bisnis.
   - Bertanggung jawab mengembalikan status HTTP dan struktur DTO response standar.

2. **Layer Service (`internal/service/`)**:
   - Berisi business logic murni.
   - **DILARANG** mengimpor atau menggunakan `github.com/gin-gonic/gin` atau `*gin.Context` di layer ini. Gunakan Go standard `context.Context`.
   - Mengorkestrasi pemanggilan antar repository atau external service (misal: TURN credential generator).

3. **Layer Repository (`internal/repository/`)**:
   - Hanya bertanggung jawab melakukan eksekusi SQL ke MySQL menggunakan connection pool `sqlx.DB` atau `database/sql`.
   - Selalu sertakan `context.Context` pada query database (`QueryContext`, `ExecContext`) untuk mendukung pembatalan timeout dan graceful shutdown.

---

## 2. Standar Penanganan WebSocket & WebRTC Signaling

1. **Model Hub & Client:**
   - Gunakan pola _Gorilla Hub_ yang memisahkan channel: `register`, `unregister`, dan `broadcast`.
   - Pisahkan goroutine baca dan tulis untuk setiap client (`readPump` dan `writePump`) guna mencegah race condition pada socket.
2. **Thread-Safety:**
   - Hub Room Manager harus dilindungi menggunakan `sync.RWMutex` saat membaca atau memodifikasi daftar client/peer yang aktif.
3. **Standar Pesan (Envelope Contract):**
   - Jangan mengirim raw string. Seluruh komunikasi WebSocket wajib divalidasi ke dalam struct pesan standar:
     ```go
     type SignalMessage struct {
         Type     string          `json:"type"`
         RoomID   string          `json:"room_id"`
         SenderID string          `json:"sender_id"`
         TargetID string          `json:"target_id,omitempty"`
         Payload  json.RawMessage `json:"payload"`
     }
     ```
4. **Validasi Handshake:**
   - Token JWT **wajib** divalidasi sebelum koneksi di-upgrade dari HTTP ke WebSocket (`upgrader.Upgrade`). Jika tidak valid, kembalikan HTTP 401 Unauthorized secara langsung.

---

## 3. Konvensi Database & SQL (MySQL)

1. Gunakan nama tabel jamak dalam format `snake_case` (contoh: `users`, `rooms`, `call_sessions`).
2. Gunakan tipe data `CHAR(36)` untuk Primary Key bertipe UUID v4.
3. Pastikan relasi antar tabel dilindungi `FOREIGN KEY` dengan batasan `ON DELETE CASCADE` atau `ON DELETE SET NULL` yang tepat.
4. Gunakan `sql.NullString`, `sql.NullTime`, atau pointer struct untuk kolom yang memperbolehkan nilai `NULL`.
5. Semua operasi database yang mengubah status kritis (misal: penutupan room dan pencatatan akhir durasi call session) wajib dibungkus dalam **Database Transaction (`sqlx.Tx`)**.

---

## 4. Standar Kode & Konvensi Go

1. **Error Handling:**
   - Jangan pernah mengabaikan return `error` (hindari penulisan `_ = doSomething()`).
   - Bungkus error dengan konteks menggunakan `fmt.Errorf("failed to ...: %w", err)` agar _stack trace_ mudah dilacak.
2. **Konkurensi:**
   - Pastikan setiap goroutine yang dibuat memiliki siklus hidup yang jelas untuk menghindari _goroutine leaks_.
   - Gunakan `select` dengan `ctx.Done()` atau sinyal quit channel pada loop berulang.
3. **Nama Variabel & Interface:**
   - Deklarasikan interface pada sisi konsumen (_consumer_), bukan pada sisi penyedia (_provider_).
   - Nama method dan struct harus ringkas, jelas, dan idiomatis Go (_idiomatic Go_).

---

## 5. ⚡ Protokol Catatan Rilis: 1 Berkas JSON Per Hari

Setiap penambahan fitur atau perbaikan bug yang **dirasakan pengguna (warga, operator, admin)** harus dicatat di direktori `releases/` dengan aturan berikut:

### Aturan Konsolidasi Harian (1 JSON Per Hari)

1. **Periksa Tanggal Hari Ini**: Cek apakah berkas `releases/YYYY-MM-DD-*.json` sudah ada.
2. **Jika SUDAH ADA**:
   - **JANGAN** membuat berkas baru!
   - Buka berkas tersebut dan tambahkan poin baru ke dalam array `items: [...]`.
   - Perbarui ringkasan `title` dan `body` harian secara terpadu.
   - Naikkan `type` jika terdapat kategori yang lebih tinggi (`important` > `feature` > `fix` > `improvement`).
3. **Jika BELUM ADA**:
   - Buat berkas baru dengan nama `releases/YYYY-MM-DD-nama-singkat.json`.
4. **Format Dwibahasa (`id` & `en`)**:
   - Wajib menyertakan `id` (Bahasa Indonesia) dan `en` (Bahasa Inggris) untuk `title`, `body`, dan tiap elemen dalam `items`.
5. **Wajib Bahasa Publik & Bebas Istilah Teknis Internal (Security-Safe & Citizen-Facing)**:
   - **DILARANG KERAS** menuliskan istilah teknis arsitektur inti, nama library/package pihak ketiga, nama framework, database, middleware, nama file kode sumber, atau mekanisme internal (seperti _Gin, Gorilla, WebSocket, WebRTC, Coturn, STUN, TURN, MySQL, sqlx, JWT, bcrypt, air, goroutine, mutex, endpoint `/api/v1/...`_, dsb.) di dalam berkas rilis.
   - **Tujuan Keamanan**: Mencegah pihak luar menganalisis atau menebak struktur komponen internal dan potensi celah keamanan sistem dari catatan rilis publik.
   - **Orientasi Pelayanan Warga**: Seluruh `title`, `body`, dan butir `items` wajib dirumuskan sebagai **pemberitahuan pelayanan publik** yang ramah, sopan, dan berfokus pada dampak langsung bagi masyarakat (misalnya: kemudahan konsultasi daring, kejernihan audio/video, kecepatan membuka ruang pertemuan, penghematan kuota ponsel, dan perlindungan kerahasiaan data warga).

### Contoh Format Catatan Rilis:
```json
{
  "date": "2026-10-08",
  "type": "feature",
  "title": {
    "id": "Penyediaan Fondasi Layanan Komunikasi Interaktif & Konsultasi Warga",
    "en": "Provision of Interactive Communication & Citizen Consultation Service Foundation"
  },
  "body": {
    "id": "Penyediaan fondasi sistem layanan komunikasi tatap muka daring untuk mendukung pelaksanaan konsultasi dan interaksi warga secara langsung, stabil, dan aman.",
    "en": "Establishment of real-time online consultation service foundations to support direct, stable, and secure citizen interactions."
  },
  "items": [
    {
      "type": "new",
      "id": "Penyediaan jalur layanan pemantauan keaktifan sistem secara berkala untuk memastikan keandalan operasional layanan publik setiap saat.",
      "en": "Provision of periodic system health monitoring channels to ensure public service operational reliability at all times."
    }
  ]
}
```

---

## 6. 🛠️ Kumpulan Skills Agen (`.agents/skills/`)

Untuk menyelesaikan tugas secara presisi, gunakan referensi skill yang tersedia di direktori `.agents/skills/`:

1. **`release-notes`** ([.agents/skills/release-notes/SKILL.md](file:///.agents/skills/release-notes/SKILL.md)):
   - Panduan pencatatan berkas rilis harian, validasi konsolidasi 1 JSON/hari, skema dwibahasa, dan standardisasi bahasa publik yang aman.
2. **`webrtc-signaling`** ([.agents/skills/webrtc-signaling/SKILL.md](file:///.agents/skills/webrtc-signaling/SKILL.md)):
   - Panduan arsitektur WebSocket Hub, negosiasi SDP Offer/Answer, pertukaran ICE Candidates, proteksi mutex, dan pembuatan kredensial ephemeral TURN.
3. **`database-migration`** ([.agents/skills/database-migration/SKILL.md](file:///.agents/skills/database-migration/SKILL.md)):
   - Panduan pembuatan skema migrasi MySQL di `scripts/migrations/`, standar UUID `CHAR(36)`, relasi foreign key, dan transaksi database.
4. **`api-testing`** ([.agents/skills/api-testing/SKILL.md](file:///.agents/skills/api-testing/SKILL.md)):
   - Panduan pengujian HTTP handler Gin, otentikasi JWT, handshake WebSocket via CLI, dan unit testing Go.

---

## 7. Checklist Verifikasi Sebelum Menyerahkan Kode

- [ ] Kode dapat dikompilasi tanpa error (`go build ./...`).
- [ ] Lolos pengujian linter (`golangci-lint run` atau `go vet ./...`).
- [ ] Tidak ada dependensi cyclic (_cycle import_).
- [ ] Parameter sensitif (secret key JWT, password database, TURN shared secret) diambil dari file environment, bukan di-_hardcode_.
- [ ] Koneksi WebSocket dan channel ditutup dengan aman (_defer close_) untuk mencegah kebocoran memori.
- [ ] **Catatan Rilis Publik (Wajib)**: Jika perubahan berdampak pada fungsionalitas pengguna, pastikan berkas rilis harian di `releases/YYYY-MM-DD-*.json` sudah dibuat atau diperbarui sesuai protokol di Bagian 5.
