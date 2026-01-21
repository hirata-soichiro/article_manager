import Link from 'next/link'

export default function Sidebar() {
    return (
        <aside className="w-64 bg-gray-100 min-h-screen p-6">
            {/* aside: サイドバー用のHTML要素 */}
            {/* w-64: 幅を64単位（256px） */}
            {/* bg-gray-100: 薄いグレー背景 */}
            {/* min-h-screen: 最小高さを画面いっぱい */}
            {/* p-6: padding 6単位 */}

            <nav>
                <h2 className="text-lg font-semibold mb-4 text-gray-700">
                    {/* text-lg: テキストサイズ大 */}
                    {/* font-semibold: やや太字 */}
                    {/* mb-4: 下マージン4単位 */}
                    {/* text-gray-700: 濃いグレーテキスト */}
                    メニュー
                </h2>

                <ul className="space-y-2">
                    {/* space-y-2: 縦方向の要素間隔2単位 */}

                    <li>
                        <Link
                            href="/articles"
                            className="block p-2 rounded hover:bg-gray-200 transition"
                        >
                            {/* block: ブロック要素（幅いっぱい） */}
                            {/* rounded: 角を丸く */}
                            {/* hover:bg-gray-200: ホバー時に背景色変更 */}
                            {/* transition: 変化をスムーズに */}
                            🔍 記事検索
                        </Link>
                    </li>

                    <li>
                        <Link
                            href="/articles/new"
                            className="block p-2 rounded hover:bg-gray-200 transition"
                        >
                            ✏️ 記事登録
                        </Link>
                    </li>

                    <li>
                        <Link
                            href="/tags"
                            className="block p-2 rounded hover:bg-gray-200 transition"
                        >
                            🏷️ タグ管理
                        </Link>
                    </li>
                </ul>
            </nav>
        </aside>
    )
}
