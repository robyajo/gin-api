#!/bin/bash

# Target folder
PROJECT_DIR="gin-api"

echo "==> Membuat folder project: $PROJECT_DIR"
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR" || exit 1

echo "==> Membuat struktur direktori..."
mkdir -p cmd/api \
  config \
  internal/delivery/http/handler \
  internal/delivery/http/middleware \
  internal/delivery/ws \
  internal/model \
  internal/repository \
  internal/service \
  pkg/database \
  pkg/utils \
  scripts/migrations

echo "==> Membuat file-file arsitektur awal..."

# Entry Point (main.go)
cat << 'FILE_EOF' > cmd/api/main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service berjalan normal",
		})
	})

	fmt.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}
FILE_EOF

# Config (config.go)
cat << 'FILE_EOF' > config/config.go
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "webrtc_db"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
FILE_EOF

# Database Connector (mysql.go)
cat << 'FILE_EOF' > pkg/database/mysql.go
package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gin-api/config"
)

func InitMySQL(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
FILE_EOF

# WebSocket Messaging & Hub Placeholders
touch internal/delivery/ws/client.go
touch internal/delivery/ws/hub.go
touch internal/delivery/ws/message.go

# Router & Middleware Placeholders
touch internal/delivery/http/router.go
touch internal/delivery/http/middleware/auth_middleware.go

# Model, Repository, Service Placeholders
touch internal/model/user.go
touch internal/model/room.go
touch internal/repository/user_repository.go
touch internal/repository/room_repository.go
touch internal/service/user_service.go
touch internal/service/turn_service.go

# Migration Script Sample
touch scripts/migrations/000001_init_schema.up.sql
touch scripts/migrations/000001_init_schema.down.sql

# .env & .env.example
cat << 'FILE_EOF' > .env.example
APP_PORT=8080
APP_ENV=development

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=webrtc_db

JWT_SECRET=supersecretkey123
TURN_SECRET=turnsharedsecretkey
TURN_HOST=turn.example.com
FILE_EOF

cp .env.example .env

# .gitignore
cat << 'FILE_EOF' > .gitignore
bin/
dist/
build/
main
tmp/
*.exe
*.log
.env
.idea/
.vscode/
.DS_Store
FILE_EOF

# Inisialisasi go.mod jika belum ada
if [ ! -f "go.mod" ]; then
  echo "==> Menginisialisasi go.mod..."
  go mod init gin-api
fi

echo "==> Mengambil dependensi inti..."
go get -u github.com/gin-gonic/gin
go get -u github.com/gorilla/websocket
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jmoiron/sqlx
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/google/uuid

go mod tidy

echo "==> Struktur berhasil dibuat di folder: $PROJECT_DIR"
