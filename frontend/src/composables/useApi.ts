import { ref } from 'vue'
import { stockApi } from '@/services/stockApi'

export function useApi() {
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Health check
  const checkHealth = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await stockApi.checkHealth()
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Health check failed'
      console.error('Error checking health:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Sync all stocks
  const syncAllStocks = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await stockApi.syncAllStocks()
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to sync stocks'
      console.error('Error syncing stocks:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    // State
    loading,
    error,

    // Methods
    checkHealth,
    syncAllStocks
  }
}