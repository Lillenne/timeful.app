# Podman Quadlets for Timeful

This directory contains example Podman Quadlet files for running Timeful as systemd services.

## Files

- `timeful-mongodb.container` - MongoDB database service
- `timeful-backend.container` - Backend API server service
- `timeful-frontend.container` - Frontend web server service
- `timeful.network` - Network configuration for inter-container communication

## Installation

### For User Services (Rootless)

```bash
# Create systemd directory
mkdir -p ~/.config/containers/systemd

# Copy quadlet files
cp quadlets/*.container ~/.config/containers/systemd/
cp quadlets/*.network ~/.config/containers/systemd/

# Create environment file for backend
mkdir -p ~/.config/timeful
cp .env.example ~/.config/timeful/backend.env
# Edit backend.env with your configuration
nano ~/.config/timeful/backend.env

# Build the images
podman build -f Dockerfile.backend -t localhost/timeful-backend:latest .
podman build -f Dockerfile.frontend -t localhost/timeful-frontend:latest .

# Reload systemd
systemctl --user daemon-reload

# Enable and start services
systemctl --user enable --now timeful-mongodb.service
systemctl --user enable --now timeful-backend.service
systemctl --user enable --now timeful-frontend.service

# Enable linger to keep services running when logged out
loginctl enable-linger $USER
```

### For System Services (Requires Root)

```bash
# Copy quadlet files
sudo cp quadlets/*.container /etc/containers/systemd/
sudo cp quadlets/*.network /etc/containers/systemd/

# Create environment file for backend
sudo mkdir -p /etc/timeful
sudo cp .env.example /etc/timeful/backend.env
# Edit backend.env with your configuration
sudo nano /etc/timeful/backend.env

# Update backend.container to use system path
sudo sed -i 's|%h/.config/timeful|/etc/timeful|g' /etc/containers/systemd/timeful-backend.container

# Build the images
sudo podman build -f Dockerfile.backend -t localhost/timeful-backend:latest .
sudo podman build -f Dockerfile.frontend -t localhost/timeful-frontend:latest .

# Reload systemd
sudo systemctl daemon-reload

# Enable and start services
sudo systemctl enable --now timeful-mongodb.service
sudo systemctl enable --now timeful-backend.service
sudo systemctl enable --now timeful-frontend.service
```

## Management

Check status:
```bash
systemctl --user status timeful-backend.service
systemctl --user status timeful-frontend.service
```

View logs:
```bash
journalctl --user -u timeful-backend.service -f
journalctl --user -u timeful-frontend.service -f
```

Restart:
```bash
systemctl --user restart timeful-backend.service
systemctl --user restart timeful-frontend.service
```

For system services, remove `--user` flag and add `sudo`.

## More Information

See [PODMAN.md](../PODMAN.md) for complete documentation.
