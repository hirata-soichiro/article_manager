import { Article, CreateArticleInput, UpdateArticleInput } from '@/types/article'
import { BaseApiClient } from './baseClient'

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

// バックエンドのAPIと通信するクライアント
class ArticleClient extends BaseApiClient {
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
