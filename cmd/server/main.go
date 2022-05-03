package main

import (
	"context"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/db"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
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
	fmt.Println(plantService.GetPlant(context.Background(),
		"8f94b0f4-1630-4727-a60b-8bfb715a985a"))
	return nil
}

func main() {
	fmt.Println("Nectar REST API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
