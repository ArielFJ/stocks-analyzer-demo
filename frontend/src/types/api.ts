// Base API types
export interface ApiResponse<T> {
  success: boolean
  data: T
  error?: string
}

// Pagination types
export interface PaginationMeta {
  page: number
  page_size: number
  total_items: number
  total_pages: number
  has_next: boolean
  has_previous: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  meta: PaginationMeta
}

export interface PaginationParams {
  page?: number
  page_size?: number
}

export interface StockFilterParams {
  action_type?: string
  brokerage?: string
  sort_by?: string
}

export interface FilterOption {
  label: string
  value: string
}

export interface FilterOptions {
  action_types: FilterOption[]
  brokerages: FilterOption[]
  sort_by: FilterOption[]
}

// Stock types
export interface Stock {
  id: number
  symbol: string
  name: string
  created_at: string
  updated_at: string
}

export interface StockAnalysis {
  id: number
  stock_id: number
  target_from: string
  target_to: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  analysis_date: string
  created_at: string
}

export interface StockWithAnalysis {
  id: number
  symbol: string
  name: string
  created_at: string
  updated_at: string
  latest_analysis?: StockAnalysis[]
}

export interface StockRecommendation {
  stock: StockWithAnalysis
  score: number
  reason: string
  confidence: string
}

// Health check response
export interface HealthResponse {
  status: string
  version: string
}

// Sync response
export interface SyncResponse {
  message: string
}

// Refresh response
export interface RefreshResponse {
  message: string
  symbol: string
}

// API Error
export interface ApiError {
  success: false
  error: string
}

// Analytics types
export interface BrokerageAnalytics {
  brokerage: string
  analysis_count: number
  percentage: number
}

export interface ActionTypeAnalytics {
  action_type: string
  count: number
  percentage: number
}

export interface ActivityTrendPoint {
  date: string
  count: number
}

export interface MarketIntelligenceOverview {
  total_stocks: number
  total_recommendations: number
  recent_analysis: number
  upgrades: number
  downgrades: number
  high_confidence_recs: number
  selection_rate: number
  top_brokerages: BrokerageAnalytics[]
  top_action_types: ActionTypeAnalytics[]
  recent_activity_trend: ActivityTrendPoint[]
  average_recommendation_score: number
}