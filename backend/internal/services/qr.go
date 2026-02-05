package services

import (
	"encoding/base64"
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

const (
	defaultQRCodeSize = 256
	minQRCodeSize     = 128
	maxQRCodeSize     = 1024
)

var recoveryLevelByName = map[string]qrcode.RecoveryLevel{
	"L": qrcode.Low,
	"M": qrcode.Medium,
	"Q": qrcode.High,
	"H": qrcode.Highest,
}

// GenerateQRCodeDataURL creates a base64 data URL for a QR image.
// It also returns the normalized QR level and size that were applied.
func GenerateQRCodeDataURL(content string, level string, size int) (string, string, int, error) {
	recoveryLevel, normalizedLevel, normalizedSize, err := normalizeQROptions(level, size)
	if err != nil {
		return "", "", 0, err
	}

	pngBytes, err := qrcode.Encode(content, recoveryLevel, normalizedSize)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed generating qr code: %w", err)
	}

	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
	return dataURL, normalizedLevel, normalizedSize, nil
}

func normalizeQROptions(level string, size int) (qrcode.RecoveryLevel, string, int, error) {
	normalizedLevel := strings.ToUpper(strings.TrimSpace(level))
	if normalizedLevel == "" {
		normalizedLevel = "M"
	}

	recoveryLevel, ok := recoveryLevelByName[normalizedLevel]
	if !ok {
		return qrcode.Medium, "", 0, fmt.Errorf("invalid qr error correction level %q (use L, M, Q, or H)", level)
	}

	normalizedSize := size
	if normalizedSize == 0 {
		normalizedSize = defaultQRCodeSize
	}
	if normalizedSize < minQRCodeSize || normalizedSize > maxQRCodeSize {
		return qrcode.Medium, "", 0, fmt.Errorf("invalid qr size %d (use %d-%d)", normalizedSize, minQRCodeSize, maxQRCodeSize)
	}

	return recoveryLevel, normalizedLevel, normalizedSize, nil
}
