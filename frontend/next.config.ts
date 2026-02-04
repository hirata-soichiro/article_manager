import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // React Strict Modeを有効化（開発時のバグ検出）
  reactStrictMode: true,

  // コンパイラ最適化
  compiler: {
    // 本番環境でconsole.logを削除
    removeConsole: process.env.NODE_ENV === 'production' ? {
      exclude: ['error', 'warn'],
    } : false,
  },

  // 画像最適化設定
  images: {
    formats: ['image/avif', 'image/webp'], // 次世代フォーマットを優先
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    minimumCacheTTL: 60 * 60 * 24 * 365, // 1年間キャッシュ
  },

  // Turbopack設定（Next.js 16+）
  turbopack: {},

  // 実験的機能
  experimental: {
    // 最適化されたパッケージインポート
    optimizePackageImports: ['@/components', '@/lib'],
  },

  // パフォーマンス予算
  onDemandEntries: {
    // 開発時のページキャッシュ時間
    maxInactiveAge: 60 * 1000, // 60秒
    pagesBufferLength: 5, // 同時に保持するページ数
  },
};

export default nextConfig;
