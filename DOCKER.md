# Timeful Docker Deployment Guide

This guide will help you self-host Timeful using Docker Compose.

## Prerequisites

- Docker (version 20.10 or later)
- Docker Compose (version 2.0 or later)
- A Google Cloud account (for OAuth and Calendar integration)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/schej-it/timeful.app.git
cd timeful.app
```

### 2. Set Up Environment Variables

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

### 3. Start the Application

```bash
docker-compose up -d
```

This will:
- Pull the MongoDB image
- Build the Timeful application (frontend + backend)
- Start all services
- Create persistent volumes for data storage

### 4. Access the Application

Open your browser and navigate to:
- **Application**: http://localhost:3002
- **API Documentation**: http://localhost:3002/swagger/index.html

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

## Management Commands

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app
docker-compose logs -f mongodb
```

### Stop the Application

```bash
docker-compose down
```

### Stop and Remove All Data

```bash
docker-compose down -v
```

### Rebuild the Application

If you pull updates from the repository:

```bash
git pull
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Backup MongoDB Data

```bash
# Create backup
docker-compose exec mongodb mongodump --db=schej-it --out=/data/backup

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

For systemd integration with Podman, you can create quadlet files. See [Podman Quadlets documentation](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html) for more information.

Example quadlet setup is coming soon.

## Architecture

The Docker setup consists of:

- **MongoDB**: Database for storing events, users, and application data
- **Timeful App**: Combined frontend (Vue.js) and backend (Go) in a single container
  - Frontend: Built as static files and served by the Go server
  - Backend: Go server handling API requests and serving frontend
  - Port 3002: Main application access point

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
