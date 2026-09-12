package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/RakhaYandra/shiftbase/internal/repository"
)

type fakeAtt struct{ mode string }

func (f *fakeAtt) CheckIn(_ int64, _ time.Time) error {
	if f.mode == "dup" {
		return repository.ErrDuplicate
	}
	return nil
}
func (f *fakeAtt) CheckOut(_ int64, _ time.Time) error {
	if f.mode == "empty" {
		return repository.ErrNoOpen
	}
	return nil
}

func TestCheckInDuplicate(t *testing.T) {
	svc := &AttendanceService{Attendance: &fakeAtt{mode: "dup"}}
	err := svc.CheckIn(1, BusinessNow())
	if err == nil || !strings.Contains(err.Error(), "sudah check-in") {
		t.Fatalf("mau error duplikat, dapat %v", err)
	}
}

func TestCheckInOK(t *testing.T) {
	svc := &AttendanceService{Attendance: &fakeAtt{}}
	if err := svc.CheckIn(1, BusinessNow()); err != nil {
		t.Fatalf("mau nil, dapat %v", err)
	}
}

func TestCheckOutNoOpen(t *testing.T) {
	svc := &AttendanceService{Attendance: &fakeAtt{mode: "empty"}}
	if err := svc.CheckOut(1, BusinessNow()); !errors.Is(err, repository.ErrNoOpen) {
		t.Fatalf("mau ErrNoOpen, dapat %v", err)
	}
}

func TestImportCSVMixed(t *testing.T) {
	svc := &ImportService{Employees: &fakeEmp{dupOn: "dup@example.com"}}
	in := "name,email,phone,position,hire_date\n" +
		"Ayu,ayu@example.com,0812,kasir,2024-01-15\n" +
		",noname@example.com,0812,kasir,2024-01-15\n" +
		"Budi,bukan-email,0812,kasir,2024-01-15\n" +
		"Citra,citra@example.com,0812,kasir,15-01-2024\n" +
		"Dedi,dup@example.com,0812,koki,2024-01-15\n" +
		"Eka,eka@example.com\n"
	res, err := svc.ImportCSV(strings.NewReader(in))
	if err != nil {
		t.Fatalf("mau nil, dapat %v", err)
	}
	if res.Imported != 1 || res.Failed != 5 {
		t.Fatalf("mau imported=1 failed=5, dapat %+v", res)
	}
	if len(res.Errors) != 5 || res.Errors[0].Row != 3 {
		t.Fatalf("nomor baris salah: %+v", res.Errors)
	}
}

func TestImportCSVBadHeader(t *testing.T) {
	svc := &ImportService{Employees: &fakeEmp{}}
	_, err := svc.ImportCSV(strings.NewReader("nama,email\nAyu,a@b.c\n"))
	if err == nil {
		t.Fatal("mau error header")
	}
}

func TestReportBadDate(t *testing.T) {
	svc := &ReportService{}
	if _, err := svc.Overtime("2024-13-40", ""); err == nil {
		t.Fatal("mau error from")
	}
	if _, err := svc.Coverage("kemarin"); err == nil {
		t.Fatal("mau error date")
	}
}
