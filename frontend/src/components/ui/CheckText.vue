<script setup lang="ts">
import { CheckCircle } from 'lucide-vue-next'

defineProps<{
  text?: string
}>()
</script>
<template>
  <div class="flex items-center gap-1">
    <CheckCircle class="w-4 h-4 text-success mr-1" />
    <span class="text-sm">{{ text }}</span>
  </div>
</template>
