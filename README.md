# url-qr-shortner

A learning-focused URL shortener that now supports SQLite persistence and customizable QR code generation.

## Current Features

- Short URL creation with random short codes
- Redirect by short code
- Click counter persistence in SQLite
- In-memory cache for hot short codes
- QR code generation on URL creation
- QR customization:
  - Error correction level: `L`, `M`, `Q`, `H`
  - Size range: `128` to `1024` px

## API Endpoints

- `POST /api/v1/urls`
  - Body:
    ```json
    {
      "original_url": "https://example.com",
      "qr_error_correction": "M",
      "qr_size": 256
    }
    ```
- `GET /api/v1/urls/{shortCode}`
- `GET /api/v1/urls/{shortCode}/qr?level=Q&size=320`
- `GET /{shortCode}`

## Project Goal

This project is a sandbox for understanding routing, ID generation, persistence, analytics, and eventually advanced link-management features.

See `SETUP.md` for local setup and deployment notes.
