package service

import (
	"errors"
	"time"

	"github.com/RakhaYandra/shiftbase/internal/repository"
)

type ReportService struct {
	Reports *repository.ReportRepository
}

func checkDate(v string) error {
	if v == "" {
		return nil
	}
	_, err := time.Parse("2006-01-02", v)
	return err
}

func (s *ReportService) Overtime(from, to string) ([]*repository.OvertimeRow, error) {
	if err := checkDate(from); err != nil {
		return nil, errors.New("from harus YYYY-MM-DD")
	}
	if err := checkDate(to); err != nil {
		return nil, errors.New("to harus YYYY-MM-DD")
	}
	return s.Reports.Overtime(from, to)
}

func (s *ReportService) Coverage(date string) ([]*repository.CoverageRow, error) {
	if err := checkDate(date); err != nil {
		return nil, errors.New("date harus YYYY-MM-DD")
	}
	return s.Reports.Coverage(date)
}
