---
name: release-notes
description: "ACTIVATE setiap kali selesai mengerjakan perubahan yang berdampak langsung atau dirasakan pengguna: menambah fitur baru, memperbaiki bug, mengubah perilaku API/layanan, atau memperbarui keamanan. Triggers on: selesai perubahan, tambah fitur, ubah perilaku, perbaikan bug, releases/, changelog, catatan rilis, release-notes. Menulis atau memperbarui berkas rilis JSON harian di folder releases/. Do NOT activate untuk refactor internal tanpa perubahan perilaku, formatting kode, atau saat hanya membaca rilis."
metadata:
    author: project
---

# Catatan Rilis (Changelog) - Layanan Komunikasi Warga (gin-api)

Setiap perubahan yang **dirasakan pengguna (warga, operator, atau admin)** harus tercatat di dalam folder `releases/`. Berkas JSON ini merupakan catatan resmi riwayat rilis publik aplikasi.

---

## Kapan WAJIB Menulis/Memperbarui Berkas Rilis

| Jenis Perubahan | Wajib Berkas Rilis? | Kategori Tipe (`type`) |
| :--- | :---: | :--- |
| Fitur baru (endpoint publik, alur konsultasi baru, dsb.) | **Ya** | `feature` |
| Perbaikan kendala/bug yang dirasakan warga atau operator | **Ya** | `fix` |
| Peningkatan kecepatan respon, latensi audio/video, atau efisiensi data | **Ya** | `improvement` |
| Penambalan keamanan & perlindungan data warga | **Ya** | `important` / `security` |
| Refactor internal tanpa perubahan perilaku, komentar, format kode | **Tidak** | - |

> **Prinsip Utama**: Tanyakan *"Apakah warga atau operator pengguna layanan akan menyadari, terbantu, atau merasakan dampak perubahan ini?"* Jika ya, catat ke berkas rilis!

---

## ⚡ ATURAN UTAMA: 1 Berkas JSON Per Hari (Konsolidasi Harian)

- **MAKSIMAL 1 BERKAS JSON PER HARI**:
  - Jangan membuat berkas baru jika pada hari yang sama sudah ada berkas rilis yang dibuat.
  - **Langkah Kerja AI**:
    1. Periksa folder `releases/` untuk tanggal hari ini (`releases/YYYY-MM-DD-*.json`).
    2. **Jika SUDAH ADA**:
       - Buka berkas JSON tersebut.
       - Tambahkan poin perubahan baru ke dalam array `items: [...]`.
       - Perbarui `title` dan `body` agar mencerminkan rangkuman harian terpadu.
       - Naikkan nilai `type` jika ada poin dengan prioritas lebih tinggi (`important` > `feature` > `fix` > `improvement`).
    3. **Jika BELUM ADA**:
       - Buat berkas baru dengan nama berpola: `releases/YYYY-MM-DD-nama-singkat.json`.

---

## Format JSON Dwibahasa

```json
{
  "date": "2026-10-08",
  "type": "feature",
  "title": {
    "id": "Judul perubahan dalam Bahasa Indonesia",
    "en": "Title of the change in English"
  },
  "body": {
    "id": "Satu atau dua kalimat narasi ringkasan perubahan.",
    "en": "One or two sentences summarizing the changes."
  },
  "items": [
    {
      "type": "new",
      "id": "Deskripsi peningkatan atau fitur yang dirasakan pengguna.",
      "en": "Description of the enhancement or feature experienced by users."
    },
    {
      "type": "fix",
      "id": "Deskripsi kendala yang diselesaikan dan kenyamanan yang didapat.",
      "en": "Description of the resolved issue and its benefit to users."
    }
  ]
}
```

### Pedoman Kolom

| Kolom | Aturan Format |
| :--- | :--- |
| `date` | Format `YYYY-MM-DD`, tanggal hari kerja saat perubahan dirilis. |
| `type` | Prioritas: `important` > `feature` > `fix` > `improvement`. |
| `title` | Wajib dwibahasa (`id` & `en`), ringkas dan jelas (maks. 80 karakter). |
| `body` | Narasi 1–2 kalimat dwibahasa: ringkasan dampak positif bagi kenyamanan warga. |
| `items` | Array butir perubahan dwibahasa dengan subtipe `type` (`new`, `fix`, `improvement`, `security`, `removed`). |

---

## Bahasa Publik & Bebas Istilah Teknis Internal (Security-Safe & Citizen-Facing)

Catatan rilis ini dipublikasikan untuk warga masyarakat. **DILARANG KERAS** memuat istilah teknis arsitektur inti, nama library, package pihak ketiga, framework, database, middleware, nama berkas, atau mekanisme internal (seperti *Gin, Gorilla, WebSocket, WebRTC, Coturn, STUN, TURN, MySQL, sqlx, JWT, bcrypt, air, goroutine, mutex, endpoint `/api/v1/...`* dsb.):

| ❌ Dilarang (Terlalu Teknis / Celah Keamanan) | ✅ Wajib (Bahasa Pelayanan Warga & Aman) |
| :--- | :--- |
| "Implementasi endpoint WebSocket di `/ws/signaling` dengan Gorilla Hub" | "Penyediaan Jalur Komunikasi Tatap Muka: Warga kini dapat terhubung ke ruang konsultasi video secara langsung dan interaktif." |
| "Integrasi Coturn server dengan kredensial ephemeral HMAC-SHA1" | "Peningkatan Kelancaran Sambungan: Panggilan video tetap tersambung dengan stabil meski warga menggunakan jaringan internet terbatas." |
| "Perbaikan race condition mutex pada map client room" | "Peningkatan Kestabilan Ruang Temu: Sesi konsultasi tatap muka bebas dari risiko terputus saat banyak peserta bergabung." |
| "Validasi token JWT pada HTTP upgrade handshake" | "Perlindungan Keamanan Sesi Konsultasi: Penguatan verifikasi identitas untuk memastikan hanya pihak berwenang yang dapat masuk ke ruang temu." |
| "Optimasi connection pool MySQL sqlx 50 open conns" | "Peningkatan Kecepatan Akses: Waktu tunggu pembukaan ruang layanan daring menjadi jauh lebih cepat dan responsif." |

