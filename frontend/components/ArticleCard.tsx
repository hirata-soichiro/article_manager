'use client'

import { Article } from '@/types/article'
import TagList from './TagList'
import Link from 'next/link'

interface ArticleCardProps {
    article: Article
}

export default function ArticleCard({ article }: ArticleCardProps) {
    return (
        <div className="card hover:shadow-lg transition-shadow duration-300">
            <div className="flex justify-between items-start mb-2">
                <h3 className="text-xl font-bold text-gray-800 flex-1">
                    {article.title}
                </h3>

                {/* 編集ボタン */}
                <Link 
                    href={`/articles/${article.id}/edit`}
                    className="ml-2 px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors duration-200"
                >
                    編集
                </Link>
            </div>

            <a
                href={article.url}
                target="_blank"
                rel="noopener noreferrer"
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

            <div className="text-xs text-gray-400 mt-2">
                <p>作成: {new Date(article.createdAt).toLocaleDateString('ja-JP')}</p>
                <p>更新: {new Date(article.updatedAt).toLocaleDateString('ja-JP')}</p>
            </div>
        </div>
    )
}
