FactoryBot.define do
  factory :authentication do
    user { nil }
    provider { "MyString" }
    uid { "MyString" }
  end
end
