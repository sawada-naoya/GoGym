if Rails.env.production?
  puts "AWS_ACCESS_KEY_ID: #{ENV['AWS_ACCESS_KEY_ID']}"
  puts "AWS_SECRET_ACCESS_KEY: #{ENV['AWS_SECRET_ACCESS_KEY']}"
  puts "AWS_BUCKET_NAME: #{ENV['AWS_BUCKET_NAME']}"
end
