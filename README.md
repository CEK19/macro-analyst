# Macro Analyst

Real-time cryptocurrency market data streaming from **Binance** and macroeconomic data from **FRED API** (Federal Reserve Economic Data).

## Features

### Cryptocurrency (Binance WebSocket)
- 🚀 **Real-time Data**: Live prices from Binance using `adshao/go-binance` SDK
- 📊 **Multi-Symbol Tracking**: BTC, ETH, BNB, SOL, ADA, XRP (all vs USDT)
- 🔄 **Auto-Reconnect**: Automatic reconnection on connection loss
- ⚡ **Throttling**: Smart throttling to prevent overwhelming clients
- 🧵 **Concurrent**: Hub-and-Spoke pattern with goroutine-safe broadcasting

### Macroeconomic Data (FRED API)
- 🏛️ **Federal Reserve Data**: Access to official US economic indicators
- 📈 **Key Metrics**: Fed assets, interest rates, inflation, USD index
- 🔌 **RESTful API**: Clean HTTP endpoints for data access
- ✅ **95.3% Test Coverage**: Production-ready with comprehensive tests
- 🎯 **Clean Architecture**: Interface-based, testable, maintainable code

### General
- 🛡️ **Production-Ready**: Graceful shutdown, context cancellation, comprehensive tests
- 📝 **Well Documented**: Detailed API documentation and examples
- 🧪 **Highly Tested**: 70%+ coverage across all packages

## Architecture

```
┌─────────────┐          ┌─────────────┐
│  Binance    │          │  FRED API   │
│  WebSocket  │          │ (St. Louis  │
└──────┬──────┘          │     Fed)    │
       │                 └──────┬──────┘
       ▼                        ▼
┌─────────────┐          ┌─────────────┐
│  Ingestor   │          │ FRED Client │
└──────┬──────┘          └──────┬──────┘
       │                        │
       ▼                        ▼
┌─────────────┐          ┌─────────────┐
│     Hub     │          │   HTTP API  │
└──────┬──────┘          └─────────────┘
       │
       ├───────┬───────┬───────┐
       ▼       ▼       ▼       ▼
   Client  Client  Client  Client
```

## Quick Start

### 1. Setup Environment

```bash
# Copy example env file
cp .env.example .env

# Edit .env and add your FRED API key
# Get free key at: https://fred.stlouisfed.org/docs/api/api_key.html
```

### 2. Run Server

```bash
# Install dependencies
go mod download

# Run tests
make test

# Start server
make run
```

Server starts on `http://localhost:8080`

## API Endpoints

### WebSocket (Cryptocurrency)
- `ws://localhost:8080/ws/prices` - Real-time crypto prices

### HTTP (General)
- `GET /` - API information
- `GET /health` - Health check with active client count

### HTTP (FRED Macroeconomic Data)
- `GET /api/v1/fred/tickers` - List all available tickers
- `GET /api/v1/fred/latest` - Get all latest values
- `GET /api/v1/fred/latest/:symbol` - Get latest value for specific ticker
- `GET /api/v1/fred/ticker/:symbol` - Get historical data

**Supported Tickers:**
| Symbol | Description |
|--------|-------------|
| `WALCL` | Federal Reserve Total Assets (QE/QT) |
| `WTREGEN` | Treasury General Account |
| `RRPONTSYD` | Overnight Reverse Repo |
| `FEDFUNDS` | Federal Funds Rate |
| `CPIAUCSL` | Consumer Price Index (Inflation) |
| `DTWEXBGS` | US Dollar Index |

### Example API Calls

```bash
# Get all supported tickers
curl http://localhost:8080/api/v1/fred/tickers

# Get latest Fed assets value
curl http://localhost:8080/api/v1/fred/latest/WALCL

# Get all latest macroeconomic data
curl http://localhost:8080/api/v1/fred/latest

# Get historical CPI data (last 10 months)
curl "http://localhost:8080/api/v1/fred/ticker/CPIAUCSL?limit=10"
```

**Response Example:**
```json
{
  "ticker": "WALCL",
  "description": "Federal Reserve Total Assets",
  "value": "7500000000000",
  "date": "2024-02-15",
  "updated_at": "2024-02-20T10:00:00Z"
}
```

## Documentation

- 📚 [FRED API Guide](docs/FRED_API.md) - Comprehensive FRED integration documentation
- 🧪 Test coverage: **95.3%** (FRED), **70.7%** (WebSocket)

## Testing

```bash
# Run all tests with coverage
make test

# Generate HTML coverage report
make test-coverage

# Test specific package
go test ./internal/fred/... -v
go test ./internal/ws/... -v
```

## Test WebSocket

Open `test-ws-client.html` in browser and click Connect.

### Expected Data Format

**Single Symbol Update:**
```json
{
  "symbol": "BTCUSDT",
  "price": 94250.50,
  "change": 125.30,
  "changePercent": 0.13,
  "volume": 15234567890,
  "timestamp": "14:23:45.123"
}
```

**Multi-Symbol Update:**
```json
{
  "type": "multi_update",
  "data": [
    {
      "symbol": "BTCUSDT",
      "price": 94250.50,
      "change": 125.30,
      "changePercent": 0.13,
      "volume": 15234567890,
      "timestamp": "14:23:45.123"
    },
    {
      "symbol": "ETHUSDT",
      "price": 2635.80,
      ...
    }
  ]
}
```

## Environment

Create `.env` file:
```
PORT=8080
APP_ENV=local
```

## Configuration

Customize the Ingestor in `cmd/api/main.go`:

```go
// Change throttle interval
ingestor := ws.NewIngestor(hub,
    ws.WithThrottleInterval(1 * time.Second),
)

// Add more symbols
ingestor.AddSymbol("DOGEUSDT")
```

## Dependencies

- **gofiber/fiber** - Fast HTTP framework
- **adshao/go-binance** - Official Binance SDK
- **gofiber/contrib/websocket** - WebSocket support

## Performance

- **Throttle Interval**: 500ms (configurable)
- **Update Rate**: ~10 updates/second (6 symbols)
- **Auto-Reconnect**: Built-in with exponential backoff
- **Graceful Shutdown**: Clean disconnection on SIGINT/SIGTERM
