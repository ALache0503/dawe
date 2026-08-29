import type { Protein, ProteinDetails, ProteinListItem } from '@/types/protein'

export const proteinListItem: ProteinListItem = {
  accession: 'TST0000001',
  entryName: 'TEST_HUMAN',
  proteinName: 'Test protein for frontend tests',
  organismName: 'Human',
  reviewed: true,
  annotationScore: 4,
  length: 120,
}

export const protein: Protein = {
  accession: 'TST0000001',
  taxonId: 9606,
  entryName: 'TEST_HUMAN',
  proteinName: 'Test protein for frontend tests',
  reviewed: true,
  annotationScore: 4,
  mass: 12345,
  length: 120,
  sequence: 'MTESTPROTEINSEQUENCEFORFRONTENDTESTING',
  proteinExistence: 'Predicted',
  geneNames: 'TESTGENE1',
}

export const proteinWithoutMass: Protein = {
  ...protein,
  accession: 'TST0000002',
  entryName: 'TEST_NULL',
  proteinName: 'Test protein without known mass',
  mass: null,
  proteinExistence: null,
  geneNames: null,
}

export const proteinDetails: ProteinDetails = {
  protein,
  organism: {
    taxonId: 9606,
    commonName: 'Human',
    scientificName: 'Homo sapiens',
  },
  comments: [],
  features: [],
  goTerms: [],
  keywords: [],
  crossReferences: [],
}
