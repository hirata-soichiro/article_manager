'use client'

import { useBookRecommendations } from '@/hooks/useBookRecommendations'
import { memo, useCallback } from 'react'
import type { BookRecommendation } from '@/types/book'

// 書籍カードコンポーネント
interface BookCardProps {
    book: BookRecommendation
}

const BookCard = memo(function BookCard({ book }: BookCardProps) {
    // Amazon購入リンクを優先、なければ楽天
    const purchaseLink = book.purchaseLinks.amazon || book.purchaseLinks.rakuten
    const purchaseLinkLabel = book.purchaseLinks.amazon ? 'Amazonで購入' : '楽天で購入'

    const handlePurchaseClick = useCallback((e: React.MouseEvent) => {
        // カード全体のクリックイベントを停止しない（外部リンクなので）
    }, [])

    return (
        <div
            className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 hover:shadow-md hover:scale-[1.02] transition-all duration-200 flex flex-col"
        >
            {/* タイトル */}
            <h3 className="text-base font-bold text-gray-800 mb-2 line-clamp-2 min-h-[3rem]">
                {book.title}
            </h3>

            {/* 著者名 */}
            <p className="text-sm text-gray-600 mb-3 line-clamp-1">
                {book.author}
            </p>

            {/* 購入ボタン */}
            {purchaseLink ? (
                <a
                    href={purchaseLink}
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={handlePurchaseClick}
                    className="mt-auto w-full text-center px-4 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white text-sm font-medium rounded-lg hover:shadow-lg hover:scale-105 transition-all duration-200 flex items-center justify-center gap-2"
                    aria-label={`${book.title}を購入`}
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    <span>{purchaseLinkLabel}</span>
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                </a>
            ) : (
                <div className="mt-auto w-full text-center px-4 py-2 bg-gray-200 text-gray-500 text-sm font-medium rounded-lg cursor-not-allowed">
                    購入リンクなし
                </div>
            )}
        </div>
    )
})

// スケルトンローディング用のカードコンポーネント
const SkeletonBookCard = memo(function SkeletonBookCard() {
    return (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 flex flex-col animate-pulse">
            {/* タイトル部分 */}
            <div className="h-6 bg-gray-200 rounded mb-2"></div>
            <div className="h-6 bg-gray-200 rounded mb-2 w-3/4"></div>

            {/* 著者名部分 */}
            <div className="h-4 bg-gray-200 rounded mb-3 w-1/2"></div>

            {/* ボタン部分 */}
            <div className="mt-auto h-10 bg-gray-200 rounded"></div>
        </div>
    )
})

// メインコンポーネント
export default function BookRecommendations() {
    const { books, loading, error } = useBookRecommendations()

    // ローディング中
    if (loading) {
        return (
            <section
                className="mb-8 bg-white rounded-lg shadow-sm border border-gray-200 p-6"
                aria-label="書籍推薦を読み込み中"
            >
                {/* セクションタイトル */}
                <div className="flex items-center gap-2 mb-4">
                    <svg className="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                    </svg>
                    <h2 className="text-lg font-semibold text-gray-800">
                        あなたにおすすめの書籍
                    </h2>
                </div>

                {/* スケルトンUI（5個） */}
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                    {[...Array(5)].map((_, index) => (
                        <SkeletonBookCard key={index} />
                    ))}
                </div>
            </section>
        )
    }

    // エラー時または書籍が0件の場合は静かに失敗（セクション非表示）
    if (error || books.length === 0) {
        return null
    }

    // 書籍推薦を表示
    return (
        <section
            className="mb-8 bg-white rounded-lg shadow-sm border border-gray-200 p-6"
            aria-label="あなたにおすすめの書籍"
        >
            {/* セクションタイトル */}
            <div className="flex items-center gap-2 mb-4">
                <svg className="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
                <h2 className="text-lg font-semibold text-gray-800">
                    あなたにおすすめの書籍
                </h2>
            </div>

            {/* 書籍カード一覧 */}
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                {books.map((book, index) => (
                    <BookCard key={`${book.title}-${index}`} book={book} />
                ))}
            </div>
        </section>
    )
}