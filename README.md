# Macro Analyst

Real-time cryptocurrency market data streaming from **Binance** via WebSocket.

## Features

- 🚀 **Real-time Data**: Live prices from Binance using `adshao/go-binance` SDK
- 📊 **Multi-Symbol Tracking**: BTC, ETH, BNB, SOL, ADA, XRP (all vs USDT)
- 🔄 **Auto-Reconnect**: Automatic reconnection on connection loss
- ⚡ **Throttling**: Smart throttling to prevent overwhelming clients
- 🧵 **Concurrent**: Hub-and-Spoke pattern with goroutine-safe broadcasting
- 🛡️ **Production-Ready**: Graceful shutdown, context cancellation, comprehensive tests

## Architecture

```
┌─────────────┐
│  Binance    │  ← Real-time WebSocket data
│  WebSocket  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Ingestor   │  ← Data ingestion + throttling (500ms)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│     Hub     │  ← Central message broker
└──────┬──────┘
       │
       ├───────┬───────┬───────┐
       ▼       ▼       ▼       ▼
   Client  Client  Client  Client  ← WebSocket connections
```

## Run Server

```bash
# Install dependencies
go mod download

# Start with hot reload
air

# Or direct run
go run cmd/api/main.go
```

Server starts on `http://localhost:8080`

## Endpoints

- `ws://localhost:8080/ws/prices` - WebSocket for real-time prices
- `http://localhost:8080/health` - Health check
- `http://localhost:8080/` - API info

## Test

```bash
# Run all tests
go test ./...

# With coverage
go test ./... -cover

# Specific package
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
