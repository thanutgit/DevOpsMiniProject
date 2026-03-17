package service_status_home

import (
	"fmt"
	"os"
	"time"

	"DevOpsMiniProject/util"

	"DevOpsMiniProject/repository"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type StatusHome interface {
	GetAllInfo(c fiber.Ctx) error
}

type statusHome struct {
	Status         string
	Version        string
	Environment    string
	Hostname       string
	Uptime         string
	TotalRequests  int64
	db             *gorm.DB
	RedisStatus    string
	BuildTime      string
	CurrentTime    string
	StartTime      string
	ProcessID      int
	userRepository repository.UserRepository
}

func (s statusHome) GetAllInfo(c fiber.Ctx) error {
	util.IncrementRequest()
	environment := os.Getenv("APP_ENV")
	status := "Running"
	version := util.Version
	hostname, _ := os.Hostname()
	uptime := time.Since(util.StartTime).String()
	totalRequests := int64(util.GetTotalRequests())
	dbStatus := s.userRepository.GetStatusDB()
	redisStatus := "🟡Not Configured"
	buildTime := util.BuildTime
	startTime := util.StartTime
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
