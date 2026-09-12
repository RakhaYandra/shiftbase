package service

import (
	"errors"
	"time"

	"github.com/RakhaYandra/shiftbase/internal/domain"
	"github.com/RakhaYandra/shiftbase/internal/repository"
)

var ErrShiftConflict = errors.New("shift bertabrakan dengan jadwal lain")

type EmployeeService struct {
	Employees EmployeeCreator
}

func (s *EmployeeService) Create(e *domain.Employee) (int64, error) {
	if _, err := time.Parse("2006-01-02", e.HireDate); err != nil {
		return 0, errors.New("hire_date harus YYYY-MM-DD")
	}
	return s.Employees.Create(e)
}

type ShiftService struct {
	Shifts ShiftStore
}

func validShift(date, start, end string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return errors.New("date harus YYYY-MM-DD")
	}
	ts, err1 := time.Parse("15:04", start)
	te, err2 := time.Parse("15:04", end)
	if err1 != nil || err2 != nil {
		return errors.New("start_time/end_time harus HH:MM")
	}
	if !te.After(ts) {
		return errors.New("end_time harus setelah start_time")
	}
	return nil
}

func (s *ShiftService) Create(r *repository.ShiftRow) (int64, error) {
	if err := validShift(r.Date, r.StartTime, r.EndTime); err != nil {
		return 0, err
	}
	ok, err := s.Shifts.Overlaps(r.EmployeeID, r.Date, r.StartTime, r.EndTime, 0)
	if err != nil {
		return 0, err
	}
	if ok {
		return 0, ErrShiftConflict
	}
	id, err := s.Shifts.Create(r)
	if errors.Is(err, repository.ErrDuplicate) {
		return 0, ErrShiftConflict
	}
	return id, err
}

func (s *ShiftService) Update(r *repository.ShiftRow) error {
	if err := validShift(r.Date, r.StartTime, r.EndTime); err != nil {
		return err
	}
	ok, err := s.Shifts.Overlaps(r.EmployeeID, r.Date, r.StartTime, r.EndTime, r.ID)
	if err != nil {
		return err
	}
	if ok {
		return ErrShiftConflict
	}
	err = s.Shifts.Update(r)
	if errors.Is(err, repository.ErrDuplicate) {
		return ErrShiftConflict
	}
	return err
}
