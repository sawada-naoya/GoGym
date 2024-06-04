namespace :export do
  desc "レビューデータをCSVにエクスポート"
  task reviews: :environment do
    require 'csv'

    # Reviewモデルからuser_id, gym_id, ratingの3つのフィールドを選択して全てのレコードを取得
    reviews = Review.select(:user_id, :gym_id, :rating)
    # CSVデータを生成する。headers: trueにより、CSVの最初の行にヘッダーが追加される
    csv_data = CSV.generate(headers: true) do |csv|
      # CSVのヘッダー行を追加
      csv << ["user_id", "gym_id", "rating"]
      reviews.each do |review|
        # 各レビューのuser_id, gym_id, ratingをCSVの行として追加
        csv << [review.user_id, review.gym_id, review.rating]
      end
    end

    File.write('lib/python/reviews.csv', csv_data)
    puts 'レビューデータをlib/python/reviews.csvにエクスポートしました'
  end
end

# このRakeタスクは、Reviewモデルから特定のフィールドを選択し、それらのデータをCSV形式でlib/python/reviews.csvにエクスポートするものです。タスクを実行するには、以下のコマンドを実行します。「bundle exec rake export:reviews」これにより、指定されたCSVファイルが生成され、レビューのデータがエクスポートされます
