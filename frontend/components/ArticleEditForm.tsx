'use client'

import { useState, useEffect, FormEvent } from 'react'
import { useRouter, useParams } from 'next/navigation'
import { articleClient } from '@/lib/api/articleClient'
import { tagClient } from '@/lib/api/tagClient'
import { Tag } from '@/types/tag'
import { UpdateArticleInput } from '@/types/article'
import DeleteConfirmDialog from './DeleteConfirmDialog'

export default function ArticleEditForm() {
    const router = useRouter()
    const params = useParams()
    const articleId = Number(params.id)

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

    // 記事データの読み込み状態
    const [articleLoading, setArticleLoading] = useState(true)
    const [articleError, setArticleError] = useState<string | null>(null)

    // フォーム送信の状態
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [formError, setFormError] = useState<string | null>(null)

    // バリデーションエラーの状態
    const [titleError, setTitleError] = useState<string | null>(null)
    const [urlError, setUrlError] = useState<string | null>(null)
    const [summaryError, setSummaryError] = useState<string | null>(null)

    // 削除ダイアログの状態
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
    const [isDeleting, setIsDeleting] = useState(false)
    const [deleteError, setDeleteError] = useState<string | null>(null)

    // 削除処理
    const handleDeleteClick = () => {
        setIsDeleteDialogOpen(true)
        setDeleteError(null)
    }

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

    // 記事データを取得
    useEffect(() => {
        const fetchArticle = async () => {
            try {
                setArticleLoading(true)
                const article = await articleClient.getById(articleId)

                setTitle(article.title)
                setUrl(article.url)
                setSummary(article.summary)
                setMemo(article.memo)
                setSelectedTags(article.tags)

                setArticleError(null)
            } catch (err) {
                setArticleError('記事の読み込みに失敗しました')
            } finally {
                setArticleLoading(false)
            }
        }

        if (articleId) {
            fetchArticle()
        }
    }, [articleId])

    // バリデーション関数
    const validateTitle = (value: string) => {
        if (!value.trim()) {
            setTitleError('タイトルは必須です')
            return false
        }
        setTitleError(null)
        return true
    }

    const validateUrl = (value: string) => {
        if (!value.trim()) {
            setUrlError('URLは必須です')
            return false
        }
        try {
            new URL(value)
            setUrlError(null)
            return true
        } catch {
            setUrlError('正しいURL形式で入力してください')
            return false
        }
    }

    const validateSummary = (value: string) => {
        if (!value.trim()) {
            setSummaryError('要約は必須です')
            return false
        }
        setSummaryError(null)
        return true
    }

    const isFormValid = () => {
        return title.trim() !== '' && url.trim() !== '' && summary.trim() !== '' &&
                !titleError && !urlError && !summaryError
    }

    // タグの選択/解除
    const toggleTag = (tagName: string) => {
        setSelectedTags((prev) =>
            prev.includes(tagName) ? prev.filter((t) => t !== tagName) : [...prev, tagName]
        )
    }

    // フォーム送信
    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()

        const isTitleValid = validateTitle(title)
        const isUrlValid = validateUrl(url)
        const isSummaryValid = validateSummary(summary)

        if (!isTitleValid || !isUrlValid || !isSummaryValid) {
            return
        }

        try {
            setIsSubmitting(true)
            setFormError(null)

            const input: UpdateArticleInput = {
                title: title.trim(),
                url: url.trim(),
                summary: summary.trim(),
                tags: selectedTags,
                memo: memo.trim(),
            }

            await articleClient.update(articleId, input)
            router.push('/articles')
        } catch (err) {
            setFormError('記事の更新に失敗しました')
        } finally {
            setIsSubmitting(false)
        }
    }

    const handleCancel = () => {
        router.push('/articles')
    }

    const handleDeleteConfirm = async () => {
        try {
            setIsDeleting(true)
            setDeleteError(null)
            await articleClient.delete(articleId)
            router.push('/articles')
        } catch (err) {
            setDeleteError('記事の削除に失敗しました')
            setIsDeleting(false)
        }
    }

    const handleDeleteCancel = () => {
        setIsDeleteDialogOpen(false)
        setDeleteError(null)
    }

    if (tagsError) {
        return (
            <div className="max-w-4xl mx-auto">
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                    {tagsError}
                </div>
            </div>
        )
    }

    if (articleError) {
        return (
            <div className="max-w-4xl mx-auto">
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                    {articleError}
                </div>
            </div>
        )
    }

    if (articleLoading || tagsLoading) {
        return (
            <div className="max-w-4xl mx-auto">
                <div className="text-gray-500">読み込み中...</div>
            </div>
        )
    }

    return (
        <div className="max-w-4xl mx-auto">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-800 mb-2">記事編集</h1>
                <p className="text-gray-600">記事情報を編集します</p>
            </div>

            <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow-md p-6">
                {formError && (
                    <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                        {formError}
                    </div>
                )}

                {deleteError && (
                    <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                        {deleteError}
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
                            titleError ? 'border-red-500 focus:ring-red-500' : 'border-gray-300 focus:ring-blue-500'
                        }`}
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
                            urlError ? 'border-red-500 focus:ring-red-500' : 'border-gray-300 focus:ring-blue-500'
                        }`}
                        placeholder="https://example.com"
                    />
                    {urlError && <p className="mt-1 text-sm text-red-500">{urlError}</p>}
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
                            summaryError ? 'border-red-500 focus:ring-red-500' : 'border-gray-300 focus:ring-blue-500'
                        }`}
                        placeholder="記事の要約を入力"
                    />
                    {summaryError && <p className="mt-1 text-sm text-red-500">{summaryError}</p>}
                </div>

                {/* タグ選択 */}
                <div className="mb-6">
                    <label className="block text-sm font-medium text-gray-700 mb-2">タグ</label>
                    <div className="flex flex-wrap gap-2">
                        {tags.map((tag) => {
                            const isSelected = selectedTags.includes(tag.name)
                            return (
                                <button
                                    key={tag.id}
                                    type="button"
                                    onClick={() => toggleTag(tag.name)}
                                    className={`px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${
                                        isSelected
                                            ? 'bg-blue-600 text-white'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                                >
                                    {tag.name}
                                </button>
                            )
                        })}
                    </div>
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
                        {isSubmitting ? '更新中...' : '更新'}
                    </button>
                    <button
                        type="button"
                        onClick={handleCancel}
                        disabled={isSubmitting}
                        className="px-6 py-3 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 disabled:bg-gray-200 transition font-medium"
                    >
                        キャンセル
                    </button>
                    <button
                        type="button"
                        onClick={handleDeleteClick}
                        disabled={isSubmitting || isDeleting}
                        className="px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-red-400 disabled:cursor-not-allowed transition font-medium"
                    >
                        削除
                    </button>
                </div>
            </form>

            {/* 削除確認ダイアログ */}
            <DeleteConfirmDialog
                isOpen={isDeleteDialogOpen}
                title="記事を削除しますか？"
                message={`「${title}」を削除します。この操作は取り消せません。`}
                onConfirm={handleDeleteConfirm}
                onCancel={handleDeleteCancel}
                isDeleting={isDeleting}
            />
        </div>
    )
}
