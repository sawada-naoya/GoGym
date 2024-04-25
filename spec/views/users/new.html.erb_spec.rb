require 'rails_helper'

RSpec.describe "新規登録", type: :view do
  context '入力情報正常系' do
    it 'ユーザーが新規作成できること' do
      visit new_user_path
      expect {
        fill_in 'user[name]', with: '田中 太郎'
        fill_in 'user[email]', with: 'example@example.com'
        fill_in 'user[password]', with: '12345678'
        fill_in 'user[password_confirmation]', with: '12345678'
        click_button 'commit'
        Capybara.assert_current_path("/", ignore_query: true)
      }.to change { User.count }.by(1)
      expect(page).to have_content('会員登録が完了しました'), 'フラッシュメッセージ「会員登録が完了しました」が表示されていません'
    end
  end

  context '入力情報異常系' do
    it 'ユーザーが新規作成できない' do
      visit new_user_path
      expect {
        fill_in 'user[email]', with: 'example@example.com'
        click_button 'commit'
      }.to change { User.count }.by(0)
      expect(page).to have_content('会員登録に失敗しました'), 'フラッシュメッセージ「会員登録に失敗しました」が表示されていません'
      expect(page).to have_content('アカウント名を入力してください'), 'エラーメッセージ「アカウント名を入力してください」が表示されていません'
      expect(page).to have_content('パスワードを入力してください'), 'エラーメッセージ「パスワードを入力してください」が表示されていません'
      expect(page).to have_content('パスワードは3文字以上で入力してください'), 'エラーメッセージ「パスワードは3文字以上で入力してください」が表示されていません'
    end
  end
end
