FactoryBot.define do
  factory :gym do
    name { Faker::Company.name }
    membership_fee { Faker::Commerce.price }
    business_hours { '9:00 - 21:00' }
    access { Faker::Address.full_address }
    remarks { Faker::Lorem.sentence }
    website { Faker::Internet.url }
    after(:create) do |gym|
      create(:location, gym: gym)
    end
  end
end
