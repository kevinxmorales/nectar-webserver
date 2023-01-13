//go:build integration

package db

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"testing"
)

func TestCareLogDatabase(t *testing.T) {

	t.Run("test create care log entry", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		userId := uuid.NewV4().String()
		ctx := context.WithValue(context.Background(), "userId", userId)
		insertedPlant, err := db.AddPlant(ctx, plant.Plant{
			CommonName:     "testPlant",
			ScientificName: "scientificName",
			Toxicity:       "not toxic",
			UserId:         userId,
		}, []string{"imageUrl1", "imageUrl2"})
		assert.NoError(t, err)

		//Log a plant care entry
		notes := "I think I may have over-watered this plant today"
		entry := care.LogEntry{
			PlantId:       insertedPlant.PlantId,
			Notes:         notes,
			WasFertilized: false,
			WasWatered:    true,
		}
		logEntry, err := db.AddCareLogEntry(context.Background(), entry)
		assert.NoError(t, err)
		assert.True(t, notes == logEntry.Notes)
		assert.Equal(t, insertedPlant.PlantId, logEntry.PlantId)
		assert.NotNil(t, logEntry.Date)
		assert.True(t, logEntry.WasWatered)
		assert.False(t, logEntry.WasFertilized)
	})

	t.Run("test create care log entry with no notes", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			CommonName:     "testPlant",
			ScientificName: "scientificName",
			Toxicity:       "not toxic",
			UserId:         uuid.NewV4().String(),
		}, []string{})
		assert.NoError(t, err)

		//Log a plant care entry with no notes
		entry := care.LogEntry{
			PlantId:       insertedPlant.PlantId,
			WasFertilized: false,
			WasWatered:    true,
		}
		logEntry, err := db.AddCareLogEntry(context.Background(), entry)
		assert.NoError(t, err)
		assert.Equal(t, "", logEntry.Notes)
		db.DeletePlant(context.Background(), insertedPlant.PlantId)
		db.DeleteCareLogEntry(context.Background(), logEntry.Id)
	})

	t.Run("test create and get many care log entries", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			CommonName:     "testPlant",
			ScientificName: "scientificName",
			Toxicity:       "not toxic",
			UserId:         uuid.NewV4().String(),
		}, []string{})
		assert.NoError(t, err)

		numEntries := 5
		entries := make([]care.LogEntry, numEntries)
		//Insert 5 care log entries
		for i := 0; i < numEntries; i++ {
			notes := fmt.Sprintf("iteration: %d", i)
			//Log a plant care entry with no notes
			entry := care.LogEntry{
				PlantId:       insertedPlant.PlantId,
				WasFertilized: false,
				WasWatered:    true,
				Notes:         notes,
			}
			logEntry, err := db.AddCareLogEntry(context.Background(), entry)
			assert.NoError(t, err)
			entries[i] = *logEntry
		}
		//Query for all those entries again
		queriedEntries, err := db.GetCareLogsEntries(context.Background(), insertedPlant.PlantId)
		assert.NoError(t, err)
		assert.Equal(t, numEntries, len(queriedEntries))

		//Confirm that these are all the same entries we made earlier
		for i := 0; i < numEntries; i++ {
			assert.Equal(t, entries[i], queriedEntries[i])
		}
	})

	t.Run("test delete a care log entry", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			CommonName:     "testPlant",
			ScientificName: "scientificName",
			Toxicity:       "not toxic",
			UserId:         uuid.NewV4().String(),
		}, []string{})
		assert.NoError(t, err)

		//Log a plant care entry
		notes := "I think I may have over-watered this plant today"
		entry := care.LogEntry{
			PlantId:       insertedPlant.PlantId,
			Notes:         notes,
			WasFertilized: false,
			WasWatered:    true,
		}
		logEntry, err := db.AddCareLogEntry(context.Background(), entry)
		assert.NoError(t, err)

		//Delete this care log entry
		err = db.DeleteCareLogEntry(context.Background(), logEntry.Id)
		assert.NoError(t, err)

		//Try to query it again, this should be empty
		entries, err := db.GetCareLogsEntries(context.Background(), insertedPlant.PlantId)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(entries))
	})

	t.Run("test update a care log entry", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			CommonName:     "testPlant",
			ScientificName: "scientificName",
			Toxicity:       "not toxic",
			UserId:         uuid.NewV4().String(),
		}, []string{})
		assert.NoError(t, err)

		//Log a plant care entry
		notes := "I think I may have over-watered this plant today"
		entry := care.LogEntry{
			PlantId:       insertedPlant.PlantId,
			Notes:         notes,
			WasFertilized: false,
			WasWatered:    true,
		}
		logEntry, err := db.AddCareLogEntry(context.Background(), entry)
		assert.NoError(t, err)

		newNote := "Nope this plant was not over-watered"
		newLogEntry := care.LogEntry{
			PlantId:       insertedPlant.PlantId,
			Notes:         newNote,
			WasFertilized: true,
			WasWatered:    false,
		}
		updatedEntry, err := db.UpdateCareLogEntry(context.Background(), logEntry.Id, newLogEntry)
		assert.NoError(t, err)

		assert.NotEqual(t, logEntry, updatedEntry)
		assert.Equal(t, newNote, updatedEntry.Notes)
	})
}
