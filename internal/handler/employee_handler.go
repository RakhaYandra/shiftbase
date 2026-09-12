package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/RakhaYandra/shiftbase/internal/domain"
	"github.com/RakhaYandra/shiftbase/internal/repository"
	"github.com/RakhaYandra/shiftbase/internal/service"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	Svc   *service.EmployeeService
	Repos *repository.EmployeeRepository
}

type employeeIn struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Position string `json:"position"`
	HireDate string `json:"hire_date" binding:"required"`
}

func empOut(e *domain.Employee) gin.H {
	return gin.H{"id": e.ID, "name": e.Name, "email": e.Email, "phone": e.Phone,
		"position": e.Position, "hire_date": e.HireDate}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var in employeeIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.Svc.Create(&domain.Employee{Name: in.Name, Email: in.Email, Phone: in.Phone,
		Position: orDefault(in.Position, "staff"), HireDate: in.HireDate})
	if errors.Is(err, repository.ErrDuplicate) {
		c.JSON(http.StatusConflict, gin.H{"error": "email_taken"})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, _ := h.Repos.GetByID(id)
	c.JSON(http.StatusCreated, empOut(e))
}

func (h *EmployeeHandler) List(c *gin.Context) {
	out, err := h.Repos.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	res := make([]gin.H, 0, len(out))
	for _, e := range out {
		res = append(res, empOut(e))
	}
	c.JSON(http.StatusOK, res)
}

func (h *EmployeeHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	e, err := h.Repos.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, empOut(e))
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var in employeeIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := parseDate(in.HireDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hire_date harus YYYY-MM-DD"})
		return
	}
	err := h.Repos.Update(&domain.Employee{ID: id, Name: in.Name, Email: in.Email,
		Phone: in.Phone, Position: orDefault(in.Position, "staff"), HireDate: in.HireDate})
	if errors.Is(err, repository.ErrDuplicate) {
		c.JSON(http.StatusConflict, gin.H{"error": "email_taken"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	e, err := h.Repos.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, empOut(e))
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Repos.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
