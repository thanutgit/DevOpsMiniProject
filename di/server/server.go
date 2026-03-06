package server

import (
	"DevOpsMiniProject/di/config"
	router "DevOpsMiniProject/service/router"
	"DevOpsMiniProject/template"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitApiServer(db *gorm.DB) error {
	app := fiber.New()
	cfg := config.GetConfig()

	app.Get("/", template.Home)
	app.Get("/about", template.About)

	router.InitRouter(app, db)
	app.Get("/user", func(c fiber.Ctx) error {
		return template.User(c, db)
	})

	log.Fatal(app.Listen(":" + cfg.Server.AppPort))
	return nil
}
