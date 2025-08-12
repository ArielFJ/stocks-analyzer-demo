<template>
  <div v-if="totalPages > 1" class="flex items-center justify-between gap-4">
    <!-- Previous Button -->
    <button
      @click="goToPrevious"
      :disabled="!hasPreviousPage"
      class="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white disabled:hover:text-gray-500"
    >
      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      Previous
    </button>

    <!-- Page Numbers -->
    <div class="flex items-center gap-1">
      <!-- First Page -->
      <template v-if="showFirstPage">
        <button
          @click="goToPage(1)"
          :class="[
            'px-3 py-2 text-sm font-medium rounded-lg transition-colors',
            currentPage === 1
              ? 'bg-blue-600 text-white'
              : 'text-gray-500 bg-white border border-gray-300 hover:bg-gray-50 hover:text-gray-700'
          ]"
        >
          1
        </button>
        <span v-if="showStartEllipsis" class="px-2 py-2 text-gray-500">...</span>
      </template>

      <!-- Page Numbers Around Current Page -->
      <template v-for="page in visiblePages" :key="page">
        <button
          @click="goToPage(page)"
          :class="[
            'px-3 py-2 text-sm font-medium rounded-lg transition-colors',
            currentPage === page
              ? 'bg-blue-600 text-white'
              : 'text-gray-500 bg-white border border-gray-300 hover:bg-gray-50 hover:text-gray-700'
          ]"
        >
          {{ page }}
        </button>
      </template>

      <!-- Last Page -->
      <template v-if="showLastPage">
        <span v-if="showEndEllipsis" class="px-2 py-2 text-gray-500">...</span>
        <button
          @click="goToPage(totalPages)"
          :class="[
            'px-3 py-2 text-sm font-medium rounded-lg transition-colors',
            currentPage === totalPages
              ? 'bg-blue-600 text-white'
              : 'text-gray-500 bg-white border border-gray-300 hover:bg-gray-50 hover:text-gray-700'
          ]"
        >
          {{ totalPages }}
        </button>
      </template>
    </div>

    <!-- Next Button -->
    <button
      @click="goToNext"
      :disabled="!hasNextPage"
      class="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white disabled:hover:text-gray-500"
    >
      Next
      <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>
  </div>

  <!-- Page Info -->
  <div class="text-sm text-gray-700 mt-3 text-center">
    Showing {{ startItem }} to {{ endItem }} of {{ totalItems }} results
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  currentPage: number
  totalPages: number
  totalItems: number
  pageSize: number
  hasNextPage: boolean
  hasPreviousPage: boolean
}

interface Emits {
  (e: 'page-change', page: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Computed properties for pagination logic
const visiblePages = computed(() => {
  const delta = 2 // Number of pages to show on each side of current page
  const range = []
  const rangeWithDots = []

  for (
    let i = Math.max(2, props.currentPage - delta);
    i <= Math.min(props.totalPages - 1, props.currentPage + delta);
    i++
  ) {
    range.push(i)
  }

  if (props.currentPage - delta > 2) {
    rangeWithDots.push(props.currentPage - delta - 1)
  }

  rangeWithDots.push(...range)

  if (props.currentPage + delta < props.totalPages - 1) {
    rangeWithDots.push(props.currentPage + delta + 1)
  }

  return rangeWithDots
})

const showFirstPage = computed(() => props.totalPages > 1 && !visiblePages.value.includes(1))
const showLastPage = computed(() => props.totalPages > 1 && !visiblePages.value.includes(props.totalPages))
const showStartEllipsis = computed(() => visiblePages.value.length > 0 && visiblePages.value[0] > 2)
const showEndEllipsis = computed(() => {
  const lastVisible = visiblePages.value[visiblePages.value.length - 1]
  return visiblePages.value.length > 0 && lastVisible < props.totalPages - 1
})

const startItem = computed(() => (props.currentPage - 1) * props.pageSize + 1)
const endItem = computed(() => Math.min(props.currentPage * props.pageSize, props.totalItems))

// Methods
const goToPage = (page: number) => {
  if (page !== props.currentPage && page >= 1 && page <= props.totalPages) {
    emit('page-change', page)
  }
}

const goToPrevious = () => {
  if (props.hasPreviousPage) {
    goToPage(props.currentPage - 1)
  }
}

const goToNext = () => {
  if (props.hasNextPage) {
    goToPage(props.currentPage + 1)
  }
}
</script>