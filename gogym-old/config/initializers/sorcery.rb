Rails.application.config.sorcery.submodules = [:external]

Rails.application.config.sorcery.configure do |config|
  config.not_authenticated_action = :not_authenticated
  config.external_providers = %i[google]

  config.google.key = ENV['GOOGLE_KEY']
  config.google.secret = ENV['GOOGLE_SECRET']
  config.google.callback_url = "http://localhost:3000/oauth/callback?provider=google"
  config.google.user_info_mapping = {:email => "email", :name => "name"}

  config.user_config do |user|
    user.authentications_class = Authentication
    user.stretches = 1 if Rails.env.test?

  end
  config.user_class = "User"
end
