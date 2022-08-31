package db

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"strconv"
	"time"
)

type PlantRow struct {
	PlantName  string    `db:"plant_name"`
	CategoryID string    `db:"plant_category_id"`
	Images     Images    `db:"plant_images"`
	CreatedAt  time.Time `db:"plant_created_at"`
	Id         int       `db:"id"`
	PlantId    string    `db:"plantId"`
	UserId     int       `db:"user_id"`
	Username   string    `db:"user_username"`
}

type PlantCategory struct {
	ID    uint   `db:"id"`
	Icon  string `db:"icon"`
	Label string `db:"label"`
	Color string `db:"color"`
}

type Images struct {
	Files []string `json:"files" db:"files"`
}

func convertPlantRowToPlant(p PlantRow) *plant.Plant {
	return &plant.Plant{
		Id:         p.Id,
		IdStr:      p.PlantId,
		Name:       p.PlantName,
		UserId:     p.UserId,
		Username:   p.Username,
		CategoryID: p.CategoryID,
		Images:     imageMapper(p),
		CreatedAt:  p.CreatedAt,
	}
}

func imageMapper(p PlantRow) []plant.ImageUrls {
	imagesArray := p.Images.Files
	//do not use nil slice declaration
	imagesUrls := []plant.ImageUrls{}
	for i := 0; i < len(imagesArray); i++ {
		var imageUrl plant.ImageUrls
		fileName := imagesArray[i]
		imageUrl.Url = fileName
		imageUrl.ThumbnailUrl = fileName
		imagesUrls = append(imagesUrls, imageUrl)
	}
	return imagesUrls
}

func (d *Database) GetPlant(ctx context.Context, id int) (*plant.Plant, error) {
	tag := "db.plant.GetPlant"
	var pr PlantRow
	var images string
	query := `SELECT 
    				plants.id, 
    				plants.plnt_nm, 
    				plants.plnt_usr_id, 
    				plants.plnt_ctgry_id, 
    				plants.plnt_urls, 
    				plants.plnt_created_at, 
    				users.usr_username
				FROM plants
				JOIN users ON plnt_usr_id = users.id
				WHERE plants.id = $1`
	row := d.Client.QueryRowContext(ctx, query, id)
	if err := row.Scan(&pr.Id, &pr.PlantName, &pr.UserId, &pr.CategoryID, &images, &pr.CreatedAt, &pr.Username); err != nil {
		return nil, fmt.Errorf("QueryRowContext in %s failed for %v", tag, err)
	}
	var plantImages Images
	if err := json.Unmarshal([]byte(images), &plantImages); err != nil {
		return nil, fmt.Errorf("json.Unmarshal in %s failed for %v", tag, err)
	}
	pr.Images = plantImages
	p := convertPlantRowToPlant(pr)
	return p, nil
}

func (d *Database) GetPlantsByUserId(ctx context.Context, id int) ([]plant.Plant, error) {
	tag := "db.plant.GetPlantsByUserId"
	query := `SELECT id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_urls, plnt_created_at
				FROM plants 
				WHERE plnt_usr_id = $1`
	rows, err := d.Client.QueryContext(ctx, query, id)
	defer closeDbRows(rows, query)
	if err != nil {
		return nil, fmt.Errorf("QueryContext in %s failed for %v", tag, err)
	}
	var plantList []plant.Plant
	for rows.Next() {
		pr := PlantRow{}
		var images string
		if err := rows.Scan(&pr.Id, &pr.PlantName, &pr.UserId, &pr.CategoryID, &images, &pr.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan in %s failed for %v", tag, err)
		}
		var plantImages Images
		if err := json.Unmarshal([]byte(images), &plantImages); err != nil {
			return nil, fmt.Errorf("json.Unmarshal in %s failed for %v", tag, err)
		}
		pr.Images = plantImages
		p := convertPlantRowToPlant(pr)
		plantList = append(plantList, *p)
	}
	return plantList, nil
}

func (d *Database) AddPlant(ctx context.Context, p plant.Plant) (*plant.Plant, error) {
	tag := "db.plant.AddPlant"
	query := `INSERT INTO plants (
                    plnt_nm, 
                    plnt_usr_id, 
                    plnt_ctgry_id, 
                    plnt_urls) 
                VALUES ($1, $2, $3, $4) 
                RETURNING id`
	images := Images{
		Files: p.FileNames,
	}
	imagesJson, err := json.Marshal(images)
	if err != nil {
		return nil, fmt.Errorf("json.Marshall in %s failed for %v", tag, err)
	}
	catId, err := strconv.Atoi(p.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("strconv.Atoi in %s failed for %v", tag, err)
	}
	rows, err := d.Client.QueryContext(ctx, query, p.Name, p.UserId, catId, imagesJson)
	if err != nil {
		return nil, fmt.Errorf("QueryContext in %s failed for %v", tag, err)
	}
	var plantID int
	for rows.Next() {
		if err := rows.Scan(&plantID); err != nil {
			return nil, fmt.Errorf("rows.Scan in %s failed for %v", tag, err)
		}
	}
	p.Id = plantID
	return &p, nil
}

func (d *Database) DeletePlant(ctx context.Context, id int) error {
	tag := "db.plant.DeletePlant"
	query := `DELETE FROM plants 
				WHERE id = $1`
	if _, err := d.Client.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("ExecContext in %s failed for %v", tag, err)
	}
	return nil
}

func (d *Database) UpdatePlant(ctx context.Context, id int, p plant.Plant) (*plant.Plant, error) {
	tag := "db.plant.UpdatePlant"
	query := `UPDATE plants SET
		plnt_nm = :plant_name
		WHERE plants.id = :id`
	plantRow := PlantRow{
		Id:        id,
		PlantName: p.Name,
		UserId:    p.UserId,
	}
	_, err := d.Client.NamedQueryContext(ctx, query, plantRow)
	if err != nil {
		return nil, fmt.Errorf("NamedQueryContext in %s failed for %v", tag, err)
	}
	updatedPlant := convertPlantRowToPlant(plantRow)
	return updatedPlant, nil
}
