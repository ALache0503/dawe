<script setup lang="ts">
import { onMounted } from 'vue'

import AppAlert from '@/components/AppAlert.vue'
import LoadingIndicator from '@/components/LoadingIndicator.vue'
import PaginationControls from '@/components/PaginationControls.vue'
import ProteinSearchForm from '@/components/ProteinSearchForm.vue'
import ProteinTable from '@/components/ProteinTable.vue'
import { useProteinStore } from '@/stores/proteinStore'

const proteinStore = useProteinStore()

onMounted(() => {
  void proteinStore.loadProteins()
})
</script>

<template>
  <section class="page">
    <div class="page-heading">
      <div>
        <h1>Proteins</h1>
        <p>Search, inspect, create, update, and delete protein records.</p>
      </div>

      <RouterLink class="button button-primary" to="/proteins/new"> Create protein </RouterLink>
    </div>

    <ProteinSearchForm
      :initial-search="proteinStore.search"
      :initial-reviewed="proteinStore.reviewed"
      :initial-taxon-id="proteinStore.taxonId"
      @search="(search, reviewed, taxonId) => proteinStore.applySearch(search, reviewed, taxonId)"
    />

    <AppAlert
      v-if="proteinStore.errorMessage"
      variant="error"
      :message="proteinStore.errorMessage"
    />

    <LoadingIndicator v-if="proteinStore.isLoading" />

    <template v-else>
      <p class="result-count">
        {{ proteinStore.totalItems }} result<span v-if="proteinStore.totalItems !== 1">s</span>
      </p>

      <ProteinTable :proteins="proteinStore.items" />

      <PaginationControls
        v-if="proteinStore.totalPages > 0"
        :page="proteinStore.page"
        :total-pages="proteinStore.totalPages"
        :has-previous-page="proteinStore.hasPreviousPage"
        :has-next-page="proteinStore.hasNextPage"
        @previous="proteinStore.goToPage(proteinStore.page - 1)"
        @next="proteinStore.goToPage(proteinStore.page + 1)"
      />
    </template>
  </section>
</template>

<style scoped>
.page {
  max-width: 75rem;
  margin: 0 auto;
  padding: 2rem;
}

.page-heading {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.page-heading h1 {
  margin: 0;
}

.page-heading p {
  margin: 0.5rem 0 0;
  color: #475569;
}

.result-count {
  color: #475569;
}
</style>
