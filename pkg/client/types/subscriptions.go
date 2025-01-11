package types

type GenericMessage struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type AllMidsMessage struct {
	Channel string `json:"channel"`
	Data    Mids   `json:"data"`
}

type Mids struct {
	Prices map[string]string `json:"mids"`
}

type OrderUpdatesMessage struct {
	Channel string  `json:"channel"`
	Data    []Order `json:"data"`
}
type Order struct {
	Order           BasicOrder `json:"order"`
	Status          string     `json:"status"`          // Status: "open", "filled", "canceled", "triggered", "rejected", etc.
	StatusTimestamp int64      `json:"statusTimestamp"` // Last status update timestamp (in ms)
}

type BasicOrder struct {
	Coin      string `json:"coin"`      // Coin symbol (e.g., BTC)
	Side      string `json:"side"`      // "B" for buy, "S" for sell
	LimitPx   string `json:"limitPx"`   // Limit price of the order
	Sz        string `json:"sz"`        // Size of the order
	Oid       int64  `json:"oid"`       // Unique order ID
	Timestamp int64  `json:"timestamp"` // Order creation timestamp
	OrigSz    string `json:"origSz"`    // Original order size
	Cloid     string `json:"cloid"`     // Client order ID (changed from optional to mandatory)
}
