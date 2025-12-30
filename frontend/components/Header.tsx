export default function Header() {
    return (
        <header className="bg-blue-600 text-white p-4 shadow-md">
            {/* bg-blue-600: 青色背景 */}
            {/* text-white: 白色テキスト */}
            {/* p-4: padding（内側の余白）4単位 */}
            {/* shadow-md: 中程度の影 */}

            <div className="container mx-auto flex justify-between items-center">
                {/* container: コンテナ幅に制限 */}
                {/* mx-auto: 左右中央揃え */}
                {/* flex: フレックスボックスレイアウト */}
                {/* justify-between: 左右に要素を配置 */}
                {/* items-center: 垂直方向中央揃え */}

                <h1 className="text-2xl font-bold">
                    {/* text-2xl: テキストサイズ大 */}
                    {/* font-bold: 太字 */}
                    記事管理システム
                </h1>

                <nav>
                    {/* ナビゲーション（今後拡張予定） */}
                </nav>
            </div>
        </header>
    )
}
