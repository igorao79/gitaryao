package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gitserv/internal/api"
	"gitserv/internal/archive"
	"gitserv/internal/auth"
	"gitserv/internal/config"
	"gitserv/internal/database"
	"gitserv/internal/models"
)

func main() {
	cfg := config.Load()

	// Ensure data directories exist
	reposDir := cfg.DataDir + "/repos"
	if err := os.MkdirAll(reposDir, 0755); err != nil {
		log.Fatalf("Failed to create repos dir: %v", err)
	}

	// Initialize database (Turso if TURSO_DATABASE_URL set, otherwise local SQLite)
	dbPath := cfg.DataDir + "/gitserv.db"
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Restore repos from database (for ephemeral filesystems like Render)
	archiver := archive.NewManager(db, reposDir)
	if err := archiver.RestoreAll(); err != nil {
		fmt.Printf("Warning: restore repos: %v\n", err)
	}

	// JWT manager
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
		fmt.Println("WARNING: Using default JWT secret. Set JWT_SECRET env var in production!")
	}
	jwtMgr := auth.NewJWTManager(jwtSecret)

	// OAuth providers
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost" + cfg.ListenAddr
	}
	providers := auth.LoadProviders(backendURL)

	userStore := &models.UserStore{DB: db}
	oauthHandler := auth.NewOAuthHandler(providers, jwtMgr, userStore, cfg.FrontendURL)

	// Build router
	router := api.NewRouter(cfg, db, jwtMgr, oauthHandler, archiver)

	fmt.Printf("GitServ starting on %s\n", cfg.ListenAddr)
	fmt.Printf("  Data dir:     %s\n", cfg.DataDir)
	fmt.Printf("  Database:     %s\n", dbPath)
	fmt.Printf("  Frontend URL: %s\n", cfg.FrontendURL)
	fmt.Println()
	fmt.Printf("API:   http://localhost%s/api/repos\n", cfg.ListenAddr)
	fmt.Printf("Git:   http://localhost%s/{owner}/{repo}.git\n", cfg.ListenAddr)
	fmt.Printf("Auth:  http://localhost%s/auth/github\n", cfg.ListenAddr)
	fmt.Printf("Auth:  http://localhost%s/auth/google\n", cfg.ListenAddr)

	if err := http.ListenAndServe(cfg.ListenAddr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
