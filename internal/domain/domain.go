package domain

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleStaff   Role = "staff"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

type Employee struct {
	ID       int64
	UserID   *int64
	Name     string
	Email    string
	Phone    string
	Position string
	HireDate string // YYYY-MM-DD
}
