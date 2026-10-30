package main

import (
	"DevOpsMiniProject/di/config"
	"DevOpsMiniProject/di/database"
	"DevOpsMiniProject/di/server"
	"fmt"
	"log"
	"os"
	_ "time/tzdata"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.GetConfig()

	db, err := database.InitDatabase(&cfg.Database)
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	defer sqlDB.Close()

	switch cfg.Server.Service {
	case "server":
		if err := server.InitApiServer(db, &cfg); err != nil {
			return fmt.Errorf("server: %w", err)
		}
		log.Println("server stopped gracefully")

	case "migrator":
		if err := database.MigrateDatabase(db); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		log.Println("migration completed")

	default:
		return fmt.Errorf("unknown APP_SERVICE %q (want server or migrator)", cfg.Server.Service)
	}
	return nil
}