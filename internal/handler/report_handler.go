package handler

import (
	"net/http"

	"github.com/RakhaYandra/shiftbase/internal/service"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Svc *service.ReportService
}

func (h *ReportHandler) Overtime(c *gin.Context) {
	out, err := h.Svc.Overtime(c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := make([]gin.H, 0, len(out))
	for _, o := range out {
		res = append(res, gin.H{"employee_id": o.EmployeeID, "name": o.Name,
			"total_hours": o.TotalHours, "overtime_hours": o.OvertimeHours})
	}
	c.JSON(http.StatusOK, res)
}

func (h *ReportHandler) Coverage(c *gin.Context) {
	out, err := h.Svc.Coverage(c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := make([]gin.H, 0, len(out))
	for _, r := range out {
		res = append(res, gin.H{"date": r.Date, "headcount": r.Headcount})
	}
	c.JSON(http.StatusOK, res)
}
