package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/RakhaYandra/shiftbase/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeUsers struct{ hash string }

func (f *fakeUsers) FindByEmail(email string) (*repository.UserRow, error) {
	if email != "ada@example.com" {
		return nil, sql.ErrNoRows
	}
	return &repository.UserRow{ID: 9, Email: email, PasswordHash: f.hash, Role: "staff"}, nil
}
func (f *fakeUsers) Create(email, hash, role string) (int64, error) { return 10, nil }

func TestRegisterTaken(t *testing.T) {
	svc := &AuthService{Users: &fakeUsers{}, Secret: "s"}
	if _, err := svc.Register("ada@example.com", "Rahasia123"); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("mau ErrEmailTaken, dapat %v", err)
	}
}

func TestRegisterOK(t *testing.T) {
	svc := &AuthService{Users: &fakeUsers{}, Secret: "s"}
	id, err := svc.Register("baru@example.com", "Rahasia123")
	if err != nil || id != 10 {
		t.Fatalf("mau (10,nil), dapat (%d,%v)", id, err)
	}
}

func TestLogin(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.MinCost)
	svc := &AuthService{Users: &fakeUsers{hash: string(hash)}, Secret: "s3cr3t"}
	if _, err := svc.Login("tidakada@example.com", "x"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("mau invalid untuk user tak ada, dapat %v", err)
	}
	if _, err := svc.Login("ada@example.com", "salah"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("mau invalid untuk password salah, dapat %v", err)
	}
	tok, err := svc.Login("ada@example.com", "Rahasia123")
	if err != nil {
		t.Fatalf("mau nil, dapat %v", err)
	}
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(tok, claims, func(*jwt.Token) (any, error) { return []byte("s3cr3t"), nil }); err != nil {
		t.Fatalf("token tak valid: %v", err)
	}
	if claims["role"] != "staff" {
		t.Fatalf("mau role staff, dapat %v", claims["role"])
	}
}
