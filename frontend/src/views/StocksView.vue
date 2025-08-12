<template>
  <div class="space-y-6">
    <!-- Header with Actions -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900">Stock Analysis</h1>
        <p class="mt-1 text-sm text-gray-600">
          Real-time analyst coverage and recommendations
        </p>
      </div>
      <div class="flex gap-3">
        <button
          @click="syncStocks"
          :disabled="syncLoading"
          class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg v-if="syncLoading" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ syncLoading ? 'Syncing...' : 'Sync Data' }}
        </button>
        <button
          @click="refreshWithFilters"
          :disabled="loading"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <!-- Filters -->
    <FiltersCard />

    <!-- Error Message -->
    <div v-if="error" class="rounded-md bg-red-50 p-4">
      <div class="flex">
        <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error</h3>
          <p class="mt-1 text-sm text-red-700">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Success Message -->
    <div v-if="successMessage" class="rounded-md bg-green-50 p-4">
      <div class="flex">
        <svg class="h-5 w-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <p class="text-sm font-medium text-green-800">{{ successMessage }}</p>
        </div>
      </div>
    </div>

    <!-- Stocks Grid -->
    <div v-if="hasStocks" class="space-y-4">
      <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        <StockAnalysisCard
          v-for="stock in stocks"
          :key="stock.id"
          :stock="stock"
        />
      </div>

      <!-- Pagination -->
      <VPagination
        v-if="paginationMeta"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="totalItems"
        :page-size="paginationMeta.page_size"
        :has-next-page="hasNextPage"
        :has-previous-page="hasPreviousPage"
        @page-change="goToPageWithFilters"
      />
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No stocks found</h3>
      <p class="mt-1 text-sm text-gray-500">
        Get started by syncing data from the external API.
      </p>
      <div class="mt-6">
        <button
          @click="syncStocks"
          :disabled="syncLoading"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {{ syncLoading ? 'Syncing...' : 'Sync Stock Data' }}
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-else class="text-center py-12">
      <svg class="animate-spin mx-auto h-8 w-8 text-blue-600" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="mt-2 text-sm text-gray-500">Loading stocks...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useStocks } from '@/composables/useStocks'
import { useApi } from '@/composables/useApi'
import { useFiltersStore } from '@/stores/filters'
import StockAnalysisCard from '@/components/ui/StockAnalysisCard.vue'
import VPagination from '@/components/ui/VPagination.vue'
import FiltersCard from '@/components/ui/FiltersCard.vue'
import type { StockFilterParams } from '@/types/api'

// Composables
const {
  stocks,
  loading,
  error,
  currentPage,
  paginationMeta,
  hasStocks,
  totalPages,
  totalItems,
  hasNextPage,
  hasPreviousPage,
  fetchStocksPaginated,
  goToPage,
  refresh
} = useStocks({ pageSize: 12 })

const { loading: syncLoading, syncAllStocks } = useApi()
const filtersStore = useFiltersStore()

// Local state
const successMessage = ref<string | null>(null)

// Convert filter store values to API format
const getApiFilters = (): StockFilterParams => {
  return {
    action_type: filtersStore.filters.actionType === 'all' ? '' : filtersStore.filters.actionType,
    brokerage: filtersStore.filters.brokerage === 'all' ? '' : filtersStore.filters.brokerage,
    sort_by: filtersStore.filters.sortBy
  }
}

// Watch for filter changes and refetch data
watch(() => filtersStore.filters, () => {
  fetchStocksPaginated({ page: 1, page_size: 12 }, getApiFilters())
}, { deep: true })

// Methods
const syncStocks = async () => {
  successMessage.value = null
  try {
    const result = await syncAllStocks()
    successMessage.value = result.message
    // Refresh the stocks data after successful sync
    setTimeout(() => {
      refreshWithFilters()
    }, 1000)
    
    // Clear success message after 5 seconds
    setTimeout(() => {
      successMessage.value = null
    }, 5000)
  } catch (err) {
    console.error('Sync failed:', err)
  }
}

// Enhanced goToPage that preserves filters
const goToPageWithFilters = async (page: number) => {
  await fetchStocksPaginated({ page, page_size: 12 }, getApiFilters())
}

// Enhanced refresh that preserves filters  
const refreshWithFilters = async () => {
  await fetchStocksPaginated({ page: currentPage.value, page_size: 12 }, getApiFilters())
}

// Load initial data
onMounted(async () => {
  // Fetch filter options first
  await filtersStore.fetchFilterOptions()
  // Then fetch stocks with current filters
  fetchStocksPaginated({ page: 1, page_size: 12 }, getApiFilters())
})
</script>