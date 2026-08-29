<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  initialSearch: string
  initialReviewed?: boolean
  initialTaxonId?: number
}>()

const emit = defineEmits<{
  search: [search: string, reviewed: boolean | undefined, taxonId: number | undefined]
}>()

const search = ref(props.initialSearch)
const reviewedValue = ref(props.initialReviewed === undefined ? '' : String(props.initialReviewed))
const taxonId = ref(props.initialTaxonId?.toString() ?? '')

const localError = ref<string | null>(null)

function submit(): void {
  const normalizedTaxonId = taxonId.value === '' ? undefined : Number(taxonId.value)

  if (
    normalizedTaxonId !== undefined &&
    (!Number.isInteger(normalizedTaxonId) || normalizedTaxonId <= 0)
  ) {
    localError.value = 'Taxon ID must be a positive integer.'
    return
  }

  const normalizedReviewed = reviewedValue.value === '' ? undefined : reviewedValue.value === 'true'

  emit('search', search.value, normalizedReviewed, normalizedTaxonId)
}

function reset(): void {
  search.value = ''
  reviewedValue.value = ''
  taxonId.value = ''

  emit('search', '', undefined, undefined)
}
</script>

<template>
  <form class="search-form" @submit.prevent="submit">
    <label>
      Search proteins
      <input v-model="search" type="search" placeholder="Accession, entry name, or protein name" />
    </label>

    <label>
      Reviewed status
      <select v-model="reviewedValue">
        <option value="">All records</option>
        <option value="true">Reviewed</option>
        <option value="false">Unreviewed</option>
      </select>
    </label>

    <label>
      Taxon ID
      <input v-model="taxonId" type="number" min="1" placeholder="e.g. 9606" />
    </label>

    <div class="actions">
      <button class="button button-primary" type="submit">Search</button>
      <button class="button" type="button" @click="reset">Reset</button>
    </div>
  </form>
</template>

<style scoped>
.search-form {
  display: grid;
  grid-template-columns: minmax(16rem, 2fr) minmax(10rem, 1fr) minmax(8rem, 1fr) auto;
  gap: 1rem;
  align-items: end;
  padding: 1rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.5rem;
  background: white;
}

label {
  display: grid;
  gap: 0.4rem;
  color: #334155;
  font-weight: 600;
}

input,
select {
  min-height: 2.5rem;
  padding: 0.5rem;
  border: 1px solid #94a3b8;
  border-radius: 0.375rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

@media (max-width: 52rem) {
  .search-form {
    grid-template-columns: 1fr;
  }
}
</style>
