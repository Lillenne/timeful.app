<div align="center">
  
<img src="./.github/assets/images/logo.svg" width="200px" alt="Timeful logo" />

</div>
<br />
<div align="center">

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-orange.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Donate](https://img.shields.io/badge/-Donate%20with%20Paypal-blue?logo=paypal)](https://www.paypal.com/donate/?hosted_button_id=KWCH6LGJCP6E6)
[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/timeful_app?label=%40timeful_app&labelColor=white)](https://x.com/timeful_app)
[![Discord](https://img.shields.io/badge/-Join%20Discord-7289DA?logo=discord&logoColor=white)](https://discord.gg/v6raNqYxx3)
[![Subreddit subscribers](https://img.shields.io/reddit/subreddit-subscribers/schej?label=join%20r%2Fschej)](https://www.reddit.com/r/schej/)

</div>

<img src="./.github/assets/images/hero.jpg" alt="Timeful hero" />

Timeful is a scheduling platform helps you find the best time for a group to meet. It is a free availability poll that is easy to use and integrates with your calendar.

Hosted version of the site: https://timeful.app

Built with [Vue 2](https://github.com/vuejs/vue), [MongoDB](https://github.com/mongodb/mongo), [Go](https://github.com/golang/go), and [TailwindCSS](https://github.com/tailwindlabs/tailwindcss)

## Demo

[![demo video](http://markdown-videos-api.jorgenkh.no/youtube/vFkBC8BrkOk)](https://www.youtube.com/watch?v=vFkBC8BrkOk)

## Features

- See when everybody's availability overlaps
- Easily specify date + time ranges to meet between
- Google calendar, Outlook, Apple calendar integration
- "Available" vs. "If needed" times
- Determine when a subset of people are available
- Schedule across different time zones
- Email notifications + reminders
- Duplicating polls
- Availability groups - stay up to date with people's real-time calendar availability
- Export availability as CSV
- Only show responses to event creator

## Self-hosting

Timeful can be easily self-hosted using Docker Compose! 🐳

### Quick Start with Pre-built Images (Recommended)

The fastest way to get started:

1. Install [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
2. Pull the pre-built images:
   ```bash
   docker pull ghcr.io/lillenne/timeful.app/backend:latest
   docker pull ghcr.io/lillenne/timeful.app/frontend:latest
   ```
3. Clone the repository: `git clone https://github.com/schej-it/timeful.app.git`
4. Copy the example environment file: `cp .env.example .env`
5. Configure your `.env` file (minimum: `ENCRYPTION_KEY`)
6. Start the application: `docker compose -f docker-compose.ghcr.yml up -d` or `make up-ghcr`
7. Access at http://localhost:3002

### Building from Source

If you prefer to build the images yourself:

1. Install [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
2. Clone the repository: `git clone https://github.com/schej-it/timeful.app.git`
3. Copy the example environment file: `cp .env.example .env`
4. Configure your `.env` file with at minimum:
   - Encryption key (generate with: `openssl rand -base64 32`)
5. Start the application: `docker compose up -d`
6. Access at http://localhost:3002

### Available Images

Pre-built Docker images are automatically published to GitHub Container Registry:
- **Backend**: `ghcr.io/lillenne/timeful.app/backend:latest`
- **Frontend**: `ghcr.io/lillenne/timeful.app/frontend:latest`

Images are built for both `linux/amd64` and `linux/arm64` platforms.

### Full Documentation

See [DOCKER.md](./DOCKER.md) for complete self-hosting instructions, including:
- Detailed setup guide for both pre-built and source builds
- Production deployment with reverse proxy (Nginx/Caddy)
- Custom domain and CORS configuration
- Optional features configuration (Google OAuth, Stripe, Email, etc.)
- **Premium features automatically unlocked** for self-hosted deployments
- Backup and maintenance procedures
- Troubleshooting tips
- Podman support with Quadlets

### Requirements

- Docker (version 20.10+) or Podman
- **Minimum**: Just an encryption key for anonymous event scheduling
- **Optional**: Google Cloud account for user accounts and calendar integration
- 2GB+ RAM recommended
- 10GB+ disk space for MongoDB data

### Self-Hosted Premium Features

By default, all premium features are **automatically unlocked** for self-hosted deployments! No payment processing or Stripe configuration required. Perfect for organizations wanting to run Timeful for their teams.
