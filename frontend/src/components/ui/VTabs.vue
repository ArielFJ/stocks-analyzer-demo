<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { Component } from 'vue'

export interface TabOption {
  label: string
  value: string
  href: string
  icon: Component
}

const { currentTab } = defineProps<{
  currentTab: string
  options: TabOption[]
}>()

const emit = defineEmits<{
  'update:currentTab': [value: string]
}>()

const isActive = (option: TabOption) => {
  return currentTab === option.value
}

const handleTabClick = (option: TabOption) => {
  if (!isActive(option)) {
    emit('update:currentTab', option.value)
  }
}
</script>

<template>
  <div class="rounded-xl bg-muted p-1.5">
    <ul class="flex justify-around">
      <li
        v-for="option in options"
        :key="option.value"
        :class="{
          'flex flex-1 transition-all duration-150 ease-in-out scale-95': true,
          'bg-white rounded-xl transform scale-100': isActive(option),
        }"
      >
        <RouterLink
          :to="option.href"
          :class="{
            'flex py-2 cursor-pointer items-center justify-center text-sm font-medium text-muted-foreground hover:text-foreground w-full': true,
            '!text-foreground': isActive(option),
          }"
          @click="handleTabClick(option)"
        >
          <!-- Dynamic Slot to support multiple tabs -->
          <slot :name="option.value">
            {{ option.label }}
          </slot>
        </RouterLink>
      </li>
    </ul>
  </div>
</template>
