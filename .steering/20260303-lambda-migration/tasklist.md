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

## Phase 0: 既存インフラの確認（ECS移行で完了済み）

### 0.1 Terraform基盤構築

- [x] `terraform/` ディレクトリ作成（✅ ECS移行で完了済み）
- [x] S3バケット作成（tfstate保存用）（✅ ECS移行で完了済み）
- [x] `terraform/backend.tf` 作成（✅ ECS移行で完了済み）
- [x] `terraform/provider.tf` 作成（✅ ECS移行で完了済み）
- [x] `terraform/variables.tf` 作成（✅ ECS移行で完了済み）
- [x] `terraform/outputs.tf` 作成（✅ ECS移行で完了済み）

### 0.2 VPCネットワーク構築

- [x] `terraform/vpc.tf` 作成（✅ ECS移行で完了済み）
  - [x] VPCリソース定義（10.0.0.0/16）
  - [x] パブリックサブネット（10.0.1.0/24、ap-northeast-1a）
  - [x] プライベートサブネット（10.0.11.0/24、ap-northeast-1a）
  - [x] インターネットゲートウェイ定義
  - [x] ルートテーブル定義

### 0.3 セキュリティグループ定義

- [x] `terraform/security_groups.tf` 作成（✅ ECS移行で完了済み）
  - [x] ECS（現在Lambda用に流用）用セキュリティグループ
  - [x] RDS用セキュリティグループ

### 0.4 RDS MySQL構築

- [x] `terraform/rds.tf` 作成（✅ ECS移行で完了済み）
  - [x] DBサブネットグループ定義
  - [x] RDSインスタンス定義（db.t4g.micro、Single-AZ）
  - [x] 自動バックアップ設定（保持期間1日）
  - [x] セキュリティグループアタッチ
- [x] Terraform apply実行（✅ ECS移行で完了済み）
- [x] RDSエンドポイント確認（✅ ECS移行で完了済み）

### 0.5 ECR構築

- [x] `terraform/ecr.tf` 作成（✅ ECS移行で完了済み）
  - [x] フロントエンド用ECRリポジトリ定義
  - [x] バックエンド用ECRリポジトリ定義
  - [x] ライフサイクルポリシー定義
- [x] Terraform apply実行（✅ ECS移行で完了済み）
- [x] ECRリポジトリURL確認（✅ ECS移行で完了済み）

### 0.6 CloudWatch Logs構築

- [x] `terraform/cloudwatch.tf` 作成（✅ ECS移行で完了済み）
  - [x] ロググループ定義（保持期間1日）

### 0.7 IAMロール構築

- [x] `terraform/iam.tf` 作成（✅ ECS移行で完了済み）
  - [x] ECSタスク実行ロール定義（Lambda移行では調整が必要）
  - [x] CloudWatch Logs書き込み権限

---

## Phase 1: 既存インフラの調整（Lambda用）

### 1.1 セキュリティグループの調整

- [ ] `terraform/security_groups.tf` を更新
  - [ ] Lambda用セキュリティグループを追加
    - [ ] アウトバウンド: All traffic（RDS、Parameter Store、インターネットアクセス用）
  - [ ] RDS用セキュリティグループを更新
    - [ ] インバウンド: MySQL (3306) from Lambda SG

### 1.2 IAMロールの調整

- [ ] `terraform/iam.tf` を更新
  - [ ] Lambda実行ロールを追加
    - [ ] VPC実行ポリシー（AWSLambdaVPCAccessExecutionRole）
    - [ ] Parameter Store読み取り権限（SSM GetParameter）
    - [ ] CloudWatch Logs書き込み権限
    - [ ] RDSアクセス権限（必要に応じて）
  - [ ] ECS固有のロールは削除または無効化

### 1.3 CloudWatch Logsの調整

- [ ] `terraform/cloudwatch.tf` を更新
  - [ ] Lambda用ロググループを追加（`/aws/lambda/article-manager-api`）
  - [ ] 既存のECS用ロググループは削除または無効化

### 1.4 ECS固有リソースの削除

- [ ] `terraform/ecs.tf` を削除（Lambda移行により不要）
- [ ] `terraform/service_discovery.tf` を削除（Lambda移行により不要）

---

## Phase 2: Go Lambda Handlerの実装

### 2.1 依存関係の追加

- [x] `go.mod`にLambda関連ライブラリを追加（✅ Lambda移行で完了済み）
  - [x] `github.com/aws/aws-lambda-go`
  - [x] `github.com/awslabs/aws-lambda-go-api-proxy`
  - [x] `github.com/aws/aws-sdk-go-v2/config`
  - [x] `github.com/aws/aws-sdk-go-v2/service/ssm`
  - [x] `go mod tidy`を実行

### 2.2 Parameter Store統合

- [ ] `internal/config/parameter_store.go`を作成
  - [ ] `LoadFromParameterStore()`関数を実装
  - [ ] SSMクライアントでParameter Storeから取得
  - [ ] DB接続情報、APIキーを構造体にマッピング

- [ ] `internal/config/config.go`を確認
  - [ ] `Config`構造体が既存のまま使用可能か確認
  - [ ] Parameter Store対応のコメントを追加

### 2.3 Lambda Entrypoint

- [ ] `cmd/lambda/main.go`を作成
  - [ ] `init()`関数でDB接続初期化（グローバル変数）
  - [ ] `init()`関数でマイグレーション実行
  - [ ] `init()`関数で依存性注入（既存`cmd/server/main.go`を参考）
  - [ ] `init()`関数でHTTPルーター設定（既存コードを移植）
  - [ ] `Handler()`関数を実装（aws-lambda-go-api-proxy統合）
  - [ ] `main()`関数を実装（`lambda.Start(Handler)`）

### 2.4 DB接続プール設定の最適化

- [ ] `internal/infrastructure/database/mysql.go`を更新
  - [ ] Lambda環境向けの接続プール設定を追加
  - [ ] `SetMaxOpenConns(5)`
  - [ ] `SetMaxIdleConns(2)`
  - [ ] `SetConnMaxLifetime(5 * time.Minute)`

### 2.5 Dockerfileとビルド設定追加

- [ ] `api/Dockerfile.lambda`を作成
  - [ ] マルチステージビルド構成（builder + runtime）
  - [ ] ビルダーステージ: `golang:1.21-alpine`でバイナリビルド
  - [ ] ランタイムステージ: `public.ecr.aws/lambda/provided:al2023`をベース
  - [ ] バイナリを`${LAMBDA_RUNTIME_DIR}/bootstrap`にコピー
- [ ] `api/Makefile`を作成（オプション）
  - [ ] `docker-build-lambda`ターゲットを追加
  - [ ] `docker-run-lambda`ターゲットを追加
  - [ ] `clean`ターゲットを追加

### 2.6 ローカルテスト（オプション）

- [ ] Dockerイメージをローカルビルド
  - [ ] `docker build -f Dockerfile.lambda -t article-manager-lambda:local .`
- [ ] Lambda Runtime Interface Emulator (RIE)でローカル起動
  - [ ] `docker run -p 9000:8080 --env-file .env.local article-manager-lambda:local`
- [ ] `curl`でAPIテスト
  - [ ] `curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'`
  - [ ] API Proxyの動作確認

---

## Phase 3: Terraform Lambda定義

### 3.1 Parameter Store定義

- [ ] `terraform/parameters.tf`を更新
  - [ ] DB接続情報をParameter Storeに定義
    - [ ] `/article-manager/db/host` - RDSエンドポイント
    - [ ] `/article-manager/db/port` - 3306
    - [ ] `/article-manager/db/user` - admin
    - [ ] `/article-manager/db/password` - パスワード（SecureString）
    - [ ] `/article-manager/db/name` - article_manager
  - [ ] Gemini API KeyをParameter Storeに定義
    - [ ] `/article-manager/gemini/api-key` - Gemini API Key（SecureString）
  - [ ] Google Books API KeyをParameter Storeに定義
    - [ ] `/article-manager/google-books/api-key` - Google Books API Key（SecureString）

### 3.2 ECRリポジトリ定義

- [ ] `terraform/ecr.tf`を更新
  - [ ] Lambda用ECRリポジトリを追加
    - [ ] リポジトリ名: `article-manager-lambda-api`
    - [ ] イメージスキャン: 有効
    - [ ] ライフサイクルポリシー: 最新10イメージを保持

### 3.3 Lambda関数定義

- [ ] `terraform/lambda.tf`を作成
  - [ ] Lambda関数を定義
    - [ ] 関数名: `article-manager-api`
    - [ ] パッケージタイプ: `Image`（Docker Container Image）
    - [ ] イメージURI: `${aws_ecr_repository.lambda_api.repository_url}:latest`
    - [ ] メモリ: 512MB
    - [ ] タイムアウト: 30秒
    - [ ] VPC設定（パブリックサブネット、Lambda SG）
    - [ ] 環境変数設定（APP_ENV=production）
    - [ ] IAMロールをアタッチ（Lambda実行ロール）
  - [ ] Function URLs設定
    - [ ] 認証: NONE
    - [ ] CORS設定（フロントエンドドメインを許可）
    - [ ] HTTPメソッド: GET, POST, PUT, DELETE

### 3.4 ECS（フロントエンド）定義

- [ ] `terraform/ecs.tf`を作成
  - [ ] ECSクラスター定義
  - [ ] フロントエンド用ECSタスク定義
    - [ ] Fargateタイプ
    - [ ] CPU/メモリ設定
    - [ ] コンテナ定義（Next.js）
    - [ ] 環境変数（`NEXT_PUBLIC_API_URL` = Lambda Function URL）
  - [ ] ECSサービス定義
    - [ ] パブリックサブネット配置
    - [ ] セキュリティグループ
    - [ ] パブリックIP割り当て

---

## Phase 4: インフラのデプロイとテスト

### 4.1 Lambda Docker Imageのビルドとプッシュ

- [ ] ECRにログイン
  - [ ] `aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com`
- [ ] Dockerイメージをビルド
  - [ ] `cd api && docker build -f Dockerfile.lambda -t article-manager-lambda-api:latest .`
- [ ] イメージにECRタグを付与
  - [ ] `docker tag article-manager-lambda-api:latest <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/article-manager-lambda-api:latest`
- [ ] ECRにプッシュ
  - [ ] `docker push <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/article-manager-lambda-api:latest`

### 4.2 Terraformデプロイ

- [ ] `terraform init`を実行（既存のバックエンドを確認）
- [ ] `terraform plan`を実行して確認
  - [ ] Lambda関数が作成される
  - [ ] Parameter Storeが設定される
  - [ ] セキュリティグループが更新される
  - [ ] IAMロールが更新される
- [ ] `terraform apply`を実行

### 4.3 Parameter Storeに値を設定

- [ ] AWS CLIまたはコンソールでParameter Storeに値を設定
  - [ ] `/article-manager/db/host` = RDSエンドポイント（Terraform outputから取得）
  - [ ] `/article-manager/db/port` = 3306
  - [ ] `/article-manager/db/user` = admin
  - [ ] `/article-manager/db/password` = （強力なパスワード）
  - [ ] `/article-manager/db/name` = article_manager
  - [ ] `/article-manager/gemini/api-key` = （Gemini APIキー）
  - [ ] `/article-manager/google-books/api-key` = （Google Books APIキー）

### 4.4 Lambda Function URLのテスト

- [ ] Function URLを取得（`terraform output`）
- [ ] `curl`でヘルスチェック
  - [ ] `GET /` が200 OKを返す
  - [ ] `GET /api/health` が200 OKを返す（ヘルスチェックエンドポイント実装済みの場合）

### 4.5 RDS接続テスト

- [ ] Lambda CloudWatch Logsを確認
  - [ ] DB接続成功ログを確認
  - [ ] マイグレーション実行ログを確認

### 4.6 API動作確認

- [ ] 記事一覧取得をテスト
  - [ ] `GET /api/articles` が200 OKを返す（空配列）
- [ ] 記事作成をテスト
  - [ ] `POST /api/articles` で記事を作成
  - [ ] レスポンスが正しく返る
- [ ] 記事生成（Gemini）をテスト
  - [ ] `POST /api/articles/generate` でテスト
  - [ ] Gemini APIが正常に応答する（5-10秒）
- [ ] 書籍推薦をテスト
  - [ ] `GET /api/book-recommendations` でテスト
  - [ ] Gemini APIが正常に応答する

---

## Phase 5: CI/CD構築

### 5.1 GitHub Actions ワークフロー作成

- [ ] `.github/workflows/deploy-lambda.yml`を作成
  - [ ] トリガー: `main`ブランチへのpush（`api/**`の変更時）
  - [ ] ステップ1: AWS認証（aws-actions/configure-aws-credentials）
  - [ ] ステップ2: ECRログイン（aws-actions/amazon-ecr-login）
  - [ ] ステップ3: Dockerイメージビルド
    - [ ] `docker build -f api/Dockerfile.lambda -t article-manager-lambda-api:$GITHUB_SHA api/`
  - [ ] ステップ4: イメージタグ付与
    - [ ] `latest`タグと`$GITHUB_SHA`タグを付与
  - [ ] ステップ5: ECRにプッシュ
    - [ ] `docker push <ecr-url>/article-manager-lambda-api:latest`
    - [ ] `docker push <ecr-url>/article-manager-lambda-api:$GITHUB_SHA`
  - [ ] ステップ6: Lambda関数イメージ更新
    - [ ] `aws lambda update-function-code --function-name article-manager-api --image-uri <ecr-url>/article-manager-lambda-api:latest`
  - [ ] ステップ7: デプロイ完了通知

### 5.2 GitHub Secrets設定

- [ ] GitHub Secretsに以下を設定
  - [ ] `AWS_ACCESS_KEY_ID`
  - [ ] `AWS_SECRET_ACCESS_KEY`
  - [ ] `AWS_REGION` = ap-northeast-1
  - [ ] `LAMBDA_FUNCTION_NAME` = article-manager-api

### 5.3 CI/CDテスト

- [ ] ダミーコミットを`main`ブランチにpush
- [ ] GitHub Actionsが自動実行される
- [ ] Lambda関数が更新される
- [ ] デプロイ後、APIが正常動作する

---

## Phase 6: フロントエンド更新

### 6.1 API URL変更

- [ ] `frontend/config/constants.ts`を更新
  - [ ] `API_BASE_URL`をLambda Function URLに変更
  - [ ] 環境変数`NEXT_PUBLIC_API_URL`で管理

### 6.2 CORS確認

- [ ] Terraformの`lambda.tf`でCORS設定を確認
  - [ ] ECS Fargateからのアクセスを許可（`allow_origins = ["*"]` または特定のURL）
  - [ ] `Access-Control-Allow-Origin`ヘッダーが正しい

### 6.3 フロントエンドデプロイ

- [ ] ECS Fargateタスク定義で環境変数を更新
  - [ ] `NEXT_PUBLIC_API_URL` = Lambda Function URL
- [ ] ECS Fargateに再デプロイ

### 6.4 E2Eテスト

- [ ] フロントエンドから記事一覧取得をテスト
- [ ] フロントエンドから記事作成をテスト
- [ ] フロントエンドから記事編集をテスト
- [ ] フロントエンドから記事削除をテスト
- [ ] フロントエンドから記事生成（Gemini）をテスト
- [ ] フロントエンドからタグ管理をテスト
- [ ] フロントエンドから検索機能をテスト
- [ ] フロントエンドから書籍推薦をテスト
- [ ] 全ての機能が正常に動作することを確認

---

## Phase 7: ドキュメント更新と振り返り

### 7.1 ドキュメント更新

- [ ] `docs/architecture.md`を更新
  - [ ] Lambda + RDSアーキテクチャ図を追加
  - [ ] システム構成図を更新
  - [ ] デプロイ手順を更新
- [ ] `CLAUDE.md`を更新
  - [ ] Lambda環境でのセットアップ手順を追加
  - [ ] 環境変数の設定方法を追加
  - [ ] 開発コマンドにLambda関連コマンドを追加

### 7.2 実装後の振り返り

- [ ] `.steering/20260303-lambda-migration/tasklist.md`の振り返りを記録
  - [ ] 実装完了日を記載
  - [ ] 計画と実績の差分を記載
  - [ ] 学んだことを記載
  - [ ] 次回への改善提案を記載
  - [ ] コスト実績を記載

### 7.3 最終確認

- [ ] 全タスクが `[x]` になっていることを確認
- [ ] スキップしたタスクがある場合、技術的理由が明記されているか確認
- [ ] ドキュメントが最新の状態に更新されているか確認
- [ ] Lambda環境が正常に動作しているか最終確認

---

## 実装後の振り返り

### 実装完了日
{YYYY-MM-DD}

### 計画と実績の差分

**計画と異なった点**:
- {計画時には想定していなかった技術的な変更点}
- {実装方針の変更とその理由}
- {ECS移行で構築済みのインフラをどのように活用したか}

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
- {Lambdaの実装を通じて学んだ技術的な知見}
- {aws-lambda-go-api-proxyの使い方}
- {Parameter Storeの活用方法}
- {既存インフラの活用方法（ECS移行からの引き継ぎ）}

**プロセス上の改善点**:
- {タスク管理で良かった点}
- {ステアリングファイルの活用方法}
- {既存インフラを活用した効率的な移行プロセス}

### 次回への改善提案
- {次回の機能追加で気をつけること}
- {より効率的な実装方法}
- {タスク計画の改善点}
- {Lambda環境での運用改善点}

### コスト実績

**月額コスト**:
- Lambda: ${実際のコスト}（無料枠内に収まったか）
- RDS: ${実際のコスト}（既存インフラ継続使用）
- その他: ${実際のコスト}
- 合計: ${実際のコスト}

**無料枠の活用状況**:
- {どの程度無料枠内で収まったか}
- {ECS移行からLambda移行でのコスト削減効果}

**ECSとの比較**:
- ECS月額コスト（想定）: $22-25
- Lambda月額コスト（実績）: ${実際のコスト}
- コスト削減額: ${削減額}
