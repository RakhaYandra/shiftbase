package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"

	"github.com/RakhaYandra/shiftbase/internal/domain"
	"github.com/RakhaYandra/shiftbase/internal/repository"
)

type ImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

type ImportResult struct {
	Imported int           `json:"imported"`
	Failed   int           `json:"failed"`
	Errors   []ImportError `json:"errors"`
}

type ImportService struct {
	Employees EmployeeCreator
}

var csvHeader = []string{"name", "email", "phone", "position", "hire_date"}

func (s *ImportService) ImportCSV(r io.Reader) (*ImportResult, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true
	head, err := cr.Read()
	if err != nil {
		return nil, errors.New("csv kosong")
	}
	for i := range head {
		head[i] = strings.ToLower(strings.TrimSpace(head[i]))
	}
	if fmt.Sprint(head) != fmt.Sprint(csvHeader) {
		return nil, fmt.Errorf("header harus: %s", strings.Join(csvHeader, ","))
	}
	res := &ImportResult{}
	row := 1 // header = baris 1
	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		row++
		if err != nil {
			res.Failed++
			res.Errors = append(res.Errors, ImportError{row, "baris rusak"})
			continue
		}
		if len(rec) != 5 {
			res.Failed++
			res.Errors = append(res.Errors, ImportError{row, "kolom harus 5"})
			continue
		}
		for i := range rec {
			rec[i] = strings.TrimSpace(rec[i])
		}
		if rec[0] == "" {
			res.Failed++
			res.Errors = append(res.Errors, ImportError{row, "name wajib"})
			continue
		}
		if _, err := mail.ParseAddress(rec[1]); err != nil {
			res.Failed++
			res.Errors = append(res.Errors, ImportError{row, "email tidak valid"})
			continue
		}
		if _, err := time.Parse("2006-01-02", rec[4]); err != nil {
			res.Failed++
			res.Errors = append(res.Errors, ImportError{row, "hire_date harus YYYY-MM-DD"})
			continue
		}
		pos := rec[3]
		if pos == "" {
			pos = "staff"
		}
		_, err = s.Employees.Create(&domain.Employee{Name: rec[0], Email: rec[1],
			Phone: rec[2], Position: pos, HireDate: rec[4]})
		if errors.Is(err, repository.ErrDuplicate) {
			res.Failed++
			res.Errors = append(res.Errors, ImportError{row, "email sudah ada"})
			continue
		} else if err != nil {
			res.Failed++
			res.Errors = append(res.Errors, ImportError{row, err.Error()})
			continue
		}
		res.Imported++
	}
	return res, nil
}
