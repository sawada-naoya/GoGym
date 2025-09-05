require 'pycall'
require 'pandas'
# include PyCall::Import

sur = PyCall.import_module('surprise')
pickle = PyCall.import_module('pickle')

module Recommend
    # CSVファイルからデータを読み込み、Surpriseライブラリのデータセット形式に変換するメソッド
    def self.load_data
        # pandasを使ってCSVファイルをデータフレームに読み込み
        df = Pandas.read_csv('lib/python/reviews.csv')
        # SurpriseのReaderオブジェクトを初期化。rating_scaleはデータの評価スケールを指定。
        reader = sur.Reader(rating_scale: [df['rating'].min, df['rating'].max])
        # データフレームからSurpriseのデータセットを作成
        data = sur.Dataset.load_from_df(df[['user_id', 'gym_id', 'rating']], reader)
        data
    end

    # SVDアルゴリズムを使ってモデルをトレーニングし、モデルを保存するメソッド
    def self.train_model
        data = load_data
        # データセットをトレーニングセットに変換
        trainset = data.build_full_trainset

        # SVDアルゴリズムでモデルをトレーニング
        algo = sur.SVD.new
        # トレーニングセットを使ってモデルをトレーニング
        algo.fit(trainset)

        # トレーニング済みモデルの保存
        File.open('lib/python/svd_model.pkl', 'wb') do |f|
            # トレーニング済みモデルをpickleを使ってファイルに保存
            pickle.dump(algo, f)
        end
    end

    # トレーニング済みモデルを読み込み、特定のユーザーとジムに対してレコメンデーションを行うメソッド
    def self.recommend(user_id, gym_id)
        # pickleを使ってトレーニング済みモデルをファイルから読み込み
        algo = File.open('lib/python/svd_model.pkl', 'rb') do |f|
            pickle.load(f)
        end
        # 特定のユーザーとジムに対するレコメンデーションを生成
        prediction = algo.predict(user_id, gym_id)
        # 予測された評価値を返す
        prediction.est
    end

    # CSVファイルの存在を確認してトレーニングを実行するメソッド
    def self.train_if_csv_exists
        csv_path = 'lib/python/reviews.csv'
        if File.exist?(csv_path)
            puts "CSVファイルが存在します。トレーニングを開始します。"
            train_model
        else
            puts "CSVファイルが存在しません。トレーニングをスキップします。"
        end
    end
end

# 自動的に実行される部分
Recommend.train_if_csv_exists
