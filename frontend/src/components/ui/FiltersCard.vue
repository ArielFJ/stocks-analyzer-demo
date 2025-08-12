<script setup lang="ts">
import VCard from '@/components/ui/VCard.vue'
import VSelect from '@/components/ui/VSelect.vue'
import { SlidersHorizontal } from 'lucide-vue-next'
import { useFiltersStore } from '@/stores/filters'

// Use the Pinia store
const filtersStore = useFiltersStore()
</script>

<template>
  <VCard no-elevation content-class-name="!p-4">
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2">
        <SlidersHorizontal class="w-5 h-5 text-muted-foreground" />
        <p class="font-semibold text-muted-foreground">Filters:</p>
      </div>

      <!-- Action Type Filter -->
      <div class="w-64">
        <VSelect
          :options="filtersStore.actionTypeOptions"
          v-model="filtersStore.filters.actionType"
          placeholder="Select Action Type"
        />
      </div>
      
      <!-- Brokerage Filter -->
      <div class="w-64">
        <VSelect
          :options="filtersStore.brokerageOptions"
          v-model="filtersStore.filters.brokerage"
          placeholder="Select Brokerage"
        />
      </div>
      
      <!-- Sort By Filter -->
      <div class="w-64">
        <VSelect
          :options="filtersStore.sortByOptions"
          v-model="filtersStore.filters.sortBy"
          placeholder="Sort By"
        />
      </div>
      
      <!-- Clear Filters Button -->
      <button
        class="px-3 py-2 bg-background text-foreground border border-border text-sm rounded-md hover:bg-accent/90 hover:text-white transition-colors disabled:opacity-50"
        :disabled="!filtersStore.hasActiveFilters"
        @click="filtersStore.clearFilters"
      >
        Clear all ({{ filtersStore.activeFiltersCount }})
      </button>
    </div>
  </VCard>
</template>
