package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-api/internal/models"
	"stock-api/internal/services"

	"github.com/gorilla/mux"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	response := Response{
		Success: false,
		Error:   message,
	}
	writeJSONResponse(w, status, response)
}

func writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeSuccessResponse(w, map[string]string{
			"status":  "healthy",
			"version": "1.0.0",
		})
	}
}

func GetStocksHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse pagination parameters
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")

		if pageStr == "" {
			pageStr = "1"
		}

		if pageSizeStr == "" {
			pageSizeStr = "20"
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			pageSize = 20
		}

		// Parse filter parameters
		filters := models.StockFilterParams{
			ActionType: r.URL.Query().Get("action_type"),
			Brokerage:  r.URL.Query().Get("brokerage"),
			SortBy:     r.URL.Query().Get("sort_by"),
		}

		paginatedStocks, err := stockService.GetStocksWithMetricsPaginated(page, pageSize, filters)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch stocks: "+err.Error())
			return
		}

		writeSuccessResponse(w, paginatedStocks)
	}
}

func SyncAllStocksHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if process can start
		canStart, err := stockService.CanStartStockSync()
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to check if stock sync can start: "+err.Error())
			return
		}

		if !canStart {
			writeErrorResponse(w, http.StatusConflict, "Stock sync process is already running or interval time hasn't passed")
			return
		}

		go stockService.SyncAllStocks()

		writeSuccessResponse(w, map[string]string{
			"message": "Syncing all stocks in the background from KarenAI API... this may take a while",
		})
	}
}

func GetStockBySymbolHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		if symbol == "" {
			writeErrorResponse(w, http.StatusBadRequest, "Symbol is required")
			return
		}

		stock, err := stockService.GetStockWithMetrics(symbol)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch stock: "+err.Error())
			return
		}

		if stock == nil {
			writeErrorResponse(w, http.StatusNotFound, "Stock not found")
			return
		}

		writeSuccessResponse(w, stock)
	}
}

func RefreshStockDataHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		if symbol == "" {
			writeErrorResponse(w, http.StatusBadRequest, "Symbol is required")
			return
		}

		err := stockService.RefreshStockData(symbol)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to refresh stock data: "+err.Error())
			return
		}

		writeSuccessResponse(w, map[string]string{
			"message": "Stock data refreshed successfully",
			"symbol":  symbol,
		})
	}
}

func GetRecommendationsHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if pagination parameters are provided
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")

		if pageStr != "" || pageSizeStr != "" {
			// Use paginated response
			page, err := strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				page = 1
			}

			pageSize, err := strconv.Atoi(pageSizeStr)
			if err != nil || pageSize < 1 {
				pageSize = 20
			}

			paginatedRecommendations, err := stockService.GetRecommendationsPaginated(page, pageSize)
			if err != nil {
				writeErrorResponse(w, http.StatusInternalServerError, "Failed to get recommendations: "+err.Error())
				return
			}

			writeSuccessResponse(w, paginatedRecommendations)
		} else {
			// Use non-paginated response for backward compatibility
			recommendations, err := stockService.GetRecommendations()
			if err != nil {
				writeErrorResponse(w, http.StatusInternalServerError, "Failed to get recommendations: "+err.Error())
				return
			}

			writeSuccessResponse(w, recommendations)
		}
	}
}

func SearchStockHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		if symbol == "" {
			writeErrorResponse(w, http.StatusBadRequest, "Symbol is required")
			return
		}

		stock, err := stockService.SearchAndAddStock(symbol)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to search stock: "+err.Error())
			return
		}

		writeSuccessResponse(w, stock)
	}
}

func GetFilterOptionsHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filterOptions, err := stockService.GetFilterOptions()
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to get filter options: "+err.Error())
			return
		}

		writeSuccessResponse(w, filterOptions)
	}
}

func GetMarketIntelligenceOverviewHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		overview, err := stockService.GetMarketIntelligenceOverview()
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to get market intelligence overview: "+err.Error())
			return
		}

		writeSuccessResponse(w, overview)
	}
}
