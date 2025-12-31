export interface Article {
    id: number
    title: string
    url: string
    summary: string
    tags: string[]
    memo?: string
    createdAt: string
    updatedAt: string
}
