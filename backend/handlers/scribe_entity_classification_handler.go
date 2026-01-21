package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"your_project/repository"
)

type AnalyzeRequest struct {
	TranscriptID   string `json:"transcript_id"`
	TranscriptText string `json:"transcript_text"`
	CreatorID      string `json:"creator_id"`
}

type PythonWorkerResponse struct {
	AnalysisID   string                 `json:"analysis_id"`
	TranscriptID string                 `json:"transcript_id"`
	CreatorID    string                 `json:"creator_id"`
	Entities     map[string]interface{} `json:"entities"`
	Tone         map[string]interface{} `json:"tone"`
	Style        map[string]interface{} `json:"style"`
	SafetyFlags  map[string]interface{} `json:"safety_flags"`
	CreatedAt    string                 `json:"created_at"`
	Error        string                 `json:"error,omitempty"`
	Status       string                 `json:"status,omitempty"`
}

//
// -------------------------
// POST /api/entity-classification/analyze
// (JSON transcript – unchanged)
// -------------------------
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.TranscriptID == "" || req.TranscriptText == "" || req.CreatorID == "" {
		http.Error(w, "transcript_id, transcript_text and creator_id are required", http.StatusBadRequest)
		return
	}

	analysisID := uuid.New().String()

	tempDir := "./tmp/transcripts"
	_ = os.MkdirAll(tempDir, 0755)

	tempFile := filepath.Join(tempDir, analysisID+".txt")
	if err := os.WriteFile(tempFile, []byte(req.TranscriptText), 0644); err != nil {
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile)

	pythonScript := os.Getenv("PYTHON_WORKER_PATH")
	pythonExec := os.Getenv("PYTHON_EXEC")
	if pythonExec == "" {
		pythonExec = "python"
	}

	cmd := exec.Command(
		pythonExec,
		pythonScript,
		"--file", tempFile,
		"--analysis-id", analysisID,
		"--transcript-id", req.TranscriptID,
		"--creator-id", req.CreatorID,
	)

	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()

	if err != nil {
		http.Error(w, "Python worker failed:\n"+string(output), http.StatusInternalServerError)
		return
	}

	var pythonResp PythonWorkerResponse
	if err := json.Unmarshal(output, &pythonResp); err != nil {
		http.Error(w, "Invalid Python worker response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"analysis_id":  pythonResp.AnalysisID,
		"entities":     pythonResp.Entities,
		"tone":         pythonResp.Tone,
		"style":        pythonResp.Style,
		"safety_flags": pythonResp.SafetyFlags,
	})
}

//
// -------------------------
// POST /api/entity-classification/upload
// (FILE upload – FIXED)
// -------------------------
func UploadAndProcessFile(w http.ResponseWriter, r *http.Request) {

	log.Println("Upload request received")

	r.Body = http.MaxBytesReader(w, r.Body, 500<<20)
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempDir := "./tmp/uploads"
	_ = os.MkdirAll(tempDir, 0755)

	analysisID := uuid.New().String() // 🔑 SINGLE SOURCE OF TRUTH
	ext := filepath.Ext(header.Filename)
	filePath := filepath.Join(tempDir, analysisID+ext)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		os.Remove(filePath)
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}
	out.Close()

	pythonScript := os.Getenv("PYTHON_WORKER_PATH")
	if pythonScript == "" {
		os.Remove(filePath)
		http.Error(w, "PYTHON_WORKER_PATH env variable not set", http.StatusInternalServerError)
		return
	}

	pythonExec := os.Getenv("PYTHON_EXEC")
	if pythonExec == "" {
		pythonExec = "python"
	}

	cmd := exec.Command(
		pythonExec,
		pythonScript,
		"--file", filePath,
		"--analysis-id", analysisID,
	)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	_ = os.Remove(filePath)

	if err != nil {
		http.Error(w, "Python worker failed:\n"+string(output), http.StatusInternalServerError)
		return
	}

	// 🔒 Extract JSON only
	start := bytes.IndexByte(output, '{')
	end := bytes.LastIndexByte(output, '}')
	if start == -1 || end == -1 || end < start {
		http.Error(w, "Invalid Python worker response", http.StatusInternalServerError)
		return
	}

	var pythonResp PythonWorkerResponse
	if err := json.Unmarshal(output[start:end+1], &pythonResp); err != nil {
		http.Error(w, "Invalid Python worker response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"result": pythonResp,
	})
}

//
// -------------------------
// GET handlers
// -------------------------
func GetResultHandler(w http.ResponseWriter, r *http.Request) {
	analysisID := mux.Vars(r)["analysis_id"]
	result, err := repository.GetAnalysisByID(analysisID)
	if err != nil {
		http.Error(w, "Analysis not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func ListResultsHandler(w http.ResponseWriter, r *http.Request) {
	results, err := repository.GetAllAnalyses()
	if err != nil {
		http.Error(w, "Failed to fetch results", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}
