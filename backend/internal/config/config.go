package config

import "os"

type Config struct {
	DataDir     string
	ListenAddr  string
	FrontendURL string
}

func Load() *Config {
	return &Config{
		DataDir:     getEnv("DATA_DIR", "./data"),
		ListenAddr:  getEnv("LISTEN_ADDR", ":8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
