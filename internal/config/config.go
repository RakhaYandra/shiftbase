package config

import "os"

type Config struct {
	Port         string
	DBDSN        string
	JWTSecret    string
	FrontendURL  string
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DBDSN:       getenv("DB_DSN", "shift:shiftpass@tcp(localhost:3306)/shiftbase?parseTime=true&loc=Local"),
		JWTSecret:   getenv("JWT_SECRET", "dev-secret-ganti-di-prod"),
		FrontendURL: getenv("FRONTEND_URL", "http://localhost:5173"),
	}
}
