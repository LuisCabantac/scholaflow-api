# ScholaFlow API

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![Chi](https://img.shields.io/badge/Chi-v5.2-blue)](https://github.com/go-chi/chi)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org/)
[![pgx](https://img.shields.io/badge/pgx-v5.7-blue)](https://github.com/jackc/pgx)
[![sqlc](https://img.shields.io/badge/sqlc-v1.31-teal)](https://sqlc.dev/)
[![AWS SDK v2](https://img.shields.io/badge/AWS_SDK_v2-S3-FF9900?logo=amazon-aws)](https://aws.amazon.com/sdk-for-go/)
[![Go-Mail](https://img.shields.io/badge/Go--Mail-v0.8-green)](https://github.com/wneessen/go-mail)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A high-performance, modular REST API backend for the [ScholaFlow](https://github.com/LuisCabantac/scholaflow) LMS ecosystem. Built with Go 1.24, Chi router, PostgreSQL (`pgxpool`), sqlc type-safe query generation, AWS S3 / Supabase Storage, and transactional SMTP mail delivery.

---

## 1. Overview & Key Capabilities

ScholaFlow API serves as the core backend service responsible for academic domain logic, coursework lifecycle, assignment submissions, storage interactions, and notification processing. Engineered for low latency, memory efficiency, and serverless/container deployment flexibility, it pairs Go's concurrency primitives with sqlc compile-time SQL validation.

### Core Capabilities

- **Compile-Time Type-Safe SQL (`sqlc`):** Generates zero-allocation, type-safe Go query code from raw PostgreSQL schema migrations and queries with `pgx/v5` pool integration.
- **Connection Pooling & Supabase Compatibility:** Configured with `pgxpool` using simple protocol mode (`QueryExecModeSimpleProtocol`) for zero-fault execution behind transaction poolers (Supavisor / PgBouncer on port `6543`) as well as direct connections.
- **S3 Object Storage Integration:** Direct binary stream upload and presigned asset handling for classroom materials, stream attachments, submissions, and avatars via AWS SDK v2 over Supabase S3 storage.
- **SMTP Transactional Mailer:** Secure TLS/STARTTLS email dispatch powered by `go-mail` for classroom notifications, invitations, and alerts.
- **Robust HTTP & Error Architecture:** Standardized HTTP response encoding, strict size-limited JSON payload decoders, and unified error mapping via domain `apperror`.
- **Serverless & Edge Ready:** Pre-configured with custom build flags for deployment on Vercel Serverless Functions or standard long-running containers.

---

## 2. Architecture / How it Works

ScholaFlow API follows a modular clean-architecture layout. Inbound HTTP requests traverse middleware (request ID tracing, structured logging, panic recovery, CORS, and contextual timeouts) before reaching domain routing handlers. Handlers invoke services that leverage `sqlc` generated queries atop connection-pooled PostgreSQL transactions.

### System Architecture Flow

```mermaid
flowchart TD
    subgraph Clients ["Frontend & Client Ecosystem"]
        WebClient["ScholaFlow Web (Next.js / React 19)"]
        MobileClient["ScholaFlow Mobile (Expo / React Native)"]
        DesktopClient["ScholaFlow Desktop (Wails v3)"]
    end

    subgraph API ["ScholaFlow API (Go HTTP Server)"]
        ChiRouter["Chi Router & Middleware (cors, logger, timeout)"]
        
        subgraph Core ["Application Core (cmd/api)"]
            AppStruct["Application Container (api.go)"]
            Handlers["HTTP Route Handlers"]
        end

        subgraph Infrastructure ["Internal Packages (internal/*)"]
            AppError["apperror (Domain Errors)"]
            Response["response (JSON / Text Serializer)"]
            Request["request (JSON Decoder & File Parser)"]
            Env["env (Environment Configuration)"]
        end
    end

    subgraph ExternalServices ["External Infrastructure"]
        SupabaseS3["Supabase S3 Storage (AWS SDK v2)"]
        GmailSMTP["Gmail SMTP / Mail Transport (go-mail)"]
        SupabaseAuth["ScholaFlow Auth / JWKS Gateway"]
    end

    subgraph DataLayer ["Database & Persistence"]
        PGXPool["pgxpool (Simple Protocol Pooler)"]
        SQLC["sqlc Generated Queries (internal/adapters/postgresql/sqlc)"]
        Postgres[("PostgreSQL Database")]
    end

    WebClient -->|"REST API Requests"| ChiRouter
    MobileClient -->|"REST API Requests"| ChiRouter
    DesktopClient -->|"REST API Requests"| ChiRouter

    ChiRouter --> Handlers
    Handlers --> Infrastructure
    Handlers -->|"Query Execution"| SQLC
    SQLC --> PGXPool
    PGXPool --> Postgres
    
    Handlers -->|"File Storage & Retrieval"| SupabaseS3
    Handlers -->|"Dispatch Notifications"| GmailSMTP
    Handlers -.->|"Verify JWT / Auth"| SupabaseAuth
```

### Request Flow Lifecycle

1. **Routing & Middleware Pipeline:** Incoming requests are assigned a unique UUID request ID, logged via structured `log/slog`, and bounded by a 30-second context timeout.
2. **Input Decoding & Validation:** Request bodies are strictly decoded using `request.DecodeJSON` with a 1MB payload ceiling and unknown-field rejections, or parsed via `request.FormFile` for multipart binaries.
3. **Database Execution:** Handlers execute type-safe SQL operations via `queries.*`. For multi-step atomic operations, handlers acquire a transaction from `pool.Begin(ctx)` and bind it via `queries.WithTx(tx)`.
4. **Binary & Email Pipelines:** Uploads are streamed directly to Supabase S3 buckets without disk overhead; notifications trigger concurrent goroutine SMTP transmissions.
5. **Standardized Response:** Responses are marshaled to JSON through `response.JSON` or mapped through `response.Error` for type-safe error contracts.

---

## 3. Tech Stack

### Core Framework & Runtime

- **Language:** [Go 1.24+](https://golang.org/)
- **Router:** [go-chi/chi/v5](https://github.com/go-chi/chi) (Lightweight, composable HTTP router)
- **Middleware:** `chi/middleware`, `go-chi/cors`
- **Structured Logging:** Standard library `log/slog`

### Database & Persistence

- **Database:** [PostgreSQL 16](https://www.postgresql.org/) (Compatible with Supabase, AWS RDS, Neon)
- **Driver & Pool:** [pgx/v5](https://github.com/jackc/pgx) (`github.com/jackc/pgx/v5/pgxpool`)
- **Query Generator:** [sqlc v1.31+](https://sqlc.dev/) (Compiles SQL to type-safe Go)
- **Schema Migrations:** [Goose](https://github.com/pressly/goose) (`internal/adapters/postgresql/migrations`)

### Cloud, Storage & Transport

- **Object Storage:** [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) (`service/s3`) with Supabase S3 endpoint
- **Mail Service:** [go-mail](https://github.com/wneessen/go-mail) (`github.com/wneessen/go-mail`) with STARTTLS support
- **Environment Management:** [godotenv](https://github.com/joho/godotenv) & [go-env](https://github.com/allisson/go-env)

---

## 4. Project Structure

```text
scholaflow-api/
├── cmd/
│   └── api/
│       ├── api.go                              # Application struct, router mounting, and server runner
│       └── main.go                             # Entrypoint: client initialization (S3, Mail, DB pool)
├── internal/
│   ├── adapters/
│   │   └── postgresql/
│   │       ├── migrations/                     # Goose SQL migration files
│   │       │   ├── 00001_user.sql
│   │       │   ├── 00002_session.sql
│   │       │   ├── 00003_account.sql
│   │       │   ├── 00004_verification.sql
│   │       │   ├── 00005_jwks.sql
│   │       │   ├── 00006_notification.sql
│   │       │   ├── 00007_note.sql
│   │       │   ├── 00008_classroom.sql
│   │       │   ├── 00009_classroom_enrollment.sql
│   │       │   ├── 00010_classroom_message.sql
│   │       │   ├── 00011_classroom_topic.sql
│   │       │   ├── 00012_stream.sql
│   │       │   ├── 00013_submission.sql
│   │       │   ├── 00014_stream_comment.sql
│   │       │   └── 00015_submission_comment.sql
│   │       └── sqlc/                           # sqlc generated Go packages
│   │           ├── queries/                    # Input SQL query definitions (*.sql)
│   │           ├── db.go                       # Generated DBTX interface & query constructor
│   │           └── models.go                   # Generated Go model structs
│   ├── apperror/
│   │   └── apperror.go                         # Centralized domain error types & constructors
│   ├── env/
│   │   └── env.go                              # Environment variable helpers & fallbacks
│   ├── request/
│   │   └── request.go                          # Size-limited JSON decoding & multipart parsing
│   └── response/
│       └── response.go                         # Standard JSON/text writers & error responders
├── go.mod                                      # Go module definition
├── go.sum                                      # Go module checksums
├── sqlc.yaml                                   # sqlc code generation configuration
├── vercel.json                                 # Vercel deployment & build flag configuration
└── LICENSE                                     # MIT License
```

---

## 5. Getting Started

### Prerequisites

Ensure you have the following installed on your development machine:

- **Go:** `1.24.x` or higher
- **PostgreSQL Database:** Running instance (or a managed [Supabase](https://supabase.com/) project)
- **Goose CLI:** For database migrations (`go install github.com/pressly/goose/v3/cmd/goose@latest`)
- **sqlc CLI:** For query code generation (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- **Google App Password:** For Gmail SMTP transactional mail delivery

---

### Step-by-Step Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/LuisCabantac/scholaflow-api.git
cd scholaflow-api
```

#### 2. Install Go Dependencies

```bash
go mod download
```

#### 3. Configure Environment Variables

Create a `.env` file in the project root:

```bash
touch .env
```

Populate `.env` with your credentials:

| Variable                | Required | Description                                                    | Example / Default                                                         |
| :---------------------- | :------: | :------------------------------------------------------------- | :------------------------------------------------------------------------ |
| `PORT`                  | Optional | Port for the HTTP server                                       | `8080`                                                                    |
| `APP_ENV`               | Optional | Runtime environment (`development` / `production`)             | `development`                                                             |
| `APP_URL`               | **Yes**  | Canonical URL of your primary frontend                         | `http://localhost:3000`                                                   |
| `GOOSE_DBSTRING`        | **Yes**  | PostgreSQL DSN (Supabase pooler on port 6543 or direct on 5432)| `postgresql://postgres.[ref]:[pass]@aws-0-[region].pooler.supabase.com:6543/postgres` |
| `AWS_ENDPOINT_URL_S3`   | **Yes**  | Supabase Storage S3 endpoint URL                               | `https://[project-id].supabase.co/storage/v1/s3`                          |
| `AWS_ACCESS_KEY_ID`     | **Yes**  | Supabase Storage S3 Access Key ID                              | `your-s3-access-key-id`                                                   |
| `AWS_SECRET_ACCESS_KEY` | **Yes**  | Supabase Storage S3 Secret Access Key                          | `your-s3-secret-access-key`                                               |
| `AWS_REGION`            | **Yes**  | S3 Region                                                      | `ap-northeast-1` / `us-east-1`                                            |
| `SMTP_EMAIL`            | **Yes**  | Sender email address for SMTP                                  | `your-app@gmail.com`                                                      |
| `SMTP_PASSWORD`         | **Yes**  | 16-character Google Account App Password                       | `xxxx xxxx xxxx xxxx`                                                     |

#### 4. Run Database Migrations

Apply the migration schema to your PostgreSQL database using Goose (use direct connection port `5432` for migrations):

```bash
goose -dir internal/adapters/postgresql/migrations postgres "$GOOSE_DBSTRING" up
```

#### 5. Generate sqlc Code

If you modify or create new SQL queries in `internal/adapters/postgresql/sqlc/queries`:

```bash
sqlc generate
```

#### 6. Launch the Development Server

```bash
go run ./cmd/api
```

The server will start listening at [http://localhost:8080](http://localhost:8080).

#### 7. Verify Healthcheck

```bash
# Root info endpoint
curl http://localhost:8080/
# Response: {"message":"ScholaFlow API","version":"1.0.0","status":"running","environment":"development"}

# Health check endpoint
curl http://localhost:8080/healthz
# Response: OK
```

---

## 6. Development & Usage Reference

### 1. Adding New SQL Queries

1. Add your SQL query file under `internal/adapters/postgresql/sqlc/queries/`:
   ```sql
   -- name: GetClassroomByID :one
   SELECT * FROM classroom
   WHERE id = $1 LIMIT 1;
   ```
2. Run code generation:
   ```bash
   sqlc generate
   ```
3. Use the generated method directly on `app.queries`:
   ```go
   classroom, err := app.queries.GetClassroomByID(ctx, classroomID)
   ```

### 2. Transaction Management

When performing multi-table atomic updates:

```go
tx, err := app.db.Begin(ctx)
if err != nil {
    return response.Error(w, apperror.Internal(err))
}
defer tx.Rollback(ctx)

qtx := app.queries.WithTx(tx)

// Run queries on qtx...
if err := qtx.CreateEnrollment(ctx, params); err != nil {
    return response.Error(w, apperror.Internal(err))
}

if err := tx.Commit(ctx); err != nil {
    return response.Error(w, apperror.Internal(err))
}
```

---

## 7. Troubleshooting / Error Handling

### 1. Prepared Statement Error on Supabase (`SQLSTATE 26000`)

- **Symptom:** `ERROR: prepared statement "..." does not exist (SQLSTATE 26000)`.
- **Resolution:** You are connected to Supabase's Transaction Pooler (port `6543`). Verify `poolCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol` is set in [cmd/api/main.go](file:///home/luis/dev/projects/SCHOLAFLOW/scholaflow-api/cmd/api/main.go).

### 2. Goose Migrations Hanging or Failing on Port 6543

- **Symptom:** Migration commands stall, hang, or error during advisory locks.
- **Resolution:** Goose requires a direct PostgreSQL connection or session pooler (port `5432`) to acquire transaction advisory locks. Do not run schema migrations over port `6543`.

### 3. SMTP Authentication Failure

- **Symptom:** `535-5.7.8 Username and Password not accepted`.
- **Resolution:** Ensure 2-Step Verification is active on the sender Google Account and `SMTP_PASSWORD` is a 16-character **Google App Password**, not your primary Gmail password.

### 4. CORS Blocked During Frontend Development

- **Symptom:** `Cross-Origin Request Blocked: The Same Origin Policy disallows reading the remote resource`.
- **Resolution:** Ensure `APP_ENV=development` is set in your `.env` so localhost dev origins (`localhost:3000`, `localhost:5173`) are permitted by the CORS middleware.

---

## License

This project is licensed under the [MIT License](LICENSE).
