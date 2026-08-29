import { afterEach, describe, expect, it, vi } from 'vitest'

import { proteinApi } from './proteinApi'

describe('proteinApi', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('builds a list request with search, pagination, and filters', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        items: [],
        page: 2,
        pageSize: 10,
        totalItems: 0,
        totalPages: 0,
      }),
    })

    vi.stubGlobal('fetch', fetchMock)

    await proteinApi.list({
      search: 'hemoglobin',
      page: 2,
      pageSize: 10,
      reviewed: true,
      taxonId: 9606,
    })

    expect(fetchMock).toHaveBeenCalledTimes(1)

    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:8080/api/v1/proteins?search=hemoglobin&page=2&pageSize=10&reviewed=true&taxonId=9606',
      expect.objectContaining({
        headers: expect.objectContaining({
          Accept: 'application/json',
        }),
      }),
    )
  })

  it('encodes the accession in detail requests', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        protein: {},
        organism: {},
        comments: [],
        features: [],
        goTerms: [],
        keywords: [],
        crossReferences: [],
      }),
    })

    vi.stubGlobal('fetch', fetchMock)

    await proteinApi.getDetails('TST0000001')

    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:8080/api/v1/proteins/TST0000001',
      expect.any(Object),
    )
  })

  it('sends a POST request when creating a protein', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: async () => ({
        accession: 'TST0000001',
      }),
    })

    vi.stubGlobal('fetch', fetchMock)

    const payload = {
      accession: 'TST0000001',
      taxonId: 9606,
      entryName: 'TEST_HUMAN',
      proteinName: 'Test protein',
      reviewed: false,
      annotationScore: 2,
      mass: 12345,
      length: 120,
      sequence: 'MTESTSEQUENCE',
      proteinExistence: null,
      geneNames: null,
    }

    await proteinApi.create(payload)

    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:8080/api/v1/proteins',
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({
          'Content-Type': 'application/json',
        }),
        body: JSON.stringify(payload),
      }),
    )
  })

  it('sends a DELETE request when deleting a protein', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 204,
    })

    vi.stubGlobal('fetch', fetchMock)

    await proteinApi.remove('TST0000001')

    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:8080/api/v1/proteins/TST0000001',
      expect.objectContaining({
        method: 'DELETE',
      }),
    )
  })
})
