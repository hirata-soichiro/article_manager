export class ApiError extends Error {
    public statusCode: number
    public endpoint: string
    public method: string
    public details?: unknown

    constructor(
        message: string,
        statusCode: number,
        endpoint: string,
        method: string,
        details?: unknown
    ) {
        super(message)
        this.name = 'ApiError'
        this.statusCode = statusCode
        this.endpoint = endpoint
        this.method = method
        this.details = details
        Object.setPrototypeOf(this, ApiError.prototype)
    }

    isClientError(): boolean {
        return this.statusCode >= 400 && this.statusCode < 500
    }

    isServerError(): boolean {
        return this.statusCode >= 500
    }

    isNotFound(): boolean {
        return this.statusCode === 404
    }

    isUnauthorized(): boolean {
        return this.statusCode === 401
    }

    isForbidden(): boolean {
        return this.statusCode === 403
    }

    isValidationError(): boolean {
        return this.statusCode === 422 || this.statusCode === 400
    }

    getUserMessage(): string {
        if (this.isNotFound()) {
            return 'リクエストされたリソースが見つかりませんでした'
        }
        if (this.isUnauthorized()) {
            return '認証が必要です'
        }
        if (this.isForbidden()) {
            return 'このリソースへのアクセス権限がありません'
        }
        if (this.isValidationError()) {
            return '入力内容に誤りがあります'
        }
        if (this.isServerError()) {
            return 'サーバーエラーが発生しました。時間をおいて再度お試しください'
        }
        return this.message || '予期しないエラーが発生しました'
    }
}
