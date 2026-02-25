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

WORKDIR /root/

# ビルドステージからバイナリをコピー
COPY --from=builder /app/server .

# ポート8080を公開
EXPOSE 8080

# バイナリを実行
CMD ["./server"]
