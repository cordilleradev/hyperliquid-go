package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/cordilleradev/hyperliquid-go/pkg/client/types"
	"github.com/gorilla/websocket"
)

type HyperliquidWebsocketClient struct {
	websocketUrl string
	conn         *websocket.Conn
	AllMidsChan  chan map[string]float64
	OrderChan    chan []types.Order
	ErrorChan    chan error
	mutex        sync.Mutex
	lastRequest  time.Time
}

func NewHyperliquidWebsocketClient(rawUrl string) (*HyperliquidWebsocketClient, error) {
	if _, err := url.ParseRequestURI(rawUrl); err != nil {
		return nil, errors.New("Invalid websocket URL")
	}
	conn, _, err := websocket.DefaultDialer.Dial(rawUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &HyperliquidWebsocketClient{
		websocketUrl: rawUrl,
		conn:         conn,
		AllMidsChan:  make(chan map[string]float64),
		OrderChan:    make(chan []types.Order),
		ErrorChan:    make(chan error),
	}

	go func() {
		ticker := time.NewTicker(50 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				client.mutex.Lock()
				pingMessage := map[string]interface{}{
					"method": "ping",
				}
				if err := client.sendMessage(pingMessage); err != nil {
					fmt.Println("Error sending ping:", err)
					client.ErrorChan <- err
				}
				client.mutex.Unlock()
			}
		}
	}()

	go func() {
		for {
			_, msg, err := client.conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				client.ErrorChan <- err
				return
			}

			var response types.GenericMessage
			if err := json.Unmarshal(msg, &response); err != nil {
				fmt.Println("Unmarshal error:", err)
				client.ErrorChan <- err
				continue
			}

			switch response.Channel {
			case "allMids":
				client.handleAllMidsMessage(msg)
			case "orderUpdates":
				client.handleOrderUpdatesMessage(msg)
			}
		}
	}()

	return client, nil
}

func (client *HyperliquidWebsocketClient) handleAllMidsMessage(msg json.RawMessage) {
	var midsResponse types.AllMidsMessage
	if err := json.Unmarshal(msg, &midsResponse); err != nil {
		fmt.Println("Unmarshal error:", err)
		client.ErrorChan <- err
		return
	}

	prices := make(map[string]float64)
	for k, v := range midsResponse.Data.Prices {
		if k[0] != '@' {
			floatVal, _ := strconv.ParseFloat(v, 64)
			prices[k] = floatVal
		}
	}

	client.AllMidsChan <- prices
}

func (client *HyperliquidWebsocketClient) handleOrderUpdatesMessage(msg json.RawMessage) {
	var orderResponse types.OrderUpdatesMessage
	if err := json.Unmarshal(msg, &orderResponse); err != nil {
		fmt.Println("Unmarshal error:", err)
		client.ErrorChan <- err
		return
	}

	client.OrderChan <- orderResponse.Data
}

func (client *HyperliquidWebsocketClient) sendMessage(message interface{}) error {
	// Ensure at least 50ms between requests
	timeSinceLastRequest := time.Since(client.lastRequest)
	if timeSinceLastRequest < 50*time.Millisecond {
		time.Sleep(50*time.Millisecond - timeSinceLastRequest)
	}
	client.lastRequest = time.Now()

	return client.conn.WriteJSON(message)
}

func (client *HyperliquidWebsocketClient) StreamAllMids() error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	subscriptionMessage := map[string]interface{}{
		"method": "subscribe",
		"subscription": map[string]interface{}{
			"type": "allMids",
		},
	}
	if err := client.sendMessage(subscriptionMessage); err != nil {
		close(client.AllMidsChan)
		client.ErrorChan <- err
		return fmt.Errorf("Subscription error: %w", err)
	}
	return nil
}

func (client *HyperliquidWebsocketClient) StreamOrderUpdates(userId string) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	subscriptionMessage := map[string]interface{}{
		"method": "subscribe",
		"subscription": map[string]interface{}{
			"type": "orderUpdates",
			"user": userId,
		},
	}
	if err := client.sendMessage(subscriptionMessage); err != nil {
		close(client.OrderChan)
		client.ErrorChan <- err
		return fmt.Errorf("Subscription error: %w", err)
	}
	return nil
}
