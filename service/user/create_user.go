package user_service

import (
	"DevOpsMiniProject/entity"
	"DevOpsMiniProject/util"

	"github.com/gofiber/fiber/v3"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age"`
}

func (u userService) HandleCreateUser(c fiber.Ctx) error {
	var request CreateUserRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		return err
	}

	createdUser, err := u.userRepository.CreateUser(entity.User{
		Username: request.Username,
		Name:     request.Name,
		Age:      request.Age,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(createdUser)
}
