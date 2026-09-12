package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/RakhaYandra/shiftbase/internal/repository"
	"github.com/RakhaYandra/shiftbase/internal/service"
	"github.com/gin-gonic/gin"
)

type ShiftHandler struct {
	Svc       *service.ShiftService
	Shifts    *repository.ShiftRepository
	Employees *repository.EmployeeRepository
}

type shiftIn struct {
	EmployeeID int64  `json:"employee_id" binding:"required"`
	Date       string `json:"date" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
}

func shiftOut(s *repository.ShiftRow) gin.H {
	return gin.H{"id": s.ID, "employee_id": s.EmployeeID, "date": s.Date,
		"start_time": s.StartTime, "end_time": s.EndTime}
}

func parseDate(v string) (time.Time, error) { return time.Parse("2006-01-02", v) }

// scopeEmployee: staff hanya boleh akses datanya sendiri (via employees.user_id).
func (h *ShiftHandler) scopeEmployee(c *gin.Context) (int64, bool) {
	if c.GetString("role") != "staff" {
		if v := c.Query("employee_id"); v != "" {
			id, _ := strconv.ParseInt(v, 10, 64)
			return id, true
		}
		return 0, true
	}
	uid, _ := c.Get("userID")
	e, err := h.Employees.GetByUserID(uid.(int64))
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false
	} else if err != nil {
		return 0, false
	}
	return e.ID, true
}

func (h *ShiftHandler) Create(c *gin.Context) {
	var in shiftIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.Svc.Create(&repository.ShiftRow{EmployeeID: in.EmployeeID, Date: in.Date,
		StartTime: in.StartTime, EndTime: in.EndTime})
	if errors.Is(err, service.ErrShiftConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "shift_conflict"})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, _ := h.Shifts.GetByID(id)
	c.JSON(http.StatusCreated, shiftOut(s))
}

func (h *ShiftHandler) List(c *gin.Context) {
	empID, ok := h.scopeEmployee(c)
	if !ok {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	out, err := h.Shifts.List(c.Query("date"), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	res := make([]gin.H, 0, len(out))
	for _, s := range out {
		if c.GetString("role") == "staff" && s.EmployeeID != empID {
			continue
		}
		res = append(res, shiftOut(s))
	}
	c.JSON(http.StatusOK, res)
}

func (h *ShiftHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	s, err := h.Shifts.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	if empID, ok := h.scopeEmployee(c); !ok || (c.GetString("role") == "staff" && s.EmployeeID != empID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, shiftOut(s))
}

func (h *ShiftHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var in shiftIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Svc.Update(&repository.ShiftRow{ID: id, EmployeeID: in.EmployeeID, Date: in.Date,
		StartTime: in.StartTime, EndTime: in.EndTime})
	if errors.Is(err, service.ErrShiftConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "shift_conflict"})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.Shifts.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, shiftOut(s))
}

func (h *ShiftHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Shifts.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
