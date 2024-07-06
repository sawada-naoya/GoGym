require 'json'
require 'net/http'
require 'uri'

Faker::Config.locale = 'ja'

file_path = Rails.root.join('db', 'gyms.json')
gyms = JSON.parse(File.read(file_path))

def get_coordinates(address)
  base_url = "https://maps.googleapis.com/maps/api/geocode/json"
  api_key = ENV['GOOGLE_MAPS_API_KEY']
  url = "#{base_url}?address=#{URI.encode_www_form_component(address)}&key=#{api_key}"

  response = Net::HTTP.get(URI(url))
  result = JSON.parse(response)

  if result['status'] == 'OK'
    location = result['results'][0]['geometry']['location']
    return [location['lat'], location['lng']]
  else
    raise "Geocoding API error: #{result['status']}"
  end
end

tags = [
  '駐車場',
  '駐輪場',
  'パーソナルトレーニング',
  'タンニング',
  '体組成計',
  'サウナ',
  'プール',
  '酸素カプセル',
  '24時間営業',
  'お風呂',
  'スタジオ',
  'プロテインラウンジ',
  'ウォーターサーバー'
]

tags.each do |tag_name|
  Tag.find_or_create_by!(name: tag_name)
end

# 本番環境でのみジムデータを作成
if Rails.env.production?
  # JSONファイルから読み込んだジムデータをデータベースに保存
  gyms.each do |gym|
    # 既存のジムデータを確認
    existing_gym = Gym.joins(:location).find_by(name: gym['name'], locations: { address: gym['address'] })
    next if existing_gym

    lat, lng = get_coordinates(gym['address'])

    created_gym = Gym.create!(
      name: gym['name'],
      access: gym['access'] || 'アクセス情報がありません',
      membership_fee: gym['membership_fee'] || '料金情報がありません',
      business_hours: gym['business_hours'] || '営業時間情報がありません',
      website: gym['website'] || ' '
    )

    # 各ジムに対応するロケーションデータを作成
    Location.create!(
      address: gym['address'], # JSONファイルの住所を使用
      latitude: lat, # 住所から取得した緯度
      longitude: lng, # 住所から取得した経度
      gym_id: created_gym.id # gyms と locations を関連付けるための gym_id
    )
  end
else
  # 開発環境でのみダミーデータを作成
  # ユーザーを10人作成
  10.times do |n|
    User.create!(
      name: Faker::JapaneseMedia::StudioGhibli.character,
      email: Faker::Internet.unique.email,
      password: 'test',
      password_confirmation: 'test'
    )
  end

  # ユーザーのIDを取得
  user_ids = User.pluck(:id)

  # タグの作成
  tags = [
    '駐車場', '駐輪場', 'パーソナルトレーニング', 'タンニング',
    '体組成計', 'サウナ', 'プール', '酸素カプセル',
    '24時間営業', 'お風呂', 'スタジオ', 'プロテインラウンジ',
    'ウォーターサーバー'
  ]

  tag_records = tags.map { |tag_name| Tag.find_or_create_by!(name: tag_name) }

  # JSONファイルから読み込んだジムデータをデータベースに保存
  gyms.each do |gym|
    lat, lng = get_coordinates(gym['address'])

    created_gym = Gym.create!(
      name: gym['name'],
      access: gym['access'] || 'アクセス情報がありません',
      membership_fee: Faker::Commerce.price(range: 1000..10000),
      business_hours: '9:00 - 21:00',
      website: Faker::Internet.url,
      user_id: user_ids.sample
    )

    # 各ジムに対応するロケーションデータを作成
    Location.create!(
      address: gym['address'], # JSONファイルの住所を使用
      latitude: lat, # 住所から取得した緯度
      longitude: lng, # 住所から取得した経度
      gym_id: created_gym.id # gyms と locations を関連付けるための gym_id
    )

    # 各ジムに対して20個ずつ口コミと評価を作成
    20.times do
      Review.create!(
        title: Faker::JapaneseMedia::StudioGhibli.movie,
        content: Faker::JapaneseMedia::StudioGhibli.quote,
        rating: (rand(9) * 0.5 + 1.0).round(1), # 1.0から5.0のランダムな評価を追加
        image: 'app/assets/images/fake.jpg',
        user_id: user_ids.sample,
        gym_id: created_gym.id # 正しい gym_id を使用
      )
    end

    # 各ジムにランダムに4つのタグを関連付け
    created_gym.tags << tag_records.sample(4)
  end
end
