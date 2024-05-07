Faker::Config.locale = 'ja'

20.times do
  gym = Gym.create!(
    name: Faker::Games::Pokemon.name,
    membership_fee: Faker::Commerce.price(range: 1000..10000),
    business_hours: '9:00 - 21:00',
    access: '東京駅から徒歩5分',
    photos: 'app/assets/images/fake.jpg',
    website: Faker::Internet.url,
  )
  location = Location.create!(
    address: Faker::Address.full_address, # ダミーの住所
    latitude: Faker::Address.latitude.to_f, # ダミーの緯度
    longitude: Faker::Address.longitude.to_f, # ダミーの経度
    gym_id: gym.id # gyms と locations を関連付けるための gym_id
  )
end

10.times do |n|
  User.create!(
    name: Faker::Name.unique.name,
    email: Faker::Internet.unique.email,
    password: "test",
    password_confirmation: "test"
  )
end

user_ids = User.pluck(:id)
Gym.all.each do |gym|
  20.times do
    Review.create!(
      title: Faker::Movies::HarryPotter.character,
      content: Faker::JapaneseMedia::StudioGhibli.quote,
      user_id: user_ids.sample,
      gym_id: gym.id
    )
  end
end
