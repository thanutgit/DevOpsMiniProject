package template

import (
	service_status_home "DevOpsMiniProject/service/status_home"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func Home(c fiber.Ctx) error {
	getInfo := service_status_home.ProvideStatusService()
	status, err := getInfo.GetAllInfo()
	if err != nil {
		return c.Status(500).SendString("Error getting status")
	}

	output := fmt.Sprintf(`
	Mini Production Platform
	===========================

	Runtime Overview
	---------------------------
	
	Status: 🟢%s
	Version: %s
	Environment: %s
	Build Time: %s

	Pod Information
	---------------------------
	Hostname(Pod): %s
	Process ID: %d
	Start Time: %s
	Uptime: %s
	Total Requests: %d

	Dependencies
	---------------------------
	Database Status: 🟢%s
	Redis Status: 🟢%s

	============================
	`, status.Status, status.Version, status.Environment, status.BuildTime, status.Hostname, status.ProcessID, status.StartTime, status.Uptime, status.TotalRequests, status.DatabaseStatus, status.RedisStatus)

	return c.SendString(output)
}
