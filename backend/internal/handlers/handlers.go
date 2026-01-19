package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Shiv-Kiran/url-qr-shortner/internal/database"
	"github.com/yourusername/url-qr-shortner/internal/services"
)

// HandleHome serves a simple home page
func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head><title>URL Shortener</title></head>
		<body>
		<h1>URL Shortener API</h1>
		<p>POST /api/v1/urls - Create shortened URL</p>
		<p>GET /api/v1/urls/:shortCode - Get URL info</p>
		<p>GET /:shortCode - Redirect to URL</p>
		</body>
		</html>
	`))
}

// CreateURLRequest represents the request body
type CreateURLRequest struct {
	OriginalURL string `json:"original_url"`
}

// CreateURLResponse represents the response
type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// HandleCreateURL handles POST /api/v1/urls
func HandleCreateURL(w http.ResponseWriter, r *http.Request) {
	var req CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	// Shorten URL
	urlModel, err := services.ShortenURL(req.OriginalURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateURLResponse{
		ShortCode:   urlModel.ShortCode,
		OriginalURL: urlModel.OriginalURL,
		ShortURL:    "http://localhost:8080/" + urlModel.ShortCode,
	})
}

// HandleGetURL handles GET /api/v1/urls/:shortCode
func HandleGetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	urlModel, err := services.GetURL(shortCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urlModel)
}

// HandleRedirect handles GET /:shortCode
func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	urlModel, err := services.GetURL(shortCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL not found"})
		return
	}

	// Increment clicks
	database.IncrementClicks(shortCode)

	// Redirect
	w.Header().Set("Location", urlModel.OriginalURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
