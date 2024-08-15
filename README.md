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

## サービス概要
トレーニングをしている人、始めたい人、パーソナルトレーニングを受けてみたい人など、そのような人たちに最適のジムを見つけることができるサービスです。

## このサービスへの思い・作りたい理由  
トレーニングマシンの数や広さ、立地、金額、混雑具合、接客などジムによって、トレーニングのしやすさには違いがあります。私は普段同じジムに行っていますが、より快適にトレーニングをできる環境があれば、様々なジムを利用してみたいと思っています。  
そのような思いから、自分に合ったジムを見つけるサービスがあると便利だと思いこのサービスを作るに至りました。

## ユーザー層について  
#### メインターゲット：トレーニングを始めたい人  
→ 軽い運動から初めてみようと思ってもどのようなジムがあり、どこに行けば良いのかわからないと考えている人にとって有益な情報を届ける事ができると思うから
#### サブターゲット：トレーニングをしている人  
→ 住む場所や職場が変わり今まで行っていたジムに行けなくなった。今より金額の安いジムを探している。設備のいいジムに変えたい。などそのような人のためにも自分に合うジムを見つける事が出来るようになるから

## サービスの利用イメージ

1. ジムに行きたい、ジムを探している人が web ブラウザから利用
1. ジムに行き、そのジムの情報についてのレビューをする
1. 1,2 を繰り返してジムの情報が増える

## ユーザーの獲得について  
X、SEO など

## サービスの差別化ポイント・推しポイント

- 地図上からの検索機能があること
- 協調フィルタリングによるレコメンデーション機能で自分に合うジムを見つけられるようにサポートできる

## 機能一覧 
#### MVPリリース
- トップページ
- 検索
- タグ
- 詳細(閲覧・編集)
- 口コミ・写真投稿(閲覧・編集・削除)
- 会員登録・ログイン
- 位置情報

#### 本リリース
- レーティング
- レコメンド
- ソーシャルログイン
- パスワード変更
- お気に入り

## ER 図

![スクリーンショット 2024-08-15 22 00 44](https://github.com/user-attachments/assets/96906ea7-954e-48e1-9e3b-4ec94fade68d)


## 使用技術
| **カテゴリ** | **技術** |
----|---- 
| フロントエンド | JavaScript、HotWire CSS |
| サーバサイド | Ruby on Rail、Ruby |
| インフラ | render |
| データベース | PostgreSQL |
| 開発環境 | Docker |
| WebAPI | Geocoding API、Maps JavaScript API |
| その他 | VCS: GitHub、CI/CD: GitHub Actions |
