'use client'

import { useState, useMemo } from 'react'
import { useArticles } from '@/hooks/useArticles'
import { useTags } from '@/hooks/useTags'
import ArticleCard from "@/components/ArticleCard"
import { Article } from '@/types/article'

export default function ArticlesPage() {
    const { articles, loading: articlesLoading, error: articlesError } = useArticles()
    const { tags, loading: tagsLoading } = useTags()
    const [selectedTag, setSelectedTag] = useState<string | null>(null)

    // ローディング状態を統合
    const loading = articlesLoading || tagsLoading

    // 選択されたタグに基づいて記事をフィルタリング
    const filteredArticles = useMemo(() => {
        if (!selectedTag) {
            // タグが選択されていない場合は全記事を表示
            return articles
        }
        return articles.filter(article => article.tags.includes(selectedTag))
    }, [articles, selectedTag])

    // ローディング中の表示
    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                {/* flex: フレックスボックスで配置 */}
                {/* justify-center: 横方向の中央揃え */}
                {/* items-center: 縦方向の中央揃え */}
                {/* min-h-[400px]: 最小の高さを400pxに設定 */}

                <div className="text-gray-600">
                    {/* text-gray-600: グレーの文字色 */}

                    {/* くるくる回るスピナー（読み込み中アニメーション） */}
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    {/* animate-spin: 回転アニメーション */}
                    {/* rounded-full: 完全な円形 */}
                    {/* h-12 w-12: 高さ・幅12単位（48px） */}
                    {/* border-b-2: 下側のボーダー（線）を2px */}
                    {/* border-blue-600: ボーダーの色を青に */}
                    {/* mx-auto: 左右の余白を自動（中央揃え） */}
                    {/* mb-4: 下の余白4単位 */}

                    <p>記事を読み込んでいます...</p>
                </div>
            </div>
        )
    }

    // エラー発生時の表示
    if (articlesError) {
        return (
            <div className="flex jusify-center items-center min-h-[400px]">
                <div className="text-center">
                    {/* text-center: テキストを中央揃え */}

                    <div className="text-red-600 mb-4">
                        {/* text-red-600: 赤色の文字（エラーを表現） */}

                        {/* エラーアイコン（丸い警告マーク） */}
                        <svg
                            className="w-16 h-16 mx-auto mb-2"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M12 8v4m0 4h.01M21 12a9 0 11-18 0 9 9 0 0118 0z"
                            />
                        </svg>

                        <h2 className="text-xl font-bold mb-2">エラーが発生しました</h2>
                        {/* text-xl: テキストサイズ大 */}
                        {/* font-bold: 太字 */}

                        <p className="text-gray-600">{articlesError.message}</p>
                        {/* エラーメッセージを表示 */}
                    </div>

                    {/* トラブルシューティング情報 */}
                    <div className="mt-4 text-sm text-gray-500">
                        {/* mt-4: 上の余白4単位 */}
                        {/* text-sm: 小さめのテキスト */}

                        <p>バックエンドAPIが起動しているか確認してください。</p>
                        <p className="mt-1">期待されるエンドポイント: http://localhost:8080/api/articles</p>
                    </div>
                </div>
            </div>
        )
    }

    // 記事が0件の場合の表示
    if (articles.length === 0) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <div className="text-center text-gray-600">
                    {/* 書類アイコン */}
                    <svg
                        className="w-16 h-16 mx-auto mb-4 text-gray-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                        />
                    </svg>

                    <h2 className="text-xl font-bold mb-2">記事がありません</h2>
                    <p>最初の記事を追加してみましょう</p>
                </div>
            </div>
        )
    }

    // 記事一覧の表示
    return (
        <div>
            {/* ヘッダー部分 */}
            <div className="mb-6">
                {/* mb-6: 下の余白6単位 */}

                <h1 className="text-3xl font-bold text-gray-800">
                    {/* text-3xl: 大きなテキスト */}
                    記事一覧
                </h1>

                <p className="text-gray-600 mt-2">
                    {/* mt-2: 上の余白2単位 */}
                    {selectedTag
                        ? `「${selectedTag}」のタグを持つ記事: ${filteredArticles.length}件`
                        : `全 ${articles.length} 件の記事`
                    }
                </p>
            </div>

            {/* タグフィルター */}
            {tags.length > 0 && (
                <div className="mb-6 bg-white rounded-lg shadow-sm p-4">
                    <h2 className="text-sm font-semibold text-gray-700 mb-3">タグで絞り込み</h2>
                    <div className="flex flex-wrap gap-2">
                        {/* 「すべて」ボタン */}
                        <button
                            onClick={() => setSelectedTag(null)}
                            className={`
                                px-4 py-2 rounded-full text-sm font-medium transition-all duration-200
                                ${selectedTag === null
                                    ? 'bg-blue-600 text-white shadow-md'
                                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                }
                            `}
                        >
                            すべて ({articles.length})
                        </button>

                        {/* 各タグのボタン */}
                        {tags.map(tag => {
                            const count = articles.filter(article =>
                                article.tags.includes(tag.name)
                            ).length

                            return (
                                <button
                                    key={tag.id}
                                    onClick={() => setSelectedTag(tag.name)}
                                    className={`
                                        px-4 py-2 rounded-full text-sm font-medium transition-all duration-200
                                        ${selectedTag === tag.name
                                            ? 'bg-blue-600 text-white shadow-md'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200 hover:shadow-sm'
                                        }
                                    `}
                                >
                                    {tag.name} ({count})
                                </button>
                            )
                        })}
                    </div>
                </div>
            )}

            {/* フィルター後の記事一覧 */}
            {filteredArticles.length === 0 ? (
                <div className="flex justify-center items-center min-h-[200px]">
                    <div className="text-center text-gray-500">
                        <p className="text-lg">「{selectedTag}」タグの記事が見つかりません</p>
                        <button
                            onClick={() => setSelectedTag(null)}
                            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
                        >
                            すべての記事を表示
                        </button>
                    </div>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {filteredArticles.map((article) => (
                        <ArticleCard 
                            key={article.id}
                            article={article}
                        />
                    ))}
                </div>
            )}
        </div>
    )
}
