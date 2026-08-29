import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'

import ProteinForm from './ProteinForm.vue'

import type { Protein } from '@/types/protein'

const existingProtein: Protein = {
  accession: 'TST0000001',
  taxonId: 9606,
  entryName: 'TEST_HUMAN',
  proteinName: 'Test protein',
  reviewed: false,
  annotationScore: 2,
  mass: 12345,
  length: 120,
  sequence: 'MTESTPROTEINSEQUENCE',
  proteinExistence: 'Predicted',
  geneNames: 'TESTGENE1',
}

describe('ProteinForm', () => {
  async function fillRequiredCreateFields(wrapper: ReturnType<typeof mount>): Promise<void> {
    const inputs = wrapper.findAll('input')

    await inputs[0].setValue('TST0000001')
    await inputs[1].setValue('9606')
    await inputs[2].setValue('TEST_HUMAN')
    await inputs[3].setValue('Test protein')
    await wrapper.find('select').setValue('2')
    await inputs[5].setValue('12345')
    await inputs[6].setValue('120')
    await wrapper.find('textarea').setValue('MTESTPROTEINSEQUENCE')
  }

  it('emits a create event with accession and numeric mass', async () => {
    const wrapper = mount(ProteinForm, {
      props: {
        submitLabel: 'Create protein',
      },
    })

    await fillRequiredCreateFields(wrapper)
    await wrapper.find('form').trigger('submit.prevent')

    const events = wrapper.emitted('create')

    expect(events).toHaveLength(1)
    expect(events?.[0]?.[0]).toEqual({
      accession: 'TST0000001',
      taxonId: 9606,
      entryName: 'TEST_HUMAN',
      proteinName: 'Test protein',
      reviewed: false,
      annotationScore: 2,
      mass: 12345,
      length: 120,
      sequence: 'MTESTPROTEINSEQUENCE',
      proteinExistence: null,
      geneNames: null,
    })
  })

  it('emits a create event with null for an empty optional mass', async () => {
    const wrapper = mount(ProteinForm, {
      props: {
        submitLabel: 'Create protein',
      },
    })

    const inputs = wrapper.findAll('input')

    await inputs[0].setValue('TST0000002')
    await inputs[1].setValue('9606')
    await inputs[2].setValue('TEST_NULL')
    await inputs[3].setValue('Protein without mass')
    await wrapper.find('select').setValue('3')
    await inputs[6].setValue('85')
    await wrapper.find('textarea').setValue('MNULLABLEMASSTESTSEQUENCE')

    await wrapper.find('form').trigger('submit.prevent')

    const events = wrapper.emitted('create')

    expect(events).toHaveLength(1)

    const payload = events?.[0]?.[0] as {
      mass: number | null
      proteinExistence: string | null
      geneNames: string | null
    }

    expect(payload.mass).toBeNull()
    expect(payload.proteinExistence).toBeNull()
    expect(payload.geneNames).toBeNull()
  })

  it('prefills the edit form from an existing protein', () => {
    const wrapper = mount(ProteinForm, {
      props: {
        protein: existingProtein,
        includeAccession: false,
        submitLabel: 'Save changes',
      },
    })

    const inputs = wrapper.findAll('input')

    expect(inputs[0].element.value).toBe('9606')
    expect(inputs[1].element.value).toBe('TEST_HUMAN')
    expect(inputs[2].element.value).toBe('Test protein')
    expect(inputs[4].element.value).toBe('12345')
    expect(inputs[5].element.value).toBe('120')
    expect(wrapper.find('textarea').element.value).toBe('MTESTPROTEINSEQUENCE')
  })

  it('emits update without accession in edit mode', async () => {
    const wrapper = mount(ProteinForm, {
      props: {
        protein: existingProtein,
        includeAccession: false,
        submitLabel: 'Save changes',
      },
    })

    await wrapper.find('form').trigger('submit.prevent')

    const events = wrapper.emitted('update')

    expect(events).toHaveLength(1)

    const payload = events?.[0]?.[0] as Record<string, unknown>

    expect(payload).not.toHaveProperty('accession')
    expect(payload).toMatchObject({
      taxonId: 9606,
      entryName: 'TEST_HUMAN',
      proteinName: 'Test protein',
      mass: 12345,
    })
  })

  it('does not emit create if accession is missing', async () => {
    const wrapper = mount(ProteinForm, {
      props: {
        submitLabel: 'Create protein',
      },
    })

    const inputs = wrapper.findAll('input')

    await inputs[1].setValue('9606')
    await inputs[2].setValue('TEST_HUMAN')
    await inputs[3].setValue('Test protein')
    await inputs[6].setValue('120')

    await wrapper.find('select').setValue('2')
    await wrapper.find('textarea').setValue('MTESTPROTEINSEQUENCE')

    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.emitted('create')).toBeUndefined()
    expect(wrapper.text()).toContain('Accession is required.')
  })
})
