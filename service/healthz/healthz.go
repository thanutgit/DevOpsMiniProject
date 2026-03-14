package service_healthz

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type HealthzService interface {
	HandleHealthCheck(c fiber.Ctx) error
}

type healthzService struct {
	db *gorm.DB
}

func (h *healthzService) HandleHealthCheck(c fiber.Ctx) error {
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "unhealthy",
			"checks": fiber.Map{
				"database": "fail",
			},
		})
	} else {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
			"checks": fiber.Map{
				"database": "ok",
			},
		})
	}
}

func ProvideHealthzService(db *gorm.DB) HealthzService {
	return &healthzService{
		db: db,
	}
}
