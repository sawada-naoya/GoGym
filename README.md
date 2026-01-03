# GoGym

GoGym は、トレーニング記録を管理することを目的とした Web アプリケーションです。

---

## 主な機能

- トレーニング記録管理（部位 / 種目 / 重量 / レップ）
- 月別ビューによる履歴確認
- モバイル UI
- 認証機能
- 多言語対応（日本語 / 英語）

---

## 設計上の特徴

- Go による API サーバーと Next.js Web フロントエンドの分離構成
- Clean Architecture をベースにした Go バックエンド設計
- Next.js App Router + Server Actions を利用した更新フロー
- features 単位で責務を分離したフロントエンド構成
- actions / apis 分離による BFF 肥大化の抑制

---

## 技術スタック

### Frontend

- Next.js 15 (App Router)
- React
- TypeScript
- Tailwind CSS
- NextAuth v5
- Server Actions

### Backend

- Go 1.25
- Echo (Web Framework)
- Clean Architecture
- Air（ホットリロード）

### Infrastructure

- PostgreSQL 16
- Docker / Docker Compose

---

## 開発環境のセットアップ

### 前提条件

- Docker & Docker Compose
- Node.js 20 以上
- Go 1.25
- Air（Go ホットリロード）

### 起動方法

#### 1. PostgreSQL を起動

```bash
# リポジトリをクローン
git clone <repository-url>
cd GoGym

# 環境変数ファイルをコピー
cd infra
cp .env.sample .env

# PostgreSQL コンテナを起動
docker-compose up -d postgres
```

#### 2. API サーバーを起動（別ターミナル）

```bash
cd apps/api
air
```

#### 3. Web フロントエンドを起動（別ターミナル）

```bash
cd apps/web
npm install
npm run dev
```

### アクセス

- **Web**: http://localhost:3003
- **API**: http://localhost:8081
- **PostgreSQL**: localhost:5433

### 停止方法

```bash
# PostgreSQL を停止
cd infra
docker-compose down

# API / Web は各ターミナルで Ctrl + C
```

---

## ライセンス

MIT
