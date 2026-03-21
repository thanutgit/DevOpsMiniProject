package service_status_home

import (
	"fmt"
	"os"

	"DevOpsMiniProject/util"

	"DevOpsMiniProject/repository"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type StatusHome interface {
	GetAllInfo(c fiber.Ctx) error
}

type statusHome struct {
	db             *gorm.DB
	userRepository repository.UserRepository
}

func (s statusHome) GetAllInfo(c fiber.Ctx) error {
	util.IncrementRequest()
	environment := os.Getenv("APP_ENV")
	status := "Running"
	version := util.Version
	hostname, _ := os.Hostname()
	uptime := util.Uptime()
	totalRequests := int64(util.GetTotalRequests())
	dbStatus := s.userRepository.GetStatusDB()
	redisStatus := "🟡Not Configured"
	buildTime := util.Buildtime()
	startTime := util.StartTime()
	processID := os.Getpid()

	output := fmt.Sprintf(`
	Mini Production Platform
	=========================

	Runtime Overview
	-------------------------
	Status: %s
	Environment: %s
	Version: %s
	Build Time: %s 

	Pod Information
	-------------------------
	Hostname (Pod): %s
	Process ID: %d
	Start Time: %s
	Uptime: %s
	Total Requests: %d

	Dependencies
	-------------------------
	Database Status: %s
	Redis Status: %s
	`, status, environment, version, buildTime, hostname, processID, startTime, uptime, totalRequests, dbStatus, redisStatus)

	return c.SendString(output)
}

func ProvideHomeService(db *gorm.DB) StatusHome {
	return &statusHome{
		db:             db,
		userRepository: repository.ProvideUserRepository(db),
	}
}
