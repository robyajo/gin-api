---
name: api-testing
description: "Gunakan skill ini saat menguji endpoint REST API Gin, alur otentikasi JWT, validasi health check, pengujian integrasi WebSocket signaling, atau menjalankan test suite Go (go test)."
metadata:
    author: project
---

# API & WebSocket Integration Testing Guide

Panduan operasional bagi AI Agent untuk melakukan verifikasi fungsionalitas HTTP REST dan WebSocket Signaling pada `gin-api`.

---

## 1. Format Standar Respons REST API

Setiap respons HTTP dari handler wajib mematuhi standar amplop JSON terpadu:

```json
{
  "code": 200,
  "status": "success",
  "message": "Operasi berhasil dilakukan",
  "data": { ... }
}
```

Format respons error:
```json
{
  "code": 400,
  "status": "error",
  "message": "Validasi input gagal",
  "errors": [ ... ]
}
```

---

## 2. Pengujian Alur REST via CLI

1. **Pemeriksaan Health Check:**
   ```bash
   curl -s http://localhost:8080/health
   ```
2. **Pengujian Register & Login:**
   ```bash
   # Login dan ambil token JWT
   TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"password123"}' | jq -r .data.token)
   ```
3. **Panggilan Endpoint Terproteksi:**
   ```bash
   curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/rooms
   ```

---

## 3. Pengujian WebSocket Signaling

Gunakan alat bantu CLI seperti `websocat` atau script test Go:
```bash
# Handshake dengan token JWT
websocat "ws://localhost:8080/ws?token=$TOKEN&room_id=test-room"
```

Kirim pesan join test:
```json
{"type":"join","room_id":"test-room","sender_id":"user-1","payload":{}}
```

---

## 4. Eksekusi Unit & Integration Tests

Jalankan test bawaan Go:
```bash
# Menjalankan seluruh test dengan ringkasan verbositas
go test -v ./...

# Menjalankan test dengan deteksi race condition
go test -race ./...
```
