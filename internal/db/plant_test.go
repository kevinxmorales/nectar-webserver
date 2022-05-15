//go:build integration

package db

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"testing"
)

func TestPlantDatabase(t *testing.T) {
	t.Run("test create plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		plantName := "testPlant"
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			Name:   plantName,
			UserId: uuid.NewV4().String(),
		})
		assert.NoError(t, err)

		newPlant, err := db.GetPlant(context.Background(), insertedPlant.ID)
		assert.Equal(t, plantName, newPlant.Name)
	})

	t.Run("test delete plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		plantName := "testingPlant"
		p, err := db.AddPlant(context.Background(), plant.Plant{
			Name:   plantName,
			UserId: uuid.NewV4().String(),
		})
		assert.NoError(t, err)

		err = db.DeletePlant(context.Background(), p.ID)
		assert.NoError(t, err)

		_, err = db.GetPlant(context.Background(), p.ID)
		assert.Error(t, err)
	})

	t.Run("test updating plant", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		originalName := "testingPlant"
		userId := uuid.NewV4().String()
		p, err := db.AddPlant(context.Background(), plant.Plant{
			Name:   originalName,
			UserId: userId,
		})
		assert.NoError(t, err)

		newName := "newPlantName"
		updatedPlant, err := db.UpdatePlant(context.Background(), p.ID, plant.Plant{
			Name:   newName,
			UserId: userId,
		})
		assert.NoError(t, err)

		assert.Equal(t, newName, updatedPlant.Name)
	})

	t.Run("test getting a plant that does not exist", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		idNotInDB := uuid.NewV4().String()
		_, err = db.GetPlant(context.Background(), idNotInDB)
		assert.Error(t, err)
	})

	t.Run("test get plants by user id", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		userID := uuid.NewV4().String()
		numPlants := 3
		var plantList []plant.Plant
		for i := 0; i < numPlants; i++ {
			name := "testPlant"
			p := plant.Plant{
				Name:   name,
				UserId: userID,
			}
			plantList = append(plantList, p)
		}

		//insert 3 plants belonging to the same user
		for i := 0; i < numPlants; i++ {
			insertedPlant, err := db.AddPlant(context.Background(), plantList[i])
			assert.NoError(t, err)
			plantList[i].ID = insertedPlant.ID
		}

		//insert 1 plant that does not belong to this user
		notMyPlant, err := db.AddPlant(context.Background(), plant.Plant{
			Name:   "testPlant2",
			UserId: uuid.NewV4().String(),
		})

		userPlants, err := db.GetPlantsByUserId(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, numPlants, len(userPlants))

		// Assert that the array only contains plants that
		// belong to the user
		for i := 0; i < numPlants; i++ {
			assert.NotEqual(t, notMyPlant.ID, plantList[i].ID)
		}
	})

	t.Run("test getting plants by user id, where user has no plants", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		userId := uuid.NewV4().String()
		_, err = db.AddPlant(context.Background(), plant.Plant{
			Name:   "testPlant",
			UserId: userId,
		})
		assert.NoError(t, err)

		otherId := uuid.NewV4().String()
		plantList, err := db.GetPlantsByUserId(context.Background(), otherId)
		assert.NoError(t, err)
		// User has no plants, should return an empty slice
		assert.Equal(t, 0, len(plantList))
	})

}
