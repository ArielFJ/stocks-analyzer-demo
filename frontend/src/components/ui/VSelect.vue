<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { ChevronDown, Check } from 'lucide-vue-next'

export interface SelectOption {
  value: string
  label: string
}

interface Props {
  modelValue?: string
  options: SelectOption[]
  placeholder?: string
  className?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const isOpen = ref(false)
const selectRef = ref<HTMLDivElement>()

const selectedOption = computed(() =>
  props.options.find((option) => option.value === props.modelValue),
)

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

const selectOption = (option: SelectOption) => {
  emit('update:modelValue', option.value)
  isOpen.value = false
}

const handleClickOutside = (event: Event) => {
  if (selectRef.value && !selectRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

// Add click outside listener
if (typeof window !== 'undefined') {
  window.addEventListener('click', handleClickOutside)
}

onUnmounted(() => {
  // Clean up the event listener
  if (typeof window !== 'undefined') {
    window.removeEventListener('click', handleClickOutside)
  }
})
</script>

<template>
  <div ref="selectRef" class="relative" :class="className">
    <button
      type="button"
      @click="toggleDropdown"
      class="flex items-center justify-between w-full px-3 py-2 text-sm rounded-md bg-transparent hover:bg-muted transition-colors"
    >
      <span class="text-muted-foreground">
        {{ selectedOption?.label || placeholder || 'Select option...' }}
      </span>
      <ChevronDown
        class="w-4 h-4 text-muted-foreground transition-transform"
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute z-50 w-full mt-1 bg-popover border border-border rounded-md shadow-lg"
      >
        <div class="py-1 px-1 max-h-60 overflow-auto">
          <button
            v-for="option in options"
            :key="option.value"
            type="button"
            @click="selectOption(option)"
            class="flex items-center rounded justify-between w-full px-3 py-2 text-sm text-left hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus:outline-none transition-colors"
          >
            <span>{{ option.label }}</span>
            <Check v-if="option.value === modelValue" class="w-4 h-4 text-primary" />
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>
