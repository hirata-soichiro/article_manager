# 開発ガイドライン（軽量版）

**バージョン:** 3.0.0 (個人開発向け軽量版)
**作成日:** 2026-02-14
**更新日:** 2026-02-14
**対象:** 個人開発者（自分自身）

---

## 目次

1. [概要](#1-概要)
2. [開発コマンド集](#2-開発コマンド集)
3. [コーディング規約](#3-コーディング規約)
4. [命名規則](#4-命名規則)
5. [テスト規約](#5-テスト規約)
6. [Git規約](#6-git規約)
7. [セキュリティチェックリスト](#7-セキュリティチェックリスト)
8. [トラブルシューティング](#8-トラブルシューティング)

---

## 1. 概要

### 1.1 このドキュメントについて

個人開発における最低限の開発ルールを定義します。**「ルール」ではなく「ガイド」** として活用してください。

### 1.2 基本原則

1. **可読性優先**: 将来の自分が理解できるコードを書く
2. **シンプルさ**: 過度に複雑な実装は避ける
3. **一貫性**: プロジェクト全体で統一されたスタイルを維持
4. **テスタビリティ**: テスト可能な設計を心がける

---

## 2. 開発コマンド集

### 2.1 バックエンド（Go）

```bash
# ホットリロードで実行
cd api && air -c .air.toml

# テスト実行
cd api && go test ./...
cd api && go test -v ./...              # 詳細表示
cd api && go test -cover ./...          # カバレッジ付き

# リント
cd api && golangci-lint run
cd api && golangci-lint run --fix       # 自動修正

# ビルド
cd api && go build -o server cmd/server/main.go
```

### 2.2 フロントエンド（Next.js）

```bash
cd frontend

# 開発サーバー起動
npm run dev

# テスト実行
npm run test              # すべてのテストを実行
npm run test:ui           # UIでテストを実行
npm run test:coverage     # カバレッジレポート付き

# リント・型チェック
npm run lint              # ESLint
npx tsc --noEmit          # TypeScript型チェック

# ビルド
npm run build
npm run start             # 本番サーバー起動
```

### 2.3 Docker

```bash
# すべてのサービスを起動
docker-compose up -d

# ログを表示
docker-compose logs -f [api|frontend|db]

# すべてのサービスを停止
docker-compose down

# コンテナを再ビルド
docker-compose up -d --build

# データベースに接続
docker exec -it db mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}
```

### 2.4 テストデータベース

```bash
# 1. テストデータベースを起動
docker-compose -f docker-compose-test.yml up -d

# 2. マイグレーション実行（初回のみ）
docker exec -i {container_id} mysql -utest_user -ptest_password article_manager_test < api/internal/infrastructure/database/migrations/000001_create_articles_table.up.sql

# 3. リポジトリテスト実行
cd api && go test -v ./internal/infrastructure/repository/

# 4. テストデータベースを停止
docker-compose -f docker-compose-test.yml down
```

---

## 3. コーディング規約

### 3.1 バックエンド（Go）

#### 基本方針
- [Effective Go](https://go.dev/doc/effective_go)に従う
- `gofmt`でフォーマット（保存時に自動実行推奨）
- `golangci-lint`でリント（コミット前に実行）

#### ファイルサイズ基準
- コメント・空行を除いた実装コードが**100行を超えたら分割を検討**
- 分割後は各ファイル**50-80行**を目標

#### Clean Architectureの遵守
依存関係の方向:
```
cmd/server → interface/handler → usecase → domain ← infrastructure
```

**重要**: domain層は他に依存しない

#### エラーハンドリング
```go
// ✅ Good: エラーをラップして詳細情報を追加
func (uc *ArticleUseCase) CreateArticle(ctx context.Context, input *ArticleInput) (*entity.Article, error) {
    article := &entity.Article{
        Title:   input.Title,
        URL:     input.URL,
        Summary: input.Summary,
    }

    if err := uc.repo.Create(ctx, article); err != nil {
        return nil, fmt.Errorf("failed to create article: %w", err)
    }

    return article, nil
}
```

#### インポート順序
標準ライブラリ → 外部ライブラリ → 内部パッケージ（空行で区切る）

```go
import (
    "context"
    "fmt"

    "github.com/jmoiron/sqlx"

    "article_manager/internal/domain/entity"
)
```

#### コメント規約
- 公開関数・型には必ずコメントを付ける
- 複雑なロジックには説明コメント（「なぜ」を書く）

---

### 3.2 フロントエンド（TypeScript/React）

#### 基本方針
- TypeScript Strict モードを使用
- ESLint + Prettier でフォーマット
- Next.js App Router の規約に従う

#### ファイルサイズ基準
- コンポーネント実装が**100行を超えたら分割を検討**
- 状態管理ロジックはカスタムフックに切り出す

#### 型定義
```typescript
// ✅ Good: 明確な型定義
interface ArticleCardProps {
  article: Article;
  onDelete: (id: number) => Promise<void>;
  className?: string;
}

// ❌ Bad: any の多用
interface ArticleCardProps {
  article: any;  // ❌
  onDelete: any;  // ❌
}
```

#### コンポーネント設計
```typescript
// ✅ Good: カスタムフックで状態管理を分離
export function ArticleList() {
  const { articles, loading, error, deleteArticle } = useArticles();

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage error={error} />;

  return (
    <div className="article-list">
      {articles.map(article => (
        <ArticleCard key={article.id} article={article} onDelete={deleteArticle} />
      ))}
    </div>
  );
}
```

#### Tailwind CSS
- Tailwind CSS を優先的に使用
- カスタムCSSは必要最小限に（`::-webkit-scrollbar`等、Tailwindで表現不可能なもののみ）
- クラス名の順序: レイアウト → 配置 → サイズ → スペーシング → 背景・ボーダー → テキスト → エフェクト

---

## 4. 命名規則

### 4.1 バックエンド（Go）

| 種類 | 規則 | 例 |
|------|------|-----|
| ファイル | スネークケース | `article.go`, `article_repository.go` |
| パッケージ | 小文字、単数形 | `entity`, `usecase`, `handler` |
| 公開変数・関数 | パスカルケース | `CreateArticle`, `ArticleID` |
| 非公開変数・関数 | キャメルケース | `validateInput`, `articleID` |
| 定数 | パスカルケース | `MaxRetries`, `DefaultTimeout` |
| テストファイル | `_test.go` サフィックス | `article_test.go` |

### 4.2 フロントエンド（TypeScript）

| 種類 | 規則 | 例 |
|------|------|-----|
| コンポーネント | パスカルケース | `ArticleCard.tsx` |
| カスタムフック | `use`プレフィックス + キャメルケース | `useArticles.ts` |
| ユーティリティ | キャメルケース | `formatDate.ts` |
| 型定義ファイル | スネークケース | `article.ts` |
| 変数 | キャメルケース | `articleId`, `userName` |
| 定数 | アッパースネークケース | `API_BASE_URL` |
| 型・インターフェース | パスカルケース | `Article`, `ArticleFormData` |
| イベントハンドラ | `handle`プレフィックス | `handleClick`, `handleSubmit` |
| Boolean変数 | `is/has/should`プレフィックス | `isLoading`, `hasError` |
| テストファイル | `.test.ts(x)` サフィックス | `ArticleCard.test.tsx` |

---

## 5. テスト規約

### 5.1 カバレッジ目標

- **バックエンド全体: 70%以上**
- domain/entity: **90%以上**（ビジネスロジックのため）
- usecase: **80%以上**
- handler: **60%以上**

### 5.2 バックエンドテスト（Go）

#### テーブル駆動テスト
```go
func TestArticle_Validate(t *testing.T) {
    tests := []struct {
        name    string
        article entity.Article
        wantErr bool
    }{
        {
            name: "有効な記事",
            article: entity.Article{
                Title:   "テスト記事",
                URL:     "https://example.com",
            },
            wantErr: false,
        },
        {
            name: "タイトルが空",
            article: entity.Article{
                Title: "",
                URL:   "https://example.com",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.article.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 5.3 フロントエンドテスト（TypeScript）

#### コンポーネントテスト
```typescript
import { render, screen, fireEvent } from '@testing-library/react';

describe('ArticleCard', () => {
  it('記事情報が正しく表示される', () => {
    const mockArticle = {
      id: 1,
      title: 'テスト記事',
      url: 'https://example.com',
    };

    render(<ArticleCard article={mockArticle} onDelete={jest.fn()} />);

    expect(screen.getByText('テスト記事')).toBeInTheDocument();
  });
});
```

#### カスタムフックテスト
```typescript
import { renderHook, waitFor } from '@testing-library/react';

describe('useArticles', () => {
  beforeEach(() => {
    global.fetch = vi.fn();
  });

  it('記事一覧を取得できる', async () => {
    const mockArticles = [
      { id: 1, title: 'Test 1' },
      { id: 2, title: 'Test 2' },
    ];

    (global.fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => mockArticles,
    });

    const { result } = renderHook(() => useArticles());

    await waitFor(() => {
      expect(result.current.articles).toEqual(mockArticles);
    });
  });
});
```

---

## 6. Git規約

### 6.1 ブランチ命名規則

```
main              ← 本番環境
feature/xxx       ← 新機能開発
fix/xxx           ← バグ修正
hotfix/xxx        ← 緊急修正
refactor/xxx      ← リファクタリング
```

**例**:
```bash
feature/add-tag-hierarchy
fix/search-performance
hotfix/security-vulnerability
```

### 6.2 コミットメッセージ

**フォーマット**:
```
<type>: <subject>

<body>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

**Type の種類**:
| Type | 説明 | 例 |
|------|------|-----|
| `feat` | 新機能追加 | `feat: add tag hierarchy feature` |
| `fix` | バグ修正 | `fix: resolve search performance issue` |
| `docs` | ドキュメント変更 | `docs: update API documentation` |
| `style` | コードフォーマット | `style: format code with prettier` |
| `refactor` | リファクタリング | `refactor: extract validation logic` |
| `test` | テスト追加・修正 | `test: add unit tests for ArticleUseCase` |
| `chore` | ビルド・設定変更 | `chore: update dependencies` |

**良いコミットメッセージの例**:
```bash
feat: add AI-powered article generation from URL

Implemented Gemini API integration to automatically generate
article title, summary, and tags from a given URL.

- Add GeminiClient in infrastructure layer
- Add ArticleGeneratorUseCase
- Add POST /articles/generate endpoint

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

### 6.3 コミット前チェックリスト

個人開発でも以下を確認してからコミット:

- [ ] テストがすべてパス（`go test ./...` / `npm run test`）
- [ ] リントエラーなし（`golangci-lint run` / `npm run lint`）
- [ ] TypeScript型エラーなし（`npx tsc --noEmit`）
- [ ] 動作確認済み（手動テスト）
- [ ] コミットメッセージが明確

### 6.4 プルリクエスト（PR）

**PRテンプレート**:
```markdown
## 概要
この変更の目的を簡潔に説明してください。

## 変更内容
- [ ] 機能A を追加
- [ ] バグB を修正

## テスト
- [ ] ユニットテスト追加済み
- [ ] 手動テスト実施済み

## チェックリスト
- [ ] リント・型チェックをパス
- [ ] テストカバレッジが低下していない
```

---

## 7. セキュリティチェックリスト

### 7.1 バックエンド

- [ ] **SQL Injection対策**: プレースホルダーを使用（`sqlx`の`?`）
- [ ] **XSS対策**: HTMLエスケープ（フロントエンドで実施）
- [ ] **認証・認可**: 必要に応じてミドルウェアで実装
- [ ] **機密情報**: 環境変数で管理（コードにハードコードしない）
- [ ] **エラーメッセージ**: 詳細なスタックトレースをクライアントに返さない
- [ ] **CORS設定**: 本番環境では適切に制限

### 7.2 フロントエンド

- [ ] **XSS対策**: React はデフォルトでエスケープ（`dangerouslySetInnerHTML`は避ける）
- [ ] **API Key**: 環境変数で管理（`.env.local`、Gitにコミットしない）
- [ ] **入力バリデーション**: クライアント側とサーバー側の両方で実施
- [ ] **依存関係**: 定期的に`npm audit`で脆弱性チェック

### 7.3 環境変数管理

**必須ルール**:
- `.env`ファイルは`.gitignore`に追加
- `.env.example`をリポジトリにコミット（ダミー値のみ）
- 本番環境では別途環境変数を設定

---

## 8. トラブルシューティング

### 8.1 バックエンド

#### マイグレーションエラー

```bash
# マイグレーション状態の確認
docker-compose logs api | grep -i migration

# データベーススキーマを確認
docker exec -it db mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}
mysql> SHOW TABLES;
mysql> DESCRIBE articles;
```

#### テストデータベース接続エラー

```bash
# コンテナが起動しているか確認
docker ps | grep mysql_test

# マイグレーションが実行されているか確認
docker exec -it mysql_test mysql -utest_user -ptest_password article_manager_test -e "SHOW TABLES;"
```

#### Air が起動しない

```bash
# .air.toml の設定を確認
cat api/.air.toml

# Air を再インストール
go install github.com/cosmtrek/air@latest

# 手動で起動
cd api && go run cmd/server/main.go
```

---

### 8.2 フロントエンド

#### Next.js が起動しない

```bash
# node_modules を削除して再インストール
cd frontend
rm -rf node_modules package-lock.json
npm install

# ポートが使用中か確認
lsof -i :3000
kill -9 <PID>

# キャッシュをクリア
rm -rf .next
npm run dev
```

#### API 接続エラー

```bash
# API_BASE_URL を確認
cat frontend/.env.local

# バックエンドが起動しているか確認
curl http://localhost:8080/api/articles

# CORS エラーの場合: バックエンドの CORS 設定を確認
```

#### テストが失敗する

```bash
# キャッシュをクリア
npm run test -- --clearCache

# 特定のテストファイルのみ実行
npm run test -- ArticleCard.test.tsx

# デバッグモード
npm run test -- --no-coverage --verbose
```

---

### 8.3 Docker

#### コンテナが起動しない

```bash
# ログを確認
docker-compose logs -f [api|frontend|db]

# コンテナを削除して再作成
docker-compose down -v
docker-compose up -d --build

# ポートが使用中か確認
lsof -i :8080  # API
lsof -i :3000  # Frontend
lsof -i :3306  # MySQL
```

#### データベースに接続できない

```bash
# MySQL コンテナが起動しているか確認
docker ps | grep db

# 環境変数を確認
cat .env

# 手動で接続テスト
docker exec -it db mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}
```

---

## 付録: よく使うコマンド

### A. 開発フロー（毎日）

```bash
# 1. 最新コードを取得
git pull origin main

# 2. 開発環境起動
docker-compose up -d

# 3. 開発サーバー起動（別ターミナル）
cd api && air -c .air.toml
cd frontend && npm run dev

# 4. 変更後: テスト実行
cd api && go test ./...
cd frontend && npm run test

# 5. コミット
git add .
git commit -m "feat: add new feature

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
"
git push origin feature/xxx
```

### B. リント・型チェック（コミット前）

```bash
# バックエンド
cd api && golangci-lint run --fix
cd api && go test ./...

# フロントエンド
cd frontend && npm run lint
cd frontend && npx tsc --noEmit
cd frontend && npm run test
```

### C. データベース操作

```bash
# マイグレーション確認
docker-compose logs api | grep -i migration

# データベースに接続
docker exec -it db mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}

# テーブル確認
mysql> SHOW TABLES;
mysql> DESCRIBE articles;

# データ確認
mysql> SELECT * FROM articles;
```

---

## 変更履歴

| バージョン | 日付 | 変更内容 |
|-----------|------|---------|
| 3.0.0 | 2026-02-14 | 個人開発向け軽量版に全面刷新（2,765行→約700行） |
| 2.0.0 | 2026-02-14 | 詳細版（チーム開発向け、2,765行） |
| 1.0.0 | 2026-01-XX | 初版 |

---

**注意**: このドキュメントは「ガイド」であり、厳格な「ルール」ではありません。状況に応じて柔軟に対応してください。
