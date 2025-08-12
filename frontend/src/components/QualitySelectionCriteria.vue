<script setup lang="ts">
import { CheckCircle } from 'lucide-vue-next'
import CardContainer from './ui/CardContainer.vue'
import CheckText from './ui/CheckText.vue'
</script>

<template>
  <CardContainer label="Quality Selection Criteria" :icon="CheckCircle" className="mt-6">
    <h3 class="text-sm text-muted-foreground mb-2 mt-1">
      Only stocks meeting multiple positive criteria are selected
    </h3>
    <div class="grid grid-cols-2 gap-3 mt-4">
      <CheckText text="Recent analyst upgrades (upgraded rating)" />
      <CheckText text="Buy or Strong Buy current ratings" />
      <CheckText text="Price target increases" />
      <CheckText text="Recommendations from top-tier brokerages" />
      <CheckText text="Recent activity (last 30-60 days)" />
      <CheckText text="New initiations with positive ratings" />
    </div>
  </CardContainer>
</template>
