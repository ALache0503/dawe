<template>
  <section class="page">
    <h1>UniProt Data Warehouse</h1>
    <p>
      A fullstack application for searching, managing, and inspecting representative UniProtKB
      protein records.
    </p>

    <RouterLink class="button button-primary" to="/proteins"> Browse proteins </RouterLink>
  </section>
</template>

<style scoped>
.page {
  max-width: 75rem;
  margin: 0 auto;
  padding: 2rem;
}
</style>
