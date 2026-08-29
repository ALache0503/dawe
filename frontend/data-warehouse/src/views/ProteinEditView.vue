<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import AppAlert from '@/components/AppAlert.vue'
import LoadingIndicator from '@/components/LoadingIndicator.vue'
import ProteinForm from '@/components/ProteinForm.vue'
import { proteinApi } from '@/services/proteinApi'
import type { Protein, UpdateProteinRequest } from '@/types/protein'

const props = defineProps<{
  accession: string
}>()

const router = useRouter()

const protein = ref<Protein | null>(null)
const isLoading = ref(false)
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

async function loadProtein(): Promise<void> {
  isLoading.value = true
  errorMessage.value = null

  try {
    const details = await proteinApi.getDetails(props.accession)
    protein.value = details.protein
  } catch (error) {
    errorMessage.value =
      error instanceof Error
        ? error.message
        : 'Protein could not be loaded.'
  } finally {
    isLoading.value = false
  }
}

async function updateProtein(payload: UpdateProteinRequest): Promise<void> {
  isSubmitting.value = true
  errorMessage.value = null

  try {
    await proteinApi.update(props.accession, payload)

    await router.push({
      name: 'protein-detail',
      params: {
        accession: props.accession,
      },
    })
  } catch (error) {
    errorMessage.value =
      error instanceof Error
        ? error.message
        : 'Protein could not be updated.'
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  void loadProtein()
})
</script>

<template>
  <section class="page">
    <h1>Edit protein</h1>

    <LoadingIndicator v-if="isLoading" />

    <AppAlert
      v-else-if="errorMessage"
      variant="error"
      :message="errorMessage"
    />

    <ProteinForm
    v-else-if="protein"
      :protein="protein"
      :include-accession="false"
      submit-label="Save changes"
      :is-submitting="isSubmitting"
      @update="updateProtein"
      @cancel="
        router.push({
          name: 'protein-detail',
          params: { accession },
        })
      "
    />
  </section>
</template>

<style scoped>
.page {
  max-width: 75rem;
  margin: 0 auto;
  padding: 2rem;
}
</style>