package types

type UniverseItem struct {
	Name         string `json:"name"`
	SzDecimals   int    `json:"szDecimals"`
	MaxLeverage  int    `json:"maxLeverage"`
	OnlyIsolated bool   `json:"onlyIsolated"`
}

type Universe struct {
	Universe []UniverseItem `json:"universe"`
}

type ImpactPrices []string

type AssetContext struct {
	DayNtlVlm    string       `json:"dayNtlVlm"`
	Funding      string       `json:"funding"`
	ImpactPxs    ImpactPrices `json:"impactPxs"`
	MarkPx       string       `json:"markPx"`
	MidPx        string       `json:"midPx"`
	OpenInterest string       `json:"openInterest"`
	OraclePx     string       `json:"oraclePx"`
	Premium      string       `json:"premium"`
	PrevDayPx    string       `json:"prevDayPx"`
}
