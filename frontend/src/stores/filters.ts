import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { stockApi } from '@/services/stockApi'
import type { FilterOption } from '@/types/api'

export interface FilterOptions {
  actionType: string
  brokerage: string
  sortBy: string
}

export const useFiltersStore = defineStore('filters', () => {
  // State (reactive data)
  const filters = ref<FilterOptions>({
    actionType: 'all',
    brokerage: 'all',
    sortBy: 'newest'
  })

  // Dynamic filter options from API
  const actionTypeOptions = ref<FilterOption[]>([
    { label: 'All actions', value: 'all' }
  ])

  const brokerageOptions = ref<FilterOption[]>([
    { label: 'All brokerages', value: 'all' }
  ])

  const sortByOptions = ref<FilterOption[]>([
    { label: 'Newest', value: 'newest' }
  ])

  // Loading state
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Fetch filter options from API
  const fetchFilterOptions = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await stockApi.getFilterOptions()
      actionTypeOptions.value = response.data.action_types
      brokerageOptions.value = response.data.brokerages
      sortByOptions.value = response.data.sort_by
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch filter options'
      console.error('Error fetching filter options:', err)
    } finally {
      loading.value = false
    }
  }

  // Getters (computed properties)
  const activeFiltersCount = computed(() => {
    let count = 0
    if (filters.value.actionType !== 'all') count++
    if (filters.value.brokerage !== 'all') count++
    if (filters.value.sortBy !== 'newest') count++
    return count
  })

  const hasActiveFilters = computed(() => activeFiltersCount.value > 0)

  const currentActionTypeLabel = computed(() =>
    actionTypeOptions.value.find(option => option.value === filters.value.actionType)?.label || 'All actions'
  )

  const currentBrokerageLabel = computed(() =>
    brokerageOptions.value.find(option => option.value === filters.value.brokerage)?.label || 'All brokerages'
  )

  const currentSortByLabel = computed(() =>
    sortByOptions.value.find(option => option.value === filters.value.sortBy)?.label || 'Newest'
  )

  // Actions (methods that modify state)
  const setActionType = (value: string) => {
    filters.value.actionType = value
  }

  const setBrokerage = (value: string) => {
    filters.value.brokerage = value
  }

  const setSortBy = (value: string) => {
    filters.value.sortBy = value
  }

  const updateFilters = (newFilters: Partial<FilterOptions>) => {
    Object.assign(filters.value, newFilters)
  }

  const clearFilters = () => {
    filters.value = {
      actionType: 'all',
      brokerage: 'all',
      sortBy: 'newest'
    }
  }

  const resetToDefaults = () => {
    clearFilters()
  }

  // Filter logic functions
  const filterTickers = <T extends { analystAction: string; brokerage: string; date: string; symbol: string; companyName: string; priceTarget: string }>(tickers: T[]): T[] => {
    let filtered = [...tickers]

    // Filter by action type
    if (filters.value.actionType !== 'all') {
      filtered = filtered.filter(ticker => {
        const action = ticker.analystAction?.toLowerCase() || ''
        switch (filters.value.actionType) {
          case 'upgrade':
            return action.includes('buy') || action.includes('raised') || action.includes('upgrade')
          case 'downgrade':
            return action.includes('sell') || action.includes('downgrade') || action.includes('cut')
          case 'neutral':
            return !action.includes('buy') && !action.includes('raised') && !action.includes('upgrade') &&
                   !action.includes('sell') && !action.includes('downgrade') && !action.includes('cut')
          default:
            return true
        }
      })
    }

    // Filter by brokerage
    if (filters.value.brokerage !== 'all') {
      filtered = filtered.filter(ticker => {
        const brokerage = ticker.brokerage?.toLowerCase().replace(/\s+/g, '-') || ''
        return brokerage === filters.value.brokerage
      })
    }

    return filtered
  }

  const sortTickers = <T extends { date: string; symbol: string; companyName: string; priceTarget: string }>(tickers: T[]): T[] => {
    const sorted = [...tickers]

    switch (filters.value.sortBy) {
      case 'newest':
        return sorted.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
      case 'oldest':
        return sorted.sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())
      case 'ticker-a-z':
        return sorted.sort((a, b) => a.symbol.localeCompare(b.symbol))
      case 'company-a-z':
        return sorted.sort((a, b) => a.companyName.localeCompare(b.companyName))
      case 'price-high-low':
        return sorted.sort((a, b) => {
          const priceA = parseFloat(a.priceTarget.replace(/[^0-9.-]+/g, '')) || 0
          const priceB = parseFloat(b.priceTarget.replace(/[^0-9.-]+/g, '')) || 0
          return priceB - priceA
        })
      case 'price-low-high':
        return sorted.sort((a, b) => {
          const priceA = parseFloat(a.priceTarget.replace(/[^0-9.-]+/g, '')) || 0
          const priceB = parseFloat(b.priceTarget.replace(/[^0-9.-]+/g, '')) || 0
          return priceA - priceB
        })
      default:
        return sorted
    }
  }

  const processTickerData = <T extends { analystAction: string; brokerage: string; date: string; symbol: string; companyName: string; priceTarget: string }>(tickers: T[]): T[] => {
    const filtered = filterTickers(tickers)
    return sortTickers(filtered)
  }

  // Return all reactive state and methods
  return {
    // State
    filters,
    loading,
    error,

    // Options
    actionTypeOptions,
    brokerageOptions,
    sortByOptions,

    // Getters
    activeFiltersCount,
    hasActiveFilters,
    currentActionTypeLabel,
    currentBrokerageLabel,
    currentSortByLabel,

    // Actions
    setActionType,
    setBrokerage,
    setSortBy,
    updateFilters,
    clearFilters,
    resetToDefaults,
    fetchFilterOptions,

    // Utility functions
    filterTickers,
    sortTickers,
    processTickerData
  }
})
