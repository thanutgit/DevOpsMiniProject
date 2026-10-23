package healthz

import (
	"DevOpsMiniProject/repository"
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

type HealthzService interface {
	HandleHealthCheck(c fiber.Ctx) error
	HandleLiveCheck(c fiber.Ctx) error
}

type healthzService struct {
	userRepository repository.UserRepository
}

func (h *healthzService) HandleHealthCheck(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 800*time.Millisecond)
	defer cancel()

	if err := h.userRepository.PingDB(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
			"checks": fiber.Map{"database": "fail"},
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
		"checks": fiber.Map{"database": "ok"},
	})
}

func (h *healthzService) HandleLiveCheck(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "alive"})
}

func ProvideHealthzService(repo repository.UserRepository) HealthzService {
	return &healthzService{userRepository: repo}
}
