FactoryBot.define do
  factory :location do
    address { Faker::Address.full_address }
    latitude { Faker::Address.latitude }
    longitude { Faker::Address.longitude }
    association :gym
  end
end
