package server

import (
	"DevOpsMiniProject/di/config"
	router "DevOpsMiniProject/service/router"
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/gorm"
)

func InitApiServer(db *gorm.DB, cfg *config.Config) error {
	app := fiber.New()
	app.Use(logger.New())
	router.InitRouter(app, db)

	// ctx จะถูก cancel เมื่อได้ SIGINT (Ctrl+C) หรือ SIGTERM (จาก k8s)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return app.Listen(":"+cfg.Server.AppPort, fiber.ListenConfig{
		GracefulContext: ctx,              // ctx ถูก cancel → หยุดรับ connection ใหม่ แล้วรอ request ที่ค้างอยู่
		ShutdownTimeout: 20 * time.Second, // รอได้ไม่เกิน 20s
	})
}
