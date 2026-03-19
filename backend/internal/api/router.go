package api

import (
	"database/sql"

	"gitserv/internal/auth"
	"gitserv/internal/config"
	"gitserv/internal/githttp"
	"gitserv/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server holds shared dependencies for API handlers.
type Server struct {
	ReposDir string
	Config   *config.Config
	DB       *sql.DB
	Users    *models.UserStore
	Repos    *models.RepoStore
	JWT      *auth.JWTManager
}

// NewRouter creates the main HTTP router with all routes mounted.
func NewRouter(cfg *config.Config, db *sql.DB, jwtMgr *auth.JWTManager, oauthHandler *auth.OAuthHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL, "http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	reposDir := cfg.DataDir + "/repos"
	srv := &Server{
		ReposDir: reposDir,
		Config:   cfg,
		DB:       db,
		Users:    &models.UserStore{DB: db},
		Repos:    &models.RepoStore{DB: db},
		JWT:      jwtMgr,
	}
	gitHandler := githttp.NewHandler(reposDir)

	// OAuth routes (no auth needed)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/github", oauthHandler.GithubLogin)
		r.Get("/github/callback", oauthHandler.GithubCallback)
		r.Get("/google", oauthHandler.GoogleLogin)
		r.Get("/google/callback", oauthHandler.GoogleCallback)
	})

	// REST API
	r.Route("/api", func(r chi.Router) {
		// Apply auth middleware to parse JWT (doesn't block)
		r.Use(AuthMiddleware(jwtMgr))

		// Public endpoints
		r.Get("/repos/public", srv.ListPublicRepos)

		// Repository browsing (public, no auth required)
		r.Get("/repos/{owner}/{name}/tree/{ref}", srv.GetTree)
		r.Get("/repos/{owner}/{name}/tree/{ref}/*", srv.GetTree)
		r.Get("/repos/{owner}/{name}/blob/{ref}/*", srv.GetBlob)
		r.Get("/repos/{owner}/{name}/commits/{ref}", srv.GetCommits)
		r.Get("/repos/{owner}/{name}/branches", srv.GetBranches)

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(RequireAuth)
			r.Post("/repos", srv.CreateRepo)
			r.Get("/repos", srv.ListMyRepos)
			r.Get("/user", oauthHandler.CurrentUser)
		})
	})

	// Git Smart HTTP Protocol
	r.Route("/{owner}/{repo}.git", func(r chi.Router) {
		r.Get("/info/refs", gitHandler.InfoRefs)
		r.Post("/git-upload-pack", gitHandler.UploadPack)
		r.Post("/git-receive-pack", gitHandler.ReceivePack)
	})

	return r
}
