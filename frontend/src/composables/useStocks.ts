import { ref, computed, type Ref } from 'vue'
import { stockApi } from '@/services/stockApi'
import type {
  StockWithAnalysis,
  PaginationParams,
  StockFilterParams,
  PaginationMeta
} from '@/types/api'

export interface UseStocksOptions {
  autoLoad?: boolean
  pageSize?: number
}

export function useStocks(options: UseStocksOptions = {}) {
  const { autoLoad = false, pageSize = 20 } = options

  // State
  const stocks: Ref<StockWithAnalysis[]> = ref([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Pagination state
  const paginationMeta: Ref<PaginationMeta | null> = ref(null)
  const currentPage = ref(1)

  // Computed
  const hasStocks = computed(() => stocks.value.length > 0)
  const totalPages = computed(() => paginationMeta.value?.total_pages || 0)
  const totalItems = computed(() => paginationMeta.value?.total_items || 0)
  const hasNextPage = computed(() => paginationMeta.value?.has_next || false)
  const hasPreviousPage = computed(() => paginationMeta.value?.has_previous || false)

  // Methods
  const fetchStocks = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await stockApi.getAllStocks()
      stocks.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stocks'
      console.error('Error fetching stocks:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchStocksPaginated = async (params: PaginationParams = {}, filters: StockFilterParams = {}) => {
    loading.value = true
    error.value = null

    const paginationParams = {
      page: params.page || currentPage.value,
      page_size: params.page_size || pageSize,
      ...params
    }

    try {
      const response = await stockApi.getStocksPaginated(paginationParams, filters)
      stocks.value = response.data.data
      paginationMeta.value = response.data.meta
      currentPage.value = response.data.meta.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch paginated stocks'
      console.error('Error fetching paginated stocks:', err)
    } finally {
      loading.value = false
    }
  }

  const goToPage = async (page: number) => {
    if (page < 1 || page > totalPages.value) {
      return
    }

    await fetchStocksPaginated({ page, page_size: pageSize })
  }

  const nextPage = async () => {
    if (hasNextPage.value) {
      await goToPage(currentPage.value + 1)
    }
  }

  const previousPage = async () => {
    if (hasPreviousPage.value) {
      await goToPage(currentPage.value - 1)
    }
  }

  const refresh = async () => {
    if (paginationMeta.value) {
      await fetchStocksPaginated({ page: currentPage.value, page_size: pageSize })
    } else {
      await fetchStocks()
    }
  }

  // Auto-load data if enabled
  if (autoLoad) {
    fetchStocksPaginated()
  }

  return {
    // State
    stocks,
    loading,
    error,
    currentPage,
    paginationMeta,

    // Computed
    hasStocks,
    totalPages,
    totalItems,
    hasNextPage,
    hasPreviousPage,

    // Methods
    fetchStocks,
    fetchStocksPaginated,
    goToPage,
    nextPage,
    previousPage,
    refresh
  }
}

export function useStock(symbol: string) {
  const stock: Ref<StockWithAnalysis | null> = ref(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchStock = async () => {
    if (!symbol) return

    loading.value = true
    error.value = null

    try {
      const response = await stockApi.getStockBySymbol(symbol)
      stock.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stock'
      console.error('Error fetching stock:', err)
    } finally {
      loading.value = false
    }
  }

  const searchStock = async () => {
    if (!symbol) return

    loading.value = true
    error.value = null

    try {
      const response = await stockApi.searchStock(symbol)
      stock.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to search stock'
      console.error('Error searching stock:', err)
    } finally {
      loading.value = false
    }
  }

  const refreshStock = async () => {
    if (!symbol) return

    loading.value = true
    error.value = null

    try {
      await stockApi.refreshStockData(symbol)
      // Refresh the stock data after successful refresh
      await fetchStock()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to refresh stock data'
      console.error('Error refreshing stock:', err)
    } finally {
      loading.value = false
    }
  }

  return {
    // State
    stock,
    loading,
    error,

    // Methods
    fetchStock,
    searchStock,
    refreshStock
  }
}
