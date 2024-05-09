FactoryBot.define do
  factory :gym do
    name { Faker::Company.name }
    membership_fee { "#{rand(1..10) * 1000}円" }
    business_hours { "10:00 - 22:00" }
    access { "#{Faker::Address.street_address}, #{Faker::Address.city}" }
    remarks { Faker::Lorem.sentence }
    website { Faker::Internet.url }
  end
end
