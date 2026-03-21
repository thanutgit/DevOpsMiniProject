package server

import (
	"DevOpsMiniProject/di/config"
	router "DevOpsMiniProject/service/router"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitApiServer(db *gorm.DB) error {
	app := fiber.New()
	cfg := config.GetConfig()

	router.InitRouter(app, db)

	log.Fatal(app.Listen(":" + cfg.Server.AppPort))
	return nil
}
