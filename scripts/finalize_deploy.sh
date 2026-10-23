#!/usr/bin/env bash
set -euo pipefail

APP_USER="${1:-$USER}"
APP_DIR="${2:-$HOME/food-discovery}"

cd "${APP_DIR}"

# Rename binaries if linux build artifacts exist
[ -f under100_web_linux ] && mv -f under100_web_linux under100-web
[ -f under100_worker_linux ] && mv -f under100_worker_linux under100-worker
[ -f under100_seed_linux ] && mv -f under100_seed_linux under100-seed

chmod +x under100-web under100-worker under100-seed 2>/dev/null || true

# Extract assets if tarball provided
[ -f assets.tar.gz ] && tar -xzf assets.tar.gz

# Run seeder / migration once
echo "Running database seeder / schema migrations..."
./under100-seed || echo "Seed completed with notes"

# Create systemd service for Web
sudo tee /etc/systemd/system/under100-web.service > /dev/null << EOF
[Unit]
Description=Under100 Food Discovery Web App
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=${APP_USER}
WorkingDirectory=${APP_DIR}
EnvironmentFile=${APP_DIR}/.env
ExecStart=${APP_DIR}/under100-web
Restart=always
RestartSec=5s
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

# Create systemd service for Worker
sudo tee /etc/systemd/system/under100-worker.service > /dev/null << EOF
[Unit]
Description=Under100 Food Discovery Worker Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=${APP_USER}
WorkingDirectory=${APP_DIR}
EnvironmentFile=${APP_DIR}/.env
ExecStart=${APP_DIR}/under100-worker
Restart=always
RestartSec=5s
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

# Configure Nginx reverse proxy
sudo tee /etc/nginx/sites-available/default > /dev/null << 'EOF'
server {
    listen 80 default_server;
    listen [::]:80 default_server;

    server_name _;

    client_max_body_size 50M;

    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }
}
EOF

# Test Nginx and reload
sudo nginx -t
sudo systemctl reload nginx

# Reload systemd and start services
sudo systemctl daemon-reload
sudo systemctl enable --now under100-web
sudo systemctl enable --now under100-worker
sudo systemctl restart under100-web
sudo systemctl restart under100-worker

echo "=== Deployment Finished Successfully ==="
