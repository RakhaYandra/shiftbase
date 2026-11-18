package config

import (
	"strings"
	"testing"
)

func TestLoadRejectsMissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("JWT_SECRET kosong harus ditolak — tidak boleh ada default yang bisa dipakai")
	}
	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("pesan error tidak menyebut JWT_SECRET: %v", err)
	}
}

func TestLoadRejectsRepoDefault(t *testing.T) {
	t.Setenv("JWT_SECRET", "dev-secret-ganti-di-prod") // nilai lama yang tertanam di repo

	_, err := Load()
	if err == nil {
		t.Fatal("kunci bawaan repo harus ditolak")
	}
	if !strings.Contains(err.Error(), "terlalu pendek") {
		t.Fatalf("error harus menyebut panjang, bukan sekadar kosong: %v", err)
	}
}

func TestLoadAcceptsStrongSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("a", MinJWTSecretBytes))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("kunci %d byte harus diterima: %v", MinJWTSecretBytes, err)
	}
	if len(cfg.JWTSecret) != MinJWTSecretBytes {
		t.Fatalf("secret tidak diteruskan: %d", len(cfg.JWTSecret))
	}
}

func TestLoadBoundaryLength(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("a", MinJWTSecretBytes-1))
	if _, err := Load(); err == nil {
		t.Fatal("1 byte di bawah minimum harus ditolak")
	}

	t.Setenv("JWT_SECRET", strings.Repeat("a", MinJWTSecretBytes))
	if _, err := Load(); err != nil {
		t.Fatalf("tepat di minimum harus diterima: %v", err)
	}
}
