<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900">Stock Recommendations</h1>
        <p class="mt-1 text-sm text-gray-600">
          AI-powered recommendations based on analyst sentiment and market data
        </p>
      </div>
      <button
        @click="refresh"
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

    <!-- Summary Stats -->
    <div v-if="hasRecommendations" class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white rounded-lg border p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-500">Total Recommendations</p>
            <p class="text-2xl font-semibold text-gray-900">{{ totalItems }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg border p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-500">High Confidence</p>
            <p class="text-2xl font-semibold text-gray-900">{{ highConfidenceRecommendations.length }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg border p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-500">Average Score</p>
            <p class="text-2xl font-semibold text-gray-900">
              {{ recommendations.length > 0 ? (recommendations.reduce((sum, rec) => sum + rec.score, 0) / recommendations.length).toFixed(1) : '0' }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Recommendations List -->
    <div v-if="hasRecommendations" class="space-y-4">
      <div class="space-y-4">
        <RecommendationCard
          v-for="(recommendation, index) in recommendations"
          :key="`${recommendation.stock.id}-${index}`"
          :recommendation="recommendation"
          :rank="paginationMeta ? (currentPage - 1) * paginationMeta.page_size + index + 1 : index + 1"
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
        @page-change="goToPage"
      />
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No recommendations available</h3>
      <p class="mt-1 text-sm text-gray-500">
        Recommendations will appear here once stock data is synced and analyzed.
      </p>
    </div>

    <!-- Loading State -->
    <div v-else class="text-center py-12">
      <svg class="animate-spin mx-auto h-8 w-8 text-blue-600" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="mt-2 text-sm text-gray-500">Loading recommendations...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRecommendations } from '@/composables/useRecommendations'
import RecommendationCard from '@/components/ui/RecommendationCard.vue'
import VPagination from '@/components/ui/VPagination.vue'

// Composables
const {
  recommendations,
  loading,
  error,
  currentPage,
  paginationMeta,
  hasRecommendations,
  totalPages,
  totalItems,
  hasNextPage,
  hasPreviousPage,
  highConfidenceRecommendations,
  fetchRecommendationsPaginated,
  goToPage,
  refresh
} = useRecommendations({ pageSize: 10 })

// Load initial data
onMounted(() => {
  fetchRecommendationsPaginated({ page: 1, page_size: 10 })
})
</script>