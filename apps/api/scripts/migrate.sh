#!/bin/bash
set -e

# マイグレーション実行スクリプト
# Render/Fly.ioのデプロイ前コマンドとして実行

# golang-migrateをインストール（まだ入っていない場合）
if ! command -v migrate &> /dev/null; then
    echo "Installing golang-migrate..."
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
    mv migrate /usr/local/bin/migrate || mv migrate ./migrate
    MIGRATE_CMD="./migrate"
else
    MIGRATE_CMD="migrate"
fi

# 環境変数チェック
if [ -z "$DB_HOST" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo "Error: Database environment variables not set"
    exit 1
fi

# DB_PORTが設定されていない場合は3306（MySQL）を使用
DB_PORT=${DB_PORT:-3306}

# DSNを構築
DSN="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?multiStatements=true"

echo "Running database migrations..."
$MIGRATE_CMD -path ./internal/infra/db/migrations -database "$DSN" up

echo "Migrations completed successfully"
