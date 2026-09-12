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
	defer db.Close()

	users := &repository.UserRepository{DB: db}
	authSvc := &service.AuthService{Users: users, Secret: cfg.JWTSecret}
	authH := &handler.AuthHandler{Svc: authSvc, Users: users}

	r := handler.NewRouter(authH, cfg.JWTSecret)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
