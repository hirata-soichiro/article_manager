import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Header from "@/components/Header";
import Sidebar from "@/components/Sidebar";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "記事管理システム",
  description: "記事を管理するWebアプリケーション",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      {/* lang="ja": 日本語ページであることを宣言 */}

      <body>
        <Header />

        <div className="flex">
          {/* flex: 横並びレイアウト */}

          <Sidebar />

          <main className="flex-1 p-8">
            {/* flex-1: 残りのスペースを全て使う */}
            {/* p-8: padding 8単位 */}
            {children}
            {/* ここに各ページの内容（page.tsx）が表示される */}
          </main>
        </div>
      </body>
    </html>
  );
}
