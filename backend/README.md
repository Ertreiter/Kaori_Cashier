# Kaori POS Backend

A Point of Sale (POS) backend API for cafe operations with multi-store support.

## Features

- 🏪 Multi-store management
- 📱 3 order sources (Table QR, Client App, Cashier)
- 💳 Payment integration (Midtrans)
- 🔄 Real-time order updates (WebSocket)
- 📊 Daily reports
- 🎫 Voucher system
- 👥 Membership (schema ready)

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL (Supabase)
- **Auth**: JWT
- **Real-time**: WebSocket
- **Container**: Docker

## Project Structure

```
backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration loading
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # Auth, CORS, logging
│   ├── model/           # Database models
│   ├── repository/      # Database operations
│   ├── service/         # Business logic
│   └── websocket/       # Real-time hub
├── pkg/
│   ├── database/        # Database connection
│   ├── jwt/             # JWT utilities
│   └── response/        # API response helpers
├── migrations/          # SQL migrations
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Quick Start

### Using Docker (Recommended)

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your Supabase credentials

# Start with Docker
docker-compose up -d

# API available at http://localhost:8080
```

### Local Development

```bash
# Install dependencies
go mod download

# Run migrations
go run cmd/migrate/main.go up

# Start server
go run cmd/server/main.go
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | Server port (default: 8080) |
| `DATABASE_URL` | Supabase PostgreSQL connection string |
| `JWT_SECRET` | Secret key for JWT signing |
| `MIDTRANS_SERVER_KEY` | Midtrans server key |
| `MIDTRANS_CLIENT_KEY` | Midtrans client key |
| `MIDTRANS_IS_PRODUCTION` | true/false |

## API Documentation

See [API Endpoints](../docs/api.md)

## License

Private - All rights reserved
