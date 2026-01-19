package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yourusername/url-qr-shortner/internal/models"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dsn string) error {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// SaveURL saves a shortened URL to the database
func SaveURL(url *models.URL) error {
	query := `INSERT INTO urls (original_url, short_code) VALUES (?, ?)`
	result, err := db.Exec(query, url.OriginalURL, url.ShortCode)
	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	url.ID = id
	return nil
}

// GetURLByShortCode retrieves a URL by its short code
func GetURLByShortCode(shortCode string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at, clicks FROM urls WHERE short_code = ?`
	row := db.QueryRow(query, shortCode)

	url := &models.URL{}
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CreatedAt, &url.Clicks)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("URL not found")
		}
		return nil, err
	}

	return url, nil
}

// IncrementClicks increments the click count for a URL
func IncrementClicks(shortCode string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?`
	_, err := db.Exec(query, shortCode)
	return err
}

// URLExists checks if a short code already exists
func URLExists(shortCode string) (bool, error) {
	query := `SELECT COUNT(*) FROM urls WHERE short_code = ?`
	var count int
	err := db.QueryRow(query, shortCode).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
