package api

import (
	"stock-api/internal/services"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, stockService *services.StockService) {
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/stocks", GetStocksHandler(stockService)).Methods("GET")
	api.HandleFunc("/stocks/sync", SyncAllStocksHandler(stockService)).Methods("POST")
	api.HandleFunc("/stocks/filter-options", GetFilterOptionsHandler(stockService)).Methods("GET")
	api.HandleFunc("/stocks/recommendations", GetRecommendationsHandler(stockService)).Methods("GET")
	api.HandleFunc("/analytics/market-intelligence-overview", GetMarketIntelligenceOverviewHandler(stockService)).Methods("GET")
	api.HandleFunc("/stocks/{symbol}", GetStockBySymbolHandler(stockService)).Methods("GET")
	api.HandleFunc("/stocks/{symbol}/refresh", RefreshStockDataHandler(stockService)).Methods("POST")
	api.HandleFunc("/stocks/search/{symbol}", SearchStockHandler(stockService)).Methods("GET")

	api.HandleFunc("/health", HealthHandler()).Methods("GET")
}
