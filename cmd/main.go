package main

import (
	"DevOpsMiniProject/di/config"
	"DevOpsMiniProject/di/database"
	"DevOpsMiniProject/di/server"
	"log"
	_ "time/tzdata"
)

func main() {
	cfg := config.GetConfig()
	db, err := database.InitDatabase(&cfg.Database)
	if err != nil {
		log.Panic(err.Error())
		panic(err)
	}

	if cfg.Server.Service == "server" {
		err := server.InitApiServer(db, &cfg)
		if err != nil {
			log.Panic(err.Error())
			panic(err)

		}
	} else if cfg.Server.Service == "migrator" {
		err := database.MigrateDatabase(db)
		if err != nil {
			log.Panic(err.Error())
			panic(err)
		}
	}
}
