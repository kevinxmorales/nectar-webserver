package main

import (
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/db"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	transportHttp "gitlab.com/kevinmorales/nectar-rest-api/internal/transport/http"
)

func Run() error {
	fmt.Println("starting up application")
	database, err := db.NewDatabase()
	if err != nil {
		fmt.Println("FAILED to connect to the database")
		return err
	}
	if err := database.MigrateDB(); err != nil {
		fmt.Println("FAILED to migrate database")
		return err
	}

	plantService := plant.NewService(database)
	httpHandler := transportHttp.NewHandler(plantService)

	fmt.Println("service has successfully started :)")
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Nectar REST API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
