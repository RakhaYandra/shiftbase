package handler

import (
	"net/http"

	"github.com/RakhaYandra/shiftbase/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(auth *AuthHandler, secret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	v1 := r.Group("/v1")
	v1.POST("/auth/register", auth.Register)
	v1.POST("/auth/login", auth.Login)
	v1.GET("/me", middleware.Auth(secret), auth.Me)
	return r
}
