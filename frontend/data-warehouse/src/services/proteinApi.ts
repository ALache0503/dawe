import { apiRequest } from './apiClient'

import type {
  CreateProteinRequest,
  ProteinDetails,
  ProteinPage,
  ProteinSearchParams,
  UpdateProteinRequest,
} from '@/types/protein'

function buildQueryString(params: ProteinSearchParams): string {
  const searchParams = new URLSearchParams()

  if (params.search?.trim()) {
    searchParams.set('search', params.search.trim())
  }

  if (params.page !== undefined) {
    searchParams.set('page', String(params.page))
  }

  if (params.pageSize !== undefined) {
    searchParams.set('pageSize', String(params.pageSize))
  }

  if (params.reviewed !== undefined) {
    searchParams.set('reviewed', String(params.reviewed))
  }

  if (params.taxonId !== undefined) {
    searchParams.set('taxonId', String(params.taxonId))
  }

  const query = searchParams.toString()
  return query ? `?${query}` : ''
}

export const proteinApi = {
  list(params: ProteinSearchParams): Promise<ProteinPage> {
    return apiRequest<ProteinPage>(`/proteins${buildQueryString(params)}`)
  },

  getDetails(accession: string): Promise<ProteinDetails> {
    return apiRequest<ProteinDetails>(
      `/proteins/${encodeURIComponent(accession)}`,
    )
  },

  create(payload: CreateProteinRequest): Promise<void> {
    return apiRequest<void>('/proteins', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
  },

  update(
    accession: string,
    payload: UpdateProteinRequest,
  ): Promise<void> {
    return apiRequest<void>(`/proteins/${encodeURIComponent(accession)}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
  },

  remove(accession: string): Promise<void> {
    return apiRequest<void>(`/proteins/${encodeURIComponent(accession)}`, {
      method: 'DELETE',
    })
  },
}