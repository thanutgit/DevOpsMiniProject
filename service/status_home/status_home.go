package service_status_home

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3"
)

type statusHome struct {
	Status         string
	Version        string
	Environment    string
	Hostname       string
	Uptime         string
	TotalRequests  int
	DatabaseStatus string
	RedisStatus    string
	BuildTime      string
	CurrentTime    string
	StartTime      string
	ProcessID      int
}

var (
	startTime     = time.Now()
	totalRequests uint64
	buildTime     string
)

func GetRuntimeUptime() string {
	duration := time.Since(startTime)

	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
}

func IncrementRequest() {
	atomic.AddUint64(&totalRequests, 1)
}

func GetTotalRequests() uint64 {
	return atomic.LoadUint64(&totalRequests)
}

func (s statusHome) GetAllInfo(c fiber.Ctx) error {
	IncrementRequest()
	env := os.Getenv("APP_ENV")
	s.Environment = env
	s.Status = "Running"
	s.Version = "v.1.0.0"
	s.Hostname, _ = os.Hostname()
	s.Uptime = GetRuntimeUptime()
	s.TotalRequests = int(GetTotalRequests())
	s.DatabaseStatus = "Connected"
	s.RedisStatus = "Connected"
	s.BuildTime = "2026-03-03T10:15:00Z"
	s.CurrentTime = "2026-03-03T10:15:00Z"
	s.StartTime = "2026-03-03T10:15:00Z"
	s.ProcessID = os.Getpid()

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
	Current Time: %s
	Uptime: %s
	Total Requests: %d

	Dependencies
	-------------------------
	Database Status: %s
	Redis Status: %s
	`, s.Status, s.Environment, s.Version, s.BuildTime, s.Hostname, s.ProcessID,
		s.StartTime, s.CurrentTime, s.Uptime, s.TotalRequests, s.DatabaseStatus, s.RedisStatus)

	return c.SendString(output)
}

func ProvideHomeService() *statusHome { return &statusHome{} }
