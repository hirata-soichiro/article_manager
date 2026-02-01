'use client'

import { Article } from '@/types/article'
import TagList from './TagList'
import Link from 'next/link'
import { useRouter } from 'next/navigation'

interface ArticleCardProps {
    article: Article
    onDelete?: (id: number) => void
}

export default function ArticleCard({ article, onDelete }: ArticleCardProps) {
    const router = useRouter()

    const handleCardClick = () => {
        router.push(`/articles/${article.id}`)
    }

    return (
        <div
            onClick={handleCardClick}
            className="bg-white rounded-lg shadow-md p-6 border border-gray-200 hover:shadow-xl hover:scale-[1.02] transition-all duration-300 relative cursor-pointer"
        >
            <h3 className="text-xl font-bold text-gray-800 mb-2">
                {article.title}
            </h3>

            <a
                href={article.url}
                target="_blank"
                rel="noopener noreferrer"
                onClick={(e) => e.stopPropagation()}
                className="text-blue-600 hover:underline text-sm mb-2 block break-all"
            >
                {article.url}
            </a>

            <p className="text-gray-600 mb-3 line-clamp-2">
                {article.summary}
            </p>

            <div className="mb-3">
                <TagList tags={article.tags} />
            </div>

            {article.memo && (
                <p className="text-gray-500 text-sm italic mb-2">
                    メモ： {article.memo}
                </p>
            )}

            {/* 編集・削除ボタン */}
            <div className="absolute bottom-4 right-4 flex gap-2">
                {/* 編集ボタン */}
                <Link
                    href={`/articles/${article.id}/edit`}
                    onClick={(e) => e.stopPropagation()}
                    className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
                    aria-label="記事を編集"
                >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                </Link>

                {/* 削除ボタン */}
                {onDelete && (
                    <button
                        onClick={(e) => {
                            e.stopPropagation()
                            onDelete(article.id)
                        }}
                        className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
                        aria-label="記事を削除"
                    >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                    </button>
                )}
            </div>
        </div>
    )
}
