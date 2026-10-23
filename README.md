# Under100 (Bite100) — Bengaluru Hyperlocal Food Discovery Under ₹100

> **Bengaluru's strictly curated food discovery platform where every single dish is priced at or below ₹100.** Discover legendary heritage darshinis, iconic breakfast joints, street food hubs, and hidden gems across Bengaluru without breaking the bank.

---

## 1. Project Overview & Use Cases

**Under100** (branded consumer-facing as **Bite100**) was engineered specifically for the culinary landscape of Bengaluru, India. It solves a real-world dilemma: popular food delivery platforms prioritize high average order values, burying authentic ₹30–₹90 tiffin meals and street delicacies under surge pricing, delivery fees, and platform surcharges.

### Target Personas & Use Cases:
1. **Budget-Conscious Diners & Foodies (Customer Persona):**
   - Browse nearby restaurants based on GPS geolocation, auto-calculated distances, and travel times across 11 Bengaluru neighborhoods (Indiranagar, Koramangala, Malleshwaram, Jayanagar, Basavanagudi, etc.).
   - Filter strictly dishes $\le$ ₹100 across 12 canonical cuisines, meal types, and food types.
   - Real-time **"Open Now"** schedule filtering with timezone precision (`Asia/Kolkata`).
   - Pure Vegetarian and Non-Vegetarian diet toggles.
   - Customer favorites bookmarking and 1-to-5 star restaurant reviews with rating breakdown.
   - "Order for Someone Else" delivery address coordination.
   - In-place cart counter steppers and slide-out e-commerce drawer with zero-hidden-fee guarantee.

2. **Local Food Outlets & Darshini Merchants (Merchant Persona):**
   - Self-service restaurant registration submitted to admin review queue.
   - Menu management: add, edit, and toggle availability of dishes strictly $\le$ ₹100 (backend price validation blocks anything > ₹100).
   - 7-day operating schedule editor with multi-period shifts and temporary closure scheduler with IDOR protection.

3. **Platform Administrators (Admin Persona):**
   - Restaurant approval/rejection moderation queue.
   - Food category lifecycle management (create, update, activate/deactivate without data loss).
   - Customer review moderation (flagged review triage, hide/restore toggles).

---

## 2. Technology Stack

- **Backend Runtime:** Go 1.24
- **Web Framework:** Echo v4 (High-performance, minimalist Go web framework)
- **Database:** PostgreSQL 15+ with Ent ORM (`entgo.io/ent`)
- **Cache & Session Store:** Redis 5+ (Port `6379`)
- **Frontend Engine:** Server-Side Rendered HTML5 + HTMX 1.9.10
- **Styling:** Vanilla CSS3 + Design System Tokens (No bloated client-side JS runtime dependencies)
- **Interactive UI:** Vanilla JavaScript + Client-side Haversine Distance Engine
- **Asset Integrity:** 19 curated, SHA256-verified static images and SVGs

---

## 3. Directory Structure

```text
C:\Users\astha\Desktop\food-discovery
├── cmd/
│   ├── web/                     # Web application entry point (main.go)
│   ├── seed/                    # Database seeder (main.go)
│   └── worker/                  # Background worker queue (main.go)
├── config/                      # Configuration YAML and Go loader (config.go)
├── ent/                         # Ent ORM schemas, generated entities, and migrations
├── pkg/
│   ├── repos/fooddiscovery/     # Availability engine, category repo, discovery queries
│   ├── routing/fooddiscovery/   # HTTP handlers for customer, merchant, and admin
│   ├── routing/routes/          # Route registration and middleware mounting
│   ├── services/                # Dependency injection container, database, cache, auth
│   └── middleware/              # Auth, CSRF, rate-limiting, and error middleware
├── static/                      # Static CSS bundles, icons, and 19 dish photographs
├── templates/                   # Go HTML templates (pages, layouts, HTMX components)
├── verify_final_handover.ps1    # Master 13-workflow regression test suite
└── under100_server.exe          # Compiled production Go web server binary
```

---

## 4. Prerequisites

Before running the application locally, ensure the following software is installed on Windows:
1. **Go 1.24+** (`go version`)
2. **PostgreSQL 15+** (Installed locally e.g. at `C:\Users\astha\pgsql` or standard PostgreSQL service)
3. **Redis 5+** (Installed locally e.g. at `C:\Users\astha\redis-server.exe` or Windows service)
4. **PowerShell 5.1+** or modern terminal
5. **curl.exe** (Built into Windows 10/11)

---

## 5. Rebuilding & Running from this Git Repository

### 5.1 Clone & Verify Dependencies
```bash
git clone <repo-url> food-discovery
cd food-discovery
go mod download
go mod verify
```

### 5.2 Start Dependencies (PostgreSQL 16 & Redis)
You can run the required database and cache services using Docker Compose or standalone services:

**Option A — Using Docker Compose (Recommended):**
```bash
docker compose up -d
```
This spins up PostgreSQL 16 on port `5432` (`goship_db`), Redis on port `6379`, and Mailpit on port `8025`.

**Option B — Using Local Standalone Services (Windows):**
```powershell
# Start PostgreSQL
& "C:\Users\astha\pgsql\bin\pg_ctl.exe" start -D "C:\Users\astha\pgdata" -l "C:\Users\astha\pgsql\logfile.log"

# Start Redis
Start-Process -FilePath "C:\Users\astha\redis-server.exe" -WindowStyle Hidden
```

### 5.3 Initialize Database Schema & Seed Data
Run the seeder once to apply Ent ORM schema migrations and seed all restaurants, menus, and accounts:
```bash
go run ./cmd/seed
```

### 5.4 Build the Application Targets
All 3 binaries compile cleanly directly from the repository source code:

**1. Web Application Server (`cmd/web`):**
```bash
# Windows
go build -o under100_server.exe ./cmd/web

# Linux / macOS
go build -o under100_server ./cmd/web
```

**2. Asynchronous Background Worker (`cmd/worker`):**
```bash
# Windows
go build -o under100_worker.exe ./cmd/worker

# Linux / macOS
go build -o under100_worker ./cmd/worker
```

**3. Linux Deployment Cross-Compilation (from Windows/macOS):**
```powershell
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o ./build_tmp/under100_web_linux ./cmd/web
go build -ldflags="-s -w" -o ./build_tmp/under100_worker_linux ./cmd/worker
go build -ldflags="-s -w" -o ./build_tmp/under100_seed_linux ./cmd/seed
```

### 5.5 Launch & Verify
Start the compiled web server:
```bash
# Windows
.\under100_server.exe

# Linux
./under100_server
```

Open your browser at `http://127.0.0.1:8000`.

---

## 6. Default Credentials

The development database is pre-configured with three operational personas:

| Role | Email | Password | Access Level |
| :--- | :--- | :--- | :--- |
| **Admin** | `admin@fooddiscovery.com` | `Admin@123456` | Full platform admin (`/admin`) |
| **Merchant** | `merchant.ctr@fooddiscovery.com` | `Merchant@123456` | Owns Central Tiffin Room (`/merchant`) |
| **Customer** | `user@fooddiscovery.com` | `User@123456` | Customer discovery, favorites, reviews |

---

## 7. Key Application URLs

- **Storefront / Homepage:** `http://127.0.0.1:8000/`
- **Nearby Restaurant Discovery:** `http://127.0.0.1:8000/nearby`
- **"Open Now" Filtered Discovery:** `http://127.0.0.1:8000/nearby?open_now=true`
- **Category Filter (e.g. South Indian):** `http://127.0.0.1:8000/nearby?category=south-indian`
- **Search Dishes / Outlets:** `http://127.0.0.1:8000/nearby?q=Dosa`
- **Customer Favorites:** `http://127.0.0.1:8000/favorites`
- **Merchant Owner Portal:** `http://127.0.0.1:8000/merchant`
- **Merchant Weekly Operating Hours:** `http://127.0.0.1:8000/merchant/hours`
- **Admin Dashboard:** `http://127.0.0.1:8000/admin`
- **Admin Category Management:** `http://127.0.0.1:8000/admin/categories`
- **Admin Review Moderation:** `http://127.0.0.1:8000/admin/reviews`

---

## 8. Running Tests & Automated Verification

### Full Regression Suite
Execute the comprehensive handover verification suite:
```powershell
powershell.exe -ExecutionPolicy Bypass -File .\verify_final_handover.ps1
```
This tests all 13 critical workflows (Auth, Discovery, ₹100 Price Cap, Categories, Favorites, Reviews, Merchant Onboarding, Admin Approval, IDOR Protection, Open Now Hours, 19 Static Asset Hashes, and Database Count Preservation).

### Unit Tests
Run Go availability service and category repository unit tests:
```powershell
go test -v ./pkg/repos/fooddiscovery -run TestAvailabilityService
go test -v ./pkg/repos/fooddiscovery -run "TestSlugify|TestMatchCategoryKeywords"
```

### Code Quality Audit
```powershell
go vet ./...
```

---

## 9. Troubleshooting & Gotchas

1. **PostgreSQL "Connection Refused":**
   Ensure `postgres.exe` is running and accepting TCP connections on port 5432. Check `C:\Users\astha\pgsql\logfile.log` if recovery is in progress.
2. **Redis Session Error:**
   If authentication fails with session storage errors, verify `redis-server.exe` is running on `127.0.0.1:6379`.
3. **Hard Price Cap Enforcement:**
   Dishes priced > ₹100 are strictly rejected upon submission by merchants (HTTP 400 Bad Request) and filtered out of discovery queries.
4. **IDOR Protection:**
   Merchants can only modify dishes, operating hours, and profile details for restaurants they explicitly own. Cross-merchant modifications are rejected with HTTP 403 Forbidden.
