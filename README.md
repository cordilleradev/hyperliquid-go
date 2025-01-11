# Hyperliquid Go Client

This Go package provides a client for accessing the Hyperliquid API, supporting both REST and WebSocket connections for fetching various types of data related to perpetuals trading.

## Features

- **REST Client**: Allows for interaction with Hyperliquid API endpoints for fetching perpetuals metadata, account summaries, and leaderboard information.
- **WebSocket Client**: Handles real-time data streaming for all market mid-prices and order updates.
- **Utilities**: Includes utility functions for handling API responses, such as parsing and error management.

## Installation

To install the Hyperliquid Go client, use `go get`:

```bash
go get github.com/cordilleradev/hyperliquid-go/pkg/client
```

## Usage

### REST Client

First, create a new instance of the `HyperliquidRestClient` with the raw URL and leaderboard URL:

```go
import "github.com/cordilleradev/hyperliquid-go/pkg/client/rest"

client, err := rest.NewHyperliquidRestClient("https://api.hyperliquid.com", "https://leaderboard.hyperliquid.com")
if err != nil {
    // handle error
}
```

#### Fetch Perpetuals Metadata

```go
metadata, err := client.PerpetualsMetadata()
if err != nil {
    // handle error
}
// use metadata
```

#### Fetch Account Summary

```go
user := "user_public_key"
accountSummary, err := client.PerpetualsAccountSummary(user)
if err != nil {
    // handle error
}
// use account summary
```

#### Fetch Leaderboard Data

```go
leaderboardData, err := client.LeaderboardCall()
if err != nil {
    // handle error
}
// use leaderboard data
```

### WebSocket Client

Create an instance of `HyperliquidWebsocketClient`:

```go
import "github.com/cordilleradev/hyperliquid-go/pkg/client/stream"

wsClient, err := stream.NewHyperliquidWebsocketClient("wss://ws.hyperliquid.com")
if err != nil {
    // handle error
}
```

#### Streaming All Market Mids

To start streaming all current market mid-prices:

```go
err = wsClient.StreamAllMids()
if err != nil {
    // handle error
}

for mids := range wsClient.AllMidsChan {
    // process mids
}
```

#### Streaming Order Updates

To stream order updates for a specific user:

```go
userId := "user_public_key"
err = wsClient.StreamOrderUpdates(userId)
if err != nil {
    // handle error
}

for orders := range wsClient.OrderChan {
    // process order updates
}
```

## Error Handling

Both REST and WebSocket clients send errors through a dedicated channel, `ErrorChan`, which can be monitored:

```go
go func() {
    for err := range wsClient.ErrorChan {
        // handle error
    }
}()
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

```
