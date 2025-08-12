package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type KarenAIClient struct {
	apiToken string
	baseURL  string
	client   *http.Client
}

type StockAnalysis struct {
	Ticker     string    `json:"ticker"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

type StockListResponse struct {
	Items    []StockAnalysis `json:"items"`
	NextPage string          `json:"next_page"`
}

func NewKarenAIClient(apiToken string) *KarenAIClient {
	return &KarenAIClient{
		apiToken: apiToken,
		baseURL:  "https://api.karenai.click/swechallenge",
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *KarenAIClient) GetStocksList(nextPage string) (*StockListResponse, error) {
	url := c.baseURL + "/list"
	if nextPage != "" {
		url += "?next_page=" + nextPage
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stocks list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var stockList StockListResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stockList, nil
}

func (c *KarenAIClient) GetAllStocks() ([]StockAnalysis, error) {
	var allStocks []StockAnalysis
	nextPage := ""

	for {
		response, err := c.GetStocksList(nextPage)
		if err != nil {
			return nil, err
		}

		allStocks = append(allStocks, response.Items...)

		if response.NextPage == "" {
			break
		}

		nextPage = response.NextPage

		// Add a small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	return allStocks, nil
}
