# ビルドステージ
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 依存関係をコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド（最適化）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -trimpath -o server ./cmd/server/main.go

# 実行ステージ
FROM alpine:3.23

# CA証明書を追加（HTTPS通信用）
RUN apk --no-cache add ca-certificates

# 非rootユーザーを作成
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# ビルドステージからバイナリをコピー
COPY --from=builder /app/server .

# バイナリの所有権を変更
RUN chown appuser:appuser server

# 非rootユーザーに切り替え
USER appuser

# ポート8080を公開
EXPOSE 8080

# バイナリを実行
CMD ["./server"]
