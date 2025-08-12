<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  value: string | number
  label?: string
  subtitle?: string
  variant?: 'default' | 'compact' | 'detailed'
  colorClass?: string
  withBorder?: boolean
  loading?: boolean
  error?: boolean
  formatValue?: (value: string | number) => string
}

const props = withDefaults(defineProps<Props>(), {
  value: '',
  label: '',
  subtitle: '',
  variant: 'default',
  colorClass: 'text-primary',
  withBorder: false,
  loading: false,
  error: false,
  formatValue: (value: string | number) => String(value),
})

// Computed properties for better maintainability
const displayValue = computed(() => {
  if (props.loading) return '...'
  if (props.error) return 'Error'
  return props.formatValue(props.value)
})

const containerClasses = computed(() => ({
  'py-5 px-3 bg-muted/30 rounded-lg transition-colors duration-200 relative': true,
  '!bg-card': props.variant === 'detailed' || !!props.subtitle,
  'border border-border': props.withBorder,
  'border-destructive/50 bg-destructive/5': props.error,
  'animate-pulse': props.loading,
}))

const valueClasses = computed(() => [
  'font-bold transition-colors duration-200',
  props.error ? 'text-destructive' : props.colorClass,
  props.value.toString().length > 1 ? 'text-lg' : 'text-2xl',
])

// Determine layout based on variant or subtitle presence
const isCompactLayout = computed(
  () => props.variant === 'compact' || (props.variant === 'default' && !!props.subtitle),
)

// Generate accessible labels
const ariaLabel = computed(() => {
  const parts = [props.label, props.subtitle, displayValue.value].filter(Boolean)
  return parts.join(', ')
})
</script>

<template>
  <div
    :class="containerClasses"
    :aria-label="ariaLabel"
    role="group"
    :aria-live="loading ? 'polite' : undefined"
  >
    <!-- Compact Layout (when subtitle exists or variant is compact) -->
    <div v-if="isCompactLayout" class="flex items-center justify-between">
      <div class="grid gap-1 min-w-0 flex-1">
        <h4 v-if="label" class="text-sm font-semibold text-foreground truncate">
          {{ label }}
        </h4>
        <p v-if="subtitle" class="text-xs text-muted-foreground line-clamp-2">
          {{ subtitle }}
        </p>
      </div>

      <div class="flex-shrink-0 ml-3">
        <p :class="valueClasses" :title="String(value)">
          {{ displayValue }}
        </p>
      </div>
    </div>

    <!-- Centered Layout (default when no subtitle) -->
    <div v-else class="flex flex-col items-center justify-center text-center">
      <p :class="valueClasses" :title="String(value)">
        {{ displayValue }}
      </p>

      <h4 v-if="label" class="text-xs text-muted-foreground mt-1">
        {{ label }}
      </h4>
    </div>

    <!-- Loading indicator -->
    <div
      v-if="loading"
      class="absolute inset-0 left-0 right-0 flex items-center justify-center bg-background/50 rounded-lg"
    >
      <div class="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin" />
    </div>
  </div>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media (prefers-reduced-motion: reduce) {
  .transition-colors,
  .animate-pulse,
  .animate-spin {
    animation: none;
    transition: none;
  }
}
</style>
