package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"your_project/models"
)

var db *sql.DB

// InitRepository initializes the database connection
func InitRepository(database *sql.DB) {
	db = database
}

// getTableName allows environment-based table configuration
func getTableName() string {
	table := os.Getenv("ENTITY_CLASSIFICATION_TABLE")
	if table == "" {
		table = "scribe_entity_classification_dev"
	}
	return table
}

// ----------------------------------------------------
// INSERT ANALYSIS (JSON ONLY)
// ----------------------------------------------------
func InsertAnalysis(ec *models.EntityClassification) error {
	if db == nil {
		return errors.New("database not initialized")
	}

	entities, err := json.Marshal(ec.Entities)
	if err != nil {
		return err
	}

	tone, err := json.Marshal(ec.Tone)
	if err != nil {
		return err
	}

	style, err := json.Marshal(ec.Style)
	if err != nil {
		return err
	}

	safety, err := json.Marshal(ec.SafetyFlags)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
		(
			analysis_id,
			transcript_id,
			creator_id,
			entities,
			tone,
			style,
			safety_flags,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, getTableName())

	_, err = db.Exec(
		query,
		ec.AnalysisID,
		ec.TranscriptID,
		ec.CreatorID,
		entities,
		tone,
		style,
		safety,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	return err
}

// ----------------------------------------------------
// GET ANALYSIS BY ID
// ----------------------------------------------------
func GetAnalysisByID(id string) (*models.EntityClassification, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	query := fmt.Sprintf(`
		SELECT
			analysis_id,
			transcript_id,
			creator_id,
			entities,
			tone,
			style,
			safety_flags,
			created_at,
			updated_at
		FROM %s
		WHERE analysis_id = $1
	`, getTableName())

	row := db.QueryRow(query, id)

	var ec models.EntityClassification
	var entities, tone, style, safety []byte

	err := row.Scan(
		&ec.AnalysisID,
		&ec.TranscriptID,
		&ec.CreatorID,
		&entities,
		&tone,
		&style,
		&safety,
		&ec.CreatedAt,
		&ec.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("analysis not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(entities, &ec.Entities); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(tone, &ec.Tone); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(style, &ec.Style); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(safety, &ec.SafetyFlags); err != nil {
		return nil, err
	}

	return &ec, nil
}

// ----------------------------------------------------
// LIST ALL ANALYSES
// ----------------------------------------------------
func GetAllAnalyses() ([]*models.EntityClassification, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	query := fmt.Sprintf(`
		SELECT
			analysis_id,
			transcript_id,
			creator_id,
			entities,
			tone,
			style,
			safety_flags,
			created_at,
			updated_at
		FROM %s
		ORDER BY created_at DESC
	`, getTableName())

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.EntityClassification

	for rows.Next() {
		var ec models.EntityClassification
		var entities, tone, style, safety []byte

		if err := rows.Scan(
			&ec.AnalysisID,
			&ec.TranscriptID,
			&ec.CreatorID,
			&entities,
			&tone,
			&style,
			&safety,
			&ec.CreatedAt,
			&ec.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(entities, &ec.Entities); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(tone, &ec.Tone); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(style, &ec.Style); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(safety, &ec.SafetyFlags); err != nil {
			return nil, err
		}

		results = append(results, &ec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
