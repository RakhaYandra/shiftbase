package repository

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoOpen = errors.New("tidak ada check-in terbuka hari ini")

type AttendanceRow struct {
	ID         int64
	EmployeeID int64
	Date       string
	CheckIn    string
	CheckOut   string
}

type AttendanceRepository struct{ DB *sql.DB }

func (r *AttendanceRepository) CheckIn(employeeID int64, at time.Time) error {
	_, err := r.DB.Exec(`INSERT INTO attendance(employee_id,date,check_in) VALUES(?,?,?)`,
		employeeID, at.Format("2006-01-02"), at)
	if isDup(err) {
		return ErrDuplicate
	}
	return err
}

func (r *AttendanceRepository) CheckOut(employeeID int64, at time.Time) error {
	res, err := r.DB.Exec(`UPDATE attendance SET check_out=? WHERE employee_id=? AND date=? AND check_out IS NULL`,
		at, employeeID, at.Format("2006-01-02"))
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNoOpen
	}
	return nil
}

const attCols = `id,employee_id,DATE_FORMAT(date,'%Y-%m-%d'),DATE_FORMAT(check_in,'%Y-%m-%d %H:%i:%s'),DATE_FORMAT(check_out,'%Y-%m-%d %H:%i:%s')`

func (r *AttendanceRepository) List(employeeID int64, from, to string) ([]*AttendanceRow, error) {
	rows, err := r.DB.Query(`SELECT `+attCols+` FROM attendance
		WHERE (?=0 OR employee_id=?) AND (?='' OR date>=?) AND (?='' OR date<=?) ORDER BY date`,
		employeeID, employeeID, from, from, to, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*AttendanceRow
	for rows.Next() {
		var a AttendanceRow
		var ci, co sql.NullString
		if err := rows.Scan(&a.ID, &a.EmployeeID, &a.Date, &ci, &co); err != nil {
			return nil, err
		}
		a.CheckIn, a.CheckOut = ci.String, co.String
		out = append(out, &a)
	}
	return out, rows.Err()
}
