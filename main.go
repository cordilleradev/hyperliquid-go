package main

import (
	"fmt"

	"github.com/cordilleradev/hyperliquid-go/pkg/client/stream"
)

func main() {
	// restClient, _ := rest.NewHyperliquidRestClient(
	// 	"https://api.hyperliquid.xyz",
	// 	"https://stats-data.hyperliquid.xyz/Mainnet/leaderboard",
	// )

	ethAddress := "0x880ac484a1743862989a441d6d867238c7aa311c"

	websocketClient, _ := stream.NewHyperliquidWebsocketClient("wss://api.hyperliquid.xyz/ws")
	webSocketErr := websocketClient.StreamOrderUpdates(ethAddress)
	if webSocketErr != nil {
		fmt.Println("Error streaming all mids:", webSocketErr)
		return
	}
	websocketClient.StreamAllMids()

	go func() {
		for range websocketClient.OrderChan {
			fmt.Println("x")
		}
	}()

	go func() {
		for range websocketClient.AllMidsChan {
			fmt.Println("y")
		}
	}()

	for e := range websocketClient.ErrorChan {
		fmt.Println(e)
	}
}
