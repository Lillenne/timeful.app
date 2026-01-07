# Timeful Docker Deployment Guide

This guide will help you self-host Timeful using Docker Compose.

## Prerequisites

- Docker (version 20.10 or later)
- Docker Compose (version 2.0 or later)
- A Google Cloud account (for OAuth and Calendar integration - optional)

## Deployment Options

Timeful can be deployed in two ways:

1. **Using Pre-built Images from GitHub Container Registry** (Recommended) - Faster, no build required
2. **Building from Source** - For development or customization

## Quick Start with Pre-built Images

### 1. Pull Pre-built Images

The easiest and fastest way to deploy Timeful:

```bash
# Download the pre-built images
docker pull ghcr.io/lillenne/timeful.app/backend:latest
docker pull ghcr.io/lillenne/timeful.app/frontend:latest

# Clone the repository (only for configuration files)
git clone https://github.com/schej-it/timeful.app.git
cd timeful.app
```

### 2. Configure Environment

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and set at minimum:
# ENCRYPTION_KEY=$(openssl rand -base64 32)
nano .env
```

### 3. Start with Pre-built Images

```bash
# Use the GHCR compose file
docker compose -f docker-compose.ghcr.yml up -d

# Or use the Makefile
make up-ghcr
```

The application will be available at http://localhost:3002

### 4. Update to Latest Version

```bash
# Pull latest images
docker compose -f docker-compose.ghcr.yml pull

# Restart services
docker compose -f docker-compose.ghcr.yml up -d
```

## Quick Start - Building from Source

### 1. Clone the Repository

```bash
git clone https://github.com/schej-it/timeful.app.git
cd timeful.app
```

### 2. Run the Setup Script (Recommended)

The easiest way to get started:

```bash
./docker-setup.sh
```

This script will:
- Check if Docker and Docker Compose are installed
- Create a `.env` file from the template
- Guide you through the configuration
- Start the application

**Or use the Makefile:**

```bash
make setup
```

### 3. Manual Setup (Alternative)

If you prefer to set up manually:

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` and configure the required settings:

#### Required Configuration

1. **Google OAuth Credentials** (Required for calendar integration)
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select an existing one
   - Enable the following APIs:
     - Google Calendar API
     - Google People API
   - Create OAuth 2.0 credentials (Web application type)
   - Add authorized redirect URI: `http://localhost:3002/api/auth/google/callback`
     - For production, use your domain: `https://yourdomain.com/api/auth/google/callback`
   - Required OAuth scopes:
     - `https://www.googleapis.com/auth/calendar.events.readonly`
     - `https://www.googleapis.com/auth/calendar.calendarlist.readonly`
     - `https://www.googleapis.com/auth/contacts.readonly`
   - Copy the Client ID and Client Secret to your `.env` file:
     ```
     CLIENT_ID=your_client_id_here
     CLIENT_SECRET=your_client_secret_here
     ```

2. **Encryption Key** (Required)
   - Generate a secure encryption key:
     ```bash
     openssl rand -base64 32
     ```
   - Add it to your `.env` file:
     ```
     ENCRYPTION_KEY=your_generated_key_here
     ```

### 4. Start the Application

**Using the Makefile (recommended):**

```bash
make up
```

**Or using Docker Compose directly:**

```bash
docker compose up -d
```

This will:
- Pull the MongoDB image
- Build the Timeful application (frontend + backend)
- Start all services
- Create persistent volumes for data storage

### 5. Access the Application

Open your browser and navigate to:
- **Application**: http://localhost:3002
- **API Documentation**: http://localhost:3002/swagger/index.html

## Makefile Commands

For easier management, use the included Makefile:

```bash
make help          # Show all available commands
make setup         # Run initial setup
make up            # Start the application
make down          # Stop the application
make restart       # Restart the application
make logs          # View logs (all services)
make logs-app      # View app logs only
make logs-db       # View database logs only
make backup        # Backup MongoDB database
make restore       # Restore from latest backup
make pull          # Pull updates and restart
make status        # Show container status
make clean         # Stop and remove everything (⚠️  deletes data!)
```

## Optional Features

### Email Notifications

To enable email notifications, configure Gmail in your `.env`:

```env
GMAIL_APP_PASSWORD=your_app_password
SCHEJ_EMAIL_ADDRESS=your_email@gmail.com
```

**Note**: You need to create a Gmail App Password:
1. Enable 2-factor authentication on your Google account
2. Go to https://myaccount.google.com/apppasswords
3. Create a new app password
4. Use this password in your `.env` file

### Email Campaigns (Listmonk)

If you want to use Listmonk for email campaigns:

```env
LISTMONK_URL=http://your-listmonk-instance
LISTMONK_USERNAME=admin
LISTMONK_PASSWORD=your_password
LISTMONK_LIST_ID=1
```

### Slack Notifications

For monitoring and alerts via Slack:

```env
SLACK_PROD_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

### Payment Processing (Stripe)

To enable premium features with payments:

```env
STRIPE_API_KEY=sk_live_your_stripe_key
```

## Docker Images

### Pre-built Images on GitHub Container Registry

Docker images are automatically built and published to GitHub Container Registry (GHCR) on every push to the `main` branch:

- **Backend**: `ghcr.io/lillenne/timeful.app/backend:latest`
- **Frontend**: `ghcr.io/lillenne/timeful.app/frontend:latest`

**Benefits of using pre-built images:**
- ✅ No build time - start immediately
- ✅ Automatically updated with latest changes
- ✅ Multi-platform support (amd64 and arm64)
- ✅ Smaller download size (layers cached)
- ✅ No need to clone the entire repository

**Available tags:**
- `latest` - Latest stable build from main branch
- `main` - Alias for latest
- `<version>` - Specific version tags (when releases are published)
- `<branch>-<sha>` - Specific commit builds

### Using Pre-built Images

To use pre-built images instead of building from source:

```bash
# Pull the latest images
docker compose -f docker-compose.ghcr.yml pull

# Start the application
docker compose -f docker-compose.ghcr.yml up -d
```

Or use the Makefile:

```bash
make up-ghcr
```

### Updating to Latest Version

With pre-built images, updates are simple:

```bash
# Pull latest images and restart
make pull-ghcr

# Or manually:
docker compose -f docker-compose.ghcr.yml pull
docker compose -f docker-compose.ghcr.yml up -d
```

### Building from Source

If you want to build the images yourself (for development or customization):

```bash
# Build images
docker compose build

# Start with built images
docker compose up -d
```

## Production Deployment

### Using a Reverse Proxy

For production, use a reverse proxy like Nginx or Caddy with SSL:

#### Example Nginx Configuration

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:3002;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

#### Example Caddy Configuration

```caddyfile
yourdomain.com {
    reverse_proxy localhost:3002
}
```

### Docker Compose with Reverse Proxy

You can also add Caddy to your `docker-compose.yml`:

```yaml
  caddy:
    image: caddy:2-alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
      - caddy_config:/config
    networks:
      - timeful-network
```

### Configuring for Custom Domains

If you're using a custom domain or reverse proxy, you need to configure both the base URL and CORS:

#### 1. Set the Base URL

**IMPORTANT**: This is required for Google OAuth to work with custom domains.

Edit your `.env` file and set `BASE_URL` to your domain:

```bash
BASE_URL=https://yourdomain.com
```

For local Docker development:
```bash
BASE_URL=http://localhost:3002
```

The `BASE_URL` is used for:
- OAuth redirect URIs (e.g., `{BASE_URL}/api/auth/google/callback`)
- Email links and event URLs
- Stripe payment redirects

#### 2. Configure CORS

Add your domain to the CORS allowed origins in `.env`:

```bash
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

For multiple domains (separate with commas, no spaces):
```bash
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com,https://staging.yourdomain.com
```

**Default behavior**:
- If `CORS_ALLOWED_ORIGINS` is not set: Uses default domains (schej.it, timeful.app) plus localhost
- If `CORS_ALLOWED_ORIGINS` is set: Uses your custom domains plus localhost (replaces default domains)

**Important**: 
- Always use the full URL including protocol (`https://`)
- Don't include trailing slashes
- Localhost origins (`:3002`, `:8080`) are always allowed automatically

#### 3. Update Google OAuth Credentials

**Critical**: Your Google OAuth credentials must match your BASE_URL.

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Navigate to APIs & Services > Credentials
3. Edit your OAuth 2.0 Client ID
4. Add authorized redirect URI: `{BASE_URL}/api/auth/google/callback`
   - Example: `https://yourdomain.com/api/auth/google/callback`
   - For Docker: `http://localhost:3002/api/auth/google/callback`

#### 4. Complete Example Configuration

For a custom domain deployment:

```bash
# .env file
ENCRYPTION_KEY=your_encryption_key_here
BASE_URL=https://yourdomain.com
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
CLIENT_ID=your_google_client_id
CLIENT_SECRET=your_google_client_secret
```

For Docker local development:

```bash
# .env file
ENCRYPTION_KEY=your_encryption_key_here
BASE_URL=http://localhost:3002
CLIENT_ID=your_google_client_id
CLIENT_SECRET=your_google_client_secret
```

#### 5. Restart the Application

After changing these settings:

```bash
docker compose restart backend
```

Or with GHCR images:

```bash
docker compose -f docker-compose.ghcr.yml restart backend
```

## Management Commands

### Using Makefile (Recommended)

```bash
make logs          # View all logs
make logs-app      # View app logs only
make logs-db       # View database logs only
make restart       # Restart services
make down          # Stop services
make status        # Check container status
make backup        # Create database backup
make restore       # Restore from backup
```

### Using Docker Compose Directly

#### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f app
docker compose logs -f mongodb
```

### Stop the Application

```bash
make down
# or
docker compose down
```

### Stop and Remove All Data

```bash
make clean
# or
docker compose down -v
```

### Rebuild the Application

If you pull updates from the repository:

```bash
make pull
# or manually:
git pull
docker compose down
docker compose build --no-cache
docker compose up -d
```

### Backup MongoDB Data

```bash
make backup
# or manually:
docker compose exec mongodb mongodump --db=schej-it --out=/data/backup

# Copy backup to host
docker cp timeful-mongodb:/data/backup ./mongodb-backup

# Restore from backup
docker cp ./mongodb-backup timeful-mongodb:/data/backup
docker-compose exec mongodb mongorestore --db=schej-it /data/backup/schej-it --drop
```

## Updating

To update to the latest version:

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Troubleshooting

### Application won't start

1. Check logs:
   ```bash
   docker-compose logs app
   ```

2. Verify environment variables are set correctly in `.env`

3. Ensure MongoDB is healthy:
   ```bash
   docker-compose ps mongodb
   ```

### Google OAuth errors

1. Verify your OAuth credentials in Google Cloud Console
2. Check that redirect URIs match your domain
3. Ensure all required APIs are enabled

### Port conflicts

If port 3002 is already in use, modify `docker-compose.yml`:

```yaml
services:
  app:
    ports:
      - "8080:3002"  # Change 8080 to any available port
```

### MongoDB connection issues

Check MongoDB logs:
```bash
docker-compose logs mongodb
```

Verify the connection string in your `.env` is correct.

## Using Podman Instead of Docker

Timeful can also be deployed using Podman and Podman Compose:

```bash
# Install podman-compose
pip3 install podman-compose

# Use podman-compose instead of docker-compose
podman-compose up -d
podman-compose logs -f
podman-compose down
```

### Podman Quadlets (systemd integration)

For systemd integration with Podman, you can use Quadlets to run Timeful as system services.

**Complete Quadlet setup is available in [PODMAN.md](./PODMAN.md)**

Quick start:
```bash
# Copy quadlet files
mkdir -p ~/.config/containers/systemd
cp quadlets/*.{container,network} ~/.config/containers/systemd/

# Build and enable
podman build -t localhost/timeful:latest .
systemctl --user daemon-reload
systemctl --user enable --now timeful-mongodb.service timeful-app.service
```

See [quadlets/README.md](./quadlets/README.md) for example files and [PODMAN.md](./PODMAN.md) for complete documentation.

## Architecture

The Docker setup consists of three separate containers:

- **MongoDB**: Database for storing events, users, and application data
  - Port 27017 (internal)
  
- **Backend (Go API Server)**: Handles API requests and business logic
  - Port 3002 (internal only, not exposed to host)
  - Connects to MongoDB
  
- **Frontend (Nginx + Vue.js)**: Serves static files and proxies API requests
  - Port 80 → Host port 3002
  - Proxies `/api` requests to backend container
  - Proxies `/sockets/` for WebSocket connections
  - Serves Vue.js static files with caching

**Benefits of separation:**
- Independent scaling of frontend and backend
- Backend not directly exposed to internet
- Nginx efficiently serves static files
- Standard microservices architecture

## Security Considerations

1. **Change default ports** in production
2. **Use strong encryption key** (generate with `openssl rand -base64 32`)
3. **Keep Google OAuth credentials secure** - never commit to version control
4. **Use Docker secrets** for sensitive environment variables in production
5. **Regular backups** of MongoDB data
6. **Use HTTPS** in production (via reverse proxy)
7. **Keep Docker images updated** regularly

## Support

- GitHub Issues: https://github.com/schej-it/timeful.app/issues
- Discord: https://discord.gg/v6raNqYxx3

## License

This project is licensed under the AGPL v3 License.
