import Link from 'next/link'

export default function Sidebar() {
    return (
        <aside className="w-72 bg-gradient-to-b from-slate-50 to-slate-100 min-h-screen p-6 border-r border-slate-200">
            <nav>
                <h2 className="text-sm font-semibold mb-6 text-slate-500 uppercase tracking-wider">
                    Navigation
                </h2>

                <ul className="space-y-3">
                    <li>
                        <Link
                            href="/articles"
                            className="group flex items-center gap-4 p-4 rounded-xl hover:bg-white hover:shadow-md transition-all duration-200 border border-transparent hover:border-slate-200"
                        >
                            <div className="flex items-center justify-center w-12 h-12 rounded-lg bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-sm group-hover:shadow-md group-hover:scale-110 transition-all duration-200">
                                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                </svg>
                            </div>
                            <div className="flex flex-col">
                                <span className="font-semibold text-slate-800 group-hover:text-blue-600 transition-colors">記事検索</span>
                                <span className="text-xs text-slate-500">Search Articles</span>
                            </div>
                        </Link>
                    </li>

                    <li>
                        <Link
                            href="/articles/new"
                            className="group flex items-center gap-4 p-4 rounded-xl hover:bg-white hover:shadow-md transition-all duration-200 border border-transparent hover:border-slate-200"
                        >
                            <div className="flex items-center justify-center w-12 h-12 rounded-lg bg-gradient-to-br from-emerald-400 to-emerald-600 text-white shadow-sm group-hover:shadow-md group-hover:scale-110 transition-all duration-200">
                                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                                </svg>
                            </div>
                            <div className="flex flex-col">
                                <span className="font-semibold text-slate-800 group-hover:text-emerald-600 transition-colors">記事登録</span>
                                <span className="text-xs text-slate-500">New Article</span>
                            </div>
                        </Link>
                    </li>

                    <li>
                        <Link
                            href="/tags"
                            className="group flex items-center gap-4 p-4 rounded-xl hover:bg-white hover:shadow-md transition-all duration-200 border border-transparent hover:border-slate-200"
                        >
                            <div className="flex items-center justify-center w-12 h-12 rounded-lg bg-gradient-to-br from-purple-400 to-purple-600 text-white shadow-sm group-hover:shadow-md group-hover:scale-110 transition-all duration-200">
                                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                                </svg>
                            </div>
                            <div className="flex flex-col">
                                <span className="font-semibold text-slate-800 group-hover:text-purple-600 transition-colors">タグ管理</span>
                                <span className="text-xs text-slate-500">Manage Tags</span>
                            </div>
                        </Link>
                    </li>
                </ul>
            </nav>
        </aside>
    )
}
