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
		newName := "newPlantName"
		p, err := db.AddPlant(context.Background(), plant.Plant{
			Name:   originalName,
			UserId: userId,
		})
		assert.NoError(t, err)

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

}
