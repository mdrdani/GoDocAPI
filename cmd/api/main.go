package main

import (
	"log"

	"godocapi/internal/config"
	"godocapi/internal/database"
	"godocapi/internal/http"
	"godocapi/internal/repository"
	"godocapi/internal/service"
	"godocapi/internal/storage"

	_ "godocapi/docs"
)

// @title DocAPI
// @version 1.1
// @description High-performance document management API
// @description High-performance document management API
// @BasePath /api/v1
func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Connect to Database
	dbPool, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 3. Initialize Storage (RustFS/S3)
	store, err := storage.NewRustFSStore(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// 4. Initialize Repository
	repo := repository.NewDocumentRepository(dbPool)

	// 5. Initialize Service
	svc := service.NewDocumentService(repo, store)

	// 6. Initialize HTTP Server
	server := http.NewServer(cfg, svc)

	// 7. Start Server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
