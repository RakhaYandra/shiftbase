package repository

import "database/sql"

type ShiftRow struct {
	ID         int64
	EmployeeID int64
	Date       string
	StartTime  string
	EndTime    string
}

type ShiftRepository struct{ DB *sql.DB }

func (r *ShiftRepository) Create(s *ShiftRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO shifts(employee_id,date,start_time,end_time) VALUES(?,?,?,?)`,
		s.EmployeeID, s.Date, s.StartTime, s.EndTime)
	if err != nil {
		if isDup(err) {
			return 0, ErrDuplicate
		}
		return 0, err
	}
	return res.LastInsertId()
}

const shiftCols = `id,employee_id,DATE_FORMAT(date,'%Y-%m-%d'),DATE_FORMAT(start_time,'%H:%i'),DATE_FORMAT(end_time,'%H:%i')`

func scanShift(row interface{ Scan(...any) error }) (*ShiftRow, error) {
	var s ShiftRow
	if err := row.Scan(&s.ID, &s.EmployeeID, &s.Date, &s.StartTime, &s.EndTime); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShiftRepository) GetByID(id int64) (*ShiftRow, error) {
	return scanShift(r.DB.QueryRow(`SELECT `+shiftCols+` FROM shifts WHERE id=?`, id))
}

func (r *ShiftRepository) List(date string, employeeID int64) ([]*ShiftRow, error) {
	q := `SELECT ` + shiftCols + ` FROM shifts WHERE (?='' OR date=?) AND (?=0 OR employee_id=?) ORDER BY date,start_time`
	rows, err := r.DB.Query(q, date, date, employeeID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ShiftRow
	for rows.Next() {
		s, err := scanShift(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// Overlaps: shift lain milik employee sama di tanggal sama yang jamnya tabrakan.
func (r *ShiftRepository) Overlaps(employeeID int64, date, start, end string, excludeID int64) (bool, error) {
	var id int64
	err := r.DB.QueryRow(`SELECT id FROM shifts WHERE employee_id=? AND date=? AND id<>?
		AND start_time<? AND end_time>? LIMIT 1`, employeeID, date, excludeID, end, start).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *ShiftRepository) Update(s *ShiftRow) error {
	_, err := r.DB.Exec(`UPDATE shifts SET employee_id=?,date=?,start_time=?,end_time=? WHERE id=?`,
		s.EmployeeID, s.Date, s.StartTime, s.EndTime, s.ID)
	if isDup(err) {
		return ErrDuplicate
	}
	return err
}

func (r *ShiftRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM shifts WHERE id=?`, id)
	return err
}
