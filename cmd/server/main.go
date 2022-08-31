package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/auth"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/blob"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/db"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/env"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	transportHttp "gitlab.com/kevinmorales/nectar-rest-api/internal/transport/http"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
)

func Run() error {
	log.Info("loading environment variables")
	if err := env.LoadEnvironment(); err != nil {
		return fmt.Errorf("FAILED to load environment variables %v", err)
	}
	log.Info("attempting to connect to database")
	database, err := db.NewDatabase()
	if err != nil {
		return fmt.Errorf("FAILED to connect to the database %v", err)
	}
	log.Info("attempting to run migrations")
	if err := database.MigrateDB(); err != nil {
		return fmt.Errorf("FAILED to migrate database %v", err)
	}
	log.Info("attempting to get s3 connection")
	blobStoreSession, err := blob.NewBlobStoreSession()
	if err != nil {
		return fmt.Errorf("FAILED to connect to the blob store %v", err)
	}
	log.Info("ready to start up server")
	plantService := plant.NewService(database, blobStoreSession)
	userService := user.NewService(database)
	authService := auth.NewService(database)
	careService := care.NewService(database)
	httpHandler := transportHttp.NewHandler(plantService, userService, careService, authService)

	printBanner()
	log.Info("service is ready to start :)")
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func printBanner() {
	fmt.Println(",--.  ,--.                  ,--.                    ")
	fmt.Println("|  ,'.|  |  ,---.   ,---. ,-'  '-.  ,--,--. ,--.--. ")
	fmt.Println("|  |' '  | | .-. : | .--' '-.  .-' ' ,-.  | |  .--' ")
	fmt.Println("|  | `   | \\   --. \\ `--.   |  |   \\ '-'  | |  |    ")
	fmt.Println("`--'  `--'  `----'  `---'   `--'    `--`--' `--'    ")
	fmt.Println("----------------- Nectar REST API -----------------")
}

func main() {
	log.Info("starting up application")
	if err := Run(); err != nil {
		log.Error(err)
	}
}
