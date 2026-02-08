# Podman Quadlets Setup for Timeful

Podman Quadlets allow you to run containers as systemd services, providing better integration with system management. Podman runs rootless by default, providing enhanced security.

## Prerequisites

- Podman 4.4 or later
- systemd

## Security Features

All container images are hardened with security best practices:
- **Non-root users**: Containers run as dedicated non-root users
- **Minimal capabilities**: Only essential Linux capabilities are enabled
- **No privilege escalation**: Security options prevent privilege escalation
- **Rootless by default**: Podman runs without requiring root privileges

## Setup Instructions

### 1. Create Quadlet Files

Quadlet files should be placed in one of these directories:
- System-wide: `/etc/containers/systemd/` (requires root)
- User-specific: `~/.config/containers/systemd/` (recommended - runs rootless)

For user services (recommended for non-root deployment):

```bash
mkdir -p ~/.config/containers/systemd
```

### 2. MongoDB Quadlet

Create `~/.config/containers/systemd/timeful-mongodb.container`:

```ini
[Unit]
Description=Timeful MongoDB Database
After=network-online.target
Wants=network-online.target

[Container]
Image=docker.io/library/mongo:6.0
ContainerName=timeful-mongodb
AutoUpdate=registry

# Run as non-root user (mongodb user in official image)
User=999:999

# Environment
Environment=MONGO_INITDB_DATABASE=schej-it

# Networking
Network=timeful.network
PublishPort=27017:27017

# Volumes
Volume=timeful-mongodb-data:/data/db
Volume=timeful-mongodb-config:/data/configdb

# Security options
SecurityLabelDisable=false
NoNewPrivileges=true

# Health check
HealthCmd=mongosh --eval "db.adminCommand('ping')" --quiet
HealthInterval=10s
HealthTimeout=5s
HealthRetries=5

[Service]
Restart=unless-stopped
TimeoutStartSec=900

[Install]
WantedBy=default.target
```

### 3. Timeful Backend Quadlet

Create `~/.config/containers/systemd/timeful-backend.container`:

```ini
[Unit]
Description=Timeful Backend API Server
After=timeful-mongodb.service
Requires=timeful-mongodb.service
After=network-online.target
Wants=network-online.target

[Container]
Image=localhost/timeful-backend:latest
ContainerName=timeful-backend
AutoUpdate=registry

# Run as non-root user (appuser defined in image)
User=1000:1000

# Environment variables
Environment=MONGO_URI=mongodb://timeful-mongodb:27017
Environment=MONGO_DB_NAME=schej-it
Environment=GIN_MODE=release
EnvironmentFile=%h/.config/timeful/backend.env

# Networking
Network=timeful.network
# Backend port exposed internally only

# Volumes
Volume=timeful-backend-logs:/app/logs

# Security options
SecurityLabelDisable=false
NoNewPrivileges=true

[Service]
Restart=unless-stopped
TimeoutStartSec=900

[Install]
WantedBy=default.target
```

### 3. Timeful Frontend Quadlet

Create `~/.config/containers/systemd/timeful-frontend.container`:

```ini
[Unit]
Description=Timeful Frontend Web Server
After=timeful-backend.service
Requires=timeful-backend.service
After=network-online.target
Wants=network-online.target

[Container]
Image=localhost/timeful-frontend:latest
ContainerName=timeful-frontend
AutoUpdate=registry

# Run as non-root user (nginx user in nginx:alpine image)
User=101:101

# Networking
Network=timeful.network
PublishPort=3002:80

# Environment (optional)
Environment=BACKEND_HOST=timeful-backend
Environment=BACKEND_PORT=3002

# Security options
SecurityLabelDisable=false
NoNewPrivileges=true

[Service]
Restart=unless-stopped
TimeoutStartSec=900

[Install]
WantedBy=default.target
```

### 4. Network Quadlet

Create `~/.config/containers/systemd/timeful.network`:

```ini
[Network]
NetworkName=timeful
Driver=bridge
```

### 5. Volume Quadlets

Create `~/.config/containers/systemd/timeful-mongodb-data.volume`:

```ini
[Volume]
VolumeName=timeful-mongodb-data
```

Create `~/.config/containers/systemd/timeful-mongodb-config.volume`:

```ini
[Volume]
VolumeName=timeful-mongodb-config
```

Create `~/.config/containers/systemd/timeful-backend-logs.volume`:

```ini
[Volume]
VolumeName=timeful-backend-logs
```

### 6. Environment File

Create `~/.config/timeful/backend.env` with your configuration:

```env
# Required
CLIENT_ID=your_google_client_id
CLIENT_SECRET=your_google_client_secret
ENCRYPTION_KEY=your_encryption_key

# Optional - Microsoft OAuth (for Outlook calendar integration)
MICROSOFT_CLIENT_ID=your_microsoft_client_id
MICROSOFT_CLIENT_SECRET=your_microsoft_client_secret

# Optional - Other services
GMAIL_APP_PASSWORD=
SCHEJ_EMAIL_ADDRESS=
STRIPE_API_KEY=
```

Make sure the file has proper permissions:

```bash
chmod 600 ~/.config/timeful/backend.env
```

### 7. Build the Images

Before starting services, build the backend and frontend images:

```bash
cd /path/to/timeful.app
podman build -f Dockerfile.backend -t localhost/timeful-backend:latest .
podman build -f Dockerfile.frontend -t localhost/timeful-frontend:latest .
```

### 8. Reload systemd and Start Services

```bash
# Reload systemd to detect new quadlet files
systemctl --user daemon-reload

# Enable services to start on boot
systemctl --user enable timeful-mongodb.service
systemctl --user enable timeful-backend.service
systemctl --user enable timeful-frontend.service

# Start services
systemctl --user start timeful-mongodb.service
systemctl --user start timeful-backend.service
systemctl --user start timeful-frontend.service
```

### 9. Enable Linger (Optional)

To keep services running even when not logged in:

```bash
loginctl enable-linger $USER
```

## Managing Services

### Check Status

```bash
systemctl --user status timeful-mongodb.service
systemctl --user status timeful-backend.service
```

### View Logs

```bash
journalctl --user -u timeful-backend.service -f
journalctl --user -u timeful-mongodb.service -f
```

### Restart Services

```bash
systemctl --user restart timeful-backend.service
systemctl --user restart timeful-mongodb.service
```

### Stop Services

```bash
systemctl --user stop timeful-backend.service
systemctl --user stop timeful-mongodb.service
```

### Disable Services

```bash
systemctl --user disable timeful-backend.service
systemctl --user disable timeful-mongodb.service
```

## Updating the Application

When you want to update to a new version:

```bash
# Stop services
systemctl --user stop timeful-backend.service

# Pull latest code
cd /path/to/timeful.app
git pull

# Rebuild images
podman build -f Dockerfile.backend -t localhost/timeful-backend:latest .
podman build -f Dockerfile.frontend -t localhost/timeful-frontend:latest .

# Start services
systemctl --user start timeful-backend.service
systemctl --user start timeful-frontend.service
```

## Backup and Restore

### Backup

```bash
# Create backup directory
mkdir -p ~/timeful-backups

# Backup MongoDB data
podman exec timeful-mongodb mongodump --db=schej-it --out=/data/backup
podman cp timeful-mongodb:/data/backup ~/timeful-backups/backup-$(date +%Y%m%d-%H%M%S)
```

### Restore

```bash
# Copy backup to container
podman cp ~/timeful-backups/backup-YYYYMMDD-HHMMSS timeful-mongodb:/data/restore

# Restore database
podman exec timeful-mongodb mongorestore --db=schej-it /data/restore/schej-it --drop
```

## Troubleshooting

### Check if quadlet files are detected

```bash
systemctl --user list-units "timeful-*"
```

### View container logs directly

```bash
podman logs timeful-backend
podman logs timeful-mongodb
```

### Inspect containers

```bash
podman inspect timeful-backend
podman inspect timeful-mongodb
```

### Check network connectivity

```bash
podman exec timeful-backend ping -c 3 timeful-mongodb
```

### Verify security configuration

Check that containers are running as non-root:

```bash
# Check user IDs in running containers
podman exec timeful-backend id
podman exec timeful-frontend id
podman exec timeful-mongodb id

# Verify no-new-privileges setting
podman inspect timeful-backend | grep -i "NoNewPrivileges"
podman inspect timeful-frontend | grep -i "NoNewPrivileges"
```

## Rootless Podman Benefits

Podman runs rootless by default when used as a regular user, providing several security advantages:

### Security Advantages

1. **No Root Required**: Containers run entirely within your user namespace
   - No need for sudo or root privileges
   - Container breaches cannot compromise the host system
   - Natural isolation between users

2. **User Namespace Mapping**: UIDs/GIDs are automatically mapped
   - Container UID 1000 (appuser) maps to your user's subuid range
   - Host filesystem is protected from container processes
   - Works transparently with our fixed UIDs

3. **Enhanced Security**: Multiple layers of protection
   - SELinux/AppArmor integration (where available)
   - Seccomp filtering
   - No new privileges by default
   - Capability dropping

### Compatibility Notes

The hardened Docker images work seamlessly with rootless Podman:

- **Fixed UIDs**: Backend (1000), Frontend (101), MongoDB (999)
  - These are automatically mapped to your user's subuid range
  - Volume permissions work correctly out of the box
  - No manual chown operations needed

- **Volume Mounts**: Named volumes handle permissions automatically
  - Podman manages UID/GID mappings transparently
  - Files created by containers are owned by your user on the host
  - Backup and restore operations work normally

- **Port Binding**: Ports < 1024 work with rootless Podman
  - Frontend binds to port 80 internally (mapped to 3002 on host)
  - Podman allows unprivileged port binding with slirp4netns
  - No special configuration needed

### Using Podman Compose

For simpler deployment without systemd:

```bash
# Install podman-compose
pip3 install podman-compose

# Use with any docker-compose.yml file
podman-compose -f docker-compose.yml up -d
podman-compose -f docker-compose.yml logs -f
podman-compose -f docker-compose.yml down

# All security features work automatically
# No root required, completely rootless
```

### Rootless vs Rootful

**Rootless (Recommended)**:
- Runs as regular user
- More secure
- No sudo needed
- User systemd services
- Limited to your user's processes

**Rootful** (System-wide):
- Requires root/sudo
- System-wide services
- Traditional Docker-like behavior
- Place quadlets in `/etc/containers/systemd/`

For most use cases, rootless mode is recommended.

## System-wide Installation

To run as system services (requires root), place quadlet files in `/etc/containers/systemd/` and use:

```bash
sudo systemctl daemon-reload
sudo systemctl enable timeful-mongodb.service
sudo systemctl enable timeful-backend.service
sudo systemctl start timeful-mongodb.service
sudo systemctl start timeful-backend.service
```

Use `systemctl` instead of `systemctl --user` for all management commands.

## Benefits of Quadlets

- **systemd integration**: Standard Linux service management
- **Auto-restart**: Services restart on failure
- **Dependency management**: Services start in correct order
- **Logging**: Centralized logging via journald
- **Resource management**: Use systemd resource controls
- **Boot integration**: Start services at system boot

## Self-Hosting Listmonk with Podman

Timeful can be integrated with [Listmonk](https://listmonk.app/), a self-hosted newsletter and mailing list manager. The Listmonk services are included in the Docker Compose files and work seamlessly with Podman.

### Using Podman Compose with Listmonk

The easiest way to run Listmonk with Podman is using `podman-compose`:

```bash
# First, copy the example configuration file
cp listmonk-config-example.toml listmonk-config.toml

# Edit the configuration file with your secure passwords
nano listmonk-config.toml

# Start all services including Listmonk
podman-compose up -d

# Or start only Listmonk services
podman-compose up -d listmonk-db listmonk
```

The database schema is automatically initialized on first startup. Access Listmonk at http://localhost:9000 with the credentials you configured in `listmonk-config.toml`.

For complete Listmonk configuration instructions, see the [DOCKER.md](./DOCKER.md#self-hosting-listmonk-with-docker-compose) guide. All instructions apply to Podman as well.

### Adding Listmonk to Quadlets

If you're using Podman Quadlets (systemd integration), you can add Listmonk services:

#### 1. PostgreSQL for Listmonk

Create `~/.config/containers/systemd/timeful-listmonk-db.container`:

```ini
[Unit]
Description=Timeful Listmonk PostgreSQL Database
After=network-online.target
Wants=network-online.target

[Container]
Image=docker.io/library/postgres:15-alpine
ContainerName=timeful-listmonk-db
AutoUpdate=registry

# Run as postgres user (uid 70 in postgres:alpine image)
User=70:70

# Environment
Environment=POSTGRES_PASSWORD=listmonk
Environment=POSTGRES_USER=listmonk
Environment=POSTGRES_DB=listmonk
EnvironmentFile=%h/.config/timeful/listmonk-db.env

# Networking
Network=timeful.network

# Volumes
Volume=timeful-listmonk-db-data:/var/lib/postgresql/data

# Health check
HealthCmd=pg_isready -U listmonk
HealthInterval=10s
HealthTimeout=5s
HealthRetries=5

# Security options
SecurityLabelDisable=false
NoNewPrivileges=true

[Service]
Restart=unless-stopped
TimeoutStartSec=900

[Install]
WantedBy=default.target
```

#### 2. Listmonk Service

Create `~/.config/containers/systemd/timeful-listmonk.container`:

```ini
[Unit]
Description=Timeful Listmonk Newsletter Manager
After=timeful-listmonk-db.service
Requires=timeful-listmonk-db.service
After=network-online.target
Wants=network-online.target

[Container]
Image=docker.io/listmonk/listmonk:latest
ContainerName=timeful-listmonk
AutoUpdate=registry

# Run as non-root user
User=1000:1000

# Networking
Network=timeful.network
PublishPort=9000:9000

# Volumes - mount the config file
Volume=%h/.config/timeful/listmonk-config.toml:/listmonk/config.toml:ro

# Environment
Environment=TZ=Etc/UTC

# Command - automatically initialize/upgrade database on startup
Exec=sh -c "./listmonk --install --idempotent --yes && ./listmonk --upgrade --yes && ./listmonk"

# Security options
SecurityLabelDisable=false
NoNewPrivileges=true

[Service]
Restart=unless-stopped
TimeoutStartSec=900

[Install]
WantedBy=default.target
```

#### 3. Volume Quadlet

Create `~/.config/containers/systemd/timeful-listmonk-db-data.volume`:

```ini
[Volume]
VolumeName=timeful-listmonk-db-data
```

#### 4. Configuration Files

Create `~/.config/timeful/listmonk-config.toml`:

```toml
[app]
address = "0.0.0.0:9000"
admin_username = "admin"
admin_password = "listmonk"

[db]
host = "timeful-listmonk-db"
port = 5432
user = "listmonk"
password = "listmonk"
database = "listmonk"
ssl_mode = "disable"
max_open = 25
max_idle = 25
max_lifetime = "300s"
```

Create `~/.config/timeful/listmonk-db.env` (optional, for custom password):

```env
POSTGRES_PASSWORD=your_secure_password
```

Make sure files have proper permissions:

```bash
chmod 600 ~/.config/timeful/listmonk-config.toml
chmod 600 ~/.config/timeful/listmonk-db.env
```

#### 5. Enable and Start Services

```bash
# Reload systemd
systemctl --user daemon-reload

# Enable services
systemctl --user enable timeful-listmonk-db.service
systemctl --user enable timeful-listmonk.service

# Start services
systemctl --user start timeful-listmonk-db.service
systemctl --user start timeful-listmonk.service

# Initialize Listmonk on first run
podman exec timeful-listmonk ./listmonk --install
```

#### 6. Configure Timeful Backend

Update `~/.config/timeful/backend.env` to include Listmonk configuration:

```env
# ... existing configuration ...

# Listmonk Configuration
LISTMONK_URL=http://timeful-listmonk:9000
LISTMONK_USERNAME=admin
LISTMONK_PASSWORD=listmonk
LISTMONK_LIST_ID=1
```

Restart the backend service:

```bash
systemctl --user restart timeful-backend.service
```

### Managing Listmonk with Quadlets

```bash
# Check status
systemctl --user status timeful-listmonk.service
systemctl --user status timeful-listmonk-db.service

# View logs
journalctl --user -u timeful-listmonk.service -f
journalctl --user -u timeful-listmonk-db.service -f

# Restart services
systemctl --user restart timeful-listmonk.service

# Stop services
systemctl --user stop timeful-listmonk.service timeful-listmonk-db.service

# Disable services (won't start on boot)
systemctl --user disable timeful-listmonk.service timeful-listmonk-db.service
```

### Backup Listmonk Database (Podman)

```bash
# Create backup directory
mkdir -p ~/timeful-backups

# Backup Listmonk database
podman exec timeful-listmonk-db pg_dump -U listmonk listmonk > ~/timeful-backups/listmonk-backup-$(date +%Y%m%d-%H%M%S).sql

# Restore from backup
podman exec -i timeful-listmonk-db psql -U listmonk listmonk < ~/timeful-backups/listmonk-backup-YYYYMMDD-HHMMSS.sql
```

### Rootless Podman with Listmonk

Listmonk works seamlessly with rootless Podman:

- **PostgreSQL** runs as UID 70 (postgres user), automatically mapped to your subuid range
- **Listmonk** runs as UID 1000 (non-root user), mapped transparently
- **Volumes** handle permissions automatically
- **Port 9000** binds successfully with slirp4netns (rootless port binding)

No special configuration needed - everything works out of the box with rootless Podman!

### Troubleshooting Listmonk with Podman

**Check if containers are running**:
```bash
podman ps | grep listmonk
```

**Check database connectivity**:
```bash
podman exec timeful-listmonk-db pg_isready -U listmonk
```

**View detailed logs**:
```bash
podman logs timeful-listmonk --tail 50
podman logs timeful-listmonk-db --tail 50
```

**Test network connectivity**:
```bash
podman exec timeful-listmonk ping -c 3 timeful-listmonk-db
```

**Verify security settings**:
```bash
# Check user IDs
podman exec timeful-listmonk id
podman exec timeful-listmonk-db id

# Verify no-new-privileges
podman inspect timeful-listmonk | grep -i "NoNewPrivileges"
```

For complete Listmonk configuration and usage instructions, see the [DOCKER.md](./DOCKER.md#self-hosting-listmonk-with-docker-compose) guide.

## Additional Resources

- [Podman Quadlets Documentation](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html)
- [systemd.unit Documentation](https://www.freedesktop.org/software/systemd/man/systemd.unit.html)
- [Listmonk Documentation](https://listmonk.app/docs/)
