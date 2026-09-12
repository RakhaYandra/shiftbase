package service

import (
	"errors"
	"time"

	"github.com/RakhaYandra/shiftbase/internal/repository"
)

// WIB tanpa tzdata (distroless-safe): UTC+7 fixed.
var WIB = time.FixedZone("WIB", 7*3600)

func BusinessNow() time.Time { return time.Now().In(WIB) }

type AttendanceService struct {
	Attendance *repository.AttendanceRepository
}

func (s *AttendanceService) CheckIn(employeeID int64, now time.Time) error {
	err := s.Attendance.CheckIn(employeeID, now)
	if errors.Is(err, repository.ErrDuplicate) {
		return errors.New("sudah check-in hari ini")
	}
	return err
}

func (s *AttendanceService) CheckOut(employeeID int64, now time.Time) error {
	err := s.Attendance.CheckOut(employeeID, now)
	if errors.Is(err, repository.ErrNoOpen) {
		return repository.ErrNoOpen
	}
	return err
}
