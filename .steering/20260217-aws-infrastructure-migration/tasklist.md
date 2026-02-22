# タスクリスト

## 🚨 タスク完全完了の原則

**このファイルの全タスクが完了するまで作業を継続すること**

### 必須ルール
- **全てのタスクを`[x]`にすること**
- 「時間の都合により別タスクとして実施予定」は禁止
- 「実装が複雑すぎるため後回し」は禁止
- 未完了タスク（`[ ]`）を残したまま作業を終了しない

### 実装可能なタスクのみを計画
- 計画段階で「実装可能なタスク」のみをリストアップ
- 「将来やるかもしれないタスク」は含めない
- 「検討中のタスク」は含めない

### タスクスキップが許可される唯一のケース
以下の技術的理由に該当する場合のみスキップ可能:
- 実装方針の変更により、機能自体が不要になった
- アーキテクチャ変更により、別の実装方法に置き換わった
- 依存関係の変更により、タスクが実行不可能になった

スキップ時は必ず理由を明記:
```markdown
- [x] ~~タスク名~~（実装方針変更により不要: 具体的な技術的理由）
```

### タスクが大きすぎる場合
- タスクを小さなサブタスクに分割
- 分割したサブタスクをこのファイルに追加
- サブタスクを1つずつ完了させる

---

## Phase 1: 準備・Terraform基盤構築 (Week 1, Day 1-3)

### 1.1 プロジェクト構造の準備

- [x] `terraform/` ディレクトリを作成
- [x] `.github/workflows/` ディレクトリを作成
- [x] `scripts/` ディレクトリを作成
- [x] `.gitignore` に Terraform関連ファイルを追加
  - [x] `terraform/.terraform/` を追加
  - [x] `terraform/*.tfstate` を追加
  - [x] `terraform/*.tfstate.backup` を追加
  - [x] `terraform/.terraform.lock.hcl` を除外しない (バージョン固定のため)

### 1.2 OIDC Identity Provider作成（手動）

- [x] AWS Management Consoleにログイン（ルートユーザーまたは管理者権限のあるIAMユーザー）
- [x] IAM → Identity providers → Add providerを選択
  - [x] プロバイダータイプ: `OpenID Connect`
  - [x] プロバイダーURL: `https://token.actions.githubusercontent.com`
  - [x] Audience: `sts.amazonaws.com`
  - [x] "Add provider"をクリック
- [x] 作成されたプロバイダーのARNをメモ（後で使用）

### 1.3 GitHub Actions用IAMロール作成（手動）

- [x] IAM → Roles → Create roleを選択
- [x] Trusted entity type: `Web identity`を選択
  - [x] Identity provider: 先ほど作成したOIDCプロバイダーを選択
  - [x] Audience: `sts.amazonaws.com`
  - [x] GitHub organization: 自分のGitHubユーザー名または組織名
  - [x] GitHub repository: `article_manager`（リポジトリ名）
  - [x] GitHub branch: `main`（または`feature/infra`）
- [x] 権限設定（以下のいずれかを選択）
  - [x] オプションA（簡単・学習用）: `AdministratorAccess`ポリシーをアタッチ
  - [x] オプションB（推奨・本番用）: カスタムポリシーで最小権限を設定
    - [x] 必要な権限: EC2, VPC, RDS, ECS, ECR, S3, Secrets Manager, Route53, IAM（ロール作成用）, CloudWatch Logs, Systems Manager
- [x] ロール名: `github-actions-terraform-role`
- [x] ロールARNをメモ（GitHub Secretsに登録する）

### 1.4 GitHub Secrets設定（Phase 1: OIDC用）

- [x] GitHubリポジトリのSettings → Secrets and variables → Actionsにアクセス
- [x] 以下のSecretsを登録:
  - [x] `AWS_ROLE_ARN` - タスク1.3で作成したIAMロールのARN
  - [x] `AWS_REGION` - `ap-northeast-1`

### 1.5 Terraformバックエンド設定（S3ネイティブロック使用）

- [x] S3バケットを作成 (terraform.tfstate保存用)
  - [x] AWS Management Consoleまたは一時的にローカルでAWS CLIを使用
  - [x] バケット名: `article-manager-terraform-state`
  - [x] リージョン: `ap-northeast-1`
  - [x] バージョニングを有効化
  - [x] 暗号化を有効化
- [x] `terraform/backend.tf` を作成
  - [x] S3バケット名を定義
  - [x] `use_lockfile = true` を設定（S3ネイティブステートロック）

### 1.6 Terraform基本設定ファイル

- [x] `terraform/provider.tf` を作成
  - [x] AWS Providerバージョンを `~> 5.0` に設定
  - [x] リージョンを `ap-northeast-1` に設定
- [x] `terraform/variables.tf` を作成
  - [x] プロジェクト名 (`project_name` = "article-manager")
  - [x] 環境名 (`environment` = "production")
  - [x] AWS リージョン (`aws_region`)
  - [x] VPC CIDR (`vpc_cidr` = "10.0.0.0/16")
  - [x] ドメイン名 (`domain_name` - Phase 6で有効化)

### 1.7 Terraform実行用GitHub Actionsワークフロー作成

- [x] `.github/workflows/terraform-apply.yml` を作成
  - [x] トリガー: `workflow_dispatch`（手動実行）
  - [x] ジョブ: Terraform実行
    - [x] `permissions: id-token: write, contents: read` を設定（OIDC用）
    - [x] チェックアウト
    - [x] AWS認証情報を設定（OIDC方式）
      - [x] `aws-actions/configure-aws-credentials@v6` を使用（v4から最新版に更新）
      - [x] `role-to-assume: ${{ secrets.AWS_ROLE_ARN }}` を指定
      - [x] `aws-region: ${{ secrets.AWS_REGION }}` を指定
    - [x] Terraformをセットアップ（hashicorp/setup-terraform@v3）
    - [x] `terraform fmt -check` を追加（コードスタイルチェック）
    - [x] `terraform init`
    - [x] `terraform validate`
    - [x] `terraform plan`
    - [x] `terraform apply -auto-approve`（planが成功した場合）
- [ ] ワークフローファイルをコミット・プッシュ

---

## Phase 2: ネットワーク・セキュリティ構築 (Week 1, Day 4-5)

### 2.1 VPCネットワーク構築（Single-AZ、シンプル構成）

- [x] `terraform/vpc.tf` を作成
  - [x] VPCリソースを定義 (CIDR: `10.0.0.0/16`)
  - [x] パブリックサブネット (ap-northeast-1a)
    - [x] `10.0.1.0/24`
  - [x] プライベートサブネット (ap-northeast-1a)
    - [x] `10.0.11.0/24`
  - [x] インターネットゲートウェイを定義
  - [x] ルートテーブルを定義
    - [x] パブリックサブネット用 (0.0.0.0/0 → Internet Gateway)
    - [x] プライベートサブネット用 (デフォルトのローカルルートのみ)

### 2.2 セキュリティグループ定義

- [x] `terraform/security-groups.tf` を作成
  - [x] ECSタスク用セキュリティグループ
    - [x] インバウンド: HTTP (80) from 0.0.0.0/0
    - [x] インバウンド: ポート3000 from 0.0.0.0/0 (Frontend直接アクセス用、オプション)
    - [x] インバウンド: ポート8080 from 0.0.0.0/0 (Backend直接アクセス用、オプション)
    - [x] アウトバウンド: All traffic
  - [x] RDS用セキュリティグループ
    - [x] インバウンド: MySQL (3306) from ECS SG
    - [x] アウトバウンド: なし

### 2.3 Terraform初期化・検証（GitHub Actions経由）

- [x] GitHub Actionsで `terraform-apply.yml` ワークフローを実行
- [x] ワークフローログで以下を確認:
  - [x] `terraform init` が成功
  - [x] `terraform validate` が成功
  - [x] `terraform plan` で実行計画を確認
- [x] VPC、サブネット、セキュリティグループが作成されることを確認
  - [x] 必要に応じて段階的に適用（terraform targetを使用）

---

## Phase 3: RDS・ECR・Secrets Manager構築 (Week 1, Day 6-7)

### 3.0 Parameter Store手動設定（事前準備）

- [x] AWS Management ConsoleまたはAWS CLIでParameter Storeにアクセス
- [x] 機密情報（SecureString型）を作成
  - [x] `/article-manager/db/admin-password` - `.env`の`MYSQL_ROOT_PASSWORD`
  - [x] `/article-manager/db/app-password` - `.env`の`MYSQL_PASSWORD`
  - [x] `/article-manager/api/gemini-api-key` - `.env`の`GEMINI_API_KEY`
  - [x] `/article-manager/api/google-books-api-key` - `.env`の`GOOGLE_BOOKS_API_KEY`

### 3.1 Terraform - Parameter Store参照設定

- [x] `terraform/parameters.tf` を作成
  - [x] 機密情報（SecureString）をdata sourceで参照
    - [x] `/article-manager/db/admin-password`
    - [x] `/article-manager/db/app-password`
    - [x] `/article-manager/api/gemini-api-key`
    - [x] `/article-manager/api/google-books-api-key`
  - [x] 非機密情報をTerraformで作成
    - [x] `/article-manager/db/name` = "article_manager"
    - [x] `/article-manager/db/admin-user` = "admin"
    - [x] `/article-manager/db/app-user` = "article_user"

### 3.2 RDS Parameter Group作成

- [x] `terraform/rds.tf` にパラメータグループを追加
  - [x] `character_set_server` = `utf8mb4`
  - [x] `collation_server` = `utf8mb4_unicode_ci`
  - [x] `innodb_ft_min_token_size` = `2`
  - [x] `ft_min_word_len` = `2`
  - [x] 注: `[client]`と`[mysql]`セクションの文字セット設定はRDSパラメータグループでは不可。アプリケーション接続文字列で`charset=utf8mb4`を指定する

### 3.3 RDS for MySQL構築

- [x] `terraform/rds.tf` を作成（または既存ファイルに追加）
  - [x] DBサブネットグループを定義 (プライベートサブネット)
  - [x] RDSインスタンスを定義
    - [x] インスタンスクラス: `db.t4g.micro`
    - [x] エンジン: `mysql` バージョン `8.0`
    - [x] ストレージ: `gp3`, 20 GB
    - [x] Multi-AZ: `false`
    - [x] 自動バックアップ: `false` (保持期間0日)
    - [x] データベース名: `article_manager`
    - [x] マスターユーザー名: `admin`
    - [x] パスワード: Parameter Storeから取得 (`aws_ssm_parameter.db_admin_password.value`)
    - [x] セキュリティグループ: RDS用SG
    - [x] パブリックアクセス: `false`
    - [x] 暗号化: `true` (KMSデフォルトキー)
    - [x] パラメータグループ: カスタムパラメータグループを使用
- [x] GitHub Actionsで `terraform-apply.yml` ワークフローを実行
- [x] ワークフローログで RDS作成を確認
- [x] AWS Management ConsoleでRDSエンドポイントを確認

### 3.4 ECR (Elastic Container Registry) 構築

- [x] `terraform/ecr.tf` を作成
  - [x] フロントエンド用ECRリポジトリを定義
    - [x] リポジトリ名: `article-manager-frontend`
    - [x] イメージタグの可変性: `MUTABLE`
    - [x] スキャン設定: オンプッシュスキャン有効化
  - [x] バックエンド用ECRリポジトリを定義
    - [x] リポジトリ名: `article-manager-api`
    - [x] イメージタグの可変性: `MUTABLE`
    - [x] スキャン設定: オンプッシュスキャン有効化
  - [x] ライフサイクルポリシーを定義
    - [x] 未使用イメージを30日後に削除
    - [x] `latest` タグは削除しない
- [x] GitHub Actionsで `terraform-apply.yml` ワークフローを実行
- [x] AWS Management ConsoleでECRリポジトリURLを確認

---

## Phase 4: Dockerイメージ最適化・ECRプッシュ (Week 2, Day 1-2)

### 4.1 バックエンド (Go API) Dockerfile最適化

- [ ] `api/Dockerfile` をマルチステージビルドに変更
  - [ ] ビルドステージ: `golang:1.25-alpine`
  - [ ] 実行ステージ: `alpine:latest`
  - [ ] ポート8080を公開
- [ ] ヘルスチェックエンドポイント `/api/health` を実装
  - [ ] `api/internal/interface/handler/health_handler.go` を作成
  - [ ] HTTPステータス200を返すシンプルなハンドラー
  - [ ] `main.go` でルートを登録
- [ ] ローカルでビルドテスト
  - [ ] `cd api && docker build -t article-manager-api:local .`
  - [ ] `docker run -p 8080:8080 article-manager-api:local`
  - [ ] `curl http://localhost:8080/api/health` で動作確認

### 4.2 フロントエンド (Next.js) Dockerfile最適化

- [ ] `frontend/Dockerfile` をマルチステージビルドに変更
  - [ ] ビルドステージ: `node:20-alpine`
  - [ ] 実行ステージ: `node:20-alpine`
  - [ ] ポート3000を公開
- [ ] ローカルでビルドテスト
  - [ ] `cd frontend && docker build -t article-manager-frontend:local .`
  - [ ] `docker run -p 3000:3000 article-manager-frontend:local`
  - [ ] ブラウザで `http://localhost:3000` にアクセスして動作確認

### 4.3 ECRへのイメージプッシュ

- [ ] AWS ECRにログイン
  - [ ] `aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com`
- [ ] バックエンドイメージをビルド・タグ付け・プッシュ
  - [ ] `cd api && docker build -t article-manager-api:latest .`
  - [ ] `docker tag article-manager-api:latest <ecr-backend-url>:latest`
  - [ ] `docker push <ecr-backend-url>:latest`
- [ ] フロントエンドイメージをビルド・タグ付け・プッシュ
  - [ ] `cd frontend && docker build -t article-manager-frontend:latest .`
  - [ ] `docker tag article-manager-frontend:latest <ecr-frontend-url>:latest`
  - [ ] `docker push <ecr-frontend-url>:latest`
- [ ] ECRコンソールでイメージが正しくpushされたことを確認

---

## Phase 5: ECS Fargate構築 (Week 2, Day 3-4)

### 5.1 IAMロール・ポリシー設定

- [ ] `terraform/iam.tf` を作成
  - [ ] ECSタスク実行ロール (Task Execution Role) を定義
    - [ ] `ecs-tasks.amazonaws.com` を信頼エンティティに設定
    - [ ] ポリシー: `AmazonECSTaskExecutionRolePolicy` をアタッチ
    - [ ] ポリシー: Parameter Store読み取り権限を追加（`ssm:GetParameters`, `ssm:GetParameter`）
    - [ ] ポリシー: KMS復号権限を追加（SecureString復号用、`kms:Decrypt`）
  - [ ] ECSタスクロール (Task Role) を定義
    - [ ] CloudWatch Logs書き込み権限

### 5.2 CloudWatch Logsグループ作成

- [ ] `terraform/cloudwatch.tf` を作成
  - [ ] アプリケーション用ロググループ `/ecs/article-manager-app`
    - [ ] 保持期間: 1日
  - [ ] ~~CloudWatch Metrics~~（不要: コスト削減）
  - [ ] ~~CloudWatch Alarms~~（不要: コスト削減）

### 5.3 ECS Fargateタスク定義

- [ ] `terraform/ecs.tf` を作成
  - [ ] ECSクラスターを定義: `article-manager-cluster`
  - [ ] バックエンドタスク定義を作成
    - [ ] ファミリー名: `article-manager-api`
    - [ ] ネットワークモード: `awsvpc`
    - [ ] 必要なCPU: `256` (.25 vCPU)
    - [ ] 必要なメモリ: `512` MB
    - [ ] コンテナ定義:
      - [ ] 名前: `api`
      - [ ] イメージ: ECR API URL
      - [ ] ポートマッピング: `8080`
      - [ ] ログ設定: CloudWatch Logs (`/ecs/article-manager-app`)
      - [ ] 環境変数（非機密情報）:
        - [ ] `DB_HOST` - RDSエンドポイントを直接指定
        - [ ] `PORT=8080` - 環境変数として直接指定
      - [ ] シークレット (Parameter Storeから取得、`secrets`フィールド使用):
        - [ ] `DB_NAME` - `arn:aws:ssm:region:account:parameter/article-manager/db/name`
        - [ ] `DB_USER` - `arn:aws:ssm:region:account:parameter/article-manager/db/app-user` ← **article_user**
        - [ ] `DB_PASSWORD` - `arn:aws:ssm:region:account:parameter/article-manager/db/app-password` ← **SecureString**
        - [ ] `GEMINI_API_KEY` - `arn:aws:ssm:region:account:parameter/article-manager/api/gemini-api-key` ← **SecureString**
        - [ ] `GOOGLE_BOOKS_API_KEY` - `arn:aws:ssm:region:account:parameter/article-manager/api/google-books-api-key` ← **SecureString**
    - [ ] タスク実行ロール: 上記で作成したIAMロール
  - [ ] フロントエンドタスク定義を作成
    - [ ] ファミリー名: `article-manager-frontend`
    - [ ] ネットワークモード: `awsvpc`
    - [ ] 必要なCPU: `256` (.25 vCPU)
    - [ ] 必要なメモリ: `512` MB
    - [ ] コンテナ定義:
      - [ ] 名前: `frontend`
      - [ ] イメージ: ECR Frontend URL
      - [ ] ポートマッピング: `3000`
      - [ ] ログ設定: CloudWatch Logs (`/ecs/article-manager-app`)
      - [ ] 環境変数:
        - [ ] `NEXT_PUBLIC_API_URL` = `http://localhost:8080`（同一タスク内通信）
    - [ ] タスク実行ロール: 上記で作成したIAMロール

### 5.4 ECSサービス作成

- [ ] `terraform/ecs.tf` に ECSサービスを追加
  - [ ] バックエンドECSサービスを定義
    - [ ] サービス名: `article-manager-api-service`
    - [ ] クラスター: `article-manager-cluster`
    - [ ] タスク定義: `article-manager-api`
    - [ ] 起動タイプ: `FARGATE`
    - [ ] Desired Count: 1（固定、Auto Scaling なし）
    - [ ] ネットワーク設定:
      - [ ] サブネット: パブリックサブネット
      - [ ] セキュリティグループ: ECS SG
      - [ ] パブリックIPの割り当て: `true`
  - [ ] フロントエンドECSサービスを定義
    - [ ] サービス名: `article-manager-frontend-service`
    - [ ] クラスター: `article-manager-cluster`
    - [ ] タスク定義: `article-manager-frontend`
    - [ ] 起動タイプ: `FARGATE`
    - [ ] Desired Count: 1（固定）
    - [ ] ネットワーク設定:
      - [ ] サブネット: パブリックサブネット
      - [ ] セキュリティグループ: ECS SG
      - [ ] パブリックIPの割り当て: `true`

### 5.5 Terraform Apply - ECS構築（GitHub Actions経由）

- [ ] GitHub Actionsで `terraform-apply.yml` ワークフローを実行
- [ ] ワークフローログでECS構築を確認
- [ ] ECSサービスのタスクが正常に起動することを確認
  - [ ] AWS Management Console → ECS → Clusters → `article-manager-cluster`
  - [ ] サービスのタスク数が `1/1 RUNNING` になることを確認
- [ ] CloudWatch Logsでアプリケーションログを確認
  - [ ] エラーがないことを確認

### 5.6 パブリックIP経由でのアクセステスト

- [ ] ECSタスクのパブリックIPを取得
  - [ ] AWS CLI: `aws ecs describe-tasks --cluster article-manager-cluster --tasks <task-arn>`
- [ ] ブラウザでパブリックIPにアクセス
  - [ ] `http://<public-ip>:3000` でフロントエンドが表示されることを確認
  - [ ] `http://<public-ip>:8080/api/health` でバックエンドAPIが動作することを確認

---

## Phase 6: Route 53設定 (Week 2, Day 5)

### 6.0 GitHub Secrets設定（Phase 6: ドメイン関連）

- [ ] GitHubリポジトリのSettings → Secrets and variables → Actionsにアクセス
- [ ] 以下のSecretsを追加登録:
  - [ ] `TF_VAR_domain_name` - 使用するドメイン名（例: `app.example.com`）

### 6.1 ドメイン取得とRoute 53ホストゾーン作成

- [ ] ドメインを取得（Route 53またはお名前.comなど）
- [ ] Route 53でホストゾーンを作成
  - [ ] ドメイン名: `example.com` (実際のドメイン名)
  - [ ] タイプ: パブリック
- [ ] ドメインレジストラのネームサーバーをRoute 53のNSレコードに変更

### 6.2 Terraform でRoute 53リソース定義

- [ ] `terraform/route53.tf` を作成
  - [ ] Route 53ホストゾーンを定義（または既存をdata sourceで参照）
  - [ ] Aレコードを定義
    - [ ] レコード名: `app.example.com`
    - [ ] タイプ: A
    - [ ] TTL: 60秒（変更の反映を早める）
    - [ ] 値: ECSタスクのパブリックIP（初期値、後で動的更新）

### 6.3 手動でRoute 53のAレコードを設定

- [ ] ECSタスクのパブリックIPを取得
- [ ] Route 53コンソールでAレコードを手動作成
  - [ ] レコード名: `app.example.com`
  - [ ] 値: ECSタスクのパブリックIP
  - [ ] TTL: 60秒
- [ ] ブラウザで `http://app.example.com` にアクセスして動作確認

---

## Phase 7: データベースマイグレーション (Week 2, Day 6-7)

### 7.1 ローカルMySQLからデータエクスポート

- [ ] ローカルMySQLコンテナが起動していることを確認
- [ ] mysqldumpでデータをエクスポート
  - [ ] `docker exec db mysqldump -u root -p${MYSQL_ROOT_PASSWORD} article_manager > local_dump.sql`
- [ ] dumpファイルのサイズとレコード数を確認

### 7.2 一時的なS3バケット作成とアップロード

- [ ] S3バケットを作成 (`article-manager-temp-migration-bucket`)
- [ ] dumpファイルをS3にアップロード
  - [ ] `aws s3 cp local_dump.sql s3://article-manager-temp-migration-bucket/`
- [ ] S3にアップロードされたことを確認

### 7.3 踏み台EC2インスタンスの起動 (一時的)

- [ ] EC2インスタンスを起動 (踏み台用)
  - [ ] AMI: Amazon Linux 2023
  - [ ] インスタンスタイプ: `t3.micro`
  - [ ] サブネット: プライベートサブネット (ap-northeast-1a)
  - [ ] セキュリティグループ: RDSアクセス可能なSG
  - [ ] IAMロール: S3読み取り権限付与
- [ ] EC2に SSH接続 (Systems Manager Session Managerを使用)

### 7.4 RDSへのデータインポート

- [ ] EC2インスタンス上でMySQL Clientをインストール
- [ ] S3からdumpファイルをダウンロード
- [ ] RDSにデータをインポート
  - [ ] `mysql -h <rds-endpoint> -u admin -p article_manager < local_dump.sql`
- [ ] データ整合性を確認
  - [ ] テーブル一覧を確認
  - [ ] レコード数を確認

### 7.5 クリーンアップ

- [ ] EC2インスタンスを削除
- [ ] S3バケットを削除
- [ ] ローカルのdumpファイルを削除

### 7.6 アプリケーション用ユーザー（article_user）作成

- [ ] EC2インスタンスからRDSにadminユーザーで接続
  - [ ] `mysql -h <rds-endpoint> -u admin -p article_manager`
- [ ] アプリケーション用ユーザーを作成
  - [ ] `CREATE USER 'article_user'@'%' IDENTIFIED BY '<app-password>';`
  - [ ] パスワードは`.env`の`MYSQL_PASSWORD`（`hEKLvsNNXTmGGEq1`）を使用
- [ ] 必要最小限の権限を付与
  - [ ] `GRANT SELECT, INSERT, UPDATE, DELETE ON article_manager.* TO 'article_user'@'%';`
  - [ ] `GRANT CREATE, ALTER, DROP, INDEX ON article_manager.* TO 'article_user'@'%';`
  - [ ] `FLUSH PRIVILEGES;`
- [ ] article_userで接続テスト
  - [ ] `mysql -h <rds-endpoint> -u article_user -p article_manager`
  - [ ] `SELECT * FROM articles LIMIT 1;` でデータ取得確認

**SQL実行例**:
```sql
-- article_user作成
CREATE USER 'article_user'@'%' IDENTIFIED BY 'hEKLvsNNXTmGGEq1';

-- 権限付与
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, INDEX
  ON article_manager.* TO 'article_user'@'%';

FLUSH PRIVILEGES;

-- 確認
SHOW GRANTS FOR 'article_user'@'%';
```

### 7.7 ECSタスクからRDS接続テスト（article_user使用）

- [ ] CloudWatch Logsでバックエンドログを確認
- [ ] ブラウザで記事一覧ページにアクセス
  - [ ] RDSから正しくデータが取得されることを確認
  - [ ] article_userでの接続が成功していることを確認

---

## Phase 8: 起動/停止ワークフロー構築 (Week 3, Day 1-3)

### 8.1 起動スクリプト作成

- [ ] `scripts/update-route53.sh` を作成
  - [ ] ECSタスクのパブリックIPを取得
  - [ ] Route 53のAレコードを更新
  - [ ] 完了メッセージを表示

### 8.2 GitHub Actions - 起動ワークフロー

- [ ] `.github/workflows/start-infrastructure.yml` を作成
  - [ ] トリガー: `workflow_dispatch`（手動実行）
  - [ ] ステップ1: RDS起動
    - [ ] `aws rds start-db-instance --db-instance-identifier article-manager-db`
  - [ ] ステップ2: RDS起動完了待機
    - [ ] `aws rds wait db-instance-available --db-instance-identifier article-manager-db`
  - [ ] ステップ3: ECSサービスのdesired_countを1に設定
    - [ ] `aws ecs update-service --cluster article-manager-cluster --service article-manager-api-service --desired-count 1`
    - [ ] `aws ecs update-service --cluster article-manager-cluster --service article-manager-frontend-service --desired-count 1`
  - [ ] ステップ4: ECSタスク起動完了待機
  - [ ] ステップ5: ECSタスクのパブリックIPを取得
  - [ ] ステップ6: Route 53のAレコードを更新（`scripts/update-route53.sh` 実行）
  - [ ] ステップ7: 完了通知（ドメインURLを表示）
- [ ] GitHub Secretsに以下を追加設定（AWS認証情報は設定済み）
  - [ ] `ROUTE53_HOSTED_ZONE_ID` - Route 53ホストゾーンID（Terraform outputまたはAWSコンソールから取得）

### 8.3 GitHub Actions - 停止ワークフロー

- [ ] `.github/workflows/stop-infrastructure.yml` を作成
  - [ ] トリガー: `workflow_dispatch`（手動実行）
  - [ ] ステップ1: ECSサービスのdesired_countを0に設定
    - [ ] `aws ecs update-service --cluster article-manager-cluster --service article-manager-api-service --desired-count 0`
    - [ ] `aws ecs update-service --cluster article-manager-cluster --service article-manager-frontend-service --desired-count 0`
  - [ ] ステップ2: ECSタスク停止完了待機
  - [ ] ステップ3: RDS停止
    - [ ] `aws rds stop-db-instance --db-instance-identifier article-manager-db`
  - [ ] ステップ4: 完了通知

### 8.4 起動/停止テスト

- [ ] GitHub Actionsで起動ワークフローを実行
  - [ ] RDSが起動することを確認
  - [ ] ECSタスクが起動することを確認
  - [ ] Route 53が更新されることを確認
  - [ ] ドメインでアクセス可能になることを確認
- [ ] GitHub Actionsで停止ワークフローを実行
  - [ ] ECSタスクが停止することを確認
  - [ ] RDSが停止することを確認

---

## Phase 9: CI/CDパイプライン構築 (Week 3, Day 4-5)

### 9.1 GitHub Actions - アプリケーションデプロイ

- [ ] `.github/workflows/deploy.yml` を作成
  - [ ] トリガー: `main` ブランチへのpush
  - [ ] ジョブ1: Build & Push Backend
    - [ ] AWS認証情報を設定
    - [ ] ECRにログイン
    - [ ] Dockerイメージをビルド (`api/`)
    - [ ] イメージにタグ付け (`latest`, `${GITHUB_SHA}`)
    - [ ] ECRにpush
  - [ ] ジョブ2: Build & Push Frontend
    - [ ] AWS認証情報を設定
    - [ ] ECRにログイン
    - [ ] Dockerイメージをビルド (`frontend/`)
    - [ ] イメージにタグ付け (`latest`, `${GITHUB_SHA}`)
    - [ ] ECRにpush
  - [ ] ジョブ3: Deploy to ECS
    - [ ] 新しいタスク定義を登録（バックエンド）
    - [ ] ECSサービスを更新 (`force-new-deployment`)
    - [ ] 新しいタスク定義を登録（フロントエンド）
    - [ ] ECSサービスを更新
    - [ ] デプロイ完了を待機
- [ ] GitHub Secretsに以下を追加設定（Terraform outputまたはAWSコンソールから取得）
  - [ ] `ECR_FRONTEND_REPOSITORY` - フロントエンドECRリポジトリURL
  - [ ] `ECR_BACKEND_REPOSITORY` - バックエンドECRリポジトリURL
  - [ ] `ECS_CLUSTER` - ECSクラスター名（`article-manager-cluster`）
  - [ ] `ECS_SERVICE_FRONTEND` - フロントエンドECSサービス名
  - [ ] `ECS_SERVICE_BACKEND` - バックエンドECSサービス名

### 9.2 CI/CD動作確認

- [ ] バックエンドコードに軽微な変更を加えてcommit・push
  - [ ] GitHub Actionsでワークフローが起動することを確認
  - [ ] ビルド成功、ECRへのpush成功、ECSデプロイ成功を確認
- [ ] フロントエンドコードに軽微な変更を加えてcommit・push
  - [ ] GitHub Actionsで自動デプロイ成功を確認

---

## Phase 10: ドキュメント更新・運用準備 (Week 3, Day 6-7)

### 10.1 運用手順書作成

- [ ] `docs/runbook.md` を作成
  - [ ] 起動手順（GitHub Actions経由）
  - [ ] 停止手順（GitHub Actions経由）
  - [ ] デプロイ手順（GitHub Actions経由）
  - [ ] トラブルシューティング
    - [ ] ECSタスクが起動しない場合
    - [ ] RDS接続エラーの場合
    - [ ] Route 53更新が失敗する場合
  - [ ] コスト確認方法
  - [ ] ログ確認手順

### 10.2 AWS移行ガイド作成

- [ ] `docs/aws-migration-guide.md` を作成
  - [ ] 移行手順の概要
  - [ ] 前提条件
  - [ ] Terraform実行手順
  - [ ] データマイグレーション手順
  - [ ] ロールバック手順

### 10.3 プロジェクトドキュメント更新

- [ ] `CLAUDE.md` を更新
  - [ ] プロジェクト概要にAWS環境を追加
  - [ ] 開発コマンドにAWS関連コマンドを追加
- [ ] `README.md` を更新（必要に応じて）
  - [ ] デプロイセクションを追加

---

## Phase 11: 総合テストと品質チェック (Week 4, Day 1-2)

### 11.1 インフラテスト

- [ ] Terraform構成が正しいことを確認
  - [ ] GitHub Actionsワークフローで `terraform validate` が成功
  - [ ] GitHub Actionsワークフローで `terraform plan` がエラーなく完了
- [ ] すべてのAWSリソースが正しく作成されていることを確認
  - [ ] VPC、サブネット、セキュリティグループ
  - [ ] RDS（Single-AZ、db.t4g.micro）
  - [ ] ECS（Fargate、パブリックIP割り当て）
  - [ ] Route 53
  - [ ] ECR

### 11.2 E2E動作確認

- [ ] 起動ワークフローを実行し、全機能が動作することを確認
  - [ ] RDSが起動
  - [ ] ECSタスクが起動
  - [ ] Route 53が更新
  - [ ] ドメインでアクセス可能
- [ ] ブラウザで実際にアプリケーションを操作
  - [ ] 記事一覧ページにアクセス
  - [ ] 記事作成（AI生成含む）
  - [ ] 記事編集
  - [ ] 記事削除
  - [ ] タグ管理
  - [ ] 検索機能
  - [ ] 書籍推薦機能
- [ ] すべての機能が正常に動作することを確認

### 11.3 コスト確認

- [ ] AWS Cost Explorerで実際のコストを確認
- [ ] 目標コスト（月$10以内）に収まっているか確認
- [ ] 停止機能が正しく動作し、コストが削減されることを確認

---

## Phase 12: 実装後の振り返り (Week 4, Day 3)

### 12.1 振り返り記録

- [ ] このファイル (`tasklist.md`) の最下部に振り返りを記載
  - [ ] 実装完了日を記載
  - [ ] 計画と実績の差分を記載
  - [ ] 学んだことを記載
  - [ ] 次回への改善提案を記載

### 12.2 最終確認

- [ ] 全タスクが `[x]` になっていることを確認
- [ ] スキップしたタスクがある場合、技術的理由が明記されているか確認
- [ ] ドキュメントが最新の状態に更新されているか確認
- [ ] AWS環境が正常に動作しているか最終確認

---

## 実装後の振り返り

### 実装完了日
{YYYY-MM-DD}

### 計画と実績の差分

**計画と異なった点**:
- {実装中に発生した技術的な問題と解決策}
- {当初の想定と異なったアーキテクチャ選択}
- {追加で必要になったAWSリソース}

**新たに必要になったタスク**:
- {実装中に追加したタスク}
- {なぜ追加が必要だったか}

**技術的理由でスキップしたタスク**（該当する場合のみ）:
- {タスク名}
  - スキップ理由: {具体的な技術的理由}
  - 代替実装: {何に置き換わったか}

**⚠️ 注意**: 「時間の都合」「難しい」などの理由でスキップしたタスクはここに記載しないこと。全タスク完了が原則。

### 学んだこと

**技術的な学び**:
- {Terraform、ECS、RDS、Route 53を使った実装で学んだこと}
- {起動/停止機能の実装で学んだこと}
- {コスト最適化のテクニック}

**プロセス上の改善点**:
- {ステアリングファイルの活用方法}
- {タスクの粒度とスケジュール管理}
- {ドキュメント駆動開発の効果}

### 次回への改善提案
- {次回のインフラ構築で気をつけること}
- {Terraformモジュール化やより効率的なCI/CD設計}
- {コスト最適化やパフォーマンスチューニングの計画}
