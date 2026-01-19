package models

import "time"

// URL represents a shortened URL record
type URL struct {
	ID        int64     `json:"id"`
	OriginalURL string   `json:"original_url"`
	ShortCode  string    `json:"short_code"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Clicks     int64     `json:"clicks"`
}
