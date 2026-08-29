<script setup lang="ts">
import type { ProteinDetails } from '@/types/protein'

defineProps<{
  details: ProteinDetails
}>()
</script>

<template>
  <section class="overview">
    <h2>Overview</h2>

    <dl>
      <div>
        <dt>Accession</dt>
        <dd>{{ details.protein.accession }}</dd>
      </div>

      <div>
        <dt>Organism</dt>
        <dd>
          {{ details.organism.scientificName }}
          <span v-if="details.organism.commonName"> ({{ details.organism.commonName }}) </span>
        </dd>
      </div>

      <div>
        <dt>Taxon ID</dt>
        <dd>{{ details.organism.taxonId }}</dd>
      </div>

      <div>
        <dt>Reviewed</dt>
        <dd>{{ details.protein.reviewed ? 'Yes' : 'No' }}</dd>
      </div>

      <div>
        <dt>Annotation score</dt>
        <dd>{{ details.protein.annotationScore }}/5</dd>
      </div>

      <div>
        <dt>Mass</dt>
        <dd>
          {{
            details.protein.mass != null
              ? `${details.protein.mass.toLocaleString()} Da`
              : 'Not available'
          }}
        </dd>
      </div>

      <div>
        <dt>Length</dt>
        <dd>{{ details.protein.length }} amino acids</dd>
      </div>

      <div>
        <dt>Protein existence</dt>
        <dd>{{ details.protein.proteinExistence ?? 'Not available' }}</dd>
      </div>

      <div>
        <dt>Gene names</dt>
        <dd>{{ details.protein.geneNames ?? 'Not available' }}</dd>
      </div>
    </dl>
  </section>
</template>

<style scoped>
.overview {
  padding: 1.25rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.5rem;
  background: white;
}

.overview h2 {
  margin-top: 0;
}

dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem 2rem;
  margin: 0;
}

dt {
  color: #64748b;
  font-size: 0.875rem;
}

dd {
  margin: 0.25rem 0 0;
  color: #0f172a;
  overflow-wrap: anywhere;
}

@media (max-width: 42rem) {
  dl {
    grid-template-columns: 1fr;
  }
}
</style>
