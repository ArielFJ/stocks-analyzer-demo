import { ref, type Ref } from 'vue'
import { stockApi } from '@/services/stockApi'
import type { MarketIntelligenceOverview } from '@/types/api'

export function useAnalytics() {
  // State
  const analytics: Ref<MarketIntelligenceOverview | null> = ref(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Methods
  const fetchAnalytics = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await stockApi.getMarketIntelligenceOverview()
      analytics.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch analytics'
      console.error('Error fetching analytics:', err)
    } finally {
      loading.value = false
    }
  }

  const refresh = async () => {
    await fetchAnalytics()
  }

  return {
    // State
    analytics,
    loading,
    error,

    // Methods
    fetchAnalytics,
    refresh
  }
}