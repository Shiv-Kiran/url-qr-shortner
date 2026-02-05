# Local Setup Guide

## Prerequisites
- Go 1.21+
- Node.js 18+

## Backend Setup

1. Install dependencies:
   ```bash
   cd backend
   go mod download
   ```

2. Create `.env` from example:
   ```bash
   cp .env.example .env
   ```

3. Default backend environment values:
   ```env
   PORT=8080
   DB_DRIVER=sqlite3
   DB_DSN=url_shortener.db
   BASE_URL=http://localhost:8080
   ```

4. Run backend:
   ```bash
   go run cmd/api/main.go
   ```

Backend runs on `http://localhost:8080`.

## Frontend Setup

1. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```

2. Create `.env` from example:
   ```bash
   cp .env.example .env
   ```

3. Run frontend:
   ```bash
   npm run dev
   ```

Frontend runs on `http://localhost:5173`.

## API Quick Test

Create short URL with QR customization:

```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "original_url": "https://example.com",
    "qr_error_correction": "H",
    "qr_size": 320
  }'
```

Fetch a QR code for an existing short code:

```bash
curl "http://localhost:8080/api/v1/urls/<shortCode>/qr?level=Q&size=256"
```

## Deployment Notes

- Deploy backend first (Render, Railway, Fly.io, etc.) and keep `BASE_URL` set to your backend public URL.
- Deploy frontend to Vercel.
- In Vercel project env vars, set:
  - `VITE_API_URL=https://<your-backend-domain>/api/v1`
