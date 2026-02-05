package services

import (
	"strings"
	"testing"
)

func TestGenerateQRCodeDataURLDefaults(t *testing.T) {
	dataURL, level, size, err := GenerateQRCodeDataURL("https://example.com", "", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if level != "M" {
		t.Fatalf("expected default level M, got %s", level)
	}
	if size != defaultQRCodeSize {
		t.Fatalf("expected default size %d, got %d", defaultQRCodeSize, size)
	}
	if !strings.HasPrefix(dataURL, "data:image/png;base64,") {
		t.Fatalf("expected png data URL prefix, got %s", dataURL)
	}
}

func TestGenerateQRCodeDataURLInvalidLevel(t *testing.T) {
	_, _, _, err := GenerateQRCodeDataURL("https://example.com", "X", 256)
	if err == nil {
		t.Fatal("expected error for invalid level, got nil")
	}
}

func TestGenerateQRCodeDataURLInvalidSize(t *testing.T) {
	_, _, _, err := GenerateQRCodeDataURL("https://example.com", "H", 64)
	if err == nil {
		t.Fatal("expected error for invalid size, got nil")
	}
}
