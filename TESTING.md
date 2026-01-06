# Testing Timeful Docker Deployment

This guide will help you test the Docker deployment to ensure everything works correctly.

## Prerequisites

- Docker 20.10+ and Docker Compose v2
- At least 4GB free disk space
- Google OAuth credentials (see DOCKER.md for setup)

## Quick Test (Development Mode)

For a quick test without full OAuth setup:

```bash
# Use development compose file (minimal config required)
docker compose -f docker-compose.dev.yml up --build
```

This will:
- Start MongoDB
- Build and start the application
- Use a development encryption key
- Run in debug mode

Access at http://localhost:3002

**Note**: Some features requiring OAuth (calendar integration) won't work without proper credentials.

## Full Production Test

### 1. Setup Environment

```bash
# Create .env from template
cp .env.example .env

# Generate encryption key
echo "ENCRYPTION_KEY=$(openssl rand -base64 32)" >> .env

# Add your Google OAuth credentials
nano .env  # or your preferred editor
```

### 2. Build the Image

```bash
# Build without cache for clean build
docker compose build --no-cache
```

Expected output:
- ✅ Frontend build completes successfully
- ✅ Backend build completes successfully
- ✅ Final image is created

### 3. Start Services

```bash
docker compose up -d
```

Wait for services to start (~30 seconds):

```bash
# Check status
docker compose ps

# Should show:
# - timeful-mongodb: healthy
# - timeful-app: running
```

### 4. Check Logs

```bash
# View all logs
docker compose logs

# Follow app logs
docker compose logs -f app

# Check for startup messages:
# ✅ "Server Restarted"
# ✅ "Connected to MongoDB"
# ✅ Server listening on :3002
```

### 5. Test Application

#### A. Access Web Interface

Open http://localhost:3002 in your browser.

**Expected**:
- ✅ Homepage loads
- ✅ No console errors
- ✅ Can navigate to create event page

#### B. Test API Endpoint

```bash
# Health check (should return HTML)
curl -I http://localhost:3002/

# API documentation
curl http://localhost:3002/swagger/index.html
```

#### C. Test Database Connection

```bash
# Connect to MongoDB
docker compose exec mongodb mongosh schej-it

# In mongo shell:
# > show collections
# > db.events.countDocuments()
# > exit
```

### 6. Test OAuth (if configured)

1. Navigate to http://localhost:3002
2. Click "Sign In with Google"
3. Should redirect to Google OAuth
4. After auth, should redirect back to app

**Note**: Make sure your redirect URI in Google Console matches:
- Local: `http://localhost:3002/api/auth/google/callback`
- Production: `https://yourdomain.com/api/auth/google/callback`

### 7. Test Event Creation

1. Create a new event
2. Add some availability
3. Check if data persists:

```bash
# Query events in database
docker compose exec mongodb mongosh schej-it --eval "db.events.find().pretty()"
```

### 8. Performance Test

```bash
# Simple load test (requires Apache Bench)
ab -n 100 -c 10 http://localhost:3002/

# Expected:
# - No failed requests
# - Reasonable response time (<500ms for simple requests)
```

## Testing Different Scenarios

### With External MongoDB

Update `.env`:
```env
MONGO_URI=mongodb://your-external-host:27017
```

Then restart:
```bash
docker compose restart app
```

### With Custom Port

Edit `docker-compose.yml`:
```yaml
services:
  app:
    ports:
      - "8080:3002"  # Change 8080 to desired port
```

### Behind Reverse Proxy

See DOCKER.md for Nginx/Caddy configuration examples.

## Troubleshooting Tests

### Application Won't Start

```bash
# Check logs
docker compose logs app

# Common issues:
# - Missing environment variables
# - MongoDB not ready (wait 30s and retry)
# - Port 3002 already in use
```

### Database Connection Failed

```bash
# Test MongoDB connectivity
docker compose exec app ping -c 3 mongodb

# Check MongoDB logs
docker compose logs mongodb

# Verify environment
docker compose exec app env | grep MONGO
```

### Frontend Not Loading

```bash
# Check if frontend files exist in container
docker compose exec app ls -la /app/frontend/dist

# Should see index.html and asset files
```

### OAuth Not Working

Common issues:
1. Wrong redirect URI in Google Console
2. Missing environment variables (CLIENT_ID, CLIENT_SECRET)
3. Incorrect domain in CORS settings

```bash
# Verify env vars are set
docker compose exec app env | grep CLIENT
```

## Automated Testing Script

Create `test-deployment.sh`:

```bash
#!/bin/bash

echo "Starting Timeful deployment test..."

# Start services
docker compose up -d

# Wait for services
echo "Waiting for services to start..."
sleep 30

# Check if app is responding
if curl -f http://localhost:3002/ > /dev/null 2>&1; then
    echo "✅ Application is responding"
else
    echo "❌ Application is not responding"
    docker compose logs app
    exit 1
fi

# Check MongoDB
if docker compose exec mongodb mongosh schej-it --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    echo "✅ MongoDB is healthy"
else
    echo "❌ MongoDB is not healthy"
    docker compose logs mongodb
    exit 1
fi

echo "✅ All tests passed!"
echo "Access the application at: http://localhost:3002"
```

Make it executable and run:
```bash
chmod +x test-deployment.sh
./test-deployment.sh
```

## Cleanup After Testing

```bash
# Stop and remove containers
docker compose down

# Remove volumes (deletes data)
docker compose down -v

# Remove images
docker rmi timeful-app:latest
docker rmi mongo:6.0
```

## Performance Benchmarks

Expected performance on modern hardware:

- **Startup time**: < 60 seconds
- **Homepage load**: < 1 second
- **API response**: < 200ms (simple queries)
- **Memory usage**: 
  - MongoDB: ~200MB
  - App: ~50MB
- **Disk space**: ~2GB (with data)

## Success Criteria

✅ All containers start successfully
✅ No errors in logs
✅ Homepage loads without errors
✅ API documentation accessible
✅ Database queries work
✅ OAuth flow completes (if configured)
✅ Events can be created and retrieved
✅ Data persists after restart

## Reporting Issues

If you encounter issues during testing:

1. Collect logs: `docker compose logs > logs.txt`
2. Check versions: `docker version && docker compose version`
3. Share your setup (sanitized .env, docker-compose.yml)
4. Report on GitHub: https://github.com/schej-it/timeful.app/issues

Include:
- Operating system
- Docker version
- Error messages
- Steps to reproduce
