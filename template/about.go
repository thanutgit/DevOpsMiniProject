package template

import (
	service_status_about "DevOpsMiniProject/service/status_about"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func About(c fiber.Ctx) error {
	getInfo := service_status_about.ProvideStatusAbout()
	status, err := getInfo.GetAllInfo()
	if err != nil {
		return c.Status(500).SendString("Error getting about info")
	}
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
	`, status.ServiceName, status.Description, status.Owner, status.Repository, status.Version, status.CommitSHA, status.BuildTime, status.GoVersion, status.Environment, status.ClusterInfo, status.Namespace)

	return c.SendString(output)
}
