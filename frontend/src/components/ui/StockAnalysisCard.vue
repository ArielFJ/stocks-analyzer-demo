<script setup lang="ts">
import { Calendar, Building2, Target, TrendingUp, Star } from 'lucide-vue-next'
import VCard from './VCard.vue'
import VBadge from './VBadge.vue'
import type { StockWithAnalysis } from '@/types/api'

interface Props {
  stock: StockWithAnalysis
  showFullAnalysis?: boolean
}

const { stock, showFullAnalysis = true } = defineProps<Props>()

// Helper functions
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now.getTime() - date.getTime())
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return '1 day ago'
  if (diffDays < 30) return `${diffDays} days ago`
  
  const diffMonths = Math.floor(diffDays / 30)
  return diffMonths === 1 ? '1 month ago' : `${diffMonths} months ago`
}

const extractPriceValue = (priceStr: string): number => {
  if (!priceStr) return 0
  const matches = priceStr.match(/[\d.]+/)
  return matches ? parseFloat(matches[0]) : 0
}

const calculatePriceTargetChange = (from: string, to: string): number => {
  const fromPrice = extractPriceValue(from)
  const toPrice = extractPriceValue(to)
  
  if (fromPrice && toPrice && fromPrice !== toPrice) {
    return ((toPrice - fromPrice) / fromPrice) * 100
  }
  return 0
}

const getRatingVariant = (rating: string): 'default' | 'secondary' | 'outline' => {
  const lowerRating = rating.toLowerCase()
  if (lowerRating.includes('buy') || lowerRating.includes('strong')) return 'default'
  if (lowerRating.includes('hold') || lowerRating.includes('neutral')) return 'secondary'
  return 'outline'
}

const getLatestAnalysis = () => {
  if (!stock.latest_analysis || stock.latest_analysis.length === 0) return null
  return stock.latest_analysis[0] // Most recent analysis
}

const latestAnalysis = getLatestAnalysis()
const priceTargetChange = latestAnalysis 
  ? calculatePriceTargetChange(latestAnalysis.target_from, latestAnalysis.target_to)
  : 0
</script>

<template>
  <VCard className="hover:shadow-md transition-shadow">
    <div class="flex items-start justify-between mb-4">
      <div class="flex-1">
        <div class="flex items-center gap-3 mb-2">
          <h3 class="text-xl font-bold">{{ stock.symbol }}</h3>
          <VBadge 
            v-if="latestAnalysis?.rating_to" 
            :variant="getRatingVariant(latestAnalysis.rating_to)"
          >
            {{ latestAnalysis.rating_to }}
          </VBadge>
          <Star v-if="priceTargetChange > 10" class="w-4 h-4 text-yellow-500" />
        </div>
        <p class="text-muted-foreground text-sm">{{ stock.name }}</p>
      </div>
      
      <div class="text-right text-sm">
        <div v-if="latestAnalysis" class="flex items-center gap-1 mb-1">
          <Calendar class="w-3 h-3 text-muted-foreground" />
          <span class="text-muted-foreground">
            {{ formatDate(latestAnalysis.analysis_date) }}
          </span>
        </div>
        <div v-if="latestAnalysis" class="flex items-center gap-1">
          <Building2 class="w-3 h-3 text-muted-foreground" />
          <span class="text-xs text-muted-foreground truncate max-w-[100px]">
            {{ latestAnalysis.brokerage }}
          </span>
        </div>
      </div>
    </div>

    <!-- Analysis Summary -->
    <div v-if="latestAnalysis && showFullAnalysis" class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-4">
      <!-- Price Target Change -->
      <div class="bg-muted/50 rounded-lg p-3">
        <div class="flex items-center gap-2 mb-1">
          <Target class="w-4 h-4 text-muted-foreground" />
          <span class="text-xs text-muted-foreground">Price Target</span>
        </div>
        <div class="space-y-1">
          <div v-if="priceTargetChange !== 0" class="flex items-center gap-1">
            <TrendingUp 
              :class="`w-3 h-3 ${priceTargetChange > 0 ? 'text-green-600' : 'text-red-600'}`" 
            />
            <span 
              :class="`text-sm font-semibold ${priceTargetChange > 0 ? 'text-green-600' : 'text-red-600'}`"
            >
              {{ priceTargetChange > 0 ? '+' : '' }}{{ priceTargetChange.toFixed(1) }}%
            </span>
          </div>
          <div class="text-xs text-muted-foreground">
            {{ latestAnalysis.target_from }} → {{ latestAnalysis.target_to }}
          </div>
        </div>
      </div>

      <!-- Rating Change -->
      <div class="bg-muted/50 rounded-lg p-3">
        <div class="flex items-center gap-2 mb-1">
          <Star class="w-4 h-4 text-muted-foreground" />
          <span class="text-xs text-muted-foreground">Rating</span>
        </div>
        <div class="space-y-1">
          <p class="text-sm font-semibold">{{ latestAnalysis.rating_to }}</p>
          <div v-if="latestAnalysis.rating_from !== latestAnalysis.rating_to" class="text-xs text-muted-foreground">
            From: {{ latestAnalysis.rating_from }}
          </div>
        </div>
      </div>
    </div>

    <!-- Latest Action -->
    <div v-if="latestAnalysis" class="mb-4">
      <span class="text-sm text-muted-foreground">Latest Action: </span>
      <span class="text-sm font-medium">{{ latestAnalysis.action }}</span>
    </div>

    <!-- Analysis Count -->
    <div v-if="stock.latest_analysis && stock.latest_analysis.length > 1" class="text-xs text-muted-foreground">
      {{ stock.latest_analysis.length }} recent analyst updates
    </div>

    <!-- No Analysis State -->
    <div v-else-if="!latestAnalysis" class="text-center py-4">
      <p class="text-sm text-muted-foreground">No recent analyst coverage</p>
    </div>
  </VCard>
</template>