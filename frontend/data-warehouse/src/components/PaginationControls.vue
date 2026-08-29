<script setup lang="ts">
defineProps<{
  page: number
  totalPages: number
  hasPreviousPage: boolean
  hasNextPage: boolean
}>()

defineEmits<{
  previous: []
  next: []
}>()
</script>

<template>
  <nav class="pagination" aria-label="Pagination">
    <button class="button" type="button" :disabled="!hasPreviousPage" @click="$emit('previous')">
      Previous
    </button>

    <span>Page {{ page }} of {{ totalPages }}</span>

    <button class="button" type="button" :disabled="!hasNextPage" @click="$emit('next')">
      Next
    </button>
  </nav>
</template>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-top: 1rem;
}
</style>
