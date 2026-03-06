package router

import (
	user_service "DevOpsMiniProject/service/user"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitRouter(server *fiber.App, db *gorm.DB) {
	userService := user_service.ProvideUserService(db)

	server.Post("/user", userService.HandleCreateUser)
	server.Delete("/user", userService.HandleDeleteUser)
}
