# 要求内容

## 概要

Docker Composeベースのローカル開発環境からAWSクラウドインフラへの移行を実施し、**個人学習用のコスト最適化されたインフラ基盤**を構築します。

## 背景

現在のArticle Managerは、docker-compose.ymlを使用したローカル環境でのみ動作しています。この構成では以下の課題があります:

- **本番デプロイ不可**: ローカル環境専用の設定のため、外部からのアクセスができない
- **AWSインフラの学習機会がない**: Docker Composeのみで、クラウドインフラの知識が身につかない
- **運用の手動化**: バックアップ、モニタリング、ログ管理などが手動運用
- **CI/CDパイプラインがない**: デプロイの自動化がされていない

AWSへの移行により、**クラウドインフラの学習**と**実践的なCI/CD構築**を実現しつつ、**個人利用のためコストを最小限**に抑えます。

## 移行の目的

### 主な目的
- **AWS学習**: ECS Fargate、RDS、Terraform、GitHub Actionsなどの実践的な学習
- **個人利用**: 自分専用のアプリケーション環境として使用
- **コスト最小化**: 使わない時は停止し、実質月額コストを**$10以内**に抑える

### 学習目標
- Infrastructure as Code (Terraform) の実践
- コンテナオーケストレーション (ECS Fargate) の理解
- マネージドデータベース (RDS) の運用
- CI/CDパイプライン (GitHub Actions) の構築
- AWSネットワーク設計 (VPC、セキュリティグループ) の学習

## 実装対象の機能

### 1. AWS ECS Fargateによるコンテナ実行（最小構成）

- Docker ComposeからAWS ECS (Elastic Container Service) への移行
- Fargateを使用したサーバーレスコンテナ実行環境の構築
- **固定1タスク構成**（Auto Scaling なし、コスト削減）
- **ALB不使用**（Route 53による動的DNS更新で直接アクセス）
- タスク構成: 0.25 vCPU / 512 MB メモリ

### 2. Amazon RDS for MySQLへのデータベース移行（最小構成）

- ローカルMySQLコンテナからAmazon RDS for MySQLへの移行
- **Single-AZ構成**（Multi-AZ は使わない、コスト削減）
- インスタンスタイプ: **db.t4g.micro**（無料枠対象）
- 自動バックアップとポイントインタイムリカバリの設定（保持期間1日）
- セキュリティグループによるネットワークアクセス制御（ECSからのみ）

### 3. Infrastructure as Code (IaC) の導入

- Terraformによるインフラ定義のコード化
- 環境変数管理 (AWS Systems Manager Parameter Store / Secrets Manager)
- VPCネットワーク設計 (パブリック/プライベートサブネット分離)
- セキュリティグループの設定

### 4. Route 53による動的DNS管理（ALB不使用）

- ドメイン取得とRoute 53によるDNS管理
- ECSタスク起動時にパブリックIPを取得
- Route 53のAレコードを自動更新（例: `app.example.com` → タスクのパブリックIP）
- **HTTPのみ対応**（HTTPS はPhase 2で検討）

### 5. 起動/停止機能（コスト最適化の要）

- **GitHub Actions ワークフローによる起動/停止**
  - 手動トリガーでECSタスクとRDSを起動
  - 手動トリガーでECSタスクとRDSを停止
  - 起動時にRoute 53のAレコードを自動更新
- **使わない時は完全停止**して課金を最小化
- 起動スクリプト: `.github/workflows/start-infrastructure.yml`
- 停止スクリプト: `.github/workflows/stop-infrastructure.yml`

### 6. CI/CDパイプライン構築

- GitHub Actions による自動ビルド・デプロイパイプライン
- ECR (Elastic Container Registry) へのDockerイメージ push
- ECSサービスの自動更新（`main`ブランチpush時）
- デプロイワークフロー: `.github/workflows/deploy.yml`

### 7. モニタリングとロギング（最小構成）

- CloudWatch Logsへのアプリケーションログ集約（デバッグ用）
- ログ保持期間: **1日間**（コスト最小化）
- CloudWatch Metrics、Alarmsは使用しない（コスト削減）

## 受け入れ条件

### Infrastructure as Code (Terraform)

- [ ] Terraformコード (`terraform/`) が正常にapply可能
- [ ] VPC、サブネット、セキュリティグループが正しく作成される
- [ ] ECSクラスター、タスク定義、サービスが作成される
- [ ] RDSインスタンス（Single-AZ）がプライベートサブネットに作成される
- [ ] Route 53ホストゾーンとAレコードが作成される
- [ ] terraform.tfstateがS3バックエンドで管理される

### コンテナデプロイ

- [ ] フロントエンドコンテナがECS Fargateで起動する
- [ ] バックエンドAPIコンテナがECS Fargateで起動する
- [ ] 環境変数がAWS Systems Manager Parameter Store / Secrets Managerから注入される
- [ ] RDSへのデータベース接続が成功する
- [ ] ドメイン（例: `app.example.com`）経由でフロントエンドにアクセス可能

### ネットワーク・セキュリティ

- [ ] ECSタスクはパブリックサブネットに配置され、パブリックIPが割り当てられる
- [ ] RDSはプライベートサブネットに配置され、ECSからのみアクセス可能
- [ ] セキュリティグループで必要最小限のポート（80, 3000, 8080, 3306）のみ許可
- [ ] HTTPで通信可能（HTTPSはPhase 2で実装）

### 起動/停止機能

- [ ] GitHub Actionsで起動ワークフローが正常に実行される
  - [ ] RDSインスタンスが起動する
  - [ ] ECSタスクが起動する
  - [ ] タスクのパブリックIPを取得できる
  - [ ] Route 53のAレコードが自動更新される
- [ ] GitHub Actionsで停止ワークフローが正常に実行される
  - [ ] ECSタスクが停止する
  - [ ] RDSインスタンスが停止する
- [ ] 停止中は課金が最小化される（RDS停止中は約90%削減）

### CI/CD

- [ ] GitHub Actionsで`main`ブランチpush時に自動デプロイが実行される
- [ ] Dockerイメージが正しくECRにpushされる
- [ ] ECSサービスが新しいイメージで自動更新される
- [ ] デプロイ後、アプリケーションが正常に動作する

### モニタリング・ロギング

- [ ] アプリケーションログがCloudWatch Logsに送信される
- [ ] ログ保持期間が1日間に設定される

### データベース

- [ ] ローカルのMySQLデータがRDSにマイグレーションされる
- [ ] RDSの自動バックアップが有効化される（保持期間1日）
- [ ] Single-AZ構成で動作する
- [ ] データベース接続文字列が環境変数で管理される

## 成功指標

### 技術的指標

- **学習目標達成**: Terraform、ECS、RDS、GitHub Actionsの基本的な理解
- **動作確認**: ドメイン経由でアプリケーションにアクセス可能
- **起動/停止**: GitHub Actions経由で5分以内に起動・停止が完了
- **デプロイ**: CI/CDパイプラインで自動デプロイが5分以内に完了

### コスト指標

- **月額コスト**: $20-25/月（常時起動の場合）
  - ECS Fargate: $6.5/月（0.25 vCPU × 730時間）
  - RDS (db.t4g.micro, Single-AZ): $15/月
  - Route 53: $0.50/月
  - その他（ECR, CloudWatch Logs等）: $1以内
- **実質コスト**: 週20時間稼働なら**$5-6/月**
  - 使わない時は停止することでコスト最小化
  - RDS停止中は約90%のコスト削減
- **目標**: 月$10以内（起動/停止を活用）

### 運用改善指標

- **起動/停止時間**: 手動操作 → GitHub Actions経由で5分以内
- **デプロイ時間**: 手動デプロイ30分 → 自動デプロイ5分以内
- **インフラ変更の再現性**: コード化により100%再現可能

## スコープ外

以下はこのフェーズでは実装しません:

### Phase 2以降に延期

- **Let's Encrypt によるHTTPS対応**: ECSタスク内での証明書管理
- **CloudFlare統合**: 無料HTTPSプロキシ
- **Multi-AZ構成**: 高可用性は個人学習では不要
- **Application Load Balancer (ALB)**: コストが高い（$20/月）
- **Auto Scaling**: 固定1タスクで十分
- **CloudWatch Metrics/Alarms**: 個人学習では監視は最小限で十分
- **マルチリージョン構成**: 単一リージョン (ap-northeast-1) のみで運用
- **CDN (CloudFront)**: 静的コンテンツ配信の最適化は後回し
- **ElastiCache (Redis)**: キャッシュ層の導入は後回し
- **WAF (Web Application Firewall)**: セキュリティ強化は次フェーズ
- **Lambda統合**: サーバーレス処理の導入は将来検討
- **X-Rayによる分散トレーシング**: 個人学習では不要

### 運用最適化 (後回し)

- **コスト最適化**: Reserved InstancesやSavings Plansの活用
- **パフォーマンスチューニング**: 詳細なクエリ最適化やインデックス調整
- **スケジュール起動/停止**: Lambda + EventBridgeでの自動化
- **CloudWatch Metrics/Alarms**: パフォーマンス監視やアラート設定

### セキュリティ強化 (Phase 2)

- **AWS GuardDuty**: 脅威検出サービスの導入
- **AWS Config**: コンプライアンス監査
- **VPCフローログ**: ネットワークトラフィックの詳細分析
- **AWS Security Hub**: セキュリティ統合管理

## 参照ドキュメント

### プロジェクト内ドキュメント

- `docs/product-requirements.md` - プロダクト要求定義書
- `docs/functional-design.md` - 機能設計書
- `docs/architecture.md` - アーキテクチャ設計書 (ローカル環境)
- `CLAUDE.md` - プロジェクトメモリ

### AWSドキュメント

- [AWS ECS Best Practices](https://docs.aws.amazon.com/AmazonECS/latest/bestpracticesguide/intro.html)
- [Amazon RDS Best Practices](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_BestPractices.html)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [GitHub Actions for AWS](https://github.com/aws-actions)

### 移行ガイド

- Docker Compose → ECS移行パターン（ALB不使用）
- オンプレミスMySQL → RDS移行手順
- Route 53動的DNS更新のベストプラクティス
- 環境変数管理のベストプラクティス (Parameter Store vs Secrets Manager)
- 起動/停止スクリプトの実装パターン

## コスト試算の詳細

### 常時起動の場合（月額）

| サービス | 構成 | 月額コスト |
|---------|------|-----------|
| ECS Fargate | 0.25 vCPU × 730h | $5.84 |
| ECS Fargate | 512 MB × 730h | $0.64 |
| RDS MySQL | db.t4g.micro Single-AZ | $12.41 |
| RDS ストレージ | 20 GB gp3 | $2.53 |
| Route 53 | ホストゾーン1個 | $0.50 |
| ECR | イメージストレージ | $0.50 |
| CloudWatch Logs | 1日保持（最小限） | $0.10 |
| **合計** | | **約$22.5/月** |

### 実質コスト（週20時間稼働の場合）

| サービス | 稼働時間 | 月額コスト |
|---------|---------|-----------|
| ECS Fargate | 80時間/月 | $0.74 |
| RDS MySQL | 80時間/月 | $1.36 |
| RDS ストレージ | 常時 | $2.53 |
| Route 53 | 常時 | $0.50 |
| その他（ECR等） | 常時 | $0.60 |
| **合計** | | **約$5.7/月** |

**無料枠を活用すれば初年度は更に安くなる可能性あり**
