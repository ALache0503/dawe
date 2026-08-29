<script setup lang="ts">
import type { ProteinListItem } from '@/types/protein'

defineProps<{
  proteins: ProteinListItem[]
}>()
</script>

<template>
  <div class="table-wrapper">
    <table>
      <thead>
        <tr>
          <th>Accession</th>
          <th>Protein name</th>
          <th>Organism</th>
          <th>Reviewed</th>
          <th>Score</th>
          <th>Length</th>
          <th><span class="sr-only">Actions</span></th>
        </tr>
      </thead>

      <tbody>
        <tr v-if="proteins.length === 0">
          <td colspan="7" class="empty-state">
            No proteins match the current search criteria.
          </td>
        </tr>

        <tr v-for="protein in proteins" :key="protein.accession">
          <td>
            <RouterLink
              :to="{
                name: 'protein-detail',
                params: { accession: protein.accession },
              }"
            >
              {{ protein.accession }}
            </RouterLink>
          </td>
          <td>{{ protein.proteinName }}</td>
          <td>{{ protein.organismName }}</td>
          <td>{{ protein.reviewed ? 'Yes' : 'No' }}</td>
          <td>{{ protein.annotationScore }}/5</td>
          <td>{{ protein.length }} aa</td>
          <td>
            <RouterLink
              class="table-action"
              :to="{
                name: 'protein-detail',
                params: { accession: protein.accession },
              }"
            >
              Details
            </RouterLink>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-wrapper {
  overflow-x: auto;
  margin-top: 1rem;
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
  padding: 0.875rem;
  border-bottom: 1px solid #e2e8f0;
  text-align: left;
  vertical-align: top;
}

th {
  color: #334155;
  background: #f8fafc;
}

tbody tr:last-child td {
  border-bottom: 0;
}

.empty-state {
  padding: 2rem;
  color: #64748b;
  text-align: center;
}

.table-action {
  white-space: nowrap;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
}
</style>