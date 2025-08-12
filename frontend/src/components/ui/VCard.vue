<script setup lang="ts">
defineProps<{
  className?: string
  contentClassName?: string
  noElevation?: boolean
}>()
</script>

<template>
  <div
    :class="[
      'rounded-lg border bg-card text-card-foreground',
      className,
      {
        'shadow-md hover:shadow-lg transition-shadow': !noElevation,
        'shadow-none': noElevation,
      },
    ]"
  >
    <div :class="['p-6', contentClassName]">
      <slot />
    </div>
  </div>
</template>
