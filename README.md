# GoGym

トレーニング記録管理アプリケーション

## 開発環境のセットアップ

### 前提条件

- Docker & Docker Compose

### 起動方法

```bash
# 1. リポジトリをクローン
git clone <repository-url>
cd GoGym

# 2. 環境変数ファイルをコピー
cd infra
cp .env.sample .env.local

# 3. コンテナを起動
docker-compose up -d
```

### アクセス

- **Web**: http://localhost:3003
- **API**: http://localhost:8081

### 停止方法

```bash
cd infra
docker-compose down
```

## 技術スタック

- **Frontend**: Next.js 15, React, TypeScript, Tailwind CSS
- **Backend**: Go 1.24, Echo
- **Database**: MySQL 8.0
- **Auth**: NextAuth v5

## ライセンス

MIT
