package handler

import (
	"net/http"

	"github.com/RakhaYandra/shiftbase/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Auth     *AuthHandler
	Employee *EmployeeHandler
	Shift    *ShiftHandler
}

func NewRouter(d *Deps, secret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	v1 := r.Group("/v1")
	v1.POST("/auth/register", d.Auth.Register)
	v1.POST("/auth/login", d.Auth.Login)

	auth := v1.Group("", middleware.Auth(secret))
	auth.GET("/me", d.Auth.Me)

	emp := auth.Group("/employees")
	emp.POST("", middleware.RequireRole("admin"), d.Employee.Create)
	emp.GET("", middleware.RequireRole("admin", "manager"), d.Employee.List)
	emp.GET("/:id", middleware.RequireRole("admin", "manager"), d.Employee.Get)
	emp.PUT("/:id", middleware.RequireRole("admin"), d.Employee.Update)
	emp.DELETE("/:id", middleware.RequireRole("admin"), d.Employee.Delete)

	sh := auth.Group("/shifts")
	sh.POST("", middleware.RequireRole("admin", "manager"), d.Shift.Create)
	sh.GET("", d.Shift.List)
	sh.GET("/:id", d.Shift.Get)
	sh.PUT("/:id", middleware.RequireRole("admin", "manager"), d.Shift.Update)
	sh.DELETE("/:id", middleware.RequireRole("admin", "manager"), d.Shift.Delete)
	return r
}
