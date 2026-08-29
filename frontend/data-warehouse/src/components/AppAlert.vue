<script setup lang="ts">
defineProps<{
  variant: 'error' | 'success' | 'info'
  message: string
}>()
</script>

<template>
  <p class="alert" :class="`alert-${variant}`" role="alert">
    {{ message }}
  </p>
</template>

<style scoped>
.alert {
  margin: 1rem 0;
  padding: 0.875rem 1rem;
  border-radius: 0.375rem;
}

.alert-error {
  color: #991b1b;
  background: #fee2e2;
}

.alert-success {
  color: #166534;
  background: #dcfce7;
}

.alert-info {
  color: #1e40af;
  background: #dbeafe;
}
</style>
