import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { proteinApi } from '@/services/proteinApi'
import { useProteinStore } from './proteinStore'

vi.mock('@/services/proteinApi', () => ({
  proteinApi: {
    list: vi.fn(),
  },
}))

describe('proteinStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads proteins and pagination metadata', async () => {
    vi.mocked(proteinApi.list).mockResolvedValue({
      items: [
        {
          accession: 'P69905',
          entryName: 'HBA_HUMAN',
          proteinName: 'Hemoglobin subunit alpha',
          organismName: 'Human',
          reviewed: true,
          annotationScore: 5,
          length: 142,
        },
      ],
      page: 1,
      pageSize: 20,
      totalItems: 1,
      totalPages: 1,
    })

    const store = useProteinStore()

    await store.loadProteins()

    expect(proteinApi.list).toHaveBeenCalledWith({
      search: '',
      page: 1,
      pageSize: 20,
      reviewed: undefined,
      taxonId: undefined,
    })

    expect(store.items).toHaveLength(1)
    expect(store.totalItems).toBe(1)
    expect(store.totalPages).toBe(1)
    expect(store.isLoading).toBe(false)
    expect(store.errorMessage).toBeNull()
  })

  it('resets the page to one when applying a search', async () => {
    vi.mocked(proteinApi.list).mockResolvedValue({
      items: [],
      page: 1,
      pageSize: 20,
      totalItems: 0,
      totalPages: 0,
    })

    const store = useProteinStore()
    store.page = 3

    await store.applySearch('hemoglobin', true, 9606)

    expect(store.page).toBe(1)
    expect(proteinApi.list).toHaveBeenCalledWith({
      search: 'hemoglobin',
      page: 1,
      pageSize: 20,
      reviewed: true,
      taxonId: 9606,
    })
  })

  it('stores an error when loading proteins fails', async () => {
    vi.mocked(proteinApi.list).mockRejectedValue(new Error('Network error'))

    const store = useProteinStore()

    await store.loadProteins()

    expect(store.items).toEqual([])
    expect(store.errorMessage).toBe('Network error')
    expect(store.isLoading).toBe(false)
  })
})
