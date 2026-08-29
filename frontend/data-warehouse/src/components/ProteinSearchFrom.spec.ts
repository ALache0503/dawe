import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'

import ProteinSearchForm from './ProteinSearchForm.vue'

describe('ProteinSearchForm', () => {
  function createWrapper() {
    return mount(ProteinSearchForm, {
      props: {
        initialSearch: '',
        initialReviewed: undefined,
        initialTaxonId: undefined,
      },
    })
  }

  it('emits search with text, reviewed filter, and taxon ID', async () => {
    const wrapper = createWrapper()

    const inputs = wrapper.findAll('input')
    const searchInput = inputs[0]
    const taxonIdInput = inputs[1]

    await searchInput.setValue('hemoglobin')
    await wrapper.find('select').setValue('true')
    await taxonIdInput.setValue('9606')

    await wrapper.find('form').trigger('submit.prevent')

    const events = wrapper.emitted('search')

    expect(events).toHaveLength(1)
    expect(events?.[0]).toEqual(['hemoglobin', true, 9606])
  })

  it('emits an undefined taxon ID when the field is empty', async () => {
    const wrapper = createWrapper()

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('hemoglobin')

    await wrapper.find('form').trigger('submit.prevent')

    const events = wrapper.emitted('search')

    expect(events).toHaveLength(1)
    expect(events?.[0]).toEqual(['hemoglobin', undefined, undefined])
  })

  it('emits empty filters when reset is clicked', async () => {
    const wrapper = mount(ProteinSearchForm, {
      props: {
        initialSearch: 'hemoglobin',
        initialReviewed: true,
        initialTaxonId: 9606,
      },
    })

    await wrapper.find('button[type="button"]').trigger('click')

    const events = wrapper.emitted('search')

    expect(events).toHaveLength(1)
    expect(events?.[0]).toEqual(['', undefined, undefined])
  })
})
