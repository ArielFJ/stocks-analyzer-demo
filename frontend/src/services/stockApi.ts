import type {
  ApiResponse,
  PaginatedResponse,
  PaginationParams,
  StockFilterParams,
  StockWithAnalysis,
  StockRecommendation,
  HealthResponse,
  SyncResponse,
  RefreshResponse,
  FilterOptions,
  MarketIntelligenceOverview
} from '@/types/api'

class StockApiService {
  private baseUrl: string

  constructor(baseUrl: string = 'http://localhost:8080/api/v1') {
    this.baseUrl = baseUrl
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      
      if (!data.success) {
        throw new Error(data.error || 'API request failed')
      }

      return data
    } catch (error) {
      console.error(`API request failed for ${url}:`, error)
      throw error
    }
  }

  // Health check
  async checkHealth(): Promise<ApiResponse<HealthResponse>> {
    return this.request<ApiResponse<HealthResponse>>('/health')
  }

  // Stocks endpoints
  async getAllStocks(): Promise<ApiResponse<StockWithAnalysis[]>> {
    return this.request<ApiResponse<StockWithAnalysis[]>>('/stocks')
  }

  async getStocksPaginated(params: PaginationParams = {}, filters: StockFilterParams = {}): Promise<ApiResponse<PaginatedResponse<StockWithAnalysis>>> {
    const searchParams = new URLSearchParams()
    
    if (params.page) {
      searchParams.append('page', params.page.toString())
    }
    if (params.page_size) {
      searchParams.append('page_size', params.page_size.toString())
    }
    if (filters.action_type) {
      searchParams.append('action_type', filters.action_type)
    }
    if (filters.brokerage) {
      searchParams.append('brokerage', filters.brokerage)
    }
    if (filters.sort_by) {
      searchParams.append('sort_by', filters.sort_by)
    }

    const endpoint = `/stocks${searchParams.toString() ? `?${searchParams.toString()}` : ''}`
    return this.request<ApiResponse<PaginatedResponse<StockWithAnalysis>>>(endpoint)
  }

  async getFilterOptions(): Promise<ApiResponse<FilterOptions>> {
    return this.request<ApiResponse<FilterOptions>>('/stocks/filter-options')
  }

  async getMarketIntelligenceOverview(): Promise<ApiResponse<MarketIntelligenceOverview>> {
    return this.request<ApiResponse<MarketIntelligenceOverview>>('/analytics/market-intelligence-overview')
  }

  async getStockBySymbol(symbol: string): Promise<ApiResponse<StockWithAnalysis>> {
    return this.request<ApiResponse<StockWithAnalysis>>(`/stocks/${symbol}`)
  }

  async searchStock(symbol: string): Promise<ApiResponse<StockWithAnalysis>> {
    return this.request<ApiResponse<StockWithAnalysis>>(`/stocks/search/${symbol}`)
  }

  async refreshStockData(symbol: string): Promise<ApiResponse<RefreshResponse>> {
    return this.request<ApiResponse<RefreshResponse>>(`/stocks/${symbol}/refresh`, {
      method: 'POST'
    })
  }

  // Sync endpoint
  async syncAllStocks(): Promise<ApiResponse<SyncResponse>> {
    return this.request<ApiResponse<SyncResponse>>('/stocks/sync', {
      method: 'POST'
    })
  }

  // Recommendations endpoints
  async getAllRecommendations(): Promise<ApiResponse<StockRecommendation[]>> {
    return this.request<ApiResponse<StockRecommendation[]>>('/stocks/recommendations')
  }

  async getRecommendationsPaginated(params: PaginationParams = {}): Promise<ApiResponse<PaginatedResponse<StockRecommendation>>> {
    const searchParams = new URLSearchParams()
    
    if (params.page) {
      searchParams.append('page', params.page.toString())
    }
    if (params.page_size) {
      searchParams.append('page_size', params.page_size.toString())
    }

    const endpoint = `/stocks/recommendations${searchParams.toString() ? `?${searchParams.toString()}` : ''}`
    return this.request<ApiResponse<PaginatedResponse<StockRecommendation>>>(endpoint)
  }
}

// Create and export a singleton instance
export const stockApi = new StockApiService()
export default stockApi