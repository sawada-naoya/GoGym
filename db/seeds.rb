Faker::Config.locale = 'ja'

10.times do
  Gym.create!(
    name: Faker::Games::Pokemon.name,
    membership_fee: Faker::Commerce.price(range: 1000..10000),
    business_hours: '9:00 - 21:00',
    access: '東京駅から徒歩5分',
    photos: 'app/assets/images/fake.jpg',
    website: Faker::Internet.url,
    location_attributes: {
      address: Faker::Address.full_address,
      latitude: Faker::Address.latitude,
      longitude: Faker::Address.longitude
    }
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
