# Podman Quadlets Setup for Timeful

Podman Quadlets allow you to run containers as systemd services, providing better integration with system management.

## Prerequisites

- Podman 4.4 or later
- systemd

## Setup Instructions

### 1. Create Quadlet Files

Quadlet files should be placed in one of these directories:
- System-wide: `/etc/containers/systemd/`
- User-specific: `~/.config/containers/systemd/`

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

# Environment
Environment=MONGO_INITDB_DATABASE=schej-it

# Networking
Network=timeful.network
PublishPort=27017:27017

# Volumes
Volume=timeful-mongodb-data:/data/db
Volume=timeful-mongodb-config:/data/configdb

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

# Networking
Network=timeful.network
PublishPort=3002:80

# Environment (optional)
Environment=VUE_APP_POSTHOG_API_KEY=${VUE_APP_POSTHOG_API_KEY:-}

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

# Optional
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

## Additional Resources

- [Podman Quadlets Documentation](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html)
- [systemd.unit Documentation](https://www.freedesktop.org/software/systemd/man/systemd.unit.html)
