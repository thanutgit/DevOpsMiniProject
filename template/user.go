package template

import (
	"DevOpsMiniProject/entity"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func User(ctx fiber.Ctx, db *gorm.DB) error {
	var users []entity.User
	db.Find(&users)
	result := ""

	for _, user := range users {
		result += fmt.Sprintf("Username : %s | Name : %s | Age : %d\n", user.Username, user.Name, user.Age)
	}
	return ctx.SendString(result)
}
