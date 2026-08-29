<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import AppAlert from '@/components/AppAlert.vue'
import LoadingIndicator from '@/components/LoadingIndicator.vue'
import ProteinOverview from '@/components/ProteinOverview.vue'
import ProteinSequence from '@/components/ProteinSequence.vue'
import { ApiClientError } from '@/services/apiClient'
import { proteinApi } from '@/services/proteinApi'
import type { ProteinDetails } from '@/types/protein'

const props = defineProps<{
  accession: string
}>()

const router = useRouter()

const details = ref<ProteinDetails | null>(null)
const isLoading = ref(false)
const isDeleting = ref(false)
const errorMessage = ref<string | null>(null)

async function loadDetails(): Promise<void> {
  isLoading.value = true
  errorMessage.value = null

  try {
    details.value = await proteinApi.getDetails(props.accession)
  } catch (error) {
    if (error instanceof ApiClientError && error.status === 404) {
      errorMessage.value = `Protein ${props.accession} was not found.`
    } else {
      errorMessage.value =
        error instanceof Error ? error.message : 'Protein details could not be loaded.'
    }
  } finally {
    isLoading.value = false
  }
}

async function deleteProtein(): Promise<void> {
  const confirmed = window.confirm(
    `Delete protein ${props.accession}? This action cannot be undone.`,
  )

  if (!confirmed) {
    return
  }

  isDeleting.value = true
  errorMessage.value = null

  try {
    await proteinApi.remove(props.accession)
    await router.push({ name: 'protein-list' })
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Protein could not be deleted.'
  } finally {
    isDeleting.value = false
  }
}

onMounted(() => {
  void loadDetails()
})

watch(
  () => props.accession,
  () => {
    void loadDetails()
  },
)
</script>

<template>
  <section class="page">
    <div class="toolbar">
      <RouterLink class="button" to="/proteins">← Back to proteins</RouterLink>

      <div class="toolbar-actions">
        <RouterLink
          class="button button-primary"
          :to="{
            name: 'protein-edit',
            params: { accession },
          }"
        >
          Edit
        </RouterLink>

        <button
          class="button button-danger"
          type="button"
          :disabled="isDeleting"
          @click="deleteProtein"
        >
          {{ isDeleting ? 'Deleting…' : 'Delete' }}
        </button>
      </div>
    </div>

    <LoadingIndicator v-if="isLoading" />

    <AppAlert v-else-if="errorMessage" variant="error" :message="errorMessage" />

    <template v-else-if="details">
      <header class="detail-heading">
        <p class="accession">{{ details.protein.accession }}</p>
        <h1>{{ details.protein.proteinName }}</h1>
        <p>{{ details.protein.entryName }}</p>
      </header>

      <ProteinOverview :details="details" />
      <ProteinSequence :sequence="details.protein.sequence" />

      <section class="detail-section">
        <h2>Comments</h2>

        <p v-if="details.comments.length === 0">No comments available.</p>

        <article v-for="comment in details.comments" :key="comment.commentId" class="card">
          <h3>{{ comment.typeName }}</h3>
          <p>{{ comment.commentText }}</p>
        </article>
      </section>

      <section class="detail-section">
        <h2>Features</h2>

        <p v-if="details.features.length === 0">No features available.</p>

        <div v-else class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Type</th>
                <th>Description</th>
                <th>Position</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="feature in details.features" :key="feature.featureId">
                <td>{{ feature.typeName }}</td>
                <td>{{ feature.description || '—' }}</td>
                <td>{{ feature.startPosition }}–{{ feature.endPosition }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="detail-section">
        <h2>GO Terms</h2>

        <p v-if="details.goTerms.length === 0">No GO terms available.</p>

        <div v-else class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>GO ID</th>
                <th>Term</th>
                <th>Category</th>
                <th>Aspect</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="term in details.goTerms" :key="term.goId">
                <td>{{ term.goId }}</td>
                <td>{{ term.termName }}</td>
                <td>{{ term.category || '—' }}</td>
                <td>{{ term.aspect }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="detail-section">
        <h2>Keywords</h2>

        <p v-if="details.keywords.length === 0">No keywords available.</p>

        <ul v-else class="tag-list">
          <li v-for="keyword in details.keywords" :key="keyword.keywordId">
            {{ keyword.keywordName }}
          </li>
        </ul>
      </section>

      <section class="detail-section">
        <h2>Cross references</h2>

        <p v-if="details.crossReferences.length === 0">No cross references available.</p>

        <div v-else class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Database</th>
                <th>Reference ID</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="reference in details.crossReferences"
                :key="`${reference.databaseName}-${reference.referenceId}`"
              >
                <td>{{ reference.databaseName }}</td>
                <td>{{ reference.referenceId }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </template>
  </section>
</template>

<style scoped>
.page {
  max-width: 75rem;
  margin: 0 auto;
  padding: 2rem;
}

.toolbar,
.toolbar-actions {
  display: flex;
  gap: 0.75rem;
}

.toolbar {
  justify-content: space-between;
}

.detail-heading {
  margin: 2rem 0;
}

.detail-heading h1 {
  margin: 0.25rem 0;
}

.accession {
  margin: 0;
  color: #1d4ed8;
  font-family: monospace;
  font-size: 1.1rem;
  font-weight: 700;
}

.detail-section {
  margin-top: 2rem;
}

.card {
  padding: 1rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.5rem;
  background: white;
}

.card + .card {
  margin-top: 0.75rem;
}

.card h3 {
  margin-top: 0;
}

.table-wrapper {
  overflow-x: auto;
  border: 1px solid #cbd5e1;
  border-radius: 0.5rem;
  background: white;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 0.75rem;
  border-bottom: 1px solid #e2e8f0;
  text-align: left;
}

th {
  background: #f8fafc;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  padding: 0;
  list-style: none;
}

.tag-list li {
  padding: 0.4rem 0.6rem;
  border-radius: 999px;
  color: #1e3a8a;
  background: #dbeafe;
}

@media (max-width: 40rem) {
  .toolbar {
    flex-direction: column;
  }
}
</style>
