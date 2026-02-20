# FRED API Integration

## Overview

The FRED (Federal Reserve Economic Data) API integration provides access to critical macroeconomic data from the Federal Reserve Bank of St. Louis. This is a clean, testable, and production-ready implementation following senior-level best practices.

## Architecture

### Clean Code Principles Applied

1. **Interface-based Design**: `Client` interface allows easy mocking for tests
2. **Dependency Injection**: HTTP client injected for testability
3. **Single Responsibility**: Each function has one clear purpose
4. **DRY (Don't Repeat Yourself)**: Shared logic extracted into helper methods
5. **Error Handling**: Consistent error wrapping with context
6. **Context Support**: All API calls respect context cancellation/timeout

### Package Structure

```
internal/fred/
├── constants.go       # Ticker definitions and descriptions
├── models.go         # Data structures and JSON models
├── client.go         # HTTP client implementation
├── constants_test.go # Tests for constants
├── models_test.go    # Tests for data models
└── client_test.go    # Tests for client (95.3% coverage)
```

## Supported Macroeconomic Indicators

### The Global Foundation Tickers

| Ticker | Description | What It Measures |
|--------|-------------|------------------|
| `WALCL` | Federal Reserve Total Assets | Money printing/Quantitative tightening |
| `WTREGEN` | Treasury General Account | US Treasury cash balance |
| `RRPONTSYD` | Overnight Reverse Repo | Money withdrawn from the system |
| `FEDFUNDS` | Federal Funds Rate | Fed's operating interest rate |
| `CPIAUCSL` | Consumer Price Index | Inflation rate |
| `DTWEXBGS` | US Dollar Index | Strength of US dollar |

## API Endpoints

### 1. Get All Tickers
```http
GET /api/v1/fred/tickers
```

**Response:**
```json
{
  "tickers": [
    {
      "symbol": "WALCL",
      "description": "Federal Reserve Total Assets"
    },
    ...
  ],
  "count": 6
}
```

### 2. Get Latest Value for Specific Ticker
```http
GET /api/v1/fred/latest/:symbol
```

**Example:**
```bash
curl http://localhost:8080/api/v1/fred/latest/WALCL
```

**Response:**
```json
{
  "ticker": "WALCL",
  "description": "Federal Reserve Total Assets",
  "value": "7500000000000",
  "date": "2024-02-15",
  "updated_at": "2024-02-20T10:00:00Z"
}
```

### 3. Get All Latest Values
```http
GET /api/v1/fred/latest
```

**Response:**
```json
{
  "data": [
    {
      "ticker": "WALCL",
      "description": "Federal Reserve Total Assets",
      "value": "7500000000000",
      "date": "2024-02-15",
      "updated_at": "2024-02-20T10:00:00Z"
    },
    ...
  ],
  "timestamp": "2024-02-20T10:00:00Z"
}
```

### 4. Get Historical Data
```http
GET /api/v1/fred/ticker/:symbol?start_date=2024-01-01&end_date=2024-02-01&limit=50&sort_order=desc
```

**Query Parameters:**
- `start_date` (optional): Start date in YYYY-MM-DD format
- `end_date` (optional): End date in YYYY-MM-DD format
- `limit` (optional): Number of observations (default: 100)
- `sort_order` (optional): "asc" or "desc" (default: "desc")

**Example:**
```bash
curl "http://localhost:8080/api/v1/fred/ticker/CPIAUCSL?limit=10"
```

**Response:**
```json
{
  "ticker": "CPIAUCSL",
  "description": "Consumer Price Index (CPI)",
  "observations": [
    {
      "date": "2024-01-15",
      "value": "310.5"
    },
    ...
  ],
  "units": "Index 1982-1984=100",
  "frequency": "Monthly",
  "last_updated": "2024-02-20T10:00:00Z"
}
```

## Setup

### 1. Get FRED API Key

1. Visit: https://fred.stlouisfed.org/docs/api/api_key.html
2. Create a free account
3. Request an API key

### 2. Configure Environment

Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

Update `.env` with your API key:
```bash
FRED_API_KEY=your_actual_api_key_here
PORT=8080
```

### 3. Run the Server

```bash
make run
```

Or directly:
```bash
go run cmd/api/main.go
```

## Code Examples

### Using the Client Directly

```go
package main

import (
    "context"
    "fmt"
    "macro-analyst/internal/fred"
)

func main() {
    // Create client
    client := fred.NewClient("your-api-key")
    
    // Get latest value
    ctx := context.Background()
    latest, err := client.GetLatestValue(ctx, fred.TickerWALCL)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Fed Assets: %s (as of %s)\n", latest.Value, latest.Date)
    
    // Get historical data
    opts := &fred.QueryOptions{
        Limit: 10,
        SortOrder: "desc",
    }
    
    data, err := client.GetSeriesObservations(ctx, fred.TickerCPIAUCSL, opts)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("CPI observations: %d\n", len(data.Observations))
}
```

### Testing with Mock Client

```go
package mypackage

import (
    "context"
    "testing"
    "macro-analyst/internal/fred"
)

type MockFREDClient struct {
    GetLatestValueFunc func(ctx context.Context, ticker fred.Ticker) (*fred.LatestValue, error)
}

func (m *MockFREDClient) GetLatestValue(ctx context.Context, ticker fred.Ticker) (*fred.LatestValue, error) {
    return m.GetLatestValueFunc(ctx, ticker)
}

func TestMyFunction(t *testing.T) {
    mock := &MockFREDClient{
        GetLatestValueFunc: func(ctx context.Context, ticker fred.Ticker) (*fred.LatestValue, error) {
            return &fred.LatestValue{
                Ticker: ticker,
                Value: "100.5",
                Date: "2024-01-15",
            }, nil
        },
    }
    
    // Use mock in your tests
    result, err := mock.GetLatestValue(context.Background(), fred.TickerWALCL)
    if err != nil {
        t.Fatal(err)
    }
    
    if result.Value != "100.5" {
        t.Errorf("Expected 100.5, got %s", result.Value)
    }
}
```

## Testing

### Run All Tests
```bash
make test
```

### Run with Coverage
```bash
make test-coverage
```

### Coverage Results
- **FRED Package**: 95.3% coverage
- **WebSocket Package**: 70.7% coverage
- **Overall**: Comprehensive test coverage with edge cases

### Test Structure
- Unit tests for all public functions
- Mock HTTP client for isolated testing
- Context cancellation tests
- Error handling tests
- JSON serialization tests
- Concurrent access tests

## Design Decisions

### Why Interface-based?
- **Testability**: Easy to mock in tests without external dependencies
- **Flexibility**: Can swap implementations without changing consumers
- **Maintainability**: Clear contracts between components

### Why Separate Models?
- **Clean JSON**: Separate request/response structures
- **Type Safety**: Strong typing for all data
- **Documentation**: Self-documenting code with struct tags

### Why Context?
- **Timeout Control**: Every API call respects deadlines
- **Cancellation**: Can abort long-running requests
- **Best Practice**: Standard Go pattern for network operations

### Why Custom HTTP Client Interface?
- **Testing**: Mock HTTP responses without real network calls
- **Flexibility**: Can add retry logic, circuit breakers later
- **Observability**: Easy to add metrics/logging

## Performance Considerations

1. **Timeout**: Default 10s timeout for API calls
2. **Connection Reuse**: HTTP client reuses connections
3. **Rate Limiting**: FRED API has rate limits (check their docs)
4. **Caching**: Consider adding Redis cache for frequently accessed data

## Error Handling

All errors are wrapped with context:
```go
return nil, fmt.Errorf("failed to fetch observations for %s: %w", ticker, err)
```

This provides:
- Clear error messages
- Error chain for debugging
- Proper error types for handling

## Future Enhancements

1. **Caching Layer**: Add Redis/in-memory cache
2. **Rate Limiting**: Client-side rate limiting
3. **Retry Logic**: Exponential backoff for transient failures
4. **Metrics**: Prometheus metrics for API calls
5. **More Endpoints**: Support for FRED series metadata
6. **Batch Requests**: Parallel fetching with error aggregation

## References

- [FRED API Documentation](https://fred.stlouisfed.org/docs/api/fred/)
- [FRED Data Series](https://fred.stlouisfed.org/categories)
- [Go Context Package](https://pkg.go.dev/context)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
