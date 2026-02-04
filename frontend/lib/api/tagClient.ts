import { Tag, CreateTagInput, UpdateTagInput } from '@/types/tag'
import { ApiError } from '@/lib/errors/ApiError'

// APIから返却される記事データの型
interface ApiTag {
    id: number
    name: string
    created_at: string
    updated_at: string
}

// APIのベースURL
const API_BASE_URL = 'http://localhost:8080'

// バックエンドのAPIと通信するクライアント
class TagClient {
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
    // 全タグを取得
    async getAll(): Promise<Tag[]> {
        const data = await this.fetchWithErrorHandling<ApiTag[]>('/api/tags')
        return data.map(this.convertToCamelCase)
    }

    // 指定IDのタグを取得
    async getById(id: number): Promise<Tag> {
        const data = await this.fetchWithErrorHandling<ApiTag>(`/api/tags/${id}`)
        return this.convertToCamelCase(data)
    }

    // タグを作成
    async create(tag: CreateTagInput): Promise<Tag> {
        const data = await this.fetchWithErrorHandling<ApiTag>('/api/tags', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(tag),
        })
        return this.convertToCamelCase(data)
    }

    // タグを更新
    async update(id: number, tag: UpdateTagInput): Promise<Tag> {
        const data = await this.fetchWithErrorHandling<ApiTag>(`/api/tags/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(tag),
        })
        return this.convertToCamelCase(data)
    }

    // タグを削除
    async delete(id: number): Promise<void> {
        await this.fetchWithErrorHandling<void>(`/api/tags/${id}`, {
            method: 'DELETE',
        })
    }

    // APIレスポンスをフロントエンド用に変換
    private convertToCamelCase(apiTag: ApiTag): Tag {
        return {
            id: apiTag.id,
            name: apiTag.name,
            createdAt: apiTag.created_at,
            updatedAt: apiTag.updated_at,
        }
    }
}

// シングルトンインスタンスをエクスポート
export const tagClient = new TagClient()
