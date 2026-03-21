package service_status_about

import (
	"DevOpsMiniProject/util"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type StatusAbout interface {
	GetAllInfo(c fiber.Ctx) error
}

type statusAbout struct{}

var routeDescriptions = map[string]string{
	"GET /":        "Runtime overview and status",
	"GET /about":   "Service information and metadata",
	"GET /user":    "Get all users",
	"POST /user":   "Create user",
	"DELETE /user": "Delete user",
	"GET /healthz": "Health check",
}

func (s statusAbout) GetAllInfo(c fiber.Ctx) error {
	//loop route
	routes := c.App().GetRoutes()
	var endPoints strings.Builder
	for _, route := range routes {
		key := fmt.Sprintf("%s %s", route.Method, route.Path)
		desc := routeDescriptions[key]
		endPoints.WriteString(fmt.Sprintf("\t%s %s - %s\n", route.Method, route.Path, desc))
	}

	env := os.Getenv("APP_ENV")
	serviceName := "Mini_Devops_Project"
	description := "A simple Go API for DevOps mini project"
	owner := "Thanut"
	repository := "https://github.com/thanutgit/DevOpsMiniProject"
	version := util.Version
	commitSHA := util.CommitSHA
	buildTime := util.Buildtime()
	goVersion := runtime.Version()
	environment := env
	clusterInfo := "k3s-local"
	namespace := os.Getenv("POD_NAMESPACE")

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
	%s

	=========================
	`, serviceName, description, owner, repository, version,
		commitSHA, buildTime, goVersion, environment, clusterInfo, namespace, endPoints.String())

	return c.SendString(output)
}

func ProvideStatusAbout() StatusAbout { return &statusAbout{} }
