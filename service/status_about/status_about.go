package service_status_about

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
)

type statusAbout struct {
	ServiceName string
	Description string
	Owner       string
	Repository  string
	Version     string
	CommitSHA   string
	BuildTime   string
	GoVersion   string
	Environment string
	ClusterInfo string
	Namespace   string
}

type StatusAbout interface {
	GetAllInfo() (statusAbout, error)
}

func (s statusAbout) GetAllInfo(c fiber.Ctx) error {
	env := os.Getenv("APP_ENV")
	s.ServiceName = "go-api"
	s.Description = "A simple Go API for DevOps mini project"
	s.Owner = "Thanusu"
	s.Repository = "https://github.com/thanuth/go-api"
	s.Version = "v.1.0.0"
	s.CommitSHA = "abc123def456"
	s.BuildTime = "2024-06-01T12:00:00Z"
	s.GoVersion = "go1.20"
	s.Environment = env
	s.ClusterInfo = "local-cluster"
	s.Namespace = "default"

	output := fmt.Sprintf(`
	Mini Production Platform
	=========================

	Service Information
	-------------------------
	Service Name: %s
	Description: %s
	Owner: %s
	Repository: %s
	
	Build Information
	-------------------------
	Version: %s
	Commit SHA: %s
	Build Time: %s
	Go Version: %s

	Deployment Information
	-------------------------
	Environment: %s
	Cluster Info: %s
	Namespace: %s

	Available Endpoints
	-------------------------
	GET /home - Runtime overview and status
	GET /about - Service information and metadata

	=========================
	`, s.ServiceName, s.Description, s.Owner, s.Repository, s.Version,
		s.CommitSHA, s.BuildTime, s.GoVersion, s.Environment, s.ClusterInfo, s.Namespace)

	return c.SendString(output)
}

func ProvideStatusAbout() *statusAbout { return &statusAbout{} }
