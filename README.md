# Vibework Chat

Web chat MVP with React, Go, and MySQL.

## Development

1. Start MySQL:

   ```bash
   docker compose up -d mysql
   ```

2. Run backend migrations after Task 2 is complete:

   ```bash
   cd backend
   go run ./cmd/migrate
   ```

3. Start the API:

   ```bash
   cd backend
   go run ./cmd/api
   ```

4. Start the frontend:

   ```bash
   cd frontend
   npm install
   npm run dev
   ```
