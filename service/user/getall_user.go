package user

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
)

func (u userService) HandleGetAllUser(c fiber.Ctx) error {
	users, err := u.userRepository.GetAllUser()
	if err != nil {
		log.Printf("get all users failed: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	var result string
	if len(users) == 0 {
		return c.SendString("No user found")
	}

	for _, user := range users {
		result += fmt.Sprintf("Username : %s | Name : %s | Age : %d\n", user.Username, user.Name, user.Age)
	}

	return c.SendString(result)
}
