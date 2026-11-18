package config

import (
	"fmt"
	"os"
)

// MinJWTSecretBytes adalah panjang minimum kunci HS256. HMAC-SHA256 memakai
// kunci 256-bit; kunci yang lebih pendek diperpanjang secara implisit dan
// lemah. 32 byte = 32 karakter ASCII.
const MinJWTSecretBytes = 32

type Config struct {
	Port        string
	DBDSN       string
	JWTSecret   string
	FrontendURL string
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// Load mengembalikan error kalau JWT_SECRET tidak diisi. Sengaja tidak ada
// default: kunci yang tertanam di repo berarti siapa pun yang membacanya
// bisa men forging token dengan role admin.
func Load() (Config, error) {
	cfg := Config{
		Port:        getenv("PORT", "8080"),
		DBDSN:       getenv("DB_DSN", "shift:shiftpass@tcp(localhost:3306)/shiftbase?parseTime=true&loc=Local"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		FrontendURL: getenv("FRONTEND_URL", "http://localhost:5173"),
	}

	if cfg.JWTSecret == "" {
		return cfg, fmt.Errorf("JWT_SECRET wajib diisi (minimal %d karakter, mis. openssl rand -hex 32)", MinJWTSecretBytes)
	}
	if len(cfg.JWTSecret) < MinJWTSecretBytes {
		return cfg, fmt.Errorf("JWT_SECRET terlalu pendek: %d byte, minimal %d", len(cfg.JWTSecret), MinJWTSecretBytes)
	}

	return cfg, nil
}
