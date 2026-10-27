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

### Authentication & Security Unit Tests
Run the complete auth test suite covering phone OTP, simulated zero-cost provider, email verification tokens, rate limiting, and role-based access control:
```powershell
# 1. SMS Provider & Security Guards
go test -v ./pkg/services/sms

# 2. Token Hashing, Account Recovery & Masking
go test -v ./pkg/repos/authrepo

# 3. HTTP Handlers, Validation, & Access Control
go test -v ./pkg/routing/fooddiscovery -run "TestAuthUI_PagesRender|TestAuthorization_RoleAccessControls|TestRegistration_ValidationRules|TestPhoneOTP_Logic|TestEmailVerification_UnverifiedLoginRejection|TestAccountRecovery_EmailMasking"
```

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

## 9. Complete Authentication System (Zero-Cost-First Architecture)

Under100 implements an enterprise-grade, zero-cost-first authentication system designed for Bangalore diners and food outlet merchants:

### 9.1 Key Features & Capabilities
1. **Unified Dual-Mode Login (`/login`):**
   - **Option A (Phone Number + OTP):** Instant SMS login defaulting to India (+91) with 6-digit verification code, 60-second resend countdown, 5-minute expiry, and max 5 verification attempts.
   - **Option B (Email & Password):** Secure login requiring verified email, Argon2id/bcrypt hash verification, show/hide password toggle, and optional "Remember me" session longevity.
2. **Email Registration & Verification (`/register`, `/verify-email`):**
   - New accounts created with `verified = false`.
   - Single-use cryptographically secure 32-byte tokens hashed via SHA-256 with 30-minute expiration.
   - Rate-limited resend endpoint with 60-second cooldown (`/api/auth/resend-verification`).
   - Unverified accounts cannot log in until the link received via email is clicked.
3. **Forgot & Reset Password Workflow (`/forgot-password`, `/reset-password`):**
   - Anti-enumeration generic messaging protects user privacy.
   - Single-use reset tokens with trusted `APP_BASE_URL` link construction (immune to Host-header poisoning).
   - Instant session revocation on password update via password signature rotation (`pass_sig`).
4. **Account Recovery for Forgotten Emails (`/recover-email`):**
   - Customers who forget their registered email can verify their previously linked phone number via OTP.
   - Upon successful verification, securely returns a masked email (e.g. `a***v@e***e.com`).
5. **Zero-Cost Development Mode (₹0 Setup):**
   - **Simulated OTP Provider:** When running in local development (`PAGODA_APP_ENVIRONMENT != prod`), generated OTPs are printed to the console terminal (`[DEV SIMULATED OTP]`) without telecom gateway charges.
   - **Strict Production Refusal:** The simulated provider fails immediately if accidentally configured in production mode.
   - **Mock Mailer Fallback:** In development mode without a `RESEND_API_KEY`, verification links and reset tokens are safely logged to standard output.

### 9.2 Database Schema & Migrations
Database migrations are applied via PostgreSQL:
```powershell
# Apply auth migration (creates phone_otp_challenges, email_verification_tokens, and adds users.phone)
psql -U admin -d goship_db -f dbs/migrations/20261027_auth_system.sql
```
Rollback procedure (if required):
```sql
DROP TABLE IF EXISTS phone_otp_challenges;
DROP TABLE IF EXISTS email_verification_tokens;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### 9.3 Transactional Email Delivery (Resend API)
The platform uses [Resend](https://resend.com) as its official transactional email delivery provider for:
- **Registration Confirmation:** Branded welcome email confirming account creation.
- **Email Verification:** Cryptographically secure 32-byte tokens (SHA-256 hashed in PostgreSQL), 30-minute expiry, and single-use validation.
- **Forgot & Reset Password:** Anti-enumeration generic responses, single-use reset links, 30-minute expiry, and automatic session invalidation upon password update.
- **Account Security Alerts:** Immediate email alert sent whenever an account password is changed.
- **Resend Verification:** 60-second client-side countdown + server-side rate-limiting cooldown.

#### Development / Zero-Cost Mode:
If `RESEND_API_KEY` is not provided (or in local development), the system automatically uses a `MockMailClient` that logs all outgoing verification links and password reset URLs directly to the terminal, allowing 100% test coverage with zero provider cost.

#### Production Configuration:
Configure the following environment variables:
```bash
# Resend API Key (free tier: 3,000 emails/month free)
RESEND_API_KEY=re_123456789_abcdef

# Sender address (must match verified domain in Resend)
EMAIL_FROM="Food Under ₹100 <onboarding@resend.dev>"
# Or for your custom domain:
# EMAIL_FROM="Food Under ₹100 <noreply@foodunder100.in>"

# Application Base URL (for link generation in emails; no trailing slash)
APP_BASE_URL=https://your-food-app.onrender.com
```

#### Resend Custom Sending Domain & DNS Setup:
1. Log in to [resend.com/domains](https://resend.com/domains) and click **Add Domain** (e.g., `foodunder100.in` or subdomain `mail.foodunder100.in`).
2. Add the DNS records provided by Resend to your domain registrar (GoDaddy, Cloudflare, Namecheap):
   - **DKIM (CNAME or TXT):** Authenticates outgoing emails via cryptographic signature.
   - **SPF (TXT):** Authorizes Resend mail servers to send emails on behalf of your domain (`v=spf1 include:resend.com ~all`).
   - **DMARC (TXT):** Specifies email alignment policy (`_dmarc.yourdomain.com`, `v=DMARC1; p=none; rua=mailto:dmarc-reports@yourdomain.com`).
3. Click **Verify Domain** in Resend dashboard. Once verified (status turns green), set `EMAIL_FROM="Food Under ₹100 <notifications@yourdomain.com>"`.

---

## 10. Render Deployment Configuration

When deploying the Food Under ₹100 application to [Render](https://render.com):

### 10.1 Service Type & Build Settings:
- **Service Type:** Web Service
- **Environment:** Go / Docker
- **Build Command:** `go build -o under100-web ./cmd/web`
- **Start Command:** `./under100-web`

### 10.2 Required Environment Variables on Render:
In the Render Dashboard, navigate to **Environment** > **Environment Variables** and set:

| Variable Name | Description | Example / Recommended Value |
|---|---|---|
| `PAGODA_APP_ENVIRONMENT` | Environment mode | `prod` |
| `APP_BASE_URL` | Canonical public URL (no trailing slash) | `https://under100.onrender.com` |
| `PAGODA_HTTP_PORT` | HTTP port for Render proxy | `10000` (or Render's `$PORT`) |
| `PAGODA_APP_ENCRYPTIONKEY` | 32+ char session secret | `openssl rand -base64 32` |
| `PAGODA_DATABASE_HOSTNAME` | Render PostgreSQL host | `dpg-xxxxxx.oregon-postgres.render.com` |
| `PAGODA_DATABASE_PORT` | PostgreSQL port | `5432` |
| `PAGODA_DATABASE_USER` | PostgreSQL user | `goship_db_user` |
| `PAGODA_DATABASE_PASSWORD` | PostgreSQL password | Strong random password |
| `PAGODA_DATABASE_DATABASENAMEPROD` | PostgreSQL database name | `goship_db` |
| `PAGODA_DATABASE_SSLMODE` | SSL Mode for cloud PostgreSQL | `require` |
| `PAGODA_CACHE_HOSTNAME` | Render Redis host | `red-xxxxxx.oregon-redis.render.com` |
| `PAGODA_CACHE_PORT` | Redis port | `6379` |
| `RESEND_API_KEY` | Resend API key | `re_xxxxxxxxxxxxxxxx` |
| `EMAIL_FROM` | Verified sender address | `Food Under ₹100 <onboarding@resend.dev>` |

> **Security Note:** Never commit `.env` or paste real secrets into the repository. Configure all secrets strictly inside Render's Environment settings.


## 10. Troubleshooting & Gotchas

1. **PostgreSQL "Connection Refused":**
   Ensure `postgres.exe` is running and accepting TCP connections on port 5432. Check `C:\Users\astha\pgsql\logfile.log` if recovery is in progress.
2. **Redis Session Error:**
   If authentication fails with session storage errors, verify `redis-server.exe` is running on `127.0.0.1:6379`.
3. **Hard Price Cap Enforcement:**
   Dishes priced > ₹100 are strictly rejected upon submission by merchants (HTTP 400 Bad Request) and filtered out of discovery queries.
4. **IDOR Protection:**
   Merchants can only modify dishes, operating hours, and profile details for restaurants they explicitly own. Cross-merchant modifications are rejected with HTTP 403 Forbidden.
