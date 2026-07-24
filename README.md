# CodeZone

Web chat MVP with React, Go, and MySQL.

## Development

1. Start MySQL:

   ```bash
   docker compose up -d mysql
   ```

2. Configure the backend (YAML) and optionally the frontend env:

   ```bash
   cp backend/config.local.yaml backend/config.yaml
   cp frontend/.env.example frontend/.env   # optional; only needed for non-dev builds
   ```

   The backend reads `backend/config.yaml` by default. For other environments
   point `CODEZONE_CONFIG` at the appropriate template:

   ```bash
   CODEZONE_CONFIG=config.test.yaml ./api   # test / staging
   CODEZONE_CONFIG=config.gray.yaml ./api   # gray / canary
   CODEZONE_CONFIG=config.prod.yaml ./api   # production
   ```

   `config.local.yaml` / `config.test.yaml` / `config.gray.yaml` /
   `config.prod.yaml` are versioned templates; `config.yaml` is git-ignored
   so real secrets never enter the repo.

3. Run backend migrations after Task 2 is complete:

   ```bash
   cd backend
   go run ./cmd/migrate
   ```

4. Start the API:

   ```bash
   cd backend
   go run ./cmd/api
   ```

5. Start the frontend:

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
