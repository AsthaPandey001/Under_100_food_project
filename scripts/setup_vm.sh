#!/usr/bin/env bash
set -euo pipefail

APP_USER="${1:-$USER}"
APP_DIR="${2:-$HOME/food-discovery}"

echo "=== 1. Checking Sudo Permissions for ${APP_USER} ==="
if ! sudo -n true 2>/dev/null; then
    echo "Configuring passwordless sudo for ${APP_USER}..."
    sudo bash -c "echo '${APP_USER} ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/${APP_USER} && chmod 0440 /etc/sudoers.d/${APP_USER}"
fi
echo "Sudo permissions active."

echo "=== 2. Updating System and Installing Packages ==="
sudo apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y \
    postgresql postgresql-contrib \
    redis-server \
    nginx \
    curl \
    tar \
    rsync \
    ca-certificates

echo "=== 3. Configuring PostgreSQL ==="
sudo systemctl enable --now postgresql
sudo -u postgres psql -c "DO \$\$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'admin') THEN CREATE ROLE admin WITH LOGIN SUPERUSER PASSWORD 'admin'; END IF; END \$\$;"
sudo -u postgres psql -c "ALTER ROLE admin WITH PASSWORD 'admin';"
sudo -u postgres psql -c "SELECT 'CREATE DATABASE goship_db OWNER admin' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'goship_db')\gexec"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE goship_db TO admin;"
echo "PostgreSQL configured successfully."

echo "=== 4. Configuring Redis ==="
sudo systemctl enable --now redis-server
sudo systemctl is-active redis-server
echo "Redis active."

echo "=== 5. Setting up Application Directory ==="
mkdir -p "${APP_DIR}/logs"
mkdir -p "${APP_DIR}/config"

echo "=== 6. Environment File ==="
if [ ! -f "${APP_DIR}/.env" ]; then
    cat << 'EOF' > "${APP_DIR}/.env"
PAGODA_APP_ENVIRONMENT=dev
PAGODA_HTTP_HOSTNAME=0.0.0.0
PAGODA_HTTP_PORT=8000
PAGODA_HTTP_DOMAIN=http://127.0.0.1:8000
PAGODA_APP_ENCRYPTIONKEY=abcdefghijklmnopqrstuvwxyz012345
PAGODA_DATABASE_DBMODE=standalone
PAGODA_DATABASE_HOSTNAME=127.0.0.1
PAGODA_DATABASE_PORT=5432
PAGODA_DATABASE_USER=admin
PAGODA_DATABASE_PASSWORD=admin
PAGODA_DATABASE_DATABASENAMELOCAL=goship_db
PAGODA_DATABASE_DATABASENAMEPROD=goship_db
PAGODA_DATABASE_SSLMODE=disable
PAGODA_CACHE_HOSTNAME=127.0.0.1
PAGODA_CACHE_PORT=6379
PAGODA_CACHE_PASSWORD=
PAGODA_CACHE_DATABASE=0
PAGODA_MAIL_HOSTNAME=mock
EOF
fi

echo "Setup completed successfully for ${APP_USER} in ${APP_DIR}."
