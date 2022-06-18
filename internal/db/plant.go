package db

import (
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"os"
	"time"
)

type PlantRow struct {
	ID         string    `db:"id"`
	PlantName  string    `db:"plant_name"`
	UserID     string    `db:"user_id"`
	CategoryID string    `db:"plant_category_id"`
	Images     Images    `db:"plant_images"`
	CreatedAt  time.Time `db:"plant_created_at"`
}

type PlantCategory struct {
	ID    uint   `db:"id"`
	Icon  string `db:"icon"`
	Label string `db:"label"`
	Color string `db:"color"`
}

type Images struct {
	Files []ImageFiles `json:"files"`
}

type ImageFiles struct {
	FileName string `json:"fileName"`
}

func convertPlantRowToPlant(p PlantRow) *plant.Plant {
	return &plant.Plant{
		ID:         p.ID,
		Name:       p.PlantName,
		UserId:     p.UserID,
		CategoryID: p.CategoryID,
		Images:     imageMapper(p),
		CreatedAt:  p.CreatedAt,
	}
}

func imageMapper(p PlantRow) []plant.ImageUrls {
	baseURL := os.Getenv("FILE_SERVER_URL")
	imagesArray := p.Images.Files
	var imagesUrls []plant.ImageUrls
	for i := 0; i < len(imagesArray); i++ {
		var imageUrl plant.ImageUrls
		fileName := imagesArray[i].FileName
		url := fmt.Sprintf("%s%s_full.jpg", baseURL, fileName)
		thumbnailUrl := fmt.Sprintf("%s%s_full.jpg", baseURL, fileName)
		imageUrl.Url = url
		imageUrl.ThumbnailUrl = thumbnailUrl
		imagesUrls = append(imagesUrls, imageUrl)
	}
	return imagesUrls
}

func (d *Database) GetPlant(ctx context.Context, uuid string) (*plant.Plant, error) {
	var pr PlantRow
	var images string
	query := `SELECT plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs, plnt_created_at
				FROM plants 
				WHERE plnt_id = $1`
	row := d.Client.QueryRowContext(ctx, query, uuid)
	err := row.Scan(&pr.ID, &pr.PlantName, &pr.UserID, &pr.CategoryID, &images, &pr.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error fetching plant by uuid. %w", err)
	}
	var plantImages Images
	err = json.Unmarshal([]byte(images), &plantImages)
	if err != nil {
		return nil, err
	}
	pr.Images = plantImages
	return convertPlantRowToPlant(pr), nil
}

func (d *Database) GetPlantsByUserId(ctx context.Context, uuid string) ([]plant.Plant, error) {
	query := `SELECT plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs, plnt_created_at
				FROM plants 
				WHERE plnt_usr_id = $1`
	rows, err := d.Client.QueryContext(ctx, query, uuid)
	if err != nil {
		return nil, err
	}
	var plantList []plant.Plant
	for rows.Next() {
		pr := PlantRow{}
		var images string
		if err := rows.Scan(&pr.ID, &pr.PlantName, &pr.UserID, &pr.CategoryID, &images, &pr.CreatedAt); err != nil {
			return nil, err
		}
		var plantImages Images
		err = json.Unmarshal([]byte(images), &plantImages)
		if err != nil {
			return nil, err
		}
		pr.Images = plantImages
		p := convertPlantRowToPlant(pr)
		plantList = append(plantList, *p)
	}
	return plantList, nil
}

/*
'INSERT INTO plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ($1, $2, $3, $4, $5)'
*/

func (d *Database) AddPlant(ctx context.Context, p plant.Plant) (*plant.Plant, error) {
	query := `INSERT INTO plants (
                    plnt_id, 
                    plnt_nm, 
                    plnt_usr_id, 
                    plnt_ctgry_id, 
                    plnt_imgs) 
                VALUES ($1, $2, $3, $4, $5)`
	log.Info("db.AddPlant: Attempting to save plant to database")
	p.ID = uuid.NewV4().String()
	plantRow := PlantRow{
		ID:        p.ID,
		PlantName: p.Name,
		UserID:    p.UserId,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, plantRow)
	if err != nil {
		return nil, fmt.Errorf("FAILED to insert plant: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("FAILED to close rows: %w", err)
	}
	return &p, nil
}

func (d *Database) DeletePlant(ctx context.Context, id string) error {
	query := `DELETE FROM plants 
				WHERE plnt_id = $1`
	_, err := d.Client.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("FAILED to delete plant from database: %w", err)
	}
	return nil
}

func (d *Database) UpdatePlant(ctx context.Context, id string, p plant.Plant) (*plant.Plant, error) {
	query := `UPDATE plants SET
		plnt_nm = :plant_name
		WHERE plnt_id = :id`
	plantRow := PlantRow{
		ID:        id,
		PlantName: p.Name,
		UserID:    p.UserId,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, plantRow)
	if err != nil {
		return nil, fmt.Errorf("FAILED to update plant: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("FAILED to close rows: %w", err)
	}
	return convertPlantRowToPlant(plantRow), nil
}
