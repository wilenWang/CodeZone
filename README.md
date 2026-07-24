# CodeZone

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
   ENABLE_DEV_LOGIN=true go run ./cmd/api
   ```

4. Start the frontend:

   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Verification

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm run test
npm run build
```

End-to-end:

```bash
docker compose up -d mysql
cd backend && go run ./cmd/migrate
cd backend && ENABLE_DEV_LOGIN=true go run ./cmd/api
cd frontend && npm run dev
cd frontend && npm run test:e2e
```
