package server

import (
	"DevOpsMiniProject/di/config"
	router "DevOpsMiniProject/service/router"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

func InitApiServer(db *gorm.DB) error {
	app := fiber.New()
	cfg := config.GetConfig()

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	router.InitRouter(app, db)

	log.Fatal(app.Listen(":" + cfg.Server.AppPort))
	return nil
}
