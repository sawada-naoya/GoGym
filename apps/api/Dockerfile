# Go API サーバー用 Dockerfile
# マルチステージビルドで本番用の軽量・セキュアなイメージを作成

# ==============================
# ビルドステージ
# ==============================
FROM golang:1.24-alpine AS builder

# 開発ツールとSSL証明書をインストール
# git: プライベートリポジトリ依存関係の取得用
# ca-certificates: HTTPS通信用
# tzdata: タイムゾーンデータ
RUN apk add --no-cache git ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係管理ファイルを先にコピー（Docker レイヤーキャッシュ最適化）
COPY go.mod go.sum ./

# Go モジュールの依存関係をダウンロード
# このステップはgo.mod/go.sumが変更されない限りキャッシュされる
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# Go アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o bin/gogym-api \
    ./cmd/main.go

# ==============================
# 実行ステージ
# ==============================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

RUN addgroup -g 1001 -S gogymuser && \
    adduser -S gogymuser -u 1001 -G gogymuser

ENV TZ=Asia/Tokyo
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/bin/gogym-api ./gogym-api

RUN chown -R gogymuser:gogymuser /app

USER gogymuser

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/healthz || exit 1

CMD ["./gogym-api"]