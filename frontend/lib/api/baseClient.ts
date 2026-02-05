import { ApiError } from '@/lib/errors/ApiError'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

// 全APIクライアントの基底クラス
export class BaseApiClient {
    // レスポンス処理（エラー時はApiErrorをスロー）
    protected async handleResponse<T>(
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

        // 204 No Content（DELETE操作など）の処理
        if (response.status === 204) {
            return undefined as T
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

    // fetch実行とエラーハンドリング
    protected async fetchWithErrorHandling<T>(
        endpoint: string,
        options?: RequestInit
    ): Promise<T> {
        const url = `${API_BASE_URL}${endpoint}`
        const method = options?.method || 'GET'

        try {
            const response = options !== undefined
                ? await fetch(url, options)
                : await fetch(url)
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

    // APIベースURLを取得
    protected getBaseUrl(): string {
        return API_BASE_URL
    }
}
