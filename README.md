# canopy (backend)

**The Go service behind the canopy support agent.** CRUD for Silvanet devices, a downlink command queue, API-key auth, and Postgres persistence.

This repo is deliberately boring and trustworthy: it is the **source of truth** the agent is allowed to believe. The agent (a separate Python repo) never guesses device state — it asks this service, and this service either returns verified data or an honest error. Nothing here knows or cares that an LLM is the caller.

---

## What it exposes

A small REST API over HTTP, guarded by an API key (`X-API-Key` header).

| Method | Path | Does | Auth |
|---|---|---|---|
| `GET` | `/devices/:id` | Fetch one device's current state | ✅ |
| `POST` | `/devices` | Create a device | ✅ |
| `POST` | `/devices/:id/downlinks` | Queue a control command for a device | ✅ |
| `GET` | `/healthz` | Liveness probe (no auth) | — |
| `GET` | `/swagger/*` | Interactive API docs | — |

The agent only uses `GET /devices/:id` (status) and `POST /devices/:id/downlinks` (the human-approved write). The rest exists so the system is operable and seedable.

---

## Design stance

Three rules this service never breaks — they're what make it safe to put an agent in front of:

1. **The server owns server-decided fields.** `id`, timestamps, and `status` are set by the service, never accepted from the request. A request DTO is not a response DTO. A caller cannot POST a device "pre-marked online."
2. **Commands are whitelisted.** A downlink `command` must be one of a fixed set (`reset`, `calibrate`, `set_interval`). Anything else is rejected at the service layer, before it touches the queue — so even a compromised caller can't inject arbitrary commands.
3. **A downlink is *queued*, not executed.** `POST .../downlinks` writes a row with `status = queued`. This service does not talk to hardware; it records intent. That boundary is what makes the write safe to expose.

---

## Tech stack

- **Go** + **Gin** (HTTP router)
- **PostgreSQL** via **sqlx**
- **golang-migrate** for schema migrations
- **swaggo** for Swagger/OpenAPI docs
- **Docker Compose** for local run (service + Postgres)

---

## Quick start

```bash
# 1. Secrets — never committed
cp .envrc.example .envrc     # fill DB_DSN and CANOPY_API_KEY
direnv allow

# 2. Schema + seed data
make migrate-up              # creates devices + downlinks tables, seeds demo devices

# 3. Run (service on :8080, Postgres alongside)
docker compose up --build
```

Then open **http://localhost:8080/swagger/index.html** and authorise with your `CANOPY_API_KEY` (the green **Authorize** button) to try the endpoints.

> If Swagger doesn't reflect your latest code after an edit, the image is stale — `docker compose up --build` (occasionally `--no-cache`). Regenerate the spec with `make swagger` before building.

---

## Project layout

```
canopy/
├── app/
│   ├── main.go
│   ├── internal/
│   │   ├── models/        # Device, Downlink, DeviceStatus, UpdateFields
│   │   ├── service/       # business rules: owns id/timestamps/status, command whitelist
│   │   └── repository/    # RepositoryInterface + sqlx implementation
│   └── infra/
│       └── http/
│           ├── handlers/  # Gin handlers
│           ├── dto/       # request DTOs (inbound) + response DTOs (outbound) + mappers
│           └── middleware # X-API-Key auth
├── migrations/            # golang-migrate: devices, seed, downlinks
├── docker-compose.yml
├── Makefile               # migrate-up / migrate-down / swagger / run
└── .envrc.example
```

The layering is intentional: **handlers** parse HTTP, **service** holds the rules, **repository** touches the DB behind an interface (`var _ RepositoryInterface = (*Repository)(nil)` enforces the contract at compile time). Request and response DTOs are separate types so inbound and outbound shapes can never silently share a field.

---

## Data model (the bits the agent reads)

A `Device` carries the raw facts the agent reasons over — note it reports `status` and `battery` and `last_seen`, but it does **not** report "healthy." That judgment is the agent's job, computed from these fields; this service just reports what the hardware last told it.

| Field | Meaning |
|---|---|
| `id` | Device identifier (e.g. `dev-001`) |
| `name` | Human label (e.g. "North Ridge Sensor") |
| `status` | What the hardware last reported (`online` / `inactive`) |
| `battery` | Percentage, 0–100 |
| `last_seen` | Timestamp of last report (ISO-8601, with offset) |

A `Downlink` records a queued command: `device_id`, `command`, `status` (`queued` / `failed`), and timestamps.

---

### 2. Set up the database

```bash
# start postgres (Homebrew)
brew services start postgresql@18

# create DB and app user (uses .envrc values)
bash app/scripts/setup_db.sh
```
---

## Configuration

All config is environment-driven (via `direnv`), nothing committed:

| Var | Purpose |
|---|---|
| `DB_DSN` | Postgres connection string |
| `CANOPY_API_KEY` | The key callers must send as `X-API-Key` |

`.envrc` / `.env` are gitignored; `.envrc.example` ships with placeholders.

---

## Makefile targets

| Target | Does |
|---|---|
| `make migrate-up` | Apply all migrations (tables + seed) |
| `make migrate-down` | Roll back |
| `make swagger` | Regenerate the OpenAPI spec from code annotations |
| `make run` | Run the service locally (without compose) |

---

### 5. Open Swagger UI

```
http://localhost:8080/swagger/v1/index.html
```
---

## How the agent uses this service

The agent repo (`agentic-ai-90day/phase4_project`) wraps every call to this API in a **gate**: a response is either verified data or nothing-with-a-reason (401 → `auth_failed`, 404 → `not_found`, bad JSON → `invalid_json`, missing field → `missing_field:x`). The agent's LLM is only ever allowed to describe a *verified* payload. That means the contract this service must honour is simple but strict: **correct status codes and a stable field shape.** A 404 that returns a 200 with an error body would let the agent narrate a missing device as a present one. The agent's safety leans on this service being honest about failure.

See the agent repo's README for the full system story (client → MCP server → **this backend** → Postgres), the gate pattern, and the human-in-the-loop write flow.

---

*Part of the canopy capstone. This service is the trustworthy edge; the intelligence lives in the agent, the safety lives in the boundary between them.*
