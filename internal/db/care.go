package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"time"
)

type LogEntryRow struct {
	Id            string         `db:"id"`
	PlantId       string         `db:"plant_id"`
	Notes         sql.NullString `db:"notes"`
	WasFertilized bool           `db:"was_fertilized"`
	WasWatered    bool           `db:"was_watered"`
	Date          time.Time      `db:"created_at"`
}

func mapRowsToLogEntries(rows SqlRows) ([]care.LogEntry, error) {
	// Use this method of creating a slice to ensure an empty slice is
	// returned, instead of nil, in case there is no entries in the database
	logEntries := []care.LogEntry{}

	for rows.Next() {
		var row LogEntryRow
		if err := rows.Scan(&row.Id, &row.PlantId, &row.Notes, &row.WasFertilized, &row.WasWatered, &row.Date); err != nil {
			return nil, fmt.Errorf("rows.Scan failed in db.care.mapRowsToLogEntries for %v", err)
		}
		logEntry := convertRowsToLogEntry(row)
		logEntries = append(logEntries, logEntry)
	}
	return logEntries, nil
}

// order of rows: pcl_id, pcl_plnt_id, pcl_notes, pcl_was_fertilized, pcl_was_watered, pcl_date
func mapRowsToLogEntry(rows SqlRows) (*care.LogEntry, error) {
	var row LogEntryRow
	for rows.Next() {
		if err := rows.Scan(&row.Id, &row.PlantId, &row.Notes, &row.WasFertilized, &row.WasWatered, &row.Date); err != nil {
			return nil, fmt.Errorf("rows.Scan failed in db.care.mapRowsToLogEntry for %v", err)
		}
	}
	logEntry := convertRowsToLogEntry(row)
	return &logEntry, nil
}

func convertRowsToLogEntry(row LogEntryRow) care.LogEntry {
	return care.LogEntry{
		Id:            row.Id,
		PlantId:       row.PlantId,
		WasWatered:    row.WasWatered,
		WasFertilized: row.WasFertilized,
		Date:          row.Date,
		Notes:         row.Notes.String,
	}
}

func (d *Database) GetCareLogsEntries(ctx context.Context, plantId string) ([]care.LogEntry, error) {
	query := `SELECT 
    			id, 
    			plant_id, 
    			notes, 
    			was_fertilized, 
    			was_watered, 
    			created_at
				FROM care_log
				WHERE plant_id = $1
				ORDER BY created_at`
	rows, err := d.Client.QueryContext(ctx, query, plantId)
	if err != nil {
		return nil, fmt.Errorf("QueryContext in db.care.GetCareLogsEntries for %v", err)
	}
	defer closeDbRows(rows, query)
	entries, err := mapRowsToLogEntries(rows)
	if err != nil {
		return nil, fmt.Errorf("mapRowsToLogEntries in db.care.GetCareLogsEntries for %v", err)
	}
	return entries, nil
}

func (d *Database) AddCareLogEntry(ctx context.Context, entry care.LogEntry) (*care.LogEntry, error) {
	query := `INSERT INTO care_log
				(plant_id, 
				notes, 
				was_watered, 
				was_fertilized)
				VALUES 
					(:plant_id,
					:notes,
					:was_watered,
					:was_fertilized)
				RETURNING id, plant_id, notes, was_fertilized, was_watered, created_at`
	row := LogEntryRow{
		PlantId:       entry.PlantId,
		Notes:         sql.NullString{String: entry.Notes, Valid: true},
		WasWatered:    entry.WasWatered,
		WasFertilized: entry.WasFertilized,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, row)
	if err != nil {
		return nil, fmt.Errorf("NamedQueryContext in db.care.AddCareLogEntry failed for %v", err)
	}
	defer closeDbRows(rows, query)
	logEntry, err := mapRowsToLogEntry(rows)
	if err != nil {
		return nil, fmt.Errorf("mapRowsToLogEntry in db.care.AddCareLogEntry failed for %v", err)
	}
	return logEntry, nil
}

func (d *Database) DeleteCareLogEntry(ctx context.Context, logEntryId string) error {
	query := `DELETE FROM care_log
				WHERE id = $1`
	if _, err := d.Client.ExecContext(ctx, query, logEntryId); err != nil {
		return fmt.Errorf("ExecContext in db.care.DeleteCareLogEntry failed for %v", err)
	}
	return nil
}

func (d *Database) UpdateCareLogEntry(ctx context.Context, logEntryId string, entry care.LogEntry) (*care.LogEntry, error) {
	query := `UPDATE care_log 
				SET 
				notes = :notes,
				was_watered = :was_watered,
				was_fertilized = :was_fertilized
				WHERE id = :id
				RETURNING id, plant_id, notes, was_fertilized, was_watered, created_at`
	row := LogEntryRow{
		Id:            logEntryId,
		Notes:         sql.NullString{String: entry.Notes, Valid: true},
		WasWatered:    entry.WasWatered,
		WasFertilized: entry.WasFertilized,
	}

	rows, err := d.Client.NamedQueryContext(ctx, query, row)
	if err != nil {
		return nil, fmt.Errorf("NamedQueryContext in db.care.UpdateCareLogEntry failed for %v", err)
	}
	defer closeDbRows(rows, query)
	updatedEntry, err := mapRowsToLogEntry(rows)
	if err != nil {
		return nil, fmt.Errorf("mapRowsToLogEntry in db.care.UpdateCareLogEntry failed for %v", err)
	}
	return updatedEntry, nil
}
