# 🐳 Timeful Self-Hosting Quick Reference

## 🚀 Quick Start (5 minutes)

### Option 1: Automated Setup (Recommended)

```bash
git clone https://github.com/schej-it/timeful.app.git
cd timeful.app
./docker-setup.sh
```

### Option 2: Manual Setup

```bash
# 1. Clone repository
git clone https://github.com/schej-it/timeful.app.git
cd timeful.app

# 2. Create environment file
cp .env.example .env

# 3. Generate encryption key
echo "ENCRYPTION_KEY=$(openssl rand -base64 32)" >> .env

# 4. Edit .env and add Google OAuth credentials
nano .env

# 5. Start application
make up
# or: docker compose up -d
```

### Option 3: Using Makefile

```bash
make setup    # Interactive setup
make up       # Start services
make logs     # View logs
```

## 📋 Prerequisites Checklist

- [ ] Docker 20.10+ installed
- [ ] Docker Compose v2 installed
- [ ] Google Cloud account (for OAuth)
- [ ] 4GB+ free disk space
- [ ] Ports 3002 and 27017 available

## 🔑 Required Configuration

### 1. Google OAuth (Required)

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create/select project
3. Enable APIs:
   - Google Calendar API
   - Google People API
4. Create OAuth 2.0 credentials (Web application)
5. Add redirect URI: `http://localhost:3002/api/auth/google/callback`
6. Copy Client ID and Secret to `.env`

### 2. Encryption Key (Required)

```bash
openssl rand -base64 32
```

Add to `.env`:
```
ENCRYPTION_KEY=your_generated_key_here
```

### 3. Microsoft OAuth (Optional - for Outlook calendar integration)

1. Go to [Azure Portal](https://portal.azure.com/)
2. Navigate to **Azure Active Directory** → **App registrations** → **New registration**
3. Configure:
   - **Name:** `Timeful Self-Hosted`
   - **Account types:** Multitenant + Personal
   - **Redirect URI (Web):** `http://localhost:3002/auth`
4. Copy **Application (client) ID** to `.env` as `MICROSOFT_CLIENT_ID`
5. Go to **Certificates & secrets** → **New client secret**
6. Copy secret **Value** to `.env` as `MICROSOFT_CLIENT_SECRET`
7. Go to **API permissions** → Add these Microsoft Graph delegated permissions:
   - `offline_access`
   - `User.Read`
   - `Calendars.Read`
8. Update `public/config.js` with your `microsoftClientId`

**Note:** If you skip this, Outlook calendar integration will not work, but Google Calendar will still function.

## 🎯 Common Commands

| Action | Command | Description |
|--------|---------|-------------|
| **Setup** | `make setup` | Interactive setup wizard |
| **Start** | `make up` | Start all services |
| **Stop** | `make down` | Stop all services |
| **Logs** | `make logs` | View all logs |
| **Restart** | `make restart` | Restart services |
| **Status** | `make status` | Check service status |
| **Backup** | `make backup` | Backup database |
| **Update** | `make pull` | Pull updates & restart |
| **Clean** | `make clean` | Remove all data ⚠️ |

## 🌐 Access Points

After starting:

- **Application**: http://localhost:3002
- **API Docs**: http://localhost:3002/swagger/index.html
- **MongoDB**: mongodb://localhost:27017

## 🔍 Troubleshooting

### Services won't start

```bash
make logs-app    # Check app logs
make logs-db     # Check database logs
make status      # Check container status
```

### OAuth errors

- ✅ Verify redirect URI matches
- ✅ Check environment variables
- ✅ Ensure APIs are enabled

### Port conflicts

Change port in `docker-compose.yml`:
```yaml
ports:
  - "8080:3002"  # Change 8080 to any free port
```

### Database connection failed

```bash
# Test connectivity
docker compose exec app ping -c 3 mongodb

# Check MongoDB health
docker compose exec mongodb mongosh --eval "db.adminCommand('ping')"
```

## 📚 Documentation

- **[DOCKER.md](./DOCKER.md)** - Complete Docker guide
- **[PODMAN.md](./PODMAN.md)** - Podman/systemd guide
- **[TESTING.md](./TESTING.md)** - Testing procedures
- **[.env.example](./.env.example)** - Configuration template

## 🎨 Deployment Scenarios

### Development

```bash
docker compose -f docker-compose.dev.yml up
```

### Production with Nginx

See [DOCKER.md](./DOCKER.md#production-deployment) for Nginx configuration.

### Podman with systemd

```bash
# Copy quadlet files
cp quadlets/*.{container,network} ~/.config/containers/systemd/

# Enable services
systemctl --user enable --now timeful-app.service
```

See [PODMAN.md](./PODMAN.md) for complete guide.

## 🔐 Security Checklist

- [ ] Strong encryption key generated
- [ ] Google OAuth credentials secured
- [ ] .env file not committed to git
- [ ] Using HTTPS in production (reverse proxy)
- [ ] Regular database backups enabled
- [ ] Docker images kept updated

## 📊 System Requirements

**Minimum:**
- 2 CPU cores
- 2GB RAM
- 10GB disk space

**Recommended:**
- 4 CPU cores
- 4GB RAM
- 20GB disk space

## 🆘 Getting Help

- **Documentation**: See [DOCKER.md](./DOCKER.md)
- **Issues**: https://github.com/schej-it/timeful.app/issues
- **Discord**: https://discord.gg/v6raNqYxx3
- **Testing**: See [TESTING.md](./TESTING.md)

## ✅ Success Indicators

After starting, verify:

```bash
# All containers running
docker compose ps

# Application responding
curl -I http://localhost:3002/

# Database healthy
docker compose exec mongodb mongosh --eval "db.adminCommand('ping')"
```

**Expected:**
- ✅ MongoDB status: healthy
- ✅ App status: running
- ✅ HTTP response: 200 OK
- ✅ No errors in logs

## 🔄 Update Workflow

```bash
# 1. Pull latest code
git pull

# 2. Rebuild and restart
make pull

# 3. Verify
make logs
```

## 💾 Backup Workflow

```bash
# Create backup
make backup

# Backups saved to: ./backups/backup-YYYYMMDD-HHMMSS/

# To restore
make restore
```

## 🐳 Docker vs Podman

Both are supported!

**Docker:**
```bash
docker compose up -d
```

**Podman:**
```bash
podman-compose up -d
# or with quadlets (systemd):
systemctl --user start timeful-app.service
```

See [PODMAN.md](./PODMAN.md) for Podman-specific instructions.

## 📦 What's Included

- ✅ Docker Compose configuration
- ✅ Multi-stage optimized Dockerfile
- ✅ MongoDB with persistent storage
- ✅ Environment-based configuration
- ✅ Management Makefile
- ✅ Interactive setup script
- ✅ Podman Quadlets for systemd
- ✅ Comprehensive documentation
- ✅ Testing guides
- ✅ Backup/restore tools

## 🎉 You're Ready!

Pick a method above and start self-hosting Timeful in minutes!

For detailed instructions, see [DOCKER.md](./DOCKER.md).

---

**License:** AGPL v3  
**Project:** https://github.com/schej-it/timeful.app  
**Website:** https://timeful.app
