import { afterEach, describe, expect, it, vi } from 'vitest'

import { ApiClientError, apiRequest } from './apiClient'

describe('apiRequest', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('returns parsed JSON for a successful response', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({
          status: 'ok',
        }),
      }),
    )

    const result = await apiRequest<{ status: string }>('/health')

    expect(result).toEqual({
      status: 'ok',
    })
  })

  it('throws ApiClientError for a structured API error', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 404,
        json: async () => ({
          code: 'not_found',
          message: 'Resource was not found.',
        }),
      }),
    )

    await expect(apiRequest('/proteins/UNKNOWN')).rejects.toMatchObject({
      name: 'ApiClientError',
      status: 404,
      code: 'not_found',
      message: 'Resource was not found.',
    })
  })

  it('throws a fallback error if an error response is not JSON', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: async () => {
          throw new Error('invalid JSON')
        },
      }),
    )

    await expect(apiRequest('/proteins')).rejects.toBeInstanceOf(ApiClientError)

    await expect(apiRequest('/proteins')).rejects.toMatchObject({
      status: 500,
      code: 'unknown_error',
    })
  })
})
