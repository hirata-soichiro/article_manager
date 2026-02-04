import dynamic from 'next/dynamic'

// ArticleFormを動的インポート（記事作成ページでのみロード）
const ArticleForm = dynamic(() => import('@/components/ArticleForm'), {
    loading: () => (
        <div className="flex justify-center items-center min-h-[400px]">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
    ),
})

export default function NewArticlePage() {
    return <ArticleForm />
}
