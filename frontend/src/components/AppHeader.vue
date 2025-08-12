<script setup lang="ts">
import { ref } from 'vue'
import { BarChart3, Search } from 'lucide-vue-next'
import VBadge from './ui/VBadge.vue'
import VInput from './ui/VInput.vue'

// import { RouterLink } from 'vue-router'
const filteredCount = 10
const totalCount = 100

const { currentTab } = defineProps({
  currentTab: String,
  // searchQuery: String,
  onSearchChange: Function,
})

const searchQuery = ref('')

const getBadgeText = () => {
  if (currentTab === 'raw-data') {
    return `${filteredCount} of ${totalCount} analyst reports`
  }
  return `${totalCount} total analyst reports`
}
</script>

<template>
  <header class="border-b border-border bg-card">
    <div class="container mx-auto px-4 py-4">
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <!-- App Title -->
        <div class="flex items-center gap-3 w-full">
          <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <BarChart3 class="w-5 h-5 text-primary-foreground" />
          </div>
          <div>
            <h1 class="text-xl font-semibold">Stock Analyst Tracker</h1>
            <p class="text-sm text-muted-foreground">
              <span v-if="currentTab === 'manual'">Quality-focused stock selection</span>
              <span v-else-if="currentTab === 'algorithm'">AI-powered investment analysis</span>
              <span v-else-if="currentTab === 'raw-data'">Raw analyst recommendation data</span>
              <span v-else>Investment recommendations & analysis</span>
            </p>
          </div>
        </div>

        <!-- Search and Stats -->
        <div
          v-show="currentTab === 'raw-data'"
          class="flex flex-col sm:flex-row gap-3 sm:items-center"
        >
          <VBadge variant="outline">
            {{ getBadgeText() }}
          </VBadge>

          <div class="relative">
            <Search
              class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4"
            />
            <VInput
              placeholder="Search by ticker or company..."
              v-model="searchQuery"
              classes="pl-10 w-full sm:w-64"
            />
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
