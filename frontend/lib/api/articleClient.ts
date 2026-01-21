import { Article, CreateArticleInput, UpdateArticleInput } from '@/types/article'

// APIから返却される記事データの型
interface ApiArticle {
    id: number
    title: string
    url: string
    summary: string
    tags: string[] | null
    memo: string | null
    created_at: string
    updated_at: string
}

// APIのベースURL
const API_BASE_URL = 'http://localhost:8080'

// バックエンドのAPIと通信するクライアント
class ArticleClient {
    // 全記事を取得
    async getAll(): Promise<Article[]> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/articles`)

            if (!response.ok) {
                throw new Error('Failed to fetch articles')
            }

            const data: ApiArticle[] = await response.json()
            return data.map(this.convertToCamelCase)
        } catch (error) {
            throw error
        }
    }

    // 指定IDの記事を取得
    async getById(id: number): Promise<Article> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/articles/${id}`)

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to fetch article')
            }

            const data: ApiArticle = await response.json()
            return this.convertToCamelCase(data)
        } catch (error) {
            throw error
        }
    }

    // 記事を作成
    async create(article: CreateArticleInput): Promise<Article> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/articles`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(article),
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to create article')
            }

            const data: ApiArticle = await response.json()
            return this.convertToCamelCase(data)
        } catch (error) {
            throw error
        }
    }

    // 記事を更新
    async update(id: number, article: UpdateArticleInput): Promise<Article> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/articles/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(article),
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to update article')
            }

            const data: ApiArticle = await response.json()
            return this.convertToCamelCase(data)
        } catch (error) {
            throw error
        }
    }

    // 記事を削除
    async delete(id: number): Promise<void> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/articles/${id}`, {
                method: 'DELETE',
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to delete article')
            }
        } catch (error) {
            throw error
        }
    }

    // APIレスポンスをフロントエンド用に変換
    private convertToCamelCase(apiArticle: ApiArticle): Article {
        return {
            id: apiArticle.id,
            title: apiArticle.title,
            url: apiArticle.url,
            summary: apiArticle.summary,
            tags: apiArticle.tags ?? [],
            memo: apiArticle.memo ?? '',
            createdAt: apiArticle.created_at,
            updatedAt: apiArticle.updated_at,
        }
    }
}

// シングルトンインスタンスをエクスポート
export const articleClient = new ArticleClient()
