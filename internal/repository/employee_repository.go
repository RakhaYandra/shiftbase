package repository

import (
	"database/sql"
	"errors"

	"github.com/RakhaYandra/shiftbase/internal/domain"
	"github.com/go-sql-driver/mysql"
)

var ErrDuplicate = errors.New("duplicate entry")

func isDup(err error) bool {
	var e *mysql.MySQLError
	return errors.As(err, &e) && e.Number == 1062
}

type EmployeeRepository struct{ DB *sql.DB }

func (r *EmployeeRepository) Create(e *domain.Employee) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO employees(user_id,name,email,phone,position,hire_date) VALUES(?,?,?,?,?,?)`,
		e.UserID, e.Name, e.Email, e.Phone, e.Position, e.HireDate)
	if err != nil {
		if isDup(err) {
			return 0, ErrDuplicate
		}
		return 0, err
	}
	return res.LastInsertId()
}

func scanEmployee(row interface{ Scan(...any) error }) (*domain.Employee, error) {
	var e domain.Employee
	if err := row.Scan(&e.ID, &e.UserID, &e.Name, &e.Email, &e.Phone, &e.Position, &e.HireDate); err != nil {
		return nil, err
	}
	return &e, nil
}

const empCols = `id,user_id,name,email,phone,position,DATE_FORMAT(hire_date,'%Y-%m-%d')`

func (r *EmployeeRepository) List() ([]*domain.Employee, error) {
	rows, err := r.DB.Query(`SELECT ` + empCols + ` FROM employees ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*domain.Employee
	for rows.Next() {
		e, err := scanEmployee(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *EmployeeRepository) GetByID(id int64) (*domain.Employee, error) {
	return scanEmployee(r.DB.QueryRow(`SELECT `+empCols+` FROM employees WHERE id=?`, id))
}

func (r *EmployeeRepository) GetByUserID(userID int64) (*domain.Employee, error) {
	return scanEmployee(r.DB.QueryRow(`SELECT `+empCols+` FROM employees WHERE user_id=?`, userID))
}

func (r *EmployeeRepository) Update(e *domain.Employee) error {
	_, err := r.DB.Exec(`UPDATE employees SET name=?,email=?,phone=?,position=?,hire_date=? WHERE id=?`,
		e.Name, e.Email, e.Phone, e.Position, e.HireDate, e.ID)
	if isDup(err) {
		return ErrDuplicate
	}
	return err
}

func (r *EmployeeRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM employees WHERE id=?`, id)
	return err
}
