<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import AppAlert from '@/components/AppAlert.vue'
import ProteinForm from '@/components/ProteinForm.vue'
import { proteinApi } from '@/services/proteinApi'
import type { CreateProteinRequest } from '@/types/protein'

const router = useRouter()

const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

async function createProtein(payload: CreateProteinRequest): Promise<void> {
  isSubmitting.value = true
  errorMessage.value = null

  try {
    await proteinApi.create(payload)

    await router.push({
      name: 'protein-detail',
      params: {
        accession: payload.accession,
      },
    })
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Protein could not be created.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <section class="page">
    <h1>Create protein</h1>
    <p>Create a new protein record in the data warehouse.</p>

    <AppAlert v-if="errorMessage" variant="error" :message="errorMessage" />

    <ProteinForm
      submit-label="Create protein"
      :is-submitting="isSubmitting"
      @create="createProtein"
      @cancel="router.push({ name: 'protein-list' })"
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
