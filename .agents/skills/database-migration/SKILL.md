---
name: database-migration
description: "Gunakan skill ini saat merancang skema database MySQL, membuat atau mengubah file migrasi up/down di scripts/migrations/, memastikan konvensi UUID CHAR(36), indeks foreign key, transaksi basis data, dan integritas data relasional."
metadata:
    author: project
---

# Database Migration & MySQL Standards Guide

Panduan operasional bagi AI Agent untuk mengelola evolusi skema basis data MySQL pada aplikasi `gin-api`.

---

## 1. Konvensi Penamaan Berkas Migrasi

Semua script migrasi disimpan di direktori `scripts/migrations/` dengan pola:

```text
scripts/migrations/
├── 000001_init_schema.up.sql
├── 000001_init_schema.down.sql
├── 000002_add_room_participants.up.sql
└── 000002_add_room_participants.down.sql
```

- Penomoran wajib berurutan 6 digit (`000001_...`).
- Setiap berkas `.up.sql` (eksekusi perubahan) WAJIB memiliki pasangan `.down.sql` (rollback pemulihan).

---

## 2. Standar DDL MySQL

1. **Engine & Charset:** Selalu gunakan `ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`.
2. **Primary Key UUID:** Gunakan `CHAR(36)` untuk ID berformat UUID v4 string.
3. **Audit Timestamps:** Sertakan kolom:
   ```sql
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   ```
4. **Relasi Foreign Key:**
   - Berikan nama constraint yang jelas: `fk_<table_name>_<reference_table>`.
   - Tentukan aksi `ON DELETE` secara eksplisit (`CASCADE` atau `RESTRICT`/`SET NULL`).
5. **Indeks Performa:**
   - Tambahkan indeks pada kolom yang sering dicari (`room_id`, `user_id`, `status`, `created_at`).

---

## 3. Aturan Keamanan & Transaksi Data

- Operasi DDL yang memengaruhi data produksi tidak boleh menghapus data tanpa mekanisme backup.
- Operasi multi-tabel dalam kode Go WAJIB menggunakan `sqlx.Tx` (Database Transaction):
  ```go
  tx, err := db.Beginx()
  if err != nil { return err }
  defer tx.Rollback()
  // Operasi SQL...
  return tx.Commit()
  ```

