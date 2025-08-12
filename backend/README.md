# Stock Analysis API Backend

A Go-based REST API for stock analyst recommendations retrieval, analysis, and investment recommendations.

## Features

- **Stock Data Retrieval**: Fetch stock analyst recommendations from KarenAI API
- **Database Storage**: Store stock information and analyst coverage in CockroachDB
- **Stock Recommendations**: Advanced algorithm to recommend best stocks based on analyst sentiment
- **REST API**: Clean RESTful endpoints for frontend integration
- **Real-time Sync**: Sync all stock data from external API on demand

## Tech Stack

- **Backend**: Go 1.19+
- **Database**: CockroachDB (PostgreSQL compatible)
- **External API**: KarenAI Stock Challenge API
- **HTTP Router**: Gorilla Mux
- **CORS**: rs/cors

## Setup

### Prerequisites

1. Go 1.19 or higher
2. CockroachDB instance running
3. KarenAI API token

### Installation

1. Copy environment variables:
   ```bash
   cp .env.example .env
   ```

2. Update `.env` with your configuration (token is pre-configured):
   ```
   DATABASE_URL=postgresql://root@localhost:26257/stockdb?sslmode=disable
   KAREN_AI_TOKEN=your_api_token_here
   PORT=8080
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. Sync initial stock data:
   ```bash
   curl -X POST http://localhost:8080/api/v1/stocks/sync
   ```

### CockroachDB Setup

1. Start CockroachDB:
   ```bash
   cockroach start-single-node --insecure --listen-addr=localhost:26257
   ```

2. Create database:
   ```bash
   cockroach sql --insecure --execute="CREATE DATABASE stockdb;"
   ```

## API Endpoints

### Health Check
- `GET /api/v1/health` - Check API health

### Stocks
- `GET /api/v1/stocks` - Get all stocks with latest analyst coverage
- `GET /api/v1/stocks/{symbol}` - Get specific stock by symbol with analysis history
- `POST /api/v1/stocks/sync` - Sync all stocks from KarenAI API (recommended first step)
- `POST /api/v1/stocks/{symbol}/refresh` - Refresh specific stock data
- `GET /api/v1/stocks/search/{symbol}` - Search for existing stock

### Recommendations
- `GET /api/v1/stocks/recommendations` - Get top stock recommendations based on analyst sentiment

## Response Format

All endpoints return responses in this format:
```json
{
  "success": true,
  "data": {...},
  "error": null
}
```

## Database Schema

### Tables

1. **stocks** - Basic stock information (symbol, company name)
2. **stock_analysis** - Analyst recommendations and target price changes

## Recommendation Algorithm

The recommendation engine analyzes stocks based on analyst sentiment:

- **Rating Quality**: Buy/Strong Buy ratings score highest (80 points)
- **Rating Changes**: Upgrades add 15 points, downgrades subtract 10
- **Target Price Changes**: Raises >10% add 20 points, >5% add 10 points
- **Action Types**: Initiations add 10 points, raises add 12 points
- **Coverage Consistency**: Multiple recent positive analyses add 8 points
- **Recent Activity**: Stocks with 3+ recent analyses get 5 point bonus

Stocks are scored 0-100 and ranked by total score. Top 10 recommendations are returned.

## Data Source

Stock data comes from KarenAI API which provides:
- Company ticker symbols and names
- Analyst ratings (Buy, Outperform, Hold, etc.)
- Price targets and target changes
- Brokerage firm recommendations
- Analysis timestamps

## Development

### Project Structure
```
backend/
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routes
│   ├── clients/           # KarenAI API client
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and migration
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   └── services/          # Business logic and recommendation engine
├── Makefile               # Development commands
├── Dockerfile             # Container configuration
└── README.md
```

### Quick Start Commands

```bash
# Setup development environment
make setup

# Run the application
make run

# Sync stock data (run after first startup)
curl -X POST http://localhost:8080/api/v1/stocks/sync

# Get recommendations
curl http://localhost:8080/api/v1/stocks/recommendations

# Build for production
make build-prod
```

### Adding New Features

1. Add models in `internal/models/`
2. Create repository methods in `internal/repository/`
3. Implement business logic in `internal/services/`
4. Add API endpoints in `internal/api/`
5. Update database schema in `internal/database/`