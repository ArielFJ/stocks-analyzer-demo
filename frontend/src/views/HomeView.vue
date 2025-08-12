<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRecommendations } from '@/composables/useRecommendations'
import { useApi } from '@/composables/useApi'
import DataOverview from '@/components/DataOverview.vue'
import QualitySelectionCriteria from '@/components/QualitySelectionCriteria.vue'
import { Trophy, Loader2, RefreshCw } from 'lucide-vue-next'
import VBadge from '@/components/ui/VBadge.vue'
import RecommendationCard from '@/components/ui/RecommendationCard.vue'

// Get top recommendations
const {
  recommendations,
  loading,
  error,
  hasRecommendations,
  fetchRecommendationsPaginated
} = useRecommendations()

const { loading: syncLoading, syncAllStocks } = useApi()
const successMessage = ref<string | null>(null)

const syncData = async () => {
  successMessage.value = null
  try {
    const result = await syncAllStocks()
    successMessage.value = result.message
    // Refresh recommendations after sync
    setTimeout(() => {
      fetchRecommendationsPaginated({ page: 1, page_size: 100 }) // Load all for home page
    }, 1000)
    
    // Clear success message after 5 seconds
    setTimeout(() => {
      successMessage.value = null
    }, 5000)
  } catch (err) {
    console.error('Sync failed:', err)
  }
}

onMounted(() => {
  fetchRecommendationsPaginated({ page: 1, page_size: 100 }) // Load all for home page
})

// Get top 5 recommendations for display (computed to be reactive)
const topRecommendations = computed(() => recommendations.value.slice(0, 5))
</script>

<template>
  <DataOverview />
  <QualitySelectionCriteria />

  <div class="flex justify-between items-center gap-3 w-full mt-8">
    <div class="flex items-center gap-2">
      <Trophy class="w-5 h-5 text-primary" />
      <div>
        <h1 class="text-xl font-semibold">Qualified Investment Opportunities</h1>
      </div>
    </div>
    <div class="flex items-center gap-3">
      <button
        @click="syncData"
        :disabled="syncLoading"
        class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <RefreshCw v-if="!syncLoading" class="w-4 h-4 mr-2" />
        <Loader2 v-else class="w-4 h-4 mr-2 animate-spin" />
        {{ syncLoading ? 'Syncing...' : 'Sync Data' }}
      </button>
      <VBadge variant="outline"> 
        {{ hasRecommendations ? `${topRecommendations.length} of ${recommendations.length}` : 'Loading...' }}
      </VBadge>
    </div>
  </div>

  <!-- Success Message -->
  <div v-if="successMessage" class="rounded-md bg-green-50 p-4 mt-4">
    <div class="flex">
      <svg class="h-5 w-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <div class="ml-3">
        <p class="text-sm font-medium text-green-800">{{ successMessage }}</p>
      </div>
    </div>
  </div>

  <!-- Loading State -->
  <div v-if="loading" class="flex items-center justify-center py-12">
    <Loader2 class="w-6 h-6 animate-spin text-blue-600" />
    <span class="ml-2 text-sm text-gray-600">Loading recommendations...</span>
  </div>

  <!-- Error State -->
  <div v-else-if="error" class="rounded-md bg-red-50 p-4 mt-4">
    <div class="flex">
      <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
      </svg>
      <div class="ml-3">
        <h3 class="text-sm font-medium text-red-800">Unable to load recommendations</h3>
        <p class="mt-1 text-sm text-red-700">{{ error }}</p>
      </div>
    </div>
  </div>

  <!-- Recommendations List -->
  <div v-else-if="hasRecommendations" class="grid grid-cols-1 gap-4 mt-4">
    <RecommendationCard
      v-for="(recommendation, index) in topRecommendations"
      :key="`home-${recommendation.stock.id}`"
      :recommendation="recommendation"
      :rank="index + 1"
    />
    
    <!-- Show more link -->
    <div class="text-center mt-6">
      <router-link
        to="/recommendations"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-blue-600 bg-blue-50 hover:bg-blue-100 transition-colors"
      >
        View All Recommendations
        <svg class="ml-2 w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </router-link>
    </div>
  </div>

  <!-- Empty State -->
  <div v-else class="text-center py-12">
    <Trophy class="mx-auto h-12 w-12 text-gray-400" />
    <h3 class="mt-2 text-sm font-medium text-gray-900">No recommendations available</h3>
    <p class="mt-1 text-sm text-gray-500">
      Recommendations will appear here once stock data is synced from the API.
    </p>
    <div class="mt-6">
      <router-link
        to="/stocks"
        class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
      >
        View Stocks
      </router-link>
    </div>
  </div>
</template>
