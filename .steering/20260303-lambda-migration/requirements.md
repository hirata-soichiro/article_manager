# 要求内容

## 概要

Docker ComposeベースのローカルGo APIサーバーからAWS Lambda + Function URLsへの移行を実施し、個人開発に最適なサーバーレスアーキテクチャを構築します。

**重要**: 既にPhase 1-4（VPC、RDS、ECR、IAMロール等の基盤構築）はECS Fargate移行計画で完了済みです。このLambda移行では、既存インフラを活用しつつ、Lambda固有の実装を追加します。

## 背景

### ECS移行計画の実施状況

2026年2月に開始したECS Fargate移行計画では、以下のインフラが既に構築されています：

#### ✅ 完了済みのインフラ（ECS移行 Phase 1-4）

- **Terraform基盤**: S3バックエンド、VPC、サブネット、セキュリティグループ
- **RDS MySQL**: db.t4g.micro（Single-AZ）、プライベートサブネット配置
- **ECR**: フロントエンド・バックエンド用リポジトリ
- **IAM**: ECSタスク実行ロール、CloudWatch Logs書き込み権限
- **CloudWatch Logs**: ロググループ（保持期間1日）
- **Service Discovery**: AWS Cloud Map プライベートDNSネームスペース

#### 🔄 未完了（Phase 5以降）

- ECS Fargateタスク定義・サービス作成
- Route 53設定
- 起動/停止ワークフロー
- CI/CDパイプライン

### ECS移行案の課題

ECS Fargateでのインフラ構築を進めた結果、以下の課題が明らかになりました：

- **常時起動コスト**: 月額$22-25（使わない時も課金）
- **起動/停止の手間**: GitHub Actions経由で毎回手動起動が必要
- **インフラ管理の複雑さ**: Fargate、ALB、Route 53の設定が必要
- **オーバーエンジニアリング**: 個人開発アプリには過剰な構成

### Lambda移行のメリット

- **ほぼ無料**: AWS無料枠（月100万リクエスト）で運用可能
- **自動スケーリング**: アイドル時は完全にコストゼロ
- **シンプル**: インフラ管理が不要、Function URLsで直接アクセス
- **学習価値**: サーバーレスアーキテクチャの実践的な学習
- **既存インフラの活用**: VPC、RDS、ECR等は再利用可能

### アプリケーションの使用パターン

- 週に数回、1日100リクエスト程度
- 月間トータル約3,000リクエスト
- 同時アクセスは1-2人程度
- トラフィックが安定していない（波がある）

→ **Lambdaの無料枠内で十分に運用可能**

## 実装対象の機能

### 1. 既存インフラの活用と調整

**既に構築済み（再利用）**:
- VPC（10.0.0.0/16）
- パブリックサブネット（10.0.1.0/24）
- プライベートサブネット（10.0.11.0/24）
- RDS MySQL（db.t4g.micro、Single-AZ）
- ECR（フロントエンド・バックエンド用リポジトリ）
- CloudWatch Logs

**Lambda用に調整が必要**:
- セキュリティグループ（Lambda → RDS接続を許可）
- IAMロール（Lambda実行ロール、Parameter Store読み取り権限）
- VPC設定（Lambda用のサブネット、ENI配置）

### 2. Go APIのLambda対応への変更

- Lambda Handlerの実装（`cmd/lambda/main.go`）
- 既存のHTTPルーティングを維持
- Gemini API処理は同期処理（タイムアウト30秒）

### 3. RDS MySQL（既存インフラを継続使用）

**既に構築済み**:
- インスタンスタイプ: db.t4g.micro（無料枠対象）
- ストレージ: 20GB gp3
- 自動バックアップ: 1日保持
- プライベートサブネット配置

**Lambda移行での変更点**:
- セキュリティグループをLambdaからのアクセスに変更
- 接続元をECSからLambdaに変更

**データ移行**:
- 新規環境として構築（クリーンな状態から開始）
- マイグレーションスクリプトは初回Lambda起動時に自動実行

### 4. AWS Systems Manager Parameter Store（新規構築）

**責務**:
- DB接続情報の管理
- APIキーの安全な保管

**設定内容**:
- `/article-manager/db/host` - RDSエンドポイント
- `/article-manager/db/port` - 3306
- `/article-manager/db/user` - admin
- `/article-manager/db/password` - パスワード（SecureString）
- `/article-manager/db/name` - article_manager
- `/article-manager/gemini/api-key` - Gemini API Key（SecureString）
- `/article-manager/google-books/api-key` - Google Books API Key（SecureString）

### 5. Lambda Function URLsによる直接HTTPアクセス（新規構築）

**Function URLsの設定**:
- 認証: NONE（個人開発のため、後でIP制限を検討）
- CORS設定: フロントエンドドメインを許可
- HTTPメソッド: GET, POST, PUT, DELETE

**URL構成**:
- API: `https://<function-id>.lambda-url.ap-northeast-1.on.aws/api/*`
- フロントエンドから直接アクセス

### 6. Infrastructure as Code（Terraform）の更新

**既存Terraformの活用**:
- `terraform/vpc.tf` - そのまま使用
- `terraform/rds.tf` - そのまま使用
- `terraform/ecr.tf` - Lambda用ECRリポジトリを追加
- `terraform/security_groups.tf` - Lambda用に調整
- `terraform/iam.tf` - Lambda実行ロール追加
- `terraform/cloudwatch.tf` - そのまま使用
- `terraform/parameters.tf` - Parameter Store設定を更新

**新規追加**:
- `terraform/lambda.tf` - Lambda関数定義
- `terraform/ecs.tf` - フロントエンド用ECS設定（新規作成）

**削除**:
- `terraform/service_discovery.tf` - 不要（Lambda移行）
- `terraform/route53.tf` - 不要（Phase 2以降で検討）

### 7. CI/CDパイプライン構築

**GitHub Actionsワークフロー（Docker Image方式）**:
- Lambda用Dockerイメージをビルド・ECRにpush
- Lambda関数が自動的に新しいイメージをpull
- トリガー: `main`ブランチへのpush

### 8. フロントエンドの更新

**API URLの変更**:
- `config/constants.ts`の`API_BASE_URL`をLambda Function URLに変更
- 環境変数`NEXT_PUBLIC_API_URL`で管理

**デプロイ先**:
- ECS Fargate（フロントエンド用ECS設定を新規作成）
- CORS設定: ECS FargateのパブリックIPからLambda Function URLへのアクセスを許可

## 受け入れ条件

### Lambda関数

- [ ] Go APIがLambda Handlerとして正しく動作する
- [ ] 既存のHTTPルーティング（GET/POST/PUT/DELETE）が維持される
- [ ] Function URLsからHTTPアクセスが可能
- [ ] Cold Start時間が5秒以内
- [ ] 環境変数がAWS Systems Manager Parameter Storeから取得できる
- [ ] RDSへのデータベース接続が成功する
- [ ] マイグレーションが初回起動時に実行される

### 既存インフラの活用

- [x] VPC、サブネット、セキュリティグループが既に存在する（ECS移行で完了済み）
- [x] RDSインスタンス（Single-AZ、db.t4g.micro）が既に存在する（ECS移行で完了済み）
- [x] ECRリポジトリが既に存在する（ECS移行で完了済み）
- [ ] セキュリティグループがLambda → RDS接続を許可する
- [ ] IAMロールがLambda実行に必要な権限を持つ

### Infrastructure as Code

- [x] Terraformバックエンド（S3）が設定されている（ECS移行で完了済み）
- [ ] Terraformコードが正常にapply可能（Lambda関数追加）
- [ ] Lambda関数が作成される
- [ ] IAMロールとポリシーが正しく設定される（Lambda実行ロール追加）
- [ ] AWS Systems Manager Parameter Storeが設定される

### CI/CD

- [ ] GitHub Actionsでビルド・デプロイが自動実行される
- [ ] `main`ブランチpush時にLambda関数が更新される
- [ ] デプロイ後、アプリケーションが正常に動作する

### フロントエンド

- [ ] Lambda Function URLへのAPI呼び出しが成功する
- [ ] CORS設定が正しく動作する
- [ ] 既存の機能（記事CRUD、タグ管理、AI生成、書籍推薦）が動作する

## 成功指標

### 技術的指標

- **Lambda実行時間**: 平均300ms以内（Cold Start除く）
- **Cold Start時間**: 5秒以内
- **RDS接続**: 100ms以内
- **API応答時間**: 500ms以内（90パーセンタイル）

### コスト指標

#### 月額コスト（常時利用可能な状態）

| サービス | 構成 | 月額コスト |
|---------|------|-----------|
| Lambda | 無料枠内（3,000リクエスト） | $0.00 |
| RDS MySQL | db.t4g.micro Single-AZ（既存） | $12.41 |
| RDS ストレージ | 20 GB gp3（既存） | $2.53 |
| AWS Systems Manager | Parameter Store（無料枠） | $0.00 |
| CloudWatch Logs | 1日保持（既存） | $0.10 |
| **合計** | | **約$15/月** |

#### 実質コスト（週20時間稼働の場合）

| サービス | 稼働時間 | 月額コスト |
|---------|---------|-----------|
| Lambda | リクエストベース | $0.00 |
| RDS MySQL | 80時間/月 | $1.36 |
| RDS ストレージ | 常時 | $2.53 |
| その他 | 常時 | $0.10 |
| **合計** | | **約$4/月** |

#### ECSとの比較

| 項目 | ECS Fargate | Lambda |
|-----|------------|--------|
| 月額コスト（常時） | $22-25 | $15 |
| 月額コスト（週20h） | $5-6 | $4 |
| 起動/停止の手間 | GitHub Actions必要 | 不要（自動） |
| Cold Start | なし | 5秒程度 |
| インフラ管理 | 複雑（ALB, Route 53等） | シンプル（Function URLs） |
| 学習価値 | コンテナオーケストレーション | サーバーレスアーキテクチャ |

**結論**: 個人開発では**Lambda**が圧倒的に優れている

### 運用改善指標

- **デプロイ時間**: 5分以内（GitHub Actions経由）
- **インフラ変更の再現性**: Terraformによる100%再現可能
- **起動/停止の手間**: 不要（Lambdaは自動スケール）

## スコープ外

以下はこのフェーズでは実装しません:

### Phase 2以降に延期

- **API Gateway統合**: Function URLsで十分、カスタムドメインはPhase 2
- **Lambda@Edge**: CDN統合は後回し
- **DynamoDBへの移行**: RDSで開始、必要に応じて検討
- **ElastiCache (Redis)**: キャッシュ層の導入は後回し
- **CloudWatch Alarmsの詳細設定**: 基本的なログのみ
- **VPC Endpointsの設定**: コスト削減のため後回し
- **Lambda Layers**: 依存関係の共有は将来検討
- **Step Functions統合**: 複雑なワークフローは現状不要
- **Route 53**: カスタムドメインは後で検討

### 運用最適化 (後回し)

- **Lambda Provisioned Concurrency**: Cold Start対策（コストが高い）
- **RDS Proxy**: 接続プーリングの最適化（コストが高い）
- **CloudWatch Logs Insights**: 詳細なログ分析は後回し
- **X-Rayによる分散トレーシング**: 個人学習では不要

### セキュリティ強化 (Phase 2)

- **Lambda Function URLsの認証**: IP制限、Cognito統合は後で検討
- **AWS WAF**: Function URLs前段のファイアウォールは後回し
- **Secrets Manager**: Parameter Storeで開始、必要に応じて移行

### ECS移行で削除する部分

以下はECS移行計画で実装予定だったが、Lambda移行では不要:
- **ECS Fargate**: サーバーレスで代替
- **Application Load Balancer (ALB)**: Function URLsで代替
- **Route 53動的DNS更新**: 固定URLで十分
- **起動/停止ワークフロー**: 自動スケールで不要

## 参照ドキュメント

### プロジェクト内ドキュメント

- `docs/product-requirements.md` - プロダクト要求定義書
- `docs/functional-design.md` - 機能設計書
- `docs/architecture.md` - アーキテクチャ設計書（ローカル環境）
- `CLAUDE.md` - プロジェクトメモリ
- `.steering/20260217-aws-infrastructure-migration/` - ECS移行計画（Phase 1-4完了）

### AWSドキュメント

- [AWS Lambda Go SDK](https://github.com/aws/aws-lambda-go)
- [AWS Lambda Go API Proxy](https://github.com/awslabs/aws-lambda-go-api-proxy)
- [Lambda Function URLs](https://docs.aws.amazon.com/lambda/latest/dg/lambda-urls.html)
- [Amazon RDS for MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)

### 移行ガイド

- Go HTTP Server → Lambda Handler変換パターン
- ECS Fargate → Lambda移行手順（既存インフラ活用）
- Lambda Cold Start最適化のベストプラクティス
- AWS Systems Manager Parameter Storeの使用方法
