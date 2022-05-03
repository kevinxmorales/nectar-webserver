package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
)

type PlantRow struct {
	ID        string
	PlantName sql.NullString
	UserID    sql.NullString
}

func convertPlantRowToPlant(p PlantRow) plant.Plant {
	return plant.Plant{
		ID:     p.ID,
		Name:   p.PlantName.String,
		UserId: p.UserID.String,
	}
}

func (d *Database) GetPlant(ctx context.Context, uuid string) (plant.Plant, error) {
	var plantRow PlantRow
	query := "SELECT id, plant_name as plantName, user_id as userId" +
		" FROM plants WHERE id = $1"
	row := d.Client.QueryRowContext(ctx, query, uuid)
	err := row.Scan(&plantRow.ID, &plantRow.PlantName, &plantRow.UserID)
	if err != nil {
		return plant.Plant{}, fmt.Errorf("error fetching plant by uuid. %w", err)
	}

	return convertPlantRowToPlant(plantRow), nil
}
