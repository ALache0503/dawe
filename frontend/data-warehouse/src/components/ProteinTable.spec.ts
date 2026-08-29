import { describe, expect, it } from 'vitest'
import { RouterLinkStub, mount } from '@vue/test-utils'

import ProteinTable from './ProteinTable.vue'

import type { ProteinListItem } from '@/types/protein'

const proteins: ProteinListItem[] = [
  {
    accession: 'P69905',
    entryName: 'HBA_HUMAN',
    proteinName: 'Hemoglobin subunit alpha',
    organismName: 'Human',
    reviewed: true,
    annotationScore: 5,
    length: 142,
  },
  {
    accession: 'TST0000001',
    entryName: 'TEST_HUMAN',
    proteinName: 'Test protein',
    organismName: 'Human',
    reviewed: false,
    annotationScore: 2,
    length: 120,
  },
]

describe('ProteinTable', () => {
  function mountTable(items: ProteinListItem[]) {
    return mount(ProteinTable, {
      props: {
        proteins: items,
      },
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })
  }

  it('renders one table row per protein', () => {
    const wrapper = mountTable(proteins)

    expect(wrapper.findAll('tbody tr')).toHaveLength(2)
    expect(wrapper.text()).toContain('P69905')
    expect(wrapper.text()).toContain('Hemoglobin subunit alpha')
    expect(wrapper.text()).toContain('TST0000001')
    expect(wrapper.text()).toContain('Test protein')
  })

  it('renders reviewed status correctly', () => {
    const wrapper = mountTable(proteins)

    expect(wrapper.text()).toContain('Yes')
    expect(wrapper.text()).toContain('No')
  })

  it('renders an empty state for no proteins', () => {
    const wrapper = mountTable([])

    expect(wrapper.text()).toContain('No proteins match the current search criteria.')
  })

  it('links each protein to its detail view', () => {
    const wrapper = mountTable(proteins)

    const links = wrapper.findAllComponents(RouterLinkStub)

    expect(links[0].props('to')).toEqual({
      name: 'protein-detail',
      params: {
        accession: 'P69905',
      },
    })
  })
})
