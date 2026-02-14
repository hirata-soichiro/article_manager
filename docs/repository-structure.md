# リポジトリ構造定義書

## Article Manager - リポジトリ構成とファイル配置ルール

---

## 目次

1. [概要](#1-概要)
2. [リポジトリ全体構造](#2-リポジトリ全体構造)
3. [バックエンド（api/）](#3-バックエンドapi)
4. [フロントエンド（frontend/）](#4-フロントエンドfrontend)
5. [ドキュメント（docs/）](#5-ドキュメントdocs)
6. [設定ファイル・環境変数](#6-設定ファイル環境変数)
7. [ファイル配置ルール](#7-ファイル配置ルール)
8. [命名規則](#8-命名規則)
9. [テストファイル配置規則](#9-テストファイル配置規則)

---

## 1. 概要

このドキュメントは、Article Managerプロジェクトのリポジトリ構造、各ディレクトリの役割、ファイル配置ルールを定義します。

### 1.1 設計思想

- **関心の分離**: フロントエンド、バックエンド、ドキュメントを明確に分離
- **Clean Architecture**: バックエンドはドメイン駆動設計に基づく階層構造
- **テスタビリティ**: 実装ファイルと同じ場所にテストファイルを配置
- **保守性**: 明確な命名規則と配置ルールによる可読性向上

---

## 2. リポジトリ全体構造

```
article_manager/                    # リポジトリルート
├── api/                            # バックエンド（Go）
│   ├── cmd/                        # アプリケーションエントリーポイント
│   ├── internal/                   # 内部パッケージ（Clean Architecture）
│   ├── tmp/                        # Airホットリロード一時ファイル（Gitignore対象）
│   ├── .air.toml                   # Airホットリロード設定
│   ├── .gitignore                  # APIディレクトリGit除外設定
│   ├── Dockerfile                  # バックエンドコンテナ定義
│   ├── go.mod                      # Go依存関係管理
│   └── go.sum                      # Go依存関係ロックファイル
│
├── frontend/                       # フロントエンド（Next.js）
│   ├── app/                        # Next.js App Routerページ
│   ├── components/                 # 再利用可能なUIコンポーネント
│   ├── hooks/                      # カスタムReact Hooks
│   ├── lib/                        # ユーティリティ・ライブラリ
│   ├── types/                      # TypeScript型定義
│   ├── contexts/                   # React Context
│   ├── config/                     # 設定ファイル
│   ├── public/                     # 静的ファイル
│   ├── Dockerfile                  # フロントエンドコンテナ定義
│   ├── package.json                # npm依存関係
│   └── package-lock.json           # npm依存関係ロックファイル
│
├── docs/                           # 永続的ドキュメント
│   ├── product-requirements.md     # プロダクト要求定義書
│   ├── functional-design.md        # 機能設計書
│   ├── architecture.md             # 技術仕様書
│   ├── repository-structure.md     # リポジトリ構造定義書（本ドキュメント）
│   ├── development-guidelines.md   # 開発ガイドライン（作成予定）
│   └── glossary.md                 # ユビキタス言語定義（作成予定）
│
├── .steering/                      # 作業単位のステアリングファイル（作成予定）
│   └── [YYYYMMDD]-[title]/         # 日付とタイトルでディレクトリ作成
│       ├── requirements.md         # 作業要求定義
│       ├── design.md               # 作業設計
│       └── tasklist.md             # タスクリスト
│
├── db/                             # データベース永続化（Gitignore対象）
│   ├── data/                       # 開発環境MySQLデータ永続化（Docker Volume）
│   └── test_data/                  # テスト環境MySQLデータ永続化（Docker Volume）
│
├── .claude/                        # Claude Code設定（Gitignore対象）
│   ├── settings.local.json         # ローカル設定
│   └── agents/                     # カスタムエージェント定義
│
├── .env                            # 環境変数（Gitignore対象）
├── .env.example                    # 環境変数テンプレート（作成予定）
├── .gitignore                      # Git除外設定
├── .mcp.json                       # MCPサーバー設定
├── CLAUDE.md                       # プロジェクトメモリ（Claude Code用）
├── README.md                       # プロジェクト概要
├── docker-compose.yml              # Docker Compose設定（開発環境）
└── docker-compose-test.yml         # Docker Compose設定（テスト環境）
```

---

## 3. バックエンド（api/）

### 3.1 ディレクトリ構造

バックエンドはClean Architectureに基づく階層構造を採用しています。

```
api/
├── cmd/                                    # コマンド（アプリケーションエントリーポイント）
│   └── server/
│       └── main.go                         # メイン関数、依存性注入、サーバー起動
│
├── internal/                               # 内部パッケージ（外部から参照不可）
│   ├── domain/                             # ドメイン層（ビジネスルール）
│   │   ├── entity/                         # エンティティ（ドメインオブジェクト）
│   │   │   ├── article.go
│   │   │   ├── article_test.go
│   │   │   ├── book_recommendation.go
│   │   │   ├── tag.go
│   │   │   └── tag_test.go
│   │   ├── repository/                     # リポジトリインターフェース
│   │   │   ├── article_repository.go
│   │   │   ├── book_recommendation_repository.go
│   │   │   └── tag_repository.go
│   │   ├── service/                        # ドメインサービス
│   │   │   ├── ai_generator.go
│   │   │   └── book_recommendation_service.go
│   │   └── errors/                         # ドメイン固有のエラー
│   │       └── errors.go
│   │
│   ├── usecase/                            # ユースケース層（アプリケーションロジック）
│   │   ├── article_usecase.go
│   │   ├── article_usecase_test.go
│   │   ├── article_generator_usecase.go
│   │   ├── article_generator_usecase_test.go
│   │   ├── tag_usecase.go
│   │   ├── tag_usecase_test.go
│   │   ├── book_recommendation_usecase.go
│   │   └── book_recommendation_usecase_test.go
│   │
│   ├── interface/                          # インターフェース層（入出力）
│   │   └── handler/                        # HTTPハンドラー
│   │       ├── article_handler.go
│   │       ├── article_handler_test.go
│   │       ├── article_generator_handler.go
│   │       ├── article_generator_handler_test.go
│   │       ├── tag_handler.go
│   │       ├── tag_handler_test.go
│   │       ├── book_recommendation_handler.go
│   │       ├── book_recommendation_handler_test.go
│   │       └── error_handler.go
│   │
│   └── infrastructure/                     # インフラ層（外部システム統合）
│       ├── database/                       # データベース関連
│       │   ├── mysql.go                    # MySQL接続
│       │   ├── migration.go                # マイグレーション管理
│       │   └── migrations/                 # マイグレーションSQLファイル
│       │       ├── 000001_create_articles_table.up.sql
│       │       ├── 000001_create_articles_table.down.sql
│       │       ├── 000002_create_tags_table.up.sql
│       │       ├── 000002_create_tags_table.down.sql
│       │       ├── 000003_create_article_tags_table.up.sql
│       │       ├── 000003_create_article_tags_table.down.sql
│       │       ├── 000004_create_book_recommendations_table.up.sql
│       │       └── 000004_create_book_recommendations_table.down.sql
│       ├── repository/                     # リポジトリ実装
│       │   ├── mysql_article_repository.go
│       │   ├── mysql_article_repository_test.go
│       │   ├── mysql_tag_repository.go
│       │   ├── mysql_tag_repository_test.go
│       │   ├── mysql_book_recommendation_repository.go
│       │   ├── memory_article_repository.go  # インメモリ実装（テスト用）
│       │   └── memory_tag_repository.go
│       ├── ai/                             # AI統合
│       │   ├── gemini_client.go
│       │   └── gemini_client_test.go
│       ├── external/                       # 外部API統合
│       │   ├── google_books_client.go
│       │   └── google_books_client_test.go
│       ├── logger/                         # ロガー
│       │   └── logger.go
│       ├── service/                        # インフラサービス実装
│       │   └── book_recommendation_service_impl.go
│       └── timeutil/                       # 時刻関連ユーティリティ
│           ├── formatter.go
│           └── formatter_test.go
│
├── tmp/                                    # Airホットリロード一時ファイル（Gitignore対象）
│   └── main                                # ホットリロード用ビルドバイナリ（自動生成）
│
├── .air.toml                               # Airホットリロード設定
├── .gitignore                              # APIディレクトリGit除外設定
├── Dockerfile                              # Dockerコンテナ定義
├── go.mod                                  # Go依存関係管理
└── go.sum                                  # Go依存関係ロックファイル
```

### 3.2 各層の役割

| 層 | 役割 | 依存方向 | 配置場所 |
|----|------|---------|---------|
| **cmd/** | アプリケーションエントリーポイント | すべての層に依存可能 | `cmd/server/main.go` |
| **domain/** | ビジネスルール、エンティティ定義 | 他の層に依存しない | `internal/domain/` |
| **usecase/** | アプリケーションロジック | domainに依存 | `internal/usecase/` |
| **interface/** | 入出力処理（HTTPハンドラー） | usecase, domainに依存 | `internal/interface/` |
| **infrastructure/** | 外部システム統合 | domain, usecaseに依存 | `internal/infrastructure/` |

**依存関係の方向**:
```
cmd/server/main.go
    ↓
interface/handler → usecase → domain
    ↓                           ↑
infrastructure/ ────────────────┘
```

### 3.3 重要なファイル

| ファイル | 役割 |
|---------|------|
| `cmd/server/main.go` | アプリケーションのエントリーポイント、依存性注入、ルーティング設定、サーバー起動 |
| `internal/domain/entity/*.go` | ドメインエンティティ定義（Article, Tag, BookRecommendation） |
| `internal/domain/repository/*.go` | リポジトリインターフェース定義 |
| `internal/usecase/*.go` | ビジネスロジック実装 |
| `internal/interface/handler/*.go` | HTTPハンドラー実装 |
| `internal/infrastructure/repository/*.go` | リポジトリ実装（MySQL, インメモリ） |

---

## 4. フロントエンド（frontend/）

### 4.1 ディレクトリ構造

フロントエンドはNext.js App Routerベースの構造を採用しています。

```
frontend/
├── app/                                    # Next.js App Routerページ
│   ├── layout.tsx                          # ルートレイアウト
│   ├── page.tsx                            # ホームページ（/）
│   ├── globals.css                         # グローバルCSS（Tailwind含む）
│   ├── favicon.ico                         # ファビコン
│   ├── articles/                           # 記事関連ページ
│   │   ├── page.tsx                        # 記事一覧（/articles）
│   │   ├── new/
│   │   │   └── page.tsx                    # 記事作成（/articles/new）
│   │   └── [id]/                           # 動的ルート
│   │       ├── page.tsx                    # 記事詳細（/articles/:id）
│   │       └── edit/
│   │           └── page.tsx                # 記事編集（/articles/:id/edit）
│   └── tags/                               # タグ関連ページ
│       └── page.tsx                        # タグ管理（/tags）
│
├── components/                             # 再利用可能なUIコンポーネント
│   ├── ArticleCard.tsx                     # 記事カードコンポーネント
│   ├── ArticleForm.tsx                     # 記事作成フォーム
│   ├── ArticleForm.test.tsx                # ArticleFormのテスト
│   ├── ArticleEditForm.tsx                 # 記事編集フォーム
│   ├── ArticleEditForm.test.tsx            # ArticleEditFormのテスト
│   ├── SearchBar.tsx                       # 検索バー
│   ├── Tag.tsx                             # タグ表示コンポーネント
│   ├── TagList.tsx                         # タグリスト
│   ├── BookRecommendations.tsx             # 書籍推薦コンポーネント
│   ├── BookRecommendations.test.tsx        # BookRecommendationsのテスト
│   ├── Header.tsx                          # ヘッダー
│   ├── Sidebar.tsx                         # サイドバー
│   ├── Toast.tsx                           # トースト通知
│   ├── ErrorBoundary.tsx                   # エラーバウンダリー
│   └── DeleteConfirmDialog.tsx             # 削除確認ダイアログ
│
├── hooks/                                  # カスタムReact Hooks
│   ├── useArticles.ts                      # 記事CRUD操作フック
│   ├── useArticles.test.ts                 # useArticlesのテスト
│   ├── useArticleSearch.ts                 # 記事検索フック
│   ├── useArticleSearch.test.ts            # useArticleSearchのテスト
│   ├── useTags.ts                          # タグ操作フック
│   ├── useTags.test.ts                     # useTagsのテスト
│   ├── useBookRecommendations.ts           # 書籍推薦フック
│   └── useBookRecommendations.test.ts      # useBookRecommendationsのテスト
│
├── lib/                                    # ライブラリ・ユーティリティ
│   ├── api/                                # APIクライアント
│   │   ├── baseClient.ts                   # 基底APIクライアント
│   │   ├── articleClient.ts                # 記事APIクライアント
│   │   ├── articleClient.test.ts           # articleClientのテスト
│   │   ├── tagClient.ts                    # タグAPIクライアント
│   │   ├── tagClient.test.ts               # tagClientのテスト
│   │   └── bookRecommendationClient.ts     # 書籍推薦APIクライアント
│   ├── errors/                             # エラークラス
│   │   └── ApiError.ts                     # APIエラークラス
│   └── utils/                              # ユーティリティ関数
│       ├── sample.ts
│       └── sample.test.ts
│
├── types/                                  # TypeScript型定義
│   ├── article.ts                          # 記事型定義
│   ├── tag.ts                              # タグ型定義
│   └── book.ts                             # 書籍型定義
│
├── contexts/                               # React Context
│   └── ToastContext.tsx                    # トースト通知Context
│
├── config/                                 # 設定ファイル
│   └── constants.ts                        # 定数定義（API URL、バリデーション）
│
├── public/                                 # 静的ファイル
│   ├── next.svg
│   ├── vercel.svg
│   └── *.svg                               # その他SVGアイコン
│
├── coverage/                               # テストカバレッジレポート（Gitignore対象）
│
├── ARCHITECTURE.md                         # フロントエンドアーキテクチャ詳細
├── Dockerfile                              # Dockerコンテナ定義
├── README.md                               # フロントエンドREADME
├── package.json                            # npm依存関係
├── package-lock.json                       # npm依存関係ロックファイル
├── next.config.ts                          # Next.js設定
├── tsconfig.json                           # TypeScript設定
├── tsconfig.tsbuildinfo                    # TypeScriptビルド情報（Gitignore対象）
├── vitest.config.ts                        # Vitestテスト設定
├── vitest.setup.ts                         # Vitestセットアップ
├── eslint.config.mjs                       # ESLint設定
├── postcss.config.mjs                      # PostCSS設定
└── next-env.d.ts                           # Next.js型定義（自動生成）
```

### 4.2 各ディレクトリの役割

| ディレクトリ | 役割 | 配置ルール |
|------------|------|-----------|
| **app/** | Next.js App Routerページ | ルーティングに対応したディレクトリ構造 |
| **components/** | 再利用可能なUIコンポーネント | コンポーネント名.tsx、テストは同じディレクトリに配置 |
| **hooks/** | カスタムReact Hooks | `use`プレフィックス、テストは同じディレクトリに配置 |
| **lib/api/** | APIクライアント | バックエンドAPIとの通信ロジック |
| **lib/errors/** | エラークラス | カスタムエラークラス定義 |
| **lib/utils/** | ユーティリティ関数 | 汎用的なヘルパー関数 |
| **types/** | TypeScript型定義 | ドメインオブジェクトの型定義 |
| **contexts/** | React Context | グローバル状態管理 |
| **config/** | 設定ファイル | 定数、環境変数関連 |
| **public/** | 静的ファイル | 画像、アイコン、フォント |

### 4.3 重要なファイル

| ファイル | 役割 |
|---------|------|
| `app/layout.tsx` | アプリケーション全体のレイアウト、ToastProvider設定 |
| `app/page.tsx` | ホームページ（ダッシュボード） |
| `config/constants.ts` | API URL、バリデーション定数、アプリ設定 |
| `lib/api/baseClient.ts` | 基底APIクライアント、共通エラーハンドリング |
| `contexts/ToastContext.tsx` | トースト通知のグローバル状態管理 |

---

## 5. ドキュメント（docs/）

### 5.1 ディレクトリ構造

```
docs/
├── product-requirements.md         # プロダクト要求定義書
├── functional-design.md            # 機能設計書
├── architecture.md                 # 技術仕様書
├── repository-structure.md         # リポジトリ構造定義書（本ドキュメント）
├── development-guidelines.md       # 開発ガイドライン（作成予定）
├── glossary.md                     # ユビキタス言語定義（作成予定）
└── images/                         # ドキュメント用画像（必要に応じて作成）
    └── *.png, *.svg
```

### 5.2 各ドキュメントの役割

| ドキュメント | 役割 | 更新タイミング |
|------------|------|--------------|
| **product-requirements.md** | プロダクトビジョン、機能要件、ユーザーストーリー | 大きな機能追加時のみ |
| **functional-design.md** | API設計、データモデル、UI設計、画面遷移 | 設計変更時のみ |
| **architecture.md** | 技術スタック、システム構成、デプロイ・運用設計 | 技術選定変更、インフラ変更時のみ |
| **repository-structure.md** | ディレクトリ構造、ファイル配置ルール | ディレクトリ構造変更時のみ |
| **development-guidelines.md** | コーディング規約、命名規則、Git規約 | ルール追加・変更時のみ |
| **glossary.md** | ドメイン用語、ビジネス用語の定義 | 新規用語追加時 |

### 5.3 データベース永続化ディレクトリ（db/）

```
db/                                     # データベース永続化（Gitignore対象）
├── data/                               # 開発環境MySQLデータ永続化（Docker Volume）
│   ├── article_manager/                # アプリケーションデータベース
│   ├── mysql/                          # MySQLシステムデータベース
│   ├── performance_schema/             # パフォーマンススキーマ
│   └── sys/                            # MySQL sysスキーマ
└── test_data/                          # テスト環境MySQLデータ永続化（Docker Volume）
    ├── article_manager_test/           # テストデータベース
    ├── mysql/                          # MySQLシステムデータベース
    ├── performance_schema/             # パフォーマンススキーマ
    └── sys/                            # MySQL sysスキーマ
```

**役割**:
- Docker ComposeのMySQLコンテナのデータを永続化
- コンテナを再起動してもデータが保持される
- `docker-compose.yml`の`volumes`設定で自動マウント

**重要事項**:
- `db/`ディレクトリ全体が**Gitignore対象**（データベースファイルはバージョン管理しない）
- 開発環境（`db/data/`）とテスト環境（`db/test_data/`）で分離
- データの手動バックアップが必要（詳細は`docs/architecture.md`の6.6節参照）
- データベースのクリーンアップ: `docker-compose down -v`でボリュームごと削除可能

---

## 6. 設定ファイル・環境変数

### 6.1 ルートディレクトリの設定ファイル

| ファイル | 役割 | 配置場所 |
|---------|------|---------|
| `.env` | 環境変数定義（開発・本番）、**Gitignore対象** | ルート |
| `.env.example` | 環境変数テンプレート（作成予定） | ルート |
| `.gitignore` | Git除外設定 | ルート |
| `.claude/` | Claude Code設定・カスタムエージェント、**Gitignore対象** | ルート |
| `.mcp.json` | MCPサーバー設定 | ルート |
| `CLAUDE.md` | プロジェクトメモリ（Claude Code用） | ルート |
| `README.md` | プロジェクト概要、セットアップ手順 | ルート |
| `docker-compose.yml` | 開発環境Docker Compose設定 | ルート |
| `docker-compose-test.yml` | テスト環境Docker Compose設定 | ルート |

### 6.2 バックエンド設定ファイル（api/）

| ファイル | 役割 | 配置場所 |
|---------|------|---------|
| `.air.toml` | Airホットリロード設定（ビルドコマンド、監視パス等） | `api/.air.toml` |
| `.gitignore` | APIディレクトリGit除外設定（tmp/等） | `api/.gitignore` |
| `go.mod` | Go依存関係管理 | `api/go.mod` |
| `go.sum` | Go依存関係ロックファイル | `api/go.sum` |
| `Dockerfile` | バックエンドコンテナ定義 | `api/Dockerfile` |

### 6.3 フロントエンド設定ファイル（frontend/）

| ファイル | 役割 | 配置場所 |
|---------|------|---------|
| `next.config.ts` | Next.js設定 | `frontend/next.config.ts` |
| `tsconfig.json` | TypeScript設定 | `frontend/tsconfig.json` |
| `vitest.config.ts` | Vitestテスト設定 | `frontend/vitest.config.ts` |
| `vitest.setup.ts` | Vitestセットアップ | `frontend/vitest.setup.ts` |
| `eslint.config.mjs` | ESLint設定 | `frontend/eslint.config.mjs` |
| `postcss.config.mjs` | PostCSS設定 | `frontend/postcss.config.mjs` |
| `package.json` | npm依存関係 | `frontend/package.json` |
| `package-lock.json` | npm依存関係ロックファイル | `frontend/package-lock.json` |
| `Dockerfile` | フロントエンドコンテナ定義 | `frontend/Dockerfile` |

### 6.4 環境変数（.env）

```bash
# データベース
MYSQL_ROOT_PASSWORD=<root_password>
MYSQL_DATABASE=article_manager
MYSQL_USER=<db_user>
MYSQL_PASSWORD=<db_password>
DB_PORT=3306

# バックエンドAPI
API_PORT=8080

# フロントエンド
FRONTEND_PORT=3000
NEXT_PUBLIC_API_URL=http://localhost:8080

# 外部API
GEMINI_API_KEY=<your_gemini_api_key>
GOOGLE_BOOKS_API_KEY=<your_books_api_key>
```

**重要**: `.env`ファイルは**Gitignore対象**です。機密情報を含むため、リポジトリにコミットしないでください。

### 6.5 .claudeディレクトリの構成

```
.claude/
├── settings.local.json         # Claude Codeローカル設定（Gitignore対象）
└── agents/                     # カスタムエージェント定義
    └── *.md                    # エージェント定義ファイル
```

**用途**: Claude Code（AI開発支援ツール）の設定とカスタムエージェント定義を格納します。

**Gitignore対象理由**: 開発者ごとの設定が含まれるため、バージョン管理から除外します。

---

## 7. ファイル配置ルール

### 7.1 基本原則

1. **関心の分離**: 機能ごとにファイルを分割
2. **近接配置**: 関連ファイルは同じディレクトリに配置
3. **テスト併置**: テストファイルは実装ファイルと同じディレクトリに配置
4. **命名一貫性**: ディレクトリ名、ファイル名は統一的な命名規則に従う

### 7.2 バックエンドのファイル配置

| 配置場所 | ファイルの種類 | 例 |
|---------|--------------|-----|
| `cmd/server/` | アプリケーションエントリーポイント | `main.go` |
| `internal/domain/entity/` | エンティティ、値オブジェクト | `article.go`, `tag.go` |
| `internal/domain/repository/` | リポジトリインターフェース | `article_repository.go` |
| `internal/domain/service/` | ドメインサービスインターフェース | `ai_generator.go` |
| `internal/domain/errors/` | ドメイン固有のエラー | `errors.go` |
| `internal/usecase/` | ユースケース実装 | `article_usecase.go` |
| `internal/interface/handler/` | HTTPハンドラー | `article_handler.go` |
| `internal/infrastructure/repository/` | リポジトリ実装 | `mysql_article_repository.go` |
| `internal/infrastructure/ai/` | AI統合 | `gemini_client.go` |
| `internal/infrastructure/external/` | 外部API統合 | `google_books_client.go` |
| `internal/infrastructure/database/` | DB接続、マイグレーション | `mysql.go`, `migration.go` |
| `internal/infrastructure/logger/` | ロガー | `logger.go` |

**ファイル分割基準**:
- 1エンティティ = 1ファイル
- 1ユースケース = 1ファイル
- 1ハンドラー = 1ファイル
- 複数の関連機能をまとめない（単一責任の原則）

### 7.3 フロントエンドのファイル配置

| 配置場所 | ファイルの種類 | 例 |
|---------|--------------|-----|
| `app/` | ページコンポーネント | `page.tsx`, `layout.tsx` |
| `components/` | 再利用可能なUIコンポーネント | `ArticleCard.tsx`, `SearchBar.tsx` |
| `hooks/` | カスタムフック | `useArticles.ts`, `useTags.ts` |
| `lib/api/` | APIクライアント | `articleClient.ts`, `tagClient.ts` |
| `lib/errors/` | エラークラス | `ApiError.ts` |
| `lib/utils/` | ユーティリティ関数 | `formatDate.ts`, `validateUrl.ts` |
| `types/` | TypeScript型定義 | `article.ts`, `tag.ts` |
| `contexts/` | React Context | `ToastContext.tsx` |
| `config/` | 設定ファイル | `constants.ts` |

**ファイル分割基準**:
- 1コンポーネント = 1ファイル
- 1カスタムフック = 1ファイル
- 1エンティティの型定義 = 1ファイル
- 関連する定数は`config/constants.ts`にまとめる

### 7.4 ドキュメントのファイル配置

| 配置場所 | ファイルの種類 | 例 |
|---------|--------------|-----|
| `docs/` | 永続的ドキュメント | `architecture.md`, `functional-design.md` |
| `docs/images/` | ドキュメント用画像（必要に応じて） | `er-diagram.png`, `wireframe.svg` |
| `.steering/[YYYYMMDD]-[title]/` | 作業単位のドキュメント | `requirements.md`, `design.md`, `tasklist.md` |

**ドキュメント配置ルール**:
- 永続的ドキュメントは`docs/`に配置
- 作業単位のドキュメントは`.steering/`に配置
- 図表はMermaid記法で記述（可能な限り）
- 複雑な図表のみ画像ファイルとして`docs/images/`に配置

---

## 8. 命名規則

### 8.1 ディレクトリ命名規則

| 種類 | 命名規則 | 例 |
|-----|---------|-----|
| **パッケージディレクトリ（Go）** | 小文字、単数形、短い名前 | `entity`, `repository`, `usecase` |
| **Next.jsページディレクトリ** | 小文字、複数形（リソース名） | `articles`, `tags` |
| **コンポーネントディレクトリ** | 小文字、複数形 | `components`, `contexts` |
| **ステアリングディレクトリ** | `[YYYYMMDD]-[開発タイトル]` | `20250115-add-tag-feature` |

### 8.2 ファイル命名規則

#### 8.2.1 バックエンド（Go）

| ファイルの種類 | 命名規則 | 例 |
|-------------|---------|-----|
| **エンティティ** | スネークケース、単数形 | `article.go`, `tag.go` |
| **リポジトリインターフェース** | `<entity>_repository.go` | `article_repository.go` |
| **リポジトリ実装** | `<storage>_<entity>_repository.go` | `mysql_article_repository.go` |
| **ユースケース** | `<entity>_usecase.go` | `article_usecase.go` |
| **ハンドラー** | `<entity>_handler.go` | `article_handler.go` |
| **テストファイル** | `<filename>_test.go` | `article_test.go` |

#### 8.2.2 フロントエンド（TypeScript）

| ファイルの種類 | 命名規則 | 例 |
|-------------|---------|-----|
| **Reactコンポーネント** | パスカルケース | `ArticleCard.tsx`, `SearchBar.tsx` |
| **カスタムフック** | キャメルケース、`use`プレフィックス | `useArticles.ts`, `useTags.ts` |
| **型定義** | スネークケース、単数形 | `article.ts`, `tag.ts` |
| **APIクライアント** | キャメルケース、`Client`サフィックス | `articleClient.ts`, `tagClient.ts` |
| **ユーティリティ** | キャメルケース | `formatDate.ts`, `validateUrl.ts` |
| **ページ** | `page.tsx`, `layout.tsx` | `app/articles/page.tsx` |
| **テストファイル** | `<filename>.test.ts(x)` | `ArticleCard.test.tsx` |

#### 8.2.3 ドキュメント

| ファイルの種類 | 命名規則 | 例 |
|-------------|---------|-----|
| **永続的ドキュメント** | ケバブケース、`.md` | `product-requirements.md`, `architecture.md` |
| **ステアリングドキュメント** | 小文字、`.md` | `requirements.md`, `design.md`, `tasklist.md` |

### 8.3 変数・関数命名規則

#### 8.3.1 Go

| 種類 | 命名規則 | 例 |
|-----|---------|-----|
| **変数** | キャメルケース | `articleID`, `userName` |
| **定数** | パスカルケースまたはアッパースネークケース | `MaxRetries`, `DEFAULT_TIMEOUT` |
| **関数** | パスカルケース（公開）、キャメルケース（非公開） | `CreateArticle`, `validateInput` |
| **インターフェース** | パスカルケース | `ArticleRepository`, `AIGenerator` |
| **構造体** | パスカルケース | `Article`, `Tag` |

#### 8.3.2 TypeScript

| 種類 | 命名規則 | 例 |
|-----|---------|-----|
| **変数** | キャメルケース | `articleId`, `userName` |
| **定数** | アッパースネークケース | `API_BASE_URL`, `MAX_TAG_LENGTH` |
| **関数** | キャメルケース | `createArticle`, `validateUrl` |
| **型・インターフェース** | パスカルケース | `Article`, `Tag`, `ArticleFormData` |
| **React Component** | パスカルケース | `ArticleCard`, `SearchBar` |

---

## 9. テストファイル配置規則

### 9.1 基本原則

- **テストファイルは実装ファイルと同じディレクトリに配置**
- **テストファイル名は実装ファイル名 + `_test` サフィックス**
- **1実装ファイル = 1テストファイル**

### 9.2 バックエンド（Go）

#### 配置ルール

```
internal/
├── domain/
│   └── entity/
│       ├── article.go              # 実装
│       └── article_test.go         # テスト
├── usecase/
│   ├── article_usecase.go          # 実装
│   └── article_usecase_test.go     # テスト
└── interface/
    └── handler/
        ├── article_handler.go      # 実装
        └── article_handler_test.go # テスト
```

#### テストファイル命名

| 実装ファイル | テストファイル |
|------------|--------------|
| `article.go` | `article_test.go` |
| `article_usecase.go` | `article_usecase_test.go` |
| `mysql_article_repository.go` | `mysql_article_repository_test.go` |

#### テスト実行コマンド

```bash
# すべてのテストを実行
go test ./...

# 特定のパッケージのテストを実行
go test ./internal/usecase/

# カバレッジ付きでテストを実行
go test -cover ./...
```

### 9.3 フロントエンド（TypeScript）

#### 配置ルール

```
components/
├── ArticleCard.tsx             # 実装
└── ArticleCard.test.tsx        # テスト

hooks/
├── useArticles.ts              # 実装
└── useArticles.test.ts         # テスト

lib/
└── api/
    ├── articleClient.ts        # 実装
    └── articleClient.test.ts   # テスト
```

#### テストファイル命名

| 実装ファイル | テストファイル |
|------------|--------------|
| `ArticleCard.tsx` | `ArticleCard.test.tsx` |
| `useArticles.ts` | `useArticles.test.ts` |
| `articleClient.ts` | `articleClient.test.ts` |

#### テスト実行コマンド

```bash
# すべてのテストを実行
npm run test

# UIモードでテストを実行
npm run test:ui

# カバレッジ付きでテストを実行
npm run test:coverage
```

### 9.4 統合テスト・E2Eテスト（Phase 2予定）

統合テストやE2Eテストは、専用のディレクトリに配置します。

```
api/
└── test/
    ├── integration/            # 統合テスト
    │   ├── article_test.go
    │   └── tag_test.go
    └── e2e/                    # E2Eテスト（Phase 2）
        └── ...

frontend/
└── test/
    ├── integration/            # 統合テスト（Phase 2）
    │   └── ...
    └── e2e/                    # E2Eテスト（Phase 2）
        └── ...
```

---

## 付録A: ファイル作成チェックリスト

新しいファイルを作成する際のチェックリスト：

### バックエンド（Go）

- [ ] ファイルは適切な層（domain, usecase, interface, infrastructure）に配置されているか
- [ ] ファイル名は命名規則に従っているか（スネークケース）
- [ ] パッケージ名はディレクトリ名と一致しているか
- [ ] 公開関数・構造体はドキュメントコメントが付いているか
- [ ] テストファイル（`*_test.go`）を同じディレクトリに作成したか
- [ ] 依存関係の方向は正しいか（外側から内側への依存）

### フロントエンド（TypeScript）

- [ ] ファイルは適切なディレクトリ（components, hooks, lib等）に配置されているか
- [ ] ファイル名は命名規則に従っているか（パスカルケース/キャメルケース）
- [ ] 型定義は`types/`に配置されているか
- [ ] テストファイル（`*.test.ts(x)`）を同じディレクトリに作成したか
- [ ] コンポーネントは責務が単一か（Single Responsibility Principle）

### ドキュメント

- [ ] 永続的ドキュメントは`docs/`に配置されているか
- [ ] 作業単位のドキュメントは`.steering/[YYYYMMDD]-[title]/`に配置されているか
- [ ] Mermaid記法で図表を記述しているか（可能な場合）
- [ ] ドキュメントの目次が更新されているか

---

| 日付 | バージョン | 変更内容 | 担当者 |
|-----|-----------|---------|-------|
| 2026-02-13 | 1.0 | 初版作成 | - |
| 2026-02-13 | 1.1 | 最優先修正: `.air.toml`、`.claude/`、`api/.gitignore`追加、マイグレーションファイル4つに修正、設定ファイルセクション拡充 | - |

---

**注意**: このドキュメントは永続的ドキュメントとして管理されます。リポジトリ構造の大きな変更時のみ更新してください。
