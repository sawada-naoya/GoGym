module SystemHelper
  def login_as(user)
    visit root_path
    click_link "ログイン"
    fill_in 'email', with: user.email
    fill_in 'password', with: 'password'
    click_button 'ログイン'
    # Capybara.assert_current_path(root_path, ignore_query: true)
  end
end

RSpec.configure do |config|
  config.include SystemHelper
end
