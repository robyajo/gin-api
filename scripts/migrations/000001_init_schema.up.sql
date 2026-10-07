-- ==========================================================
-- 1. TABEL: roles
-- ==========================================================
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_roles_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed default 4 roles jika belum ada
INSERT INTO roles (id, name, description) VALUES
    ('11111111-1111-1111-1111-111111111111', 'superadmin', 'Akses penuh seluruh konfigurasi, hak akses, dan manajemen sistem'),
    ('22222222-2222-2222-2222-222222222222', 'admin', 'Pengelolaan pengguna, ruang konsultasi, dan monitoring operasional'),
    ('33333333-3333-3333-3333-333333333333', 'operator', 'Petugas operator penanggung jawab ruang sesi panggilan tatap muka warga'),
    ('44444444-4444-4444-4444-444444444444', 'user', 'Pengguna umum warga masyarakat untuk akses layanan konsultasi daring')
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;

-- ==========================================================
-- 2. TABEL: users
-- ==========================================================
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email),
    KEY idx_users_role_id (role_id),
    CONSTRAINT fk_users_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==========================================================
-- 3. TABEL: rooms
-- ==========================================================
CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(36) NOT NULL,
    title VARCHAR(100) NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    status ENUM('active', 'ended', 'locked') NOT NULL DEFAULT 'active',
    max_participants INT UNSIGNED NOT NULL DEFAULT 2,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_rooms_created_by (created_by),
    KEY idx_rooms_status (status),
    CONSTRAINT fk_rooms_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==========================================================
-- 4. TABEL: room_participants
-- ==========================================================
CREATE TABLE IF NOT EXISTS room_participants (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    room_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP NULL DEFAULT NULL,
    connection_status ENUM('connected', 'disconnected', 'failed') NOT NULL DEFAULT 'connected',
    PRIMARY KEY (id),
    KEY idx_participants_room_id (room_id),
    KEY idx_participants_user_id (user_id),
    CONSTRAINT fk_participants_room_id FOREIGN KEY (room_id) REFERENCES rooms (id) ON DELETE CASCADE,
    CONSTRAINT fk_participants_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==========================================================
-- 5. TABEL: call_logs
-- ==========================================================
CREATE TABLE IF NOT EXISTS call_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    room_id VARCHAR(36) NOT NULL,
    caller_id VARCHAR(36) NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NULL DEFAULT NULL,
    duration_seconds INT UNSIGNED NOT NULL DEFAULT 0,
    status ENUM('completed', 'missed', 'rejected', 'failed') NOT NULL DEFAULT 'completed',
    metadata JSON NULL COMMENT 'Menyimpan statistik WebRTC, jitter, codec, resolusi',
    PRIMARY KEY (id),
    KEY idx_call_logs_room_id (room_id),
    KEY idx_call_logs_caller_id (caller_id),
    CONSTRAINT fk_call_logs_room_id FOREIGN KEY (room_id) REFERENCES rooms (id) ON DELETE CASCADE,
    CONSTRAINT fk_call_logs_caller_id FOREIGN KEY (caller_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
