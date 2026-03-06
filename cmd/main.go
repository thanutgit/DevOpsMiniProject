package main

import (
	"DevOpsMiniProject/di/config"
	"DevOpsMiniProject/di/database"
	"DevOpsMiniProject/di/server"
)

func main() {
	cfg := config.GetConfig()
	db, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}

	if cfg.Server.Service == "server" {
		err := server.InitApiServer(db)
		if err != nil {
			panic(err)
		}
	} else if cfg.Server.Service == "migrator" {
		err := database.MigrateDatabase(db)
		if err != nil {
			panic(err)
		}
	}
}
