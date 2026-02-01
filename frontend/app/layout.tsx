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
  title: "ArticleHub",
  description: "あなたの記事を整理・管理するモダンなプラットフォーム",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      {/* lang="ja": 日本語ページであることを宣言 */}

      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <Header />

        <div className="flex">
          <Sidebar />

          <main className="flex-1 p-8 bg-gradient-to-br from-slate-50 via-white to-slate-50 min-h-screen">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
