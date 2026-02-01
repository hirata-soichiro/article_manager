'use client'

import { useState, useMemo } from 'react'
import { useArticles } from '@/hooks/useArticles'
import { useArticleSearch } from '@/hooks/useArticleSearch'
import { useTags } from '@/hooks/useTags'
import ArticleCard from '@/components/ArticleCard'
import SearchBar from '@/components/SearchBar'
import DeleteConfirmDialog from '@/components/DeleteConfirmDialog'
import { Article } from '@/types/article'

export default function ArticlesPage() {
    // 記事一覧フック
    const { articles, loading: articlesLoading, error: articlesError, deleteArticle } = useArticles()
    // 検索機能フック
    const { results: searchResults, loading: searchLoading, error: searchError, keyword, search, clearSearch } = useArticleSearch()
    // タグ一覧フック
    const { tags, loading: tagsLoading } = useTags()
    // タグフィルター用の状態
    const [selectedTag, setSelectedTag] = useState<string | null>(null)
    // タグ表示の展開状態
    const [showAllTags, setShowAllTags] = useState(false)
    // 削除確認ダイアログの状態
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
    const [articleToDelete, setArticleToDelete] = useState<number | null>(null)
    const [isDeleting, setIsDeleting] = useState(false)
    // ページネーションの状態
    const [currentPage, setCurrentPage] = useState(1)
    const [pageSize, setPageSize] = useState(10)
    // 検索中かどうかの判定
    const isSearching = keyword.trim().length > 0
    // 表示する記事リストを決定
    const displayArticles = isSearching ? searchResults : articles
    // ローディング状態を統合
    const loading = isSearching ? searchLoading : (articlesLoading || tagsLoading)
    // エラー状態を統合
    const error = isSearching ? searchError : articlesError

    // 選択されたタグに基づいて記事をフィルタリング
    const filteredArticles = useMemo(() => {
        if (!selectedTag) {
            // タグが選択されていない場合は全記事を表示
            return displayArticles
        }
        return displayArticles.filter(article => article.tags.includes(selectedTag))
    }, [displayArticles, selectedTag])

    // ページネーション計算
    const totalPages = Math.ceil(filteredArticles.length / pageSize)
    const startIndex = (currentPage - 1) * pageSize
    const endIndex = startIndex + pageSize
    const paginatedArticles = filteredArticles.slice(startIndex, endIndex)

    // ページサイズ変更時に1ページ目に戻る
    const handlePageSizeChange = (newPageSize: number) => {
        setPageSize(newPageSize)
        setCurrentPage(1)
    }

    // ページ変更
    const handlePageChange = (newPage: number) => {
        setCurrentPage(newPage)
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    // 検索処理のハンドラー
    const handleSearch = async (searchKeyword: string) => {
        setCurrentPage(1) // 検索時は1ページ目に戻る
        if (searchKeyword.trim().length === 0) {
            clearSearch()
            setSelectedTag(null)
        } else {
            setSelectedTag(null)
            await search(searchKeyword)
        }
    }

    // 検索クリアのハンドラー
    const handleClearSearch = () => {
        clearSearch()
        setSelectedTag(null)
    }

    // 削除ダイアログを開く
    const handleDeleteClick = (id: number) => {
        setArticleToDelete(id)
        setDeleteDialogOpen(true)
    }

    // 削除をキャンセル
    const handleCancelDelete = () => {
        setDeleteDialogOpen(false)
        setArticleToDelete(null)
    }

    // 削除を実行
    const handleConfirmDelete = async () => {
        if (articleToDelete === null) return

        try {
            setIsDeleting(true)
            await deleteArticle(articleToDelete)
            setDeleteDialogOpen(false)
            setArticleToDelete(null)

            // 削除後、現在のページに記事がなくなった場合は前のページに移動
            if (paginatedArticles.length === 1 && currentPage > 1) {
                setCurrentPage(currentPage - 1)
            }
        } catch (err) {
            console.error('削除に失敗しました', err)
        } finally {
            setIsDeleting(false)
        }
    }

    // ローディング中の表示
    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <div className="text-gray-600">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-t-2 border-purple-600 mx-auto mb-4"></div>
                    <p className="font-medium">{isSearching ? '検索中...' : '記事を読み込んでいます...'}</p>
                </div>
            </div>
        )
    }

    // エラー発生時の表示
    if (error) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <div className="text-center">
                    <div className="text-red-600 mb-4">
                        <svg className="w-16 h-16 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <h2 className="text-xl font-bold mb-2">エラーが発生しました</h2>
                        <p className="text-gray-600">{error.message}</p>
                    </div>
                    <div className="mt-4 text-sm text-gray-500">
                        <p>バックエンドAPIが起動しているか確認してください。</p>
                        <p className="mt-1">期待されるエンドポイント: http://localhost:8080/api/articles</p>
                    </div>
                    {isSearching && (
                        <button
                            onClick={handleClearSearch}
                            className="mt-4 px-4 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg hover:shadow-lg hover:scale-105 transition"
                        >
                            検索をクリア
                        </button>
                    )}
                </div>
            </div>
        )
    }

    // 記事が0件の場合の表示
    if (displayArticles.length === 0 && !isSearching) {
        return (
            <div>
                <div className="mb-8">
                    <SearchBar onSearch={handleSearch} />
                </div>
                <div className="flex justify-center items-center min-h-[500px]">
                    <div className="text-center max-w-md">
                        <div className="relative mb-8">
                            <div className="absolute inset-0 flex items-center justify-center">
                                <div className="w-32 h-32 bg-gradient-to-br from-blue-100 to-purple-100 rounded-full animate-pulse"></div>
                            </div>
                            <svg className="w-32 h-32 mx-auto relative text-gradient-to-br from-blue-500 to-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <h2 className="text-3xl font-bold mb-3 bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">まだ記事がありません</h2>
                        <p className="text-gray-600 mb-6 leading-relaxed">最初の記事を作成して、あなたのナレッジを蓄積していきましょう</p>
                        <a
                            href="/articles/new"
                            className="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-200"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                            記事を作成する
                        </a>
                    </div>
                </div>
            </div>
        )
    }

    // 記事一覧の表示
    return (
        <div>
            {/* 検索バー */}
            <div className="mb-8">
                <SearchBar onSearch={handleSearch} disabled={loading} placeholder="例: React, TypeScript, API..." />
            </div>

            {/* ヘッダー部分 */}
            <div className="mb-8">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold text-gray-800">
                            {isSearching ? '検索結果' : '記事一覧'}
                        </h1>
                        <p className="text-gray-600 mt-2">
                            {isSearching ? (
                                <>「{keyword}」の検索結果: {filteredArticles.length}件{selectedTag && ` (「${selectedTag}」でフィルタ中)`}</>
                            ) : (
                                selectedTag ? `「${selectedTag}」のタグを持つ記事: ${filteredArticles.length}件` : `全 ${articles.length} 件の記事`
                            )}
                        </p>
                    </div>
                    {isSearching && (
                        <button
                            onClick={handleClearSearch}
                            className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-all duration-200 flex items-center gap-2"
                        >
                            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                            検索をクリア
                        </button>
                    )}
                </div>
            </div>

            {/* タグフィルター */}
            {tags.length > 0 && displayArticles.length > 0 && (
                <div className="mb-8 bg-white rounded-lg shadow-sm p-6 border border-gray-200">
                    <div className="flex items-center gap-2 mb-4">
                        <svg className="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                        </svg>
                        <h2 className="text-base font-semibold text-gray-800">
                            タグで絞り込み{isSearching && <span className="text-gray-500 font-normal ml-2">（検索結果内）</span>}
                        </h2>
                    </div>
                    <div className="flex flex-wrap gap-2">
                        <button
                            onClick={() => {
                                setSelectedTag(null)
                                setCurrentPage(1)
                            }}
                            className={`px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${selectedTag === null ? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-md' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}`}
                        >
                            すべて ({displayArticles.length})
                        </button>
                        {(() => {
                            // タグを記事数でソートして、件数が0のタグを除外
                            const tagsWithCount = tags
                                .map(tag => ({
                                    ...tag,
                                    count: displayArticles.filter(article => article.tags.includes(tag.name)).length
                                }))
                                .filter(tag => tag.count > 0)
                                .sort((a, b) => b.count - a.count)

                            // 表示するタグを決定（最初は8個、もっと見るで全部）
                            const displayTags = showAllTags ? tagsWithCount : tagsWithCount.slice(0, 8)
                            const hasMoreTags = tagsWithCount.length > 8

                            return (
                                <>
                                    {displayTags.map(tag => {
                                        // 記事数に応じた背景色の濃淡を計算
                                        const opacity = Math.min(100, 40 + (tag.count * 10))
                                        const isBold = tag.count >= 5

                                        return (
                                            <button
                                                key={tag.id}
                                                onClick={() => {
                                                    setSelectedTag(tag.name)
                                                    setCurrentPage(1)
                                                }}
                                                style={{
                                                    backgroundColor: selectedTag === tag.name ? undefined : `rgba(107, 114, 128, ${opacity / 100})`
                                                }}
                                                className={`px-4 py-2 rounded-full text-sm transition-all duration-200 hover:shadow-md hover:scale-105 ${
                                                    selectedTag === tag.name
                                                        ? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-md font-semibold'
                                                        : `text-gray-700 hover:bg-gray-300 ${isBold ? 'font-semibold' : 'font-medium'}`
                                                }`}
                                            >
                                                {tag.name} ({tag.count})
                                            </button>
                                        )
                                    })}
                                    {hasMoreTags && (
                                        <button
                                            onClick={() => setShowAllTags(!showAllTags)}
                                            className="px-4 py-2 rounded-full text-sm font-medium text-indigo-600 border-2 border-indigo-600 hover:bg-indigo-50 transition-all duration-200 hover:shadow-md"
                                        >
                                            {showAllTags ? (
                                                <span className="flex items-center gap-1">
                                                    <span>閉じる</span>
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                                                    </svg>
                                                </span>
                                            ) : (
                                                <span className="flex items-center gap-1">
                                                    <span>もっと見る (+{tagsWithCount.length - 8})</span>
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                                                    </svg>
                                                </span>
                                            )}
                                        </button>
                                    )}
                                </>
                            )
                        })()}
                    </div>
                </div>
            )}

            {/* ページサイズ選択と表示情報 */}
            {filteredArticles.length > 0 && (
                <div className="mb-6 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <span className="text-sm text-gray-600">
                            {startIndex + 1} - {Math.min(endIndex, filteredArticles.length)} 件を表示 (全 {filteredArticles.length} 件)
                        </span>
                    </div>
                    <div className="flex items-center gap-3">
                        <span className="text-sm text-gray-600">表示件数:</span>
                        <select
                            value={pageSize}
                            onChange={(e) => handlePageSizeChange(Number(e.target.value))}
                            className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        >
                            <option value={10}>10件</option>
                            <option value={25}>25件</option>
                            <option value={50}>50件</option>
                            <option value={100}>100件</option>
                        </select>
                    </div>
                </div>
            )}

            {/* 記事一覧 */}
            {filteredArticles.length === 0 ? (
                <div className="flex justify-center items-center min-h-[200px]">
                    <div className="text-center text-gray-500">
                        {isSearching ? (
                            <>
                                <svg className="w-16 h-16 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                </svg>
                                <p className="text-lg mb-2">{selectedTag ? `「${keyword}」かつ「${selectedTag}」タグの記事が見つかりません` : `「${keyword}」に一致する記事が見つかりません`}</p>
                                <p className="text-sm text-gray-400">別のキーワードで検索してみてください</p>
                                <button onClick={handleClearSearch} className="mt-4 px-4 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg hover:shadow-lg hover:scale-105 transition">検索をクリア</button>
                            </>
                        ) : (
                            <>
                                <p className="text-lg">「{selectedTag}」タグの記事が見つかりません</p>
                                <button onClick={() => setSelectedTag(null)} className="mt-4 px-4 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg hover:shadow-lg hover:scale-105 transition">すべての記事を表示</button>
                            </>
                        )}
                    </div>
                </div>
            ) : (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {paginatedArticles.map((article) => (
                            <ArticleCard key={article.id} article={article} onDelete={handleDeleteClick} />
                        ))}
                    </div>

                    {/* ページネーション */}
                    {totalPages > 1 && (
                        <div className="mt-8 flex items-center justify-center gap-2">
                            {/* 前へボタン */}
                            <button
                                onClick={() => handlePageChange(currentPage - 1)}
                                disabled={currentPage === 1}
                                className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 flex items-center gap-2"
                            >
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                                </svg>
                                <span className="hidden sm:inline">前へ</span>
                            </button>

                            {/* ページ番号 */}
                            <div className="flex items-center gap-1">
                                {(() => {
                                    const pageNumbers = []
                                    const maxVisible = 5
                                    let startPage = Math.max(1, currentPage - Math.floor(maxVisible / 2))
                                    let endPage = Math.min(totalPages, startPage + maxVisible - 1)

                                    if (endPage - startPage < maxVisible - 1) {
                                        startPage = Math.max(1, endPage - maxVisible + 1)
                                    }

                                    // 最初のページ
                                    if (startPage > 1) {
                                        pageNumbers.push(
                                            <button
                                                key={1}
                                                onClick={() => handlePageChange(1)}
                                                className="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-all duration-200"
                                            >
                                                1
                                            </button>
                                        )
                                        if (startPage > 2) {
                                            pageNumbers.push(<span key="ellipsis1" className="px-2 text-gray-400">...</span>)
                                        }
                                    }

                                    // 中間のページ
                                    for (let i = startPage; i <= endPage; i++) {
                                        pageNumbers.push(
                                            <button
                                                key={i}
                                                onClick={() => handlePageChange(i)}
                                                className={`px-3 py-2 border rounded-lg transition-all duration-200 ${
                                                    currentPage === i
                                                        ? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white border-transparent shadow-md'
                                                        : 'border-gray-300 hover:bg-gray-50'
                                                }`}
                                            >
                                                {i}
                                            </button>
                                        )
                                    }

                                    // 最後のページ
                                    if (endPage < totalPages) {
                                        if (endPage < totalPages - 1) {
                                            pageNumbers.push(<span key="ellipsis2" className="px-2 text-gray-400">...</span>)
                                        }
                                        pageNumbers.push(
                                            <button
                                                key={totalPages}
                                                onClick={() => handlePageChange(totalPages)}
                                                className="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-all duration-200"
                                            >
                                                {totalPages}
                                            </button>
                                        )
                                    }

                                    return pageNumbers
                                })()}
                            </div>

                            {/* 次へボタン */}
                            <button
                                onClick={() => handlePageChange(currentPage + 1)}
                                disabled={currentPage === totalPages}
                                className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 flex items-center gap-2"
                            >
                                <span className="hidden sm:inline">次へ</span>
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </button>
                        </div>
                    )}
                </>
            )}

            {/* 削除確認ダイアログ */}
            <DeleteConfirmDialog
                isOpen={deleteDialogOpen}
                title="記事を削除しますか？"
                message="この操作は取り消せません。本当に削除してもよろしいですか？"
                onConfirm={handleConfirmDelete}
                onCancel={handleCancelDelete}
                isDeleting={isDeleting}
            />
        </div>
    )
}
