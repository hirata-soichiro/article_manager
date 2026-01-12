import Tag from './Tag'

interface TagListProps {
    tags: string[]
    onTagClick?: (tag: string) => void
}

export default function TagList({ tags, onTagClick }: TagListProps) {
    if (!tags || tags.length === 0) {
        return null
    }

    return (
        <div className="flex flex-wrap gap-2">
            {tags.map((tag, index) => (
                <Tag
                    key={`${tag}-${index}`}
                    name={tag}
                    // onTagClickが渡されている場合のみクリック処理を設定
                    onClick={onTagClick ? () => onTagClick(tag) : undefined}
                />
            ))}
        </div>
    )
}
