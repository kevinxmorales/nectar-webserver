package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/auth"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/blob"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/cache"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/db"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/health"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/messaging"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	_ "gitlab.com/kevinmorales/nectar-rest-api/internal/serialize"
	transportHttp "gitlab.com/kevinmorales/nectar-rest-api/internal/transport/http"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
)

func Run() error {
	log.Info("attempting to connect to database")
	database, err := db.NewDatabase()
	if err != nil {
		return fmt.Errorf("FAILED to connect to the database %v", err)
	}
	log.Info("attempting to run migrations")
	if err := database.MigrateDB(); err != nil {
		return fmt.Errorf("FAILED to migrate database %v", err)
	}
	log.Info("attempting to connect to cache")
	cacheClient, err := cache.NewCache()
	if err != nil {
		return fmt.Errorf("FAILED to connect to cache: %v", err)
	}
	log.Info("attempting to get s3 connection")
	blobStoreSession, err := blob.NewService()
	if err != nil {
		return fmt.Errorf("FAILED to connect to the blob store %v", err)
	}
	log.Info("attempting to set up auth client")
	authClient, err := auth.NewAuthClient()
	if err != nil {
		return fmt.Errorf("FAILED to setup the authentication client %v", err)
	}
	log.Info("ready to start up server")
	messageQueue, err := messaging.NewMessageQueue()
	if err != nil {
		return fmt.Errorf("FAILED to connect to the messaging queue %v", err)
	}
	plantService := plant.NewService(database, blobStoreSession, messageQueue)
	userService := user.NewService(database, authClient, blobStoreSession, messageQueue)
	authService := auth.NewService(database, authClient, cacheClient)
	careService := care.NewService(database)
	healthService := health.NewService(database, cacheClient)
	httpHandler := transportHttp.NewHandler(plantService, userService, careService, authService, healthService)

	printBanner()
	log.Info("service is ready to start :)")
	if err := httpHandler.Serve(); err != nil {
		return fmt.Errorf("FAILED to serve the http server: %v", err)
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
		panic(fmt.Errorf("application could not start %v", err))
	}
}
