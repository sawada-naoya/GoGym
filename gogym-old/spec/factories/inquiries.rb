FactoryBot.define do
  factory :inquiry do
    name { Faker::Name.name }
    email { Faker::Internet.email }
    inquiry_content { Faker::Lorem.paragraph(sentence_count: 50) }
  end
end
