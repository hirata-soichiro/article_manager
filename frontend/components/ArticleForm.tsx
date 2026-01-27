'use client'

import { useState, useEffect, FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import { articleClient } from '@/lib/api/articleClient'
import { tagClient } from '@/lib/api/tagClient'
import { Tag } from '@/types/tag'
import { CreateArticleInput } from '@/types/article'

export default function ArticleForm() {
    const router = useRouter()

    // フォームの状態管理
    const [title, setTitle] = useState('')
    const [url, setUrl] = useState('')
    const [summary, setSummary] = useState('')
    const [memo, setMemo] = useState('')
    const [selectedTags, setSelectedTags] = useState<string[]>([])

    // タグ一覧の状態
    const [tags, setTags] = useState<Tag[]>([])
    const [tagsLoading, setTagsLoading] = useState(true)
    const [tagsError, setTagsError] = useState<string | null>(null)

    // フォーム送信の状態
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [formError, setFormError] = useState<string | null>(null)

    // AI自動生成の状態
    const [isGenerating, setIsGenerating] = useState(false)
    const [generateError, setGenerateError] = useState<string | null>(null)

    // バリデーションエラーの状態
    const [titleError, setTitleError] = useState<string | null>(null)
    const [urlError, setUrlError] = useState<string | null>(null)
    const [summaryError, setSummaryError] = useState<string | null>(null)

    // タグ一覧を取得
    useEffect(() => {
        const fetchTags = async () => {
            try {
                setTagsLoading(true)
                const data = await tagClient.getAll()
                setTags(data)
                setTagsError(null)
            } catch (err) {
                setTagsError('タグの読み込みに失敗しました')
            } finally {
                setTagsLoading(false)
            }
        }

        fetchTags()
    }, [])

    // タイトルのバリデーション
    const validateTitle = (value: string) => {
        if (!value.trim()) {
            setTitleError('タイトルは必須です')
            return false
        }
        setTitleError(null)
        return true
    }

    // URLのバリデーション
    const validateUrl = (value: string) => {
        if (!value.trim()) {
            setUrlError('URLは必須です')
            return false
        }

        // URL形式のチェック
        try {
            new URL(value)
            setUrlError(null)
            return true
        } catch {
            setUrlError('正しいURL形式で入力してください')
            return false
        }
    }

    // 要約のバリデーション
    const validateSummary = (value: string) => {
        if (!value.trim()) {
            setSummaryError('要約は必須です')
            return false
        }
        setSummaryError(null)
        return true
    }

    // すべての入力が有効かチェック
    const isFormValid = () => {
        return (
            title.trim() !== '' &&
            url.trim() !== '' &&
            summary.trim() !== '' &&
            !titleError &&
            !urlError &&
            !summaryError
        )
    }

    // タグの選択/解除
    const toggleTag = (tagName: string) => {
        setSelectedTags((prev) => {
            if (prev.includes(tagName)) {
                // すでに選択されていたら解除
                return prev.filter((t) => t !== tagName)
            } else {
                // 選択されていなかったら追加
                return [...prev, tagName]
            }
        })
    }

    // AI自動生成
    const handleAIGenerate = async () => {
        if (!url.trim()) {
            setGenerateError('URLを入力してください')
            return
        }

        const isUrlValid = validateUrl(url)
        if (!isUrlValid) {
            setGenerateError('正しいURL形式で入力してください')
            return
        }

        try {
            setIsGenerating(true)
            setGenerateError(null)

            const response = await fetch('http://localhost:8080/api/articles/generate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    url: url.trim(),
                    memo: memo.trim(),
                }),
            })

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}))
                const errorMessage = errorData.error || 'AI生成に失敗しました'

                // エラーコード別のメッセージ
                if (response.status === 429) {
                    throw new Error('API制限に達しました。しばらく時間をおいてから再度お試しください')
                } else if (response.status === 504) {
                    throw new Error('タイムアウトしました。もう一度お試しください')
                } else if (response.status === 401) {
                    throw new Error('APIキーが無効です')
                } else if (response.status === 403) {
                    throw new Error('コンテンツがブロックされました')
                } else {
                    throw new Error(errorMessage)
                }
            }

            const data = await response.json()

            // フォームフィールドに自動入力
            setTitle(data.title)
            setSummary(data.summary)

            // タグの設定
            if (data.tags && Array.isArray(data.tags)) {
                setSelectedTags(data.tags)
            }

            // バリデーションエラーをクリア
            setTitleError(null)
            setSummaryError(null)
        } catch (err) {
            if (err instanceof Error) {
                setGenerateError(err.message)
            } else {
                setGenerateError('AI生成中にエラーが発生しました')
            }
        } finally {
            setIsGenerating(false)
        }
    }

    // フォーム送信
    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()

        // 最終バリデーション
        const isTitleValid = validateTitle(title)
        const isUrlValid = validateUrl(url)
        const isSummaryValid = validateSummary(summary)

        if (!isTitleValid || !isUrlValid || !isSummaryValid) {
            return
        }

        try {
            setIsSubmitting(true)
            setFormError(null)

            const input: CreateArticleInput = {
                title: title.trim(),
                url: url.trim(),
                summary: summary.trim(),
                tags: selectedTags,
                memo: memo.trim(),
            }

            await articleClient.create(input)

            // 成功したら記事一覧に戻る
            router.push('/articles')
        } catch (err) {
            setFormError('記事の作成に失敗しました')
        } finally {
            setIsSubmitting(false)
        }
    }

    // キャンセル
    const handleCancel = () => {
        router.push('/articles')
    }

    // タグ読み込みエラー時
    if (tagsError) {
        return (
            <div className="max-w-4xl mx-auto">
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                    {tagsError}
                </div>
            </div>
        )
    }

    return (
        <div className="max-w-4xl mx-auto">
            {/* ヘッダー */}
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-800 mb-2">記事登録</h1>
                <p className="text-gray-600">新しい記事を登録します</p>
            </div>

            {/* フォーム */}
            <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow-md p-6">
                {/* グローバルエラーメッセージ */}
                {formError && (
                    <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                        {formError}
                    </div>
                )}

                {/* タイトル */}
                <div className="mb-6">
                    <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                        タイトル <span className="text-red-500">*</span>
                    </label>
                    <input
                        type="text"
                        id="title"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        onBlur={() => validateTitle(title)}
                        className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
                            titleError
                                ? 'border-red-500 focus:ring-red-500'
                                : 'border-gray-300 focus:ring-blue-500'
                        }`}
                        placeholder="記事のタイトルを入力"
                    />
                    {titleError && <p className="mt-1 text-sm text-red-500">{titleError}</p>}
                </div>

                {/* URL */}
                <div className="mb-6">
                    <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-2">
                        URL <span className="text-red-500">*</span>
                    </label>
                    <input
                        type="text"
                        id="url"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        onBlur={() => validateUrl(url)}
                        className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
                            urlError
                                ? 'border-red-500 focus:ring-red-500'
                                : 'border-gray-300 focus:ring-blue-500'
                        }`}
                        placeholder="https://example.com"
                    />
                    {urlError && <p className="mt-1 text-sm text-red-500">{urlError}</p>}
                </div>

                {/* AI自動生成ボタン */}
                <div className="mt-3">
                    <button
                        type="button"
                        onClick={handleAIGenerate}
                        disabled={!url.trim() || isGenerating || isSubmitting}
                        className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition font-medium"
                    >
                        {isGenerating ? '生成中...' : 'AI自動生成'}
                    </button>
                    {generateError && (
                        <p className="mt-2 text-sm text-red-500">{generateError}</p>
                    )}
                    {isGenerating && (
                        <p className="mt-2 text-sm text-blue-600">URLの内容を分析しています...</p>
                    )}
                </div>

                {/* 要約 */}
                <div className="mb-6">
                    <label htmlFor="summary" className="block text-sm font-medium text-gray-700 mb-2">
                        要約 <span className="text-red-500">*</span>
                    </label>
                    <textarea
                        id="summary"
                        value={summary}
                        onChange={(e) => setSummary(e.target.value)}
                        onBlur={() => validateSummary(summary)}
                        rows={4}
                        className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
                            summaryError
                                ? 'border-red-500 focus:ring-red-500'
                                : 'border-gray-300 focus:ring-blue-500'
                        }`}
                        placeholder="記事の要約を入力"
                    />
                    {summaryError && <p className="mt-1 text-sm text-red-500">{summaryError}</p>}
                </div>

                {/* タグ選択 */}
                <div className="mb-6">
                    <label className="block text-sm font-medium text-gray-700 mb-2">タグ</label>
                    {tagsLoading ? (
                        <div className="text-gray-500">タグを読み込んでいます...</div>
                    ) : (
                        <div className="flex flex-wrap gap-2">
                            {tags.map((tag) => (
                                <button
                                    key={tag.id}
                                    type="button"
                                    onClick={() => toggleTag(tag.name)}
                                    className={`px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${
                                        selectedTags.includes(tag.name)
                                            ? 'bg-blue-600 text-white shadow-md'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                                >
                                    {tag.name}
                                </button>
                            ))}
                        </div>
                    )}
                </div>

                {/* メモ */}
                <div className="mb-6">
                    <label htmlFor="memo" className="block text-sm font-medium text-gray-700 mb-2">
                        メモ（任意）
                    </label>
                    <textarea
                        id="memo"
                        value={memo}
                        onChange={(e) => setMemo(e.target.value)}
                        rows={3}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="個人的なメモを入力"
                    />
                </div>

                {/* ボタン */}
                <div className="flex gap-4">
                    <button
                        type="submit"
                        disabled={!isFormValid() || isSubmitting}
                        className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition font-medium"
                    >
                        {isSubmitting ? '登録中...' : '登録'}
                    </button>
                    <button
                        type="button"
                        onClick={handleCancel}
                        disabled={isSubmitting}
                        className="px-6 py-3 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 disabled:bg-gray-200 transition font-medium"
                    >
                        キャンセル
                    </button>
                </div>
            </form>
        </div>
    )
}
