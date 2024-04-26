require_relative "boot"

require "rails/all"

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)

module GoGym
  class Application < Rails::Application
    config.load_defaults 7.0
    config.i18n.default_locale = :ja
    config.time_zone = 'Asia/Tokyo'
    config.active_record.default_timezone = :local
    config.i18n.available_locales = [:ja, :en]
    config.i18n.load_path += Dir[Rails.root.join('config', 'locales', 'view', '*.{rb, yml}')]
    config.decorator_class = 'Draper::Decorator'
    config.autoload_paths += %W(#{config.root}/app/decorators)
  end
end
