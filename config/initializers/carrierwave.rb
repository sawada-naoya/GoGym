require 'carrierwave/storage/abstract'
require 'carrierwave/storage/file'
require 'carrierwave/storage/fog'

CarrierWave.configure do |config|
  if Rails.env.production? # 本番環境の場合はS3へアップロード
    Rails.logger.info "S3_ACCESS_KEY_ID: #{ENV['S3_ACCESS_KEY_ID']}"
    Rails.logger.info "S3_SECRET_ACCESS_KEY: #{ENV['S3_SECRET_ACCESS_KEY']}"
    config.storage :fog
    config.fog_provider = 'fog/aws'
    config.fog_directory  = 'gogym-images' # 作成したバケット名を記述
    config.fog_credentials = {
      provider: 'AWS',
      aws_access_key_id: ENV['S3_ACCESS_KEY_ID'], # 環境変数
      aws_secret_access_key: ENV['S3_SECRET_ACCESS_KEY'], # 環境変数
      region: 'ap-northeast-1',   # アジアパシフィック(東京)を選択した場合
      path_style: true
    }
  else # 本番環境以外の場合はアプリケーション内にアップロード
    config.storage :file
    config.enable_processing = false if Rails.env.development?
  end
end
