# FRED API - Example Responses

This document provides real-world example responses from the FRED API endpoints.

## 1. GET /api/v1/fred/tickers

List all available macroeconomic tickers.

**Request:**
```bash
curl http://localhost:8080/api/v1/fred/tickers
```

**Response (200 OK):**
```json
{
  "tickers": [
    {
      "symbol": "WALCL",
      "description": "Federal Reserve Total Assets"
    },
    {
      "symbol": "WTREGEN",
      "description": "Treasury General Account"
    },
    {
      "symbol": "RRPONTSYD",
      "description": "Overnight Reverse Repo"
    },
    {
      "symbol": "FEDFUNDS",
      "description": "Federal Funds Rate"
    },
    {
      "symbol": "CPIAUCSL",
      "description": "Consumer Price Index (CPI)"
    },
    {
      "symbol": "DTWEXBGS",
      "description": "US Dollar Index"
    }
  ],
  "count": 6
}
```

---

## 2. GET /api/v1/fred/latest/:symbol

Get the most recent value for a specific ticker.

**Request:**
```bash
curl http://localhost:8080/api/v1/fred/latest/WALCL
```

**Response (200 OK):**
```json
{
  "ticker": "WALCL",
  "description": "Federal Reserve Total Assets",
  "value": "7769411000000",
  "date": "2024-02-14",
  "updated_at": "2024-02-20T10:15:30Z"
}
```

**Example - Fed Funds Rate:**
```bash
curl http://localhost:8080/api/v1/fred/latest/FEDFUNDS
```

```json
{
  "ticker": "FEDFUNDS",
  "description": "Federal Funds Rate",
  "value": "5.33",
  "date": "2024-01-01",
  "updated_at": "2024-02-20T10:15:30Z"
}
```

**Example - CPI (Inflation):**
```bash
curl http://localhost:8080/api/v1/fred/latest/CPIAUCSL
```

```json
{
  "ticker": "CPIAUCSL",
  "description": "Consumer Price Index (CPI)",
  "value": "310.326",
  "date": "2024-01-01",
  "updated_at": "2024-02-20T10:15:30Z"
}
```

---

## 3. GET /api/v1/fred/latest

Get latest values for all tickers at once.

**Request:**
```bash
curl http://localhost:8080/api/v1/fred/latest
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "ticker": "WALCL",
      "description": "Federal Reserve Total Assets",
      "value": "7769411000000",
      "date": "2024-02-14",
      "updated_at": "2024-02-20T10:15:30Z"
    },
    {
      "ticker": "WTREGEN",
      "description": "Treasury General Account",
      "value": "750000000000",
      "date": "2024-02-14",
      "updated_at": "2024-02-20T10:15:31Z"
    },
    {
      "ticker": "RRPONTSYD",
      "description": "Overnight Reverse Repo",
      "value": "500000000000",
      "date": "2024-02-14",
      "updated_at": "2024-02-20T10:15:32Z"
    },
    {
      "ticker": "FEDFUNDS",
      "description": "Federal Funds Rate",
      "value": "5.33",
      "date": "2024-01-01",
      "updated_at": "2024-02-20T10:15:33Z"
    },
    {
      "ticker": "CPIAUCSL",
      "description": "Consumer Price Index (CPI)",
      "value": "310.326",
      "date": "2024-01-01",
      "updated_at": "2024-02-20T10:15:34Z"
    },
    {
      "ticker": "DTWEXBGS",
      "description": "US Dollar Index",
      "value": "108.5432",
      "date": "2024-02-13",
      "updated_at": "2024-02-20T10:15:35Z"
    }
  ],
  "timestamp": "2024-02-20T10:15:35Z"
}
```

---

## 4. GET /api/v1/fred/ticker/:symbol

Get historical observations for a ticker with full metadata.

**Request - Last 10 CPI readings:**
```bash
curl "http://localhost:8080/api/v1/fred/ticker/CPIAUCSL?limit=10&sort_order=desc"
```

**Response (200 OK):**
```json
{
  "ticker": "CPIAUCSL",
  "description": "Consumer Price Index (CPI)",
  "title": "Consumer Price Index for All Urban Consumers: All Items in U.S. City Average",
  "observations": [
    {
      "date": "2024-01-01",
      "value": "310.326"
    },
    {
      "date": "2023-12-01",
      "value": "309.685"
    },
    {
      "date": "2023-11-01",
      "value": "308.734"
    },
    {
      "date": "2023-10-01",
      "value": "308.258"
    },
    {
      "date": "2023-09-01",
      "value": "308.038"
    },
    {
      "date": "2023-08-01",
      "value": "307.789"
    },
    {
      "date": "2023-07-01",
      "value": "306.480"
    },
    {
      "date": "2023-06-01",
      "value": "305.109"
    },
    {
      "date": "2023-05-01",
      "value": "304.127"
    },
    {
      "date": "2023-04-01",
      "value": "303.363"
    }
  ],
  "units": "Index 1982-1984=100",
  "units_short": "Index 1982-1984=100",
  "frequency": "Monthly",
  "notes": "Indexes are available for the U.S. city average, the 4 regions, 3 population-size classes...",
  "last_updated": "2024-02-20T10:20:00Z"
}
```

**Request - Fed Assets with date range:**
```bash
curl "http://localhost:8080/api/v1/fred/ticker/WALCL?start_date=2024-01-01&end_date=2024-02-01&limit=5"
```

**Response (200 OK):**
```json
{
  "ticker": "WALCL",
  "description": "Federal Reserve Total Assets",
  "title": "Assets: Total Assets: Total Assets (Less Eliminations from Consolidation): Wednesday Level",
  "observations": [
    {
      "date": "2024-01-31",
      "value": "7780000000000"
    },
    {
      "date": "2024-01-24",
      "value": "7775000000000"
    },
    {
      "date": "2024-01-17",
      "value": "7770000000000"
    },
    {
      "date": "2024-01-10",
      "value": "7765000000000"
    },
    {
      "date": "2024-01-03",
      "value": "7760000000000"
    }
  ],
  "units": "Millions of Dollars",
  "units_short": "Mil. of $",
  "frequency": "Weekly, As of Wednesday",
  "notes": "For further information regarding treasury liabilities, please see the U.S. Department of Treasury...",
  "last_updated": "2024-02-20T10:25:00Z"
}
```

---

## Query Parameters

### GET /api/v1/fred/ticker/:symbol

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `start_date` | string | No | - | Start date (YYYY-MM-DD) |
| `end_date` | string | No | - | End date (YYYY-MM-DD) |
| `limit` | integer | No | 100 | Number of observations |
| `sort_order` | string | No | "desc" | "asc" or "desc" |

**Examples:**

```bash
# Last 50 observations
curl "http://localhost:8080/api/v1/fred/ticker/WALCL?limit=50"

# Date range
curl "http://localhost:8080/api/v1/fred/ticker/CPIAUCSL?start_date=2023-01-01&end_date=2023-12-31"

# Ascending order
curl "http://localhost:8080/api/v1/fred/ticker/FEDFUNDS?sort_order=asc&limit=20"

# Combined
curl "http://localhost:8080/api/v1/fred/ticker/DTWEXBGS?start_date=2024-01-01&limit=30&sort_order=desc"
```

---

## Error Responses

### 503 Service Unavailable (FRED API not configured)

**Request:**
```bash
curl http://localhost:8080/api/v1/fred/latest/WALCL
```

**Response (503):**
```json
{
  "error": "FRED API client not configured"
}
```

**Solution:** Set `FRED_API_KEY` in your `.env` file.

---

### 500 Internal Server Error (API Error)

**Request:**
```bash
curl http://localhost:8080/api/v1/fred/latest/INVALID_TICKER
```

**Response (500):**
```json
{
  "error": "failed to fetch observations for INVALID_TICKER: API returned status 400: Bad Request - series_id not found"
}
```

---

## Understanding the Data

### Metadata Fields

All ticker responses now include full metadata from FRED:

- **`title`**: Official FRED series title (more detailed than description)
- **`units`**: Full unit description (e.g., "Millions of Dollars")
- **`units_short`**: Abbreviated unit (e.g., "Mil. of $")
- **`frequency`**: Update frequency with details (e.g., "Weekly, As of Wednesday")
- **`notes`**: Detailed notes about the data series methodology

### WALCL (Fed Assets)
- **Units**: Millions of Dollars (Mil. of $)
- **Frequency**: Weekly (Wednesday)
- **What it means**: Total size of Fed's balance sheet. Increases = money printing (QE), Decreases = tightening (QT)

### WTREGEN (Treasury Account)
- **Units**: Millions of Dollars (Mil. of $)
- **Frequency**: Daily
- **What it means**: Cash balance in US Treasury's account at the Fed

### RRPONTSYD (Reverse Repo)
- **Units**: Billions of Dollars (Bil. of $)
- **Frequency**: Daily
- **What it means**: Money parked at the Fed by financial institutions, effectively removing liquidity from markets

### FEDFUNDS (Fed Funds Rate)
- **Units**: Percent (%)
- **Frequency**: Monthly
- **What it means**: Target interest rate for overnight lending between banks

### CPIAUCSL (CPI)
- **Units**: Index 1982-1984=100
- **Frequency**: Monthly
- **What it means**: Consumer Price Index, measures inflation

### DTWEXBGS (Dollar Index)
- **Units**: Index (Jan 2006 = 100)
- **Frequency**: Daily
- **What it means**: Trade-weighted US dollar strength vs basket of currencies

---

## Rate Limits

FRED API has rate limits:
- **120 requests per minute** (free tier)
- **Hourly limit**: varies by API key type

Consider implementing caching for production use.

---

## Tips for Production

1. **Cache responses**: Use Redis or in-memory cache
2. **Error handling**: Implement retry logic with exponential backoff
3. **Rate limiting**: Client-side rate limiting to respect FRED limits
4. **Monitoring**: Add metrics for API call success/failure rates
5. **Fallback**: Have backup data sources or cached values
6. **Alerting**: Alert on API failures or stale data
