# 設計書

## アーキテクチャ概要

Docker ComposeベースのHTTPサーバーをAWS Lambda + RDSのサーバーレスアーキテクチャに移行します。

**重要**: ECS Fargate移行計画で既に構築済みのVPC、RDS、ECR等のインフラを活用します。

## 移行前後のアーキテクチャ

### 移行前（現在）

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTP
       ↓
┌──────────────────────────┐
│  Docker Compose          │
│  ┌────────────────────┐  │
│  │  Frontend          │  │
│  │  (Next.js:3000)    │  │
│  └──────┬─────────────┘  │
│         │ HTTP           │
│  ┌──────↓─────────────┐  │
│  │  Backend API       │  │
│  │  (Go:8080)         │  │
│  └──────┬─────────────┘  │
│         │ TCP            │
│  ┌──────↓─────────────┐  │
│  │  MySQL:3306        │  │
│  └────────────────────┘  │
└──────────────────────────┘
```

### 移行後（Lambda + 既存インフラ）

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTPS
       ↓
┌───────────────────────────────────────────────┐
│  ECS Fargate (Frontend)                      │
│  - Next.js                                   │
│  - Public IP                                 │
└───────────────┬───────────────────────────────┘
                │ HTTPS
                ↓
┌───────────────────────────────────────────────┐
│  AWS Lambda Function URLs                    │
│  https://<id>.lambda-url.ap-northeast-1...   │
│  ┌────────────────────────────────────────┐  │
│  │  Lambda Function (Go)                  │  │
│  │  - Memory: 512MB                       │  │
│  │  - Timeout: 30s                        │  │
│  │  - Runtime: provided.al2023            │  │
│  │  ┌──────────────────────────────────┐  │  │
│  │  │  aws-lambda-go-api-proxy         │  │  │
│  │  └────────┬─────────────────────────┘  │  │
│  │           │                            │  │
│  │  ┌────────↓─────────────────────────┐  │  │
│  │  │  Clean Architecture              │  │  │
│  │  │  (handler/usecase/domain/infra)  │  │  │
│  │  └────────┬─────────────────────────┘  │  │
│  └───────────┼─────────────────────────────┘  │
└──────────────┼─────────────────────────────────┘
               │ TCP (VPC内)
               ↓
┌───────────────────────────────────────────────┐
│  Amazon RDS for MySQL（既存）                │
│  - Instance: db.t4g.micro (Single-AZ)        │
│  - Private Subnet: 10.0.11.0/24              │
│  - Security Group: Lambda only               │
└───────────────────────────────────────────────┘
```

### ネットワーク構成（既存インフラ）

```
┌────────────────────────────────────────────────┐
│  VPC (10.0.0.0/16)（既存）                     │
│                                                │
│  ┌──────────────────────────────────────────┐ │
│  │  パブリックサブネット (10.0.1.0/24)      │ │
│  │  ap-northeast-1a                         │ │
│  │                                          │ │
│  │  - Internet Gateway                      │ │
│  │  - Lambda ENI（新規配置）                │ │
│  └──────────────────────────────────────────┘ │
│                                                │
│  ┌──────────────────────────────────────────┐ │
│  │  プライベートサブネット (10.0.11.0/24)   │ │
│  │  ap-northeast-1a                         │ │
│  │                                          │ │
│  │  ┌────────────────────────────────────┐ │ │
│  │  │  RDS MySQL（既存）                 │ │ │
│  │  │  - db.t4g.micro                    │ │ │
│  │  │  - Single-AZ                       │ │ │
│  │  │  - Security Group: Lambda only     │ │ │
│  │  └────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────┘ │
└────────────────────────────────────────────────┘
```

### 外部API統合

```
┌────────────────────┐
│  Lambda Function   │
└─────────┬──────────┘
          │ HTTPS
          ├──────────→ Google Gemini API
          │            (記事生成・書籍推薦)
          │
          └──────────→ Google Books API
                       (書籍情報検索)
```

## コンポーネント設計

### 1. 既存インフラ（ECS移行で構築済み）

#### 1.1 VPC・ネットワーク

**責務**:
- 論理的に分離されたネットワーク空間の提供
- パブリックサブネットとプライベートサブネットの分離

**実装済み（ECS移行Phase 2）**:
- **VPC CIDR**: `10.0.0.0/16`
- **パブリックサブネット**:
  - `10.0.1.0/24` (ap-northeast-1a)
  - 用途: Lambda ENI配置
- **プライベートサブネット**:
  - `10.0.11.0/24` (ap-northeast-1a)
  - 用途: RDS配置
- **インターネットゲートウェイ**: パブリックサブネット用
- **Single-AZ構成**: ap-northeast-1a のみ使用

**Lambda移行での変更点**:
- Lambda用のENI（Elastic Network Interface）をパブリックサブネットに配置
- セキュリティグループをLambda → RDS接続用に調整

#### 1.2 RDS MySQL

**責務**:
- リレーショナルデータベースのマネージドサービス提供
- 自動バックアップ（最小限）

**実装済み（ECS移行Phase 3）**:
- **インスタンスクラス**: `db.t4g.micro` (2 vCPU, 1 GB RAM)
- **ストレージ**: General Purpose SSD (gp3), 20 GB
- **エンジン**: MySQL 8.0
- **Multi-AZ**: 無効（Single-AZ、コスト削減）
- **AZ**: ap-northeast-1a
- **自動バックアップ**: 有効、保持期間1日
- **データベース名**: `article_manager`
- **マスターユーザー名**: `admin`
- **セキュリティグループ**: プライベートサブネット内

**Lambda移行での変更点**:
- セキュリティグループのインバウンドルールをLambdaからのアクセスに変更
- Parameter StoreでDB接続情報を管理

#### 1.3 ECR (Elastic Container Registry)

**実装済み（ECS移行Phase 3）**:
- **リポジトリ**:
  - `article-manager-frontend` (Next.jsイメージ)
  - `article-manager-api` (Go APIイメージ)
- **イメージタグ戦略**: `latest`, `{git-sha}`
- **ライフサイクルポリシー**: 未使用イメージを30日後に自動削除

**Lambda移行での使用**:
- Lambda用ECRリポジトリを追加（`article-manager-lambda-api`）
- フロントエンドリポジトリは継続使用（ECS Fargate）
- 既存のバックエンドリポジトリは削除または無効化（Lambda Container Imageに移行）

#### 1.4 CloudWatch Logs

**実装済み（ECS移行Phase 5）**:
- **ロググループ**: `/ecs/article-manager-app`
- **保持期間**: 1日（コスト最小化）

**Lambda移行での変更点**:
- Lambda用のロググループを追加: `/aws/lambda/article-manager-api`
- 既存のロググループは削除可能（ECS不使用のため）

#### 1.5 IAMロール

**実装済み（ECS移行Phase 5）**:
- ECSタスク実行ロール（Lambda移行では不要）
- CloudWatch Logs書き込み権限

**Lambda移行での変更点**:
- Lambda実行ロール（新規作成）
- Parameter Store読み取り権限（追加）
- VPC実行ポリシー（追加）

### 2. Lambda固有コンポーネント（新規実装）

#### 2.1 Lambda Entrypoint (`cmd/lambda/main.go`)

**責務**:
- Lambda Handlerの初期化
- DB接続プールのグローバル管理
- aws-lambda-go-api-proxyによるHTTPルーティング変換

#### 2.2 Lambda Container Image (`api/Dockerfile.lambda`)

**責務**:
- Lambda実行環境用のDockerイメージをビルド
- マルチステージビルドで最適化

#### 2.3 DB接続管理 (`infrastructure/database/mysql.go`)

**責務**:
- Lambda環境でのDB接続プール設定
- タイムアウト設定の最適化

#### 2.4 Parameter Store統合 (`internal/config/parameter_store.go`)

**責務**:
- AWS Systems Manager Parameter Storeから環境変数を取得
- DB接続情報、APIキーの管理

#### 2.5 Gemini API統合（同期処理）

**責務**:
- 記事生成（`POST /api/articles/generate`）
- 書籍推薦（`GET /api/book-recommendations`）
- Lambdaタイムアウト: 30秒

### 3. Terraform Infrastructure（更新）

#### 3.1 既存Terraformファイル（維持）

**そのまま使用**:
- `terraform/vpc.tf` - VPC、サブネット、セキュリティグループ（既存）
- `terraform/rds.tf` - RDSインスタンス（既存）
- `terraform/cloudwatch.tf` - CloudWatch Logs（既存）

**調整が必要**:
- `terraform/ecr.tf` - Lambda用ECRリポジトリを追加
- `terraform/security_groups.tf` - Lambda用のインバウンドルール追加
- `terraform/iam.tf` - Lambda実行ロール追加
- `terraform/parameters.tf` - Parameter Store設定を更新

#### 3.2 新規Terraformファイル

**新規作成**:
- `terraform/lambda.tf` - Lambda関数定義
- `terraform/ecs.tf` - フロントエンド用ECS設定

**削除**:
- `terraform/service_discovery.tf` - 不要（Lambda移行）

#### 3.3 ECRリポジトリ追加（terraform/ecr.tf）

Lambda用ECRリポジトリ（`article-manager-lambda-api`）を追加。イメージスキャン有効、ライフサイクルポリシーで最新10イメージを保持。

#### 3.4 Lambda関数定義（terraform/lambda.tf）

Lambda関数（`article-manager-api`）とFunction URLs設定。Container Image方式、メモリ512MB、タイムアウト30秒、CORS設定。

#### 3.5 Parameter Store定義（terraform/parameters.tf）

DB接続情報（host, port, user, password, name）とAPIキー（Gemini, Google Books）をParameter Storeに定義。


