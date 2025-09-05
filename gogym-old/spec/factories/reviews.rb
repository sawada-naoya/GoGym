FactoryBot.define do
  factory :review do
    title { "Test Title" }
    content { "Test Content" }
    rating { 4.5 }
    image { Rack::Test::UploadedFile.new(Rails.root.join('spec/fixtures/test_image.jpg'), 'image/jpeg') }
    association :gym
    association :user
  end
end
