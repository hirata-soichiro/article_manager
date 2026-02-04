'use client'

import React, { Component, ReactNode } from 'react'
import { ApiError } from '@/lib/errors/ApiError'

interface Props {
    children: ReactNode
}

interface State {
    hasError: boolean
    error: Error | null
}

export class ErrorBoundary extends Component<Props, State> {
    constructor(props: Props) {
        super(props)
        this.state = {
            hasError: false,
            error: null,
        }
    }

    static getDerivedStateFromError(error: Error): State {
        return {
            hasError: true,
            error,
        }
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        console.error('ErrorBoundary caught an error:', error, errorInfo)
    }

    render() {
        if (this.state.hasError && this.state.error) {
            const error = this.state.error
            const isApiError = error instanceof ApiError

            return (
                <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 via-white to-slate-50">
                    <div className="max-w-2xl w-full mx-4">
                        <div className="bg-white rounded-xl shadow-lg border border-red-100 p-8">
                            <div className="flex items-start space-x-4">
                                <div className="flex-shrink-0">
                                    <svg
                                        className="h-12 w-12 text-red-500"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            strokeLinecap="round"
                                            strokeLinejoin="round"
                                            strokeWidth={2}
                                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                                        />
                                    </svg>
                                </div>
                                <div className="flex-1">
                                    <h1 className="text-2xl font-bold text-gray-900 mb-2">
                                        エラーが発生しました
                                    </h1>
                                    <p className="text-gray-600 mb-4">
                                        {isApiError
                                            ? error.getUserMessage()
                                            : 'アプリケーションで予期しないエラーが発生しました'}
                                    </p>

                                    {isApiError && (
                                        <div className="bg-red-50 rounded-lg p-4 mb-4">
                                            <div className="text-sm space-y-1">
                                                <p className="text-red-800">
                                                    <span className="font-semibold">
                                                        ステータス:
                                                    </span>{' '}
                                                    {error.statusCode}
                                                </p>
                                                <p className="text-red-800">
                                                    <span className="font-semibold">
                                                        エンドポイント:
                                                    </span>{' '}
                                                    {error.endpoint}
                                                </p>
                                                <p className="text-red-800">
                                                    <span className="font-semibold">メソッド:</span>{' '}
                                                    {error.method}
                                                </p>
                                            </div>
                                        </div>
                                    )}

                                    <details className="mb-6">
                                        <summary className="cursor-pointer text-sm text-gray-500 hover:text-gray-700">
                                            技術的な詳細を表示
                                        </summary>
                                        <div className="mt-2 p-4 bg-gray-50 rounded-lg">
                                            <pre className="text-xs text-gray-700 overflow-x-auto">
                                                {error.stack || error.message}
                                            </pre>
                                        </div>
                                    </details>

                                    <div className="flex space-x-4">
                                        <button
                                            onClick={() => window.location.reload()}
                                            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                                        >
                                            ページを再読み込み
                                        </button>
                                        <button
                                            onClick={() => (window.location.href = '/')}
                                            className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
                                        >
                                            ホームに戻る
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )
        }

        return this.props.children
    }
}
