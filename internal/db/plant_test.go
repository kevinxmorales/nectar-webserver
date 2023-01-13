//go:build integration

package db

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"testing"
)

var testPlant = plant.Plant{
	CommonName:     "testPlant",
	ScientificName: "scientificName",
	Toxicity:       "not toxic",
	UserId:         uuid.NewV4().String(),
}

var images = []string{"https://s3imageurl.com", "https://s3imageurl2.com", "https://s3imageurl3.com"}

func TestPlantDatabase(t *testing.T) {
	t.Run("test create plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		commonName := "testPlant"
		p := plant.Plant{
			CommonName:     "testPlant",
			ScientificName: "scientificName",
			Toxicity:       "not toxic",
			UserId:         uuid.NewV4().String(),
			Images:         images,
		}
		insertedPlant, err := db.AddPlant(context.Background(), p, images)
		assert.NoError(t, err)

		newPlant, err := db.GetPlant(context.Background(), insertedPlant.PlantId)
		assert.NoError(t, err)
		spew.Dump(newPlant)
		newPlantName := newPlant.CommonName
		assert.Equal(t, commonName, newPlantName)
	})

	t.Run("test delete plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		p, err := db.AddPlant(context.Background(), testPlant, images)
		assert.NoError(t, err)

		err = db.DeletePlant(context.Background(), p.PlantId)
		assert.NoError(t, err)

		_, err = db.GetPlant(context.Background(), p.PlantId)
		assert.Error(t, err)
	})

	t.Run("test updating plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		originalName := "testPlant"
		userId := uuid.NewV4().String()
		p, err := db.AddPlant(context.Background(), plant.Plant{
			CommonName:     originalName,
			ScientificName: "scientificName",
			Toxicity:       "very toxic to pets",
			UserId:         userId,
		}, images)
		assert.NoError(t, err)

		newName := "newPlantName"
		updatedPlant, err := db.UpdatePlant(context.Background(), p.PlantId, plant.Plant{
			CommonName:     newName,
			UserId:         userId,
			ScientificName: "scientificName",
			Toxicity:       "very toxic to pets",
		})
		assert.NoError(t, err)

		assert.Equal(t, newName, updatedPlant.CommonName)
	})

	t.Run("test getting a plant that does not exist", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		_, err = db.GetPlant(context.Background(), uuid.NewV4().String())
		assert.Error(t, err)
	})

	t.Run("test get plants by user id", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		numPlants := 3
		var plantList []plant.Plant
		userID := uuid.NewV4().String()
		for i := 0; i < numPlants; i++ {
			name := "testPlant"
			p := plant.Plant{
				CommonName:     fmt.Sprintf("%s%d", name, i),
				ScientificName: fmt.Sprintf("scientificName%d", i),
				UserId:         userID,
				Toxicity:       "not toxic to pets",
			}
			plantList = append(plantList, p)
		}
		//insert 3 plants belonging to the same user
		for i := 0; i < numPlants; i++ {
			insertedPlant, err := db.AddPlant(context.Background(), plantList[i], images)
			assert.NoError(t, err)
			plantList[i].PlantId = insertedPlant.PlantId
		}
		differentUserID := uuid.NewV4().String()
		//insert 1 plant that does not belong to this user
		np := testPlant
		np.UserId = differentUserID
		notMyPlant, err := db.AddPlant(context.Background(), np, images)

		userPlants, err := db.GetPlantsByUserId(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, numPlants, len(userPlants))

		// Assert that the array only contains plants that
		// belong to the user
		for i := 0; i < numPlants; i++ {
			assert.NotEqual(t, notMyPlant.PlantId, plantList[i].PlantId)
		}
	})

	t.Run("test getting plants by user id, where user has no plants", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		_, err = db.AddPlant(context.Background(), testPlant, images)
		assert.NoError(t, err)

		otherId := uuid.NewV4().String()
		plantList, err := db.GetPlantsByUserId(context.Background(), otherId)
		assert.NoError(t, err)
		// User has no plants, should return an empty slice
		assert.Equal(t, 0, len(plantList))
	})

}
