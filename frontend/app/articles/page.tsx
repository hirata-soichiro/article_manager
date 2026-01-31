'use client'

import { useState, useMemo } from 'react'
import { useArticles } from '@/hooks/useArticles'
import { useArticleSearch } from '@/hooks/useArticleSearch'
import { useTags } from '@/hooks/useTags'
import ArticleCard from '@/components/ArticleCard'
import SearchBar from '@/components/SearchBar'
import { Article } from '@/types/article'

export default function ArticlesPage() {
    // 記事一覧フック
    const { articles, loading: articlesLoading, error: articlesError } = useArticles()
    // 検索機能フック
    const { results: searchResults, loading: searchLoading, error: searchError, keyword, search, clearSearch } = useArticleSearch()
    // タグ一覧フック
    const { tags, loading: tagsLoading } = useTags()
    // タグフィルター用の状態
    const [selectedTag, setSelectedTag] = useState<string | null>(null)
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

    // 検索処理のハンドラー
    const handleSearch = async (searchKeyword: string) => {
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

    // ローディング中の表示
    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <div className="text-gray-600">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p>{isSearching ? '検索中...' : '記事を読み込んでいます...'}</p>
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
                            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
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
                <div className="mb-6">
                    <SearchBar onSearch={handleSearch} />
                </div>
                <div className="flex justify-center items-center min-h-[400px]">
                    <div className="text-center text-gray-600">
                        <svg className="w-16 h-16 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                        <h2 className="text-xl font-bold mb-2">記事がありません</h2>
                        <p>最初の記事を追加してみましょう</p>
                    </div>
                </div>
            </div>
        )
    }

    // 記事一覧の表示
    return (
        <div>
            {/* 検索バー */}
            <div className="mb-6">
                <SearchBar onSearch={handleSearch} disabled={loading} />
            </div>

            {/* ヘッダー部分 */}
            <div className="mb-6">
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
                <div className="mb-6 bg-white rounded-lg shadow-sm p-4">
                    <h2 className="text-sm font-semibold text-gray-700 mb-3">
                        タグで絞り込み{isSearching && <span className="text-gray-500 font-normal ml-2">（検索結果内）</span>}
                    </h2>
                    <div className="flex flex-wrap gap-2">
                        <button
                            onClick={() => setSelectedTag(null)}
                            className={`px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${selectedTag === null ? 'bg-blue-600 text-white shadow-md' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}`}
                        >
                            すべて ({displayArticles.length})
                        </button>
                        {tags.map(tag => {
                            const count = displayArticles.filter(article => article.tags.includes(tag.name)).length
                            if (count === 0) return null
                            return (
                                <button
                                    key={tag.id}
                                    onClick={() => setSelectedTag(tag.name)}
                                    className={`px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${selectedTag === tag.name ? 'bg-blue-600 text-white shadow-md' : 'bg-gray-100 text-gray-700 hover:bg-gray-200 hover:shadow-sm'}`}
                                >
                                    {tag.name} ({count})
                                </button>
                            )
                        })}
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
                                <button onClick={handleClearSearch} className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">検索をクリア</button>
                            </>
                        ) : (
                            <>
                                <p className="text-lg">「{selectedTag}」タグの記事が見つかりません</p>
                                <button onClick={() => setSelectedTag(null)} className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">すべての記事を表示</button>
                            </>
                        )}
                    </div>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {filteredArticles.map((article) => (
                        <ArticleCard key={article.id} article={article} />
                    ))}
                </div>
            )}
        </div>
    )
}
