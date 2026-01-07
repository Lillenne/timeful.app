# Docker Security Hardening - Summary

This document summarizes the security hardening changes made to the Docker images and configurations.

## Overview

All Docker images have been hardened to run as non-root users with minimal capabilities, ensuring compatibility with both rootful and rootless container runtimes (Docker and Podman).

## Changes Made

### 1. Dockerfiles Updated

#### Dockerfile.backend
- **Added non-root user**: `appuser` with fixed UID/GID (1000:1000)
- **Ownership**: All application files owned by appuser
- **Permissions**: Writable directories (logs) properly configured
- **USER directive**: Container runs as appuser, not root

#### Dockerfile.frontend  
- **Non-root user**: Uses nginx user (UID 101 in nginx:alpine)
- **Ownership**: All nginx files, cache, logs owned by nginx user
- **Permissions**: Writable directories configured for nginx user
- **USER directive**: Container runs as nginx, not root

### 2. Docker Compose Files Updated

All three compose files (docker-compose.yml, docker-compose.dev.yml, docker-compose.ghcr.yml) have been updated with:

#### User Directives
- **MongoDB**: `user: "999:999"` (official mongo image user)
- **Backend**: `user: "1000:1000"` (appuser from Dockerfile)
- **Frontend**: `user: "101:101"` (nginx user from nginx:alpine)

#### Security Options
- **no-new-privileges**: Prevents privilege escalation
- **cap_drop: ALL**: Drops all Linux capabilities by default
- **cap_add**: Only adds back essential capabilities:
  - MongoDB: CHOWN, SETGID, SETUID
  - Backend: None needed
  - Frontend: NET_BIND_SERVICE, CHOWN, SETGID, SETUID

#### Other Improvements
- **Read-only mounts**: config.js mounted as `:ro` (read-only)
- **Volume permissions**: Automatically handled via fixed UIDs

### 3. Documentation Updated

#### DOCKER.md
- New comprehensive security section
- Rootless Docker setup instructions
- Security verification commands
- Security best practices list

#### PODMAN.md
- Updated quadlet configurations with user directives
- Security options added to all quadlets
- New "Rootless Podman Benefits" section
- User namespace mapping explanation
- Compatibility notes for rootless operation

### 4. Testing Added

#### test-docker-security.sh
- Automated security verification script
- Checks for USER directives in Dockerfiles
- Validates security options in compose files
- Confirms all user directives present (1000, 101, 999)
- All checks passing ✓

## Security Features

### Container Security
1. **Non-root execution**: All containers run as dedicated non-root users
2. **Minimal capabilities**: Only essential Linux capabilities enabled
3. **No privilege escalation**: Security options prevent escalation
4. **Read-only where possible**: Configuration files mounted read-only

### Rootless Compatibility
1. **Fixed UIDs**: Consistent user IDs across environments
2. **Docker rootless**: Works with Docker in rootless mode
3. **Podman native**: Works with Podman (rootless by default)
4. **User namespace mapping**: Automatic UID/GID mapping
5. **Volume permissions**: Handled transparently by runtime

## Benefits

### Security
- Reduced attack surface with minimal capabilities
- Container breaches can't escalate to root
- Better isolation between containers and host
- Compatible with security policies requiring non-root

### Compatibility
- Works with Docker (rootful and rootless)
- Works with Podman (rootless by default)
- No sudo/root required for deployment
- Standard across development and production

### Operational
- No manual permission management needed
- Volumes work correctly out of the box
- Backup/restore operations unchanged
- Standard Docker commands work as expected

## Testing

### Automated Testing
Run the security verification script:
```bash
./test-docker-security.sh
```

Expected output: All checks pass ✓

### Manual Testing
Start containers and verify:
```bash
# Start services
docker compose up -d

# Check user IDs in running containers
docker compose exec backend id
# Expected: uid=1000(appuser) gid=1000(appuser)

docker compose exec frontend id
# Expected: uid=101(nginx) gid=101(nginx)

docker compose exec mongodb id
# Expected: uid=999(mongodb) gid=999(mongodb)

# Verify security options
docker inspect timeful-backend | grep -A 10 SecurityOpt
docker inspect timeful-backend | grep -A 5 CapDrop
```

### Rootless Testing

#### Docker Rootless
```bash
# Install Docker rootless (if not done)
dockerd-rootless-setuptool.sh install

# Use normally - works automatically
docker compose up -d
```

#### Podman
```bash
# Install podman-compose
pip3 install podman-compose

# Use normally - rootless by default
podman-compose up -d
```

## Migration Notes

### For Existing Deployments

1. **Backup your data** before updating:
   ```bash
   docker compose exec mongodb mongodump --out=/data/backup
   ```

2. **Pull/rebuild images**:
   ```bash
   # If using pre-built images
   docker compose -f docker-compose.ghcr.yml pull
   
   # If building from source
   docker compose build --no-cache
   ```

3. **Restart services**:
   ```bash
   docker compose down
   docker compose up -d
   ```

4. **Verify containers are running as non-root**:
   ```bash
   docker compose exec backend id
   docker compose exec frontend id
   ```

### Volume Permissions

Named volumes automatically handle permissions correctly. No manual chown operations needed.

If using bind mounts (local directories), you may need to adjust permissions:
```bash
# For backend logs (if using bind mount)
chown -R 1000:1000 ./logs

# For MongoDB data (if using bind mount)  
chown -R 999:999 ./mongodb-data
```

Named volumes (recommended) don't require this.

## Troubleshooting

### Permission Denied Errors

If you see permission denied errors:
1. Verify you're using named volumes (not bind mounts)
2. Check directory ownership if using bind mounts
3. Ensure Docker/Podman has proper subuid/subgid configuration

### Port Binding Issues

Rootless containers may have port binding restrictions:
- Docker rootless: Works automatically
- Podman rootless: Works automatically with slirp4netns
- Ports < 1024 work fine (mapped internally)

### Container Won't Start

Check logs:
```bash
docker compose logs backend
docker compose logs frontend
docker compose logs mongodb
```

Common issues:
- Volume permission problems (use named volumes)
- Missing environment variables (check .env file)
- Port conflicts (change host ports in compose file)

## Additional Resources

- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [Rootless Docker](https://docs.docker.com/engine/security/rootless/)
- [Podman Rootless](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md)
- [Linux Capabilities](https://man7.org/linux/man-pages/man7/capabilities.7.html)

## Support

If you encounter issues:
- Check the [DOCKER.md](./DOCKER.md) documentation
- Check the [PODMAN.md](./PODMAN.md) documentation
- Run `./test-docker-security.sh` to verify configuration
- Open an issue on GitHub with logs and error messages
