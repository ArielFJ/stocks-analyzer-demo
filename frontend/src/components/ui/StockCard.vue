<script setup lang="ts">
import { Calendar, Building2, Target, CheckCircle } from 'lucide-vue-next'
import VCard from './VCard.vue'
import VBadge from './VBadge.vue'

export interface StockPick {
  ticker: string
  company: string
  rating: string
  daysAgo: number
  brokerage: string
  priceTargetIncrease: number
  latestAction: string
  reasons: string[]
}

interface Props {
  pick: StockPick
  isTopPick?: boolean
}

const { pick, isTopPick = false } = defineProps<Props>()

const formatDaysAgo = (days: number) => {
  if (days === 0) return 'Today'
  if (days === 1) return '1 day ago'
  if (days < 30) return `${days} days ago`
  const months = Math.floor(days / 30)
  return months === 1 ? '1 month ago' : `${months} months ago`
}

const getTimelinessRating = (days: number) => {
  if (days <= 7) return 'Excellent'
  if (days <= 30) return 'Good'
  return 'Fair'
}
</script>

<template>
  <VCard
    :className="`${isTopPick ? 'border-primary/30 shadow-md' : ''} hover:shadow-md transition-shadow`"
  >
      <div class="flex items-start justify-between mb-4">
        <div>
          <div class="flex items-center gap-3 mb-2">
            <h3 class="text-xl font-bold">{{ pick.ticker }}</h3>
            <VBadge variant="secondary">{{ pick.rating }}</VBadge>
            <VBadge v-if="isTopPick" variant="outline" className="border-primary text-primary">
              Top Qualified Pick
            </VBadge>
          </div>
          <p class="text-muted-foreground">{{ pick.company }}</p>
        </div>
        <div class="text-right">
          <div class="flex items-center gap-1 mb-1">
            <Calendar class="w-4 h-4 text-muted-foreground" />
            <span class="text-sm text-muted-foreground">
              {{ formatDaysAgo(pick.daysAgo) }}
            </span>
          </div>
          <div class="flex items-center gap-1">
            <Building2 class="w-4 h-4 text-muted-foreground" />
            <span class="text-xs text-muted-foreground">
              {{
                pick.brokerage.length > 20
                  ? pick.brokerage.substring(0, 20) + '...'
                  : pick.brokerage
              }}
            </span>
          </div>
        </div>
      </div>

      <!-- Key Metrics -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <div class="bg-muted/50 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-1">
            <Target class="w-4 h-4 text-muted-foreground" />
            <span class="text-xs text-muted-foreground">Price Target Change</span>
          </div>
          <p
            :class="`font-semibold ${pick.priceTargetIncrease > 0 ? 'text-green-600' : 'text-muted-foreground'}`"
          >
            {{ pick.priceTargetIncrease > 0 ? '+' : '' }}{{ pick.priceTargetIncrease.toFixed(1) }}%
          </p>
        </div>

        <div class="bg-muted/50 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-1">
            <CheckCircle class="w-4 h-4 text-muted-foreground" />
            <span class="text-xs text-muted-foreground">Criteria Met</span>
          </div>
          <p class="font-semibold text-green-600">{{ pick.reasons.length }} of 6</p>
        </div>

        <div class="bg-muted/50 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-1">
            <Calendar class="w-4 h-4 text-muted-foreground" />
            <span class="text-xs text-muted-foreground">Timeliness</span>
          </div>
          <p class="font-semibold">{{ getTimelinessRating(pick.daysAgo) }}</p>
        </div>
      </div>

      <!-- Latest Action -->
      <div class="mb-4">
        <span class="text-sm text-muted-foreground">Latest Action: </span>
        <span class="text-sm font-medium">{{ pick.latestAction }}</span>
      </div>

      <!-- Selection Reasons -->
      <div>
        <h4 class="text-sm font-medium mb-2">Why This Stock Qualified</h4>
        <ul class="text-sm text-muted-foreground space-y-1">
          <li v-for="(reason, idx) in pick.reasons" :key="idx" class="flex items-start gap-2">
            <span class="text-green-600 mt-1">✓</span>
            <span>{{ reason }}</span>
          </li>
        </ul>
      </div>
  </VCard>
</template>
