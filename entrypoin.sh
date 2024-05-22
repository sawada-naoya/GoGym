#!/bin/bash
set -e

# Remove a potentially pre-existing server.pid for Rails.
rm -f /GoGym/tmp/pids/server.pid

# データベースのマイグレーションを実行
bundle exec rails db:migrate

# シードデータを読み込む
bundle exec rails db:seed

# アセットのプリコンパイルを実行
bundle exec rails assets:precompile

# Then exec the container's main process (what's set as CMD in the Dockerfile).
exec "$@"

chmod +x entrypoint.sh
