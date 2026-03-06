package util

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ValidateRequest(c fiber.Ctx, entity any) error {
	body := c.Body()

	if err := json.Unmarshal(body, entity); err != nil {
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(entity); err != nil {
		return err
	}

	return nil
}
