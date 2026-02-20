# レポジトリ層のユニットテスト方法
1. `docker-compose -f docker-compose-test.yml up -d`でmysqlコンテナ起動
2. mysqlコンテナが起動していることを確認したら、`docker exec -i {コンテナID} mysql -utest_user -ptest_password article_manager_test < api/internal/infrastructure/database/migrations/000001_create_articles_table.up.sql`でマイグレーション実行(テーブルの構造が変わらない限りは初回のみ実行)
3. apiディレクトリに移動し`go test -v ./internal/infrastructure/repository/`を実行
5. `docker-compose -f docker-compose-test.yml down`でmysqlコンテナ停止

# context7の利用方法
claudeプロンプトで、最後に'use context7'と付与するだけ

# 将来実装したいこと
・auth認証
・opentelemetry
・本推薦機能の見直し
・マイクロサービス化
