'use client'

import { useState, useEffect } from 'react'
import { Tag, CreateTagInput, UpdateTagInput } from '@/types/tag'

export default function TagsPage() {
    const [tags, setTags] = useState<Tag[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [newTagName, setNewTagName] = useState('')
    const [editingTag, setEditingTag] = useState<Tag | null>(null)
    const [editingName, setEditingName] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)

    // タグ一覧の表示
    const fetchTags = async () => {
        try {
            setLoading(true)
            const response = await fetch(`http://localhost:8080/api/tags`)
            if (!response.ok) {
                throw new Error('タグの取得に失敗しました')
            }
            const data = await response.json()
            setTags(data || [])
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
    const startEditing = (tag: Tag) => {
        setEditingTag(tag)
        setEditingName(tag.name)
    }

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
        <div className="max-w-4xl mx-auto">
            {/* ヘッダー */}
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-800 mb-2">タグ管理</h1>
                <p className="text-gray-600">全 {tags.length} 件のタグ</p>
            </div>

            {/* タグ追加フォーム */}
            <div className="bg-white rounded-lg shadow-md p-6 mb-8">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">新しいタグを追加</h2>
                <form onSubmit={handleCreateTag} className="flex gap-4">
                    <input
                        type="text"
                        value={newTagName}
                        onChange={(e) => setNewTagName(e.target.value)}
                        placeholder="タグ名を入力"
                        className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        disabled={isSubmitting}
                        maxLength={50}
                    />
                    <button
                        type="submit"
                        disabled={!newTagName.trim() || isSubmitting}
                        className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition"
                    >
                        {isSubmitting ? '追加中...' : '追加'}
                    </button>
                </form>
            </div>

            {/* タグ一覧 */}
            <div className="bg-white rounded-lg shadow-md p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">タグ一覧</h2>

                {tags.length === 0 ? (
                    <div className="text-center py-8 text-gray-500">
                        <p>タグがまだ登録されていません</p>
                        <p className="text-sm mt-2">上のフォームから最初のタグを追加してみましょう</p>
                    </div>
                ) : (
                    <div className="space-y-2">
                        {tags.map((tag) => (
                            <div
                                key={tag.id}
                                className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition"
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
                                                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition text-sm"
                                            >
                                                {isSubmitting ? '保存中...' : '保存'}
                                            </button>
                                            <button
                                                type="button"
                                                onClick={cancelEditing}
                                                disabled={isSubmitting}
                                                className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 disabled:bg-gray-200 transition text-sm"
                                            >
                                                キャンセル
                                            </button>
                                        </div>
                                    </form>
                                ) : (
                                    // 表示モード
                                    <>
                                        <div className="flex-1">
                                            <span className="text-lg font-medium text-gray-800">
                                                {tag.name}
                                            </span>
                                        </div>
                                        <div className="flex gap-2">
                                            <button
                                                onClick={() => startEditing(tag)}
                                                className="px-4 py-2 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600 transition text-sm"
                                            >
                                                編集
                                            </button>
                                            <button
                                                onClick={() => handleDeleteTag(tag.id)}
                                                className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm"
                                            >
                                                削除
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
    )
}
