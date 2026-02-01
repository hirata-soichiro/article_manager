'use client'

import { useState, useEffect } from 'react'
import { Tag, CreateTagInput, UpdateTagInput } from '@/types/tag'

interface TagWithCount extends Tag {
    articleCount: number
}

export default function TagsPage() {
    const [tags, setTags] = useState<TagWithCount[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [newTagName, setNewTagName] = useState('')
    const [editingTag, setEditingTag] = useState<TagWithCount | null>(null)
    const [editingName, setEditingName] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [searchQuery, setSearchQuery] = useState('')

    // タグ一覧の表示（記事数を含む）
    const fetchTags = async () => {
        try {
            setLoading(true)

            // タグと記事を並行して取得
            const [tagsResponse, articlesResponse] = await Promise.all([
                fetch(`http://localhost:8080/api/tags`),
                fetch(`http://localhost:8080/api/articles`)
            ])

            if (!tagsResponse.ok) {
                throw new Error('タグの取得に失敗しました')
            }
            if (!articlesResponse.ok) {
                throw new Error('記事の取得に失敗しました')
            }

            const tagsData = await tagsResponse.json()
            const articlesData = await articlesResponse.json()

            // 各タグの使用回数をカウント
            const tagCounts: Record<string, number> = {}
            if (Array.isArray(articlesData)) {
                articlesData.forEach((article: any) => {
                    if (article.tags && Array.isArray(article.tags)) {
                        article.tags.forEach((tagName: string) => {
                            tagCounts[tagName] = (tagCounts[tagName] || 0) + 1
                        })
                    }
                })
            }

            // タグに記事数を追加
            const tagsWithCount: TagWithCount[] = (tagsData || []).map((tag: Tag) => ({
                ...tag,
                articleCount: tagCounts[tag.name] || 0
            }))

            // 記事数の多い順にソート
            tagsWithCount.sort((a, b) => b.articleCount - a.articleCount)

            setTags(tagsWithCount)
            setError(null)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'エラーが発生しました')
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchTags()
    }, [])

    // タグを追加
    const handleCreateTag = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!newTagName.trim()) return

        try {
            setIsSubmitting(true)
            const input: CreateTagInput = { name: newTagName.trim() }
            const response = await fetch(`http://localhost:8080/api/tags`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(input),
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'タグの作成に失敗しました')
            }

            setNewTagName('')
            await fetchTags()
        } catch (err) {
            alert(err instanceof Error ? err.message : 'エラーが発生しました')
        } finally {
            setIsSubmitting(false)
        }
    }

    // タグを更新
    const handleUpdateTag = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!editingTag || !editingName.trim()) return

        try {
            setIsSubmitting(true)
            const input: UpdateTagInput = { name: editingName.trim() }
            const response = await fetch(`http://localhost:8080/api/tags/${editingTag.id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(input),
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'タグの更新に失敗しました')
            }

            setEditingTag(null)
            setEditingName('')
            await fetchTags()
        } catch (err) {
            alert(err instanceof Error ? err.message : 'エラーが発生しました')
        } finally {
            setIsSubmitting(false)
        }
    }

    // タグを削除
    const handleDeleteTag = async (id: number) => {
        if (!confirm('このタグを削除してもよろしいですか？\n関連する記事からも削除されます。')) {
            return
        }

        try {
            const response = await fetch(`http://localhost:8080/api/tags/${id}`, {
                method: 'DELETE',
            })

            if (!response.ok) {
                throw new Error('タグの削除に失敗しました')
            }

            await fetchTags()
        } catch (err) {
            alert(err instanceof Error ? err.message : 'エラーが発生しました')
        }
    }

    // 編集モードに切り替え
    const startEditing = (tag: TagWithCount) => {
        setEditingTag(tag)
        setEditingName(tag.name)
    }

    // 検索でフィルタリング
    const filteredTags = tags.filter(tag =>
        tag.name.toLowerCase().includes(searchQuery.toLowerCase())
    )

    // 編集をキャンセル
    const cancelEditing = () => {
        setEditingTag(null)
        setEditingName('')
    }

    // ローディング中の表示
    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <div className="text-gray-600">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p>タグを読み込んでいます...</p>
                </div>
            </div>
        )
    }

    // エラー発生時
    if (error) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <div className="text-center text-red-600">
                    <p className="text-xl font-bold mb-2">エラーが発生しました</p>
                    <p>{error}</p>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-50 pb-12">
            <div className="max-w-5xl mx-auto px-4 py-8">
                {/* ヘッダー */}
                <div className="mb-8">
                    <h1 className="text-4xl font-bold text-gray-900 mb-3">タグ管理</h1>
                    <div className="flex items-center gap-6 text-sm text-gray-600">
                        <div className="flex items-center gap-2">
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                            </svg>
                            <span>全 {tags.length} 件のタグ</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                            <span>{tags.filter(t => t.articleCount === 0).length} 件が未使用</span>
                        </div>
                    </div>
                </div>

                {/* タグ追加フォーム */}
                <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border border-gray-200">
                    <div className="flex items-center gap-3 mb-4">
                        <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                            <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                        </div>
                        <h2 className="text-xl font-bold text-gray-900">新しいタグを追加</h2>
                    </div>
                    <form onSubmit={handleCreateTag} className="flex gap-3">
                        <input
                            type="text"
                            value={newTagName}
                            onChange={(e) => setNewTagName(e.target.value)}
                            placeholder="タグ名を入力"
                            className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                            disabled={isSubmitting}
                            maxLength={50}
                        />
                        <button
                            type="submit"
                            disabled={!newTagName.trim() || isSubmitting}
                            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition font-semibold shadow-md hover:shadow-lg"
                        >
                            {isSubmitting ? '追加中...' : '追加'}
                        </button>
                    </form>
                </div>

                {/* タグ一覧 */}
                <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200">
                    <div className="flex items-center justify-between mb-6">
                        <h2 className="text-xl font-bold text-gray-900">タグ一覧</h2>
                        {/* 検索ボックス */}
                        <div className="relative w-64">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <svg className="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                </svg>
                            </div>
                            <input
                                type="text"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                                placeholder="タグを検索..."
                            />
                        </div>
                    </div>

                    {tags.length === 0 ? (
                        <div className="text-center py-12 text-gray-500">
                            <svg className="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                            </svg>
                            <p className="text-lg font-medium mb-2">タグがまだ登録されていません</p>
                            <p className="text-sm">上のフォームから最初のタグを追加してみましょう</p>
                        </div>
                    ) : filteredTags.length === 0 ? (
                        <div className="text-center py-12 text-gray-500">
                            <p className="text-lg">「{searchQuery}」に一致するタグが見つかりません</p>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {filteredTags.map((tag) => (
                                <div
                                    key={tag.id}
                                    className={`flex items-center justify-between p-4 border rounded-lg transition-all duration-200 ${
                                        tag.articleCount === 0
                                            ? 'border-gray-300 bg-gray-50'
                                            : 'border-gray-200 bg-white hover:shadow-md'
                                    }`}
                                >
                                    {editingTag?.id === tag.id ? (
                                        // 編集モード
                                        <form onSubmit={handleUpdateTag} className="flex-1 flex items-center gap-4">
                                            <input
                                                type="text"
                                                value={editingName}
                                                onChange={(e) => setEditingName(e.target.value)}
                                                className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                                                disabled={isSubmitting}
                                                maxLength={50}
                                                autoFocus
                                            />
                                            <div className="flex gap-2">
                                                <button
                                                    type="submit"
                                                    disabled={!editingName.trim() || isSubmitting}
                                                    className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition text-sm font-medium"
                                                >
                                                    {isSubmitting ? '保存中...' : '保存'}
                                                </button>
                                                <button
                                                    type="button"
                                                    onClick={cancelEditing}
                                                    disabled={isSubmitting}
                                                    className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 disabled:bg-gray-100 transition text-sm font-medium"
                                                >
                                                    キャンセル
                                                </button>
                                            </div>
                                        </form>
                                    ) : (
                                        // 表示モード
                                        <>
                                            <div className="flex-1 flex items-center gap-4">
                                                <div className="flex items-center gap-3">
                                                    <span className="text-lg font-semibold text-gray-900">
                                                        {tag.name}
                                                    </span>
                                                    <span className={`px-3 py-1 rounded-full text-xs font-bold ${
                                                        tag.articleCount === 0
                                                            ? 'bg-gray-200 text-gray-600'
                                                            : 'bg-blue-100 text-blue-700'
                                                    }`}>
                                                        {tag.articleCount} 件
                                                    </span>
                                                    {tag.articleCount === 0 && (
                                                        <span className="text-xs text-gray-400 italic">未使用</span>
                                                    )}
                                                </div>
                                            </div>
                                            <div className="flex gap-2">
                                                <button
                                                    onClick={() => startEditing(tag)}
                                                    className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
                                                    aria-label="編集"
                                                >
                                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                                                    </svg>
                                                </button>
                                                <button
                                                    onClick={() => handleDeleteTag(tag.id)}
                                                    className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
                                                    aria-label="削除"
                                                >
                                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                                    </svg>
                                                </button>
                                            </div>
                                        </>
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}
