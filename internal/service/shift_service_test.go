package service

import (
	"errors"
	"testing"
	"time"

	"github.com/RakhaYandra/shiftbase/internal/domain"
	"github.com/RakhaYandra/shiftbase/internal/repository"
)

// fakeShift: overlap bila start baru < 12:00 (menabrak shift 08:00-12:00).
type fakeShift struct{ dup bool }

func (f *fakeShift) Overlaps(_ int64, _, start, _ string, excludeID int64) (bool, error) {
	if excludeID == 99 {
		return false, nil // update diri sendiri: tidak konflik
	}
	return start < "12:00", nil
}
func (f *fakeShift) Create(s *repository.ShiftRow) (int64, error) {
	if f.dup {
		return 0, repository.ErrDuplicate
	}
	return 7, nil
}
func (f *fakeShift) Update(_ *repository.ShiftRow) error { return nil }

func TestShiftCreateConflict(t *testing.T) {
	svc := &ShiftService{Shifts: &fakeShift{}}
	_, err := svc.Create(&repository.ShiftRow{EmployeeID: 1, Date: "2030-01-05", StartTime: "10:00", EndTime: "14:00"})
	if !errors.Is(err, ErrShiftConflict) {
		t.Fatalf("mau ErrShiftConflict, dapat %v", err)
	}
}

func TestShiftCreateAdjacentOK(t *testing.T) {
	svc := &ShiftService{Shifts: &fakeShift{}}
	id, err := svc.Create(&repository.ShiftRow{EmployeeID: 1, Date: "2030-01-05", StartTime: "12:00", EndTime: "16:00"})
	if err != nil || id != 7 {
		t.Fatalf("mau (7,nil), dapat (%d,%v)", id, err)
	}
}

func TestShiftCreateBadInput(t *testing.T) {
	svc := &ShiftService{Shifts: &fakeShift{}}
	for _, r := range []*repository.ShiftRow{
		{Date: "05-01-2030", StartTime: "08:00", EndTime: "16:00"},
		{Date: "2030-01-05", StartTime: "16:00", EndTime: "08:00"},
		{Date: "2030-01-05", StartTime: "xx", EndTime: "16:00"},
	} {
		if _, err := svc.Create(r); err == nil {
			t.Fatalf("mau error untuk %+v", r)
		}
	}
}

func TestShiftCreateDupKeyMapsConflict(t *testing.T) {
	svc := &ShiftService{Shifts: &fakeShift{dup: true}}
	_, err := svc.Create(&repository.ShiftRow{EmployeeID: 1, Date: "2030-01-05", StartTime: "16:00", EndTime: "20:00"})
	if !errors.Is(err, ErrShiftConflict) {
		t.Fatalf("mau ErrShiftConflict, dapat %v", err)
	}
}

func TestShiftUpdateSelfOK(t *testing.T) {
	svc := &ShiftService{Shifts: &fakeShift{}}
	err := svc.Update(&repository.ShiftRow{ID: 99, EmployeeID: 1, Date: "2030-01-05", StartTime: "08:00", EndTime: "12:00"})
	if err != nil {
		t.Fatalf("update diri sendiri tak boleh konflik: %v", err)
	}
}

type fakeEmp struct{ dupOn string }

func (f *fakeEmp) Create(e *domain.Employee) (int64, error) {
	if e.Email == f.dupOn {
		return 0, repository.ErrDuplicate
	}
	if _, err := time.Parse("2006-01-02", e.HireDate); err != nil {
		return 0, err
	}
	return 3, nil
}
