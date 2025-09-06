# Go API サーバー用 Dockerfile
# マルチステージビルドで本番用の軽量・セキュアなイメージを作成

# ==============================
# ビルドステージ
# ==============================
FROM golang:1.25-alpine AS builder

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
# CGO_ENABLED=0: 静的リンクでビルド（Alpine Linux で動作可能）
# GOOS=linux: Linux用バイナリを生成
# GOARCH=amd64: 64bit アーキテクチャ用
# -ldflags: リンカーフラグでバイナリサイズを最適化
#   -w: デバッグ情報を削除
#   -s: シンボルテーブルを削除
#   -extldflags "-static": 静的リンク
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o bin/gogym-api \
    cmd/api/main.go

# ==============================
# 実行ステージ
# ==============================
FROM alpine:3.19

# 実行時に必要な最小限のパッケージをインストール
# ca-certificates: HTTPS通信用
# tzdata: タイムゾーン設定用
# curl: ヘルスチェック用
RUN apk add --no-cache ca-certificates tzdata curl

# セキュリティ向上: 非rootユーザーの作成
# 1001番のUID/GIDでgogymuserユーザーとgogymuserグループを作成
RUN addgroup -g 1001 -S gogymuser && \
    adduser -S gogymuser -u 1001 -G gogymuser

# タイムゾーンを日本時間（Asia/Tokyo）に設定
ENV TZ=Asia/Tokyo
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

# アプリケーション用作業ディレクトリを設定
WORKDIR /app

# ビルドステージからビルド済みバイナリをコピー
COPY --from=builder /app/bin/gogym-api ./gogym-api

# ファイル所有権を非rootユーザーに変更
RUN chown -R gogymuser:gogymuser /app

# 非rootユーザーに切り替え（セキュリティ強化）
USER gogymuser

# アプリケーションが使用するポートを明示的に公開
EXPOSE 8080

# ヘルスチェック設定
# 30秒間隔でアプリケーションの生存確認を実行
# 起動から5秒後にチェック開始、3秒でタイムアウト、3回連続失敗で異常と判定
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

# アプリケーションを実行
# PID 1として直接実行（シグナル処理のため）
CMD ["./gogym-api"]