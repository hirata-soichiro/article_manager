import dynamic from 'next/dynamic'

// ArticleEditFormを動的インポート（記事編集ページでのみロード）
const ArticleEditForm = dynamic(() => import('@/components/ArticleEditForm'), {
    loading: () => (
        <div className="flex justify-center items-center min-h-[400px]">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
    ),
})

export default function ArticleEditPage() {
    return (
        <div className="container mx-auto px-4 py-8">
            <ArticleEditForm />
        </div>
    )
}
