#!/bin/bash
set -e

# Remove a potentially pre-existing server.pid for Rails.
rm -f /app/tmp/pids/server.pid

# データベースマイグレーションの実行
bundle exec rails db:migrate

# シードデータの作成
bundle exec rails db:seed

# アセットのプリコンパイル
bundle exec rails assets:precompile

# Railsサーバーの起動
exec "$@"
