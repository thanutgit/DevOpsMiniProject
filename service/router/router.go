package router

import (
	"DevOpsMiniProject/repository"
	"DevOpsMiniProject/service/healthz"
	"DevOpsMiniProject/service/statusabout"
	"DevOpsMiniProject/service/statushome"
	"DevOpsMiniProject/service/user"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitRouter(server *fiber.App, db *gorm.DB) {
	userRepository := repository.ProvideUserRepository(db)
	userService := user.ProvideUserService(userRepository)
	aboutService := statusabout.ProvideStatusAbout()
	homeService := statushome.ProvideHomeService(userRepository)
	healthzService := healthz.ProvideHealthzService(userRepository)

	server.Post("/user", userService.HandleCreateUser)
	server.Delete("/user", userService.HandleDeleteUser)
	server.Get("/user", userService.HandleGetAllUser)

	server.Get("/about", aboutService.GetAllInfo)
	server.Get("/", homeService.GetAllInfo)

	server.Get("/livez", healthzService.HandleLiveCheck)    // liveness: ไม่เช็ค DB
	server.Get("/readyz", healthzService.HandleHealthCheck) // readiness: เช็ค DB
}
