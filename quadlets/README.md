# Podman Quadlets for Timeful

This directory contains example Podman Quadlet files for running Timeful as systemd services.

## Files

- `timeful-mongodb.container` - MongoDB database service
- `timeful-app.container` - Timeful application service
- `timeful.network` - Network configuration for inter-container communication

## Installation

### For User Services (Rootless)

```bash
# Create systemd directory
mkdir -p ~/.config/containers/systemd

# Copy quadlet files
cp quadlets/*.container ~/.config/containers/systemd/
cp quadlets/*.network ~/.config/containers/systemd/

# Create environment file
mkdir -p ~/.config/timeful
cp .env.example ~/.config/timeful/app.env
# Edit app.env with your configuration
nano ~/.config/timeful/app.env

# Build the image
podman build -t localhost/timeful:latest .

# Reload systemd
systemctl --user daemon-reload

# Enable and start services
systemctl --user enable --now timeful-mongodb.service
systemctl --user enable --now timeful-app.service

# Enable linger to keep services running when logged out
loginctl enable-linger $USER
```

### For System Services (Requires Root)

```bash
# Copy quadlet files
sudo cp quadlets/*.container /etc/containers/systemd/
sudo cp quadlets/*.network /etc/containers/systemd/

# Create environment file
sudo mkdir -p /etc/timeful
sudo cp .env.example /etc/timeful/app.env
# Edit app.env with your configuration
sudo nano /etc/timeful/app.env

# Update app.container to use system path
sudo sed -i 's|%h/.config/timeful|/etc/timeful|g' /etc/containers/systemd/timeful-app.container

# Build the image
sudo podman build -t localhost/timeful:latest .

# Reload systemd
sudo systemctl daemon-reload

# Enable and start services
sudo systemctl enable --now timeful-mongodb.service
sudo systemctl enable --now timeful-app.service
```

## Management

Check status:
```bash
systemctl --user status timeful-app.service
```

View logs:
```bash
journalctl --user -u timeful-app.service -f
```

Restart:
```bash
systemctl --user restart timeful-app.service
```

For system services, remove `--user` flag and add `sudo`.

## More Information

See [PODMAN.md](../PODMAN.md) for complete documentation.
