export interface Article {
    id: number
    title: string
    url: string
    summary: string
    tags: string[]
    memo: string
    createdAt: string
    updatedAt: string
}

// 記事作成時のリクエスト型
export interface CreateArticleInput {
    title: string
    url: string
    summary: string
    tags: string[]
    memo?: string
}

// 記事更新時のリクエスト型
export interface UpdateArticleInput {
    title: string
    url: string
    summary: string
    tags: string[]
    memo?: string
}
