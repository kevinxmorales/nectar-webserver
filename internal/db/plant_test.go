//go:build integration

package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"testing"
)

var testPlant = plant.Plant{
	Name:       "testPlant",
	UserId:     99,
	CategoryID: "21",
	FileNames:  []string{"file1.jpg"},
}

func TestPlantDatabase(t *testing.T) {
	t.Run("test create plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		plantName := "testPlant"
		insertedPlant, err := db.AddPlant(context.Background(), testPlant)
		assert.NoError(t, err)

		newPlant, err := db.GetPlant(context.Background(), insertedPlant.Id)
		assert.NoError(t, err)
		newPlantName := newPlant.Name
		assert.Equal(t, plantName, newPlantName)
	})

	t.Run("test delete plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		p, err := db.AddPlant(context.Background(), testPlant)
		assert.NoError(t, err)

		err = db.DeletePlant(context.Background(), p.Id)
		assert.NoError(t, err)

		_, err = db.GetPlant(context.Background(), p.Id)
		assert.Error(t, err)
	})

	t.Run("test updating plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		originalName := "testPlant"
		userId := 66
		p, err := db.AddPlant(context.Background(), plant.Plant{
			Name:       originalName,
			UserId:     userId,
			CategoryID: "21",
			FileNames:  []string{"file1.jpg"},
		})
		assert.NoError(t, err)

		newName := "newPlantName"
		updatedPlant, err := db.UpdatePlant(context.Background(), p.Id, plant.Plant{
			Name:       newName,
			UserId:     userId,
			CategoryID: "21",
			FileNames:  []string{"file1.jpg"},
		})
		assert.NoError(t, err)

		assert.Equal(t, newName, updatedPlant.Name)
	})

	t.Run("test getting a plant that does not exist", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		_, err = db.GetPlant(context.Background(), 99999)
		assert.Error(t, err)
	})

	t.Run("test get plants by user id", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		numPlants := 3
		var plantList []plant.Plant
		userID := 1
		for i := 0; i < numPlants; i++ {
			name := "testPlant"
			p := plant.Plant{
				Name:       fmt.Sprintf("%s%d", name, i),
				UserId:     userID,
				CategoryID: "21",
				FileNames:  []string{"file1.jpg"},
			}
			plantList = append(plantList, p)
		}

		//insert 3 plants belonging to the same user
		for i := 0; i < numPlants; i++ {
			insertedPlant, err := db.AddPlant(context.Background(), plantList[i])
			assert.NoError(t, err)
			plantList[i].Id = insertedPlant.Id
		}
		differentUserID := 2
		//insert 1 plant that does not belong to this user
		notMyPlant, err := db.AddPlant(context.Background(), plant.Plant{
			Name:       "testPlant2",
			UserId:     differentUserID,
			CategoryID: "21",
			FileNames:  []string{"file1.jpg"},
		})

		userPlants, err := db.GetPlantsByUserId(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, numPlants, len(userPlants))

		// Assert that the array only contains plants that
		// belong to the user
		for i := 0; i < numPlants; i++ {
			assert.NotEqual(t, notMyPlant.Id, plantList[i].Id)
		}
	})

	t.Run("test getting plants by user id, where user has no plants", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		userId := 10
		_, err = db.AddPlant(context.Background(), plant.Plant{
			Name:       "testPlant",
			UserId:     userId,
			CategoryID: "21",
			FileNames:  []string{"file1.jpg"},
		})
		assert.NoError(t, err)

		otherId := 11
		plantList, err := db.GetPlantsByUserId(context.Background(), otherId)
		assert.NoError(t, err)
		// User has no plants, should return an empty slice
		assert.Equal(t, 0, len(plantList))
	})

}
