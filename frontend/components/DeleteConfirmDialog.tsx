'use client'

interface DeleteConfirmDialogProps {
    isOpen: boolean
    title: string
    message: string
    onConfirm: () => void
    onCancel: () => void
    isDeleting?: boolean
}

export default function DeleteConfirmDialog({
    isOpen,
    title,
    message,
    onConfirm,
    onCancel,
    isDeleting = false,
}: DeleteConfirmDialogProps) {
    if (!isOpen) return null

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            {/* オーバーレイ（背景暗転） */}
            <div
                className="absolute inset-0 bg-black bg-opacity-50 transition-opacity"
                onClick={onCancel}
                aria-hidden="true"
            />

            {/* ダイアログコンテンツ */}
            <div className="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6 z-10">
                {/* アイコン */}
                <div className="flex items-center justify-center w-12 h-12 mx-auto mb-4 bg-red-100 rounded-full">
                    <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                </div>

                {/* タイトル */}
                <h3 className="text-lg font-bold text-gray-900 text-center mb-2">
                    {title}
                </h3>

                {/* メッセージ */}
                <p className="text-sm text-gray-600 text-center mb-6">
                    {message}
                </p>

                {/* ボタン */}
                <div className="flex gap-3">
                    <button
                        type="button"
                        onClick={onCancel}
                        disabled={isDeleting}
                        className="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 disabled:bg-gray-100 disabled:cursor-not-allowed transition font-medium"
                    >
                        キャンセル
                    </button>
                    <button
                        type="button"
                        onClick={onConfirm}
                        disabled={isDeleting}
                        className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-red-400 disabled:cursor-not-allowed transition font-medium"
                    >
                        {isDeleting ? '削除中...' : '削除'}
                    </button>
                </div>
            </div>
        </div>
    )
}
