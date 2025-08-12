<script setup lang="ts">
import { TrendingUp, Award, Calendar, Building2, Target, Star } from 'lucide-vue-next'
import VCard from './VCard.vue'
import VBadge from './VBadge.vue'
import type { StockRecommendation } from '@/types/api'

interface Props {
  recommendation: StockRecommendation
  rank?: number
}

const { recommendation, rank } = defineProps<Props>()

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

const getConfidenceVariant = (confidence: string): 'default' | 'secondary' | 'outline' => {
  const lowerConfidence = confidence.toLowerCase()
  if (lowerConfidence === 'high') return 'default'
  if (lowerConfidence === 'medium') return 'secondary'
  return 'outline'
}

const getScoreColor = (score: number): string => {
  if (score >= 80) return 'text-green-600'
  if (score >= 60) return 'text-blue-600'
  if (score >= 40) return 'text-yellow-600'
  return 'text-gray-600'
}

const getScoreBackground = (score: number): string => {
  if (score >= 80) return 'bg-green-50'
  if (score >= 60) return 'bg-blue-50'
  if (score >= 40) return 'bg-yellow-50'
  return 'bg-gray-50'
}

const latestAnalysis = recommendation.stock.latest_analysis?.[0]
</script>

<template>
  <VCard className="hover:shadow-lg transition-all duration-200 border-l-4 border-l-blue-500">
    <div class="flex items-start justify-between mb-4">
      <!-- Stock Info -->
      <div class="flex-1">
        <div class="flex items-center gap-3 mb-2">
          <div v-if="rank" class="flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-800 rounded-full text-sm font-bold">
            {{ rank }}
          </div>
          <h3 class="text-xl font-bold">{{ recommendation.stock.symbol }}</h3>
          <VBadge 
            v-if="latestAnalysis?.rating_to" 
            variant="default"
          >
            {{ latestAnalysis.rating_to }}
          </VBadge>
          <VBadge 
            :variant="getConfidenceVariant(recommendation.confidence)"
          >
            {{ recommendation.confidence }} Confidence
          </VBadge>
        </div>
        <p class="text-gray-600 text-sm mb-1">{{ recommendation.stock.name }}</p>
        <p class="text-gray-500 text-sm">{{ recommendation.reason }}</p>
      </div>
      
      <!-- Score -->
      <div class="text-center">
        <div 
          :class="`inline-flex items-center justify-center w-16 h-16 rounded-full ${getScoreBackground(recommendation.score)}`"
        >
          <div class="text-center">
            <div :class="`text-lg font-bold ${getScoreColor(recommendation.score)}`">
              {{ Math.round(recommendation.score) }}
            </div>
            <div class="text-xs text-gray-500">Score</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Analysis Details -->
    <div v-if="latestAnalysis" class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
      <!-- Latest Analysis Date -->
      <div class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
        <Calendar class="w-5 h-5 text-gray-400" />
        <div>
          <p class="text-xs text-gray-500 uppercase tracking-wide">Last Updated</p>
          <p class="text-sm font-medium">{{ formatDate(latestAnalysis.analysis_date) }}</p>
        </div>
      </div>

      <!-- Brokerage -->
      <div class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
        <Building2 class="w-5 h-5 text-gray-400" />
        <div>
          <p class="text-xs text-gray-500 uppercase tracking-wide">Analyst</p>
          <p class="text-sm font-medium truncate">{{ latestAnalysis.brokerage }}</p>
        </div>
      </div>

      <!-- Price Target -->
      <div v-if="latestAnalysis.target_to" class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
        <Target class="w-5 h-5 text-gray-400" />
        <div>
          <p class="text-xs text-gray-500 uppercase tracking-wide">Price Target</p>
          <p class="text-sm font-medium">{{ latestAnalysis.target_to }}</p>
          <p v-if="latestAnalysis.target_from !== latestAnalysis.target_to" class="text-xs text-gray-400">
            From: {{ latestAnalysis.target_from }}
          </p>
        </div>
      </div>
    </div>

    <!-- Action and Analysis Count -->
    <div class="flex items-center justify-between pt-4 border-t border-gray-100">
      <div v-if="latestAnalysis" class="flex items-center gap-2">
        <TrendingUp class="w-4 h-4 text-blue-500" />
        <span class="text-sm text-gray-600">Latest Action:</span>
        <span class="text-sm font-medium">{{ latestAnalysis.action }}</span>
      </div>
      
      <div v-if="recommendation.stock.latest_analysis" class="flex items-center gap-1 text-xs text-gray-500">
        <Award class="w-3 h-3" />
        <span>{{ recommendation.stock.latest_analysis.length }} analyst update{{ recommendation.stock.latest_analysis.length !== 1 ? 's' : '' }}</span>
      </div>
    </div>

    <!-- Top Recommendation Indicator -->
    <div v-if="rank === 1" class="absolute top-4 right-4">
      <div class="flex items-center gap-1 px-2 py-1 bg-yellow-100 text-yellow-800 text-xs font-medium rounded-full">
        <Star class="w-3 h-3" />
        Top Pick
      </div>
    </div>
  </VCard>
</template>