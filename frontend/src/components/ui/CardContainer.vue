<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  label: string
  icon: Component
  className?: string
}>()
</script>

<template>
  <div :class="['border border-border bg-card p-4 rounded-lg', className]">
    <div class="flex items-center">
      <icon class="w-5 h-5" />
      <h2 class="ml-2">{{ label }}</h2>
    </div>

    <slot></slot>
  </div>
</template>
