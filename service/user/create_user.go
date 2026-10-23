package user

import (
	"DevOpsMiniProject/entity"
	"DevOpsMiniProject/util"
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
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
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	createdUser, err := u.userRepository.CreateUser(entity.User{
		Username: request.Username,
		Name:     request.Name,
		Age:      request.Age,
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.NewError(fiber.StatusConflict, "username already exists") // 409
	}
	if err != nil {
		log.Printf("create user failed: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal error") // 500 อื่น ๆ
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}
