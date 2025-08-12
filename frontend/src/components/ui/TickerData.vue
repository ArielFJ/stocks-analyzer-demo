<script setup lang="ts">
import VBadge from '@/components/ui/VBadge.vue'
import VCard from '@/components/ui/VCard.vue'
import { ArrowRight, TrendingUp, TrendingDown } from 'lucide-vue-next'
import { computed } from 'vue'

export interface TickerData {
  symbol: string
  companyName: string
  analystAction: string
  brokerage: string
  priceTarget: string
  previousPriceTarget: string
  date: string
  actionType?: 'upgrade' | 'downgrade' | 'neutral'
  previousAction?: string
  actionText?: string
}

const props = defineProps<{
  ticker: TickerData
}>()

const getActionType = computed(() => {
  if (props.ticker.actionType) return props.ticker.actionType

  const action = props.ticker.analystAction?.toLowerCase() || ''
  if (action.includes('buy') || action.includes('raised') || action.includes('upgrade')) {
    return 'upgrade'
  } else if (action.includes('sell') || action.includes('downgrade') || action.includes('cut')) {
    return 'downgrade'
  }
  return 'neutral'
})

const getActionText = computed(() => {
  if (props.ticker.actionText) return props.ticker.actionText

  const type = getActionType.value
  if (type === 'upgrade') return `target raised by ${props.ticker.brokerage}`
  if (type === 'downgrade') return `downgraded by ${props.ticker.brokerage}`
  return `initiated by ${props.ticker.brokerage}`
})

const getBadgeVariant = computed(() => {
  const type = getActionType.value
  if (type === 'upgrade') return 'success'
  if (type === 'downgrade') return 'destructive'
  return 'secondary'
})
</script>

<template>
  <VCard :content-class-name="`!p-4`" no-elevation>
    <div class="grid grid-cols-2 gap-2">
      <div class="flex flex-col">
        <div class="flex items-center gap-2">
          <h2 class="text-lg font-semibold">{{ ticker.symbol }}</h2>
          <TrendingUp v-if="getActionType === 'upgrade'" class="w-4 h-4 text-green-600" />
          <TrendingDown v-else-if="getActionType === 'downgrade'" class="w-4 h-4 text-red-600" />
        </div>
        <p class="text-sm text-muted-foreground">{{ ticker.companyName }}</p>
        <p
          class="text-sm font-medium my-2"
          :class="{
            'text-green-600': getActionType === 'upgrade',
            'text-red-600': getActionType === 'downgrade',
            'text-foreground': getActionType === 'neutral',
          }"
        >
          {{ getActionText }}
        </p>
        <div class="flex items-center">
          <VBadge :variant="getActionType === 'neutral' ? 'secondary' : 'outline'">
            {{ ticker.previousAction || ticker.analystAction }}
          </VBadge>
          <ArrowRight class="w-4 text-muted-foreground mx-2" />
          <VBadge :variant="getBadgeVariant">{{ ticker.analystAction }}</VBadge>
        </div>
      </div>
      <div class="flex flex-col items-end">
        <p class="text-xs text-muted-foreground">{{ ticker.date }}</p>
        <p class="text-xs text-muted-foreground">{{ ticker.brokerage }}</p>
      </div>

      <div class="flex items-center">
        <p class="text-sm text-muted-foreground">Price Target:</p>
      </div>
      <div class="flex items-center justify-end">
        <p class="text-sm font-semibold">{{ ticker.previousPriceTarget }}</p>
        <ArrowRight class="w-4 text-muted-foreground mx-2" />
        <p class="text-sm font-semibold">
          {{ ticker.priceTarget }}
        </p>
      </div>
    </div>
  </VCard>
</template>
