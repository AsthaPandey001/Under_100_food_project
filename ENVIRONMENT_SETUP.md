# Under100 (Bite100) — Environment Setup & Configuration Guide

This guide details how to set up the local development environment for Under100 on Windows. It specifies required versus optional variables, database startup, and safe development practices.

---

## 1. Prerequisites Checklist

| Component | Minimum Version | Windows Path Example | Status |
| :--- | :--- | :--- | :--- |
| **Go** | 1.24+ | `C:\Users\astha\go` | Required |
| **PostgreSQL** | 15+ | `C:\Users\astha\pgsql` | Required |
| **Redis** | 5.0+ | `C:\Users\astha\redis-server.exe` | Required |
| **PowerShell** | 5.1+ | Built-in Windows | Required |
| **curl.exe** | Built-in | `C:\Windows\System32\curl.exe` | Required |

---

## 2. Environment Variables Specification

The application reads configuration from environment variables, falls back to `config/config.yaml`, and supports `.env` files.

### Required Settings (Development Defaults)

| Variable | Dev Default | Description | Required? |
| :--- | :--- | :--- | :--- |
| `PAGODA_APP_ENVIRONMENT` | `dev` | Application mode (`dev`, `prod`, `test`) | **Yes** |
| `PAGODA_HTTP_HOSTNAME` | `127.0.0.1` | Network interface to bind to | **Yes** |
| `PAGODA_HTTP_PORT` | `8000` | Port for the HTTP web server | **Yes** |
| `PAGODA_APP_ENCRYPTIONKEY` | `abcdefghijklmnopqrstuvwxyz012345` | Min 32-char session encryption key | **Yes** |
| `PAGODA_DATABASE_HOSTNAME` | `localhost` | PostgreSQL host | **Yes** |
| `PAGODA_DATABASE_PORT` | `5432` | PostgreSQL port | **Yes** |
| `PAGODA_DATABASE_USER` | `admin` | Database username | **Yes** |
| `PAGODA_DATABASE_PASSWORD` | `admin` | Database user password | **Yes** |
| `PAGODA_DATABASE_DATABASENAMEPROD` | `goship_db` | Primary database name | **Yes** |
| `PAGODA_DATABASE_SSLMODE` | `disable` | SSL mode (`disable` for local) | **Yes** |
| `PAGODA_CACHE_HOSTNAME` | `localhost` | Redis server host | **Yes** |
| `PAGODA_CACHE_PORT` | `6379` | Redis server port | **Yes** |

### Optional Settings (External Third-Party Services)

| Variable | Default | Purpose |
| :--- | :--- | :--- |
| `PAGODA_MAIL_FROMADDRESS` | `noreply@bite100.example.com` | Outgoing email sender address |
| `PAGODA_MAIL_RESENDAPIKEY` | *(empty / placeholder)* | Resend transactional email API key (mocked in dev) |
| `PAGODA_APP_PUBLICSTRIPEKEY` | *(empty / test key)* | Stripe publishable key |
| `PAGODA_APP_PRIVATESTRIPEKEY` | *(empty / test key)* | Stripe secret key |
| `PAGODA_APP_SENTRYDSN` | *(empty)* | Error monitoring DSN |

> [!IMPORTANT]
> **Security Guideline:** Never commit real secrets, private keys, or API tokens to version control. Production secrets must be injected through encrypted environment secret managers (e.g., Kamal secrets, AWS Secrets Manager, or Doppler).

---

## 3. Local Service Startup Procedures

### 3.1 PostgreSQL 15 Startup
Start PostgreSQL using the data directory `C:\Users\astha\pgdata`:
```powershell
& "C:\Users\astha\pgsql\bin\pg_ctl.exe" start -D "C:\Users\astha\pgdata" -l "C:\Users\astha\pgsql\logfile.log"
```

Verify that PostgreSQL is running and accepting connections:
```powershell
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "SELECT 'PostgreSQL Healthy' as status;"
```

### 3.2 Redis 5 Startup
Start the Redis server process in the background:
```powershell
Start-Process -FilePath "C:\Users\astha\redis-server.exe" -WindowStyle Hidden
```

Verify Redis is running:
```powershell
Get-Process redis-server
```

---

## 4. Building and Launching the Application

Compile the web application binary from source:
```powershell
go build -o under100_server.exe ./cmd/web
```

Run the compiled executable:
```powershell
.\under100_server.exe
```

When started, Echo will output:
```text
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.11.4
High performance, minimalist Go web framework
⇨ http server started on 127.0.0.1:8000
```

---

## 5. Automated Health Check & Verification

Run the master verification script to validate all subsystems:
```powershell
powershell.exe -ExecutionPolicy Bypass -File .\verify_final_handover.ps1
```

If all 13 tests return `[PASS]`, the local development environment is 100% operational.
