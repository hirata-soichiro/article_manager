import { Tag, CreateTagInput, UpdateTagInput } from '@/types/tag'
import { BaseApiClient } from './baseClient'

// APIから返却されるタグデータの型
interface ApiTag {
    id: number
    name: string
    created_at: string
    updated_at: string
}

// バックエンドのAPIと通信するクライアント
class TagClient extends BaseApiClient {
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
