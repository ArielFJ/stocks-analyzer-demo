<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { Database, Filter } from 'lucide-vue-next'
import DataCard from './ui/DataCard.vue'
import CardContainer from './ui/CardContainer.vue'
import { useAnalytics } from '@/composables/useAnalytics'

// Get analytics data from dedicated endpoint
const { analytics, fetchAnalytics } = useAnalytics()

// Computed values for template
const analyticsData = computed(() => {
  if (!analytics.value) {
    return {
      totalStocks: 0,
      totalRecommendations: 0,
      recentAnalysis: 0,
      upgrades: 0,
      downgrades: 0,
      highConfidenceRecs: 0,
      selectionRate: 0
    }
  }
  
  return {
    totalStocks: analytics.value.total_stocks,
    totalRecommendations: analytics.value.total_recommendations,
    recentAnalysis: analytics.value.recent_analysis,
    upgrades: analytics.value.upgrades,
    downgrades: analytics.value.downgrades,
    highConfidenceRecs: analytics.value.high_confidence_recs,
    selectionRate: Math.round(analytics.value.selection_rate)
  }
})

onMounted(() => {
  fetchAnalytics()
})
</script>

<template>
  <CardContainer label="Market Intelligence Overview" :icon="Database">
    <div class="grid grid-cols-2 gap-4 mt-4">
      <div>
        <h3 class="text-sm font-semibold text-muted-foreground mb-2">Raw Analytics Data</h3>

        <div class="grid grid-cols-2 gap-2">
          <DataCard 
            :value="analyticsData.totalStocks.toString()" 
            label="Total Stocks" 
            colorClass="text-primary" 
          />
          <DataCard 
            :value="analyticsData.recentAnalysis.toString()" 
            label="Recent Analysis" 
            colorClass="text-accent" 
          />
          <DataCard 
            :value="analyticsData.upgrades.toString()" 
            label="Upgrades" 
            colorClass="text-success" 
          />
          <DataCard 
            :value="analyticsData.downgrades.toString()" 
            label="Downgrades" 
            colorClass="text-destructive" 
          />
        </div>
      </div>

      <div>
        <div class="flex items-center mb-2">
          <Filter class="w-4 h-4 text-muted-foreground mr-1" />
          <h3 class="text-sm font-semibold text-muted-foreground">AI-Powered Recommendations</h3>
        </div>

        <div class="grid grid-cols-1 gap-2">
          <DataCard
            :value="analyticsData.totalRecommendations.toString()"
            label="Quality Recommendations"
            :subtitle="`Based on ${analyticsData.totalStocks} stocks analyzed`"
            colorClass="text-success"
            withBorder
          />
          <DataCard
            :value="analyticsData.highConfidenceRecs.toString()"
            label="High Confidence"
            subtitle="Top-tier investment opportunities"
            colorClass="text-primary"
          />
        </div>
      </div>
    </div>
  </CardContainer>
</template>
