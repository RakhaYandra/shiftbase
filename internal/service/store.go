package service

import (
	"time"

	"github.com/RakhaYandra/shiftbase/internal/domain"
	"github.com/RakhaYandra/shiftbase/internal/repository"
)

// Seam unit test: service tergantung interface kecil, bukan *sql.DB.
// *repository.* tetap memenuhi interface ini (duck-typing), main.go tanpa perubahan.

// UserStore dipakai AuthService.
type UserStore interface {
	FindByEmail(email string) (*repository.UserRow, error)
	Create(email, hash, role string) (int64, error)
}

// EmployeeCreator dipakai EmployeeService & ImportService.
type EmployeeCreator interface {
	Create(e *domain.Employee) (int64, error)
}

// ShiftStore dipakai ShiftService.
type ShiftStore interface {
	Overlaps(employeeID int64, date, start, end string, excludeID int64) (bool, error)
	Create(s *repository.ShiftRow) (int64, error)
	Update(s *repository.ShiftRow) error
}

// AttendanceStore dipakai AttendanceService.
type AttendanceStore interface {
	CheckIn(employeeID int64, at time.Time) error
	CheckOut(employeeID int64, at time.Time) error
}
