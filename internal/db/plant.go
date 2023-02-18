package db

import (
	"context"
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	errs "gitlab.com/kevinmorales/nectar-rest-api/internal/nectar_errors"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"time"
)

type PlantRow struct {
	PlantId          string         `db:"id"`
	UserId           string         `db:"user_id"`
	Username         string         `db:"user_name"`
	CommonName       string         `db:"common_name"`
	ScientificName   sql.NullString `db:"scientific_name"`
	Toxicity         sql.NullString `db:"toxicity"`
	CreatedAt        time.Time      `db:"created_at"`
	UserProfileImage sql.NullString `db:"profile_image"`
}

type imagesRow struct {
	Image          string `db:"image" sql:"type:text"`
	PlantId        string `db:"plant_id" sql:"type:uuid"`
	IsPrimaryImage bool   `db:"is_primary_image" sql:"type:bool"`
}

func convertPlantRowToPlant(p PlantRow) *plant.Plant {
	return &plant.Plant{
		PlantId:          p.PlantId,
		UserId:           p.UserId,
		Username:         p.Username,
		CommonName:       p.CommonName,
		ScientificName:   p.ScientificName.String,
		Toxicity:         p.Toxicity.String,
		CreatedAt:        p.CreatedAt,
		UserProfileImage: p.UserProfileImage.String,
	}
}

func (d *Database) GetPlant(ctx context.Context, id string) (*plant.Plant, error) {
	tag := "db.plant.GetPlant"
	var pr PlantRow
	query := `SELECT 
    				plant.id, 
    				plant.user_id, 
    				plant.common_name, 
    				plant.scientific_name, 
    				plant.toxicity, 
    				plant.created_at, 
    				nectar_users.username,
    				nectar_users.profile_image
				FROM plant
				JOIN nectar_users ON plant.user_id = nectar_users.id
				WHERE 1 = 1
				AND	plant.deletion_date > CURRENT_TIMESTAMP
				AND plant.id = $1`
	row := d.Client.QueryRowContext(ctx, query, id)
	err := row.Scan(&pr.PlantId, &pr.UserId, &pr.CommonName, &pr.ScientificName, &pr.Toxicity, &pr.CreatedAt, &pr.Username, &pr.UserProfileImage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &errs.NoEntityError{Message: fmt.Sprintf("no records with id: %s", id)}
		}
		return nil, fmt.Errorf("sqlx.SelectContext in %s failed for %v", tag, err)
	}
	pl := convertPlantRowToPlant(pr)
	images, err := d.getPlantImages(ctx, pl.PlantId)
	if err != nil {
		return nil, fmt.Errorf("db.plant.getPlantImages in %s failed for %v", tag, err)
	}
	pl.Images = images
	return pl, nil
}

func (d *Database) GetPlantsByUserId(ctx context.Context, id string) ([]plant.Plant, error) {
	tag := "db.plant.GetPlantsByUserId"
	query := `SELECT 
    				plant.id, 
    				plant.user_id, 
    				plant.common_name, 
    				plant.scientific_name, 
    				plant.toxicity, 
    				plant.created_at, 
    				nectar_users.username,
    				nectar_users.profile_image
				FROM plant
				JOIN nectar_users ON plant.user_id = nectar_users.id
				WHERE 1 = 1
				AND	plant.deletion_date > CURRENT_TIMESTAMP
				AND plant.user_id = $1`
	rows, err := d.Client.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("sqlx.QueryContext in %s failed for %s", tag, err.Error())
	}
	defer closeDbRows(rows, query)
	plantList := []plant.Plant{}
	for rows.Next() {
		pr := PlantRow{}
		err := rows.Scan(&pr.PlantId, &pr.UserId, &pr.CommonName, &pr.ScientificName, &pr.Toxicity, &pr.CreatedAt, &pr.Username, &pr.UserProfileImage)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan in %s failed for %v", tag, err)
		}
		p := convertPlantRowToPlant(pr)
		plantList = append(plantList, *p)
	}
	for i, p := range plantList {
		images, err := d.getPlantImages(ctx, p.PlantId)
		if err != nil {
			return nil, err
		}
		plantList[i].Images = images
	}
	return plantList, nil
}

func (d *Database) getPlantImages(ctx context.Context, plantId string) ([]string, error) {
	tag := "db.plant.getPlantImages"
	query := `SELECT image 
			  FROM plant_images 
			  WHERE 1=1
			  AND plant_id = $1
			  AND deletion_date > CURRENT_TIMESTAMP
			  LIMIT 3`
	rows, err := d.Client.QueryContext(ctx, query, plantId)
	if err != nil {
		return nil, fmt.Errorf("sqlx.QueryContext in %s failed for %v", tag, err)
	}
	defer closeDbRows(rows, query)
	images := []string{}
	for rows.Next() {
		var im string
		if err := rows.Scan(&im); err != nil {
			return nil, fmt.Errorf("rows.Scan in %s failed for %v", tag, err)
		}
		images = append(images, im)
	}
	return images, nil
}

func (d *Database) AddPlant(ctx context.Context, p plant.Plant, images []string) (*plant.Plant, error) {
	tag := "db.plant.AddPlant"
	queryToInsertPlant := `INSERT INTO plant (
                   			id,
		                    common_name,
		                    scientific_name,
		                    user_id,
                   			toxicity)
							VALUES ($1, $2, $3, $4, $5)`
	tx, err := d.Client.Beginx()
	if err != nil {
		return nil, fmt.Errorf("sqlx.Begin in %s failed for %v", tag, err)
	}
	id := uuid.NewV4().String()
	_, err = tx.ExecContext(ctx, queryToInsertPlant, id, p.CommonName, p.ScientificName, p.UserId, p.Toxicity)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("sqlx.tx.ExecContext in %s failed for %v", tag, err)
	}
	irList := []imagesRow{}
	insertImagesQuery := "INSERT INTO plant_images (image, plant_id, is_primary_image) VALUES (:image, :plant_id, :is_primary_image)"
	for i, _ := range images {
		isPrimaryImage := i == 0
		irList = append(irList, imagesRow{Image: images[i], PlantId: id, IsPrimaryImage: isPrimaryImage})
	}
	if _, err := tx.NamedExecContext(ctx, insertImagesQuery, irList); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("sqlx.tx.NamedExecContext in %s failed for %v", tag, err)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("sqlx.tx.Commit in %s failed for %v", tag, err)
	}
	p.PlantId = id
	p.Images = images
	return &p, nil
}

func (d *Database) DeletePlant(ctx context.Context, id string) error {
	tag := "db.plant.DeletePlant"
	query := `UPDATE plant
				SET deletion_date = current_timestamp
				WHERE 1=1
				AND plant.id = $1
				AND plant.user_id = $2`
	if _, err := d.Client.ExecContext(ctx, query, id, ctx.Value("userId")); err != nil {
		return fmt.Errorf("sqlx.ExecContext in %s failed for %v", tag, err)
	}
	return nil
}

func (d *Database) UpdatePlant(ctx context.Context, id string, p plant.Plant) (*plant.Plant, error) {
	tag := "db.plant.UpdatePlant"
	query := `UPDATE plant SET
				common_name = $1,
            	scientific_name = $2,
            	toxicity = $3,
            	last_update_date = current_timestamp
				WHERE 1=1
				AND plant.id = $4
				AND plant.user_id = $5`
	tx, err := d.Client.Beginx()
	if err != nil {
		return nil, fmt.Errorf("sqlx.Begin in %s failed for %v", tag, err)
	}
	tx.MustExecContext(ctx, query, p.CommonName, p.ScientificName, p.Toxicity, p.PlantId, ctx.Value("userId"))
	insertImagesQuery := "INSERT INTO plant_images (image, plant_id, is_primary_image) VALUES (:image, :plant_id, :is_primary_image) ON CONFLICT DO NOTHING"
	images := p.Images
	for i, _ := range images {
		isPrimaryImage := i == 0
		if _, err := tx.NamedExecContext(ctx, insertImagesQuery, []imagesRow{{Image: images[i], PlantId: id, IsPrimaryImage: isPrimaryImage}}); err != nil {
			log.Error(err)
			tx.Rollback()
			return nil, fmt.Errorf("sqlx.tx.NamedExecContext in %s failed for %v", tag, err)
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("sqlx.tx.Commit in %s failed for %v", tag, err)
	}
	return &p, nil
}

func (d *Database) AddPlantImageWithId(ctx context.Context, plantId string, uri string) (string, error) {
	tag := "db.plant.AddImageToPlant"
	query := `INSERT INTO plant_images(
				image,
				plant_id) 
				VALUES ($1, $2)`
	if _, err := d.Client.ExecContext(ctx, query, uri, plantId); err != nil {
		return uri, fmt.Errorf("sqlx.ExecContext in %s failed for %v", tag, err)
	}
	return uri, nil
}

func (d *Database) DeletePlantImage(ctx context.Context, plantId string, uri string) error {
	tag := "db.plant.DeletePlantImage"
	log.Infof("deleting image: %s in %s", uri, tag)
	query := `UPDATE plant_images
				SET deletion_date = current_timestamp
       			WHERE 1=1
       			AND image = $1
       			AND plant_id = $2`
	if _, err := d.Client.ExecContext(ctx, query, uri, plantId); err != nil {
		return fmt.Errorf("sqlx.ExecContext in %s failed for %v", tag, err)
	}
	return nil
}
