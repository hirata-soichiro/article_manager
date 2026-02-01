'use client'

import { on } from 'events'
import { HTMLInputElement } from 'happy-dom'
import { useState, FormEvent, KeyboardEvent } from 'react'

interface SearchBarProps {
    onSearch: (keyword: string) => void
    placeholder?: string
    disabled?: boolean
}

export default function SearchBar({
    onSearch,
    placeholder = 'キーワードで検索...',
    disabled = false,
}: SearchBarProps) {
    const [keyword, setKeyword] = useState('')

    // 検索実行
    const handleSearch = () => {
        onSearch(keyword.trim())
    }

    // フォーム送信時の処理
    const handleSubmit = (e: FormEvent) => {
        e.preventDefault()
        handleSearch()
    }

    // Enterキーでの検索
    const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            e.preventDefault()
            handleSearch()
        }
    }

    // クリアボタン
    const handleClear = () => {
        setKeyword('')
        onSearch('')
    }

    return (
        <form onSubmit={handleSubmit} className="w-full max-w-3xl mx-auto">
            <div className="flex gap-2 items-center">
                {/* 検索アイコン付き入力フィールド */}
                <div className="relative flex-1">
                    {/* 検索アイコン */}
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <svg
                            className="h-5 w-5 text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                            />
                        </svg>
                    </div>

                    {/* テキスト入力 */}
                    <input
                        type="text"
                        value={keyword}
                        onChange={(e) => setKeyword(e.target.value)}
                        onKeyDown={handleKeyDown}
                        disabled={disabled}
                        placeholder={placeholder}
                        className={`
                            w-full pl-10 pr-10 py-2
                            border border-gray-300 rounded-lg
                            focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                            disabled:bg-gray-100 disabled:cursor-not-allowed
                            transition-all duration-200
                        `}
                    />

                    {/* クリアボタン（入力がある場合のみ表示） */}
                    {keyword && !disabled && (
                        <button
                            type="button"
                            onClick={handleClear}
                            className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 transition-colors"
                            aria-label="検索キーワードをクリア"
                        >
                            <svg
                                className="h-5 w-5"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M6 18L18 6M6 6l12 12"
                                />
                            </svg>
                        </button>
                    )}
                </div>

                {/* 検索ボタン */}
                <button
                    type="submit"
                    disabled={disabled}
                    className={`
                        px-6 py-2 rounded-lg font-medium transition-all duration-200
                        flex items-center gap-2 whitespace-nowrap shadow-md
                        ${disabled
                            ? 'bg-gray-400 text-white cursor-not-allowed'
                            : 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white hover:shadow-lg hover:scale-105 active:scale-95'
                        }
                    `}
                >
                    <svg
                        className="h-5 w-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                        />
                    </svg>
                    <span className="hidden sm:inline">検索</span>
                </button>
            </div>
        </form>
    )
}
