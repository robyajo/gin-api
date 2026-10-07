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
   - Bertanggung jawab mengembalikan status HTTP dan struktur DTO response.

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
2. Gunakan tipe data `CHAR(36)` atau `BINARY(16)` untuk Primary Key bertipe UUID.
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

## 5. Checklist Verifikasi Sebelum Menyerahkan Kode

- [ ] Kode dapat dikompilasi tanpa error (`go build ./...`).
- [ ] Lolos pengujian linter (`golangci-lint run`).
- [ ] Tidak ada dependensi cyclic (_cycle import_).
- [ ] Parameter sensitif (secret key JWT, password database, TURN shared secret) diambil dari file environment, bukan di-_hardcode_.
- [ ] Koneksi WebSocket dan channel ditutup dengan aman (_defer close_) untuk mencegah kebocoran memori.
