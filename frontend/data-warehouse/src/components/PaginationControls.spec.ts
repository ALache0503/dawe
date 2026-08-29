import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'

import PaginationControls from './PaginationControls.vue'

describe('PaginationControls', () => {
  function createWrapper(overrides = {}) {
    return mount(PaginationControls, {
      props: {
        page: 2,
        totalPages: 5,
        hasPreviousPage: true,
        hasNextPage: true,
        ...overrides,
      },
    })
  }

  it('shows current and total page count', () => {
    const wrapper = createWrapper()

    expect(wrapper.text()).toContain('Page 2 of 5')
  })

  it('emits previous after clicking Previous', async () => {
    const wrapper = createWrapper()

    await wrapper.findAll('button')[0].trigger('click')

    expect(wrapper.emitted('previous')).toHaveLength(1)
  })

  it('emits next after clicking Next', async () => {
    const wrapper = createWrapper()

    await wrapper.findAll('button')[1].trigger('click')

    expect(wrapper.emitted('next')).toHaveLength(1)
  })

  it('disables Previous on the first page', () => {
    const wrapper = createWrapper({
      page: 1,
      hasPreviousPage: false,
    })

    expect((wrapper.findAll('button')[0].element as HTMLButtonElement).disabled).toBe(true)
  })

  it('disables Next on the last page', () => {
    const wrapper = createWrapper({
      page: 5,
      totalPages: 5,
      hasNextPage: false,
    })

    expect((wrapper.findAll('button')[1].element as HTMLButtonElement).disabled).toBe(true)
  })
})
