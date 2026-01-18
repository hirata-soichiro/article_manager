import { Tag, CreateTagInput, UpdateTagInput } from '@/types/tag'

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
    // 全タグを取得
    async getAll(): Promise<Tag[]> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/tags`)

            if (!response.ok) {
                throw new Error('Failed to fetch tags')
            }

            const data: ApiTag[] = await response.json()
            return data.map(this.convertToCamelCase)
        } catch (error) {
            throw error
        }
    }

    // 指定IDのタグを取得
    async getById(id: number): Promise<Tag> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/tags/${id}`)

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to fetch tag')
            }

            const data: ApiTag = await response.json()
            return this.convertToCamelCase(data)
        } catch (error) {
            throw error
        }
    }

    // タグを作成
    async create(tag: CreateTagInput): Promise<Tag> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/tags`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(tag),
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to create tag')
            }

            const data: ApiTag = await response.json()
            return this.convertToCamelCase(data)
        } catch (error) {
            throw error
        }
    }

    // タグを更新
    async update(id: number, tag: UpdateTagInput): Promise<Tag> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/tags/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(tag),
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to update tag')
            }

            const data: ApiTag = await response.json()
            return this.convertToCamelCase(data)
        } catch (error) {
            throw error
        }
    }

    // タグを削除
    async delete(id: number): Promise<void> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/tags/${id}`, {
                method: 'DELETE',
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Failed to delete tag')
            }
        } catch (error) {
            throw error
        }
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
