package repository

import "database/sql"

type UserRow struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         string
}

type UserRepository struct{ DB *sql.DB }

func (r *UserRepository) FindByEmail(email string) (*UserRow, error) {
	var u UserRow
	err := r.DB.QueryRow(`SELECT id,email,password_hash,role FROM users WHERE email=?`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*UserRow, error) {
	var u UserRow
	err := r.DB.QueryRow(`SELECT id,email,password_hash,role FROM users WHERE id=?`, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(email, hash, role string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO users(email,password_hash,role) VALUES(?,?,?)`, email, hash, role)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
