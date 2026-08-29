<script setup lang="ts">
import { computed, reactive, watch } from 'vue'

import type {
  CreateProteinRequest,
  Protein,
  UpdateProteinRequest,
} from '@/types/protein'


const props = withDefaults(
  defineProps<{
    protein?: Protein
    isSubmitting?: boolean
    submitLabel: string
    includeAccession?: boolean
  }>(),
  {
    protein: undefined,
    isSubmitting: false,
    includeAccession: true,
  },
)

const emit = defineEmits<{
  create: [payload: CreateProteinRequest]
  update: [payload: UpdateProteinRequest]
  cancel: []
}>()

type ProteinFormState = {
  accession: string
  taxonId: string
  entryName: string
  proteinName: string
  reviewed: boolean
  annotationScore: string
  mass: number | ''
  length: string
  sequence: string
  proteinExistence: string
  geneNames: string
}

function createInitialState(protein?: Protein): ProteinFormState {
  return {
    accession: protein?.accession ?? '',
    taxonId: protein?.taxonId.toString() ?? '',
    entryName: protein?.entryName ?? '',
    proteinName: protein?.proteinName ?? '',
    reviewed: protein?.reviewed ?? false,
    annotationScore: protein?.annotationScore.toString() ?? '1',
    mass: protein?.mass ?? '',
    length: protein?.length.toString() ?? '',
    sequence: protein?.sequence ?? '',
    proteinExistence: protein?.proteinExistence ?? '',
    geneNames: protein?.geneNames ?? '',
  }
}

const form = reactive<ProteinFormState>(createInitialState(props.protein))
const localError = computed(() => {
  if (!form.taxonId || Number(form.taxonId) <= 0) {
    return 'Taxon ID must be a positive number.'
  }

  if (!form.entryName.trim()) {
    return 'Entry name is required.'
  }

  if (!form.proteinName.trim()) {
    return 'Protein name is required.'
  }

  const score = Number(form.annotationScore)
  if (!Number.isInteger(score) || score < 1 || score > 5) {
    return 'Annotation score must be an integer between 1 and 5.'
  }

  if (!form.length || Number(form.length) <= 0) {
    return 'Length must be a positive number.'
  }

  if (!form.sequence.trim()) {
    return 'Sequence is required.'
  }

  if (form.mass !== '' && Number(form.mass) <= 0) {
    return 'Mass must be positive when provided.'
  }

  return null
})

watch(
  () => props.protein,
  (protein) => {
    Object.assign(form, createInitialState(protein))
  },
)

function submit(): void {
  if (localError.value) {
    return
  }

  const commonPayload: UpdateProteinRequest = {
    taxonId: Number(form.taxonId),
    entryName: form.entryName.trim(),
    proteinName: form.proteinName.trim(),
    reviewed: form.reviewed,
    annotationScore: Number(form.annotationScore),
    mass: form.mass === '' ? null : Number(form.mass),
    length: Number(form.length),
    sequence: form.sequence.trim(),
    proteinExistence: form.proteinExistence.trim() || null,
    geneNames: form.geneNames.trim() || null,
  }

  if (props.includeAccession) {
    const createPayload: CreateProteinRequest = {
      accession: form.accession.trim(),
      ...commonPayload,
    }

    emit('create', createPayload)
    return
  }

  emit('update', commonPayload)
}
</script>

<template>
  <form class="protein-form" @submit.prevent="submit">
    <label v-if="includeAccession">
      Accession
      <input
        v-model.trim="form.accession"
        maxlength="10"
        required
        placeholder="e.g. P69905"
      />
    </label>

    <label>
      Taxon ID
      <input v-model="form.taxonId" type="number" min="1" required />
    </label>

    <label>
      Entry name
      <input v-model.trim="form.entryName" maxlength="50" required />
    </label>

    <label class="full-width">
      Protein name
      <input v-model.trim="form.proteinName" required />
    </label>

    <label>
      Reviewed
      <input v-model="form.reviewed" type="checkbox" />
    </label>

    <label>
      Annotation score
      <select v-model="form.annotationScore">
        <option value="1">1</option>
        <option value="2">2</option>
        <option value="3">3</option>
        <option value="4">4</option>
        <option value="5">5</option>
      </select>
    </label>

    <label>
      Mass in Da
      <input
        v-model="form.mass"
        type="number"
        min="1"
        placeholder="Optional"
      />
    </label>

    <label>
      Length
      <input v-model="form.length" type="number" min="1" required />
    </label>

    <label class="full-width">
      Sequence
      <textarea v-model.trim="form.sequence" rows="8" required />
    </label>

    <label class="full-width">
      Protein existence
      <input
        v-model.trim="form.proteinExistence"
        placeholder="Optional"
      />
    </label>

    <label class="full-width">
      Gene names
      <input v-model.trim="form.geneNames" placeholder="Optional" />
    </label>

    <p v-if="localError" class="validation-error">
      {{ localError }}
    </p>

    <div class="actions full-width">
      <button
        class="button button-primary"
        type="submit"
        :disabled="Boolean(localError) || isSubmitting"
      >
        {{ isSubmitting ? 'Saving…' : submitLabel }}
      </button>

      <button class="button" type="button" :disabled="isSubmitting" @click="$emit('cancel')">
        Cancel
      </button>
    </div>
  </form>
</template>

<style scoped>
.protein-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
  max-width: 60rem;
  padding: 1.5rem;
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
select,
textarea {
  width: 100%;
  padding: 0.55rem;
  border: 1px solid #94a3b8;
  border-radius: 0.375rem;
  font: inherit;
}

textarea {
  resize: vertical;
  font-family: monospace;
}

.full-width {
  grid-column: 1 / -1;
}

.actions {
  display: flex;
  gap: 0.75rem;
}

.validation-error {
  grid-column: 1 / -1;
  margin: 0;
  color: #b91c1c;
}

@media (max-width: 42rem) {
  .protein-form {
    grid-template-columns: 1fr;
  }
}
</style>