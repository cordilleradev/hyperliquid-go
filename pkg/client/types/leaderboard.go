package types

type TraderPerformance struct {
	EthAddress   string
	AccountValue float64
	Day          PeriodPerformance
	Week         PeriodPerformance
	Month        PeriodPerformance
	AllTime      PeriodPerformance
}

type PeriodPerformance struct {
	Pnl float64 `json:"pnl"`
	Roi float64 `json:"roi"`
	Vlm float64 `json:"vlm"`
}

type LeaderboardRow struct {
	EthAddress         string          `json:"ethAddress"`
	AccountValue       string          `json:"accountValue"`
	WindowPerformances [][]interface{} `json:"windowPerformances"`
}

type RawLeaderboardResponse struct {
	LeaderboardRows []LeaderboardRow `json:"leaderboardRows"`
}
