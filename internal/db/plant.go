package db

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
)

type PlantRow struct {
	ID        string `db:"id"`
	PlantName string `db:"plant_name"`
	UserID    string `db:"user_id"`
}

func convertPlantRowToPlant(p PlantRow) plant.Plant {
	return plant.Plant{
		ID:     p.ID,
		Name:   p.PlantName,
		UserId: p.UserID,
	}
}

func (d *Database) GetPlant(ctx context.Context, uuid string) (plant.Plant, error) {
	var plantRow PlantRow
	query := `SELECT id, plant_name as plantName, user_id as userId
				FROM plants WHERE id = $1`
	row := d.Client.QueryRowContext(ctx, query, uuid)
	err := row.Scan(&plantRow.ID, &plantRow.PlantName, &plantRow.UserID)
	if err != nil {
		return plant.Plant{}, fmt.Errorf("error fetching plant by uuid. %w", err)
	}

	return convertPlantRowToPlant(plantRow), nil
}

func (d *Database) AddPlant(ctx context.Context, p plant.Plant) (plant.Plant, error) {
	query := `INSERT INTO plants (id, plant_name, user_id)
				VALUES (:id, :plant_name, :user_id)`
	log.Info("db.AddPlant: Attempting to save plant to database")
	p.ID = uuid.NewV4().String()
	plantRow := PlantRow{
		ID:        p.ID,
		PlantName: p.Name,
		UserID:    p.UserId,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, plantRow)
	if err != nil {
		return plant.Plant{}, fmt.Errorf("FAILED to insert plant: %w", err)
	}
	if err := rows.Close(); err != nil {
		return plant.Plant{}, fmt.Errorf("FAILED to close rows: %w", err)
	}
	return p, nil
}

func (d *Database) DeletePlant(ctx context.Context, id string) error {
	query := "DELETE FROM plants where id = $1"
	_, err := d.Client.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("FAILED to delete plant from database: %w", err)
	}
	return nil
}

func (d *Database) UpdatePlant(ctx context.Context, id string, p plant.Plant) (plant.Plant, error) {
	query := `UPDATE plants SET
		plant_name = :plant_name
		WHERE id = :id`
	plantRow := PlantRow{
		ID:        id,
		PlantName: p.Name,
		UserID:    p.UserId,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, plantRow)
	if err != nil {
		return plant.Plant{}, fmt.Errorf("FAILED to update plant: %w", err)
	}
	if err := rows.Close(); err != nil {
		return plant.Plant{}, fmt.Errorf("FAILED to close rows: %w", err)
	}
	return convertPlantRowToPlant(plantRow), nil
}
