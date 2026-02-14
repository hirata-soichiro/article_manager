# フロントエンド実装ガイド - ArticleHub（新卒1年目向け）

このドキュメントは、プログラミングを始めて間もない方、特に新卒1年目のエンジニアが、このフロントエンドアプリケーションの仕組みを**本当に理解できる**ことを目的としています。

「このコードは何をしているか」だけでなく、「なぜこのように書かれているのか」「どのような問題を解決しているのか」を詳しく説明します。

---

## 目次

1. [このアプリケーションは何をするものか](#このアプリケーションは何をするものか)
2. [全体の仕組みを理解する](#全体の仕組みを理解する)
3. [使用技術の基礎知識](#使用技術の基礎知識)
4. [実装を段階的に理解する](#実装を段階的に理解する)
5. [実際のコード例で学ぶ](#実際のコード例で学ぶ)
6. [実際の複雑なコード - ArticleFormの全体像を理解する](#例6-実際の複雑なコード---articleformの全体像を理解する)
7. [高度なReact Hooks - useCallback と useRef の実践](#例7-高度なreact-hooks---usecallback-と-useref-の実践)
8. [コンポーネント設計のガイドライン](#コンポーネント設計のガイドライン)
9. [テストの書き方 - 品質を保証する](#テストの書き方---品質を保証する)
10. [トラブルシューティングガイド](#トラブルシューティングガイド)
11. [実践的な開発フロー - 新機能を追加する手順](#実践的な開発フロー---新機能を追加する手順)
12. [まとめ: 学んだことを振り返る](#まとめ-学んだことを振り返る)

---

## このアプリケーションは何をするものか

### 機能概要

ArticleHubは、Webの記事を管理するアプリケーションです。ユーザーは以下のことができます：

1. **記事を追加する** - URLを入力すると、AIが自動でタイトルや要約、タグを生成
2. **記事を検索する** - キーワードやタグで記事を探せる
3. **記事を編集・削除する** - 後から情報を更新できる
4. **記事を一覧表示する** - カード形式で見やすく表示

### なぜこのアプリが必要なのか？

技術記事やブログを読んでいると、「あとで読みたい」「この記事は重要だ」と思うことがあります。ブラウザのブックマークだけでは管理が大変です。このアプリを使えば：

- 記事にタグをつけて分類できる
- 要約を見て内容を思い出せる
- キーワードで素早く検索できる

**公式ドキュメント参考:** なし（これはサンプルアプリケーションです）

---

## 全体の仕組みを理解する

### Webアプリケーションの基本構造

Webアプリケーションは、大きく分けて2つの部分で構成されています：

```
┌─────────────────────────┐
│  フロントエンド           │  ← 今回説明する部分
│  (ブラウザで動く)         │     ユーザーが見る画面
│  - Next.js + React       │
└───────────┬─────────────┘
            │ HTTP通信
            │ (データのやり取り)
┌───────────▼─────────────┐
│  バックエンド             │
│  (サーバーで動く)         │     データベースと連携
│  - Go言語                │     ビジネスロジック
└─────────────────────────┘
```

**フロントエンドの役割:**
- ユーザーに画面を見せる（HTML/CSS）
- ユーザーの操作を受け付ける（クリック、入力など）
- バックエンドにデータを要求する（APIリクエスト）
- 受け取ったデータを画面に表示する

**参考:** [MDN - What is a web server?](https://developer.mozilla.org/en-US/docs/Learn/Common_questions/Web_mechanics/What_is_a_web_server)

---

## 使用技術の基礎知識

### 1. React - UIを作る基本の仕組み

**Reactとは何か:**
Reactは、Facebookが作った「画面を作るためのライブラリ」です。

**従来のWeb開発の問題:**
昔のWeb開発では、HTMLを直接操作（DOM操作）していました。例えば：

```javascript
// 昔のやり方（jQuery時代）
document.getElementById('counter').innerHTML = '0';
document.getElementById('button').addEventListener('click', function() {
  const counter = document.getElementById('counter');
  const currentValue = parseInt(counter.innerHTML);
  counter.innerHTML = currentValue + 1;
});
```

問題点：
- HTMLの要素を探す処理が多い
- 状態（カウンターの値）がどこにあるか分かりにくい
- コードが複雑になると管理が大変

**Reactの解決策:**
Reactでは「状態（データ）」を管理し、画面は自動的に更新されます。

```javascript
// Reactのやり方
function Counter() {
  const [count, setCount] = useState(0);  // 状態を管理

  return (
    <div>
      <p>{count}</p>
      <button onClick={() => setCount(count + 1)}>増やす</button>
    </div>
  );
}
```

**Reactの3つの重要な概念:**

#### 1-1. コンポーネント（部品）

画面を「部品」に分解して作ります。例えば：

```
記事一覧ページ
├── ヘッダー（Header）
├── 検索バー（SearchBar）
├── 記事カード（ArticleCard）× 複数
└── フッター（Footer）
```

各部品（コンポーネント）は独立していて、再利用できます。

#### 1-2. 状態（State）

「今どんなデータを持っているか」を管理します。

```javascript
const [articles, setArticles] = useState([]);  // 記事一覧を保持
const [loading, setLoading] = useState(true);   // ローディング中かどうか
```

状態が変わると、Reactが自動的に画面を更新してくれます。

#### 1-3. Props（プロパティ）

親コンポーネントから子コンポーネントにデータを渡す仕組み。

```javascript
// 親コンポーネント
<ArticleCard article={article} onDelete={handleDelete} />

// 子コンポーネント
function ArticleCard({ article, onDelete }) {
  // articleとonDeleteを使える
}
```

**参考:** [React - Quick Start](https://react.dev/learn)

---

### 2. Next.js - Reactをもっと便利にする

**Next.jsとは何か:**
Reactだけだと設定が大変です。Next.jsは「Reactでアプリを作りやすくするフレームワーク」です。

**Next.jsが提供する便利な機能:**

#### 2-1. ファイルベースのルーティング

URLとファイルの配置が対応します。

```
app/
├── page.tsx              → http://localhost:3000/
├── articles/
│   └── page.tsx          → http://localhost:3000/articles
│   └── [id]/
│       └── page.tsx      → http://localhost:3000/articles/123
```

#### 2-2. App Router（新しいルーティングシステム）

Next.js 13から導入された新しい仕組み。特殊なファイル名に意味があります：

- `layout.tsx` - 全ページ共通のレイアウト
- `page.tsx` - 各ページの内容
- `loading.tsx` - ローディング画面
- `error.tsx` - エラー画面

#### 2-3. サーバーコンポーネントとクライアントコンポーネント

**サーバーコンポーネント（デフォルト）:**
サーバー側で実行され、HTMLとして送られます。高速です。

**クライアントコンポーネント（`'use client'`を書く）:**
ブラウザで実行されます。インタラクティブな操作（クリック、入力など）が必要な場合に使います。

```typescript
'use client'  // これを書くとクライアントコンポーネントになる

import { useState } from 'react'

export default function Counter() {
  const [count, setCount] = useState(0)
  // ...
}
```

**参考:** [Next.js - Getting Started](https://nextjs.org/docs)

---

### 3. TypeScript - 型で安全にする

**TypeScriptとは何か:**
JavaScriptに「型」を追加した言語です。

**JavaScriptの問題:**

```javascript
// JavaScriptの例
function addNumbers(a, b) {
  return a + b;
}

addNumbers(1, 2);      // 3 (正常)
addNumbers("1", "2");  // "12" (バグ！文字列連結になる)
```

**TypeScriptの解決策:**

```typescript
// TypeScriptの例
function addNumbers(a: number, b: number): number {
  return a + b;
}

addNumbers(1, 2);      // 3 (正常)
addNumbers("1", "2");  // エラー！コンパイル時に検出
```

**型定義の例:**

```typescript
// 記事の型を定義
interface Article {
  id: number        // 数値
  title: string     // 文字列
  url: string       // 文字列
  tags: string[]    // 文字列の配列
  createdAt: string // 文字列（日付）
}

// この型を使う
function displayArticle(article: Article) {
  console.log(article.title);  // OK
  console.log(article.xxx);     // エラー！xxxというプロパティは存在しない
}
```

**なぜ型が重要なのか:**
1. **バグを早期発見** - コードを書いている時点でエラーが分かる
2. **IDE補完** - 入力候補が表示される
3. **リファクタリングが安全** - 変更の影響範囲が分かる

**参考:** [TypeScript - Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)

---

### 4. Tailwind CSS - スタイリング

**Tailwind CSSとは何か:**
CSSを書く代わりに、クラス名でスタイルを指定する方法です。

**従来のCSS:**

```html
<button class="my-button">クリック</button>

<style>
.my-button {
  background-color: blue;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
}
</style>
```

**Tailwind CSS:**

```html
<button class="bg-blue-600 text-white px-4 py-2 rounded">クリック</button>
```

- `bg-blue-600` - 背景色を青に
- `text-white` - 文字色を白に
- `px-4` - 左右のpadding
- `py-2` - 上下のpadding
- `rounded` - 角を丸く

#### なぜTailwind CSSを使うのか？従来のCSSとの比較

**問題1: 命名の苦労**

従来のCSSでは、クラス名を考えるのに時間がかかります。

```html
<!-- 従来のCSS -->
<div class="article-card">
  <h3 class="article-card__title">タイトル</h3>
</div>

<style>
.article-card { /* クラス名を考える必要がある */ }
.article-card__title { /* 命名規則を守る必要がある */ }
</style>
```

Tailwind CSSでは、**クラス名を考える必要がありません**。

```html
<!-- Tailwind CSS -->
<div class="bg-white rounded-lg p-6">
  <h3 class="text-xl font-bold">タイトル</h3>
</div>
```

**問題2: グローバルスコープの競合**

従来のCSSは、すべてグローバルスコープです。

```css
/* components/ArticleCard.css */
.title { font-size: 20px; }

/* components/UserCard.css */
.title { font-size: 24px; }  /* 競合する！ */
```

Tailwind CSSでは、**スコープの問題がありません**。

```tsx
<h3 className="text-xl">ArticleCardのタイトル</h3>
<h3 className="text-2xl">UserCardのタイトル</h3>
```

**問題3: 未使用のCSSが残る**

従来のCSSでは、削除した機能のCSSが残り続けます。

```css
.old-button { /* もう使っていないが削除していない */ }
.deprecated-card { /* いつ追加したか忘れた */ }
```

Tailwind CSSでは、**ビルド時に使われているクラスだけを抽出**します。

```bash
# 実際に使われているクラスのみが含まれる
npm run build
```

**5つの主要なメリット:**

1. **命名不要**: クラス名を考える時間がゼロ
2. **スコープ安全**: グローバルスコープの競合がない
3. **自動最適化**: 未使用CSSが自動削除される
4. **レスポンシブ簡単**: `md:w-1/2 lg:w-1/3`で画面サイズ対応
5. **デザイン統一**: 決められた値（色、スペーシング）のみ使える

**「クラス名が長い」問題の解決:**

確かにクラス名は長くなりますが、**コンポーネント化で解決**します。

```tsx
// ❌ 悪い例：毎回長いクラスを書く
{articles.map(article => (
  <div className="bg-white rounded-lg shadow-md p-6 border border-gray-200 hover:shadow-xl">
    {/* 内容 */}
  </div>
))}

// ✅ 良い例：コンポーネント化
function ArticleCard({ article }) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 border border-gray-200 hover:shadow-xl">
      {/* 内容 */}
    </div>
  )
}

// 使う側はシンプル
{articles.map(article => <ArticleCard article={article} />)}
```

**実際のプロジェクトでの例:**

このプロジェクトの`ArticleCard.tsx`を見てみましょう。

```tsx
// components/ArticleCard.tsx (実際のコード)
<article className="bg-white rounded-lg shadow-md p-6 border border-gray-200 hover:shadow-xl hover:scale-[1.02] transition-all duration-300">
  <h3 className="text-xl font-bold text-gray-800 mb-2">
    {article.title}
  </h3>
  <p className="text-gray-600 mb-3">
    {article.summary}
  </p>
</article>
```

これを従来のCSSで書くと、CSSファイルとの行き来が必要で、命名も考える必要があります。Tailwindなら**HTMLだけで完結**します。

**レスポンシブデザインの例:**

```tsx
// モバイル: 全幅、タブレット: 50%、PC: 33%
<div className="w-full md:w-1/2 lg:w-1/3">

// モバイル: 縦並び、PC: 横並び
<div className="flex flex-col lg:flex-row">

// モバイル: 隠す、PC: 表示
<div className="hidden lg:block">
```

`md:`（768px以上）、`lg:`（1024px以上）などのプレフィックスで、**1行でレスポンシブ対応**できます。

**結論:**

Tailwind CSSは、**見た目は複雑**ですが、**実際の開発は圧倒的にシンプル**です。

- 従来のCSS: 命名 + スコープ管理 + 未使用CSS + ファイル間の移動
- Tailwind CSS: コンポーネント内で完結

このプロジェクトのように、Reactと組み合わせれば、コンポーネント化で長いクラス名の問題も解決できます。

**参考:** [Tailwind CSS - Utility-First Fundamentals](https://tailwindcss.com/docs/utility-first)

---

## 実装を段階的に理解する

### ステップ1: アプリケーションの起動

#### 1-1. エントリーポイント（最初に実行される場所）

ユーザーがブラウザでアクセスすると、Next.jsは自動的に `app/layout.tsx` を読み込みます。

**なぜlayout.tsxから始まるのか？**
Next.jsのルールで、`layout.tsx`は全てのページの「外枠」になります。ヘッダーやサイドバーなど、全ページ共通の部分をここに書きます。

**`app/layout.tsx` の役割を理解する:**

```typescript
// app/layout.tsx
import type { Metadata } from "next";
import Header from "@/components/Header";
import Sidebar from "@/components/Sidebar";
import { ToastProvider } from "@/contexts/ToastContext";

// メタデータ（<head>タグに入る情報）
export const metadata: Metadata = {
  title: "ArticleHub",
  description: "記事管理アプリ",
};

export default function RootLayout({
  children,  // 各ページのコンテンツが入る
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>
        <ToastProvider>
          <Header />
          <div className="flex">
            <Sidebar />
            <main>
              {children}  {/* ここに各ページの内容が入る */}
            </main>
          </div>
        </ToastProvider>
      </body>
    </html>
  );
}
```

**ポイント解説:**

1. **`{children}`とは？**
   - 各ページのコンテンツが自動的にここに入ります
   - 例: `/articles`にアクセスすると、`app/articles/page.tsx`の内容が`{children}`に入る

2. **なぜToastProviderで囲むのか？**
   - アプリ全体で通知機能を使えるようにするため
   - 後で詳しく説明します

3. **構造を図で表すと:**
   ```
   ┌──────────────────────────────┐
   │ Header（常に表示）             │
   ├────────┬─────────────────────┤
   │Sidebar │ Main（ページ内容）   │
   │（常に  │ ← {children}が入る  │
   │ 表示） │                      │
   └────────┴─────────────────────┘
   ```

---

### ステップ2: データの流れを理解する（3層アーキテクチャ）

このアプリケーションは、コードを3つの層に分けて整理しています。

```
┌─────────────────────────────────────┐
│  Layer 1: プレゼンテーション層        │
│  場所: app/, components/             │
│  役割: 画面を表示、ユーザー操作を受付  │
│  例: ボタンをクリックしたら何かする    │
└─────────────┬───────────────────────┘
              │
              │ 「記事一覧を取得して！」
              ▼
┌─────────────────────────────────────┐
│  Layer 2: ビジネスロジック層          │
│  場所: hooks/                        │
│  役割: データの取得、状態管理         │
│  例: useArticles, useArticleSearch   │
└─────────────┬───────────────────────┘
              │
              │ 「GET /api/articles を実行」
              ▼
┌─────────────────────────────────────┐
│  Layer 3: データアクセス層            │
│  場所: lib/api/                      │
│  役割: APIとの通信                   │
│  例: articleClient.getAll()          │
└─────────────────────────────────────┘
```

**なぜ3つに分けるのか？**

1. **責任が明確になる**
   - 画面の処理はLayer 1
   - データ管理はLayer 2
   - API通信はLayer 3

2. **テストしやすい**
   - 各層を独立してテストできる

3. **変更に強い**
   - 例: APIのURLが変わっても、Layer 3だけ修正すればいい

---

### ステップ3: 記事一覧を表示する仕組み

実際のコードで、「記事一覧を表示する」処理を追っていきます。

#### 3-1. ユーザーが `/articles` にアクセス

ブラウザで `/articles` にアクセスすると、`app/articles/page.tsx` が実行されます。

```typescript
// app/articles/page.tsx
'use client'  // ← クライアントコンポーネント（ブラウザで動く）

import { useArticles } from '@/hooks/useArticles'
import ArticleCard from '@/components/ArticleCard'

export default function ArticlesPage() {
  // Layer 2（ビジネスロジック層）を呼び出す
  const { articles, loading, error } = useArticles()

  // ローディング中の表示
  if (loading) {
    return <div>読み込み中...</div>
  }

  // エラーの表示
  if (error) {
    return <div>エラー: {error.message}</div>
  }

  // 記事を表示
  return (
    <div>
      {articles.map((article) => (
        <ArticleCard key={article.id} article={article} />
      ))}
    </div>
  )
}
```

**コードの流れ:**

1. **`useArticles()`を呼ぶ**
   - これがLayer 2（ビジネスロジック層）
   - 記事データを取得してくれる

2. **状態で表示を切り替え**
   - `loading`が`true`なら「読み込み中」
   - `error`があればエラー表示
   - 成功したら記事を表示

3. **`map()`で繰り返し表示**
   - 配列の各要素に対して処理を実行
   - 各記事を`ArticleCard`コンポーネントで表示

**重要: `key`属性とは？**
```typescript
<ArticleCard key={article.id} article={article} />
```

Reactでリストを表示する時は、必ず`key`をつけます。

- **なぜ必要？** Reactが各要素を識別して、効率的に更新するため
- **何を指定する？** 一意の値（この場合は`article.id`）

**参考:** [React - Rendering Lists](https://react.dev/learn/rendering-lists)

---

#### 3-2. Layer 2: カスタムフック `useArticles` の仕組み

次に、`useArticles()`の中身を見ていきます。

```typescript
// hooks/useArticles.ts
import { useState, useEffect } from 'react'
import { articleClient } from '@/lib/api/articleClient'

export function useArticles() {
  // 状態を定義
  const [articles, setArticles] = useState([])    // 記事一覧
  const [loading, setLoading] = useState(true)    // ローディング中か
  const [error, setError] = useState(null)        // エラー

  // コンポーネントがマウントされた時に実行される
  useEffect(() => {
    // 非同期関数を定義
    async function fetchData() {
      try {
        setLoading(true)
        // Layer 3（データアクセス層）を呼ぶ
        const data = await articleClient.getAll()
        setArticles(data)  // 取得したデータを状態に保存
      } catch (err) {
        setError(err)      // エラーを状態に保存
      } finally {
        setLoading(false)  // ローディング終了
      }
    }

    fetchData()  // 関数を実行
  }, [])  // 空配列 = マウント時に1回だけ実行

  // 状態を返す
  return { articles, loading, error }
}
```

**詳細解説:**

##### useState - 状態管理

```typescript
const [articles, setArticles] = useState([])
```

- `useState`は「状態を作る」Hook
- `articles`が現在の値
- `setArticles`が値を更新する関数
- `[]`が初期値（空配列）

**なぜ状態が必要？**
JavaScriptの普通の変数だと、値を変えても画面が更新されません。

```typescript
// ダメな例
let articles = []
articles = newData  // 値は変わるが画面は更新されない

// 正しい例
const [articles, setArticles] = useState([])
setArticles(newData)  // 値が変わり、画面も自動更新される
```

##### useEffect - 副作用（データ取得）

```typescript
useEffect(() => {
  // ここに処理を書く
}, [])  // 依存配列
```

- `useEffect`は「副作用」を実行するHook
- 副作用とは: データ取得、イベント登録、タイマーなど
- 依存配列が`[]`だと、マウント時（最初に1回だけ）実行される

**なぜuseEffectを使うのか？**
コンポーネントの中で直接`async/await`を使えないからです。

```typescript
// これはできない
export default async function Component() {
  const data = await fetchData()  // エラー！
  return <div>{data}</div>
}

// useEffectを使う
export default function Component() {
  const [data, setData] = useState(null)

  useEffect(() => {
    async function fetch() {
      const result = await fetchData()
      setData(result)
    }
    fetch()
  }, [])

  return <div>{data}</div>
}
```

**参考:**
- [React - useState](https://react.dev/reference/react/useState)
- [React - useEffect](https://react.dev/reference/react/useEffect)

---

#### 3-3. Layer 3: APIクライアント `articleClient` の仕組み

最後に、実際にAPIを呼び出す部分を見ます。

```typescript
// lib/api/articleClient.ts
class ArticleClient {
  async getAll() {
    // fetchでHTTPリクエストを送る
    const response = await fetch('http://localhost:8080/api/articles')

    // レスポンスをJSONに変換
    const data = await response.json()

    // snake_case → camelCase に変換して返す
    return data.map(item => ({
      id: item.id,
      title: item.title,
      url: item.url,
      summary: item.summary,
      tags: item.tags ?? [],
      memo: item.memo ?? '',
      createdAt: item.created_at,    // snake_case → camelCase
      updatedAt: item.updated_at,    // snake_case → camelCase
    }))
  }
}

export const articleClient = new ArticleClient()
```

**詳細解説:**

##### fetch API - HTTPリクエスト

```typescript
const response = await fetch('http://localhost:8080/api/articles')
```

- `fetch`はブラウザ標準のAPI
- HTTPリクエストを送る
- `await`で結果が返ってくるまで待つ

**HTTPリクエストとは？**
```
ブラウザ → 「記事一覧をください」 → サーバー
         ← 「これが記事一覧です」  ←
```

**参考:** [MDN - Using Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch)

##### async/await - 非同期処理

```typescript
async function getAll() {
  const response = await fetch(url)
  const data = await response.json()
  return data
}
```

- `async`をつけると非同期関数になる
- `await`で結果を待つ
- `await`は`async`関数の中でのみ使える

**なぜawaitが必要？**
ネットワーク通信は時間がかかります。待たないと、データが取得できる前に次の処理が実行されてしまいます。

```typescript
// awaitなし（間違い）
const response = fetch(url)      // Promiseオブジェクトが返る
const data = response.json()     // エラー！responseはまだ取得できていない

// awaitあり（正しい）
const response = await fetch(url)     // 結果が返ってくるまで待つ
const data = await response.json()    // JSONに変換
```

**参考:** [MDN - async function](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function)

##### Nullish Coalescing Operator (??)

```typescript
tags: item.tags ?? []
```

- `??`は「左側がnullまたはundefinedなら、右側を使う」演算子
- APIから`null`が返ってきても、空配列`[]`にする

**他の方法との違い:**

```typescript
// || を使うと問題がある
const value = 0 || 10        // 10 (0はfalsyなので10になる)
const value = "" || "abc"    // "abc" (空文字もfalsyなので"abc"になる)

// ?? を使うと正しく動く
const value = 0 ?? 10        // 0 (0はnullでもundefinedでもない)
const value = "" ?? "abc"    // "" (空文字もnullでもundefinedでもない)
const value = null ?? 10     // 10 (nullなので10になる)
```

---

## 実際のコード例で学ぶ

### 例1: 記事作成フォームの実装

記事作成は、このアプリで最も複雑な処理の1つです。段階的に理解していきましょう。

#### 問題: フォームの値をどう管理するか？

**HTML/JavaScriptの場合:**

```html
<input type="text" id="title" />
<button onclick="submit()">送信</button>

<script>
function submit() {
  const title = document.getElementById('title').value
  // titleを使って何かする
}
</script>
```

問題点:
- DOMを直接操作している
- 値を取得するたびに`getElementById`を呼ぶ
- 複雑になると管理が大変

**Reactの解決策: Controlled Component（制御されたコンポーネント）**

```typescript
function ArticleForm() {
  // 状態でフォームの値を管理
  const [title, setTitle] = useState('')
  const [url, setUrl] = useState('')
  const [summary, setSummary] = useState('')

  // 送信処理
  const handleSubmit = (e) => {
    e.preventDefault()  // デフォルトのフォーム送信を防ぐ

    // ここで値を使える
    console.log({ title, url, summary })
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={title}                          // 状態から値を取得
        onChange={(e) => setTitle(e.target.value)}  // 変更時に状態を更新
      />
      <input
        type="text"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
      />
      <textarea
        value={summary}
        onChange={(e) => setSummary(e.target.value)}
      />
      <button type="submit">送信</button>
    </form>
  )
}
```

**ポイント:**

1. **`value={title}`**
   - inputの値を状態から取得
   - 常に「状態が真実」（Single Source of Truth）

2. **`onChange={(e) => setTitle(e.target.value)}`**
   - ユーザーが入力したら、状態を更新
   - `e.target.value`が入力された値

3. **`e.preventDefault()`**
   - フォームのデフォルト動作（ページリロード）を防ぐ
   - SPA（Single Page Application）では必須

**参考:** [React - Forms](https://react.dev/learn/managing-state#reacting-to-input-with-state)

---

#### AI自動生成機能の実装

ユーザーがURLを入力して「AI自動生成」ボタンを押すと、バックエンドのAI APIを呼び出します。

```typescript
function ArticleForm() {
  const [url, setUrl] = useState('')
  const [title, setTitle] = useState('')
  const [summary, setSummary] = useState('')
  const [isGenerating, setIsGenerating] = useState(false)
  const [generateError, setGenerateError] = useState(null)

  const handleAIGenerate = async () => {
    // バリデーション: URLが入力されているかチェック
    if (!url.trim()) {
      setGenerateError('URLを入力してください')
      return
    }

    try {
      setIsGenerating(true)       // ローディング開始
      setGenerateError(null)       // エラーをクリア

      // APIを呼び出す
      const data = await articleClient.generate(url)

      // 取得したデータをフォームに自動入力
      setTitle(data.title)
      setSummary(data.summary)

    } catch (err) {
      setGenerateError(err.message)
    } finally {
      setIsGenerating(false)      // ローディング終了
    }
  }

  return (
    <form>
      {/* URL入力 */}
      <input
        type="text"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="https://example.com"
      />

      {/* AI生成ボタン */}
      <button
        type="button"
        onClick={handleAIGenerate}
        disabled={!url.trim() || isGenerating}  // URLが空、または生成中は無効
      >
        {isGenerating ? '生成中...' : 'AI自動生成'}
      </button>

      {/* エラー表示 */}
      {generateError && <p className="error">{generateError}</p>}

      {/* タイトル（AI生成後に自動入力される） */}
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="タイトル"
      />

      {/* 要約（AI生成後に自動入力される） */}
      <textarea
        value={summary}
        onChange={(e) => setSummary(e.target.value)}
        placeholder="要約"
      />
    </form>
  )
}
```

**ポイント解説:**

##### 1. try-catch-finally パターン

```typescript
try {
  // 成功する可能性がある処理
  const data = await articleClient.generate(url)
} catch (err) {
  // エラーが起きたときの処理
  setGenerateError(err.message)
} finally {
  // 成功してもエラーでも必ず実行される処理
  setIsGenerating(false)
}
```

- `try`: 試したい処理を書く
- `catch`: エラーが起きたときの処理
- `finally`: 必ず実行される処理（クリーンアップ）

**なぜfinallyが必要？**
ローディング状態を必ず`false`に戻すため。

```typescript
// finallyなし（悪い例）
try {
  setIsGenerating(true)
  const data = await fetch(url)
  setIsGenerating(false)  // 成功時だけfalseになる
} catch (err) {
  setIsGenerating(false)  // エラー時にもfalseにする必要がある
}

// finallyあり（良い例）
try {
  setIsGenerating(true)
  const data = await fetch(url)
} catch (err) {
  // エラー処理
} finally {
  setIsGenerating(false)  // 必ず実行される
}
```

**参考:** [MDN - try...catch](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/try...catch)

##### 2. 条件付きレンダリング

```typescript
{isGenerating ? '生成中...' : 'AI自動生成'}
```

- 三項演算子: `条件 ? 真の場合 : 偽の場合`
- `isGenerating`が`true`なら「生成中...」、`false`なら「AI自動生成」

**他の書き方:**

```typescript
// && を使う（真の場合のみ表示）
{generateError && <p className="error">{generateError}</p>}

// if文を使う（JSXの外）
let buttonText
if (isGenerating) {
  buttonText = '生成中...'
} else {
  buttonText = 'AI自動生成'
}
return <button>{buttonText}</button>
```

##### 3. disabled属性

```typescript
disabled={!url.trim() || isGenerating}
```

- ボタンを無効化する条件:
  - `!url.trim()`: URLが空（空白のみも含む）
  - `||`: または
  - `isGenerating`: 生成中

**なぜ無効化するのか？**
- URLがないのにボタンを押せるとエラーになる
- 生成中に何度も押すと、重複リクエストが発生する

---

### 例2: Context API - グローバルな通知機能

アプリ全体で「成功しました」「エラーが発生しました」といった通知を表示したい時、どうすればいいでしょうか？

**問題: Propsのバケツリレー**

通知機能を各コンポーネントで使うには、Propsで渡していく必要があります。

```typescript
// App
<Layout showToast={showToast}>
  <ArticlesPage showToast={showToast}>
    <ArticleCard showToast={showToast} />
  </ArticlesPage>
</Layout>
```

問題点:
- 全てのコンポーネントにPropsを渡す必要がある
- 深い階層だと非常に面倒
- 途中のコンポーネントは使わないのに受け渡すだけ

**解決策: Context API**

Context APIを使うと、アプリ全体でデータを共有できます。

#### Context の作成

```typescript
// contexts/ToastContext.tsx
import { createContext, useContext, useState } from 'react'

// 1. Contextを作成
const ToastContext = createContext(undefined)

// 2. Providerコンポーネントを作成
export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([])

  // トーストを表示する関数
  const showToast = (type, message) => {
    const id = Date.now()  // 一意のID
    const newToast = { id, type, message }

    setToasts((prev) => [...prev, newToast])

    // 5秒後に自動で消す
    setTimeout(() => {
      hideToast(id)
    }, 5000)
  }

  // トーストを消す関数
  const hideToast = (id) => {
    setToasts((prev) => prev.filter(toast => toast.id !== id))
  }

  // Contextの値を提供
  return (
    <ToastContext.Provider value={{ toasts, showToast, hideToast }}>
      {children}
    </ToastContext.Provider>
  )
}

// 3. カスタムフックを作成（使いやすくするため）
export function useToast() {
  const context = useContext(ToastContext)

  if (!context) {
    throw new Error('useToast must be used within a ToastProvider')
  }

  return context
}
```

#### Contextの使用

```typescript
// app/layout.tsx（アプリ全体をラップ）
export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        <ToastProvider>
          {children}
        </ToastProvider>
      </body>
    </html>
  )
}

// どのコンポーネントからでも使える
function ArticleForm() {
  const { showToast } = useToast()  // Contextから取得

  const handleSubmit = async () => {
    try {
      await articleClient.create(data)
      showToast('success', '記事を作成しました')  // 通知を表示
    } catch (err) {
      showToast('error', 'エラーが発生しました')
    }
  }

  return <form onSubmit={handleSubmit}>...</form>
}
```

**ポイント:**

1. **createContext** - Contextを作る
2. **Provider** - 値を提供するコンポーネント
3. **useContext** - 値を取得するHook

**なぜカスタムフックを作るのか？**

```typescript
// カスタムフックなし
const context = useContext(ToastContext)
if (!context) {
  throw new Error('...')
}
const { showToast } = context

// カスタムフックあり
const { showToast } = useToast()  // 簡潔！
```

**参考:** [React - Context](https://react.dev/learn/passing-data-deeply-with-context)

---

### 例3: 楽観的UI更新 - ユーザー体験を向上させる

**問題: APIレスポンスを待つと遅い**

記事を削除する時、通常はこうなります：

```
ユーザーがクリック
  ↓
APIリクエスト送信
  ↓
レスポンスを待つ（1-2秒）
  ↓
画面から記事が消える
```

これだと、ユーザーは「本当に削除されたのか？」と不安になります。

**解決策: 楽観的UI更新（Optimistic UI Update）**

成功すると仮定して、先に画面を更新します。

```
ユーザーがクリック
  ↓
画面から記事が即座に消える
  ↓
バックグラウンドでAPIリクエスト
  ↓
（失敗したら元に戻す）
```

#### 実装例

```typescript
function useArticles() {
  const [articles, setArticles] = useState([])

  const deleteArticle = async (id) => {
    try {
      // 1. 現在の記事リストをバックアップ
      const previousArticles = articles

      // 2. 画面から即座に削除（楽観的更新）
      setArticles((prev) => prev.filter(article => article.id !== id))

      try {
        // 3. APIリクエストを送る
        await articleClient.delete(id)

      } catch (err) {
        // 4. 失敗したら元に戻す（ロールバック）
        setArticles(previousArticles)
        throw err
      }

    } catch (err) {
      console.error('削除に失敗しました', err)
    }
  }

  return { articles, deleteArticle }
}
```

**ポイント:**

1. **バックアップを取る**
   ```typescript
   const previousArticles = articles
   ```
   失敗時に元に戻すため

2. **即座に画面を更新**
   ```typescript
   setArticles((prev) => prev.filter(article => article.id !== id))
   ```
   ユーザーには瞬時に反映される

3. **失敗したらロールバック**
   ```typescript
   setArticles(previousArticles)
   ```
   元の状態に戻す

**メリット:**
- レスポンスが速く感じる
- ユーザー体験が向上

**デメリット:**
- 実装が複雑になる
- 失敗時の処理が必要

---

### 例4: キャッシング - 不要なリクエストを減らす

**問題: 同じデータを何度も取得する**

記事一覧ページに戻るたびに、毎回APIリクエストを送っていたら無駄です。

**解決策: データをキャッシュする**

一度取得したデータを保存して、次回は再利用します。

```typescript
// モジュールスコープにキャッシュを作成
const articlesCache = {
  data: null,           // キャッシュデータ
  timestamp: 0,         // 取得した時刻
  ttl: 60 * 1000,      // 有効期間（60秒）
}

function useArticles() {
  const [articles, setArticles] = useState([])
  const [loading, setLoading] = useState(true)

  const fetchArticles = async (useCache = true) => {
    // キャッシュが有効かチェック
    if (useCache && articlesCache.data) {
      const age = Date.now() - articlesCache.timestamp

      if (age < articlesCache.ttl) {
        // キャッシュが新しい（60秒以内）
        console.log('キャッシュから取得')
        setArticles(articlesCache.data)
        setLoading(false)
        return  // APIリクエストをスキップ
      }
    }

    // キャッシュがない、または古い
    console.log('APIから取得')
    setLoading(true)
    const data = await articleClient.getAll()

    // キャッシュに保存
    articlesCache.data = data
    articlesCache.timestamp = Date.now()

    setArticles(data)
    setLoading(false)
  }

  useEffect(() => {
    fetchArticles()
  }, [])

  return { articles, loading }
}
```

**ポイント:**

1. **モジュールスコープのキャッシュ**
   ```typescript
   const articlesCache = { ... }
   ```
   関数の外に定義することで、複数のコンポーネントで共有される

2. **TTL（Time To Live）**
   ```typescript
   ttl: 60 * 1000  // 60秒 = 60,000ミリ秒
   ```
   キャッシュの有効期間を設定

3. **年齢チェック**
   ```typescript
   const age = Date.now() - articlesCache.timestamp
   ```
   キャッシュが何秒前に取得されたかを計算

**キャッシュの無効化:**

記事を作成・更新・削除したら、キャッシュを無効化します。

```typescript
const createArticle = async (input) => {
  await articleClient.create(input)
  articlesCache.data = null  // キャッシュをクリア
}
```

---

### 例5: useMemo と useCallback - パフォーマンス最適化

**問題: 不要な再計算・再レンダリング**

Reactは状態が変わると、コンポーネント全体を再実行します。

```typescript
function ArticlesPage() {
  const [articles, setArticles] = useState([])
  const [keyword, setKeyword] = useState('')

  // 毎回計算される
  const filteredArticles = articles.filter(article =>
    article.title.includes(keyword)
  )

  return <div>...</div>
}
```

`keyword`が変わっていなくても、`articles`が変わると`filter`が再実行されます。

**解決策1: useMemo - 計算結果をキャッシュ**

```typescript
function ArticlesPage() {
  const [articles, setArticles] = useState([])
  const [keyword, setKeyword] = useState('')

  // articlesまたはkeywordが変わった時だけ再計算
  const filteredArticles = useMemo(() => {
    return articles.filter(article => article.title.includes(keyword))
  }, [articles, keyword])  // 依存配列

  return <div>...</div>
}
```

**useMemoの仕組み:**

1. 初回実行時、計算結果をキャッシュ
2. 再レンダリング時、依存配列をチェック
3. 依存配列の値が変わっていなければ、キャッシュを返す
4. 変わっていたら、再計算してキャッシュを更新

**解決策2: useCallback - 関数をキャッシュ**

関数も毎回新しく作られます。

```typescript
function ArticlesPage() {
  const [articles, setArticles] = useState([])

  // 毎回新しい関数が作られる
  const handleSearch = async (keyword) => {
    const results = await articleClient.search(keyword)
    setArticles(results)
  }

  return <SearchBar onSearch={handleSearch} />
}
```

子コンポーネント`SearchBar`は、`onSearch`が変わるたびに再レンダリングされます。

```typescript
function ArticlesPage() {
  const [articles, setArticles] = useState([])

  // 関数をメモ化（依存配列が変わらない限り同じ関数）
  const handleSearch = useCallback(async (keyword) => {
    const results = await articleClient.search(keyword)
    setArticles(results)
  }, [])  // 依存配列が空なので、一度だけ作られる

  return <SearchBar onSearch={handleSearch} />
}
```

**いつ使うべきか？**

- **useMemo**: 計算コストが高い処理（大量のデータのフィルタリングなど）
- **useCallback**: 子コンポーネントに関数を渡す時

**注意: 過度な最適化は不要**

全てをメモ化する必要はありません。パフォーマンス問題が実際に発生してから最適化しましょう。

**参考:**
- [React - useMemo](https://react.dev/reference/react/useMemo)
- [React - useCallback](https://react.dev/reference/react/useCallback)

---

## 例6: 実際の複雑なコード - ArticleFormの全体像を理解する

ここまでで基礎を学んできました。次は、**実際のプロジェクトのコード**を見てみましょう。

### 簡単な例から実際のコードへの段階的な理解

#### レベル1: 最もシンプルな形（50行程度）

```typescript
'use client'

import { useState } from 'react'

export default function SimpleArticleForm() {
  const [title, setTitle] = useState('')
  const [url, setUrl] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    console.log({ title, url })
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="タイトル"
      />
      <input
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="URL"
      />
      <button type="submit">送信</button>
    </form>
  )
}
```

**特徴:**
- 状態は2つだけ
- バリデーションなし
- エラーハンドリングなし
- 最小限の機能

#### レベル2: 中級レベル（150行程度）

```typescript
'use client'

import { useState, FormEvent } from 'react'
import { articleClient } from '@/lib/api/articleClient'
import { useRouter } from 'next/navigation'

export default function IntermediateArticleForm() {
  const router = useRouter()

  // 状態管理
  const [title, setTitle] = useState('')
  const [url, setUrl] = useState('')
  const [summary, setSummary] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // バリデーション
  const [titleError, setTitleError] = useState<string | null>(null)
  const [urlError, setUrlError] = useState<string | null>(null)

  const validateTitle = (value: string) => {
    if (!value.trim()) {
      setTitleError('タイトルは必須です')
      return false
    }
    setTitleError(null)
    return true
  }

  const validateUrl = (value: string) => {
    if (!value.trim()) {
      setUrlError('URLは必須です')
      return false
    }
    try {
      new URL(value)
      setUrlError(null)
      return true
    } catch {
      setUrlError('正しいURL形式で入力してください')
      return false
    }
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()

    // バリデーション
    const isTitleValid = validateTitle(title)
    const isUrlValid = validateUrl(url)

    if (!isTitleValid || !isUrlValid) {
      return
    }

    try {
      setIsSubmitting(true)
      setError(null)

      await articleClient.create({
        title: title.trim(),
        url: url.trim(),
        summary: summary.trim(),
        tags: [],
        memo: '',
      })

      router.push('/articles')
    } catch (err) {
      setError('記事の作成に失敗しました')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">記事登録</h1>

      <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow-md p-6">
        {/* エラーメッセージ */}
        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
            {error}
          </div>
        )}

        {/* タイトル */}
        <div className="mb-6">
          <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
            タイトル <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            onBlur={() => validateTitle(title)}
            className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
              titleError
                ? 'border-red-500 focus:ring-red-500'
                : 'border-gray-300 focus:ring-blue-500'
            }`}
            placeholder="記事のタイトルを入力"
          />
          {titleError && <p className="mt-1 text-sm text-red-500">{titleError}</p>}
        </div>

        {/* URL */}
        <div className="mb-6">
          <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-2">
            URL <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            id="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            onBlur={() => validateUrl(url)}
            className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
              urlError
                ? 'border-red-500 focus:ring-red-500'
                : 'border-gray-300 focus:ring-blue-500'
            }`}
            placeholder="https://example.com"
          />
          {urlError && <p className="mt-1 text-sm text-red-500">{urlError}</p>}
        </div>

        {/* 要約 */}
        <div className="mb-6">
          <label htmlFor="summary" className="block text-sm font-medium text-gray-700 mb-2">
            要約 <span className="text-red-500">*</span>
          </label>
          <textarea
            id="summary"
            value={summary}
            onChange={(e) => setSummary(e.target.value)}
            rows={4}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="記事の要約を入力"
          />
        </div>

        {/* 送信ボタン */}
        <div className="flex gap-4">
          <button
            type="submit"
            disabled={isSubmitting || !title.trim() || !url.trim() || !summary.trim()}
            className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition font-medium"
          >
            {isSubmitting ? '登録中...' : '登録'}
          </button>
          <button
            type="button"
            onClick={() => router.push('/articles')}
            disabled={isSubmitting}
            className="px-6 py-3 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition font-medium"
          >
            キャンセル
          </button>
        </div>
      </form>
    </div>
  )
}
```

**追加された機能:**
- バリデーション（onBlurで実行）
- エラーハンドリング（try-catch-finally）
- ローディング状態管理
- スタイリング（Tailwind CSS）
- ボタンの無効化ロジック

**新しい概念:**
- **`onBlur`イベント**: フィールドからフォーカスが外れた時に実行
- **条件付きクラス**: エラー時にスタイルを変更
- **`disabled`属性**: 条件に応じてボタンを無効化

#### レベル3: 実際のプロジェクトレベル（615行）

実際の`ArticleForm.tsx`では、さらに以下が追加されています：

**追加機能:**
1. **AI自動生成** (187-245行目)
   - URLからタイトル・要約・タグを自動生成
   - 独自のローディング状態とエラー処理

2. **タグ管理** (40-98行目)
   - タグ一覧の取得
   - タグ検索機能
   - 新規タグの作成
   - 複雑な表示/非表示ロジック

3. **複雑な状態管理** (15-37行目)
   - 14個の状態変数
   - 各状態の相互作用

4. **詳細なバリデーション** (123-171行目)
   - リアルタイムバリデーション
   - 複数の入力チェック

**実際のコードで学べること:**

```typescript
// components/ArticleForm.tsx の重要部分

// 1. 複数の状態をグループ化して管理
const [title, setTitle] = useState('')
const [url, setUrl] = useState('')
const [summary, setSummary] = useState('')
const [memo, setMemo] = useState('')
const [selectedTags, setSelectedTags] = useState<string[]>([])

// タグ一覧の状態
const [tags, setTags] = useState<Tag[]>([])
const [tagsLoading, setTagsLoading] = useState(true)
const [tagsError, setTagsError] = useState<string | null>(null)

// フォーム送信の状態
const [isSubmitting, setIsSubmitting] = useState(false)
const [formError, setFormError] = useState<string | null>(null)

// AI自動生成の状態
const [isGenerating, setIsGenerating] = useState(false)
const [generateError, setGenerateError] = useState<string | null>(null)
```

**状態管理の設計パターン:**
- 関連する状態をグループ化（タグ関連、フォーム送信関連など）
- それぞれにローディングとエラーの状態を持つ
- 状態名から役割が分かる命名

**複雑な処理の例: AI自動生成**

```typescript
const handleAIGenerate = async () => {
  // 1. 事前チェック
  if (!url.trim()) {
    setGenerateError('URLを入力してください')
    return
  }

  const isUrlValid = validateUrl(url)
  if (!isUrlValid) {
    setGenerateError('正しいURL形式で入力してください')
    return
  }

  try {
    setIsGenerating(true)
    setGenerateError(null)

    // 2. API呼び出し
    const data = await articleClient.generate(url.trim(), memo.trim() || undefined)

    // 3. 取得したデータをフォームに反映
    setTitle(data.title)
    setSummary(data.summary)

    // 4. 複雑なタグ処理
    if (data.tags && Array.isArray(data.tags)) {
      const existingTagNames = tags.map(tag => tag.name)
      const newTags = data.tags.filter(tagName => !existingTagNames.includes(tagName))

      // 新規タグをリストに追加
      if (newTags.length > 0) {
        const newTagObjects: Tag[] = newTags.map(tagName => ({
          id: 0,
          name: tagName,
          createdAt: '',
          updatedAt: ''
        }))
        setTags([...tags, ...newTagObjects])
        setGeneratedTags(newTags)
      }

      // 全てのタグを選択状態に
      setSelectedTags(data.tags)
    }

    // 5. バリデーションエラーをクリア
    setTitleError(null)
    setSummaryError(null)
  } catch (err) {
    if (err instanceof Error) {
      setGenerateError(err.message)
    } else {
      setGenerateError('AI生成中にエラーが発生しました')
    }
  } finally {
    setIsGenerating(false)
  }
}
```

**この処理から学べること:**
1. **複数の状態を連携させる**: 1つの処理で複数の状態を更新
2. **エラーハンドリングの詳細化**: 各段階でエラーをチェック
3. **データ変換**: APIレスポンスをUIの状態に変換
4. **配列操作**: `filter`、`map`を使った複雑なデータ処理

---

## 例7: 高度なReact Hooks - useCallback と useRef の実践

基本的なHooks（useState、useEffect）は学びました。次は、**実際のプロジェクトで頻繁に使われる高度なHooks**を理解しましょう。

### useCallback - 関数のメモ化

**問題: 関数は毎回新しく作られる**

```typescript
function ArticleCard({ article }) {
  // この関数は、ArticleCardが再レンダリングされるたびに新しく作られる
  const handleClick = () => {
    console.log(article.id)
  }

  return <button onClick={handleClick}>詳細を見る</button>
}
```

コンポーネントが再レンダリングされると、関数も毎回新しく作られます。これにより：
- 子コンポーネントが不必要に再レンダリングされる
- パフォーマンスが低下する（特に大量のリストがある場合）

**解決策: useCallback**

```typescript
import { useCallback } from 'react'

function ArticleCard({ article }) {
  // useCallbackで関数をメモ化
  // 依存配列[article.id]が変わらない限り、同じ関数インスタンスを返す
  const handleClick = useCallback(() => {
    console.log(article.id)
  }, [article.id])  // article.idが変わった時だけ関数を再生成

  return <button onClick={handleClick}>詳細を見る</button>
}
```

**実際のプロジェクトでの使用例: ArticleCard.tsx**

```typescript
// components/ArticleCard.tsx (実際のコード)
'use client'

import { memo, useCallback } from 'react'
import { useRouter } from 'next/navigation'

interface ArticleCardProps {
  article: Article
  onDelete?: (id: number) => void
}

// React.memoでpropsが変わらない限り再レンダリングを防ぐ
const ArticleCard = memo(function ArticleCard({ article, onDelete }: ArticleCardProps) {
  const router = useRouter()

  // useCallbackでハンドラをメモ化
  const handleCardClick = useCallback(() => {
    router.push(`/articles/${article.id}`)
  }, [router, article.id])

  const handleDeleteClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation()  // カードのクリックイベントを止める
    if (onDelete) {
      onDelete(article.id)
    }
  }, [article.id, onDelete])

  const handleEditClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation()
  }, [])

  return (
    <article onClick={handleCardClick}>
      <h3>{article.title}</h3>
      <button onClick={handleEditClick}>編集</button>
      <button onClick={handleDeleteClick}>削除</button>
    </article>
  )
})

export default ArticleCard
```

**ポイント:**

1. **`memo`との組み合わせ**
   ```typescript
   const ArticleCard = memo(function ArticleCard({ article, onDelete }) {
     // ...
   })
   ```
   - `memo`は「propsが変わらない限り再レンダリングしない」
   - `useCallback`は「関数が変わらない」ことを保証
   - 両方使うことで、最大限の最適化

2. **依存配列の重要性**
   ```typescript
   useCallback(() => {
     router.push(`/articles/${article.id}`)
   }, [router, article.id])  // routerとarticle.idを依存に含める
   ```
   - 関数内で使う値は全て依存配列に含める
   - 含めないと古い値を参照してしまう（バグの原因）

3. **イベント伝播の制御**
   ```typescript
   const handleDeleteClick = useCallback((e: React.MouseEvent) => {
     e.stopPropagation()  // 親要素（カード）のクリックイベントを止める
     if (onDelete) {
       onDelete(article.id)
     }
   }, [article.id, onDelete])
   ```

**いつuseCallbackを使うべきか？**

✅ **使うべき場合:**
- 子コンポーネントに関数を渡す時
- 依存配列に関数を含める時（useEffect、useMemoなど）
- React.memoと組み合わせる時

❌ **使わなくていい場合:**
- コンポーネント内だけで使う関数
- パフォーマンス問題が発生していない時
- 単純な計算の関数

### useRef - 再レンダリングを引き起こさない値の保持

**問題: 再レンダリングを避けたい**

```typescript
function Timer() {
  const [count, setCount] = useState(0)
  let intervalId  // これはダメ！再レンダリングで消える

  const start = () => {
    intervalId = setInterval(() => {
      setCount(c => c + 1)
    }, 1000)
  }

  const stop = () => {
    clearInterval(intervalId)  // intervalIdが消えているので止まらない
  }

  return (
    <div>
      <p>{count}</p>
      <button onClick={start}>開始</button>
      <button onClick={stop}>停止</button>
    </div>
  )
}
```

**解決策: useRef**

```typescript
import { useState, useRef } from 'react'

function Timer() {
  const [count, setCount] = useState(0)
  const intervalIdRef = useRef<number | null>(null)  // refで保持

  const start = () => {
    intervalIdRef.current = setInterval(() => {
      setCount(c => c + 1)
    }, 1000)
  }

  const stop = () => {
    if (intervalIdRef.current !== null) {
      clearInterval(intervalIdRef.current)
      intervalIdRef.current = null
    }
  }

  return (
    <div>
      <p>{count}</p>
      <button onClick={start}>開始</button>
      <button onClick={stop}>停止</button>
    </div>
  )
}
```

**useRefの特徴:**
1. **再レンダリングを引き起こさない**
   - 値を変更しても画面は更新されない
   - stateと違って、設定した値がすぐに使える

2. **レンダリングをまたいで値を保持**
   - コンポーネントが再レンダリングされても値が残る
   - `.current`プロパティで値にアクセス

**実際のプロジェクトでの使用例: useArticles.ts**

```typescript
// hooks/useArticles.ts (実際のコード)
export function useArticles(): UseArticlesReturn {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<ApiError | Error | null>(null)

  // 重複リクエスト防止のフラグ
  const fetchingRef = useRef(false)

  const fetchArticles = useCallback(async (useCache = true) => {
    // 重複リクエスト防止
    if (fetchingRef.current) {
      return  // 既に取得中なら何もしない
    }

    try {
      fetchingRef.current = true  // フラグをtrueに
      setLoading(true)
      setError(null)
      const data = await articleClient.getAll()
      setArticles(data)
    } catch (err) {
      setError(err as Error)
    } finally {
      setLoading(false)
      fetchingRef.current = false  // フラグをfalseに
    }
  }, [])

  useEffect(() => {
    fetchArticles()
  }, [fetchArticles])

  return { articles, loading, error }
}
```

**なぜuseRefを使うのか？**

1. **stateを使うとどうなるか**
   ```typescript
   const [fetching, setFetching] = useState(false)

   if (fetching) {
     return  // 判定時点ではfalse、設定が反映されるのは次のレンダリング
   }
   setFetching(true)  // 設定してもすぐには反映されない
   ```
   - stateは次のレンダリングで反映される
   - 連続で呼ばれると重複リクエストが発生する

2. **useRefを使うと**
   ```typescript
   if (fetchingRef.current) {
     return  // すぐに判定できる
   }
   fetchingRef.current = true  // すぐに反映される
   ```
   - 値の変更がすぐに反映される
   - 重複リクエストを確実に防げる

**useRefの他の使い道:**

1. **DOM要素への参照**
   ```typescript
   function InputForm() {
     const inputRef = useRef<HTMLInputElement>(null)

     const focusInput = () => {
       inputRef.current?.focus()  // 直接DOM操作
     }

     return (
       <div>
         <input ref={inputRef} type="text" />
         <button onClick={focusInput}>フォーカス</button>
       </div>
     )
   }
   ```

2. **前回の値を保持**
   ```typescript
   function Counter() {
     const [count, setCount] = useState(0)
     const prevCountRef = useRef(0)

     useEffect(() => {
       prevCountRef.current = count  // 前回の値を保存
     }, [count])

     return (
       <div>
         <p>現在: {count}</p>
         <p>前回: {prevCountRef.current}</p>
         <button onClick={() => setCount(count + 1)}>増やす</button>
       </div>
     )
   }
   ```

**useRefとuseStateの使い分け:**

| 用途 | useState | useRef |
|------|----------|--------|
| 画面に表示する値 | ✅ | ❌ |
| 再レンダリングが必要 | ✅ | ❌ |
| 値の変更をすぐ反映 | ❌ | ✅ |
| DOM要素への参照 | ❌ | ✅ |
| タイマーIDなどの保持 | ❌ | ✅ |

**参考:**
- [React - useCallback](https://react.dev/reference/react/useCallback)
- [React - useRef](https://react.dev/reference/react/useRef)
- [React - memo](https://react.dev/reference/react/memo)

---

## コンポーネント設計のガイドライン

実際の開発では「どうコードを分割するか」が重要です。ここでは、実践的な設計の考え方を学びます。

### いつコンポーネントを分割すべきか？

**分割すべき3つのサイン:**

#### 1. コードが長くなりすぎた（200行以上）

**悪い例:**
```typescript
// ArticlePage.tsx (600行)
function ArticlePage() {
  // 記事表示のロジック (100行)
  // 検索機能のロジック (100行)
  // フィルター機能のロジック (100行)
  // タグ管理のロジック (100行)
  // ページネーションのロジック (100行)

  return (
    <div>
      {/* 全部の UI が1つのreturnに詰まっている (200行) */}
    </div>
  )
}
```

**良い例:**
```typescript
// ArticlePage.tsx (100行)
function ArticlePage() {
  return (
    <div>
      <SearchBar />
      <FilterPanel />
      <ArticleList />
      <Pagination />
    </div>
  )
}

// SearchBar.tsx (50行)
function SearchBar() {
  // 検索機能のロジックのみ
}

// ArticleList.tsx (80行)
function ArticleList() {
  // 記事表示のロジックのみ
}
```

**判断基準:**
- 1つのファイルが200行を超えたら分割を検討
- 50-150行が理想的

#### 2. 同じUIパターンを繰り返している

**悪い例:**
```typescript
function ArticlesPage() {
  return (
    <div>
      {articles.map((article) => (
        <div key={article.id} className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-xl font-bold">{article.title}</h3>
          <p className="text-gray-600">{article.summary}</p>
          <div className="flex gap-2">
            {article.tags.map((tag) => (
              <span key={tag} className="px-2 py-1 bg-blue-100 rounded">
                {tag}
              </span>
            ))}
          </div>
          <button onClick={() => deleteArticle(article.id)}>削除</button>
        </div>
      ))}
    </div>
  )
}
```

**良い例:**
```typescript
// ArticleCard.tsx
function ArticleCard({ article, onDelete }: ArticleCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-xl font-bold">{article.title}</h3>
      <p className="text-gray-600">{article.summary}</p>
      <TagList tags={article.tags} />
      <button onClick={() => onDelete(article.id)}>削除</button>
    </div>
  )
}

// ArticlesPage.tsx
function ArticlesPage() {
  return (
    <div>
      {articles.map((article) => (
        <ArticleCard
          key={article.id}
          article={article}
          onDelete={deleteArticle}
        />
      ))}
    </div>
  )
}
```

**判断基準:**
- 同じUIを3回以上書いたらコンポーネント化
- mapで繰り返す部分は必ずコンポーネント化

#### 3. 1つのコンポーネントが複数の責任を持っている

**悪い例:**
```typescript
function ArticleForm() {
  // フォームの状態管理
  const [title, setTitle] = useState('')
  const [url, setUrl] = useState('')
  const [tags, setTags] = useState<Tag[]>([])
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const [tagSearch, setTagSearch] = useState('')

  // タグフィルタリングロジック (50行)
  const filteredTags = tags.filter(/* ... */)

  // タグ選択ロジック (30行)
  const toggleTag = (tagName: string) => {/* ... */}

  // フォーム送信ロジック (40行)
  const handleSubmit = async () => {/* ... */}

  return (
    <form onSubmit={handleSubmit}>
      {/* フォームフィールド */}
      <input value={title} onChange={(e) => setTitle(e.target.value)} />

      {/* タグ選択UI (100行) */}
      <div>
        <input
          value={tagSearch}
          onChange={(e) => setTagSearch(e.target.value)}
        />
        {filteredTags.map(/* ... */)}
      </div>

      <button type="submit">送信</button>
    </form>
  )
}
```

**良い例:**
```typescript
// ArticleForm.tsx
function ArticleForm() {
  const [title, setTitle] = useState('')
  const [url, setUrl] = useState('')
  const [selectedTags, setSelectedTags] = useState<string[]>([])

  const handleSubmit = async () => {
    // フォーム送信ロジックのみ
  }

  return (
    <form onSubmit={handleSubmit}>
      <input value={title} onChange={(e) => setTitle(e.target.value)} />
      <input value={url} onChange={(e) => setUrl(e.target.value)} />

      {/* タグ選択機能を別コンポーネントに */}
      <TagSelector
        selectedTags={selectedTags}
        onTagsChange={setSelectedTags}
      />

      <button type="submit">送信</button>
    </form>
  )
}

// TagSelector.tsx
interface TagSelectorProps {
  selectedTags: string[]
  onTagsChange: (tags: string[]) => void
}

function TagSelector({ selectedTags, onTagsChange }: TagSelectorProps) {
  const [tags, setTags] = useState<Tag[]>([])
  const [tagSearch, setTagSearch] = useState('')

  // タグ関連のロジックのみ
  const filteredTags = tags.filter(/* ... */)
  const toggleTag = (tagName: string) => {/* ... */}

  return (
    <div>
      <input
        value={tagSearch}
        onChange={(e) => setTagSearch(e.target.value)}
      />
      {filteredTags.map(/* ... */)}
    </div>
  )
}
```

**判断基準:**
- 1つのコンポーネントが複数の独立した機能を持っていたら分割
- 「〇〇と△△を担当する」と説明に"と"が入るなら分割を検討

### Props vs State の判断基準

**基本ルール:**
- **親コンポーネントが管理する値** → Props
- **自分だけが使う値** → State

**例: 検索機能**

```typescript
// ❌ 悪い例: 親子両方でstateを持つ
function ParentPage() {
  const [articles, setArticles] = useState<Article[]>([])

  return <SearchBar articles={articles} />
}

function SearchBar({ articles }: { articles: Article[] }) {
  const [keyword, setKeyword] = useState('')
  const [filtered, setFiltered] = useState<Article[]>([])

  // articlesが変わるたびに手動で更新
  useEffect(() => {
    setFiltered(articles.filter(/* ... */))
  }, [articles])

  // どちらが「正しい」データ？
}
```

```typescript
// ✅ 良い例: 親で一元管理
function ParentPage() {
  const [articles, setArticles] = useState<Article[]>([])
  const [keyword, setKeyword] = useState('')

  // 親でフィルタリング
  const filteredArticles = articles.filter((article) =>
    article.title.includes(keyword)
  )

  return (
    <div>
      <SearchBar keyword={keyword} onKeywordChange={setKeyword} />
      <ArticleList articles={filteredArticles} />
    </div>
  )
}

function SearchBar({
  keyword,
  onKeywordChange,
}: {
  keyword: string
  onKeywordChange: (keyword: string) => void
}) {
  return (
    <input
      value={keyword}
      onChange={(e) => onKeywordChange(e.target.value)}
    />
  )
}
```

**判断フローチャート:**

```
値を使うコンポーネントが1つだけ？
  ↓ Yes
  State（そのコンポーネント内で管理）

  ↓ No
複数のコンポーネントで共有する？
  ↓ Yes
  Props（親で管理して渡す）

  ↓ さらに多くのコンポーネントで共有する？
  ↓ Yes
  Context API または 状態管理ライブラリ
```

### カスタムフックを作るべきタイミング

**作るべき3つのサイン:**

#### 1. 同じロジックを複数のコンポーネントで使っている

**悪い例:**
```typescript
// ArticlesPage.tsx
function ArticlesPage() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true)
        const data = await articleClient.getAll()
        setArticles(data)
      } catch (err) {
        setError(err as Error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  // ...
}

// FavoritesPage.tsx
function FavoritesPage() {
  // 全く同じコードをコピペ
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true)
        const data = await articleClient.getFavorites()
        setArticles(data)
      } catch (err) {
        setError(err as Error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  // ...
}
```

**良い例:**
```typescript
// hooks/useArticles.ts
function useArticles() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true)
        const data = await articleClient.getAll()
        setArticles(data)
      } catch (err) {
        setError(err as Error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  return { articles, loading, error }
}

// ArticlesPage.tsx
function ArticlesPage() {
  const { articles, loading, error } = useArticles()
  // ...
}

// FavoritesPage.tsx
function FavoritesPage() {
  const { articles, loading, error } = useArticles()
  // ...
}
```

#### 2. useEffectとstateが密接に関連している

**良い例:**
```typescript
// hooks/useDebounce.ts
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState(value)

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)

    return () => clearTimeout(timer)
  }, [value, delay])

  return debouncedValue
}

// SearchBar.tsx で使用
function SearchBar() {
  const [keyword, setKeyword] = useState('')
  const debouncedKeyword = useDebounce(keyword, 500)

  useEffect(() => {
    // debouncedKeywordが変わった時だけ検索
    searchArticles(debouncedKeyword)
  }, [debouncedKeyword])

  return (
    <input
      value={keyword}
      onChange={(e) => setKeyword(e.target.value)}
    />
  )
}
```

#### 3. ビジネスロジックをテストしたい

カスタムフックにすることで、ロジックを独立してテストできます。

```typescript
// hooks/useArticles.test.ts
import { renderHook, waitFor } from '@testing-library/react'
import { useArticles } from './useArticles'

describe('useArticles', () => {
  it('記事一覧を取得できる', async () => {
    const { result } = renderHook(() => useArticles())

    expect(result.current.loading).toBe(true)

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(result.current.articles).toHaveLength(10)
  })
})
```

### 実際のプロジェクト構造

```
frontend/
├── app/                    # Next.js App Router
│   ├── layout.tsx         # 全ページ共通レイアウト
│   ├── page.tsx           # トップページ
│   └── articles/
│       ├── page.tsx       # 記事一覧ページ
│       ├── new/
│       │   └── page.tsx   # 記事作成ページ
│       └── [id]/
│           ├── page.tsx   # 記事詳細ページ
│           └── edit/
│               └── page.tsx  # 記事編集ページ
│
├── components/            # 再利用可能なコンポーネント
│   ├── ArticleCard.tsx   # 記事カード（50-100行）
│   ├── ArticleForm.tsx   # 記事フォーム（600行→分割を検討）
│   ├── SearchBar.tsx     # 検索バー（50行）
│   ├── TagList.tsx       # タグリスト（30行）
│   ├── Header.tsx        # ヘッダー（50行）
│   └── Sidebar.tsx       # サイドバー（80行）
│
├── hooks/                 # カスタムフック（ビジネスロジック）
│   ├── useArticles.ts    # 記事管理フック
│   ├── useArticleSearch.ts  # 記事検索フック
│   └── useTags.ts        # タグ管理フック
│
├── lib/                   # ユーティリティ・API
│   ├── api/
│   │   ├── baseClient.ts     # API基底クラス
│   │   ├── articleClient.ts  # 記事API
│   │   └── tagClient.ts      # タグAPI
│   ├── errors/
│   │   └── ApiError.ts       # カスタムエラー
│   └── utils/
│       └── sample.ts         # ユーティリティ関数
│
├── contexts/              # Context API
│   └── ToastContext.tsx  # 通知機能
│
├── types/                 # TypeScript型定義
│   ├── article.ts        # 記事の型
│   └── tag.ts            # タグの型
│
└── config/                # 設定ファイル
    └── constants.ts      # 定数
```

**各ディレクトリの役割:**

- **app/**: ページ定義。ルーティングに対応。ビジネスロジックは最小限に
- **components/**: UI部品。propsを受け取って表示するだけ
- **hooks/**: ビジネスロジック。データ取得や状態管理
- **lib/api/**: API通信。HTTPリクエストのみを担当
- **contexts/**: グローバルな状態。アプリ全体で使う値
- **types/**: 型定義。データ構造を定義

**設計のベストプラクティス:**

1. **1ファイル1コンポーネント**: 複数のコンポーネントを1ファイルに書かない
2. **名前でファイルを見つけられる**: ArticleCardコンポーネントはArticleCard.tsx
3. **テストファイルは隣に配置**: ArticleCard.tsx と ArticleCard.test.tsx
4. **関連するファイルをグループ化**: article関連はまとめる

---

## テストの書き方 - 品質を保証する

テストは「コードが正しく動くことを証明する」ための重要な作業です。このプロジェクトではVitestとReact Testing Libraryを使用しています。

### なぜテストが必要なのか？

**テストがない場合:**
```
あなた「記事作成機能を修正しました」
レビュアー「他の機能が壊れていないか確認しましたか？」
あなた「全部手動で確認するのは大変です...」
```

**テストがある場合:**
```
あなた「記事作成機能を修正しました」
あなた「npm test を実行して、全てのテストが通ることを確認しました」
レビュアー「素晴らしい！」
```

**テストの3つのメリット:**
1. **バグの早期発見**: コードを書いた直後にバグを見つけられる
2. **リファクタリングの安全性**: 変更後もテストが通れば安心
3. **ドキュメントとしての役割**: テストコードを見れば使い方が分かる

### テストの基本構造

**テストの3ステップ (AAA パターン):**

1. **Arrange (準備)**: テストに必要なデータや状態を準備
2. **Act (実行)**: テストしたい処理を実行
3. **Assert (検証)**: 結果が期待通りか確認

```typescript
import { describe, it, expect } from 'vitest'

describe('足し算関数', () => {
  it('2つの数を正しく足す', () => {
    // Arrange (準備)
    const a = 2
    const b = 3

    // Act (実行)
    const result = add(a, b)

    // Assert (検証)
    expect(result).toBe(5)
  })
})
```

**テストの構造:**
- `describe`: テストのグループ化（「〇〇のテスト」）
- `it`: 個別のテストケース（「△△できること」）
- `expect`: 検証（「結果が〇〇であること」）

### コンポーネントのテスト

#### 例1: 基本的なレンダリングテスト

```typescript
// ArticleCard.test.tsx
import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import ArticleCard from '@/components/ArticleCard'

describe('ArticleCard', () => {
  // テスト用のデータ
  const mockArticle = {
    id: 1,
    title: 'Go言語入門',
    url: 'https://example.com',
    summary: 'Go言語の基礎を学びます',
    tags: ['Go', 'プログラミング'],
    memo: '',
    createdAt: '2024-01-01 10:00:00',
    updatedAt: '2024-01-01 10:00:00',
  }

  it('記事の情報が正しく表示される', () => {
    // Arrange & Act: コンポーネントをレンダリング
    render(<ArticleCard article={mockArticle} />)

    // Assert: 要素が存在することを確認
    expect(screen.getByText('Go言語入門')).toBeInTheDocument()
    expect(screen.getByText('Go言語の基礎を学びます')).toBeInTheDocument()
    expect(screen.getByText('https://example.com')).toBeInTheDocument()
  })

  it('タグが表示される', () => {
    render(<ArticleCard article={mockArticle} />)

    expect(screen.getByText('Go')).toBeInTheDocument()
    expect(screen.getByText('プログラミング')).toBeInTheDocument()
  })
})
```

**重要な関数:**

- **`render()`**: コンポーネントを仮想DOMにレンダリング
- **`screen.getByText()`**: テキストで要素を見つける
- **`toBeInTheDocument()`**: 要素がDOMに存在することを確認

**要素を見つける方法（優先度順）:**

1. **`getByRole()`**: アクセシビリティロール（推奨）
   ```typescript
   screen.getByRole('button', { name: /送信/ })
   ```

2. **`getByLabelText()`**: フォームのラベル
   ```typescript
   screen.getByLabelText(/タイトル/)
   ```

3. **`getByPlaceholderText()`**: プレースホルダー
   ```typescript
   screen.getByPlaceholderText('記事のタイトルを入力')
   ```

4. **`getByText()`**: テキスト内容
   ```typescript
   screen.getByText('Go言語入門')
   ```

5. **`getByTestId()`**: 最終手段（テスト専用の属性）
   ```typescript
   <div data-testid="article-card">...</div>
   screen.getByTestId('article-card')
   ```

#### 例2: ユーザー操作のテスト

```typescript
// ArticleForm.test.tsx (実際のプロジェクトから抜粋)
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import ArticleForm from '@/components/ArticleForm'

describe('ArticleForm', () => {
  it('タイトルを入力できる', () => {
    render(<ArticleForm />)

    // 入力フィールドを取得
    const titleInput = screen.getByLabelText(/タイトル/) as HTMLInputElement

    // ユーザーが入力する動作をシミュレート
    fireEvent.change(titleInput, { target: { value: 'Go言語入門' } })

    // 入力が反映されたか確認
    expect(titleInput.value).toBe('Go言語入門')
  })

  it('タイトルが空の場合、エラーメッセージが表示される', async () => {
    render(<ArticleForm />)

    const titleInput = screen.getByLabelText(/タイトル/)

    // フォーカスして離れる（blur）
    fireEvent.focus(titleInput)
    fireEvent.blur(titleInput)

    // エラーメッセージが表示されるまで待つ
    await waitFor(() => {
      expect(screen.getByText(/タイトルは必須です/)).toBeInTheDocument()
    })
  })

  it('すべての必須項目が入力されている場合、送信ボタンが有効になる', async () => {
    render(<ArticleForm />)

    const titleInput = screen.getByLabelText(/タイトル/)
    const urlInput = screen.getByLabelText(/URL/)
    const summaryInput = screen.getByLabelText(/要約/)

    // すべてのフィールドに入力
    fireEvent.change(titleInput, { target: { value: 'Go言語入門' } })
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
    fireEvent.change(summaryInput, { target: { value: 'Go言語の基礎を学びます' } })

    // ボタンが有効になるまで待つ
    await waitFor(() => {
      const submitButton = screen.getByRole('button', { name: /登録/ })
      expect(submitButton).not.toBeDisabled()
    })
  })
})
```

**重要な関数:**

- **`fireEvent.change()`**: 入力フィールドの値を変更
- **`fireEvent.click()`**: ボタンをクリック
- **`fireEvent.blur()`**: フォーカスを外す
- **`waitFor()`**: 非同期処理が完了するまで待つ

#### 例3: APIをモック化したテスト

```typescript
// ArticleForm.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import ArticleForm from '@/components/ArticleForm'
import { articleClient } from '@/lib/api/articleClient'
import { tagClient } from '@/lib/api/tagClient'

// API クライアントをモック化
vi.mock('@/lib/api/articleClient')
vi.mock('@/lib/api/tagClient')

// Next.jsのuseRouterをモック化
const mockPush = vi.fn()
vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}))

describe('ArticleForm - API連携', () => {
  const mockTags = [
    { id: 1, name: 'Go', createdAt: '2024-01-01', updatedAt: '2024-01-01' },
    { id: 2, name: 'React', createdAt: '2024-01-02', updatedAt: '2024-01-02' },
  ]

  beforeEach(() => {
    // 各テスト前にモックをクリア
    vi.clearAllMocks()

    // タグ一覧取得のモック（デフォルト）
    vi.mocked(tagClient.getAll).mockResolvedValue(mockTags)
  })

  it('正常系：記事が正常に作成される', async () => {
    // モックの戻り値を設定
    const mockCreatedArticle = {
      id: 1,
      title: 'Go言語入門',
      url: 'https://example.com',
      summary: 'Go言語の基礎を学びます',
      tags: ['Go'],
      memo: '後で読む',
      createdAt: '2024-01-01 10:00:00',
      updatedAt: '2024-01-01 10:00:00',
    }
    vi.mocked(articleClient.create).mockResolvedValue(mockCreatedArticle)

    render(<ArticleForm />)

    // タグが読み込まれるまで待機
    await waitFor(() => {
      expect(screen.getByText('Go')).toBeInTheDocument()
    })

    // フォーム入力
    const titleInput = screen.getByLabelText(/タイトル/)
    const urlInput = screen.getByLabelText(/URL/)
    const summaryInput = screen.getByLabelText(/要約/)

    fireEvent.change(titleInput, { target: { value: 'Go言語入門' } })
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
    fireEvent.change(summaryInput, { target: { value: 'Go言語の基礎を学びます' } })

    // タグを選択
    const goTagButton = screen.getByText('Go').closest('button')
    fireEvent.click(goTagButton!)

    // フォーム送信
    const submitButton = screen.getByRole('button', { name: /登録/ })
    fireEvent.click(submitButton)

    // API が正しい引数で呼ばれたことを確認
    await waitFor(() => {
      expect(articleClient.create).toHaveBeenCalledWith({
        title: 'Go言語入門',
        url: 'https://example.com',
        summary: 'Go言語の基礎を学びます',
        tags: ['Go'],
        memo: '',
      })
    })

    // リダイレクトされることを確認
    expect(mockPush).toHaveBeenCalledWith('/articles')
  })

  it('異常系：API エラーが発生した場合、エラーメッセージが表示される', async () => {
    // エラーをモック
    vi.mocked(articleClient.create).mockRejectedValue(new Error('Failed to create article'))

    render(<ArticleForm />)

    // フォーム入力
    const titleInput = screen.getByLabelText(/タイトル/)
    const urlInput = screen.getByLabelText(/URL/)
    const summaryInput = screen.getByLabelText(/要約/)

    fireEvent.change(titleInput, { target: { value: 'テスト記事' } })
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
    fireEvent.change(summaryInput, { target: { value: 'テスト要約' } })

    // フォーム送信
    const submitButton = screen.getByRole('button', { name: /登録/ })
    fireEvent.click(submitButton)

    // エラーメッセージが表示されることを確認
    await waitFor(() => {
      expect(screen.getByText(/記事の作成に失敗しました/)).toBeInTheDocument()
    })

    // リダイレクトされないことを確認
    expect(mockPush).not.toHaveBeenCalled()
  })
})
```

**モックの重要なポイント:**

1. **`vi.mock()`**: モジュール全体をモック化
   ```typescript
   vi.mock('@/lib/api/articleClient')
   ```

2. **`vi.mocked()`**: モック関数に型をつける
   ```typescript
   vi.mocked(articleClient.create).mockResolvedValue(mockData)
   ```

3. **`mockResolvedValue()`**: 成功時の戻り値を設定
   ```typescript
   mockGetAll.mockResolvedValue(mockArticles)
   ```

4. **`mockRejectedValue()`**: エラー時の戻り値を設定
   ```typescript
   mockCreate.mockRejectedValue(new Error('Failed'))
   ```

5. **`beforeEach()`**: 各テスト前に実行される処理
   ```typescript
   beforeEach(() => {
     vi.clearAllMocks()  // モックをクリア
   })
   ```

### カスタムフックのテスト

```typescript
// useArticles.test.ts (実際のプロジェクトから抜粋)
import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useArticles, __resetArticlesCache } from './useArticles'
import { articleClient } from '@/lib/api/articleClient'

vi.mock('@/lib/api/articleClient')

describe('useArticles', () => {
  const mockArticles = [
    {
      id: 1,
      title: 'Test Article 1',
      url: 'https://example.com/1',
      summary: 'Test summary 1',
      tags: ['test'],
      memo: '',
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
    },
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    __resetArticlesCache()  // キャッシュをクリア
  })

  it('初期状態ではloading=true, articles=[]であること', () => {
    vi.mocked(articleClient.getAll).mockResolvedValue(mockArticles)

    // フックをレンダリング
    const { result } = renderHook(() => useArticles())

    // 初期状態を確認
    expect(result.current.loading).toBe(true)
    expect(result.current.articles).toEqual([])
    expect(result.current.error).toBeNull()
  })

  it('記事一覧取得成功時、articlesに記事が格納されること', async () => {
    vi.mocked(articleClient.getAll).mockResolvedValue(mockArticles)

    const { result } = renderHook(() => useArticles())

    // 非同期処理完了を待つ
    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(result.current.articles).toEqual(mockArticles)
    expect(result.current.error).toBeNull()
  })

  it('記事削除成功時、一覧から該当記事が削除されること', async () => {
    vi.mocked(articleClient.getAll).mockResolvedValue(mockArticles)
    vi.mocked(articleClient.delete).mockResolvedValue(undefined)

    const { result } = renderHook(() => useArticles())

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(result.current.articles).toHaveLength(1)

    // 記事削除を実行
    await result.current.deleteArticle(1)

    // 該当記事が削除されたか確認
    await waitFor(() => {
      expect(result.current.articles).toHaveLength(0)
    })
  })
})
```

**カスタムフックテストの重要ポイント:**

- **`renderHook()`**: フックをテスト環境でレンダリング
- **`result.current`**: フックの現在の戻り値にアクセス
- **`waitFor()`**: 状態更新を待つ
- **独立性**: 各テストでキャッシュをリセットして独立させる

### テスト実行コマンド

```bash
# 全テストを実行
npm test

# 特定のファイルのテストを実行
npm test ArticleForm.test.tsx

# watchモード（ファイル変更時に自動実行）
npm test -- --watch

# カバレッジレポートを生成
npm test -- --coverage
```

**カバレッジとは:**
- テストがコードのどれくらいをカバーしているかの指標
- 80%以上が目標
- 100%である必要はない（UIの細かい部分など）

### テスト作成の実践ガイド

**最初に書くべきテスト（優先度順）:**

1. **正常系**: 基本的な機能が動くことを確認
   ```typescript
   it('記事を作成できる', async () => {
     // 基本的な成功パターン
   })
   ```

2. **バリデーション**: 入力チェックが動くことを確認
   ```typescript
   it('タイトルが空の場合、エラーメッセージが表示される', async () => {
     // バリデーションエラーのパターン
   })
   ```

3. **異常系**: エラーが適切に処理されることを確認
   ```typescript
   it('API エラーが発生した場合、エラーメッセージが表示される', async () => {
     // エラーハンドリングのパターン
   })
   ```

4. **エッジケース**: 境界値や特殊なパターン
   ```typescript
   it('タグが100個ある場合でも正常に動作する', async () => {
     // 大量データのパターン
   })
   ```

**良いテストの特徴:**

✅ **独立している**: 他のテストに影響しない
✅ **高速**: 1秒以内に終わる
✅ **明確**: テスト名を見れば何をテストしているか分かる
✅ **安定している**: 実行するたびに同じ結果になる

❌ **悪いテストの特徴:**

❌ **不安定**: たまに失敗する（フレイキーテスト）
❌ **遅い**: 実行に時間がかかる
❌ **複雑**: 何をテストしているか分からない
❌ **依存している**: 他のテストが失敗すると失敗する

**参考:**
- [Vitest - Getting Started](https://vitest.dev/guide/)
- [React Testing Library - Intro](https://testing-library.com/docs/react-testing-library/intro/)
- [Testing Library - Queries](https://testing-library.com/docs/queries/about)

---

## トラブルシューティングガイド

実際の開発でよく遭遇するエラーと、その解決方法をまとめます。

### よくあるエラーと対処法

#### エラー1: Cannot read property 'X' of undefined

**エラーメッセージ:**
```
TypeError: Cannot read property 'title' of undefined
```

**原因:**
オブジェクトが`undefined`または`null`なのに、そのプロパティにアクセスしようとしている

**悪いコード:**
```typescript
function ArticleCard({ article }: { article: Article }) {
  return <h1>{article.title}</h1>
}

// articleがundefinedの場合にエラー
<ArticleCard article={undefined} />
```

**解決策1: Optional Chaining (`?.`)**
```typescript
function ArticleCard({ article }: { article: Article | undefined }) {
  return <h1>{article?.title}</h1>  // articleがundefinedでもエラーにならない
}
```

**解決策2: 早期リターン**
```typescript
function ArticleCard({ article }: { article: Article | undefined }) {
  if (!article) {
    return <div>記事がありません</div>
  }

  return <h1>{article.title}</h1>
}
```

**解決策3: デフォルト値**
```typescript
function ArticleCard({ article }: { article?: Article }) {
  const title = article?.title ?? '無題'
  return <h1>{title}</h1>
}
```

#### エラー2: React Hook "useXXX" is called conditionally

**エラーメッセージ:**
```
Error: React Hook "useState" is called conditionally. React Hooks must be called in the exact same order in every component render.
```

**原因:**
Hooksを条件文やループの中で呼んでいる

**悪いコード:**
```typescript
function Component() {
  if (condition) {
    const [state, setState] = useState('')  // ❌ 条件文の中でHook
  }

  return <div>...</div>
}
```

**正しいコード:**
```typescript
function Component() {
  const [state, setState] = useState('')  // ✅ トップレベルでHook

  if (condition) {
    // stateを使う処理
  }

  return <div>...</div>
}
```

**Hooksのルール:**
1. **トップレベルでのみ呼ぶ**: 条件文、ループ、ネストした関数の中で呼ばない
2. **React関数の中で呼ぶ**: 通常のJavaScript関数の中では呼ばない
3. **順序を保つ**: 毎回同じ順序でHooksを呼ぶ

#### エラー3: Hydration failed

**エラーメッセージ:**
```
Error: Hydration failed because the initial UI does not match what was rendered on the server.
```

**原因:**
サーバー側とクライアント側でレンダリング結果が異なる

**よくある原因:**

1. **日付・時刻を直接使う**
   ```typescript
   // ❌ 悪い例
   <div>{new Date().toString()}</div>
   ```

   ```typescript
   // ✅ 良い例
   'use client'  // クライアントコンポーネントにする

   function TimeDisplay() {
     const [time, setTime] = useState<string | null>(null)

     useEffect(() => {
       setTime(new Date().toString())
     }, [])

     if (!time) return null

     return <div>{time}</div>
   }
   ```

2. **localStorage/sessionStorageを使う**
   ```typescript
   // ❌ 悪い例
   const value = localStorage.getItem('key')

   // ✅ 良い例
   const [value, setValue] = useState<string | null>(null)

   useEffect(() => {
     setValue(localStorage.getItem('key'))
   }, [])
   ```

3. **ランダムな値を生成**
   ```typescript
   // ❌ 悪い例
   <div key={Math.random()}>...</div>

   // ✅ 良い例
   <div key={item.id}>...</div>
   ```

#### エラー4: Cannot update a component while rendering a different component

**エラーメッセージ:**
```
Warning: Cannot update a component (`Parent`) while rendering a different component (`Child`).
```

**原因:**
子コンポーネントのレンダリング中に親コンポーネントの状態を更新しようとしている

**悪いコード:**
```typescript
function Parent() {
  const [count, setCount] = useState(0)

  return <Child onRender={() => setCount(count + 1)} />
}

function Child({ onRender }: { onRender: () => void }) {
  onRender()  // ❌ レンダリング中に親の状態を更新
  return <div>Child</div>
}
```

**正しいコード:**
```typescript
function Parent() {
  const [count, setCount] = useState(0)

  return (
    <Child
      onMount={() => {
        setCount(count + 1)
      }}
    />
  )
}

function Child({ onMount }: { onMount: () => void }) {
  useEffect(() => {
    onMount()  // ✅ useEffectで実行
  }, [onMount])

  return <div>Child</div>
}
```

#### エラー5: Network Error / CORS Error

**エラーメッセージ:**
```
Access to fetch at 'http://localhost:8080/api/articles' from origin 'http://localhost:3000' has been blocked by CORS policy
```

**原因:**
バックエンドAPIがCORS（Cross-Origin Resource Sharing）を許可していない

**解決策1: バックエンドでCORSを許可（開発環境）**
```go
// Go言語のサーバー側
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type"},
    AllowCredentials: true,
}))
```

**解決策2: Next.jsのリライト機能を使う**
```typescript
// next.config.ts
const config = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ]
  },
}
```

**解決策3: 環境変数を確認**
```.env.local
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

確認コマンド:
```bash
echo $NEXT_PUBLIC_API_BASE_URL
```

#### エラー6: Maximum update depth exceeded

**エラーメッセージ:**
```
Error: Maximum update depth exceeded. This can happen when a component repeatedly calls setState inside componentWillUpdate or componentDidUpdate.
```

**原因:**
無限ループで状態を更新している

**悪いコード:**
```typescript
function Component() {
  const [count, setCount] = useState(0)

  // ❌ 無限ループ: レンダリングのたびに状態を更新
  setCount(count + 1)

  return <div>{count}</div>
}
```

**正しいコード:**
```typescript
function Component() {
  const [count, setCount] = useState(0)

  // ✅ useEffectで実行
  useEffect(() => {
    setCount(count + 1)
  }, [])  // 依存配列を空にして1回だけ実行

  return <div>{count}</div>
}
```

**よくあるパターン:**
```typescript
// ❌ 無限ループ
useEffect(() => {
  setCount(count + 1)
}, [count])  // countが変わるたびに実行 → countを更新 → 再実行...

// ✅ 正しい
useEffect(() => {
  setCount(prev => prev + 1)  // 関数形式で更新
}, [])  // 空配列で1回だけ実行
```

### デバッグ方法

#### 1. console.logでデバッグ

**基本:**
```typescript
function Component() {
  const [data, setData] = useState(null)

  console.log('レンダリング時のdata:', data)

  useEffect(() => {
    console.log('useEffectが実行された')
    fetchData().then((result) => {
      console.log('データ取得完了:', result)
      setData(result)
    })
  }, [])

  return <div>{data}</div>
}
```

**オブジェクトを見やすく表示:**
```typescript
console.log('データ:', JSON.stringify(data, null, 2))
```

**テーブル形式で表示:**
```typescript
console.table(articles)
```

#### 2. React Developer Toolsを使う

**インストール:**
- Chrome拡張機能: "React Developer Tools"

**使い方:**
1. ブラウザのDevToolsを開く（F12）
2. "Components"タブを選択
3. コンポーネントツリーを確認
4. 選んだコンポーネントの props と state を確認
5. フックの値を確認

**便利な機能:**
- **検索**: コンポーネント名で検索
- **Highlight**: 画面上でコンポーネントをハイライト
- **Source**: ソースコードにジャンプ

#### 3. Network タブで API 確認

**使い方:**
1. DevToolsの"Network"タブを開く
2. "Fetch/XHR"フィルターを選択
3. ページをリロード
4. APIリクエストを確認

**確認項目:**
- **Status**: 200 (成功)、404 (見つからない)、500 (サーバーエラー)
- **Request Headers**: 送信したヘッダー
- **Request Payload**: 送信したデータ
- **Response**: サーバーからの応答

#### 4. TypeScriptのエラーを確認

```bash
# 型チェックを実行
npm run type-check

# または
npx tsc --noEmit
```

**VSCodeの設定:**
- "TypeScript: Error Checking"を有効にする
- 保存時に自動チェック

### パフォーマンス問題の診断

#### React Developer Tools の Profiler

**使い方:**
1. "Profiler"タブを開く
2. 録画ボタンをクリック
3. アプリを操作
4. 録画を停止
5. 各コンポーネントのレンダリング時間を確認

**改善ポイント:**
- レンダリング時間が長いコンポーネントを特定
- 不必要な再レンダリングを見つける
- React.memoやuseCallbackで最適化

#### Lighthouse で全体のパフォーマンス測定

**使い方:**
1. DevToolsの"Lighthouse"タブを開く
2. "Generate report"をクリック
3. レポートを確認

**確認項目:**
- Performance: 読み込み速度
- Accessibility: アクセシビリティ
- Best Practices: ベストプラクティス
- SEO: 検索エンジン最適化

---

## 実践的な開発フロー - 新機能を追加する手順

実際のプロジェクトで新しい機能を追加する流れを、ステップバイステップで説明します。

### ステップ1: 要件を理解する

**例: 「記事にお気に入り機能を追加する」**

**まず確認すべきこと:**

1. **機能の詳細**
   - お気に入りボタンはどこに配置？（記事カード、記事詳細）
   - ボタンのデザインは？（ハートアイコン、星アイコン）
   - お気に入り済みの記事はどう表示？（色を変える、アイコンを塗りつぶす）

2. **データ構造**
   - バックエンドのAPIは用意されている？
   - データベースに保存される？
   - ユーザーごとに管理する？

3. **既存コードへの影響**
   - 既存のコンポーネントを修正する？
   - 新しいコンポーネントを作る？

**曖昧なら質問する:**
- 「お気に入り済みの記事一覧ページは必要ですか？」
- 「お気に入りの上限数はありますか？」
- 「お気に入りの解除はどこからできますか？」

### ステップ2: 設計を考える

**2-1. データの流れを図にする**

```
┌─────────────────────┐
│  ArticleCard        │
│  [♡ お気に入り]      │ ← ユーザーがクリック
└──────────┬──────────┘
           │ onFavorite(articleId)
           ▼
┌─────────────────────┐
│  ArticlesPage       │
│  状態管理            │ ← お気に入り状態を管理
└──────────┬──────────┘
           │ API呼び出し
           ▼
┌─────────────────────┐
│  favoriteClient     │
│  API通信            │ ← POST /api/favorites
└─────────────────────┘
```

**2-2. 必要なファイルをリストアップ**

**新規作成:**
- `lib/api/favoriteClient.ts` - お気に入りAPI
- `hooks/useFavorites.ts` - お気に入り管理フック
- `types/favorite.ts` - お気に入りの型定義

**修正:**
- `components/ArticleCard.tsx` - お気に入りボタンを追加
- `types/article.ts` - `isFavorite`プロパティを追加

**テスト:**
- `components/ArticleCard.test.tsx` - テストを追加
- `hooks/useFavorites.test.ts` - テストを追加

**2-3. 実装の順序を決める**

1. 型定義（types/favorite.ts）
2. APIクライアント（lib/api/favoriteClient.ts）
3. カスタムフック（hooks/useFavorites.ts）
4. UIコンポーネント（components/ArticleCard.tsx）
5. テスト
6. 動作確認

### ステップ3: 実装する

**3-1. 型定義を作成**

```typescript
// types/favorite.ts
export interface Favorite {
  id: number
  articleId: number
  userId: number
  createdAt: string
}

export interface CreateFavoriteInput {
  articleId: number
}
```

**3-2. APIクライアントを作成**

```typescript
// lib/api/favoriteClient.ts
import { BaseApiClient } from './baseClient'
import { Favorite, CreateFavoriteInput } from '@/types/favorite'

class FavoriteClient extends BaseApiClient {
  // お気に入り一覧を取得
  async getAll(): Promise<Favorite[]> {
    return await this.fetchWithErrorHandling<Favorite[]>('/api/favorites')
  }

  // お気に入りに追加
  async create(input: CreateFavoriteInput): Promise<Favorite> {
    return await this.fetchWithErrorHandling<Favorite>('/api/favorites', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(input),
    })
  }

  // お気に入りから削除
  async delete(articleId: number): Promise<void> {
    await this.fetchWithErrorHandling<void>(`/api/favorites/${articleId}`, {
      method: 'DELETE',
    })
  }
}

export const favoriteClient = new FavoriteClient()
```

**3-3. カスタムフックを作成**

```typescript
// hooks/useFavorites.ts
import { useState, useEffect, useCallback } from 'react'
import { favoriteClient } from '@/lib/api/favoriteClient'

export function useFavorites() {
  const [favoriteArticleIds, setFavoriteArticleIds] = useState<Set<number>>(new Set())
  const [loading, setLoading] = useState(true)

  // お気に入り一覧を取得
  useEffect(() => {
    async function fetchFavorites() {
      try {
        setLoading(true)
        const favorites = await favoriteClient.getAll()
        const ids = new Set(favorites.map(f => f.articleId))
        setFavoriteArticleIds(ids)
      } catch (err) {
        console.error('お気に入りの取得に失敗しました', err)
      } finally {
        setLoading(false)
      }
    }

    fetchFavorites()
  }, [])

  // お気に入りに追加
  const addFavorite = useCallback(async (articleId: number) => {
    try {
      await favoriteClient.create({ articleId })

      // 楽観的UI更新
      setFavoriteArticleIds(prev => new Set([...prev, articleId]))
    } catch (err) {
      console.error('お気に入りの追加に失敗しました', err)
      throw err
    }
  }, [])

  // お気に入りから削除
  const removeFavorite = useCallback(async (articleId: number) => {
    try {
      await favoriteClient.delete(articleId)

      // 楽観的UI更新
      setFavoriteArticleIds(prev => {
        const newSet = new Set(prev)
        newSet.delete(articleId)
        return newSet
      })
    } catch (err) {
      console.error('お気に入りの削除に失敗しました', err)
      throw err
    }
  }, [])

  // お気に入り判定
  const isFavorite = useCallback((articleId: number) => {
    return favoriteArticleIds.has(articleId)
  }, [favoriteArticleIds])

  return {
    favoriteArticleIds,
    loading,
    addFavorite,
    removeFavorite,
    isFavorite,
  }
}
```

**3-4. UIコンポーネントを修正**

```typescript
// components/ArticleCard.tsx
interface ArticleCardProps {
  article: Article
  isFavorite: boolean
  onFavorite: (id: number) => void
  onDelete?: (id: number) => void
}

function ArticleCard({ article, isFavorite, onFavorite, onDelete }: ArticleCardProps) {
  const handleFavoriteClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation()
    onFavorite(article.id)
  }, [article.id, onFavorite])

  return (
    <article>
      <h3>{article.title}</h3>

      {/* お気に入りボタン */}
      <button
        onClick={handleFavoriteClick}
        className={`p-2 rounded-lg transition-all ${
          isFavorite
            ? 'text-red-500 bg-red-50'
            : 'text-gray-400 hover:text-red-500 hover:bg-red-50'
        }`}
        aria-label={isFavorite ? 'お気に入りから削除' : 'お気に入りに追加'}
      >
        <svg className="w-5 h-5" fill={isFavorite ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
      </button>
    </article>
  )
}
```

**3-5. ページで使用**

```typescript
// app/articles/page.tsx
'use client'

import { useArticles } from '@/hooks/useArticles'
import { useFavorites } from '@/hooks/useFavorites'
import ArticleCard from '@/components/ArticleCard'

export default function ArticlesPage() {
  const { articles, loading: articlesLoading } = useArticles()
  const { isFavorite, addFavorite, removeFavorite, loading: favoritesLoading } = useFavorites()

  const handleFavorite = async (articleId: number) => {
    if (isFavorite(articleId)) {
      await removeFavorite(articleId)
    } else {
      await addFavorite(articleId)
    }
  }

  if (articlesLoading || favoritesLoading) {
    return <div>読み込み中...</div>
  }

  return (
    <div>
      {articles.map((article) => (
        <ArticleCard
          key={article.id}
          article={article}
          isFavorite={isFavorite(article.id)}
          onFavorite={handleFavorite}
        />
      ))}
    </div>
  )
}
```

### ステップ4: テストを書く

```typescript
// hooks/useFavorites.test.ts
import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { useFavorites } from './useFavorites'
import { favoriteClient } from '@/lib/api/favoriteClient'

vi.mock('@/lib/api/favoriteClient')

describe('useFavorites', () => {
  it('お気に入りに追加できる', async () => {
    vi.mocked(favoriteClient.getAll).mockResolvedValue([])
    vi.mocked(favoriteClient.create).mockResolvedValue({
      id: 1,
      articleId: 1,
      userId: 1,
      createdAt: '2024-01-01',
    })

    const { result } = renderHook(() => useFavorites())

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    await result.current.addFavorite(1)

    expect(result.current.isFavorite(1)).toBe(true)
  })
})
```

### ステップ5: 動作確認

**確認項目リスト:**

- [ ] お気に入りボタンが表示される
- [ ] ボタンをクリックするとお気に入りに追加される
- [ ] アイコンの色が変わる
- [ ] もう一度クリックするとお気に入りから削除される
- [ ] ページをリロードしても状態が保持される
- [ ] APIエラー時にエラーメッセージが表示される
- [ ] 複数の記事で同時に動作する
- [ ] レスポンシブデザインが崩れていない

**ブラウザで確認:**
```bash
npm run dev
```

**テストを実行:**
```bash
npm test
```

### ステップ6: コミット & プルリクエスト

**良いコミットメッセージ:**
```bash
git add .
git commit -m "$(cat <<'EOF'
feat: 記事にお気に入り機能を追加

- お気に入りAPIクライアントを実装
- useFavoritesフックを作成
- ArticleCardにお気に入りボタンを追加
- テストを追加

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
EOF
)"
```

**プルリクエストのテンプレート:**

```markdown
## 概要
記事にお気に入り機能を追加しました。

## 変更内容
- [x] お気に入りAPIクライアント (`lib/api/favoriteClient.ts`)
- [x] お気に入り管理フック (`hooks/useFavorites.ts`)
- [x] ArticleCardにお気に入りボタンを追加
- [x] テストを追加

## スクリーンショット
(お気に入りボタンの画像)

## テスト
- [x] ユニットテスト: `npm test`
- [x] 手動テスト: ブラウザで動作確認

## チェックリスト
- [x] TypeScriptの型エラーがない
- [x] テストが全て通る
- [x] ESLintエラーがない
- [x] 既存機能が壊れていない
```

---

## まとめ: 学んだことを振り返る

このドキュメントで学んだ重要な概念を振り返ります。

### Reactの基本

1. **コンポーネント** - UIを部品として作る
2. **useState** - 状態を管理する
3. **useEffect** - 副作用（データ取得など）を実行する
4. **Props** - コンポーネント間でデータを渡す
5. **useCallback** - 関数をメモ化してパフォーマンス向上
6. **useRef** - 再レンダリングを引き起こさない値の保持

### Next.jsの機能

1. **App Router** - ファイルベースのルーティング
2. **layout.tsx** - 全ページ共通のレイアウト
3. **`'use client'`** - クライアントコンポーネント
4. **Server Components** - サーバー側でレンダリング

### TypeScript

1. **型定義** - データの形を定義する
2. **interface** - オブジェクトの型を定義する
3. **型安全** - コンパイル時にエラーを検出
4. **Generics** - 再利用可能な型定義

### 実践的なパターン

1. **3層アーキテクチャ** - プレゼンテーション、ビジネスロジック、データアクセス
2. **カスタムフック** - ロジックを再利用する
3. **Context API** - グローバルな状態管理
4. **楽観的UI更新** - UX向上のテクニック
5. **キャッシング** - パフォーマンス向上

### コンポーネント設計

1. **分割の基準** - 200行以上、同じパターンの繰り返し、複数の責任
2. **Props vs State** - 親で管理するか、自分だけが使うか
3. **カスタムフックの作成** - 同じロジックの再利用、useEffectとstateの密接な関係
4. **ファイル構造** - 役割ごとにディレクトリを分ける

### テストの書き方

1. **AAA パターン** - Arrange（準備）、Act（実行）、Assert（検証）
2. **コンポーネントテスト** - render、screen、fireEvent を使う
3. **モック化** - vi.mock で外部依存を置き換える
4. **カスタムフックテスト** - renderHook を使う
5. **優先度** - 正常系 → バリデーション → 異常系 → エッジケース

### トラブルシューティング

1. **Cannot read property 'X' of undefined** - Optional Chaining (`?.`) を使う
2. **React Hook called conditionally** - Hooksはトップレベルで呼ぶ
3. **Hydration failed** - サーバーとクライアントで同じ結果にする
4. **Network Error / CORS** - バックエンドでCORSを許可
5. **デバッグ方法** - console.log、React DevTools、Networkタブ

### 実践的な開発フロー

1. **要件を理解** - 機能の詳細、データ構造、既存コードへの影響
2. **設計を考える** - データの流れ、必要なファイル、実装の順序
3. **実装する** - 型定義 → API → フック → UI → テスト
4. **テストを書く** - 各機能に対してテストを追加
5. **動作確認** - ブラウザとテストで確認
6. **コミット & PR** - 良いコミットメッセージ、詳細な説明

---

## 次のステップ: さらに学ぶために

このドキュメントで基礎から実践まで学びました。次のステップとして：

### 1. 実際に手を動かす

**おすすめの練習:**
- このプロジェクトに小さな機能を追加してみる
  - 記事の並び替え機能（タイトル順、日付順）
  - 記事のページネーション
  - タグごとのフィルタリング
- 既存のコンポーネントを別の方法で実装してみる
- テストを追加してみる

### 2. さらに深く学ぶ

**公式ドキュメント:**
- [React Documentation](https://react.dev/learn) - 公式チュートリアル
- [Next.js Documentation](https://nextjs.org/docs) - Next.jsの全機能
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html) - 型システムの詳細
- [Testing Library](https://testing-library.com/docs/react-testing-library/intro/) - テストの書き方

**おすすめの学習リソース:**
- [React + TypeScript Cheatsheets](https://react-typescript-cheatsheet.netlify.app/) - よくあるパターン集
- [Kent C. Dodds Blog](https://kentcdodds.com/blog) - React/テストのベストプラクティス
- [Patterns.dev](https://www.patterns.dev/) - デザインパターン

### 3. 実際のプロジェクトで経験を積む

- チームメンバーのコードレビューを積極的に見る
- 自分のコードをレビューしてもらう
- ペアプログラミングに参加する
- オープンソースプロジェクトにコントリビュートする

### 4. このドキュメントの活用方法

- **辞書として使う**: 忘れたことがあれば検索して読み返す
- **実装前に読む**: 新しい機能を実装する前に該当セクションを読む
- **メンターに質問する**: 分からないことがあれば遠慮なく聞く
- **フィードバックする**: このドキュメントの改善点があれば共有する

---

## 最後に

このドキュメントは、あなたが**自信を持ってコードを書けるようになる**ことを目指して作られました。

**覚えておいてほしいこと:**

1. **完璧を目指さない**: 最初から完璧なコードを書ける人はいません
2. **失敗を恐れない**: エラーは学びの機会です
3. **質問する**: 分からないことは恥ずかしいことではありません
4. **少しずつ成長**: 毎日少しずつ、着実に成長しましょう

**コーディングを楽しんでください！** 🚀