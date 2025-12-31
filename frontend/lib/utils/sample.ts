export function formatDate(date: Date): string {
    return date.toISOString().split('T')[0]
}

export function validateUrl(url: string): boolean {
    try {
        new URL(url)
        return true
    } catch {
        return false
    }
}
