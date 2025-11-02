# GoGym

ジム検索・トレーニング記録管理アプリケーション

## 開発環境のセットアップ

### 前提条件

- Docker & Docker Compose
- Node.js 18以上
- Go 1.24

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd GoGym
```

### 2. 環境変数の設定

```bash
cd infra
cp .env.sample .env.local
```

必要に応じて `.env.local` の値を編集してください。

### 3. Dockerコンテナの起動

```bash
cd infra
docker-compose up -d
```

起動するサービス：
- MySQL (port: 3307)
- API (port: 8081)
- Web (port: 3003)

### 4. コンテナの状態確認

```bash
docker ps
```

### 5. コンテナの停止

```bash
cd infra
docker-compose down
```

データを保持したまま停止する場合は上記コマンドを使用してください。
データも削除する場合：

```bash
docker-compose down -v
```

## データベース操作

### MySQLへのログイン

#### 方法1: docker exec を使用（推奨）

```bash
# gogymユーザーでログイン
docker exec -it gogym-mysql mysql -u gogym -ppassword gogym

# rootユーザーでログイン
docker exec -it gogym-mysql mysql -u root -prootpassword gogym
```

#### 方法2: SQLファイルを実行

```bash
# SQLファイルを作成
cat > /tmp/query.sql <<'EOF'
SELECT * FROM workout_parts;
EOF

# SQLファイルを実行
docker exec -i gogym-mysql mysql -u gogym -ppassword gogym < /tmp/query.sql
```

### データベース認証情報

- **Database**: `gogym`
- **User**: `gogym`
- **Password**: `password`
- **Root Password**: `rootpassword`
- **Host** (外部から): `localhost:3307`
- **Host** (コンテナ内): `mysql:3306`

### マイグレーションの実行

マイグレーションファイルは `apps/api/internal/infra/db/migrations/` に配置されています。

```bash
# APIコンテナ内でマイグレーション実行
docker exec -it gogym-api ./api migrate up
```

## アプリケーションへのアクセス

- **Web Frontend**: http://localhost:3003
- **API**: http://localhost:8081

## プロジェクト構成

```
GoGym/
├── apps/
│   ├── api/          # Go APIサーバー
│   └── web/          # Next.js Webフロントエンド
├── infra/
│   ├── docker-compose.yml
│   ├── .env.sample
│   └── .env.local    # ローカル環境変数（git管理外）
└── README.md
```

## トラブルシューティング

### コンテナが起動しない場合

```bash
# ログを確認
docker-compose logs <service-name>

# 例：APIのログを確認
docker-compose logs api

# 全サービスのログを確認
docker-compose logs -f
```

### ポートが既に使用されている場合

`.env.local` のポート番号を変更してください：

```bash
MYSQL_PORT=3307
API_PORT=8081
WEB_PORT=3003
```

### データベースをリセットしたい場合

```bash
# コンテナとボリュームを削除
docker-compose down -v

# 再起動
docker-compose up -d
```

## 開発コマンド

### APIの開発

```bash
cd apps/api
go run cmd/api/main.go
```

### Webの開発

```bash
cd apps/web
npm install
npm run dev
```

## ライセンス

MIT
