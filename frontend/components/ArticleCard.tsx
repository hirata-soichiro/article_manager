import { Article } from '@/types/article'
import TagList from './TagList' 

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

            <div className="mb-3">
                <TagList tags={article.tags} />
            </div>

            {article.memo && (
                <p className="text-gray-500 text-sm italic mb-2">
                    メモ： {article.memo}
                </p>
            )}
        </div>
    )
}
