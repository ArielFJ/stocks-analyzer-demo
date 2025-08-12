package models

import (
	"database/sql/driver"
	"time"
)

type Stock struct {
	ID        int       `json:"id"`
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StockAnalysis struct {
	ID         int       `json:"id"`
	StockID    int       `json:"stock_id"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	AnalysisDate time.Time `json:"analysis_date"`
	CreatedAt  time.Time `json:"created_at"`
}

type StockWithAnalysis struct {
	Stock
	LatestAnalysis []StockAnalysis `json:"latest_analysis,omitempty"`
}

type StockRecommendation struct {
	Stock      StockWithAnalysis `json:"stock"`
	Score      float64           `json:"score"`
	Reason     string            `json:"reason"`
	Confidence string            `json:"confidence"`
}

type RecommendationScore struct {
	ID                 int     `json:"id"`
	StockID            int     `json:"stock_id"`
	TotalScore         float64 `json:"total_score"`
	RatingScore        float64 `json:"rating_score"`
	RatingChangeScore  float64 `json:"rating_change_score"`
	TargetChangeScore  float64 `json:"target_change_score"`
	ActionScore        float64 `json:"action_score"`
	CoverageScore      float64 `json:"coverage_score"`
	Confidence         string  `json:"confidence"`
	Reason             string  `json:"reason"`
	LatestAnalysisID   *int    `json:"latest_analysis_id,omitempty"`
	CalculatedAt       time.Time `json:"calculated_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type RecommendationWithStock struct {
	RecommendationScore
	Stock StockWithAnalysis `json:"stock"`
}

type ProcessControl struct {
	ID             int       `json:"id"`
	ProcessName    string    `json:"process_name"`
	IsRunning      bool      `json:"is_running"`
	LastExecution  *time.Time `json:"last_execution"`
	IntervalMinutes int       `json:"interval_minutes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PaginationMeta struct {
	Page         int  `json:"page"`
	PageSize     int  `json:"page_size"`
	TotalItems   int  `json:"total_items"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
}

type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type StockFilterParams struct {
	ActionType string `json:"action_type" query:"action_type"`
	Brokerage  string `json:"brokerage" query:"brokerage"`
	SortBy     string `json:"sort_by" query:"sort_by"`
}

type FilterOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type FilterOptions struct {
	ActionTypes []FilterOption `json:"action_types"`
	Brokerages  []FilterOption `json:"brokerages"`
	SortBy      []FilterOption `json:"sort_by"`
}

type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

func (nf NullFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}

type MarketIntelligenceOverview struct {
	TotalStocks                int                    `json:"total_stocks"`
	TotalRecommendations       int                    `json:"total_recommendations"`
	RecentAnalysis             int                    `json:"recent_analysis"`
	Upgrades                   int                    `json:"upgrades"`
	Downgrades                 int                    `json:"downgrades"`
	HighConfidenceRecs         int                    `json:"high_confidence_recs"`
	SelectionRate              float64                `json:"selection_rate"`
	TopBrokerages              []BrokerageAnalytics   `json:"top_brokerages"`
	TopActionTypes             []ActionTypeAnalytics  `json:"top_action_types"`
	RecentActivityTrend        []ActivityTrendPoint   `json:"recent_activity_trend"`
	AverageRecommendationScore float64                `json:"average_recommendation_score"`
}

type BrokerageAnalytics struct {
	Brokerage     string  `json:"brokerage"`
	AnalysisCount int     `json:"analysis_count"`
	Percentage    float64 `json:"percentage"`
}

type ActionTypeAnalytics struct {
	ActionType    string  `json:"action_type"`
	Count         int     `json:"count"`
	Percentage    float64 `json:"percentage"`
}

type ActivityTrendPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}