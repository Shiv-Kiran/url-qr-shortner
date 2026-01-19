package handlers

import (
	"net/http"
)

// HandleCreateURL handles POST /api/v1/urls
func HandleCreateURL(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement URL shortening
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Create URL endpoint"}`))
}

// HandleRedirect handles GET /api/v1/urls/:shortCode
func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement redirect logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Redirect endpoint"}`))
}
