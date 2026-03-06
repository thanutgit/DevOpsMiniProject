package user_service

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func (u userService) HandleGetAllUser(c fiber.Ctx) error {
	users, err := u.userRepository.GetAllUser()
	var result string
	if err != nil {
		return err
	}

	for _, user := range users {
		result += fmt.Sprintf("Username : %s | Name : %s | Age : %d\n", user.Username, user.Name, user.Age)
	}
	
	return c.SendString(result)
}
