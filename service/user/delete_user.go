package user_service

import (
	"DevOpsMiniProject/util"

	"github.com/gofiber/fiber/v3"
)

type DeleteUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
}

func (u userService) HandleDeleteUser(c fiber.Ctx) error {
	var request DeleteUserRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		return err
	}

	err = u.userRepository.DeleteUser(request.Username)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString("Delete user : " + request.Username)
}
