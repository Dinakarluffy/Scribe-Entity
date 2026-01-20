package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"your_project/models"
)

var db *sql.DB

func InitRepository(database *sql.DB) {
	db = database
}

func InsertAnalysis(ec *models.EntityClassification) error {
	query := `
	INSERT INTO scribe_entity_classification_dev
	(analysis_id, transcript_id, creator_id,
	 entities, tone, style, safety_flags,
	 created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	entities, _ := json.Marshal(ec.Entities)
	tone, _ := json.Marshal(ec.Tone)
	style, _ := json.Marshal(ec.Style)
	safety, _ := json.Marshal(ec.SafetyFlags)

	_, err := db.Exec(
		query,
		ec.AnalysisID,
		ec.TranscriptID,
		ec.CreatorID,
		entities,
		tone,
		style,
		safety,
		time.Now(),
		time.Now(),
	)

	return err
}

func GetAnalysisByID(id string) (*models.EntityClassification, error) {
	row := db.QueryRow(`
		SELECT analysis_id, transcript_id, creator_id,
		       entities, tone, style, safety_flags,
		       created_at, updated_at
		FROM scribe_entity_classification_dev
		WHERE analysis_id=$1`, id)

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
		return nil, errors.New("analysis not found")
	}

	json.Unmarshal(entities, &ec.Entities)
	json.Unmarshal(tone, &ec.Tone)
	json.Unmarshal(style, &ec.Style)
	json.Unmarshal(safety, &ec.SafetyFlags)

	return &ec, nil
}

func GetAllAnalyses() ([]*models.EntityClassification, error) {
	rows, err := db.Query(`
		SELECT analysis_id, transcript_id, creator_id,
		       entities, tone, style, safety_flags,
		       created_at, updated_at
		FROM scribe_entity_classification_dev
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.EntityClassification

	for rows.Next() {
		var ec models.EntityClassification
		var entities, tone, style, safety []byte

		rows.Scan(
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

		json.Unmarshal(entities, &ec.Entities)
		json.Unmarshal(tone, &ec.Tone)
		json.Unmarshal(style, &ec.Style)
		json.Unmarshal(safety, &ec.SafetyFlags)

		results = append(results, &ec)
	}

	return results, nil
}
