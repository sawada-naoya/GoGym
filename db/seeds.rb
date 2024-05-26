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
      membership_fee: gym['membership_fee'] || ' ',
      business_hours: gym['business_hours'] || ' ',
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

  # JSONファイルから読み込んだジムデータをデータベースに保存
  gyms.each do |gym|
    lat, lng = get_coordinates(gym['address'])

    created_gym = Gym.create!(
      name: gym['name'],
      access: gym['access'] || 'アクセス情報がありません',
      membership_fee: Faker::Commerce.price(range: 1000..10000),
      business_hours: '9:00 - 21:00',
      website: Faker::Internet.url
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
        rating: rand(1.0..5.0).round(1), # 1.0から5.0のランダムな評価を追加
        image: 'app/assets/images/fake.jpg',
        user_id: user_ids.sample,
        gym_id: created_gym.id # 正しい gym_id を使用
      )
    end
  end
end
