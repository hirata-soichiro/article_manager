import { describe, it, expect } from 'vitest'
import { formatDate, validateUrl } from './sample'

describe('formatDate', () => {
    it('should format date to YYYY-MM-DD', () => {
        const date = new Date('2025-01-15T10:30:00Z')
        expect(formatDate(date)).toBe('2025-01-15')
    })
})

describe('validateUrl', () => {
    it('should return true for valid URL', () => {
        expect(validateUrl('https://example.com')).toBe(true)
    })

    it('should return false for invalid URL', () => {
        expect(validateUrl('not-a-url')).toBe(false)
    })

    it('should return false for empty string', () => {
        expect(validateUrl('')).toBe(false)
    })
})
