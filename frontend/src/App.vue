<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import { Trophy, BarChart3, TrendingUp } from 'lucide-vue-next'
import AppHeader from './components/AppHeader.vue'
import VTabs, { type TabOption } from './components/ui/VTabs.vue'
import { ref, watch } from 'vue'

const TAB_OPTIONS: TabOption[] = [
  {
    label: 'Our Selection',
    value: 'our-selection',
    href: '/',
    icon: Trophy,
  },
  {
    label: 'Stocks',
    value: 'stocks',
    href: '/stocks',
    icon: BarChart3,
  },
  {
    label: 'Recommendations',
    value: 'recommendations',
    href: '/recommendations',
    icon: TrendingUp,
  },
]

const route = useRoute()
const currentTab = ref('')

watch(route, (newRoute) => {
  currentTab.value =
    TAB_OPTIONS.find((option) => option.href === newRoute.fullPath)?.value || TAB_OPTIONS[0].value
})
</script>

<template>
  <div class="w-full min-h-screen bg-background text-foreground">
    <AppHeader :current-tab="currentTab" />

    <div class="container mx-auto px-4 py-6">
      <VTabs v-model:current-tab="currentTab" :options="TAB_OPTIONS">
        <template v-for="option in TAB_OPTIONS" :key="option.value" #[option.value]>
          <component :is="option.icon" class="w-4 h-4 mr-2" />
          {{ option.label }}
        </template>
      </VTabs>
    </div>

    <main class="container mx-auto px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
