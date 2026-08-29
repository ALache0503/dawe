<template>
  <header class="app-header">
    <RouterLink class="brand" to="/">UniProt Data Warehouse</RouterLink>

    <nav class="navigation" aria-label="Main navigation">
      <RouterLink to="/">Home</RouterLink>
      <RouterLink to="/proteins">Proteins</RouterLink>
      <RouterLink to="/proteins/new">Create protein</RouterLink>
    </nav>
  </header>
</template>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  min-height: 4rem;
  padding: 0 2rem;
  color: white;
  background: #172554;
}

.brand {
  color: white;
  font-weight: 700;
  text-decoration: none;
}

.navigation {
  display: flex;
  gap: 1rem;
}

.navigation a {
  color: #dbeafe;
  text-decoration: none;
}

.navigation a.router-link-active {
  color: white;
  font-weight: 700;
}
</style>
