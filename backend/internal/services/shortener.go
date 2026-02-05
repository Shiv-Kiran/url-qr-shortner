package services

import (
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/Shiv-Kiran/url-qr-shortner/internal/database"
	"github.com/Shiv-Kiran/url-qr-shortner/internal/models"
)

const (
	charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLen = 6
)

var (
	cache = make(map[string]*models.URL)
	mu    sync.RWMutex
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// InitCache initializes the cache
func InitCache() {
	mu.Lock()
	defer mu.Unlock()
	cache = make(map[string]*models.URL)
}

// ShortenURL generates a shortened URL
func ShortenURL(originalURL string) (*models.URL, error) {
	if !ValidateURL(originalURL) {
		return nil, fmt.Errorf("invalid URL format")
	}

	shortCode, err := generateUniqueCode()
	if err != nil {
		return nil, err
	}

	urlModel := &models.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
	}

	// Save to database
	if err := database.SaveURL(urlModel); err != nil {
		return nil, err
	}

	// Cache it
	mu.Lock()
	cache[shortCode] = urlModel
	mu.Unlock()

	return urlModel, nil
}

// GetURL retrieves URL from cache or database
func GetURL(shortCode string) (*models.URL, error) {
	// Check cache first
	mu.RLock()
	if urlModel, exists := cache[shortCode]; exists {
		mu.RUnlock()
		return urlModel, nil
	}
	mu.RUnlock()

	// Fetch from database
	urlModel, err := database.GetURLByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	// Update cache
	mu.Lock()
	cache[shortCode] = urlModel
	mu.Unlock()

	return urlModel, nil
}

// ValidateURL checks if URL is valid
func ValidateURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}
	u, err := url.Parse(urlStr)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

// generateUniqueCode generates a unique short code
func generateUniqueCode() (string, error) {
	for attempts := 0; attempts < 10; attempts++ {
		code := generateRandomCode()
		exists, err := database.URLExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique code")
}

// generateRandomCode generates a random code
func generateRandomCode() string {
	b := make([]byte, codeLen)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}
