#!/bin/sh

# ==============================================
# MinIO 初期化スクリプト
# ==============================================
# GoGymアプリケーション用のMinIO（S3互換ストレージ）を自動設定します
# - バケットの作成
# - 公開アクセス権限の設定
# - 写真アップロード用の環境準備

# 環境変数の設定（デフォルト値付き）
MINIO_ROOT_USER=${MINIO_ROOT_USER:-minioadmin}      # MinIO管理者ユーザー名
MINIO_ROOT_PASSWORD=${MINIO_ROOT_PASSWORD:-minioadmin123}  # MinIO管理者パスワード
S3_BUCKET=${S3_BUCKET:-gogym-photos}               # 作成するバケット名

# MinIOサーバーの起動完了を待機
echo "MinIOサーバーの起動を待機中..."
until mc alias set myminio http://minio:9000 $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD; do
  echo "MinIOがまだ準備できていません。5秒後にリトライします..."
  sleep 5
done

echo "MinIOが準備完了しました。バケットの設定を開始します..."

# 写真保存用バケットを作成（存在しない場合のみ）
if mc mb myminio/$S3_BUCKET 2>/dev/null; then
  echo "バケット '$S3_BUCKET' を作成しました"
else
  echo "バケット '$S3_BUCKET' は既に存在します"
fi

# バケットに公開読み取り権限を設定（写真の表示のため）
# これにより、アップロードされた写真がWebアプリから直接アクセス可能になります
mc anonymous set download myminio/$S3_BUCKET

echo "MinIOの初期化が完了しました 🚀"
echo "- バケット名: $S3_BUCKET"
echo "- アクセス権限: 公開読み取り可能"
echo "- 管理コンソール: http://localhost:9002 (${MINIO_ROOT_USER}/${MINIO_ROOT_PASSWORD})"