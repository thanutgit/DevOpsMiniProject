package service_status_about

import (
	"os"
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

func (s statusAbout) GetAllInfo() (statusAbout, error) {
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
	return s, nil
}

func ProvideStatusAbout() *statusAbout { return &statusAbout{} }
