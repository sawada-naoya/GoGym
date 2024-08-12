if Rails.env.production?
  Rails.logger.info "AWS_ACCESS_KEY_ID: #{ENV['AWS_ACCESS_KEY_ID']}"
  Rails.logger.info "AWS_SECRET_ACCESS_KEY: #{ENV['AWS_SECRET_ACCESS_KEY']}"
  Rails.logger.info "AWS_BUCKET_NAME: #{ENV['AWS_BUCKET_NAME']}"
end
