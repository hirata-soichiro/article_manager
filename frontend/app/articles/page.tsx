import ArticleCard from "@/components/ArticleCard";
import { Article } from "@/types/article";

const mockArticles: Article[] = [
    {
        id: 1,
        title: 'Docker入門 - コンテナ技術の基礎',
        url: 'https://example.com/docker-basics',
        summary: 'Dockerの基本的な使い方を解説。コンテナとイメージの違い、docker-composeの活用方法など。',
        tags: ['Docker', '技術', 'インフラ'],
        memo: '後で見返す',
        createdAt: '2025-01-01T10:00:00Z',
        updatedAt: '2025-01-01T10:00:00Z',
    },
    {
        id: 2,
        title: 'Go言語で始めるWeb開発',
        url: 'https://example.com/golang-web',
        summary: 'Golangを使ったWebアプリケーション開発の入門。HTTPサーバーの作り方、ルーティング、ミドルウェアについて。',
        tags: ['Go', 'バックエンド', '技術'],
        createdAt: '2025-01-02T14:30:00Z',
        updatedAt: '2025-01-02T14:30:00Z',
    },
    {
        id: 3,
        title: 'Next.js App Routerの完全ガイド',
        url: 'https://example.com/nextjs-app-router',
        summary: 'Next.js 13以降の新機能App Routerの使い方。Server ComponentsとClient Componentsの違い。',
        tags: ['Next.js', 'React', 'フロントエンド'],
        memo: '重要',
        createdAt: '2025-01-03T09:15:00Z',
        updatedAt: '2025-01-03T09:15:00Z',
    },
]

export default function SearchPage() {
    return (
        <div>
            <h1 className="text-3xl font-bold mb-6 text-gray-800">
                記事一覧
            </h1>

            <p className="text-gray-600 mb-4">
                全 {mockArticles.length} 件の記事
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {/* grid: グリッドレイアウト */}
                {/* grid-cols-1: 1列（モバイル） */}
                {/* md:grid-cols-2: 中サイズ画面で2列 */}
                {/* lg:grid-cols-3: 大外面で3列 */}
                {/* gap-6: グリッド間の間隔6単位 */}

                {mockArticles.map((article) => (
                    <ArticleCard 
                        key={article.id} 
                        article={article} 
                    />
                ))}
            </div>
        </div>
    )
}
