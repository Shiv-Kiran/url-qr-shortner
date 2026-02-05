package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Shiv-Kiran/url-qr-shortner/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dsn string) error {
	if dsn == "" {
		dsn = "url_shortener.db"
	}

	var err error
	db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// createTables creates necessary tables
func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		created_at TEXT NOT NULL,
		expires_at TEXT,
		clicks INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	CREATE INDEX IF NOT EXISTS idx_created_at ON urls(created_at);
	`

	_, err := db.Exec(schema)
	return err
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
	if url.CreatedAt.IsZero() {
		url.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO urls (original_url, short_code, created_at) VALUES (?, ?, ?)`
	result, err := db.Exec(query, url.OriginalURL, url.ShortCode, url.CreatedAt.UTC().Format(time.RFC3339Nano))
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
	var createdAt string
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &createdAt, &url.Clicks)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("URL not found")
		}
		return nil, err
	}

	url.CreatedAt, err = parseStoredTime(createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at for short code %q: %w", shortCode, err)
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

func parseStoredTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported datetime format %q", value)
}
