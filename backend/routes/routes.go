package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"your_project/handlers"
	"your_project/middleware"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.CORS) 
	r.Use(middleware.Logger)
	

	// Analyze transcript
	r.HandleFunc(
		"/api/entity-classification/analyze",
		handlers.AnalyzeHandler,
	).Methods(http.MethodPost)

	// Get analysis by ID
	r.HandleFunc(
		"/api/entity-classification/results/{analysis_id}",
		handlers.GetResultHandler,
	).Methods(http.MethodGet)

	// List all results (optional but correct)
	r.HandleFunc(
		"/api/entity-classification/results",
		handlers.ListResultsHandler,
	).Methods(http.MethodGet)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	return r
}
