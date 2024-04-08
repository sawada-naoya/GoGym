<!-- # README

This README would normally document whatever steps are necessary to get the
application up and running.

Things you may want to cover:

* Ruby version

* System dependencies

* Configuration

* Database creation

* Database initialization

* How to run the test suite

* Services (job queues, cache servers, search engines, etc.)

* Deployment instructions

* ... -->

# GoGym

■ サービス概要
トレーニングをしている人、始めたい人、パーソナルトレーニングを受けてみたい人など、そのような人たちに最適のジムを見つけることができるサービスです。

■ このサービスへの思い・作りたい理由  
トレーニングマシンの数や広さ、立地、金額、混雑具合、接客などジムによって、トレーニングのしやすさには違いがあります。私は普段同じジムに行っていますが、より快適にトレーニングをできる環境があれば、様々なジムを利用してみたいと思っています。  
また、以前海外や国内へ旅行に行った際、宿泊先の近くにジムがあれば行きたいと考えました。しかし地域のジムを探しても情報が少ない事が多く、利用には至りませんでした。  
そのようなことから、旅行先の有無に限らず、それ以外の条件(設備、金額、広さなど)でも自分に合ったジムを見つけるサービスがあれば便利なのになと思っていました。そこで、「自身にとって最適なジムを探せるサービスを作ってみよう！」と思い、このサービスを作るに至りました。

■ ユーザー層について  
メインターゲット：トレーニングを始めたい人  
→ 軽い運動から初めてみようと思ってもどのようなジムがあり、どこに行けば良いのかわからないと考えている人にとって有益な情報を届ける事ができると思うから

サブターゲット：トレーニングをしている人  
→ 住む場所や職場が変わり今まで行っていたジムに行けなくなった。今より金額の安いジムを探している。設備のいいジムに変えたい。などそのような人のためにも自分に合うジムを見つける事が出来るようになるから

■ サービスの利用イメージ

1. ジムに行きたい、ジムを探している人が web ブラウザから利用
1. ジムに行き、そのジムの情報についてのレビューをする
1. 1,2 を繰り返してジムの情報が増える

■ ユーザーの獲得について  
X、SEO など

■ サービスの差別化ポイント・推しポイント

- 地図上からの検索機能があること
- 協調フィルタリングによるレコメンデーション機能で自分に合うジムを見つけられるようにサポートできる

■ 機能候補  
MVP リリース

- トップページ
- 検索
- タグ
- 詳細(閲覧・編集)
- 口コミ・写真投稿(閲覧・編集・削除)
- 会員登録・ログイン
- 位置情報

本リリース

- レーティング
- レコメンド
- ソーシャルログイン
- パスワード変更
- お気に入り

■ 機能の実装方針予定  
MVP リリース

- トップページ
- 検索：ransack,kaminari,stimulus-autocomplete
- タグ
- 詳細(閲覧・編集)
- 口コミ・写真投稿(閲覧・編集・削除)：Action Text,Active Storage,CarrierWave
- 会員登録・ログイン：device
- 位置情報：Google Maps Platform,Geocoder

本リリース

- レーティング：ratyrate
- レコメンド：Amazon Personalize
- ソーシャルログイン：Facebook Login API,Google Sign-In API,Twitter API
- パスワード変更
- お気に入り

バックエンド処理

- ActiveJob,Sidekiq

■ ER 図
https://i.gyazo.com/52f5cdbf4b70a102134574b023ce4001.png

■ 技術選定案
- 開発環境: Docker
- サーバサイド: Ruby on Rails 7系、Ruby、Rails 7.0.4.3
- フロントエンド: HotWire
- CSSフレームワーク: bootstrap5系、Silicon（Bootstrapテンプレート）
- WebAPI: Google MapAPI（GoogleマップのジオロケーションAPI）、Geocoder、Amazon Personalize、Facebook Login API、Google Sign-In API、Twitter API
- インフラ:
  - Webアプリケーションサーバ: Fly.io
  - ファイルサーバ: AWS S3
  - セッションサーバ: Redis（Redis by Upstash）
  - データベースサーバ: PostgreSQL（Fly Postgres）
- その他：
  - VCS: GitHub
  - CI/CD: GitHubActions
