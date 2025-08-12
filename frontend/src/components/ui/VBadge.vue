<script setup lang="ts">
interface Props {
  variant?: 'default' | 'secondary' | 'outline' | 'destructive' | 'success'
  className?: string
}

const { variant = 'default', className } = defineProps<Props>()

const variantClasses = {
  default: 'bg-primary text-primary-foreground',
  secondary: 'bg-secondary text-secondary-foreground',
  outline: 'text-foreground border border-border bg-transparent',
  destructive: 'bg-destructive text-destructive-foreground',
  success: 'bg-success text-white',
}
</script>

<template>
  <div
    :class="[
      'inline-flex items-center justify-center rounded-md px-2 py-0.5 text-xs font-medium whitespace-nowrap shrink-0 gap-1 transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 [&>svg]:size-3 [&>svg]:pointer-events-none',
      variantClasses[variant],
      className,
    ]"
  >
    <slot />
  </div>
</template>
