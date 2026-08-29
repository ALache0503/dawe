import type { PaginatedResponse } from './api'

export type ProteinListItem = {
  accession: string
  entryName: string
  proteinName: string
  organismName: string
  reviewed: boolean
  annotationScore: number
  length: number
}

export type Protein = {
  accession: string
  taxonId: number
  entryName: string
  proteinName: string
  reviewed: boolean
  annotationScore: number
  mass: number | null
  length: number
  sequence: string
  proteinExistence: string | null
  geneNames: string | null
}

export type Organism = {
  taxonId: number
  commonName: string
  scientificName: string
}

export type ProteinComment = {
  commentId: number
  typeName: string
  commentText: string
}

export type ProteinFeature = {
  featureId: number
  typeName: string
  description: string
  startPosition: number
  endPosition: number
}

export type GoTerm = {
  goId: string
  termName: string
  category: string
  aspect: 'P' | 'F' | 'C'
}

export type Keyword = {
  keywordId: string
  keywordName: string
}

export type CrossReference = {
  referenceId: string
  databaseName: string
}

export type ProteinDetails = {
  protein: Protein
  organism: Organism
  comments: ProteinComment[]
  features: ProteinFeature[]
  goTerms: GoTerm[]
  keywords: Keyword[]
  crossReferences: CrossReference[]
}

export type ProteinSearchParams = {
  search?: string
  page?: number
  pageSize?: number
  reviewed?: boolean
  taxonId?: number
}

export type ProteinPage = PaginatedResponse<ProteinListItem>

export type CreateProteinRequest = {
  accession: string
  taxonId: number
  entryName: string
  proteinName: string
  reviewed: boolean
  annotationScore: number
  mass: number | null
  length: number
  sequence: string
  proteinExistence: string | null
  geneNames: string | null
}

export type UpdateProteinRequest = Omit<CreateProteinRequest, 'accession'>
