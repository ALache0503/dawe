import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { proteinApi } from '@/services/proteinApi'
import type {
  ProteinListItem,
  ProteinSearchParams,
} from '@/types/protein'

export const useProteinStore = defineStore('proteins', () => {
  const items = ref<ProteinListItem[]>([])
  const page = ref(1)
  const pageSize = ref(20)
  const totalItems = ref(0)
  const totalPages = ref(0)

  const search = ref('')
  const reviewed = ref<boolean | undefined>(undefined)
  const taxonId = ref<number | undefined>(undefined)

  const isLoading = ref(false)
  const errorMessage = ref<string | null>(null)

  const hasPreviousPage = computed(() => page.value > 1)
  const hasNextPage = computed(() => page.value < totalPages.value)

  function currentParams(): ProteinSearchParams {
    return {
      search: search.value,
      page: page.value,
      pageSize: pageSize.value,
      reviewed: reviewed.value,
      taxonId: taxonId.value,
    }
  }

  async function loadProteins(): Promise<void> {
    isLoading.value = true
    errorMessage.value = null

    try {
      const result = await proteinApi.list(currentParams())

      items.value = result.items
      page.value = result.page
      pageSize.value = result.pageSize
      totalItems.value = result.totalItems
      totalPages.value = result.totalPages
    } catch (error) {
      errorMessage.value =
        error instanceof Error
          ? error.message
          : 'Proteins could not be loaded.'
    } finally {
      isLoading.value = false
    }
  }

  async function applySearch(
    newSearch: string,
    newReviewed?: boolean,
    newTaxonId?: number,
  ): Promise<void> {
    search.value = newSearch
    reviewed.value = newReviewed
    taxonId.value = newTaxonId
    page.value = 1

    await loadProteins()
  }

  async function goToPage(targetPage: number): Promise<void> {
    if (targetPage < 1 || targetPage > totalPages.value) {
      return
    }

    page.value = targetPage
    await loadProteins()
  }

  return {
    items,
    page,
    pageSize,
    totalItems,
    totalPages,
    search,
    reviewed,
    taxonId,
    isLoading,
    errorMessage,
    hasPreviousPage,
    hasNextPage,
    loadProteins,
    applySearch,
    goToPage,
  }
})