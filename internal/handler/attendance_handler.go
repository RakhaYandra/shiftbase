package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/RakhaYandra/shiftbase/internal/repository"
	"github.com/RakhaYandra/shiftbase/internal/service"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Svc        *service.AttendanceService
	Attendance *repository.AttendanceRepository
	Employees  *repository.EmployeeRepository
}

type attIn struct {
	EmployeeID int64 `json:"employee_id"`
}

// resolveEmployee: staff selalu miliknya sendiri; manager/admin wajib kirim employee_id.
func (h *AttendanceHandler) resolveEmployee(c *gin.Context, in attIn) (int64, bool) {
	if c.GetString("role") == "staff" {
		uid, _ := c.Get("userID")
		e, err := h.Employees.GetByUserID(uid.(int64))
		if err != nil {
			return 0, false
		}
		return e.ID, true
	}
	if in.EmployeeID == 0 {
		return 0, false
	}
	return in.EmployeeID, true
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	var in attIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	empID, ok := h.resolveEmployee(c, in)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id wajib"})
		return
	}
	if err := h.Svc.CheckIn(empID, service.BusinessNow()); err != nil {
		if err.Error() == "sudah check-in hari ini" {
			c.JSON(http.StatusConflict, gin.H{"error": "attendance_duplicate"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"employee_id": empID, "status": "checked_in"})
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	var in attIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	empID, ok := h.resolveEmployee(c, in)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id wajib"})
		return
	}
	if err := h.Svc.CheckOut(empID, service.BusinessNow()); err != nil {
		if errors.Is(err, repository.ErrNoOpen) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no_open_checkin"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"employee_id": empID, "status": "checked_out"})
}

func (h *AttendanceHandler) List(c *gin.Context) {
	var empID int64
	if c.GetString("role") == "staff" {
		uid, _ := c.Get("userID")
		e, err := h.Employees.GetByUserID(uid.(int64))
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, []gin.H{})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}
		empID = e.ID
	} else if v := c.Query("employee_id"); v != "" {
		empID, _ = strconv.ParseInt(v, 10, 64)
	}
	out, err := h.Attendance.List(empID, c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	res := make([]gin.H, 0, len(out))
	for _, a := range out {
		res = append(res, gin.H{"id": a.ID, "employee_id": a.EmployeeID, "date": a.Date,
			"check_in": a.CheckIn, "check_out": a.CheckOut})
	}
	c.JSON(http.StatusOK, res)
}

func (h *EmployeeHandler) Import(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file wajib (multipart field 'file')"})
		return
	}
	if fh.Size > 2<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file maksimal 2MB"})
		return
	}
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file tidak bisa dibaca"})
		return
	}
	defer f.Close()
	res, err := h.Importer.ImportCSV(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"imported": res.Imported, "failed": res.Failed, "errors": res.Errors})
}
