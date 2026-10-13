# uptime

A self-hostable uptime monitor: add URLs, the Go backend checks them on a schedule, and you get response-time graphs, incident history, and public status pages.

**Stack:** Go (Golang), Postgres, React + TypeScript (Vite, TanStack Router, TanStack Query, Zustand, Tailwind, shadcn/ui).

```text
api-docs/     OpenAPI spec for the backend API
backend/      Go module: API entry point (cmd/api), feature packages (internal/), SQL migrations
frontend/     React + TypeScript application powered by Vite
compose.yaml  Postgres and the backend for Docker
```

## Development

### Run with Docker

Requires only Docker:

```powershell
docker compose up -d --build
```

Open http://localhost:3000. nginx serves the frontend and forwards `/api` to the backend.

You can also start a subset of the services. Dependencies start automatically, so `backend` also starts `postgres`:

```powershell
docker compose up -d postgres          # database only (for local development)
docker compose up -d --build backend   # database + API on port 8080
```

Stop everything with `docker compose down`. Add `-v` to also delete the database volume.

A demo account is created by the migrations: `demo@uptime.io` / `demo1234`.

### Run locally

Requirements: Go (see `backend/go.mod`), Node.js 22+, Docker.

1. Start Postgres:

   ```powershell
   docker compose up -d postgres
   ```

2. Run the backend. It applies pending migrations on startup.

   ```powershell
   cd backend
   go run ./cmd/api
   ```

3. Run the frontend:

   ```powershell
   cd frontend
   npm install
   npm run dev
   ```

Open http://localhost:8081. The Go API runs on port 8080, and Vite forwards `/api/*` to it during development.

### Configuration

| Variable       | Default                                                            |
|----------------|--------------------------------------------------------------------|
| `DATABASE_URL` | `postgres://admin:password@localhost:5432/uptime?sslmode=disable` |
| `HTTP_ADDR`    | `:8080`                                                            |
| `CHECKER_WORKERS` | `10` (number of monitors checked in parallel)                   |

The API is documented in [`api-docs/openapi.yaml`](api-docs/openapi.yaml).

### Migrations

Migrations live in `backend/migrations` as [goose](https://github.com/pressly/goose) SQL files and are embedded in the binary. Add a new one as `000NN_description.sql` with `-- +goose Up` / `-- +goose Down` sections.

### Queries

SQL queries live in `backend/internal/store/queries` and are compiled to type-safe Go in `backend/internal/store/db` by [sqlc](https://sqlc.dev). Never edit the generated code; after changing a query or migration, regenerate it from `backend`:

```powershell
docker run --rm -v "${PWD}:/src" -w /src sqlc/sqlc:1.31.1 generate
```

CI fails if the generated code is out of date.

### API client

The frontend calls the backend through a TypeScript client generated from `api-docs/openapi.yaml` by [Hey API](https://heyapi.dev) into `frontend/src/client`. The client isn't committed: `npm run dev` and `npm run build` generate it. After changing the spec while the dev server is running, regenerate it from `frontend`:

```powershell
npm run generate:api
```

## Checks

In `backend`: `go vet ./...` and `go test ./...`.

The Postgres integration tests (e.g. for the checker's `SKIP LOCKED` queue) run when `TEST_DATABASE_URL` points at a disposable database whose name ends in `_test`; they empty it first. With Docker Compose:

```powershell
docker compose exec postgres createdb -U admin uptime_test
$env:TEST_DATABASE_URL = "postgres://admin:password@localhost:5432/uptime_test?sslmode=disable"
go test ./...
```

In `frontend`: `npm run lint`, `npm test` (Vitest; `npm run test:watch` while developing) and `npm run build`.

CI runs the same checks on every push and pull request.