#!/bin/bash
set -e

echo "Running entrypoint.sh script..."

# Remove a potentially pre-existing server.pid for Rails.
rm -f /GoGym/tmp/pids/server.pid

echo "Removed server.pid if it existed."

# データベースのマイグレーションを実行
echo "Running database migrations..."
bundle exec rails db:migrate

# シードデータを読み込む
echo "Seeding the database..."
bundle exec rails db:seed

# アセットのプリコンパイルを実行
echo "Precompiling assets..."
bundle exec rails assets:precompile

# Then exec the container's main process (what's set as CMD in the Dockerfile).
echo "Starting the main process..."
exec "$@"
