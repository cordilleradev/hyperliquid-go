package types

type CumFunding struct {
	AllTime     string `json:"allTime"`
	SinceChange string `json:"sinceChange"`
	SinceOpen   string `json:"sinceOpen"`
}

type Leverage struct {
	RawUsd string `json:"rawUsd"`
	Type   string `json:"type"`
	Value  int    `json:"value"`
}

type Position struct {
	Coin           string     `json:"coin"`
	CumFunding     CumFunding `json:"cumFunding"`
	EntryPx        string     `json:"entryPx"`
	Leverage       Leverage   `json:"leverage"`
	LiquidationPx  string     `json:"liquidationPx"`
	MarginUsed     string     `json:"marginUsed"`
	MaxLeverage    int        `json:"maxLeverage"`
	PositionValue  string     `json:"positionValue"`
	ReturnOnEquity string     `json:"returnOnEquity"`
	Szi            string     `json:"szi"`
	UnrealizedPnl  string     `json:"unrealizedPnl"`
}

type AssetPosition struct {
	Position Position `json:"position"`
	Type     string   `json:"type"`
}

type CrossMarginSummary struct {
	AccountValue    string `json:"accountValue"`
	TotalMarginUsed string `json:"totalMarginUsed"`
	TotalNtlPos     string `json:"totalNtlPos"`
	TotalRawUsd     string `json:"totalRawUsd"`
}

type MarginSummary struct {
	AccountValue    string `json:"accountValue"`
	TotalMarginUsed string `json:"totalMarginUsed"`
	TotalNtlPos     string `json:"totalNtlPos"`
	TotalRawUsd     string `json:"totalRawUsd"`
}

type MarginData struct {
	AssetPositions             []AssetPosition    `json:"assetPositions"`
	CrossMaintenanceMarginUsed string             `json:"crossMaintenanceMarginUsed"`
	CrossMarginSummary         CrossMarginSummary `json:"crossMarginSummary"`
	MarginSummary              MarginSummary      `json:"marginSummary"`
	Time                       int64              `json:"time"`
	Withdrawable               string             `json:"withdrawable"`
}
