package database

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gin-api/pkg/utils"
)

// RunMigrations menjalankan pembuatan skema dasar dan seeding 4 role default
func RunMigrations(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS roles (
		id VARCHAR(36) NOT NULL,
		name VARCHAR(50) NOT NULL,
		description VARCHAR(255) NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY uq_roles_name (name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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

	CREATE TABLE IF NOT EXISTS call_logs (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		room_id VARCHAR(36) NOT NULL,
		caller_id VARCHAR(36) NOT NULL,
		started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		ended_at TIMESTAMP NULL DEFAULT NULL,
		duration_seconds INT UNSIGNED NOT NULL DEFAULT 0,
		status ENUM('completed', 'missed', 'rejected', 'failed') NOT NULL DEFAULT 'completed',
		metadata JSON NULL,
		PRIMARY KEY (id),
		KEY idx_call_logs_room_id (room_id),
		KEY idx_call_logs_caller_id (caller_id),
		CONSTRAINT fk_call_logs_room_id FOREIGN KEY (room_id) REFERENCES rooms (id) ON DELETE CASCADE,
		CONSTRAINT fk_call_logs_caller_id FOREIGN KEY (caller_id) REFERENCES users (id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	statements := []string{
		`CREATE TABLE IF NOT EXISTS roles (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(50) NOT NULL,
			description VARCHAR(255) NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uq_roles_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS users (
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
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS rooms (
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
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS room_participants (
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
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS call_logs (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			room_id VARCHAR(36) NOT NULL,
			caller_id VARCHAR(36) NOT NULL,
			started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			ended_at TIMESTAMP NULL DEFAULT NULL,
			duration_seconds INT UNSIGNED NOT NULL DEFAULT 0,
			status ENUM('completed', 'missed', 'rejected', 'failed') NOT NULL DEFAULT 'completed',
			metadata JSON NULL,
			PRIMARY KEY (id),
			KEY idx_call_logs_room_id (room_id),
			KEY idx_call_logs_caller_id (caller_id),
			CONSTRAINT fk_call_logs_room_id FOREIGN KEY (room_id) REFERENCES rooms (id) ON DELETE CASCADE,
			CONSTRAINT fk_call_logs_caller_id FOREIGN KEY (caller_id) REFERENCES users (id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	_ = schema

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migration exec failed: %w", err)
		}
	}

	// Seed 4 Role Default: superadmin, admin, operator, user
	defaultRoles := []struct {
		ID          string
		Name        string
		Description string
	}{
		{"11111111-1111-1111-1111-111111111111", "superadmin", "Akses penuh seluruh konfigurasi, hak akses, dan manajemen sistem"},
		{"22222222-2222-2222-2222-222222222222", "admin", "Pengelolaan pengguna, ruang konsultasi, dan monitoring operasional"},
		{"33333333-3333-3333-3333-333333333333", "operator", "Petugas operator penanggung jawab ruang sesi panggilan tatap muka warga"},
		{"44444444-4444-4444-4444-444444444444", "user", "Pengguna umum warga masyarakat untuk akses layanan konsultasi daring"},
	}

	for _, r := range defaultRoles {
		query := `INSERT INTO roles (id, name, description) VALUES (?, ?, ?)
				  ON DUPLICATE KEY UPDATE description = VALUES(description)`
		if _, err := db.Exec(query, r.ID, r.Name, r.Description); err != nil {
			log.Printf("Warning seeding role %s: %v", r.Name, err)
		}
	}

	// Seed default superadmin user jika belum ada user sama sekali
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = ?", "superadmin@pekanbaru.go.id")
	if err == nil && count == 0 {
		hash, err := utils.HashPassword("Password123!")
		if err == nil {
			adminID := uuid.New().String()
			insertUser := `INSERT INTO users (id, role_id, name, email, password_hash, is_active)
						   VALUES (?, ?, ?, ?, ?, 1)`
			_, _ = db.Exec(insertUser, adminID, "11111111-1111-1111-1111-111111111111", "Super Administrator", "superadmin@pekanbaru.go.id", hash)
			log.Println("==> Akun Superadmin berhasil diinisialisasi: superadmin@pekanbaru.go.id (Password123!)")
		}
	}

	return nil
}
