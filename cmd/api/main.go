package main

import (
	"log"

	"github.com/RakhaYandra/shiftbase/internal/config"
	"github.com/RakhaYandra/shiftbase/internal/handler"
	"github.com/RakhaYandra/shiftbase/internal/repository"
	"github.com/RakhaYandra/shiftbase/internal/service"
)

func main() {
	cfg := config.Load()
	db, err := repository.Open(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	users := &repository.UserRepository{DB: db}
	employees := &repository.EmployeeRepository{DB: db}
	shifts := &repository.ShiftRepository{DB: db}
	attendance := &repository.AttendanceRepository{DB: db}
	reports := &repository.ReportRepository{DB: db}

	empSvc := &service.EmployeeService{Employees: employees}

	authH := &handler.AuthHandler{
		Svc:   &service.AuthService{Users: users, Secret: cfg.JWTSecret},
		Users: users,
	}
	empH := &handler.EmployeeHandler{
		Svc:      empSvc,
		Repos:    employees,
		Importer: &service.ImportService{Employees: empSvc},
	}
	shiftH := &handler.ShiftHandler{
		Svc:       &service.ShiftService{Shifts: shifts},
		Shifts:    shifts,
		Employees: employees,
	}
	attH := &handler.AttendanceHandler{
		Svc:        &service.AttendanceService{Attendance: attendance},
		Attendance: attendance,
		Employees:  employees,
	}
	repH := &handler.ReportHandler{
		Svc: &service.ReportService{Reports: reports},
	}

	r := handler.NewRouter(&handler.Deps{Auth: authH, Employee: empH, Shift: shiftH, Attendance: attH, Report: repH}, cfg.JWTSecret)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
