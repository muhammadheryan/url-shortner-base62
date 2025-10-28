# url-shortner-base62

Simple URL shortener (base62) — minimal service to create, update and resolve short URLs.

## Features

- Create short URL records with `user_id`, `short_code`, and `original_url`.
- Update existing URL records (short_code, original_url).
- Retrieve a single URL by `id` or `short_code`.
- SQL migration included (`db/migrations/20250827113104_init_database.sql`).

## Project structure (important files)

- `cmd/main.go` — application entrypoint.
- `model/url.go` — URL entity struct.
- `repository/url/url_repository.go` — repository with Create/Update/Get methods.
- `db/migrations/20250827113104_init_database.sql` — initial migration.
- `transport/http.go` — HTTP transport (routes/handlers).

## Prerequisites

- Go 1.20+ installed (verify with `go version`).
- A SQL database (MySQL/MariaDB recommended) and `go-sql-driver/mysql` or other driver configured in project.
- `golang-migrate` or other migration tooling (optional) to run SQL migrations.

## Setup & Run (PowerShell)

1. Fetch dependencies and tidy modules:

```powershell
go mod tidy
```

2. Set environment variables (example):

```powershell
$env:DATABASE_DSN = "user:password@tcp(localhost:3306)/dbname?parseTime=true"
$env:PORT = "8080"
```

3. Run database migrations (example using `mysql` client):

```powershell
# using mysql client
mysql -u user -p password -h 127.0.0.1 -P 3306 dbname < db/migrations/20250827113104_init_database.sql
```

4. Build and run the application:

```powershell
go build -o bin/app ./cmd
.
\bin\app.exe
```

Or run directly with `go run`:

```powershell
go run ./cmd/main.go
```

5. API endpoints (example — adapt based on `transport/http.go`):

- POST /urls — create URL
- PUT /urls/{id} — update URL
- GET /urls/{short_code} — resolve or get URL

## Tests

Run unit tests:

```powershell
go test ./... -v
```

## Notes

- Adjust DSN and driver in `cmd/main.go` as necessary.
- Migration `migrate:down` in the SQL file uses `DROP TABLE urls;` ensure table name matches (`url`).

---

