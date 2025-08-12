import { ref, computed, type Ref } from 'vue'
import { stockApi } from '@/services/stockApi'
import type {
  StockRecommendation,
  PaginationParams,
  PaginationMeta
} from '@/types/api'

export interface UseRecommendationsOptions {
  autoLoad?: boolean
  pageSize?: number
}

export function useRecommendations(options: UseRecommendationsOptions = {}) {
  const { autoLoad = false, pageSize = 20 } = options

  // State
  const recommendations: Ref<StockRecommendation[]> = ref([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  
  // Pagination state
  const paginationMeta: Ref<PaginationMeta | null> = ref(null)
  const currentPage = ref(1)

  // Computed
  const hasRecommendations = computed(() => recommendations.value.length > 0)
  const totalPages = computed(() => paginationMeta.value?.total_pages || 0)
  const totalItems = computed(() => paginationMeta.value?.total_items || 0)
  const hasNextPage = computed(() => paginationMeta.value?.has_next || false)
  const hasPreviousPage = computed(() => paginationMeta.value?.has_previous || false)

  // Get top recommendations (sorted by score)
  const topRecommendations = computed(() => {
    return [...recommendations.value].sort((a, b) => b.score - a.score)
  })

  // Get high confidence recommendations
  const highConfidenceRecommendations = computed(() => {
    return recommendations.value.filter(rec => rec.confidence.toLowerCase() === 'high')
  })

  // Methods
  const fetchRecommendations = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await stockApi.getAllRecommendations()
      recommendations.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch recommendations'
      console.error('Error fetching recommendations:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchRecommendationsPaginated = async (params: PaginationParams = {}) => {
    loading.value = true
    error.value = null

    const paginationParams = {
      page: params.page || currentPage.value,
      page_size: params.page_size || pageSize,
      ...params
    }

    try {
      const response = await stockApi.getRecommendationsPaginated(paginationParams)
      recommendations.value = response.data.data
      paginationMeta.value = response.data.meta
      currentPage.value = response.data.meta.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch paginated recommendations'
      console.error('Error fetching paginated recommendations:', err)
    } finally {
      loading.value = false
    }
  }

  const goToPage = async (page: number) => {
    if (page < 1 || page > totalPages.value) {
      return
    }
    
    await fetchRecommendationsPaginated({ page, page_size: pageSize })
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
      await fetchRecommendationsPaginated({ page: currentPage.value, page_size: pageSize })
    } else {
      await fetchRecommendations()
    }
  }

  // Auto-load data if enabled
  if (autoLoad) {
    fetchRecommendationsPaginated()
  }

  return {
    // State
    recommendations,
    loading,
    error,
    currentPage,
    paginationMeta,

    // Computed
    hasRecommendations,
    totalPages,
    totalItems,
    hasNextPage,
    hasPreviousPage,
    topRecommendations,
    highConfidenceRecommendations,

    // Methods
    fetchRecommendations,
    fetchRecommendationsPaginated,
    goToPage,
    nextPage,
    previousPage,
    refresh
  }
}