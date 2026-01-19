# Local Setup Guide

## Prerequisites
- Go 1.21+
- Node.js 18+
- MySQL 5.7+

## Backend Setup

1. **Install dependencies:**
   ```bash
   cd backend
   go mod download
   ```

2. **Create MySQL database:**
   ```sql
   CREATE DATABASE url_shortener;
   ```

3. **Run migrations:**
   ```sql
   USE url_shortener;
   -- Run the SQL from backend/migrations/001_create_urls_table.sql
   ```

4. **Create .env file:**
   Copy `.env.example` to `.env` and update with your MySQL credentials:
   ```
   PORT=8080
   DB_DRIVER=mysql
   DB_DSN=root:password@tcp(localhost:3306)/url_shortener
   ```

5. **Run backend:**
   ```bash
   go run cmd/api/main.go
   ```
   Backend runs on http://localhost:8080

## Frontend Setup

1. **Install dependencies:**
   ```bash
   cd frontend
   npm install
   ```

2. **Create .env file:**
   Copy `.env.example` to `.env`

3. **Run frontend:**
   ```bash
   npm run dev
   ```
   Frontend runs on http://localhost:5173

## Testing

- Open http://localhost:5173 in your browser
- Enter a URL and click "Shorten"
- Copy the shortened URL and visit it - you'll be redirected to the original URL
- API endpoint: POST http://localhost:8080/api/v1/urls with body: `{"original_url": "https://example.com"}`
