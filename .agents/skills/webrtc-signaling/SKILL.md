---
name: webrtc-signaling
description: "Gunakan skill ini saat mengimplementasikan, menguji, atau mendebug WebSocket Hub, WebRTC signaling protocol (SDP Offer/Answer, ICE Candidate exchange), manajemen konkurensi room/peers, validasi handshake token, serta generasi kredensial ephemeral TURN server."
metadata:
    author: project
---

# WebRTC Signaling & WebSocket Hub Guide

Panduan operasional bagi AI Agent untuk menangani komponen real-time WebRTC Signaling pada aplikasi `gin-api`.

---

## 1. Arsitektur Hub & Client

Sistem menggunakan model *Hub-Client* berbasis Gorilla WebSocket:

```text
[ Client Peer A ] <--> [ Client.readPump / writePump ] <--> [ Room Hub ] <--> [ Client.readPump / writePump ] <--> [ Client Peer B ]
```

### Aturan Pokok:
1. **Goroutine Terpisah:** Setiap koneksi aktif memiliki tepat 2 goroutine:
   - `readPump()`: Membaca frame dari socket, memvalidasi ukuran payload, dan mengirimkan pesan ke channel Hub.
   - `writePump()`: Menangani pengiriman pesan dari channel buffer client ke socket dengan ticker `pingPeriod`.
2. **Koneksi Bersih (Graceful Teardown):**
   - Saat socket disconnect, unregister client dari hub, tutup channel send, dan siarkan pesan `peer-left` ke peer lain di room yang sama.
3. **Thread Safety:**
   - Gunakan `sync.RWMutex` pada struct Room / Hub saat membaca atau memperbarui peta koneksi client (`clients map[string]*Client`).

---

## 2. Kontrak Amplop Pesan (Message Envelope Contract)

Semua frame pesan WebSocket WAJIB berbentuk JSON dengan skema terstruktur:

```go
type SignalMessage struct {
    Type     string          `json:"type"`                // join, offer, answer, ice-candidate, leave, error
    RoomID   string          `json:"room_id"`             // ID room target
    SenderID string          `json:"sender_id"`           // UUID pengirim
    TargetID string          `json:"target_id,omitempty"` // UUID penerima khusus (jika direct/unicast)
    Payload  json.RawMessage `json:"payload"`             // Isi payload SDP atau ICE
}
```

### Jenis-Jenis Tipe Pesan (`Type`):
- `join`: Permintaan peer untuk masuk ke room.
- `offer`: WebRTC SDP Offer dari inisiator panggilan.
- `answer`: WebRTC SDP Answer dari penerima panggilan.
- `ice-candidate`: Pertukaran kandidat rute jaringan ICE.
- `leave`: Pemberitahuan keluar dari room.
- `error`: Notifikasi error dari server ke client.

---

## 3. Ephemeral TURN Credential Generator

Untuk menembus Symmetric NAT / Firewall ketat, backend menyediakan kredensial sementara coturn menggunakan mekanisme `TURN REST API (shared secret HMAC-SHA1)`:

```go
// Format Username: <unix_timestamp_expiry>:<user_id>
// Format Password: Base64(HMAC-SHA1(secret, username))
```

Kredensial memiliki masa berlaku terbatas (misal: 24 jam) agar tidak dapat disalahgunakan pihak luar.

---

## 4. Validasi Handshake & Keamanan

1. Token JWT wajib diekstrak dari query parameter (misal: `?token=...`) atau header `Sec-WebSocket-Protocol`.
2. Verifikasi klaim JWT sebelum memanggil `upgrader.Upgrade()`.
3. Batasi ukuran frame maksimum (`SetReadLimit(maxMessageSize)`).
4. Konfigurasikan `CheckOrigin` dengan daftar domain terpercaya, bukan `return true` terbuka tanpa filter di tahap produksi.

