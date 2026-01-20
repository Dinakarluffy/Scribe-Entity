package models

import "time"

type EntityClassification struct {
	AnalysisID   string      `json:"analysis_id"`
	TranscriptID string      `json:"transcript_id"`
	CreatorID    string      `json:"creator_id"`
	Entities     interface{} `json:"entities"`
	Tone         interface{} `json:"tone"`
	Style        interface{} `json:"style"`
	SafetyFlags  interface{} `json:"safety_flags"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
