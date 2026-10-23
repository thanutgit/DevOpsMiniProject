package user

import (
	"DevOpsMiniProject/repository"

	"github.com/gofiber/fiber/v3"
)

type UserService interface {
	HandleCreateUser(c fiber.Ctx) error
	HandleDeleteUser(c fiber.Ctx) error
	HandleGetAllUser(c fiber.Ctx) error
}

type userService struct {
	userRepository repository.UserRepository
}

func ProvideUserService(repo repository.UserRepository) UserService {
	return &userService{userRepository: repo}
}
