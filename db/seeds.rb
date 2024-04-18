10.times do
  Gym.create!(
    name: Faker::Company.name,
    membership_fee: Faker::Commerce.price(range: 1000..10000),
    business_hours: Faker::Time.between(from: DateTime.now - 1, to: DateTime.now, format: :default),
    access: Faker::Address.street_address,
    remarks: Faker::Lorem.sentence,
    photos: Faker::LoremFlickr.image(size: "400x300"),
    website: Faker::Internet.url
  )
end
