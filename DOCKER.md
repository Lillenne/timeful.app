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

#### Optional Configuration

3. **Microsoft OAuth Credentials** (Optional - for Outlook calendar integration)
   
   If you want to enable Outlook/Microsoft 365 calendar integration, you need to create a Microsoft Entra ID (formerly Azure AD) application:

   **Step-by-step setup:**
   
   a. **Create an App Registration:**
      - Go to [Azure Portal](https://portal.azure.com/)
      - Navigate to **Azure Active Directory** → **App registrations** → **New registration**
      - Enter application details:
        - **Name:** `Timeful Self-Hosted` (or any name you prefer)
        - **Supported account types:** Select "Accounts in any organizational directory and personal Microsoft accounts (Multitenant + Personal)"
        - **Redirect URI:** Select "Web" and enter: `http://localhost:3002/auth`
          - For production, use your domain: `https://yourdomain.com/auth`
      - Click **Register**

   b. **Get Client ID:**
      - After registration, you'll see the **Overview** page
      - Copy the **Application (client) ID** - this is your `MICROSOFT_CLIENT_ID`
      - Add it to your `.env` file:
        ```
        MICROSOFT_CLIENT_ID=your_microsoft_client_id_here
        ```

   c. **Create Client Secret:**
      - In the same app registration, go to **Certificates & secrets**
      - Click **New client secret**
      - Add a description (e.g., "Timeful Self-Hosted Secret")
      - Choose an expiration period (recommend: 24 months)
      - Click **Add**
      - **Important:** Copy the **Value** immediately (not the Secret ID) - this is your `MICROSOFT_CLIENT_SECRET`
      - Add it to your `.env` file:
        ```
        MICROSOFT_CLIENT_SECRET=your_microsoft_client_secret_here
        ```
      - **Note:** The secret value is only shown once. If you lose it, you'll need to create a new one.

   d. **Configure API Permissions:**
      - Go to **API permissions** in your app registration
      - Click **Add a permission** → **Microsoft Graph** → **Delegated permissions**
      - Add the following permissions:
        - `offline_access` - Maintain access to data you have given it access to
        - `User.Read` - Sign in and read user profile
        - `Calendars.Read` - Read user calendars
      - Click **Add permissions**
      - (Optional) If you're an admin, click **Grant admin consent for [Your Organization]**
        - If you're not an admin, users will be prompted to consent when they first sign in

   e. **Update config.js:**
      - Copy `config.example.js` to `config.js` in the repository root
      - Add your Microsoft Client ID to the config:
        ```javascript
        window.__TIMEFUL_CONFIG__ = {
          googleClientId: 'your_google_client_id',
          microsoftClientId: 'your_microsoft_client_id',
          // ... other settings
        }
        ```

   **Important Notes:**
   - The Microsoft OAuth credentials are separate from Google OAuth credentials
   - Both are optional - you can use just Google, just Microsoft, or both
   - If Microsoft credentials are not configured, the Outlook calendar integration button will show an error
   - For production deployments, make sure to update the redirect URI in Azure Portal to match your domain
   - Client secrets expire - set a reminder to renew them before expiration

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

## Self-Hosting Listmonk with Docker Compose

Timeful can be integrated with [Listmonk](https://listmonk.app/), a self-hosted newsletter and mailing list manager, for sending email campaigns and transactional emails.

### Why Use Listmonk?

Listmonk provides:
- **Newsletter Management**: Create and send newsletters to subscribers
- **Transactional Emails**: Send event reminders and notifications
- **Subscriber Management**: Manage email lists and preferences
- **Templates**: Create reusable email templates
- **Analytics**: Track email opens and clicks
- **Self-Hosted**: Full control over your email infrastructure

### Adding Listmonk to Your Docker Setup

The repository includes Listmonk services in both `docker-compose.yml` and `docker-compose.ghcr.yml` files. These services are **optional** and can be enabled or disabled as needed.

#### Services Included

Two additional services are defined for Listmonk:

1. **listmonk-db** - PostgreSQL database for Listmonk
   - Port: 5432 (internal only)
   - Volume: `listmonk_db_data` for data persistence

2. **listmonk** - Listmonk web application
   - Port: 9000 (default, configurable via `LISTMONK_PORT`)
   - Web interface: http://localhost:9000
   - API endpoint: http://localhost:9000/api

#### Quick Start with Listmonk

**1. Create Configuration File**

Copy the example configuration file and customize it with secure passwords:

```bash
# Copy the example configuration
cp listmonk-config-example.toml listmonk-config.toml

# Edit the configuration file
nano listmonk-config.toml
```

Update the following in `listmonk-config.toml`:
- `admin_username`: Your desired admin username
- `admin_password`: A strong password for the admin account
- `password` (in [db] section): Should match `LISTMONK_DB_PASSWORD` in your `.env` file

**Important**: Use strong passwords! Generate secure passwords with:
```bash
openssl rand -base64 32
```

**2. Configure Environment Variables**

Edit your `.env` file and set the database password:

```bash
# Listmonk Database Password
LISTMONK_DB_PASSWORD=your_secure_password_here
```

**3. Start Listmonk Services**

The Listmonk container automatically initializes the database on first run:

```bash
# Start all services including Listmonk
docker compose up -d

# Or start only Listmonk services
docker compose up -d listmonk-db listmonk
```

The database schema is automatically initialized on first startup. The initialization and upgrade commands run on every container start but are idempotent (they only make changes when needed), so subsequent startups are fast.

**4. Access Listmonk Admin Interface**

Open your browser and navigate to http://localhost:9000

Log in with the credentials you set in `listmonk-config.toml`:
- **Username**: Your configured admin username
- **Password**: Your configured admin password

**5. Configure Listmonk**

After logging in:

a. **Configure SMTP Settings** (required for sending emails):
   - Go to Settings → SMTP
   - Add your SMTP server details (e.g., Gmail, SendGrid, Mailgun, Amazon SES)
   - Test the connection

   Example for Gmail:
   ```
   Host: smtp.gmail.com
   Port: 587
   TLS: Enabled
   Username: your-email@gmail.com
   Password: your-app-password
   ```

b. **Create a Mailing List**:
   - Go to Lists → Create New
   - Name your list (e.g., "Timeful Users")
   - Note the List ID (you'll need this for the `.env` file)

c. **Create Email Templates** (required for full functionality):
   - See [LISTMONK_TEMPLATES.md](LISTMONK_TEMPLATES.md) for detailed instructions and template examples
   - You need to create 9 email templates:
     - 6 transactional notification templates (IDs: 8, 9, 10, 11, 13, 14)
     - 3 reminder email templates (configurable IDs)
   - Each template includes subject lines, HTML body, and required variables
   - Test each template after creation to ensure proper rendering

**6. Configure Timeful to Use Listmonk**

Edit your `.env` file and add/update these variables:

```env
# Listmonk Configuration
LISTMONK_ENABLED=true
LISTMONK_URL=http://listmonk:9000
LISTMONK_USERNAME=your_admin_username
LISTMONK_PASSWORD=your_admin_password
LISTMONK_LIST_ID=1

# Required: Reminder Email Template IDs
# These are the templates for scheduled reminder emails (sent immediately, after 24h, and after 72h)
# See LISTMONK_TEMPLATES.md for template examples and setup instructions
LISTMONK_INITIAL_EMAIL_REMINDER_ID=1
LISTMONK_SECOND_EMAIL_REMINDER_ID=2
LISTMONK_FINAL_EMAIL_REMINDER_ID=3

# Optional: Change Listmonk port
LISTMONK_PORT=9000
```

**Note**: The application expects certain template IDs to be hardcoded (8, 9, 10, 11, 13, 14) for transactional emails. See [LISTMONK_TEMPLATES.md](LISTMONK_TEMPLATES.md) for details on how to create templates with these IDs or how to make them configurable.

**7. Restart Services**

```bash
docker compose restart backend listmonk
```

### Reminder Email Scheduling

Timeful uses a production-ready scheduler (`robfig/cron`) for sending reminder emails instead of Google Cloud Tasks. This means:

- **No external dependencies**: All email scheduling is handled internally using a vetted cron library
- **Self-hosted friendly**: No need for Google Cloud account or service account keys
- **Reliable scheduling**: Uses `robfig/cron` (the most popular Go scheduling library) with standard cron syntax
- **Automatic scheduling**: Reminder emails are sent immediately, after 24 hours, and after 72 hours
- **Graceful cancellation**: Reminders are automatically cancelled when users respond

The backend service includes a cron-based scheduler that:
- Runs on standard cron schedule: `* * * * *` (every minute)
- Checks for pending reminder emails in MongoDB
- Sends emails at their scheduled time using Listmonk's external subscriber mode
- Marks emails as sent to prevent duplicates
- Handles user responses by cancelling remaining reminders

**Note**: If you were previously using Google Cloud Tasks (via `SERVICE_ACCOUNT_KEY_PATH`), you can remove that configuration. The application will fall back to the cron scheduler automatically. Both systems can run simultaneously for backwards compatibility during migration.

### Managing Listmonk

#### View Logs

```bash
docker compose logs -f listmonk
docker compose logs -f listmonk-db
```

#### Backup Listmonk Database

```bash
# Create backup
docker compose exec listmonk-db pg_dump -U listmonk listmonk > listmonk-backup.sql

# Restore from backup
docker compose exec -T listmonk-db psql -U listmonk listmonk < listmonk-backup.sql
```

#### Stop Listmonk (Keep Other Services Running)

```bash
docker compose stop listmonk listmonk-db
```

#### Remove Listmonk Completely

If you don't want to use Listmonk:

```bash
# Stop and remove containers
docker compose stop listmonk listmonk-db
docker compose rm -f listmonk listmonk-db

# Remove volume (⚠️  deletes all Listmonk data)
docker volume rm timeful.app_listmonk_db_data
```

Or you can simply comment out the Listmonk services in your `docker-compose.yml` file.

### Accessing Listmonk from Outside Docker

If you're running Listmonk but accessing Timeful from outside the Docker network (e.g., for development), you'll need to use `http://localhost:9000` instead of `http://listmonk:9000`:

```env
# For external access (development)
LISTMONK_URL=http://localhost:9000

# For Docker internal access (production)
LISTMONK_URL=http://listmonk:9000
```

### Production Considerations

For production deployments:

1. **Change Default Credentials**: Update both admin password and database password
2. **Use Environment Variables**: Consider using Docker secrets for sensitive data
3. **Configure SMTP**: Set up a reliable SMTP service (SendGrid, Mailgun, Amazon SES)
4. **Backup Regularly**: Set up automated backups of the PostgreSQL database
5. **Reverse Proxy**: Put Listmonk behind a reverse proxy with SSL (Nginx, Caddy, Traefik)
6. **Resource Limits**: Add memory and CPU limits in production

Example reverse proxy configuration (Nginx):

```nginx
# Listmonk admin interface
location /listmonk/ {
    proxy_pass http://localhost:9000/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

### Troubleshooting Listmonk

**Listmonk won't start**:
```bash
# Check if database is healthy
docker compose ps listmonk-db

# View logs
docker compose logs listmonk
```

**Can't access Listmonk UI**:
- Verify port 9000 is not in use: `netstat -an | grep 9000`
- Check firewall settings
- Ensure services are running: `docker compose ps`

**Database connection errors**:
- Verify PostgreSQL is healthy: `docker compose exec listmonk-db pg_isready -U listmonk`
- Check credentials in `listmonk-config.toml` match `.env` file

**Emails not sending**:
- Check SMTP configuration in Listmonk settings
- Verify SMTP credentials are correct
- Check Listmonk logs for error messages
- Test SMTP connection from within Listmonk

### Using External Listmonk Instance

If you already have a Listmonk instance running elsewhere:

1. Don't start the Listmonk services in docker-compose
2. Configure `.env` to point to your external instance:

```env
LISTMONK_URL=https://your-listmonk-domain.com
LISTMONK_USERNAME=your_username
LISTMONK_PASSWORD=your_password
LISTMONK_LIST_ID=your_list_id
```

3. Restart the Timeful backend:

```bash
docker compose restart backend
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

### Self-Hosted Premium Features

**For self-hosted deployments**, you can automatically unlock premium features for all users without requiring Stripe integration:

```env
SELF_HOSTED_PREMIUM=true
```

**What this does:**
- ✅ Automatically grants premium features to all users
- ✅ No payment processing required
- ✅ Perfect for organizations self-hosting for their teams
- ✅ No Stripe API key needed

**When to use:**
- Self-hosting for your organization or team
- Don't want to manage payments
- Want all features available immediately

**Note:** This setting defaults to `true` in the Docker Compose configurations. If you want to enable payment processing with Stripe, set `SELF_HOSTED_PREMIUM=false` and configure `STRIPE_API_KEY`.

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

The Docker images have been hardened with security best practices:

### Container Security Hardening

1. **Non-root Users**: All containers run as non-root users by default
   - Backend: Runs as `appuser` (UID 1000)
   - Frontend: Runs as `nginx` user (UID 101)
   - MongoDB: Runs as `mongodb` user (UID 999)

2. **Capability Dropping**: Containers drop all capabilities by default, only adding specific ones needed
   - Backend: Minimal capabilities (drops ALL)
   - Frontend: Only NET_BIND_SERVICE, CHOWN, SETGID, SETUID for nginx operation
   - MongoDB: Only CHOWN, SETGID, SETUID for database operations

3. **No New Privileges**: Security option `no-new-privileges:true` prevents privilege escalation

4. **Read-Only Mounts**: Configuration files are mounted read-only where possible

### Rootless Container Support

The configuration works seamlessly with both rootful and rootless container runtimes:

**Docker Rootless Mode:**
```bash
# Install Docker in rootless mode (if not already done)
dockerd-rootless-setuptool.sh install

# Use Docker normally - UID/GID mappings handled automatically
docker compose up -d
```

**Podman (Rootless by Default):**
```bash
# Podman runs rootless by default - just use it
podman-compose up -d

# Or with systemd integration
systemctl --user start timeful-backend.service
```

The fixed UIDs (1000 for backend, 101 for frontend, 999 for MongoDB) ensure consistent permissions across different environments. Volume mounts automatically work with user namespace mapping in both Docker and Podman rootless modes.

### Additional Security Best Practices

1. **Change default ports** in production
2. **Use strong encryption key** (generate with `openssl rand -base64 32`)
3. **Keep Google OAuth credentials secure** - never commit to version control
4. **Use Docker secrets** for sensitive environment variables in production
5. **Regular backups** of MongoDB data
6. **Use HTTPS** in production (via reverse proxy)
7. **Keep Docker images updated** regularly
8. **Review and audit** container logs regularly
9. **Use network segmentation** - containers communicate only through defined networks
10. **Implement resource limits** (CPU, memory) in production environments

### Security Testing

To verify the security configuration:

```bash
# Check that containers run as non-root
docker compose ps --format "{{.Name}}: {{.Image}}"
docker compose exec backend id
docker compose exec frontend id
docker compose exec mongodb id

# Verify capabilities
docker inspect timeful-backend | grep -A 20 CapDrop
docker inspect timeful-frontend | grep -A 20 CapDrop

# Check security options
docker inspect timeful-backend | grep -A 5 SecurityOpt
```

All containers should report running as non-root users with minimal capabilities.

## Support

- GitHub Issues: https://github.com/schej-it/timeful.app/issues
- Discord: https://discord.gg/v6raNqYxx3

## License

This project is licensed under the AGPL v3 License.
