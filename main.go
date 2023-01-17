package main

import (
	"app-invite-service/config"
	"app-invite-service/server"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	var serverReady = make(chan bool)
	go func() {
		<-serverReady
		close(serverReady)
	}()

	server.RunMigration(cfg.MySQL.URL)
	server.Start(serverReady, cfg)
}
