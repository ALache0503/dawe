export type ApiError = {
  code: string
  message: string
}

export type PaginatedResponse<T> = {
  items: T[]
  page: number
  pageSize: number
  totalItems: number
  totalPages: number
}