package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"your_project/models"
	"your_project/repository"
)

type AnalyzeRequest struct {
	TranscriptID   string `json:"transcript_id"`
	TranscriptText string `json:"transcript_text"`
	CreatorID      string `json:"creator_id"`
}

// POST /analyze
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Limit body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Required field validation
	if req.TranscriptID == "" || req.TranscriptText == "" || req.CreatorID == "" {
		http.Error(w, "transcript_id, transcript_text and creator_id are required", http.StatusBadRequest)
		return
	}

	analysisID := uuid.New().String()

	// ---- MOCK DATA (Replace with Python Worker later) ----
	entities := map[string]interface{}{
		"people":    []string{"Tim Ferriss"},
		"tools":     []string{"Notion"},
		"brands":    []string{},
		"products":  []string{"The 4-Hour Workweek"},
		"companies": []string{},
	}

	tone := map[string]interface{}{
		"primary":    "conversational",
		"secondary":  "educational",
		"confidence": 0.88,
	}

	style := map[string]interface{}{
		"primary":    "interview",
		"confidence": 0.92,
	}

	safety := map[string]interface{}{
		"sensitive_domains": []string{"mental health"},
		"severity":          "low",
		"requires_review":   false,
	}
	// -----------------------------------------------------

	record := &models.EntityClassification{
		AnalysisID:   analysisID,
		TranscriptID: req.TranscriptID,
		CreatorID:    req.CreatorID,
		Entities:     entities,
		Tone:         tone,
		Style:        style,
		SafetyFlags:  safety,
		CreatedAt:    time.Now(),
	}

	if err := repository.InsertAnalysis(record); err != nil {
		http.Error(w, "Failed to save analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"analysis_id":  analysisID,
		"entities":     entities,
		"tone":         tone,
		"style":        style,
		"safety_flags": safety,
	})
}

// GET /results/{analysis_id}
func GetResultHandler(w http.ResponseWriter, r *http.Request) {
	analysisID := mux.Vars(r)["analysis_id"]
	if analysisID == "" {
		http.Error(w, "analysis_id is required", http.StatusBadRequest)
		return
	}

	result, err := repository.GetAnalysisByID(analysisID)
	if err != nil {
		http.Error(w, "Analysis not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET /results
func ListResultsHandler(w http.ResponseWriter, r *http.Request) {
	results, err := repository.GetAllAnalyses()
	if err != nil {
		http.Error(w, "Failed to fetch results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
