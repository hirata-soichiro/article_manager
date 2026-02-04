import { Article, CreateArticleInput, UpdateArticleInput } from '@/types/article'
import { ApiError } from '@/lib/errors/ApiError'

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
    private async handleResponse<T>(
        response: Response,
        endpoint: string,
        method: string
    ): Promise<T> {
        if (!response.ok) {
            let errorMessage = 'An error occurred'
            let errorDetails: unknown

            try {
                const errorData = await response.json()
                errorMessage = errorData.error || errorData.message || errorMessage
                errorDetails = errorData
            } catch {
                errorMessage = response.statusText || errorMessage
            }

            throw new ApiError(
                errorMessage,
                response.status,
                endpoint,
                method,
                errorDetails
            )
        }

        try {
            return await response.json()
        } catch (error) {
            throw new ApiError(
                'Failed to parse response',
                response.status,
                endpoint,
                method,
                error
            )
        }
    }

    private async fetchWithErrorHandling<T>(
        endpoint: string,
        options?: RequestInit
    ): Promise<T> {
        const url = `${API_BASE_URL}${endpoint}`
        const method = options?.method || 'GET'

        try {
            const response = await fetch(url, options)
            return await this.handleResponse<T>(response, endpoint, method)
        } catch (error) {
            if (error instanceof ApiError) {
                throw error
            }

            if (error instanceof TypeError && error.message.includes('fetch')) {
                throw new ApiError(
                    'ネットワークエラーが発生しました。接続を確認してください',
                    0,
                    endpoint,
                    method,
                    error
                )
            }

            throw new ApiError(
                error instanceof Error ? error.message : 'Unknown error occurred',
                0,
                endpoint,
                method,
                error
            )
        }
    }
    // 全記事を取得
    async getAll(): Promise<Article[]> {
        const data = await this.fetchWithErrorHandling<ApiArticle[]>('/api/articles')
        return data.map(this.convertToCamelCase)
    }

    // 指定IDの記事を取得
    async getById(id: number): Promise<Article> {
        const data = await this.fetchWithErrorHandling<ApiArticle>(`/api/articles/${id}`)
        return this.convertToCamelCase(data)
    }

    // 記事を作成
    async create(article: CreateArticleInput): Promise<Article> {
        const data = await this.fetchWithErrorHandling<ApiArticle>('/api/articles', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(article),
        })
        return this.convertToCamelCase(data)
    }

    // 記事を更新
    async update(id: number, article: UpdateArticleInput): Promise<Article> {
        const data = await this.fetchWithErrorHandling<ApiArticle>(`/api/articles/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(article),
        })
        return this.convertToCamelCase(data)
    }

    // 記事を削除
    async delete(id: number): Promise<void> {
        await this.fetchWithErrorHandling<void>(`/api/articles/${id}`, {
            method: 'DELETE',
        })
    }

    // キーワードで記事を検索
    async searchArticles(keyword: string): Promise<Article[]> {
        const encodedKeyword = encodeURIComponent(keyword)
        const data = await this.fetchWithErrorHandling<ApiArticle[]>(
            `/api/articles/search?keyword=${encodedKeyword}`
        )
        return data.map(this.convertToCamelCase)
    }

    // URLから記事を自動生成
    async generate(url: string, memo?: string): Promise<Article> {
        const requestBody: { url: string; memo?: string } = { url }
        if (memo !== undefined) {
            requestBody.memo = memo
        }

        const data = await this.fetchWithErrorHandling<ApiArticle>('/api/articles/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
        })
        return this.convertToCamelCase(data)
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
