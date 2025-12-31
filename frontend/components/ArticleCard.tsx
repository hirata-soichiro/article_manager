import { Article } from '@/types/article'

interface ArticleCardProps {
    article: Article
}

export default function ArticleCard({ article }: ArticleCardProps) {
    return (
        <div className="card hover:shadow-lg cursor-pointer">
            {/* card: globals.cssで定義した共通スタイル */}
            {/* hover:shadow-lg: ホバー時に大きな影 */}
            {/* cursor-pointer: マウスカーソルをポインターに */}

            <h3 className="text-xl font-bold mb-2 text-gray-800">
                {/* text-xl: テキストサイズ大 */}
                {/* font-bold: 太字 */}
                {/* mb-2: 下マージン2単位 */}
                {/* text-gray-800: 濃いグレーテキスト */}
                {article.title}
            </h3>

            <a
                href={article.url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:underline text-sm mb-2 block break-all"
            >
                {/* target="_blank": 新しいタブで開く */}
                {/* rel="noopener noreferrer": セキュリティ対策 */}
                {/* hover:underline: ホバー時に下線 */}
                {/* break-all: 長いURLを折り返す */}
                {article.url}
            </a>

            <p className="text-gray-600 mb-3 line-clamp-2">
                {/* line-clamp-2: 2行まで表示、それ以上は... */}
                {article.summary}
            </p>

            <div className="flex flex-wrap gap-2 mb-3">
                {/* flex: 横並び */}
                {/* flex-wrap: 幅が足りない時に折り返す */}
                {/* gap-2: 要素間の間隔2単位 */}

                {article.tags.map((tag) => (
                    <span
                        key={tag}
                        className="px-2 py-1 bg-blue-100 text-blue-700 rounded-full text-xs"
                    >
                        {tag}
                    </span>
                ))}
            </div>

            {article.memo && (
                <p className="text-gray-500 text-sm italic mb-2">
                    メモ： {article.memo}
                </p>
            )}
        </div>
    )
}
