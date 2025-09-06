# GoGym シーダー

実際のジムデータを使ってMVP開発を行うためのシーダーシステムです。

## 📁 ディレクトリ構造

```
infra/seeds/
├── data/           # NDJSON シードデータ
│   ├── gyms.ndjson
│   ├── users.ndjson
│   └── reviews.ndjson
├── seeders/        # Go シーダー実装
│   ├── main.go
│   ├── user_seeder.go
│   ├── gym_seeder.go
│   └── review_seeder.go
└── scraping/       # スクレイピングスクリプト（今後追加）
```

## 🚀 使用方法

### 1. データベース準備
```bash
# マイグレーション実行
cd apps/api
goose -dir infra/migrations mysql "user:password@tcp(localhost:3307)/gogym?parseTime=true" up
```

### 2. シーダー実行
```bash
# シーダー実行（apps/apiディレクトリから）
go run infra/seeds/seeders/*.go
```

## 📊 データ形式

### gyms.ndjson
実際のジム情報を以下の形式で格納：
- 基本情報（名前、説明、住所、位置情報）
- 営業時間（曜日別の詳細時間）
- 料金プラン（月額、年額、一日券等）
- 設備・アメニティ情報
- SEO情報（スラッグ、メタ情報）

### users.ndjson
テスト用ユーザーアカウント：
- 基本プロフィール情報
- 居住地域、年齢層
- デフォルトパスワード: `password123`

### reviews.ndjson
リアルなレビューデータ：
- 5段階評価（総合＋詳細項目別）
- 訪問日、目的
- コメント内容

## 🔧 スクレイピング対応

### 想定するデータソース
- 大手ジムチェーンの公式サイト
- ジム検索ポータルサイト
- Googleマイビジネス情報

### スクレイピング時の注意点
1. **robots.txt遵守**
2. **リクエスト頻度制限** (1秒間隔推奨)
3. **利用規約確認**
4. **個人情報の除外**

### 推奨ツール
- **Python**: BeautifulSoup + requests
- **Go**: colly フレームワーク
- **Node.js**: puppeteer (SPA対応)

## 🛡️ データクリーニング

シーダー実行前に以下をチェック：
- 重複データの除去
- 不正な位置情報の修正
- 料金情報の正規化
- 営業時間の妥当性チェック

## 📈 MVP開発での活用

1. **検索機能テスト**: 位置情報・料金での絞り込み
2. **レビューシステム**: 評価表示・ソート機能
3. **UI/UX検証**: 実データでのユーザビリティ確認
4. **パフォーマンス**: 大量データでの動作検証

## 🔄 データ更新

```bash
# データ再投入（既存データは削除される）
go run infra/seeds/seeders/*.go

# 特定のシーダーのみ実行したい場合は、main.goを編集
```

## 📝 カスタマイズ

新しいデータソースを追加する場合：
1. `data/`に対応するNDJSONファイルを作成
2. `seeders/`に専用シーダーを実装
3. `main.go`でシーダーを登録