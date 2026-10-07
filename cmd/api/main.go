package main

import (
	"fmt"
	"log"

	"gin-api/config"
	deliveryHttp "gin-api/internal/delivery/http"
	"gin-api/internal/delivery/http/handler"
	"gin-api/internal/repository"
	"gin-api/internal/service"
	"gin-api/pkg/database"
)

func main() {
	// 1. Muat Konfigurasi Lingkungan
	cfg := config.LoadConfig()

	// 2. Hubungkan ke Basis Data MySQL
	db, err := database.InitMySQL(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v", err)
	}
	defer db.Close()
	log.Println("Berhasil terhubung ke database MySQL")

	// 3. Jalankan Migrasi & Seeding Default
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Gagal menjalankan migrasi database: %v", err)
	}
	log.Println("Migrasi dan seeding role/user berhasil diverifikasi")

	// 4. Inisialisasi Lapisan Repository
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 5. Inisialisasi Lapisan Service
	roleService := service.NewRoleService(roleRepo)
	userService := service.NewUserService(userRepo, roleRepo)

	// 6. Inisialisasi Lapisan Delivery (HTTP Handlers)
	roleHandler := handler.NewRoleHandler(roleService)
	userHandler := handler.NewUserHandler(userService)

	// 7. Setup Router Gin
	r := deliveryHttp.SetupRouter(&deliveryHttp.RouterConfig{
		RoleHandler: roleHandler,
		UserHandler: userHandler,
	})

	// 8. Jalankan Server HTTP
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server Gin API berjalan di http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server berhenti dengan error: %v", err)
	}
}
