# Self-Hosting Feature - Changes Summary

## Overview

This implementation adds complete Docker and Podman support for self-hosting Timeful, addressing issue #118.

## New Files

### Docker Support
- **Dockerfile** - Multi-stage build for frontend and backend
- **Dockerfile.backend** - Separate backend-only build (alternative)
- **Dockerfile.frontend** - Separate frontend-only build (alternative)
- **docker-compose.yml** - Production deployment configuration
- **docker-compose.dev.yml** - Development/testing configuration
- **.dockerignore** - Optimized build context
- **.env.example** - Environment variables template with documentation

### Documentation
- **DOCKER.md** - Comprehensive Docker deployment guide
  - Quick start instructions
  - Environment configuration
  - Production deployment with reverse proxy
  - Management commands
  - Backup/restore procedures
  - Troubleshooting guide
  
- **PODMAN.md** - Complete Podman Quadlets guide
  - Systemd integration
  - User and system services
  - Service management
  - Backup procedures
  
- **TESTING.md** - Testing guide for deployments
  - Quick test procedures
  - Full production test
  - Troubleshooting tests
  - Automated testing script

### Podman Support
- **quadlets/** - Directory with Podman Quadlet examples
  - **timeful-mongodb.container** - MongoDB service quadlet
  - **timeful-app.container** - Application service quadlet
  - **timeful.network** - Network configuration
  - **README.md** - Quadlet installation guide

### Convenience Tools
- **Makefile** - Common management commands
  - setup, build, up, down, restart
  - logs, backup, restore
  - dev, pull, status, clean
  
- **docker-setup.sh** - Interactive setup script
  - Checks Docker installation
  - Creates .env from template
  - Guides through configuration
  - Starts services

### CI/CD
- **.github/workflows/docker-build.yml** - Docker build testing workflow

## Modified Files

### Backend Changes

**server/db/init.go**
- Added environment variable support for MongoDB URI
- Added support for custom database name
- Default to localhost if not specified for backward compatibility

```go
mongoURI := os.Getenv("MONGO_URI")
if mongoURI == "" {
    mongoURI = "mongodb://localhost:27017"
}

dbName := os.Getenv("MONGO_DB_NAME")
if dbName == "" {
    dbName = "schej-it"
}
```

**server/main.go**
- Made .env file optional
- Environment variables can be provided directly (Docker)
- Logs warning instead of panicking if .env missing

```go
if err != nil {
    // In Docker or production environments, .env file might not exist
    logger.StdOut.Println("No .env file found, using environment variables directly")
}
```

### Documentation Updates

**README.md**
- Added self-hosting section
- Quick start instructions
- Links to Docker and Podman documentation
- System requirements

## Features Implemented

### Docker Deployment
✅ Multi-stage build combining frontend and backend
✅ Optimized Alpine-based final image
✅ MongoDB integration with health checks
✅ Persistent volumes for data
✅ Environment-based configuration
✅ Development and production modes
✅ Automated setup script
✅ Makefile for easy management
✅ Reverse proxy examples (Nginx, Caddy)

### Podman Support
✅ Complete Quadlet configuration
✅ Systemd integration
✅ User and system service modes
✅ Auto-restart on failure
✅ Dependency management
✅ Logging via journald

### Configuration
✅ Environment variable support for all settings
✅ Optional .env file
✅ Backward compatible with existing deployments
✅ Secure defaults
✅ Comprehensive .env.example with documentation

### Documentation
✅ Detailed deployment guide
✅ Production setup with SSL
✅ Backup and restore procedures
✅ Troubleshooting guide
✅ Testing procedures
✅ Management commands reference

## Breaking Changes

None - All changes are backward compatible. Existing deployments continue to work unchanged.

## Environment Variables

### Required
- `CLIENT_ID` - Google OAuth client ID
- `CLIENT_SECRET` - Google OAuth client secret
- `ENCRYPTION_KEY` - Data encryption key

### Optional
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `MONGO_DB_NAME` - Database name (default: schej-it)
- `GMAIL_APP_PASSWORD` - Gmail app password for emails
- `SCHEJ_EMAIL_ADDRESS` - Email sender address
- `STRIPE_API_KEY` - Stripe API key for payments
- `LISTMONK_*` - Listmonk configuration for campaigns
- `SLACK_*` - Slack webhooks for monitoring
- `VUE_APP_POSTHOG_API_KEY` - PostHog analytics key

## Architecture

```
┌─────────────────────────────────────────┐
│          Reverse Proxy (Optional)       │
│         Nginx / Caddy / Traefik         │
└──────────────────┬──────────────────────┘
                   │
                   │ :80/:443
                   │
┌──────────────────▼──────────────────────┐
│         Timeful Application             │
│  ┌────────────────────────────────────┐ │
│  │   Frontend (Vue.js - Static)       │ │
│  │   Served by Go server              │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │   Backend (Go)                     │ │
│  │   - API endpoints                  │ │
│  │   - Static file serving            │ │
│  │   Port: 3002                       │ │
│  └────────────────────────────────────┘ │
└──────────────────┬──────────────────────┘
                   │
                   │ mongodb://
                   │
┌──────────────────▼──────────────────────┐
│          MongoDB Database               │
│          Port: 27017                    │
│          Volume: persistent             │
└─────────────────────────────────────────┘
```

## Security Considerations

✅ No secrets in repository
✅ .env.example as template only
✅ .dockerignore prevents sensitive file inclusion
✅ Environment variables for all sensitive data
✅ Encryption key required
✅ Optional Google OAuth
✅ Persistent volumes for data protection
✅ Health checks for service reliability

## Testing

- Docker Compose configuration validated ✅
- Dockerfile syntax verified ✅
- Multi-stage build tested ✅
- Environment variables tested ✅
- Documentation reviewed ✅
- GitHub Actions workflow added ✅

## Deployment Options

1. **Docker Compose** (Recommended)
   - Simple setup with `docker compose up`
   - Production-ready configuration
   - Easy updates with `make pull`

2. **Docker Standalone**
   - Build: `docker build -t timeful .`
   - Run: `docker run -p 3002:3002 timeful`

3. **Podman Compose**
   - Drop-in replacement for Docker Compose
   - `podman-compose up -d`

4. **Podman Quadlets**
   - Systemd integration
   - Auto-start on boot
   - System service management

## Future Enhancements

Potential improvements for future releases:

- [ ] Pre-built Docker images on Docker Hub
- [ ] GitHub Container Registry (ghcr.io) images
- [ ] Kubernetes manifests (Helm chart)
- [ ] Ansible playbook for deployment
- [ ] Docker Swarm configuration
- [ ] Monitoring stack integration (Prometheus/Grafana)
- [ ] Backup automation scripts
- [ ] Multi-architecture builds (ARM support)

## Migration Guide

For existing manual deployments:

1. Pull latest code with Docker support
2. Create `.env` from `.env.example`
3. Migrate existing `.env` values
4. Stop manual services
5. Start with Docker: `docker compose up -d`
6. Verify data migration from MongoDB

## Community Impact

This implementation:
- Addresses issue #118 and community demand
- Enables easy self-hosting for organizations
- Reduces barrier to entry for new users
- Provides multiple deployment options
- Maintains backward compatibility
- Includes comprehensive documentation
- Supports both Docker and Podman ecosystems

## Credits

Based on community feedback and contributions from:
- @gambit-dross (Issue creator)
- @athamour1 (Initial Docker configuration proposal)
- Community members requesting this feature

## Support

For help with self-hosting:
- Documentation: DOCKER.md, PODMAN.md, TESTING.md
- GitHub Issues: https://github.com/schej-it/timeful.app/issues
- Discord: https://discord.gg/v6raNqYxx3
