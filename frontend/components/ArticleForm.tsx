'use client'

import React, { useState, useEffect, FormEvent } from 'react'
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
    const [generatedTags, setGeneratedTags] = useState<string[]>([])

    // タグUI の状態
    const [tagSearch, setTagSearch] = useState('')
    const [showAllTags, setShowAllTags] = useState(false)

    // 新しいタグを作成して選択
    const createAndSelectTag = () => {
        const newTagName = tagSearch.trim()

        // バリデーション
        if (!newTagName) return
        if (newTagName.length > 50) {
            alert('タグ名は50文字以内にしてください')
            return
        }

        // 重複チェック（大文字小文字を区別しない）
        const isDuplicate = tags.some(tag =>
            tag.name.toLowerCase() === newTagName.toLowerCase()
        )

        if (isDuplicate) {
            alert('このタグは既に存在します')
            return
        }

        // 新しいタグをリストに追加
        const newTag: Tag = {
            id: 0, // 新規タグはIDが0
            name: newTagName,
            createdAt: '',
            updatedAt: ''
        }
        setTags([...tags, newTag])
        setGeneratedTags([...generatedTags, newTagName])

        // 選択済みタグに追加
        setSelectedTags([...selectedTags, newTagName])

        // 検索ボックスをクリア
        setTagSearch('')
    }

    // Enterキーで新しいタグを作成
    const handleTagSearchKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            e.preventDefault()
            const trimmedSearch = tagSearch.trim()
            if (!trimmedSearch) return

            // 既存タグに完全一致するものがあるか確認
            const existingTag = tags.find(tag =>
                tag.name.toLowerCase() === trimmedSearch.toLowerCase()
            )

            if (existingTag && !selectedTags.includes(existingTag.name)) {
                // 既存タグを選択
                toggleTag(existingTag.name)
                setTagSearch('')
            } else if (!existingTag) {
                // 新しいタグを作成
                createAndSelectTag()
            }
        }
    }

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

            const data = await articleClient.generate(url.trim(), memo.trim() || undefined)

            // フォームフィールドに自動入力
            setTitle(data.title)
            setSummary(data.summary)

            // タグの設定
            if (data.tags && Array.isArray(data.tags)) {
                // 既存のタグ名リストを取得
                const existingTagNames = tags.map(tag => tag.name)

                // 新規タグを抽出（既存タグに含まれないもの）
                const newTags = data.tags.filter(tagName => !existingTagNames.includes(tagName))

                // 新規タグをtagsリストに追加
                if (newTags.length > 0) {
                    const newTagObjects: Tag[] = newTags.map(tagName => ({
                        id: 0, // 新規タグはIDが0
                        name: tagName,
                        createdAt: '',
                        updatedAt: ''
                    }))
                    setTags([...tags, ...newTagObjects])
                    setGeneratedTags(newTags)
                }

                // 全てのタグを選択状態にする
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

                {/* AI自動生成セクション */}
                <div className="mb-6 bg-gradient-to-br from-green-50 to-emerald-50 border-2 border-green-200 rounded-xl p-6">
                    <div className="flex items-start gap-4">
                        <div className="flex-shrink-0 w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center">
                            <svg className="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                            </svg>
                        </div>
                        <div className="flex-1">
                            <h3 className="text-lg font-bold text-gray-900 mb-2">AIで記事情報を自動生成</h3>
                            <p className="text-sm text-gray-600 mb-4">URLを入力するだけで、タイトル、要約、関連タグを自動的に生成します</p>
                            <button
                                type="button"
                                onClick={handleAIGenerate}
                                disabled={!url.trim() || isGenerating || isSubmitting}
                                className="w-full sm:w-auto px-8 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition font-semibold shadow-md hover:shadow-lg flex items-center justify-center gap-2"
                            >
                                {isGenerating ? (
                                    <>
                                        <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                                        <span>生成中...</span>
                                    </>
                                ) : (
                                    <>
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                                        </svg>
                                        <span>AI自動生成する</span>
                                    </>
                                )}
                            </button>
                            {generateError && (
                                <p className="mt-3 text-sm text-red-600 flex items-center gap-2">
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    {generateError}
                                </p>
                            )}
                            {isGenerating && (
                                <p className="mt-3 text-sm text-green-700 flex items-center gap-2">
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    URLの内容を分析しています...
                                </p>
                            )}
                        </div>
                    </div>
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
                    <label className="block text-sm font-medium text-gray-700 mb-3">タグ</label>
                    {tagsLoading ? (
                        <div className="text-gray-500">タグを読み込んでいます...</div>
                    ) : (
                        <div>
                            {/* 選択済みタグ表示 */}
                            {selectedTags.length > 0 && (
                                <div className="mb-4 bg-blue-50 border border-blue-200 rounded-lg p-4">
                                    <div className="flex items-center gap-2 mb-2">
                                        <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                        </svg>
                                        <span className="text-sm font-semibold text-blue-900">選択済み ({selectedTags.length})</span>
                                    </div>
                                    <div className="flex flex-wrap gap-2">
                                        {selectedTags.map((tagName) => {
                                            const isNewTag = generatedTags.includes(tagName)
                                            return (
                                                <button
                                                    key={tagName}
                                                    type="button"
                                                    onClick={() => toggleTag(tagName)}
                                                    className="px-3 py-1.5 bg-blue-600 text-white rounded-full text-sm font-medium hover:bg-blue-700 transition-all duration-200 flex items-center gap-2 shadow-sm"
                                                >
                                                    <span>{tagName}</span>
                                                    {isNewTag && (
                                                        <span className="text-xs px-1.5 py-0.5 bg-yellow-400 text-yellow-900 rounded font-bold">NEW</span>
                                                    )}
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                                    </svg>
                                                </button>
                                            )
                                        })}
                                    </div>
                                </div>
                            )}

                            {/* 検索ボックス */}
                            <div className="mb-3">
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                        <svg className="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                        </svg>
                                    </div>
                                    <input
                                        type="text"
                                        value={tagSearch}
                                        onChange={(e) => setTagSearch(e.target.value)}
                                        onKeyDown={handleTagSearchKeyDown}
                                        className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                                        placeholder="タグを検索... (Enterで新規作成)"
                                    />
                                </div>
                            </div>

                            {/* タグ一覧 */}
                            <div className="flex flex-wrap gap-2">
                                {(() => {
                                    // 検索でフィルタリング
                                    const filteredTags = tags.filter(tag =>
                                        !selectedTags.includes(tag.name) &&
                                        tag.name.toLowerCase().includes(tagSearch.toLowerCase())
                                    )

                                    // 表示するタグを決定
                                    const displayTags = showAllTags ? filteredTags : filteredTags.slice(0, 8)
                                    const hasMore = filteredTags.length > 8

                                    // 新しいタグを作成可能か判定
                                    const trimmedSearch = tagSearch.trim()
                                    const canCreateNewTag = trimmedSearch && filteredTags.length === 0 &&
                                        !tags.some(tag => tag.name.toLowerCase() === trimmedSearch.toLowerCase())

                                    return (
                                        <>
                                            {displayTags.map((tag) => {
                                                const isNewTag = generatedTags.includes(tag.name)

                                                return (
                                                    <button
                                                        key={tag.id || tag.name}
                                                        type="button"
                                                        onClick={() => toggleTag(tag.name)}
                                                        className="px-3 py-1.5 bg-gray-100 text-gray-700 rounded-full text-sm font-medium hover:bg-gray-200 hover:shadow-sm transition-all duration-200 flex items-center gap-1.5"
                                                    >
                                                        <span>{tag.name}</span>
                                                        {isNewTag && (
                                                            <span className="text-xs px-1.5 py-0.5 bg-yellow-400 text-yellow-900 rounded font-bold">NEW</span>
                                                        )}
                                                    </button>
                                                )
                                            })}

                                            {/* 新しいタグを作成 */}
                                            {canCreateNewTag && (
                                                <button
                                                    type="button"
                                                    onClick={createAndSelectTag}
                                                    className="px-3 py-1.5 bg-green-50 border-2 border-green-500 text-green-700 rounded-full text-sm font-semibold hover:bg-green-100 transition-all duration-200 flex items-center gap-1.5"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                                                    </svg>
                                                    <span>「{trimmedSearch}」を新しいタグとして作成</span>
                                                </button>
                                            )}

                                            {/* もっと見るボタン */}
                                            {hasMore && !tagSearch && (
                                                <button
                                                    type="button"
                                                    onClick={() => setShowAllTags(!showAllTags)}
                                                    className="px-3 py-1.5 border-2 border-blue-600 text-blue-600 rounded-full text-sm font-medium hover:bg-blue-50 transition-all duration-200 flex items-center gap-1"
                                                >
                                                    {showAllTags ? (
                                                        <>
                                                            <span>閉じる</span>
                                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                                                            </svg>
                                                        </>
                                                    ) : (
                                                        <>
                                                            <span>もっと見る (+{filteredTags.length - 8})</span>
                                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                                                            </svg>
                                                        </>
                                                    )}
                                                </button>
                                            )}

                                            {filteredTags.length === 0 && !canCreateNewTag && (
                                                <p className="text-sm text-gray-500 py-2">
                                                    {tagSearch ? '該当するタグが見つかりません' : 'タグがありません'}
                                                </p>
                                            )}
                                        </>
                                    )
                                })()}
                            </div>
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
