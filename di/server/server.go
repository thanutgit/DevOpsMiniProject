package server

import (
	"DevOpsMiniProject/di/config"
	router "DevOpsMiniProject/service/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/gorm"
)

func InitApiServer(db *gorm.DB, cfg *config.Config) error {
	app := fiber.New()
	app.Use(logger.New()) // log ทุก request
	router.InitRouter(app, db)

	return app.Listen(":" + cfg.Server.AppPort)
}
