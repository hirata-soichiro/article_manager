interface TagProps {
    name: string
    onClick?: () => void
}

export default function Tag({ name, onClick }: TagProps) {
    // タグ名から一貫した色を生成（同じ名前なら同じ色）
    const getColorFromTagName = (tagName: string): string => {
        const colorSets = [
            { bg: 'bg-blue-100', text: 'text-blue-700' },
            { bg: 'bg-green-100', text: 'text-green-700' },
            { bg: 'bg-purple-100', text: 'text-purple-700' },
            { bg: 'bg-pink-100', text: 'text-pink-700' },
            { bg: 'bg-yellow-100', text: 'text-yellow-700' },
            { bg: 'bg-indigo-100', text: 'text-indigo-700' },
            { bg: 'bg-red-100', text: 'text-red-700' },
            { bg: 'bg-teal-100', text: 'text-teal-700' },
        ]

        let hash = 0
        for (let i = 0; i < tagName.length; i++) {
            hash = tagName.charCodeAt(i) + hash
        }

        const colorIndex = hash % colorSets.length
        const colors = colorSets[colorIndex]
        return `${colors.bg} ${colors.text}`
    }

    const colorClasses = getColorFromTagName(name)

    return (
        <span
            onClick={onClick}
            className={`
                px-4 py-2
                ${colorClasses}
                rounded-full
                text-xs
                font-medium
                inline-flex items-center
                transition-all duration-200
                ${onClick ? 'cursor-pointer hover:scale-105 hover:shadow-md' : ''}
            `}
        >
            {name}
        </span>
    )
}
