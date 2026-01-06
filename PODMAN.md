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

### 3. Timeful App Quadlet

Create `~/.config/containers/systemd/timeful-app.container`:

```ini
[Unit]
Description=Timeful Application
After=timeful-mongodb.service
Requires=timeful-mongodb.service
After=network-online.target
Wants=network-online.target

[Container]
Image=localhost/timeful:latest
ContainerName=timeful-app
AutoUpdate=registry

# Environment variables
Environment=MONGO_URI=mongodb://timeful-mongodb:27017
Environment=MONGO_DB_NAME=schej-it
Environment=GIN_MODE=release
EnvironmentFile=%h/.config/timeful/app.env

# Networking
Network=timeful.network
PublishPort=3002:3002

# Volumes
Volume=timeful-app-logs:/app/logs.log

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

Create `~/.config/containers/systemd/timeful-app-logs.volume`:

```ini
[Volume]
VolumeName=timeful-app-logs
```

### 6. Environment File

Create `~/.config/timeful/app.env` with your configuration:

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
chmod 600 ~/.config/timeful/app.env
```

### 7. Build the Image

Before starting services, build the Timeful image:

```bash
cd /path/to/timeful.app
podman build -t localhost/timeful:latest .
```

### 8. Reload systemd and Start Services

```bash
# Reload systemd to detect new quadlet files
systemctl --user daemon-reload

# Enable services to start on boot
systemctl --user enable timeful-mongodb.service
systemctl --user enable timeful-app.service

# Start services
systemctl --user start timeful-mongodb.service
systemctl --user start timeful-app.service
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
systemctl --user status timeful-app.service
```

### View Logs

```bash
journalctl --user -u timeful-app.service -f
journalctl --user -u timeful-mongodb.service -f
```

### Restart Services

```bash
systemctl --user restart timeful-app.service
systemctl --user restart timeful-mongodb.service
```

### Stop Services

```bash
systemctl --user stop timeful-app.service
systemctl --user stop timeful-mongodb.service
```

### Disable Services

```bash
systemctl --user disable timeful-app.service
systemctl --user disable timeful-mongodb.service
```

## Updating the Application

When you want to update to a new version:

```bash
# Stop services
systemctl --user stop timeful-app.service

# Pull latest code
cd /path/to/timeful.app
git pull

# Rebuild image
podman build -t localhost/timeful:latest .

# Start service
systemctl --user start timeful-app.service
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
podman logs timeful-app
podman logs timeful-mongodb
```

### Inspect containers

```bash
podman inspect timeful-app
podman inspect timeful-mongodb
```

### Check network connectivity

```bash
podman exec timeful-app ping -c 3 timeful-mongodb
```

## System-wide Installation

To run as system services (requires root), place quadlet files in `/etc/containers/systemd/` and use:

```bash
sudo systemctl daemon-reload
sudo systemctl enable timeful-mongodb.service
sudo systemctl enable timeful-app.service
sudo systemctl start timeful-mongodb.service
sudo systemctl start timeful-app.service
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
