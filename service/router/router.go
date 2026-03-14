package router

import (
	service_healthz "DevOpsMiniProject/service/healthz"
	service_status_about "DevOpsMiniProject/service/status_about"
	service_status_home "DevOpsMiniProject/service/status_home"
	user_service "DevOpsMiniProject/service/user"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitRouter(server *fiber.App, db *gorm.DB) {
	userService := user_service.ProvideUserService(db)
	aboutService := service_status_about.ProvideStatusAbout()
	homeService := service_status_home.ProvideHomeService()
	healthzService := service_healthz.ProvideHealthzService(db)

	server.Post("/user", userService.HandleCreateUser)
	server.Delete("/user", userService.HandleDeleteUser)
	server.Get("/user", userService.HandleGetAllUser)

	server.Get("/about", aboutService.GetAllInfo)
	server.Get("/", homeService.GetAllInfo)

	server.Get("/healthz", healthzService.HandleHealthCheck)
}
