//go:build integration

package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"math/rand"
	"testing"
)

func TestCareLogDatabase(t *testing.T) {

	t.Run("test create care log entry", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		userId := 5
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			Name:       "testPlant",
			UserId:     userId,
			FileNames:  []string{"file1.jpg"},
			CategoryID: "21",
		})
		assert.NoError(t, err)

		//Log a plant care entry
		notes := "I think I may have over-watered this plant today"
		entry := care.LogEntry{
			PlantId:       insertedPlant.Id,
			Notes:         notes,
			WasFertilized: false,
			WasWatered:    true,
		}
		logEntry, err := db.AddCareLogEntry(context.Background(), entry)
		assert.NoError(t, err)
		assert.True(t, notes == logEntry.Notes)
		assert.Equal(t, insertedPlant.Id, logEntry.PlantId)
		assert.NotNil(t, logEntry.Date)
		assert.True(t, logEntry.WasWatered)
		assert.False(t, logEntry.WasFertilized)
	})

	t.Run("test create care log entry with no notes", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			Name:       "testPlant",
			UserId:     rand.Intn(100),
			FileNames:  []string{"file1.jpg"},
			CategoryID: "21",
		})
		assert.NoError(t, err)

		//Log a plant care entry with no notes
		entry := care.LogEntry{
			PlantId:       insertedPlant.Id,
			WasFertilized: false,
			WasWatered:    true,
		}
		logEntry, err := db.AddCareLogEntry(context.Background(), entry)
		assert.NoError(t, err)
		assert.Equal(t, "", logEntry.Notes)
		db.DeletePlant(context.Background(), insertedPlant.Id)
		db.DeleteCareLogEntry(context.Background(), logEntry.Id)
	})

	t.Run("test create and get many care log entries", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			Name:       "testPlant",
			UserId:     rand.Intn(100),
			FileNames:  []string{"file1.jpg"},
			CategoryID: "21",
		})
		assert.NoError(t, err)

		numEntries := 5
		entries := make([]care.LogEntry, numEntries)
		//Insert 5 care log entries
		for i := 0; i < numEntries; i++ {
			notes := fmt.Sprintf("iteration: %d", i)
			//Log a plant care entry with no notes
			entry := care.LogEntry{
				PlantId:       insertedPlant.Id,
				WasFertilized: false,
				WasWatered:    true,
				Notes:         notes,
			}
			logEntry, err := db.AddCareLogEntry(context.Background(), entry)
			assert.NoError(t, err)
			entries[i] = *logEntry
		}
		//Query for all those entries again
		queriedEntries, err := db.GetCareLogsEntries(context.Background(), insertedPlant.Id)
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
			Name:       "testPlant",
			UserId:     rand.Intn(100),
			FileNames:  []string{"file1.jpg"},
			CategoryID: "21",
		})
		assert.NoError(t, err)

		//Log a plant care entry
		notes := "I think I may have over-watered this plant today"
		entry := care.LogEntry{
			PlantId:       insertedPlant.Id,
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
		entries, err := db.GetCareLogsEntries(context.Background(), insertedPlant.Id)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(entries))
	})

	t.Run("test update a care log entry", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		//First add a new plant to db
		insertedPlant, err := db.AddPlant(context.Background(), plant.Plant{
			Name:       "testPlant",
			UserId:     rand.Intn(100),
			FileNames:  []string{"file1.jpg"},
			CategoryID: "21",
		})
		assert.NoError(t, err)

		//Log a plant care entry
		notes := "I think I may have over-watered this plant today"
		entry := care.LogEntry{
			PlantId:       insertedPlant.Id,
			Notes:         notes,
			WasFertilized: false,
			WasWatered:    true,
		}
		logEntry, err := db.AddCareLogEntry(context.Background(), entry)
		assert.NoError(t, err)

		newNote := "Nope this plant was not over-watered"
		newLogEntry := care.LogEntry{
			PlantId:       insertedPlant.Id,
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
