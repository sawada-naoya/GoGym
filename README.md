# GoGym

トレーニング記録管理アプリケーション

## 開発環境のセットアップ

### 前提条件

- Docker & Docker Compose
- Node.js 20以上
- Go 1.24
- Air（Go ホットリロード）

### 起動方法

#### 1. MySQLを起動

```bash
# リポジトリをクローン
git clone <repository-url>
cd GoGym

# 環境変数ファイルをコピー
cd infra
cp .env.sample .env

# MySQLコンテナを起動
docker-compose up -d mysql
```

#### 2. APIサーバーを起動（別ターミナル）

```bash
cd apps/api
air
```

#### 3. Webフロントエンドを起動（別ターミナル）

```bash
cd apps/web
npm install
npm run dev
```

### アクセス

- **Web**: http://localhost:3003
- **API**: http://localhost:8081
- **MySQL**: localhost:3307

### 停止方法

```bash
# MySQLを停止
cd infra
docker-compose down

# API/Webは各ターミナルで Ctrl+C
```

## 技術スタック

- **Frontend**: Next.js 15, React, TypeScript, Tailwind CSS
- **Backend**: Go 1.24, Echo, Air
- **Database**: MySQL 8.0
- **Auth**: NextAuth v5

## ライセンス

MIT
