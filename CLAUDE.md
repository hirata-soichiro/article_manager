# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

## プロジェクト概要

Article Managerは、AI機能を搭載した技術記事管理のフルスタックWebアプリケーションです。システムは以下で構成されています：

- **バックエンド (API)**: Clean Architectureに従ったGo言語ベースのREST API
- **フロントエンド**: Next.js 16 (App Router) + React 19 + TypeScript
- **データベース**: MySQL 8.0
- **AI統合**: 記事生成と書籍推薦のためのGoogle Gemini API

## アーキテクチャ

### バックエンド (Go)

APIはClean Architectureに従い、明確な関心の分離を実現しています：

```
api/
├── cmd/server/main.go           # アプリケーションエントリーポイント、依存性注入
├── internal/
│   ├── domain/                  # エンタープライズビジネスルール
│   │   ├── entity/              # ドメインエンティティ (Article, Tag, BookRecommendation)
│   │   ├── repository/          # リポジトリインターフェース
│   │   ├── service/             # ドメインサービス
│   │   └── errors/              # ドメイン固有のエラー
│   ├── usecase/                 # アプリケーションビジネスルール
│   ├── interface/handler/       # HTTPハンドラー (コントローラー)
│   └── infrastructure/          # 外部依存
│       ├── repository/          # リポジトリ実装 (MySQL)
│       ├── database/            # DB接続、マイグレーションマネージャー
│       ├── ai/                  # Geminiクライアント
│       ├── logger/              # Zapロガー
│       └── service/             # インフラストラクチャサービス
```

**主要な設計パターン:**
- Clean Architecture: 依存関係は内側を向く (domain ← usecase ← interface ← infrastructure)
- Repository Pattern: MySQLとインメモリ実装によるデータアクセスの抽象化
- Dependency Injection: すべての依存関係はmain.goで注入

### フロントエンド (Next.js)

モダンなNext.js App Routerアーキテクチャ、関心の分離を実現：

```
frontend/
├── app/                         # Next.js App Routerページ
│   ├── layout.tsx               # ルートレイアウト
│   ├── page.tsx                 # ホームページ (ダッシュボード)
│   ├── articles/                # 記事ページ
│   │   ├── page.tsx             # 記事一覧
│   │   ├── new/page.tsx         # 記事作成
│   │   └── [id]/                # 動的ルート
│   └── tags/page.tsx            # タグ管理
├── components/                  # 再利用可能なUIコンポーネント
├── hooks/                       # データ取得用のカスタムReact Hooks
│   ├── useArticles.ts           # 記事のCRUD操作
│   ├── useArticleSearch.ts      # 検索機能
│   ├── useTags.ts               # タグ操作
│   └── useBookRecommendations.ts
├── types/                       # TypeScript型定義
├── contexts/                    # React Context (ToastContext)
└── config/constants.ts          # アプリ全体の定数
```

**主要なパターン:**
- カスタムフックでAPI呼び出しと状態管理をカプセル化
- 継承よりもコンポーネント合成
- 横断的関心事（トースト）にContext使用
- TypeScriptによる型安全なAPI通信

## 開発コマンド

### バックエンド (Go API)

```bash
# ホットリロードで実行 (Air)
cd api && air -c .air.toml

# バイナリをビルド
cd api && go build -o server cmd/server/main.go

# テスト実行
cd api && go test ./...

# リポジトリテスト実行（テストデータベースが必要）
# 1. テストデータベースを起動:
docker-compose -f docker-compose-test.yml up -d

# 2. マイグレーション実行（初回のみ）:
docker exec -i {container_id} mysql -utest_user -ptest_password article_manager_test < api/internal/infrastructure/database/migrations/000001_create_articles_table.up.sql

# 3. リポジトリテスト実行:
cd api && go test -v ./internal/infrastructure/repository/

# 4. テストデータベースを停止:
docker-compose -f docker-compose-test.yml down
```

### フロントエンド (Next.js)

```bash
cd frontend

# 開発サーバー起動
npm run dev

# 本番ビルド
npm run build

# 本番サーバー起動
npm run start

# Lint実行
npm run lint

# テスト実行
npm run test              # すべてのテストを実行
npm run test:ui           # UIでテストを実行
npm run test:coverage     # カバレッジレポート付きでテストを実行
```

### フルスタック (Docker Compose)

```bash
# すべてのサービスを起動（.envファイルが必要）
docker-compose up -d

# ログを表示
docker-compose logs -f [api|frontend|db]

# すべてのサービスを停止
docker-compose down

# コンテナを再ビルド
docker-compose up -d --build
```

**必要な環境変数（ルートの.env）:**
```
MYSQL_ROOT_PASSWORD=<root_password>
MYSQL_DATABASE=article_manager
MYSQL_USER=<db_user>
MYSQL_PASSWORD=<db_password>
DB_PORT=3306
API_PORT=8080
FRONTEND_PORT=3000
GEMINI_API_KEY=<your_gemini_api_key>
GOOGLE_BOOKS_API_KEY=<your_books_api_key>
```

## テスト戦略

### バックエンドテスト

- **ユニットテスト**: ドメインエンティティとサービス
- **統合テスト**: テストデータベースを使用したリポジトリ層
- **ハンドラーテスト**: モックusecaseを使用したHTTPハンドラー
- フレームワーク: testify/assert、testify/mock

テストファイルはGoの規約に従います: 実装ファイルと並んで`*_test.go`

### フロントエンドテスト

- **コンポーネントテスト**: ユーザー操作を含むUIコンポーネント
- **フックテスト**: APIモッキングを使用したカスタムフック
- フレームワーク: Vitest、Testing Library、Happy-DOM
- 設定: `vitest.config.ts`、`vitest.setup.ts`

開発中はフロントエンドテストを頻繁に実行してください - テストスイートは高速で包括的です。

## APIエンドポイント

Base URL: `http://localhost:8080/api`

**記事（Articles）:**
- `GET /articles` - 記事一覧取得
- `GET /articles/search?q={query}&tags={tag1,tag2}` - 記事検索
- `GET /articles/{id}` - 記事詳細取得
- `POST /articles` - 記事作成
- `POST /articles/generate` - URLからAIで記事を自動生成
- `PUT /articles/{id}` - 記事更新
- `DELETE /articles/{id}` - 記事削除

**タグ（Tags）:**
- `GET /tags` - タグ一覧取得
- `GET /tags/{id}` - タグ詳細取得
- `POST /tags` - タグ作成
- `PUT /tags/{id}` - タグ更新
- `DELETE /tags/{id}` - タグ削除

**書籍推薦（Book Recommendations）:**
- `GET /book-recommendations` - 記事に基づくAI生成の書籍推薦を取得

## データベースマイグレーション

アプリケーションは`infrastructure/database/migration.go`に組み込まれたマイグレーションシステムを使用します。マイグレーションはアプリケーション起動時に自動実行されます。

マイグレーションファイルはビルド時にバイナリに埋め込まれます（埋め込まれたSQLファイルは`database/migration.go`を参照）。

開発中の手動マイグレーション:
```bash
# 開発データベースに接続
docker exec -it db mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}

# テストデータベースに接続
docker exec -it mysql_test mysql -utest_user -ptest_password article_manager_test
```

## 主要な実装詳細

### バックエンド

1. **グレースフルシャットダウン**: APIはSIGTERM/SIGINTをハンドリングし、30秒のタイムアウトでグレースフルシャットダウンを実装
2. **CORS**: すべてのオリジン（`*`）を許可するように設定 - 本番環境では制限してください
3. **ロギング**: zapによる構造化ログ、環境固有の設定
4. **エラーハンドリング**: `handler/error_handler.go`で一元的なエラーハンドリング
5. **AI統合**: `infrastructure/ai/gemini_client.go`のGeminiクライアントで記事生成と書籍推薦

### フロントエンド

1. **API設定**: `config/constants.ts`でベースURLを設定、環境変数によるオーバーライドをサポート
2. **エラーバウンダリー**: `components/ErrorBoundary.tsx`がReactエラーをキャッチ
3. **トースト通知**: `contexts/ToastContext.tsx`がアプリ全体の通知を提供
4. **検索デバウンス**: API呼び出しを減らすため検索入力をデバウンス
5. **フォームバリデーション**: `config/constants.ts`の定数を使用したクライアント側バリデーション
6. **AI機能**:
   - URLからの記事生成と自動タグ抽出
   - 保存された記事に基づく書籍推薦

### Context7統合

このプロジェクトで使用されているライブラリの最新ドキュメントを取得するには、Claudeプロンプトに`use context7`を追加してください。Go、React、Next.js、その他の依存関係の最新ドキュメントにアクセスできます。

## よくある開発タスク

### 新しいAPIエンドポイントの追加

1. `domain/entity/`でドメインエンティティを定義
2. `domain/repository/`でリポジトリインターフェースを作成
3. `infrastructure/repository/`でリポジトリを実装
4. `usecase/`でusecaseを作成
5. `interface/handler/`でハンドラーを実装
6. `cmd/server/main.go`でルートを登録
7. 各レイヤーのテストを記述

### 新しいフロントエンド機能の追加

1. `types/`でTypeScript型を定義
2. データ取得のため`hooks/`でカスタムフックを作成
3. `components/`でUIコンポーネントを構築
4. `app/`でページを作成または更新
5. コンポーネントとフックのテストを記述
6. 必要に応じて`config/constants.ts`の定数を更新

### データベース問題のデバッグ

マイグレーション状態の確認:
```bash
# API起動時のマイグレーションログを表示
docker-compose logs api | grep -i migration

# データベーススキーマを確認
docker exec -it db mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}
mysql> SHOW TABLES;
mysql> DESCRIBE articles;
```

## 重要な注意事項

- **リポジトリテスト**: MySQLテストコンテナを先に起動する必要があります（上記コマンド参照）
- **Air**: 開発中のホットリロードに使用 - 設定は`.air.toml`
- **フロントエンドアーキテクチャ**: `frontend/ARCHITECTURE.md`に詳細な説明（新卒向けの包括的ガイド）
- **Clean Architecture**: 依存関係の方向を維持 - 外側のレイヤーから内側のレイヤーへのインポートは禁止
- **APIクライアント**: フロントエンドはネイティブのfetch APIを使用、外部HTTPクライアントライブラリは不使用
