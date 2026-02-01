import Link from 'next/link'

export default function Header() {
    return (
        <header className="bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-500 text-white shadow-lg">
            <div className="container mx-auto px-6 py-4 flex justify-between items-center">
                <Link href="/articles" className="flex items-center gap-3 hover:opacity-90 transition-opacity cursor-pointer">
                    <div className="bg-white/10 backdrop-blur-sm p-2 rounded-lg">
                        <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/>
                        </svg>
                    </div>
                    <h1 className="text-3xl font-bold tracking-tight">
                        ArticleHub
                    </h1>
                </Link>

                <nav>
                    {/* ナビゲーション（今後拡張予定） */}
                </nav>
            </div>
        </header>
    )
}
