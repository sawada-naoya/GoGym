require "active_support/core_ext/integer/time"

Rails.application.configure do
  config.cache_classes = true
  config.eager_load = true

  config.consider_all_requests_local       = false
  config.action_controller.perform_caching = true

  config.public_file_server.enabled = ENV["RAILS_SERVE_STATIC_FILES"].present?

  config.assets.compile = true
  config.assets.digest = true
  config.assets.precompile += %w( application.css )
  config.assets.prefix = "/assets"

  config.active_storage.service = :amazon

  config.log_level = :info

  config.log_tags = [ :request_id ]

  config.action_mailer.perform_caching = false

  config.i18n.fallbacks = true

  config.active_support.report_deprecations = false

  config.log_formatter = ::Logger::Formatter.new

  # ログの出力先を標準出力とファイルの両方に設定する
  if ENV["RAILS_LOG_TO_STDOUT"].present?
    logger           = ActiveSupport::Logger.new(STDOUT)
    logger.formatter = config.log_formatter
    config.logger    = ActiveSupport::TaggedLogging.new(logger)
  end

  # デプロイするサービスのホストを追加する(render)
  config.hosts << 'gogym-m40u.onrender.com'
  # メール設定
  config.action_mailer.default_url_options = { host: 'gogym-m40u.onrender.com' }
  config.action_mailer.delivery_method = :smtp
  config.action_mailer.smtp_settings = {
    address: 'smtp.gmail.com',
    port: 587,
    domain: 'gmail.com',
    user_name: ENV['GMAIL_USERNAME'],
    password: ENV['GMAIL_PASSWORD'],
    authentication: :plain,
    enable_starttls_auto: true
  }
  config.action_mailer.raise_delivery_errors = true

  # Do not dump schema after migrations.
  config.active_record.dump_schema_after_migration = false

  # config.sorcery.google_callback_url = 'https://gogym-m40u.onrender.com//oauth/callback?provider=google'
end
