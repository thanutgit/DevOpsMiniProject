package user

import (
	"DevOpsMiniProject/repository"
	"DevOpsMiniProject/util"
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

type DeleteUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
}

func (u userService) HandleDeleteUser(c fiber.Ctx) error {
	var request DeleteUserRequest

	if err := util.ValidateRequest(c, &request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := u.userRepository.DeleteUser(request.Username)
	if errors.Is(err, repository.ErrUserNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if err != nil {
		log.Printf("delete user failed: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.Status(fiber.StatusOK).SendString("Delete user : " + request.Username)
}
