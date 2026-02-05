package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Shiv-Kiran/url-qr-shortner/internal/database"
	"github.com/Shiv-Kiran/url-qr-shortner/internal/services"
	"github.com/gorilla/mux"
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
		<p>GET /api/v1/urls/:shortCode/qr - Get QR code for short URL</p>
		<p>GET /:shortCode - Redirect to URL</p>
		</body>
		</html>
	`))
}

// CreateURLRequest represents the request body
type CreateURLRequest struct {
	OriginalURL       string `json:"original_url"`
	QRErrorCorrection string `json:"qr_error_correction,omitempty"`
	QRSize            int    `json:"qr_size,omitempty"`
}

// CreateURLResponse represents the response
type CreateURLResponse struct {
	ShortCode         string `json:"short_code"`
	OriginalURL       string `json:"original_url"`
	ShortURL          string `json:"short_url"`
	QRDataURL         string `json:"qr_data_url"`
	QRErrorCorrection string `json:"qr_error_correction"`
	QRSize            int    `json:"qr_size"`
}

// QRCodeResponse represents QR code response body
type QRCodeResponse struct {
	ShortCode         string `json:"short_code"`
	ShortURL          string `json:"short_url"`
	QRDataURL         string `json:"qr_data_url"`
	QRErrorCorrection string `json:"qr_error_correction"`
	QRSize            int    `json:"qr_size"`
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

	baseURL := getBaseURL(r)
	shortURL := baseURL + "/" + urlModel.ShortCode
	qrDataURL, qrLevel, qrSize, err := services.GenerateQRCodeDataURL(shortURL, req.QRErrorCorrection, req.QRSize)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateURLResponse{
		ShortCode:         urlModel.ShortCode,
		OriginalURL:       urlModel.OriginalURL,
		ShortURL:          shortURL,
		QRDataURL:         qrDataURL,
		QRErrorCorrection: qrLevel,
		QRSize:            qrSize,
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

// HandleGetQRCode handles GET /api/v1/urls/:shortCode/qr
func HandleGetQRCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if _, err := services.GetURL(shortCode); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL not found"})
		return
	}

	qrLevel := r.URL.Query().Get("level")
	qrSize, err := parseOptionalSize(r.URL.Query().Get("size"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	shortURL := getBaseURL(r) + "/" + shortCode
	qrDataURL, normalizedLevel, normalizedSize, err := services.GenerateQRCodeDataURL(shortURL, qrLevel, qrSize)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(QRCodeResponse{
		ShortCode:         shortCode,
		ShortURL:          shortURL,
		QRDataURL:         qrDataURL,
		QRErrorCorrection: normalizedLevel,
		QRSize:            normalizedSize,
	})
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
	_ = database.IncrementClicks(shortCode)

	// Redirect
	w.Header().Set("Location", urlModel.OriginalURL)
	w.WriteHeader(http.StatusMovedPermanently)
}

func getBaseURL(r *http.Request) string {
	if configured := strings.TrimSpace(os.Getenv("BASE_URL")); configured != "" {
		return strings.TrimRight(configured, "/")
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); forwardedProto != "" {
		scheme = strings.Split(forwardedProto, ",")[0]
	}

	host := strings.TrimSpace(r.Host)
	if host == "" {
		host = "localhost:8080"
	}

	return scheme + "://" + host
}

func parseOptionalSize(raw string) (int, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, nil
	}

	size, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid qr size %q", raw)
	}

	return size, nil
}
