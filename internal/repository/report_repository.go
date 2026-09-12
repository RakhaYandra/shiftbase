package repository

import "database/sql"

type OvertimeRow struct {
	EmployeeID    int64
	Name          string
	TotalHours    float64
	OvertimeHours float64
}

type CoverageRow struct {
	Date      string
	Headcount int
}

type ReportRepository struct{ DB *sql.DB }

func (r *ReportRepository) Overtime(from, to string) ([]*OvertimeRow, error) {
	rows, err := r.DB.Query(`SELECT a.employee_id, e.name,
		ROUND(SUM(TIMESTAMPDIFF(SECOND,a.check_in,a.check_out))/3600, 2),
		ROUND(SUM(GREATEST(TIMESTAMPDIFF(SECOND,a.check_in,a.check_out)/3600 - 8, 0)), 2)
		FROM attendance a JOIN employees e ON e.id=a.employee_id
		WHERE a.check_out IS NOT NULL AND (?='' OR a.date>=?) AND (?='' OR a.date<=?)
		GROUP BY a.employee_id, e.name ORDER BY 4 DESC`, from, from, to, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*OvertimeRow
	for rows.Next() {
		var o OvertimeRow
		if err := rows.Scan(&o.EmployeeID, &o.Name, &o.TotalHours, &o.OvertimeHours); err != nil {
			return nil, err
		}
		out = append(out, &o)
	}
	return out, rows.Err()
}

func (r *ReportRepository) Coverage(date string) ([]*CoverageRow, error) {
	rows, err := r.DB.Query(`SELECT DATE_FORMAT(date,'%Y-%m-%d'), COUNT(DISTINCT employee_id)
		FROM shifts WHERE (?='' OR date=?) GROUP BY date ORDER BY date`, date, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*CoverageRow
	for rows.Next() {
		var c CoverageRow
		if err := rows.Scan(&c.Date, &c.Headcount); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}
