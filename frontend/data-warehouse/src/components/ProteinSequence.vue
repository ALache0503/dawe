<script setup lang="ts">
defineProps<{
  sequence: string
}>()

async function copySequence(sequence: string): Promise<void> {
  await navigator.clipboard.writeText(sequence)
}
</script>

<template>
  <section class="sequence-section">
    <div class="heading">
      <h2>Sequence</h2>

      <button class="button" type="button" @click="copySequence(sequence)">
        Copy sequence
      </button>
    </div>

    <pre>{{ sequence }}</pre>
  </section>
</template>

<style scoped>
.sequence-section {
  margin-top: 2rem;
}

.heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

pre {
  max-height: 18rem;
  padding: 1rem;
  overflow: auto;
  border-radius: 0.5rem;
  color: #dbeafe;
  background: #0f172a;
  font-family: monospace;
  line-height: 1.6;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}
</style>