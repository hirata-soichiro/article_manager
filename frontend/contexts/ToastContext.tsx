'use client'

import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
    id: string
    type: ToastType
    message: string
    duration?: number
}

interface ToastContextType {
    toasts: Toast[]
    showToast: (type: ToastType, message: string, duration?: number) => void
    hideToast: (id: string) => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

export function ToastProvider({ children }: { children: ReactNode }) {
    const [toasts, setToasts] = useState<Toast[]>([])

    const showToast = useCallback((type: ToastType, message: string, duration = 5000) => {
        const id = `${Date.now()}-${Math.random()}`
        const newToast: Toast = { id, type, message, duration }

        setToasts((prev) => [...prev, newToast])

        if (duration > 0) {
            setTimeout(() => {
                hideToast(id)
            }, duration)
        }
    }, [])

    const hideToast = useCallback((id: string) => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id))
    }, [])

    return (
        <ToastContext.Provider value={{ toasts, showToast, hideToast }}>
            {children}
        </ToastContext.Provider>
    )
}

export function useToast() {
    const context = useContext(ToastContext)
    if (!context) {
        throw new Error('useToast must be used within a ToastProvider')
    }
    return context
}
