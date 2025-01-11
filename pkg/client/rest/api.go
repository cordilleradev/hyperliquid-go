package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cordilleradev/hyperliquid-go/pkg/client/types"
)

type HyperliquidRestClient struct {
	url            string
	leaderboardURL string
	httpClient     *http.Client
}

func NewHyperliquidRestClient(rawUrl string, leaderboardUrl string) (*HyperliquidRestClient, error) {
	urls := []string{rawUrl, leaderboardUrl}
	parsedUrls := make([]string, len(urls))

	for i, inputUrl := range urls {
		parsedUrl, err := url.Parse(inputUrl)
		if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
			return nil, fmt.Errorf("invalid URL: %s", inputUrl)
		}
		if len(parsedUrl.Path) > 0 && parsedUrl.Path[len(parsedUrl.Path)-1:] == "/" {
			parsedUrl.Path = parsedUrl.Path[:len(parsedUrl.Path)-1]
		}
		parsedUrls[i] = parsedUrl.String()
	}

	return &HyperliquidRestClient{
		url:            parsedUrls[0],
		leaderboardURL: parsedUrls[1],
		httpClient:     &http.Client{},
	}, nil
}

func (rest *HyperliquidRestClient) doRequest(endpoint string, requestType string, additionalParams map[string]string, result interface{}) error {
	reqBody := map[string]string{"type": requestType}
	for key, value := range additionalParams {
		reqBody[key] = value
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", rest.url+endpoint, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return fmt.Errorf("failed to create new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rest.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(byteData, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return nil
}

// New method for fetching metadata
func (rest *HyperliquidRestClient) PerpetualsMetadata() (types.Universe, error) {
	var metadata types.Universe
	if err := rest.doRequest("/info", "meta", nil, &metadata); err != nil {
		return types.Universe{}, err
	}
	return metadata, nil
}

// Updated PerpetualAssetContexts function
func (rest *HyperliquidRestClient) PerpetualAssetContexts() ([]types.UniverseItem, []types.AssetContext, error) {
	var respData []json.RawMessage
	if err := rest.doRequest("/info", "metaAndAssetCtxs", nil, &respData); err != nil {
		return nil, nil, err
	}

	// Parse universe part
	var universeData types.Universe
	if err := json.Unmarshal(respData[0], &universeData); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal universe data: %w", err)
	}

	// Parse asset contexts part
	var assetContexts []types.AssetContext
	if err := json.Unmarshal(respData[1], &assetContexts); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal asset contexts: %w", err)
	}

	return universeData.Universe, assetContexts, nil
}

// New method for fetching a user's perpetuals account summary
func (rest *HyperliquidRestClient) PerpetualsAccountSummary(user string) (types.MarginData, error) {
	var marginData types.MarginData
	params := map[string]string{"user": user}
	if err := rest.doRequest("/info", "clearinghouseState", params, &marginData); err != nil {
		return types.MarginData{}, err
	}
	return marginData, nil
}

func (rest *HyperliquidRestClient) LeaderboardCall() ([]types.TraderPerformance, error) {
	req, err := http.NewRequest("GET", rest.leaderboardURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := rest.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch leaderboard: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawResponse types.RawLeaderboardResponse
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		return nil, err
	}

	traderPerformances := make([]types.TraderPerformance, len(rawResponse.LeaderboardRows))
	for i, item := range rawResponse.LeaderboardRows {
		accountValue := parseStringToFloat(item.AccountValue)
		windowPerformances := item.WindowPerformances
		traderPerformance := types.TraderPerformance{
			EthAddress: item.EthAddress,
			Day: types.PeriodPerformance{
				Pnl: parseStringToFloat(windowPerformances[0][1].(map[string]interface{})["pnl"].(string)),
				Roi: parseStringToFloat(windowPerformances[0][1].(map[string]interface{})["roi"].(string)),
				Vlm: parseStringToFloat(windowPerformances[0][1].(map[string]interface{})["vlm"].(string)),
			},
			Week: types.PeriodPerformance{
				Pnl: parseStringToFloat(windowPerformances[1][1].(map[string]interface{})["pnl"].(string)),
				Roi: parseStringToFloat(windowPerformances[1][1].(map[string]interface{})["roi"].(string)),
				Vlm: parseStringToFloat(windowPerformances[1][1].(map[string]interface{})["vlm"].(string)),
			},
			Month: types.PeriodPerformance{
				Pnl: parseStringToFloat(windowPerformances[2][1].(map[string]interface{})["pnl"].(string)),
				Roi: parseStringToFloat(windowPerformances[2][1].(map[string]interface{})["roi"].(string)),
				Vlm: parseStringToFloat(windowPerformances[2][1].(map[string]interface{})["vlm"].(string)),
			},
			AllTime: types.PeriodPerformance{
				Pnl: parseStringToFloat(windowPerformances[3][1].(map[string]interface{})["pnl"].(string)),
				Roi: parseStringToFloat(windowPerformances[3][1].(map[string]interface{})["roi"].(string)),
				Vlm: parseStringToFloat(windowPerformances[3][1].(map[string]interface{})["vlm"].(string)),
			},
			AccountValue: accountValue,
		}
		traderPerformances[i] = traderPerformance
	}
	return traderPerformances, nil
}
