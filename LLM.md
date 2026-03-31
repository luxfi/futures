# Futures Module

**Module**: `github.com/luxfi/futures`
**Purpose**: Futures/FCM trading service with pluggable provider backends.

## Providers

| Provider | Env Vars | Description |
|----------|----------|-------------|
| Apex Futures | `APEX_FUTURES_API_KEY`, `APEX_FUTURES_API_SECRET`, `APEX_FUTURES_SANDBOX` | Apex Clearing FCM division. CME, NYMEX, CBOT, ICE. |
| IBKR | `IBKR_FUTURES_ENABLED`, `IBKR_FUTURES_PAPER` | Interactive Brokers Client Portal Gateway. Global futures exchanges. |

## API Endpoints

```
GET  /health                                     - Health check
GET  /providers                                  - List registered providers
GET  /v1/contracts/{provider}/{underlying}        - List futures contracts
GET  /v1/quote/{provider}/{symbol}                - Get real-time quote
POST /v1/order/{provider}/{accountID}             - Place futures order
GET  /v1/positions/{provider}/{accountID}         - List open positions
DELETE /v1/position/{provider}/{accountID}/{symbol} - Close position
GET  /v1/margin/{provider}/{accountID}/{symbol}   - Get margin requirements
```

## Running

```bash
# With Apex Futures
APEX_FUTURES_API_KEY=xxx APEX_FUTURES_API_SECRET=yyy futuresd

# With IBKR (requires Client Portal Gateway running on localhost:5000)
IBKR_FUTURES_ENABLED=true IBKR_FUTURES_PAPER=true futuresd

# Custom listen address
LISTEN_ADDR=:9090 futuresd
```

## Structure

```
cmd/futuresd/           - HTTP server entry point
pkg/types/              - Domain types (contracts, quotes, positions, margins, orders)
pkg/provider/           - FuturesProvider interface + Registry
pkg/provider/apex/      - Apex Futures FCM implementation
pkg/provider/ibkr/      - IBKR futures implementation
pkg/provider/envconfig/ - Env-based provider registration
```
