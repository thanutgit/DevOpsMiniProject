package statusabout

import (
	"DevOpsMiniProject/util"
	"fmt"
	"runtime"
	"strings"

	"DevOpsMiniProject/di/config"

	"github.com/gofiber/fiber/v3"
)

type StatusAbout interface {
	GetAllInfo(c fiber.Ctx, cfg *config.Config) error
}

type statusAbout struct{}

var routeDescriptions = map[string]string{
	"GET /":        "Runtime overview and status",
	"GET /about":   "Service information and metadata",
	"GET /user":    "Get all users",
	"POST /user":   "Create user",
	"DELETE /user": "Delete user",
	"GET /healthz": "Health check",
	"GET /Grafana": "Grafana dashboard",
}

func (s statusAbout) GetAllInfo(c fiber.Ctx, cfg *config.Config) error {
	//loop route
	routes := c.App().GetRoutes()
	var endPoints strings.Builder
	for _, route := range routes {
		if route.Method == "HEAD" {
			continue
		}
		key := fmt.Sprintf("%s %s", route.Method, route.Path)
		desc, ok := routeDescriptions[key]
		if !ok {
			continue // ← ข้าม route ที่เราไม่ได้ตั้งใจโชว์
		}
		endPoints.WriteString(fmt.Sprintf("\t%s %s - %s\n", route.Method, route.Path, desc))
	}

	env := cfg.Server.Env
	serviceName := "Mini_Devops_Project"
	description := "A simple Go API for DevOps mini project"
	owner := "Thanut"
	repository := "https://github.com/thanutgit/DevOpsMiniProject"
	version := util.Version
	commitSHA := util.CommitSHA
	buildTime := util.Buildtime()
	goVersion := runtime.Version()
	environment := env
	clusterInfo := cfg.Server.Cluster
	namespace := cfg.Server.Namespace

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
