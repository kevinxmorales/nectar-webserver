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
	CareDate      string         `db:"care_date"`
	CreatedAt     string         `db:"created_at"`
	PlantImage    string         `db:"plant_image"`
	PlantName     string         `db:"plant_name"`
	Notes         sql.NullString `db:"notes"`
	WasFertilized bool           `db:"was_fertilized"`
	WasWatered    bool           `db:"was_waterCed"`
}

func mapRowsToLogEntries(rows SqlRows) ([]care.LogEntry, error) {
	// Use this method of creating a slice to ensure an empty slice is
	// returned, instead of nil, in case there is no entries in the database
	logEntries := []care.LogEntry{}
	for rows.Next() {
		var log LogEntryRow
		var careDate time.Time
		var createdAt time.Time
		if err := rows.Scan(&log.Id, &log.PlantId, &log.Notes, &log.WasFertilized, &log.WasWatered, &careDate, &createdAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed in db.care.mapRowsToLogEntries for %v", err)
		}
		log.CareDate = careDate.Format(time.RFC1123)
		log.CreatedAt = createdAt.Format(time.RFC1123)
		logEntry := convertRowsToLogEntry(log)
		logEntries = append(logEntries, logEntry)
	}
	return logEntries, nil
}

// order of rows: pcl_id, pcl_plnt_id, pcl_notes, pcl_was_fertilized, pcl_was_watered, pcl_date
func mapRowsToLogEntry(rows SqlRows) (*care.LogEntry, error) {
	var log LogEntryRow
	for rows.Next() {
		var careDate time.Time
		if err := rows.Scan(&log.Id, &log.PlantId, &log.Notes, &log.WasFertilized, &log.WasWatered, &careDate, &log.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed in db.care.mapRowsToLogEntry for %v", err)
		}
		log.CareDate = careDate.Format(time.RFC1123)
	}
	logEntry := convertRowsToLogEntry(log)
	return &logEntry, nil
}

func convertRowsToLogEntry(row LogEntryRow) care.LogEntry {
	return care.LogEntry{
		Id:            row.Id,
		PlantId:       row.PlantId,
		WasWatered:    row.WasWatered,
		WasFertilized: row.WasFertilized,
		CareDate:      row.CareDate,
		Notes:         row.Notes.String,
		CreatedAt:     row.CreatedAt,
		PlantImage:    row.PlantImage,
		PlantName:     row.PlantName,
	}
}

func (d *Database) GetAllUsersCareLogEntries(ctx context.Context, userId string) ([]care.LogEntry, error) {
	tag := "db.care.GetAllUsersCareLogsEntries"
	findCareLogEntries := `SELECT 
							care_log.id, 
							care_log.plant_id, 
							care_log.notes, 
							care_log.was_fertilized, 
							care_log.was_watered, 
							care_log.care_date,
							care_log.created_at,
							p.common_name,
							pi.image
						   	FROM care_log
						   	INNER JOIN plant p on p.id = care_log.plant_id
						   	INNER JOIN nectar_users nu on p.user_id = nu.id
						   	INNER JOIN plant_images pi on pi.plant_id = p.id
						   	WHERE 1 = 1
						   	AND nu.id = $1
						   	AND pi.is_primary_image = true
						   	ORDER BY care_date DESC`
	rows, err := d.Client.QueryContext(ctx, findCareLogEntries, userId)
	if err != nil {
		return nil, fmt.Errorf("QueryContext in %s for %v", tag, err)
	}
	defer closeDbRows(rows, findCareLogEntries)

	logEntries := []care.LogEntry{}
	for rows.Next() {
		var log LogEntryRow
		var careDate time.Time
		var createdAt time.Time
		err := rows.Scan(
			&log.Id, &log.PlantId, &log.Notes, &log.WasFertilized, &log.WasWatered, &careDate, &createdAt, &log.PlantName, &log.PlantImage)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan failed in %s for %v", tag, err)
		}
		log.CareDate = careDate.Format(time.RFC1123)
		log.CreatedAt = createdAt.Format(time.RFC1123)
		logEntry := convertRowsToLogEntry(log)
		logEntries = append(logEntries, logEntry)
	}
	return logEntries, nil

}

func (d *Database) GetCareLogsEntries(ctx context.Context, plantId string) ([]care.LogEntry, error) {
	query := `SELECT 
    			id, 
    			plant_id, 
    			notes, 
    			was_fertilized, 
    			was_watered,
    			care_date,
    			created_at
				FROM care_log
				WHERE plant_id = $1
				ORDER BY created_at DESC`
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
